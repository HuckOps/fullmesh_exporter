package dns_check

import (
	"fmt"
	"fullmesh_exporter/pkg/conf"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"log"
	"sync"
)

func GetCheckTaskPool () {
	var wg sync.WaitGroup
	domainChan := make(chan conf.Domains, conf.Config.DnsCheck.Pool)
	for _, domain := range conf.Config.DnsCheck.Domains{
		domainChan <- domain
		wg.Add(1)
		go func(domainChan chan conf.Domains,wg *sync.WaitGroup) {
			domain := <- domainChan
			dnsCheck := GetDnsCheckTask(domain.Resolve, domain.Type, conf.Config.DnsCheck.Retry)
			dnsCheck.GetDNSResolve()

			// 查询状态推送
			if dnsCheck.Resolve == "" {
				dnsCheck.PushDNSDataToGateway("resolve_domain_status", 0)
			}else{
				// 查询时间
				dnsCheck.PushDNSDataToGateway("resolve_domain_time_cost", float64(dnsCheck.TimeCost))
				dnsCheck.PushDNSDataToGateway("resolve_domain_status", 1)
			}
			wg.Done()
		}(domainChan, &wg)
	}
	wg.Wait()
}


func (check *DNSCheck)PushDNSDataToGateway(metric string, value float64){
	metricInstance := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: metric,
	})
	metricInstance.Set(value)
	if err := push.New(conf.Config.PushGateway, conf.Config.Job).Collector(metricInstance).Grouping("domain", check.Domain).
		Grouping("type", check.Type).Grouping("instance", conf.Config.Instance).Grouping("metric", metric).Push(); err != nil {
		log.Fatalln("Push data to gateway failed")
	}else {
		log.Println(fmt.Sprintf("Ping status:  target:{\"instance\": \"%s\" ,\"domain\": \"%s\", \"type\": \"%s\", \"metric\": \"%s\"} \t" +
			" value: %f",
			conf.Config.Instance, check.Domain, check.Type, metric, value))
	}
}



