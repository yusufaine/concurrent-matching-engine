all: clean engine client

engine:
	go build -o ./engine ./cmd/engine

race:
	GORACE="halt_on_error=1" go build -race -o ./engine ./cmd/engine

client:
	go build -o ./client ./cmd/client

.PHONY: clean
clean:
	rm -f ./client ./engine
