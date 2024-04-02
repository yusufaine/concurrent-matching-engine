package cio

// The cgo code below interfaces with the struct in cio.h
// There should be no need to modify this file.

/*
#include <stdint.h>
#include "cio.h"
// Capitalized to export.
// Do not use lower caps.
typedef struct {
	enum input_type Type;
	uint32_t Order_id;
	uint32_t Price;
	uint32_t Count;
	char Instrument[9];
}cInput;
*/
import "C"
import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"unsafe"
)

type InputType byte

type Input struct {
	OrderType  InputType
	OrderId    uint32
	Price      uint32
	Count      uint32
	Instrument string
}

const (
	INPUT_BUY    InputType = 'B'
	INPUT_SELL   InputType = 'S'
	INPUT_CANCEL InputType = 'C'
)

func ReadInput(conn net.Conn) (in Input, err error) {
	buf := make([]byte, unsafe.Sizeof(C.cInput{}))
	_, err = conn.Read(buf)
	if err != nil {
		return
	}

	var cin C.cInput
	b := bytes.NewBuffer(buf)
	err = binary.Read(b, binary.LittleEndian, &cin)
	if err != nil {
		fmt.Printf("read err: %v\n", err)
		return
	}

	in.OrderType = InputType(cin.Type)
	in.OrderId = uint32(cin.Order_id)
	in.Price = uint32(cin.Price)
	in.Count = uint32(cin.Count)

	len := 0
	tmp := make([]byte, 9)
	for i, c := range cin.Instrument {
		tmp[i] = (byte)(c)
		if c != 0 {
			len++
		}
	}

	in.Instrument = string(tmp[:len])

	return
}

func WriteInput(conn net.Conn, in Input) error {
	// convert Input to cInput
	var cin C.cInput
	var cintype InputType
	switch in.OrderType {
	case INPUT_BUY:
		cintype = INPUT_BUY
	case INPUT_SELL:
		cintype = INPUT_SELL
	case INPUT_CANCEL:
		cintype = INPUT_CANCEL
	default:
		return fmt.Errorf("invalid order type: %v", in.OrderType)
	}

	cin.Type = uint32(cintype)
	cin.Order_id = C.uint32_t(in.OrderId)
	cin.Price = C.uint32_t(in.Price)
	cin.Count = C.uint32_t(in.Count)

	// copy instrument
	for i, c := range in.Instrument {
		cin.Instrument[i] = C.char(c)
		if c == 0 {
			break
		}
	}

	buf := make([]byte, unsafe.Sizeof(cin))
	binary.LittleEndian.PutUint32(buf[0:4], uint32(cin.Type))
	binary.LittleEndian.PutUint32(buf[4:8], uint32(cin.Order_id))
	binary.LittleEndian.PutUint32(buf[8:12], uint32(cin.Price))
	binary.LittleEndian.PutUint32(buf[12:16], uint32(cin.Count))
	copy(buf[16:25], []byte(in.Instrument))

	_, err := conn.Write(buf)
	return err
}

func OutputOrderDeleted(in Input, accepted bool, outTime int64) {
	acceptedTxt := "A"
	if !accepted {
		acceptedTxt = "R"
	}
	fmt.Printf("X %v %v %v\n",
		in.OrderId, acceptedTxt, outTime)
}

func OutputOrderAdded(in Input, outTime int64) {
	orderType := "S"
	if in.OrderType == INPUT_BUY {
		orderType = "B"
	}
	fmt.Printf("%v %v %v %v %v %v\n",
		orderType, in.OrderId, in.Instrument, in.Price, in.Count, outTime)
}

func OutputOrderExecuted(restingId, newId, execId, price, count uint32, outTime int64) {
	fmt.Printf("E %v %v %v %v %v %v\n",
		restingId, newId, execId, price, count, outTime)
}
