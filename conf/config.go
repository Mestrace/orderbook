package conf

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"sync"

	"github.com/Mestrace/orderbook/constant"
)

var (
	initOnce sync.Once
	cfg      = new(Config)
)

type Config struct {
	BlockchainCom struct {
		APIKey    string `json:"api_key"`
		APISecret string `json:"api_secret"`
	} `json:"blockchain_com"`
}

// Init initialize the config, will fetch the config under project_root/conf.
func Init(filename string) {
	initOnce.Do(func() {
		var err error
		filepath := path.Join(constant.ProjectRoot, "conf", filename)
		jsonFile, err := os.Open(filepath)
		if err != nil {
			panic("Init config failed")
		}
		defer jsonFile.Close()

		byteValue, _ := ioutil.ReadAll(jsonFile)

		err = json.Unmarshal(byteValue, cfg)
		if err != nil {
			panic("Init config failed")
		}
	})
}

func Get() *Config {
	return cfg
}
