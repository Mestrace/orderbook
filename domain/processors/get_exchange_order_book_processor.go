package processors

import (
	"context"

	biz_model "github.com/Mestrace/orderbook/biz/model/tradesoft/exchange/order_book"
	"github.com/Mestrace/orderbook/domain/dao"
	"github.com/Mestrace/orderbook/domain/model"
	"github.com/bytedance/gopkg/util/logger"
)

type GetExchangeOrderBookProcessor struct {
	OrderBookDAO dao.ExchangeOrderBook
}

func (p *GetExchangeOrderBookProcessor) Process(ctx context.Context, req *biz_model.GetExchangeOrderBookReq) (*biz_model.GetExchangeOrderBookResp, error) {
	var (
		resp = &biz_model.GetExchangeOrderBookResp{}
		err  error
	)

	symbolListData, err := p.OrderBookDAO.GetSymbolList(ctx, nil)
	if err != nil {
		logger.CtxErrorf(ctx, "get_symbol_list_failed|err=%+v", err)
		return nil, err
	}

	result := make(map[string]*biz_model.Symbol, len(symbolListData.Symbols))
	for _, symbol := range symbolListData.Symbols {
		symbolPriceData, err := p.OrderBookDAO.GetSymbolPrice(ctx, &model.GetSymbolPriceParams{
			Symbol: symbol,
		})
		if err != nil {
			logger.CtxWarnf(ctx, "get_symbol_price_failed|err=%+v|symbol=%s", err, symbol)
			continue
		}

		result[symbol] = &biz_model.Symbol{
			Bid: &biz_model.SymbolItem{
				PxAvg:    symbolPriceData.Bid.PriceAvg.Text('f', 2),
				QtyTotal: symbolPriceData.Bid.PriceAvg.Text('f', 2),
			},
			Ask: &biz_model.SymbolItem{
				PxAvg:    symbolPriceData.Ask.PriceAvg.Text('f', 2),
				QtyTotal: symbolPriceData.Ask.PriceAvg.Text('f', 2),
			},
		}
	}

	resp.Symbols = result

	return resp, nil
}
