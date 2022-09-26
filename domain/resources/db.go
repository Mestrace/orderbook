package resources

import (
	"fmt"
	"sync"
	"time"

	"github.com/Mestrace/orderbook/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type dbSet struct {
	Master *gorm.DB
}

var (
	initDBOnce sync.Once
	dbSetMap   map[string]*dbSet
)

// InitDB initialize the database.
func InitDB() error {
	var err error

	initDBOnce.Do(func() {
		err = initDB()
	})

	return err
}

func initDB() error {
	config := conf.Get()

	dbSetMap = make(map[string]*dbSet, len(config.Mysql))

	for dbName, dbConf := range config.Mysql {
		dbMaster, err := gorm.Open(mysql.New(mysql.Config{
			DSN: dbConf.MasterDsn,
			// disable datetime precision, which not supported before MySQL 5.6
			DisableDatetimePrecision: true,
			// drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
			DontSupportRenameIndex: true,
			// `change` when rename column, rename column not supported before MySQL 8, MariaDB
			DontSupportRenameColumn: true,
			// auto configure based on currently MySQL versiona
			SkipInitializeWithVersion: false,
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

// The keys of dbs.
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
