orderbook:
	go build -o bin/orderbook github.com/Mestrace/orderbook/cmd/orderbook

setup:
	hz update -idl idl/orderbook.thrift