package dto

import (
	"testing"

	bizModel "github.com/Mestrace/orderbook/biz/model/tradesoft/exchange/order_book"
)

func TestUnit_ValidateModelExchangeMetadata(t *testing.T) {
	type args struct {
		model *bizModel.ExchangeMetadata
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				model: &bizModel.ExchangeMetadata{
					Description: "a blockchain exchange",
					WebSite:     "www.blockchain.com",
					ExtInfo: map[string]string{
						"number_of_people": "100",
					},
				},
			},
		},
		{
			name: "desc empty",
			args: args{
				model: &bizModel.ExchangeMetadata{
					WebSite: "www.blockchain.com",
					ExtInfo: map[string]string{
						"number_of_people": "100",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "website empty",
			args: args{
				model: &bizModel.ExchangeMetadata{
					Description: "a blockchain exchange",
					ExtInfo: map[string]string{
						"number_of_people": "100",
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateModelExchangeMetadata(tt.args.model); (err != nil) != tt.wantErr {
				t.Errorf("ValidateModelExchangeMetadata() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
