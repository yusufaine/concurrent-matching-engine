#include "instrumentmap.hpp"
#include "orderbook.hpp"
#include <memory>
#include <mutex>
#include <shared_mutex>

std::shared_ptr<OrderbookPair>
InstrumentMap::get_orderbooks(std::string instrument) {
  // scoping is necessary here
  {
    std::shared_lock rlk{m_instrument_map_mutex};
    auto it = m_instrument_map->find(instrument);
    if (it != m_instrument_map->end()) {
      return it->second;
    }
  }

  std::unique_lock wlk{m_instrument_map_mutex};

  // ensure that it is not found
  auto it = m_instrument_map->find(instrument);
  if (it != m_instrument_map->end()) {
    return it->second;
  }

  auto p = std::make_shared<OrderbookPair>(
      std::make_pair(std::make_shared<BuyBook>(), //
                     std::make_shared<SellBook>()));
  m_instrument_map->insert({instrument, p});
  return p;
}
