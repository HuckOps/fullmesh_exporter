package icmp_ping

import (
	"fmt"
	"fullmesh_exporter/pkg/conf"
	"github.com/go-ping/ping"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"log"
	"sync"
	"time"
)

func (icmpPing *ICMPPing) Ping() {
	retry := 0
	p, err := ping.NewPinger(icmpPing.Dest)
	if err != nil {
		log.Fatalln("ping dest is error")
	}
	p.Count = icmpPing.Packages
	p.Size = icmpPing.Size
	p.Timeout = time.Duration(conf.Config.Ping.TimeOut) * time.Second
	p.Interval = time.Duration(icmpPing.Interval) * time.Millisecond
retryRunPing:
	err = p.Run()
	if err != nil {
		if retry < conf.Config.Ping.Retry {
			retry += 1
			goto retryRunPing
		} else {
			panic("ping error")
		}
	}
	defer func() {
		if err := recover(); err != nil {
			icmpPing.Loss = 100
		}
	}()
	e := p.Statistics()
	icmpPing.AvgRtt = float64(e.AvgRtt / time.Millisecond)
	icmpPing.MaxRtt = float64(e.MaxRtt / time.Millisecond)
	icmpPing.MinRtt = float64(e.MinRtt / time.Millisecond)
	icmpPing.Loss = e.PacketLoss
}

func GetPingTask() {
	var wg sync.WaitGroup
	pingChan := make(chan ICMPPing, conf.Config.Ping.Pool)
	for _, dest := range conf.Config.Ping.MachineList {
		ping := ICMPPing{
			Dest:     dest,
			Size:     conf.Config.Ping.Size,
			Packages: conf.Config.Ping.Packages,
			Interval: conf.Config.Ping.Interval,
		}
		pingChan <- ping
		wg.Add(1)
		go func(pingChain chan ICMPPing, wg *sync.WaitGroup) {
			ping := <-pingChain
			ping.Ping()
			ping.PushPingDataToGateway("ping_max_rtt", ping.MaxRtt)
			ping.PushPingDataToGateway("ping_min_rtt", ping.MinRtt)
			ping.PushPingDataToGateway("ping_avg_rtt", ping.AvgRtt)
			ping.PushPingDataToGateway("ping_loss", ping.Loss)
			wg.Done()
		}(pingChan, &wg)
	}
	wg.Wait()
}

func (icmpPing *ICMPPing) PushPingDataToGateway(metric string, value float64) {
	metricInstance := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: metric,
	})
	metricInstance.Set(value)
	if err := push.New(conf.Config.PushGateway, conf.Config.Job).Collector(metricInstance).Grouping("dest", icmpPing.Dest).
		Grouping("instance", conf.Config.Instance).Grouping("metric", metric).Push(); err != nil {
		log.Fatalln("Push data to gateway failed")
	} else {
		log.Println(fmt.Sprintf("Ping status:  target:{\"instance\": \"%s\" ,\"dest\": \"%s\", \"metric\": \"%s\"} \t"+
			" value: %f",
			conf.Config.Instance, icmpPing.Dest, metric, value))
	}
}
