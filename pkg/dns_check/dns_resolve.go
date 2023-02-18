package dns_check

import (
	"net"
	"time"
)

func (check *DNSCheck)GetDNSResolve (){
	defer func() {
		if err := recover(); err != nil {
			check.Resolve = ""
			check.RetryCount = check.MaxRetry
		}
	}()
	var result interface{}
	var err error
	retry := 0
	var timeCost int64
	for retry < check.MaxRetry{
		startTime := time.Now().UnixMilli()
		switch check.Type{
		case "CNAME":
			result, err = net.LookupCNAME(check.Domain)
		case "MX":
			result, err = net.LookupMX(check.Domain)
		case "A":
			result, err = net.LookupHost(check.Domain)
		case "NS":
			result, err = net.LookupNS(check.Domain)
		case "TXT":
			result, err = net.LookupTXT(check.Domain)
		default:
			panic("err")
		}
		stopTime := time.Now().UnixMilli()
		timeCost = stopTime - startTime
		if err == nil{
			goto success
		}

		retry += 1
	}
	panic("resolve error")
	success:
		check.Resolve = result
		check.RetryCount = retry
		check.TimeCost = timeCost
}
