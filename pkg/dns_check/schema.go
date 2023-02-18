package dns_check

type DNSCheck struct {
	// 基础信息设置
	Domain string
	Type string

	// 解析收敛
	MaxRetry int

	// 以本机resolv为服务器的解析结果
	Resolve interface{}

	// 解析性能分析
	TimeCost int64
	RetryCount int
}


func GetDnsCheckTask(domain string, resolveType string, maxRetry int) (DNSCheck) {
	if domain == ""{
		panic("Domain can not be a empty string")
	}
	if maxRetry < 0 || maxRetry > 100{
		panic("Retry time must in [1,100]")
	}
	return DNSCheck{
		Domain: domain,
		Type: resolveType,
		MaxRetry: maxRetry,
	}
}