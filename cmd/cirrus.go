package main

import (
	"flag"

	"github.com/lack-io/cirrus/cdiscount"
	"github.com/lack-io/cirrus/config"
	"github.com/lack-io/cirrus/internal/log"
	"github.com/lack-io/cirrus/internal/signal"
)

func main() {
	cfg := flag.String("config", "config", "cirrus.toml")
	flag.Parse()

	err := config.Init(*cfg)
	if err != nil {
		log.Fatalf("初始化备份文件失败: %v", err)
	}
	cds, err := cdiscount.NewCdiscount(config.Get())
	if err != nil {
		log.Fatalf("启动 cdiscount 失败 %v", err)
	}
	cds.Start(signal.SetupSignalHandler())
}
