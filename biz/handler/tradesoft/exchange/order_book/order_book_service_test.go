package order_book_test

import (
	"context"
	"fmt"
	"testing"

	orderBookHandler "github.com/Mestrace/orderbook/biz/handler/tradesoft/exchange/order_book"
	orderBookModel "github.com/Mestrace/orderbook/biz/model/tradesoft/exchange/order_book"
	"github.com/Mestrace/orderbook/biz/router"
	"github.com/Mestrace/orderbook/conf"
	"github.com/Mestrace/orderbook/domain/resources"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/test/assert"
	"github.com/cloudwego/hertz/pkg/common/ut"
)

func TestUnit_setResponse(t *testing.T) {
	resp := &orderBookModel.GetExchangeMetadataResp{}

	orderBookHandler.SetErrorResponse(resp, fmt.Errorf("set error"))

	assert.DeepEqual(t, resp.BizCode, int32(-1))
	assert.DeepEqual(t, resp.ErrMsg, "set error")
}

func Benchmark_GetExchangeOrderbook(b *testing.B) {
	conf.Init("config_local.json")
	resources.Init()

	r := server.Default()
	router.GeneratedRegister(r)
	r.Init()
	logger.SetLevel(logger.LevelFatal)

	b.Run("Single", func(b *testing.B) {
		clearCache()
		for i := 0; i < b.N; i++ {
			w := ut.PerformRequest(r.Engine, "GET", `exchanges/blockchain%2Ecom/order-books?order_by=1`,
				&ut.Body{}, ut.Header{Key: "Connection", Value: "close"})
			b.StopTimer()
			resp := w.Result()
			assert.DeepEqual(b, 200, resp.StatusCode())
			b.StartTimer()
		}
	})

	b.Run("Parallel-2", func(b *testing.B) {
		clearCache()
		b.SetParallelism(2)
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				ut.PerformRequest(r.Engine, "GET", `exchanges/blockchain%2Ecom/order-books?order_by=1`,
					&ut.Body{}, ut.Header{Key: "Connection", Value: "close"})
			}
		})
	})

	b.Run("Parallel-5", func(b *testing.B) {
		clearCache()
		b.SetParallelism(5)
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				ut.PerformRequest(r.Engine, "GET", `exchanges/blockchain%2Ecom/order-books?order_by=1`,
					&ut.Body{}, ut.Header{Key: "Connection", Value: "close"})
			}
		})
	})

	b.Run("Parallel-10", func(b *testing.B) {
		clearCache()
		b.SetParallelism(10)
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				ut.PerformRequest(r.Engine, "GET", `exchanges/blockchain%2Ecom/order-books?order_by=1`,
					&ut.Body{}, ut.Header{Key: "Connection", Value: "close"})
			}
		})
	})
}

func clearCache() {
	redis := resources.GetRedisClient()

	var keys []string
	var err error
	keys, _, err = redis.Scan(context.TODO(), 0, "ob:*", 200).Result()
	if err != nil {
		return
	}

	redis.Del(context.TODO(), keys...)
}
