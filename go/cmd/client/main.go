package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"
	"sync/atomic"

	"github.com/nus-cs3211-ay2324-s2/cs3211-assignment-2-a2_e0544264_e1324600/internal/pkg/cio"
)

var (
	mainIsExiting atomic.Bool
	mainMutex     sync.Mutex
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <socket path>\n < <input>\n", os.Args[0])
		os.Exit(1)
	}

	conn, err := net.DialUnix(
		"unix",
		nil,
		&net.UnixAddr{Name: os.Args[1], Net: "unix"},
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error dialing: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	go pollThread(conn)

	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		line := sc.Text()

		switch cio.InputType(line[0]) {
		case '#', '\n':
			continue
		case cio.INPUT_BUY, cio.INPUT_SELL:
			cmdType := cio.INPUT_BUY
			if line[0] == byte(cio.INPUT_SELL) {
				cmdType = cio.INPUT_SELL
			}
			cmd := cio.Input{OrderType: cio.InputType(cmdType)}
			if _, err := fmt.Sscanf(line[1:], " %d %s %d %d", &cmd.OrderId, &cmd.Instrument, &cmd.Price, &cmd.Count); err != nil {
				fmt.Fprintln(os.Stderr, "Invalid new order:", line)
				os.Exit(1)
			}
			sendCommand(conn, cmd)
		case cio.INPUT_CANCEL:
			cmd := cio.Input{OrderType: cio.INPUT_CANCEL}
			if _, err := fmt.Sscanf(line[1:], " %d", &cmd.OrderId); err != nil {
				fmt.Fprintln(os.Stderr, "Invalid cancel order:", line)
				os.Exit(1)
			}
			sendCommand(conn, cmd)
		default:
			fmt.Fprintln(os.Stderr, "Invalid command:", line[0])
			os.Exit(1)
		}
	}

	mainIsExiting.Store(true)
}

func pollThread(conn net.Conn) {
	buf := make([]byte, 1024)
	for !mainIsExiting.Load() {
		_, err := conn.Read(buf)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Connection closed by server")
			os.Exit(0)
		}
	}
}

func sendCommand(conn net.Conn, cmd cio.Input) {
	mainMutex.Lock()
	defer mainMutex.Unlock()

	if err := cio.WriteInput(conn, cmd); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to write command:", err)
		os.Exit(1)
	}
}
