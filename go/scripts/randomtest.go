package main

import (
	"flag"
	"fmt"
	"math/rand"
)

// Usage: go run randomtest.go | ./grader engine

func main() {
	var clients int
	flag.IntVar(&clients, "c", 50, "number of concurrent clients")
	flag.Parse()

	tickerSymbols := []string{
		"META", "AAPL", "NFLX", "GOOGL", "AMZN", "MSFT", "TSLA", "NVDA", "INTC", "AMD",
	}
	randomSymbol := func() string { return tickerSymbols[rand.Intn(len(tickerSymbols))] }

	prices := []int{2705, 3260, 2701, 3290}
	randomPrice := func() int { return prices[rand.Intn(len(prices))] }

	count := []int{10, 25, 30, 50}
	randomCount := func() int { return count[rand.Intn(len(count))] }

	// first line
	fmt.Println(clients)
	fmt.Printf("0-%d o\n", clients-1)
	orderID := 1
	for i := 0; i < clients; i++ {
		fmt.Printf("%v B %v %v %v %v\n", i, orderID, randomSymbol(), randomPrice(), randomCount())
		orderID++
		fmt.Printf("%v S %v %v %v %v\n", i, orderID, randomSymbol(), randomPrice(), randomCount())
		orderID++
		fmt.Printf("%v C %v\n", i, orderID-2)
		fmt.Printf("%v S %v %v %v %v\n", i, orderID, randomSymbol(), randomPrice(), randomCount())
		orderID++
		fmt.Printf("%v B %v %v %v %v\n", i, orderID, randomSymbol(), randomPrice(), randomCount())
		orderID++
		fmt.Printf("%v C %v\n", i, orderID-3)
	}

	// last 2 lines
	fmt.Println(".")
	fmt.Printf("0-%d x\n", clients-1)
}
