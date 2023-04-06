package main

import (
	"github.com/geneva-lake/tcp_pow/internal/config"
	"github.com/geneva-lake/tcp_pow/internal/wisdom"
	"github.com/geneva-lake/tcp_pow/server/service"

	"github.com/panjf2000/gnet/pkg/logging"
	"github.com/panjf2000/gnet/v2"
)

func main() {
	cfg, err := config.NewConfig[service.Config]().FromFile("server.yaml").Yaml()
	if err != nil {
		logging.Fatalf("config read error=%v", err)
	}
	w := wisdom.NewWisdom()
	s := service.NewServer(w)
	err = s.Init(cfg)
	if err != nil {
		logging.Fatalf("server init error=%v", err)
	}
	err = gnet.Run(s, "tcp://0.0.0.0:"+cfg.Port, gnet.WithMulticore(cfg.Multicore))
	logging.Infof("server run error=%v", err)
}
