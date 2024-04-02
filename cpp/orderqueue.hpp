#ifndef ORDERQUEUE_HPP
#define ORDERQUEUE_HPP

#include <memory>
#include <queue>
#include <shared_mutex>

#include "order.hpp"

/*
  Wrapper class over `std::queue` that includes a shared_mutex to allow for
  finer concurrency control.
*/
class OrderQueue {
private:
  std::shared_ptr<std::queue<std::shared_ptr<Order>>> m_queue;

public:
  std::shared_mutex m_queue_mutex;
  OrderQueue()
      : m_queue(std::make_shared<std::queue<std::shared_ptr<Order>>>()) {}

  std::shared_ptr<Order> front() { return m_queue->front(); }

  void pop() { m_queue->pop(); }

  void emplace(std::shared_ptr<Order> order) { m_queue->emplace(order); }

  bool empty() { return m_queue->empty(); }
};

#endif
