orderbook:
	go build -o bin/orderbook github.com/Mestrace/orderbook/cmd/orderbook

setup:
	thrift-fmt -d idl
	hz update -idl idl/*.thrift

run_orderbook: orderbook
	./bin/orderbook --conf config_secret.json