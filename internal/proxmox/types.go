package proxmox

type Node struct {
	Node      string  `json:"node"`
	Status    string  `json:"status"`
	CPU       float64 `json:"cpu_ratio"` // 0..1
	MaxCPU    int     `json:"max_cpu"`
	MemUsed   uint64  `json:"mem_used_bytes"`
	MaxMem    uint64  `json:"mem_total_bytes"`
	UptimeSec uint64  `json:"uptime_sec"`
}

type VM struct {
	Node     string  `json:"node"`
	VMID     int     `json:"vmid"`
	Name     string  `json:"name,omitempty"`
	Status   string  `json:"status,omitempty"`
	CPUs     int     `json:"cpus,omitempty"`
	CPU      float64 `json:"cpu,omitempty"`
	Memory   uint64  `json:"memory,omitempty"`
	MaxMem   uint64  `json:"max_mem,omitempty"`
	Disk     uint64  `json:"disk,omitempty"`
	MaxDisk  uint64  `json:"max_disk,omitempty"`
	Uptime   uint64  `json:"uptime,omitempty"`
	Template bool    `json:"template,omitempty"`
	Tags     string  `json:"tags,omitempty"`
}

type VMDetails struct {
	VM
	Config *VMConfig `json:"config,omitempty"`
}

type VMConfig struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	OSType      string `json:"os_type,omitempty"`
	Boot        string `json:"boot,omitempty"`
	OnBoot      bool   `json:"on_boot,omitempty"`
	Agent       string `json:"agent,omitempty"`
	Tags        string `json:"tags,omitempty"`
}
