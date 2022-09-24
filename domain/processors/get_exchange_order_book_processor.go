package processors

import (
	"context"
	"sort"

	biz_model "github.com/Mestrace/orderbook/biz/model/tradesoft/exchange/order_book"
	"github.com/Mestrace/orderbook/domain/dao"
	"github.com/Mestrace/orderbook/domain/model"
	"github.com/bytedance/gopkg/util/logger"
)

type GetExchangeOrderBookProcessor struct {
	OrderBookDAO dao.ExchangeOrderBook
}

func (p *GetExchangeOrderBookProcessor) Process(ctx context.Context,
	req *biz_model.GetExchangeOrderBookReq,
) (*biz_model.GetExchangeOrderBookResp, error) {
	var (
		resp    = &biz_model.GetExchangeOrderBookResp{}
		symbols []string
	)

	if req.GetSymbol() != "" {
		symbols = []string{req.GetSymbol()}
	} else {
		symbolListData, err := p.OrderBookDAO.GetSymbolList(ctx, nil)
		if err != nil {
			logger.CtxErrorf(ctx, "get_symbol_list_failed|err=%+v", err)

			return nil, err
		}
		symbols = symbolListData.Symbols
	}

	result := make([]*biz_model.Symbol, 0, len(symbols))

	for _, symbol := range symbols {
		symbolPriceData, err := p.OrderBookDAO.GetSymbolPrice(ctx, &model.GetSymbolPriceParams{
			Symbol: symbol,
		})
		if err != nil {
			logger.CtxWarnf(ctx, "get_symbol_price_failed|err=%+v|symbol=%s", err, symbol)

			continue
		}

		item := &biz_model.Symbol{
			Symbol: symbol,
		}

		if req.GetOrderType() == int32(biz_model.OrderType_All) || req.GetOrderType() == int32(biz_model.OrderType_Bids) {
			item.Bid = &biz_model.SymbolItem{
				PxAvg:    symbolPriceData.Bid.PriceAvg.Text('f', 2),
				QtyTotal: symbolPriceData.Bid.QtyTotal.Text('f', 2),
			}
		}

		if req.GetOrderType() == int32(biz_model.OrderType_All) || req.GetOrderType() == int32(biz_model.OrderType_Asks) {
			item.Ask = &biz_model.SymbolItem{
				PxAvg:    symbolPriceData.Ask.PriceAvg.Text('f', 2),
				QtyTotal: symbolPriceData.Ask.QtyTotal.Text('f', 2),
			}
		}

		result = append(result, item)

		if req.GetOrderBy() == int32(biz_model.OrderBy_Symbol) {
			sort.Slice(result, func(i, j int) bool {
				return result[i].GetSymbol() < result[j].GetSymbol()
			})
		}
	}

	resp.Symbols = result

	return resp, nil
}
