package dao

import (
	"context"

	"github.com/Mestrace/orderbook/domain/model"
)

type ExchangeMetadata interface {
	Update(ctx context.Context, param *model.UpdateMetadataParam) (*model.UpdateMetadataData, error)
	QueryByName(ctx context.Context, param *model.QueryMetadataParam) (*model.QueryMetadataData, error)
}
