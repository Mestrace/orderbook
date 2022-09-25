package processors

import (
	"context"
	"io"

	bizModel "github.com/Mestrace/orderbook/biz/model/tradesoft/exchange/order_book"
	"github.com/Mestrace/orderbook/common"
	"github.com/Mestrace/orderbook/domain/dao"
	"github.com/Mestrace/orderbook/domain/dto"
	"github.com/Mestrace/orderbook/domain/model"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/gocarina/gocsv"
)

type UpdateExchangeMetadata struct {
	MetadataDAO dao.ExchangeMetadata
}

func (p *UpdateExchangeMetadata) Process(ctx context.Context,
	req *bizModel.UpdateExchangeMetadataReq, r io.Reader) (*bizModel.GetExchangeMetadataResp, error) {
	var (
		resp = &bizModel.GetExchangeMetadataResp{}
		err  error
		rows = []*model.MetadataCSVRow{}
	)

	err = gocsv.Unmarshal(r, &rows)
	if err != nil {
		logger.CtxErrorf(ctx, "unmarshal_csv_failed|err=%+v", err)
		return resp, err
	}

	metadata := &bizModel.ExchangeMetadata{
		ExtInfo: make(map[string]string),
	}

	for _, row := range rows {
		ok, err := common.SetFieldByName(metadata, row.Key, row.Value)
		if err != nil {
			return resp, err
		} else if !ok {
			metadata.ExtInfo[row.Key] = row.Value
		}
	}

	err = dto.ValidateModelExchangeMetadata(metadata)
	if err != nil {
		logger.CtxErrorf(ctx, "validate_metadata_failed|err=%+v", err)
		return resp, err
	}

	dbModel, err := dto.ConvertModelExchangeMetadataToDB(req.GetExchangeName(), metadata)
	if err != nil {
		logger.CtxErrorf(ctx, "dto_convert_to_db_failed|err=%+v", err)
		return resp, err
	}

	_, err = p.MetadataDAO.Update(ctx, &model.UpdateMetadataParam{
		Metadata: dbModel,
	})
	if err != nil {
		logger.CtxErrorf(ctx, "update_metadata_db_failed|err=%+v", err)
		return resp, err
	}

	return resp, nil
}
