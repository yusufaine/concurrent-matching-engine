package order

import (
	"context"

	"github.com/nus-cs3211-ay2324-s2/cs3211-assignment-2-a2_e0544264_e1324600/internal/pkg/cio"
)

type Order struct {
	Ctx         context.Context
	cancel      context.CancelFunc
	Instrument  string
	Id          uint32
	ExecutionId uint32
	Price       uint32
	Count       uint32
	RestingAt   int64
	OrderType   cio.InputType
	CancelBook  cio.InputType // used for cancel order
}

func NewFromInput(ctx context.Context, in cio.Input) Order {
	c, cancel := context.WithCancel(ctx)
	return Order{
		Ctx:        c,
		cancel:     cancel,
		Instrument: string(in.Instrument),
		Id:         in.OrderId,
		Price:      in.Price,
		Count:      in.Count,
		OrderType:  in.OrderType,
	}
}

func (o Order) AsInput() cio.Input {
	return cio.Input{
		Instrument: o.Instrument,
		OrderType:  o.OrderType,
		OrderId:    o.Id,
		Price:      o.Price,
		Count:      o.Count,
	}
}

func (o Order) Done() {
	o.cancel()
}

func (active Order) CanMatch(restingPrice uint32) bool {
	switch active.OrderType {
	case cio.INPUT_BUY:
		return active.Price >= restingPrice
	case cio.INPUT_SELL:
		return active.Price <= restingPrice
	default:
		panic("unable to match with unhandled order type")
	}
}
