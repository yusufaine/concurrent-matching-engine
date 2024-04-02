#ifndef ORDERBOOK_HPP
#define ORDERBOOK_HPP

#include <map>
#include <memory>
#include <shared_mutex>
#include <utility>

#include "order.hpp"
#include "orderqueue.hpp"

/*
  Wrapper template class over `std::map` that includes a shared_mutex to allow
  for finer concurrency control. Functions are not thread-safe and require the
  caller to lock the mutex appropriately before calling.

  The template parameter `Compare` is a function that defines the comparison
  operation for the `std::map`. It defaults to `std::less<>` for a `SellBook`
  and `std::greater<>` for a `BuyBook`. The ordering of the map is used to
  define what order is at the top of the book for matching purposes.
*/
template <typename Compare> //
class Orderbook {
protected:
  std::shared_ptr<std::map<uint32_t,                    //
                           std::shared_ptr<OrderQueue>, //
                           Compare>>
      m_book;

public:
  Orderbook()
      : m_book(std::make_shared<std::map<uint32_t,                    //
                                         std::shared_ptr<OrderQueue>, //
                                         Compare>>()) {}

  std::shared_mutex m_book_mutex;
  void insert(std::shared_ptr<Order> active_order);
  bool match(std::shared_ptr<Order> active_order);
};

class BuyBook : public Orderbook<std::greater<>> {
public:
  void insert(std::shared_ptr<Order> active_buy_order);
  bool match(std::shared_ptr<Order> active_sell_order);
};

class SellBook : public Orderbook<std::less<>> {
public:
  void insert(std::shared_ptr<Order> active_sell_order);
  bool match(std::shared_ptr<Order> active_buy_order);
};

typedef std::pair<std::shared_ptr<BuyBook>, std::shared_ptr<SellBook>>
    OrderbookPair;

#endif
