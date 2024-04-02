#include <cstdint>
#include <memory>
#include <mutex>

#include "engine.hpp"
#include "orderbook.hpp"

template <typename Compare>
void Orderbook<Compare>::insert(std::shared_ptr<Order> active_order) {
  auto it = m_book->find(active_order->price);
  if (it != m_book->end()) {
    it->second->emplace(active_order);
    return;
  }

  auto new_queue = std::make_shared<OrderQueue>();
  new_queue->emplace(active_order);
  m_book->insert(std::make_pair(active_order->price, new_queue));
  return;
}

template <typename Compare>
bool Orderbook<Compare>::match(std::shared_ptr<Order> active_order) {
  while (active_order->count != 0) {
    if (m_book->empty()) {
      return false;
    }

    uint32_t best_price = m_book->begin()->first;
    std::shared_ptr<OrderQueue> best_queue = m_book->begin()->second;

    if (!active_order->can_match(best_price)) {
      return false;
    }

    if (best_queue->empty()) {
      m_book->erase(m_book->begin());
      continue;
    }

    auto resting_order = best_queue->front();

    // Lock the resting order, prevent mid-match cancellation
    std::unique_lock resting_lk{resting_order->m_order_mutex};

    if (resting_order->status == OrderStatus::CANCELLED ||
        resting_order->count == 0) {
      best_queue->pop();

      continue;
    }

    auto count_filled = std::min(active_order->count, resting_order->count);
    active_order->count -= count_filled;

    resting_order->count -= count_filled;
    resting_order->execution_id++;

    output_response_executed(resting_order, active_order, count_filled);

    // resting order drained
    if (resting_order->count == 0) {
      best_queue->pop();
    }
  }

  return true;
}

void BuyBook::insert(std::shared_ptr<Order> active_buy_order) {
  Orderbook::insert(active_buy_order);
  output_response_added(active_buy_order, false);
}

bool BuyBook::match(std::shared_ptr<Order> active_sell_order) {
  return Orderbook::match(active_sell_order);
}

void SellBook::insert(std::shared_ptr<Order> active_sell_order) {
  Orderbook::insert(active_sell_order);
  output_response_added(active_sell_order, true);
}

bool SellBook::match(std::shared_ptr<Order> active_buy_order) {
  return Orderbook::match(active_buy_order);
}
