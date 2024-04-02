import random
import time

# Running the test script to stress test the matching engine, this does not save failed tests:
#   python3 scripts/randomtest.py | ./grader engine
# To save tests that failed, refer to `randomtest.sh`

def main():
    ticker_symbols = ["META", "AAPL", "NFLX", "GOOGL", "AMZN", "MSFT", "TSLA", "NVDA", "INTC", "AMD"]
    num_of_threads=50

    # first line
    print(str(num_of_threads) + "")
    print("0-" + str(num_of_threads - 1) + " o")
    order_id = 1
    for i in range(num_of_threads):
        random.seed(time.time())
        print(f"{i} B {order_id} {random.choice(ticker_symbols)} 2705 30")
        order_id += 1
        print(f"{i} S {order_id} {random.choice(ticker_symbols)} 3260 1")
        order_id += 1
        print(f"{i} C {order_id - 2}")
        print(f"{i} S {order_id} {random.choice(ticker_symbols)} 2701 25")
        order_id += 1
        print(f"{i} B {order_id} {random.choice(ticker_symbols)} 3290 10")
        order_id += 1
        print(f"{i} C {order_id - 3}")

    # last 2 lines
    print(".")
    print("0-" + str(num_of_threads - 1) + " x")
    
if __name__ == "__main__":
    main()
