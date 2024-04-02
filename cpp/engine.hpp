// This file contains declarations for the main Engine class. You will
// need to add declarations to this file as you develop your Engine.

#ifndef ENGINE_HPP
#define ENGINE_HPP

#include <atomic>
#include <cstdint>
#include <memory>
#include <string>

#include "instrumentmap.hpp"
#include "io.hpp"
#include "order.hpp"
#include "ordermap.hpp"

struct Engine {
private:
  // Logical clock to be used for Orders' tx_timestamp
  static std::shared_ptr<std::atomic_uint32_t> m_clock;

  // Lookup table containing the buy and sell books of an instrument
  std::shared_ptr<InstrumentMap> m_instrument_map;

  // Contains all the resting currently in the engine
  std::shared_ptr<OrderMap> m_resting_map;

  void connection_thread(ClientConnection conn);

  // Main engine functions

  void execute_cancel(std::shared_ptr<Order> cancel_order);
  void match_buy_order(std::shared_ptr<Order> active_buy_order);
  void match_sell_order(std::shared_ptr<Order> active_sell_order);

public:
  Engine()
      : m_instrument_map(std::make_shared<InstrumentMap>()),
        m_resting_map(std::make_shared<OrderMap>()) {}
  void accept(ClientConnection conn);

  static uint32_t fetch_add_clock() {
    return m_clock->fetch_add(1, std::memory_order_acq_rel);
  }
};

// Wrapper output functions

inline void output_response_added(std::shared_ptr<Order> active_order,
                                  bool is_sell_side) {
  Output::OrderAdded(active_order->order_id,
                     active_order->instrument.c_str(), //
                     active_order->price,              //
                     active_order->count,              //
                     is_sell_side,                     //
                     Engine::fetch_add_clock());
}

inline void output_response_executed(std::shared_ptr<Order> resting_order,
                                     std::shared_ptr<Order> new_order,
                                     uint32_t count_filled) {
  Output::OrderExecuted(resting_order->order_id,     //
                        new_order->order_id,         //
                        resting_order->execution_id, //
                        resting_order->price,        //
                        count_filled,                //
                        Engine::fetch_add_clock());
}

inline void output_response_deleted(std::shared_ptr<Order> order,
                                    bool cancel_accepted) {
  Output::OrderDeleted(order->order_id, //
                       cancel_accepted, //
                       Engine::fetch_add_clock());
}

#endif
