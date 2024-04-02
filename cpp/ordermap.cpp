#include "ordermap.hpp"
#include <cstdint>
#include <memory>
#include <mutex>

std::shared_ptr<Order> OrderMap::find(uint32_t order_id) {
  auto it = m_order_map->find(order_id);
  if (it == m_order_map->end()) {
    return nullptr;
  }
  return it->second;
}

// this function is threadsafe
void OrderMap::insert(std::shared_ptr<Order> active_order) {
  std::lock_guard lk{m_order_map_mutex};
  m_order_map->insert({active_order->order_id, active_order});
}

void OrderMap::erase(uint32_t order_id) { m_order_map->erase(order_id); }
