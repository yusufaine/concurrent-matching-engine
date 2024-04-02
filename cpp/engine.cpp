#include <iostream>
#include <memory>
#include <mutex>
#include <shared_mutex>
#include <thread>
#include <utility>

#include "engine.hpp"
#include "io.hpp"

std::shared_ptr<std::atomic_uint32_t> Engine::m_clock{
    std::make_shared<std::atomic_uint32_t>(0)};

void Engine::accept(ClientConnection connection) {
  auto thread =
      std::thread(&Engine::connection_thread, this, std::move(connection));
  thread.detach();
}

void Engine::connection_thread(ClientConnection connection) {
  while (true) {
    ClientCommand input{};
    switch (connection.readInput(input)) {
    case ReadResult::Error:
      SyncCerr{} << "Error reading input" << std::endl;
    case ReadResult::EndOfFile:
      return;
    case ReadResult::Success:
      break;
    }

    switch (input.type) {
    case input_buy: {
      Engine::match_buy_order(std::make_shared<Order>(input));
      break;
    }
    case input_sell: {
      Engine::match_sell_order(std::make_shared<Order>(input));
      break;
    }
    case input_cancel: {
      auto cancel_order = std::make_shared<Order>(input);
      Engine::execute_cancel(cancel_order);
      break;
    }
    default: {
      SyncCerr{} << "Unhandled type: " << static_cast<char>(input.type)
                 << std::endl;
      break;
    }
    }
  }
}

void Engine::execute_cancel(std::shared_ptr<Order> cancel_order) {
  // Similar to CAS
  {
    std::shared_lock rlk{m_resting_map->m_order_map_mutex};
    auto order = m_resting_map->find(cancel_order->order_id);
    if (order == nullptr) {
      output_response_deleted(cancel_order, false);
      return;
    }
  }

  // Ensure that the order has not been added in another thread
  std::unique_lock wlk{m_resting_map->m_order_map_mutex};
  auto order = m_resting_map->find(cancel_order->order_id);
  if (order == nullptr) {
    output_response_deleted(cancel_order, false);
    return;
  }

  // Lock the current order, prevent it from being matched or cancelled partway
  std::unique_lock lk{order->m_order_mutex};

  // order has already been cancelled or has been filled
  if (order->status == OrderStatus::CANCELLED || order->count == 0) {
    output_response_deleted(order, false);
    return;
  }

  // order is either active or resting, remove from order_map and let the
  // orderbook handle the rest
  m_resting_map->erase(cancel_order->order_id);
  order->status = OrderStatus::CANCELLED;
  output_response_deleted(order, true);
}

void Engine::match_buy_order(std::shared_ptr<Order> active_buy_order) {
  auto books = m_instrument_map->get_orderbooks(active_buy_order->instrument);

  auto [buy_book, sell_book] = *books;

  std::scoped_lock lk{buy_book->m_book_mutex, sell_book->m_book_mutex};
  if (!sell_book->match(active_buy_order)) {
    buy_book->insert(active_buy_order);
    m_resting_map->insert(active_buy_order);
  }
}

void Engine::match_sell_order(std::shared_ptr<Order> active_sell_order) {
  auto books = m_instrument_map->get_orderbooks(active_sell_order->instrument);

  auto [buy_book, sell_book] = *books;

  std::scoped_lock lk{buy_book->m_book_mutex, sell_book->m_book_mutex};
  if (!buy_book->match(active_sell_order)) {
    sell_book->insert(active_sell_order);
    m_resting_map->insert(active_sell_order);
  }
}
