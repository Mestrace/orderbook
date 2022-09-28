package resources

import (
	"github.com/Mestrace/orderbook/conf"
	"go.uber.org/ratelimit"
)

var orderbookRateLimit ratelimit.Limiter

func InitRateLimit() {
	orderbookRateLimit = ratelimit.New(conf.Get().BlockchainCom.QPSLimit)
}

func GetOrderbookRateLimit() ratelimit.Limiter {
	return orderbookRateLimit
}
