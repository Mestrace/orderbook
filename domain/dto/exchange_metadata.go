package dto

import (
	"encoding/json"

	bizModel "github.com/Mestrace/orderbook/biz/model/tradesoft/exchange/order_book"
	"github.com/Mestrace/orderbook/domain/model"
	vd "github.com/bytedance/go-tagexpr/v2/validator"
)

// ValidateModelExchangeMetadata validates the exchange metadata.
func ValidateModelExchangeMetadata(model *bizModel.ExchangeMetadata) error {
	if err := vd.Validate(model); err != nil {
		return err
	}

	return nil
}

// ConvertModelExchangeMetadataToDB convert from api to db model.
func ConvertModelExchangeMetadataToDB(exchangeName string,
	apiModel *bizModel.ExchangeMetadata,
) (*model.ExchangeMetadata, error) {
	metadataByte, err := json.Marshal(apiModel)
	if err != nil {
		return nil, err
	}

	return &model.ExchangeMetadata{
		Exchange: exchangeName,
		Metadata: metadataByte,
	}, nil
}

// ConvertDBExchangeMetadataToModel convert from db to api.
func ConvertDBExchangeMetadataToModel(dbModel *model.ExchangeMetadata) (*bizModel.ExchangeMetadata, error) {
	bizModel := &bizModel.ExchangeMetadata{}

	err := json.Unmarshal(dbModel.Metadata, bizModel)
	if err != nil {
		return nil, err
	}

	return bizModel, nil
}

func ConvertSymbolStatToModel(stat *model.SymbolStat) *bizModel.SymbolItem {
	if stat == nil {
		return nil
	}

	return &bizModel.SymbolItem{
		PxAvg:    stat.PriceAvg.Text('f', 2),
		QtyTotal: stat.QtyTotal.Text('f', 2),
	}
}
