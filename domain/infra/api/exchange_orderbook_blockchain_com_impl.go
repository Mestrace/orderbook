package api

import (
	"context"
	"math/big"

	"github.com/Mestrace/orderbook/domain/dao"
	"github.com/Mestrace/orderbook/domain/model"
	blockchain_com "github.com/Mestrace/orderbook/third_party/lib-exchange-client/go"
	"github.com/bytedance/gopkg/util/logger"
	"go.uber.org/ratelimit"
)

type exchangeOrderBookBlockchainComImpl struct {
	client blockchain_com.APIClient
}

func NewExchangeOrderBookBlockchainCom(client blockchain_com.APIClient) dao.ExchangeOrderBook {
	return &exchangeOrderBookBlockchainComImpl{
		client: client,
	}
}

func (o *exchangeOrderBookBlockchainComImpl) GetSymbolList(
	ctx context.Context, param *model.GetSymbolListParams,
) (*model.GetSymbolListData, error) {
	response, _, err := o.client.UnauthenticatedApi.GetSymbols(ctx)
	if err != nil {
		logger.CtxErrorf(ctx, "blockchain_com_api_get_symbols_error|err=%+v", err)

		return nil, err
	}

	logger.CtxInfof(ctx, "blockchain_com_api_get_symbols_success|resp=%+v", response)

	result := make([]string, 0, len(response))

	for symbol := range response {
		result = append(result, symbol)
	}

	return &model.GetSymbolListData{
		Symbols: result,
	}, nil
}

func (o *exchangeOrderBookBlockchainComImpl) GetSymbolPrice(
	ctx context.Context, param *model.GetSymbolPriceParams,
) (*model.GetSymbolPriceData, error) {
	response, _, err := o.client.UnauthenticatedApi.GetL3OrderBook(ctx, param.Symbol)
	if err != nil {
		logger.CtxErrorf(ctx, "blockchain_com_api_get_l3_order_book_error|err=%+v", err)

		return nil, err
	}

	logger.CtxInfof(ctx, "blockchain_com_api_get_l3_order_book_success|resp=%+v", response)

	data := &model.GetSymbolPriceData{}

	if param.OrderType == model.OrderTypeAll ||
		param.OrderType == model.OrderTypeBids {
		asksAvgPrice, asksTotalQty := computeStateOfOrderBookEntry(response.Asks)

		data.Bid = &model.SymbolStat{
			PriceAvg: asksAvgPrice,
			QtyTotal: asksTotalQty,
		}
	}

	if param.OrderType == model.OrderTypeAll ||
		param.OrderType == model.OrderTypeAsks {
		bidsAvgPrice, bidsTotalQty := computeStateOfOrderBookEntry(response.Bids)
		data.Ask = &model.SymbolStat{
			PriceAvg: bidsAvgPrice,
			QtyTotal: bidsTotalQty,
		}
	}

	return data, nil
}

func computeStateOfOrderBookEntry(entrys []blockchain_com.OrderBookEntry) (*big.Float, *big.Float) {
	totalAmount := big.NewFloat(0)
	totalQty := big.NewFloat(0)
	priceAvg := new(big.Float)

	for _, asks := range entrys {
		qty := big.NewFloat(asks.Qty)
		px := big.NewFloat(asks.Px)
		amt := new(big.Float)
		amt.Mul(qty, px)

		totalQty.Add(totalQty, qty)

		totalAmount.Add(totalAmount, amt)
	}

	if totalAmount.Cmp(priceAvg) > 0 {
		priceAvg.Quo(totalAmount, totalQty)
	}

	return priceAvg, totalQty
}

// OrderbookWithRateLimit wrap order book instance with rate limit.
func OrderbookWithRateLimit(rate ratelimit.Limiter, orderbook dao.ExchangeOrderBook) dao.ExchangeOrderBook {
	return &orderbookWithRateLimit{
		rate:              rate,
		ExchangeOrderBook: orderbook,
	}
}

type orderbookWithRateLimit struct {
	rate ratelimit.Limiter
	dao.ExchangeOrderBook
}

func (o *orderbookWithRateLimit) GetSymbolList(
	ctx context.Context, param *model.GetSymbolListParams,
) (*model.GetSymbolListData, error) {
	o.rate.Take()

	return o.ExchangeOrderBook.GetSymbolList(ctx, param)
}

func (o *orderbookWithRateLimit) GetSymbolPrice(
	ctx context.Context, param *model.GetSymbolPriceParams,
) (*model.GetSymbolPriceData, error) {
	o.rate.Take()

	return o.ExchangeOrderBook.GetSymbolPrice(ctx, param)
}
