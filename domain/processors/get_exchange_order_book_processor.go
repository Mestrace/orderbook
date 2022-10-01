package processors

import (
	"context"
	"math/rand"
	"sort"
	"sync"

	bizModel "github.com/Mestrace/orderbook/biz/model/tradesoft/exchange/order_book"
	"github.com/Mestrace/orderbook/domain/dao"
	"github.com/Mestrace/orderbook/domain/dto"
	"github.com/Mestrace/orderbook/domain/model"
	"github.com/bytedance/gopkg/util/logger"
	"golang.org/x/sync/errgroup"
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
		symbolListData, err := p.OrderBookDAO.GetSymbolList(ctx, &model.GetSymbolListParams{
			ExchangeName: req.GetExchangeName(),
		})
		if err != nil {
			logger.CtxErrorf(ctx, "get_symbol_list_failed|err=%+v", err)

			return nil, err
		}
		symbols = symbolListData.Symbols
	}

	rand.Shuffle(len(symbols)/4, func(i, j int) {
		symbols[i], symbols[j] = symbols[j], symbols[i]
	})

	resultMu := new(sync.Mutex)
	result := make([]*bizModel.Symbol, 0, len(symbols))

	fetchGroup, gctx := errgroup.WithContext(ctx)
	fetchGroup.SetLimit(20)

	for _, symbol := range symbols {
		actualSymbol := symbol // capture local variable

		fetchGroup.Go(func() error {
			item, _ := p.fetchSymbolPrice(gctx, &model.GetSymbolPriceParams{
				ExchangeName: req.GetExchangeName(),
				Symbol:       actualSymbol,
				OrderType:    req.GetOrderType(),
			})
			if item == nil {
				return nil
			}

			resultMu.Lock()
			result = append(result, item)
			resultMu.Unlock()

			return nil
		})
	}

	if req.GetOrderBy() == int32(bizModel.OrderBy_Symbol) {
		sort.Slice(result, func(i, j int) bool {
			return result[i].GetSymbol() < result[j].GetSymbol()
		})
	}

	resp.Symbols = result

	return resp, nil
}

func (p *GetExchangeOrderBookProcessor) fetchSymbolPrice(ctx context.Context, param *model.GetSymbolPriceParams) (
	*bizModel.Symbol, error,
) {
	symbolPriceData, err := p.OrderBookDAO.GetSymbolPrice(ctx, param)
	if err != nil {
		logger.CtxWarnf(ctx, "get_symbol_price_failed|err=%+v|symbol=%s", err, param.Symbol)

		return nil, err
	}

	return &bizModel.Symbol{
		Symbol: param.Symbol,
		Ask:    dto.ConvertSymbolStatToModel(symbolPriceData.Ask),
		Bid:    dto.ConvertSymbolStatToModel(symbolPriceData.Bid),
	}, nil
}
