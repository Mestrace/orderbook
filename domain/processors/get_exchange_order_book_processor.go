package processors

import (
	"context"
	"sort"

	bizModel "github.com/Mestrace/orderbook/biz/model/tradesoft/exchange/order_book"
	"github.com/Mestrace/orderbook/domain/dao"
	"github.com/Mestrace/orderbook/domain/dto"
	"github.com/Mestrace/orderbook/domain/model"
	"github.com/bytedance/gopkg/util/logger"
)

type GetExchangeOrderBookProcessor struct {
	OrderBookDAO dao.ExchangeOrderBook
}

func (p *GetExchangeOrderBookProcessor) Process(ctx context.Context,
	req *bizModel.GetExchangeOrderBookReq,
) (*bizModel.GetExchangeOrderBookResp, error) {
	var (
		resp    = &bizModel.GetExchangeOrderBookResp{}
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

	result := make([]*bizModel.Symbol, 0, len(symbols))

	for _, symbol := range symbols {
		symbolPriceData, err := p.OrderBookDAO.GetSymbolPrice(ctx, &model.GetSymbolPriceParams{
			Symbol:    symbol,
			OrderType: req.GetOrderType(),
		})
		if err != nil {
			logger.CtxWarnf(ctx, "get_symbol_price_failed|err=%+v|symbol=%s", err, symbol)

			continue
		}

		item := &bizModel.Symbol{
			Symbol: symbol,
			Ask:    dto.ConvertSymbolStatToModel(symbolPriceData.Ask),
			Bid:    dto.ConvertSymbolStatToModel(symbolPriceData.Bid),
		}

		result = append(result, item)
	}

	if req.GetOrderBy() == int32(bizModel.OrderBy_Symbol) {
		sort.Slice(result, func(i, j int) bool {
			return result[i].GetSymbol() < result[j].GetSymbol()
		})
	}

	resp.Symbols = result

	return resp, nil
}
