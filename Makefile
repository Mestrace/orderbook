orderbook:
	go build -o bin/orderbook github.com/Mestrace/orderbook/cmd/orderbook

setup:
	hz update -idl idl/orderbook.thrift

run_orderbook:
	./bin/orderbook --conf config_secret.json