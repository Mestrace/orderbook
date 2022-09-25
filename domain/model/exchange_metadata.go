package model

import (
	"time"
)

// ExchangeMetadata db type and internal model.
type ExchangeMetadata struct {
	ID        uint   `gorm:"primarykey"`
	Exchange  string `gorm:"uniqueIndex:uniq_exchange` // exchange name
	Metadata  []byte // raw json string text column
	CreatedAt time.Time
	UpdatedAt time.Time
}

// TableName implements gorm.Tabler.
func (ExchangeMetadata) TableName() string {
	return "tab_orderbook_exchange_metadata"
}

type UpdateMetadataParam struct {
	Metadata *ExchangeMetadata
}

type UpdateMetadataData struct {
}

type QueryMetadataParam struct {
	ExchangeName string
}

type QueryMetadataData struct {
	Metadata *ExchangeMetadata
}

// MetadataCSVRow csv row for metadata csv file upload.
type MetadataCSVRow struct {
	Key   string `csv:"Key"`
	Value string `csv:"Value"`
}
