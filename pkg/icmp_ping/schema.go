package icmp_ping

type ICMPPing struct {
	// 任务元数据
	Dest string // 目的
	Size int // 包长度
	Packages int // 包数量
	Interval int // 发包间隔

	// 性能数据
	MaxRtt float64
	MinRtt float64
	AvgRtt float64
	Loss float64
}

func GetICMPPing(dest string, size int, packages int, interval int) (ICMPPing){
	return ICMPPing{
		Dest: dest,
		Size: size,
		Packages: packages,
		Interval: interval,
	}
}