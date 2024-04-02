# Concurrent Matching Engine

## Purpose

The purpose of this is to explore the different paradigms of concurrency, namely through shared resources (C++) and message passing (Go). The goal of the concurrent matching engineis to achieve as high of a level of concurrency as possible while ensuring client-level correctness which will be elaborated in later portions.

In C++, the synchronisation tools that is available as of C++20 can be used.
In Go, to get a better understanding of how message passing can be used to achieve concurrency, the `sync` package which includes `Mutex, RWMutex, atomic, WaitGroup, etc.` is not allowed to be used, the only synchronisation primitive that can be used is `chan`.
Due to these restrictions (and unfamiliarity to the language, namely C++), the code produced may not be of the most idiomatic/best quality.

As an aside, socket programming is also briefly touched to allow multiple clients to connect to the matching engine.

## Levels of Concurrency

Here, there are several levels of concurrency that can be achieved, each requiring their own way of ensuring correctness.

| Concurrency level                             | Explanation                                                                                              |
| --------------------------------------------- | -------------------------------------------------------------------------------------------------------- |
| No concurrency (mutex over the entire engine) | The engine fully serialiases the entire processing flow                                                  |
| Instrument-level                              | Orders for different instruments can execute concurrently, orders for the same instrument are serialised |
| Order-level                                   | Instrument-level + a buy and sell for the same instrument can be handled concurrently                    |

<details>
<summary><strong>[SPOILER] Achieved level of concurrency</strong></summary>
<br>

```text
C++: Instrument-level
Go:  Instrument-level
```

Order-level concurrency with correctness was difficult to achieve due the reader-writer starvation problem. The approach that was explored starved either the reader or writer of the orderbook.
</details>

## Correctness

Correctness here is defined by client/local level correctness. While multiple clients can connect to the engine and the OS scheduler may result in various interleavings, the only guarantee that needs to ensured is that the a client's order must be resolved before its subsequent orders -- Order N must be processed before Order (N+1).

## Further information

The `README.md` in the respective repositories explain the data structures used, the levels of concurrency achieved, as well as how testing was done to ensure correctness. A `grader` binary file, for Linux and ARM-based MacOS (located in [assets](./assets/)) is provided helps to simulate multiple clients and ensure that the correctness described in the previous section was achieved.

- [C++ README](./cpp/README.md)
- [Go README](./go/README.md)
