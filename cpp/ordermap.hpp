#ifndef ORDERMAP_HPP
#define ORDERMAP_HPP

#include <cstdint>
#include <memory>
#include <shared_mutex>
#include <unordered_map>

#include "order.hpp"

/*
  Wrapper class over `std::unordered_map` that includes a shared_mutex to allow
  for finer concurrency control. Functions are not thread-safe and require
  the caller to lock the mutex appropriately before calling.
*/
class OrderMap {
private:
  std::shared_ptr<std::unordered_map<uint32_t, std::shared_ptr<Order>>>
      m_order_map;

public:
  std::shared_mutex m_order_map_mutex;
  OrderMap()
      : m_order_map(std::make_shared<
                    std::unordered_map<uint32_t, std::shared_ptr<Order>>>()) {}

  std::shared_ptr<Order> find(uint32_t order_id);

  void insert(std::shared_ptr<Order> active_order);
  void erase(uint32_t order_id);
};

#endif