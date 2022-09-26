package dao

import (
	"context"

	"github.com/Mestrace/orderbook/domain/model"
)

//go:generate mockgen -package mock_$GOPACKAGE -destination mock_dao/mock_$GOFILE -source $GOFILE

// ExchangeMetadata the data interface for exchange metadata.
type ExchangeMetadata interface {
	Update(ctx context.Context, param *model.UpdateMetadataParam) (*model.UpdateMetadataData, error)
	QueryByName(ctx context.Context, param *model.QueryMetadataParam) (*model.QueryMetadataData, error)
}
