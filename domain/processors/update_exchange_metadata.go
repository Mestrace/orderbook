package processors

import (
	"context"
	"encoding/csv"
	"io"

	bizModel "github.com/Mestrace/orderbook/biz/model/tradesoft/exchange/order_book"
)

type UpdateExchangeMetadata struct {
}

func (p *UpdateExchangeMetadata) Process(ctx context.Context,
	req bizModel.UpdateExchangeMetadataReq, r io.Reader) (*bizModel.GetExchangeMetadataResp, error) {
	var (
		resp   = &bizModel.GetExchangeMetadataResp{}
		err    error
		reader = csv.NewReader(r)
	)

	

	return nil, nil
}
