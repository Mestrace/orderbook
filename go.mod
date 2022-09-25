module github.com/Mestrace/orderbook

go 1.17

replace github.com/apache/thrift => github.com/apache/thrift v0.13.0

require (
	github.com/antihax/optional v1.0.0
	github.com/apache/thrift v0.0.0-00010101000000-000000000000
	github.com/bytedance/gopkg v0.0.0-20220413063733-65bf48ffb3a7
	github.com/cloudwego/hertz v0.3.2
	golang.org/x/net v0.0.0-20220722155237-a158d28d115b // indirect
	golang.org/x/oauth2 v0.0.0-20220909003341-f21342109be1
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gorm.io/driver/mysql v1.3.6
	gorm.io/gorm v1.23.10
)
