module github.com/Mestrace/orderbook

go 1.17

replace github.com/apache/thrift => github.com/apache/thrift v0.13.0

require (
	github.com/antihax/optional v1.0.0
	github.com/apache/thrift v0.0.0-00010101000000-000000000000
	github.com/brianvoe/gofakeit/v6 v6.19.0
	github.com/bytedance/go-tagexpr/v2 v2.9.5
	github.com/bytedance/gopkg v0.0.0-20220413063733-65bf48ffb3a7
	github.com/cloudwego/hertz v0.3.2
	github.com/gocarina/gocsv v0.0.0-20220914091333-ceebdd90b590
	github.com/golang/mock v1.6.0
	github.com/hertz-contrib/monitor-prometheus v0.0.0-20220908085834-f3fe5f5e72ed
	github.com/smartystreets/goconvey v1.7.2
	golang.org/x/net v0.0.0-20220722155237-a158d28d115b // indirect
	golang.org/x/oauth2 v0.0.0-20220909003341-f21342109be1
	gorm.io/driver/mysql v1.3.6
	gorm.io/gorm v1.23.10
)
