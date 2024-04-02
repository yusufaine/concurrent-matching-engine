package orderheap

import (
	"container/heap"

	"github.com/nus-cs3211-ay2324-s2/cs3211-assignment-2-a2_e0544264_e1324600/internal/pkg/cio"
	"github.com/nus-cs3211-ay2324-s2/cs3211-assignment-2-a2_e0544264_e1324600/internal/pkg/order"
)

type OrderHeap []*order.Order

func New(orders ...order.Order) *OrderHeap {
	pq := make(OrderHeap, len(orders))
	for i, order := range orders {
		pq[i] = &order
	}
	heap.Init(&pq)
	return &pq
}

// Assert that OrderHeap implements heap.Interface.
var _ heap.Interface = (*OrderHeap)(nil)

func (pq OrderHeap) Len() int { return len(pq) }

func (pq OrderHeap) Less(i, j int) bool {
	iPrice, jPrice := pq[i].Price, pq[j].Price
	if iPrice == jPrice {
		return pq[i].RestingAt < pq[j].RestingAt
	}

	switch pq[i].OrderType {
	case cio.INPUT_BUY:
		return iPrice > jPrice
	case cio.INPUT_SELL:
		return iPrice < jPrice
	default:
		panic("unable to compare with unhandled order type")
	}
}

func (pq OrderHeap) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }

func (pq *OrderHeap) Push(x any) { *pq = append(*pq, x.(*order.Order)) }

func (pq *OrderHeap) Pop() any {
	old := *pq
	n := len(old)
	x := old[n-1]
	old[n-1] = nil // avoid memory leak
	*pq = old[0 : n-1]
	return *x
}

func (pq OrderHeap) Peek() (*order.Order, bool) {
	if pq.Len() == 0 {
		return nil, false
	}
	return pq[0], true
}
