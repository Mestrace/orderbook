package db

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/Mestrace/orderbook/conf"
	"github.com/Mestrace/orderbook/domain/model"
	"github.com/Mestrace/orderbook/domain/resources"
	"github.com/cloudwego/hertz/pkg/common/test/assert"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestIntegration_exchangeMetadataDB_QueryByName(t *testing.T) {
	conf.Init("config_local.json")
	resources.InitDB()

	var (
		db           = resources.GetMasterOrderBookMainDb()
		testExchange = "test.blockchain.com"
	)

	db.Where("exchange = ?", testExchange).Delete(&model.ExchangeMetadata{})

	edb := NewExchangeMetadataDB(db)

	{
		insertModel := &model.ExchangeMetadata{
			Exchange: testExchange,
			Metadata: []byte("{\"a\": \"b\"}"),
		}
		_, err := edb.Update(context.TODO(), &model.UpdateMetadataParam{
			Metadata: insertModel,
		})

		data, err := edb.QueryByName(context.TODO(), &model.QueryMetadataParam{
			ExchangeName: testExchange,
		})
		if err != nil {
			t.Error(err)
			return
		}

		assert.DeepEqual(t, insertModel.Metadata, data.Metadata.Metadata)
	}

	{
		insertModel := &model.ExchangeMetadata{
			Exchange: testExchange,
			Metadata: []byte("{\"c\": \"d\"}"),
		}
		_, err := edb.Update(context.TODO(), &model.UpdateMetadataParam{
			Metadata: insertModel,
		})

		data, err := edb.QueryByName(context.TODO(), &model.QueryMetadataParam{
			ExchangeName: testExchange,
		})
		if err != nil {
			t.Error(err)
			return
		}

		assert.DeepEqual(t, insertModel.Metadata, data.Metadata.Metadata)
	}
}

func TestIntegration_GenerateSql(t *testing.T) {
	t.Skip()

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,         // Disable color
		},
	)

	conf.Init("config_local.json")
	resources.InitDB()

	db := resources.GetMasterOrderBookMainDb().Session(&gorm.Session{DryRun: true, Logger: newLogger})
	err := db.Set("gorm:table_options", "ENGINE=InnoDB").Debug().Migrator().CreateTable(&model.ExchangeMetadata{})
	if err != nil {
		t.Error(err)
		return
	}
}
