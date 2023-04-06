package main

import (
	"github.com/geneva-lake/tcp_pow/client/service"
	"github.com/geneva-lake/tcp_pow/internal/config"
	"github.com/geneva-lake/tcp_pow/internal/vdf"

	"github.com/panjf2000/gnet/pkg/logging"
	"github.com/panjf2000/gnet/v2"
)

func main() {
	cfg, err := config.NewConfig[service.Config]().FromFile("client.yaml").Yaml()
	if err != nil {
		logging.Fatalf("config read err=%v", err)
	}
	v := vdf.NewVdf()
	ev := service.NewClientEvents(v)
	c, err := gnet.NewClient(ev, gnet.WithMulticore(cfg.Multicore))
	if err != nil {
		logging.Fatalf("client creating error=%v", err)
	}
	_, err = c.Dial("tcp", cfg.Address)
	if err != nil {
		logging.Fatalf("dialing error=%v", err)
	}
	err = c.Start()
	if err != nil {
		logging.Fatalf("client start error=%v", err)
	}
	<-ev.Stop()
}
