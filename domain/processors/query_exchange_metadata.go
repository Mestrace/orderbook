package processors

import (
	"context"

	bizModel "github.com/Mestrace/orderbook/biz/model/tradesoft/exchange/order_book"
	"github.com/Mestrace/orderbook/domain/dao"
	"github.com/Mestrace/orderbook/domain/dto"
	"github.com/Mestrace/orderbook/domain/model"
	"github.com/bytedance/gopkg/util/logger"
)

type GetExchangeMetadata struct {
	MetadataDAO dao.ExchangeMetadata
}

func (p *GetExchangeMetadata) Process(ctx context.Context, req *bizModel.GetExchangeMetadataReq) (
	*bizModel.GetExchangeMetadataResp, error,
) {
	var (
		resp = &bizModel.GetExchangeMetadataResp{}
		err  error
	)

	queryResp, err := p.MetadataDAO.QueryByName(ctx, &model.QueryMetadataParam{
		ExchangeName: req.GetExchangeName(),
	})
	if err != nil {
		logger.CtxErrorf(ctx, "query_exchange_by_name_failed|err=%+v", err)

		return resp, err
	}

	metadata, err := dto.ConvertDBExchangeMetadataToModel(queryResp.Metadata)
	if err != nil {
		logger.CtxErrorf(ctx, "conversion_to_resp_failed|err=%+v", err)

		return resp, err
	}

	resp.Metadata = metadata

	return resp, nil
}
