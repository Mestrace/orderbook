package processors_test

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"testing"

	bizModel "github.com/Mestrace/orderbook/biz/model/tradesoft/exchange/order_book"
	"github.com/Mestrace/orderbook/domain/dao/mock_dao"
	"github.com/Mestrace/orderbook/domain/dto"
	"github.com/Mestrace/orderbook/domain/model"
	"github.com/Mestrace/orderbook/domain/processors"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gocarina/gocsv"
	. "github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
)

func TestUnit_UpdateExchangeMetadata_Process(t *testing.T) {
	mockCtrl := NewController(t)

	faker := gofakeit.New(0)

	Convey("should process normally", t, func() {
		var (
			NumRows        = 10
			b              bytes.Buffer
			descriptionRow = &model.MetadataCSVRow{
				Key:   "Description",
				Value: "The world’s most popular way to buy, sell, and trade crypto.",
			}
			websiteRow = &model.MetadataCSVRow{
				Key:   "WebSite",
				Value: "www.blockchain.com",
			}
			actualDBMetadata     *model.ExchangeMetadata
			expectedExchangeName = "blockchain.com"
			expectedAPIMetadata  = &bizModel.ExchangeMetadata{
				Description: descriptionRow.Value,
				WebSite:     websiteRow.Value,
				ExtInfo:     make(map[string]string),
			}
		)

		// produce data
		{
			metadataCsvRows := make([]*model.MetadataCSVRow, 0, NumRows)

			metadataCsvRows = append(metadataCsvRows, descriptionRow, websiteRow)

			for len(metadataCsvRows) < NumRows {
				row := &model.MetadataCSVRow{
					Key:   faker.Noun(),
					Value: faker.Sentence(20),
				}
				metadataCsvRows = append(metadataCsvRows, row)
				expectedAPIMetadata.ExtInfo[row.Key] = row.Value
			}

			err := gocsv.Marshal(&metadataCsvRows, bufio.NewWriter(&b))
			So(err, ShouldBeNil)

		}

		mockMetadataDAO := mock_dao.NewMockExchangeMetadata(mockCtrl)
		mockMetadataDAO.EXPECT().Update(Any(), Any()).DoAndReturn(
			func(_ interface{}, param *model.UpdateMetadataParam) (*model.UpdateMetadataData, error) {
				actualDBMetadata = param.Metadata
				return nil, nil
			})

		proc := &processors.UpdateExchangeMetadata{
			MetadataDAO: mockMetadataDAO,
		}

		_, err := proc.Process(context.Background(), &bizModel.UpdateExchangeMetadataReq{
			ExchangeName: expectedExchangeName,
		}, bufio.NewReader(&b))
		So(err, ShouldBeNil)
		So(actualDBMetadata, ShouldNotBeNil)
		So(actualDBMetadata.Exchange, ShouldEqual, expectedExchangeName)

		actualExchangeMetadata, err := dto.ConvertDBExchangeMetadataToModel(actualDBMetadata)
		So(err, ShouldBeNil)
		So(actualExchangeMetadata, ShouldResemble, expectedAPIMetadata)
	})

	Convey("should return err for empty description", t, func() {
		var (
			b              bytes.Buffer
			descriptionRow = &model.MetadataCSVRow{
				Key:   "Description",
				Value: "",
			}
			websiteRow = &model.MetadataCSVRow{
				Key:   "WebSite",
				Value: "www.blockchain.com",
			}
		)
		// produce data
		{
			metadataCsvRows := []*model.MetadataCSVRow{descriptionRow, websiteRow}
			err := gocsv.Marshal(&metadataCsvRows, bufio.NewWriter(&b))
			So(err, ShouldBeNil)
		}

		proc := &processors.UpdateExchangeMetadata{
			MetadataDAO: nil,
		}

		_, err := proc.Process(context.Background(), &bizModel.UpdateExchangeMetadataReq{}, bufio.NewReader(&b))
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldContainSubstring, "Description")
	})

	Convey("should return err for csv file without incorrect column", t, func() {
		type AnotherCsvRow struct {
			Key   string `csv:"key1"`
			Value string `csv:"key2"`
		}

		var (
			b              bytes.Buffer
			descriptionRow = &AnotherCsvRow{
				Key:   "Description",
				Value: "The world’s most popular way to buy, sell, and trade crypto.",
			}
			websiteRow = &AnotherCsvRow{
				Key:   "WebSite",
				Value: "www.blockchain.com",
			}
		)
		// produce data
		{
			csvRows := []*AnotherCsvRow{descriptionRow, websiteRow}
			err := gocsv.Marshal(&csvRows, bufio.NewWriter(&b))
			So(err, ShouldBeNil)
		}

		proc := &processors.UpdateExchangeMetadata{
			MetadataDAO: nil,
		}

		_, err := proc.Process(context.Background(), &bizModel.UpdateExchangeMetadataReq{}, bufio.NewReader(&b))
		So(err, ShouldNotBeNil)
	})

	Convey("should return parse error if the type is incorrect", t, func() {
		type AnotherCsvRow struct {
			Key   string `csv:"key1"`
			Value string `csv:"key2"`
		}

		var r io.Reader
		// produce data
		{
			imgByte := faker.ImageJpeg(100, 100)
			r = bytes.NewReader(imgByte)
		}

		proc := &processors.UpdateExchangeMetadata{
			MetadataDAO: nil,
		}

		_, err := proc.Process(context.Background(), &bizModel.UpdateExchangeMetadataReq{}, r)
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldContainSubstring, "parse error")
	})
}
