package resources

import (
	"fmt"
	"time"

	"github.com/Mestrace/orderbook/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type dbSet struct {
	Master *gorm.DB
}

var (
	dbSetMap map[string]*dbSet
)

func InitDB() error {
	var (
		config = conf.Get()
	)

	dbSetMap = make(map[string]*dbSet, len(config.Mysql))

	for dbName, dbConf := range config.Mysql {
		dbMaster, err := gorm.Open(mysql.New(mysql.Config{
			DSN:                       dbConf.MasterDsn, // data source name
			DisableDatetimePrecision:  true,             // disable datetime precision, which not supported before MySQL 5.6
			DontSupportRenameIndex:    true,             // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
			DontSupportRenameColumn:   true,             // `change` when rename column, rename column not supported before MySQL 8, MariaDB
			SkipInitializeWithVersion: false,            // auto configure based on currently MySQL versiona
		}), &gorm.Config{
			SkipDefaultTransaction: true,
			PrepareStmt:            true,
		})
		if err != nil {
			return err
		}
		// set conn pool config
		{
			sqlDb, err := dbMaster.DB()
			if err != nil {
				return err
			}
			sqlDb.SetConnMaxIdleTime(time.Duration(dbConf.ConnMaxIdleTime))
			sqlDb.SetConnMaxLifetime(time.Duration(dbConf.ConnMaxLifeTime))
			sqlDb.SetMaxIdleConns(dbConf.MaxIdleConn)
			sqlDb.SetMaxOpenConns(dbConf.MaxOpenConn)
		}

		dbSetMap[dbName] = &dbSet{
			Master: dbMaster,
		}
	}
	return nil
}

// keys
const (
	keyDbOrderBookMain = "db_orderbook_main"
)

func mustGetDbSet(key string) *dbSet {
	set, ok := dbSetMap[key]
	if !ok {
		panic(fmt.Sprintf("key not exist: %s", key))
	}
	return set
}

func GetMasterOrderBookMainDb() *gorm.DB {
	set := mustGetDbSet(keyDbOrderBookMain)

	return set.Master
}
