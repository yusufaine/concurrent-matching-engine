package instrument

import (
	"container/heap"
	"context"

	"github.com/nus-cs3211-ay2324-s2/cs3211-assignment-2-a2_e0544264_e1324600/internal/pkg/cio"
	"github.com/nus-cs3211-ay2324-s2/cs3211-assignment-2-a2_e0544264_e1324600/internal/pkg/lclock"
	"github.com/nus-cs3211-ay2324-s2/cs3211-assignment-2-a2_e0544264_e1324600/internal/pkg/order"
	"github.com/nus-cs3211-ay2324-s2/cs3211-assignment-2-a2_e0544264_e1324600/internal/pkg/orderheap"
)

type Instrument struct {
	OrderProcessingChannel chan order.Order
}

func New(ctx context.Context) Instrument {
	ins := Instrument{
		OrderProcessingChannel: make(chan order.Order),
	}

	go serveProcessingChannel(ctx, ins.OrderProcessingChannel)

	return ins
}

func serveProcessingChannel(ctx context.Context, opc chan order.Order) {
	buyHeap := orderheap.New()
	sellHeap := orderheap.New()

	for {
		select {
		case <-ctx.Done():
			close(opc)
			return
		case order := <-opc:
			switch order.OrderType {
			case cio.INPUT_BUY, cio.INPUT_SELL:
				match(buyHeap, sellHeap, order)
				order.Done()
			case cio.INPUT_CANCEL:
				var h *orderheap.OrderHeap
				switch order.CancelBook {
				case cio.INPUT_BUY:
					h = buyHeap
				case cio.INPUT_SELL:
					h = sellHeap
				default:
					panic("invalid cancel order type")
				}
				cio.OutputOrderDeleted(
					order.AsInput(),
					removeFromHeap(order.Id, h),
					lclock.GetCurrentTimestamp(),
				)
				order.Done()
			default:
				panic("invalid order type")
			}
		}
	}
}

func match(buyHeap, sellHeap *orderheap.OrderHeap, activeOrder order.Order) {
	// determine resting book
	restingBook := sellHeap
	if activeOrder.OrderType == cio.INPUT_SELL {
		restingBook = buyHeap
	}

	// resting book is empty, add active order to orderbook
	restingOrder, ok := restingBook.Peek()
	if !ok {
		addToOrderbook(buyHeap, sellHeap, activeOrder)
		return
	}

	// active order cannot match resting order, add active order to orderbook
	if !activeOrder.CanMatch(restingOrder.Price) {
		addToOrderbook(buyHeap, sellHeap, activeOrder)
		return
	}

	countFilled := min(activeOrder.Count, restingOrder.Count)

	activeOrder.Count -= countFilled
	restingOrder.Count -= countFilled
	restingOrder.ExecutionId++

	cio.OutputOrderExecuted(restingOrder.Id,
		activeOrder.Id,
		restingOrder.ExecutionId,
		restingOrder.Price,
		countFilled,
		lclock.GetCurrentTimestamp(),
	)

	if restingOrder.Count == 0 {
		heap.Pop(restingBook)
	}

	if activeOrder.Count != 0 {
		match(buyHeap, sellHeap, activeOrder)
	}
}

func addToOrderbook(buyHeap, sellHeap *orderheap.OrderHeap, order order.Order) {
	pq := sellHeap
	if order.OrderType == cio.INPUT_BUY {
		pq = buyHeap
	}
	order.RestingAt = lclock.GetCurrentTimestamp()
	heap.Push(pq, &order) // heap.Push requires a pointer
	cio.OutputOrderAdded(order.AsInput(), order.RestingAt)
}

func removeFromHeap(orderId uint32, book *orderheap.OrderHeap) bool {
	for i, order := range *book {
		if order.Id == orderId {
			heap.Remove(book, i)
			return true
		}
	}
	return false
}
