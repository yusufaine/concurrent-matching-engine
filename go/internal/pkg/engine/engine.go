package engine

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/nus-cs3211-ay2324-s2/cs3211-assignment-2-a2_e0544264_e1324600/internal/pkg/cio"
	"github.com/nus-cs3211-ay2324-s2/cs3211-assignment-2-a2_e0544264_e1324600/internal/pkg/instrument"
	"github.com/nus-cs3211-ay2324-s2/cs3211-assignment-2-a2_e0544264_e1324600/internal/pkg/lclock"
	"github.com/nus-cs3211-ay2324-s2/cs3211-assignment-2-a2_e0544264_e1324600/internal/pkg/order"
)

type lookupRequest struct {
	InstrumentName string
	OrderId        uint32
	Return         chan instrument.Instrument
}

type engine struct {
	lookupChannel chan lookupRequest
}

func New(ctx context.Context) *engine {
	e := &engine{
		lookupChannel: make(chan lookupRequest),
	}

	// spawn a goroutine to respond to instrument lookup requests
	go serveLookupChannel(ctx, e.lookupChannel)

	return e
}

// main thread is responsible for maintaining and resulving the mapping of instrument
// name to their respective structs for all clients (fan-in pattern)
func serveLookupChannel(ctx context.Context, lookupChannel chan lookupRequest) {
	instrumentMap := make(map[string]instrument.Instrument)
	for {
		select {
		case <-ctx.Done():
			close(lookupChannel)
			return
		case c := <-lookupChannel:
			ins, ok := instrumentMap[c.InstrumentName]
			if !ok {
				ins = instrument.New(ctx)
				instrumentMap[c.InstrumentName] = ins
			}
			c.Return <- ins
		}
	}
}

func (e *engine) Accept(ctx context.Context, conn net.Conn) {
	go func() {
		<-ctx.Done()
		conn.Close()
	}()

	// spawn a new goroutine for each client
	go handleConn(ctx, conn, e.lookupChannel)
}

// creates a new client and reuses the connection to handle all orders from the same client
func handleConn(ctx context.Context, conn net.Conn, lookupChannel chan<- lookupRequest) {
	type clientOrderInfo struct {
		Instrument string
		OrderType  cio.InputType
	}

	// keeps track of all buy/sell orders from this client
	clientOrderLookup := map[uint32]clientOrderInfo{}

	for {
		select {
		case <-ctx.Done():
			return
		default:
			in, err := cio.ReadInput(conn)
			if err != nil {
				if err != io.EOF {
					_, _ = fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
				}
				return
			}

			newOrder := order.NewFromInput(ctx, in)
			switch newOrder.OrderType {
			case cio.INPUT_CANCEL:
				v, ok := clientOrderLookup[newOrder.Id]
				// order to cancel does not exist
				if !ok {
					cio.OutputOrderDeleted(in, false, lclock.GetCurrentTimestamp())
					continue
				}
				// newOrder is missing the instrument name
				newOrder.Instrument = v.Instrument
				newOrder.CancelBook = v.OrderType
			case cio.INPUT_BUY, cio.INPUT_SELL:
				clientOrderLookup[newOrder.Id] = clientOrderInfo{
					Instrument: newOrder.Instrument,
					OrderType:  newOrder.OrderType,
				}
			}

			// each goroutine has their own channel to query for their instrument structs
			returnChannel := make(chan instrument.Instrument)
			lookupChannel <- lookupRequest{
				InstrumentName: newOrder.Instrument,
				OrderId:        newOrder.Id,
				Return:         returnChannel,
			}
			ins := <-returnChannel
			close(returnChannel)

			ins.OrderProcessingChannel <- newOrder

			// wait for the current order to be processed
			<-newOrder.Ctx.Done()
		}
	}
}
