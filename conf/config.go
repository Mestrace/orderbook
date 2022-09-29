package conf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"sync"
	"time"

	"github.com/Mestrace/orderbook/constant"
)

var (
	initOnce sync.Once
	cfg      = new(Config)
)

type Config struct {
	BlockchainCom struct {
		APIKey                  string         `json:"api_key"`
		APISecret               string         `json:"api_secret"`
		QPSLimit                int            `json:"qps_limit"`
		SymbolListCacheDuration ConfigDuration `json:"symbol_list_cache_duration"`
		SymbolCacheDuration     ConfigDuration `json:"symbol_cache_duration"`
	} `json:"blockchain_com"`
	Mysql map[string]*struct {
		MasterDsn       string         `json:"master_dsn"`
		MaxIdleConn     int            `json:"max_idle_conn"`
		MaxOpenConn     int            `json:"max_open_conn"`
		ConnMaxLifeTime ConfigDuration `json:"conn_max_life_time"`
		ConnMaxIdleTime ConfigDuration `json:"conn_max_idle_time"`
	} `json:"mysql"`
	Redis struct {
		Addr         string         `json:"addr"`
		Password     string         `json:"password"`
		Db           int            `json:"db"`
		MaxRetries   int            `json:"max_retries"`
		DialTimeout  ConfigDuration `json:"dial_timeout"`
		ReadTimeout  ConfigDuration `json:"read_timeout"`
		WriteTimeout ConfigDuration `json:"write_timeout"`
	} `json:"redis"`
}

// Init initialize the config, will fetch the config under project_root/conf.
func Init(filename string) {
	initOnce.Do(func() {
		var err error
		filepath := path.Join(constant.ProjectRoot, "conf", filename)
		jsonFile, err := os.Open(filepath)
		if err != nil {
			panic(fmt.Sprintf("Init config failed|path=%s|err=%+v", filepath, err))
		}
		defer jsonFile.Close()

		byteValue, _ := ioutil.ReadAll(jsonFile)

		err = json.Unmarshal(byteValue, cfg)
		if err != nil {
			panic(fmt.Sprintf("Init config failed|err=%+v", err))
		}
	})
}

func Get() *Config {
	if cfg == nil {
		panic("conf not initialized")
	}

	return cfg
}

// ConfigDuration maps config duration string to actual time duration for json marshalling.
type ConfigDuration time.Duration

func (d *ConfigDuration) UnmarshalJSON(data []byte) error {
	unquoted, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}

	duration, err := time.ParseDuration(unquoted)
	if err != nil {
		return err
	}

	*d = ConfigDuration(duration)

	return nil
}
