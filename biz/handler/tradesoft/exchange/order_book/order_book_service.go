// Code generated by hertz generator.

package order_book

import (
	"context"
	"fmt"
	"time"

	order_book "github.com/Mestrace/orderbook/biz/model/tradesoft/exchange/order_book"
	"github.com/Mestrace/orderbook/conf"
	"github.com/Mestrace/orderbook/domain/infra/api"
	"github.com/Mestrace/orderbook/domain/infra/cache"
	"github.com/Mestrace/orderbook/domain/infra/db"
	"github.com/Mestrace/orderbook/domain/processors"
	"github.com/Mestrace/orderbook/domain/resources"
	blockchain_com "github.com/Mestrace/orderbook/third_party/lib-exchange-client/go"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/cloudwego/hertz/pkg/app"
)

// GetExchangeOrderBook .
// @router exchanges/:exchange_name/order-books [GET]
func GetExchangeOrderBook(ctx context.Context, c *app.RequestContext) {
	var err error
	var req order_book.GetExchangeOrderBookReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(400, err.Error())
		return
	}

	if req.GetExchangeName() != "blockchain.com" {
		err = fmt.Errorf("\"%s\" is not supported", req.GetExchangeName())
		handleResponse(ctx, c, &order_book.GetExchangeOrderBookResp{}, err)
		return
	}

	cfg := conf.Get()

	auth := context.WithValue(ctx, blockchain_com.ContextAPIKey, blockchain_com.APIKey{
		Key:    cfg.BlockchainCom.APIKey,
		Prefix: "Bearer", // Omit if not necessary.
	})

	orderbookDAO := api.NewExchangeOrderBookBlockchainCom(
		*blockchain_com.NewAPIClient(blockchain_com.NewConfiguration()))
	orderbookDAO = api.OrderbookWithRateLimit(resources.GetOrderbookRateLimit(), orderbookDAO)
	orderbookDAO = cache.OrderBookWithRedis(orderbookDAO, resources.GetRedisClient(),
		time.Duration(cfg.BlockchainCom.SymbolListCacheDuration),
		time.Duration(cfg.BlockchainCom.SymbolCacheDuration))

	proc := &processors.GetExchangeOrderBookProcessor{
		OrderBookDAO: orderbookDAO,
	}

	resp, err := proc.Process(auth, &req)
	handleResponse(ctx, c, resp, err)
}

// GetExchangeOrderBookAll .
// @router exchanges/:exchange_name/order-books [GET]
func GetExchangeOrderBookAll(ctx context.Context, c *app.RequestContext) {
	GetExchangeOrderBook(ctx, c)
}

// GetExchangeMetadata .
// @router exchanges/:exchange_name/metadata [GET]
func GetExchangeMetadata(ctx context.Context, c *app.RequestContext) {
	var err error
	var req order_book.GetExchangeMetadataReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(400, err.Error())
		return
	}

	metadataDB := resources.GetMasterOrderBookMainDb()
	metadataDAO := db.NewExchangeMetadataDB(metadataDB)

	proc := processors.GetExchangeMetadata{
		MetadataDAO: metadataDAO,
	}

	resp, err := proc.Process(ctx, &req)
	handleResponse(ctx, c, resp, err)
}

// UpdateExchangeMetadata .
// @router exchanges/:exchange_name/metadata [POST]
func UpdateExchangeMetadata(ctx context.Context, c *app.RequestContext) {
	var err error
	var req order_book.UpdateExchangeMetadataReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(400, err.Error())
		return
	}

	fileHead, err := c.FormFile("file")
	if err != nil {
		logger.CtxErrorf(ctx, "read_form_file_failed|err=%+v", err)
		handleResponse(ctx, c, &order_book.UpdateExchangeMetadataReq{}, err)
		return
	}

	file, err := fileHead.Open()
	if err != nil {
		logger.CtxErrorf(ctx, "open_form_file_failed|err=%+v", err)
		handleResponse(ctx, c, &order_book.UpdateExchangeMetadataReq{}, err)
		return
	}

	metadataDB := resources.GetMasterOrderBookMainDb()
	metadataDAO := db.NewExchangeMetadataDB(metadataDB)

	proc := processors.UpdateExchangeMetadata{
		MetadataDAO: metadataDAO,
	}

	resp, err := proc.Process(ctx, &req, file)
	handleResponse(ctx, c, resp, err)
}
