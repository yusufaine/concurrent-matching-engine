#ifndef INSTRUMENTMAP_HPP
#define INSTRUMENTMAP_HPP

#include <memory>
#include <shared_mutex>
#include <string>
#include <unordered_map>

#include "orderbook.hpp"

/*
  Wrapper class over `std::unordered_map` that includes a shared_mutex to allow
  for finer concurrency control. Functions are thread-safe.
*/
class InstrumentMap {
private:
  std::shared_ptr<
      std::unordered_map<std::string, std::shared_ptr<OrderbookPair>>>
      m_instrument_map;
  std::shared_mutex m_instrument_map_mutex;

public:
  InstrumentMap()
      : m_instrument_map(std::make_shared<std::unordered_map<
                             std::string, std::shared_ptr<OrderbookPair>>>()) {}

  std::shared_ptr<OrderbookPair> get_orderbooks(std::string instrument);
};

#endif
