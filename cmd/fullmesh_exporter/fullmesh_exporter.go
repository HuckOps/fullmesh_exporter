package main

import (
	"flag"
	"fullmesh_exporter/pkg/conf"
	"fullmesh_exporter/pkg/icmp_ping"
)

func init() {
	var configPath string
	flag.StringVar(&configPath, "conf", "./conf.yaml", "配置文件地址")
	conf.ReadConfigFile(configPath)
}

func main() {
	//c := cron.New(cron.WithSeconds())
	//c.AddFunc(conf.Config.Cron, func() {
	//	dns_check.GetCheckTaskPool()
	//	icmp_ping.GetPingTask()
	//})
	//c.Start()
	//select {}
	icmp_ping.GetPingTask()
}
