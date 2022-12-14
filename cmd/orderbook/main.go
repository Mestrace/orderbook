// Code generated by hertz generator.

package main

import (
	"flag"

	"github.com/Mestrace/orderbook/conf"
	"github.com/Mestrace/orderbook/domain/resources"
	"github.com/cloudwego/hertz/pkg/app/server"
	prometheus "github.com/hertz-contrib/monitor-prometheus"
)

func init() {
	parseFlags()
	conf.Init(confFilename)
	resources.Init()
}

var (
	confFilename string
)

func parseFlags() {
	flag.StringVar(&confFilename, "conf", "config.json", "json config, input filename, will find under project_root/conf")
	flag.Parse()
}

func main() {
	h := server.Default(server.WithTracer(prometheus.NewServerTracer(":9091", "/hertz")))

	register(h)
	h.Spin()
}
