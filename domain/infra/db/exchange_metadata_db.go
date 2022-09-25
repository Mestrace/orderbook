package db

import (
	"context"

	"github.com/Mestrace/orderbook/domain/dao"
	"github.com/Mestrace/orderbook/domain/model"
	"github.com/bytedance/gopkg/util/logger"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type exchangeMetadataDB struct {
	conn *gorm.DB
}

// NewExchangeMetadataDB creates new exchange metadata db management
func NewExchangeMetadataDB(conn *gorm.DB) dao.ExchangeMetadata {
	return &exchangeMetadataDB{
		conn: conn,
	}
}

func (db *exchangeMetadataDB) Update(ctx context.Context, param *model.UpdateMetadataParam) (*model.UpdateMetadataData, error) {
	err := db.conn.Clauses(clause.OnConflict{
		DoUpdates: clause.AssignmentColumns([]string{"metadata"}),
	}).Create(param.Metadata).Error
	if err != nil {
		logger.CtxErrorf(ctx, "update_failed|err=%+v", err)
		return nil, nil
	}

	return nil, nil
}

func (db *exchangeMetadataDB) QueryByName(ctx context.Context, param *model.QueryMetadataParam) (*model.QueryMetadataData, error) {
	result := &model.ExchangeMetadata{}

	if err := db.conn.First(
		result, "exchange = ?", param.ExchangeName,
	).Error; err != nil {
		logger.CtxErrorf(ctx, "query_by_name_failed|err=%+v", err)
		return nil, err
	}

	return &model.QueryMetadataData{
		Metadata: result,
	}, nil
}
