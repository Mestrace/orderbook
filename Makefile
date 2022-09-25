tidy:
	go mod tidy
	
orderbook: tidy
	go build -o bin/orderbook github.com/Mestrace/orderbook/cmd/orderbook

setup:
	thrift-fmt -d idl
	hz update -idl idl/*.thrift

run_orderbook: orderbook
	./bin/orderbook --conf config_local.json

test.unittest:
	go test -v --run Unit ./...

test.integration:
	go test -v --run Integration ./...