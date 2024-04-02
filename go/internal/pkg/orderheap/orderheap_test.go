package orderheap

import (
	"container/heap"
	"math/rand"
	"testing"

	"github.com/nus-cs3211-ay2324-s2/cs3211-assignment-2-a2_e0544264_e1324600/internal/pkg/cio"
	"github.com/nus-cs3211-ay2324-s2/cs3211-assignment-2-a2_e0544264_e1324600/internal/pkg/order"
)

// Head of the buy heap should be the highest price.
func Test_BuyHeap_Price(t *testing.T) {
	buyOrder := []order.Order{
		{Price: 300, OrderType: cio.INPUT_BUY},
		{Price: 100, OrderType: cio.INPUT_BUY},
		{Price: 200, OrderType: cio.INPUT_BUY},
	}

	buyHeap := New(buyOrder...)
	head := heap.Pop(buyHeap).(order.Order)
	if head.Price != 300 {
		t.Errorf("BuyHeap() failed, expected 100 but got %d", head.Price)
	}

	head = heap.Pop(buyHeap).(order.Order)
	if head.Price != 200 {
		t.Errorf("BuyHeap() failed, expected 200 but got %d", head.Price)
	}

	head = heap.Pop(buyHeap).(order.Order)
	if head.Price != 100 {
		t.Errorf("BuyHeap() failed, expected 300 but got %d", head.Price)
	}
}

// Head of the sell heap should be the lowest price.
func Test_SellHeap_Price(t *testing.T) {
	sellOrder := []order.Order{
		{Price: 300, OrderType: cio.INPUT_SELL},
		{Price: 100, OrderType: cio.INPUT_SELL},
		{Price: 200, OrderType: cio.INPUT_SELL},
	}

	sellHeap := New(sellOrder...)
	head := heap.Pop(sellHeap).(order.Order)
	if head.Price != 100 {
		t.Errorf("BuyHeap() failed, expected 100 but got %d", head.Price)
	}

	head = heap.Pop(sellHeap).(order.Order)
	if head.Price != 200 {
		t.Errorf("BuyHeap() failed, expected 200 but got %d", head.Price)
	}

	head = heap.Pop(sellHeap).(order.Order)
	if head.Price != 300 {
		t.Errorf("BuyHeap() failed, expected 300 but got %d", head.Price)
	}
}

// Head of the heap should be the earliest time.
func Test_BuyHeap_Time(t *testing.T) {
	// Test the price of a buy order.
	orders := []order.Order{
		{RestingAt: 3, OrderType: cio.INPUT_BUY},
		{RestingAt: 1, OrderType: cio.INPUT_BUY},
		{RestingAt: 2, OrderType: cio.INPUT_BUY},
	}

	for range 30 {
		// Randomise the order of the orders.
		rand.Shuffle(len(orders), func(i, j int) { orders[i], orders[j] = orders[j], orders[i] })

		buyHeap := New(orders...)
		head := heap.Pop(buyHeap).(order.Order)
		if head.RestingAt != 1 {
			t.Errorf("BuyHeap() failed, expected 100 but got %d", head.Price)
		}

		head = heap.Pop(buyHeap).(order.Order)
		if head.RestingAt != 2 {
			t.Errorf("BuyHeap() failed, expected 200 but got %d", head.Price)
		}

		head = heap.Pop(buyHeap).(order.Order)
		if head.RestingAt != 3 {
			t.Errorf("BuyHeap() failed, expected 300 but got %d", head.Price)
		}
	}
}
