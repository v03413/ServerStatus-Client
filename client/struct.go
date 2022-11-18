package client

type traffic struct {
	in  uint64
	out uint64
}

type memory struct {
	total uint64
	used  uint64
}

type swap struct {
	total uint64
	used  uint64
}

type diskStat struct {
	size uint64
	used uint64
}

type update struct {
	Cpu         float64 `json:"cpu"`
	Load1       float64 `json:"load_1"`
	Load5       float64 `json:"load_5"`
	Load15      float64 `json:"load_15"`
	Uptime      uint64  `json:"uptime"`
	MemoryTotal uint64  `json:"memory_total"`
	MemoryUsed  uint64  `json:"memory_used"`
	SwapTotal   uint64  `json:"swap_total"`
	SwapUsed    uint64  `json:"swap_used"`
	HddTotal    uint64  `json:"hdd_total"`
	HddUsed     uint64  `json:"hdd_used"`
	NetWorkIn   uint64  `json:"network_in"`
	NetWorkOut  uint64  `json:"network_out"`
	IpStatus    bool    `json:"ip_status"`
	NetWorkRx   uint64  `json:"network_rx"`
	NetWorkTx   uint64  `json:"network_tx"`
	PingCU      float64 `json:"ping_10010"`
	PingCM      float64 `json:"ping_10086"`
	PingCT      float64 `json:"ping_189"`
	TimeCU      uint    `json:"time_10010"`
	TimeCM      uint    `json:"time_10086"`
	TimeCT      uint    `json:"time_189"`
	Tcp         uint    `json:"tcp"`
	Udp         uint    `json:"udp"`
	Process     uint    `json:"process"`
	Thread      uint    `json:"thread"`
	IoRead      uint64  `json:"io_read"`
	IoWrite     uint64  `json:"io_write"`
}

type rateStat struct {
	recvBytes uint64
	sendBytes uint64
	second    uint64
}

type ioStat struct {
	readBytes  uint64
	writeBytes uint64
	second     uint64
}

type tupdStat struct {
	tcp     uint
	udp     uint
	process uint
	thread  uint
}
