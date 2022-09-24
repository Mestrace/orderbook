package dao

import (
	"context"

	"github.com/Mestrace/orderbook/domain/model"
)

type ExchangeOrderBook interface {
	GetSymbolList(ctx context.Context, param *model.GetSymbolListParams) (*model.GetSymbolListData, error)
	GetSymbolPrice(ctx context.Context, param *model.GetSymbolPriceParams) (*model.GetSymbolPriceData, error)
}
