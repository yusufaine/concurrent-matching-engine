#ifndef ORDER_HPP
#define ORDER_HPP

#include <cstdint>
#include <mutex>
#include <stdexcept>
#include <string>

#include "io.hpp"

enum class OrderStatus { VALID, CANCELLED };

/*
  Class representing an order in the system. It includes a mutex to ensure that
  the order is not modified by multiple threads at the same time, (e.g in the
  case of a cancel)
*/
class Order {
public:
  std::string instrument;
  const uint32_t order_id;
  uint32_t price;
  uint32_t count;
  OrderStatus status;
  CommandType cmd_type;
  uint32_t execution_id;
  std::mutex m_order_mutex;

  Order(ClientCommand cc)
      : instrument(cc.instrument), order_id(cc.order_id), price(cc.price),
        count(cc.count), status(OrderStatus::VALID), cmd_type(cc.type),
        execution_id(0) {}

  // returns true if the order can match the resting price based on its type
  bool can_match(uint32_t resting_price) {
    switch (cmd_type) {
    // selling is cheaper
    case CommandType::input_buy:
      return price >= resting_price;
    // buying is more expensive
    case CommandType::input_sell:
      return price <= resting_price;
    default:
      throw std::invalid_argument("unhandled command type");
    }
  }
};

#endif
