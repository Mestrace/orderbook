package model

import "math/big"

// GetSymbolListParams params for retrieving symbol list.
type GetSymbolListParams struct{}

// GetSymbolListData data for retrieving symbol data.
type GetSymbolListData struct {
	Symbols []string
}

// GetSymbolPriceParams params for retrieving symbol price.
type GetSymbolPriceParams struct {
	Symbol    string
	OrderType int32
}

// GetSymbolPriceData data fro retrieving symbol price.
type GetSymbolPriceData struct {
	Ask *SymbolStat
	Bid *SymbolStat
}

type SymbolStat struct {
	PriceAvg *big.Float
	QtyTotal *big.Float
}

const (
	OrderTypeAll int32 = iota
	OrderTypeAsks
	OrderTypeBids
)
