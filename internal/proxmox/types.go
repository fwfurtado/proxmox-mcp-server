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

type Container struct {
	Node    string `json:"node"`
	VMID    int    `json:"vmid"`
	Name    string `json:"name,omitempty"`
	Status  string `json:"status,omitempty"`
	CPUs    int    `json:"cpus,omitempty"`
	MaxMem  uint64 `json:"max_mem,omitempty"`
	MaxDisk uint64 `json:"max_disk,omitempty"`
	MaxSwap uint64 `json:"max_swap,omitempty"`
	Uptime  uint64 `json:"uptime,omitempty"`
	Tags    string `json:"tags,omitempty"`
}

type ContainerDetails struct {
	Container
	Config *ContainerConfig `json:"config,omitempty"`
}

type ContainerConfig struct {
	Hostname    string `json:"hostname,omitempty"`
	Description string `json:"description,omitempty"`
	OSType      string `json:"os_type,omitempty"`
	OnBoot      bool   `json:"on_boot,omitempty"`
	Tags        string `json:"tags,omitempty"`
	RootFS      string `json:"root_fs,omitempty"`
}

type Storage struct {
	Node         string  `json:"node"`
	Name         string  `json:"name"`
	Type         string  `json:"type,omitempty"`
	Content      string  `json:"content,omitempty"`
	Active       bool    `json:"active"`
	Enabled      bool    `json:"enabled"`
	Shared       bool    `json:"shared"`
	UsedFraction float64 `json:"used_fraction,omitempty"`
	Avail        uint64  `json:"avail_bytes,omitempty"`
	Used         uint64  `json:"used_bytes,omitempty"`
	Total        uint64  `json:"total_bytes,omitempty"`
}

type Task struct {
	UPID         string `json:"upid"`
	ID           string `json:"id,omitempty"`
	Type         string `json:"type,omitempty"`
	User         string `json:"user,omitempty"`
	Status       string `json:"status,omitempty"`
	Node         string `json:"node,omitempty"`
	PID          uint64 `json:"pid,omitempty"`
	PStart       uint64 `json:"pstart,omitempty"`
	Saved        string `json:"saved,omitempty"`
	ExitStatus   string `json:"exit_status,omitempty"`
	IsCompleted  bool   `json:"is_completed,omitempty"`
	IsRunning    bool   `json:"is_running,omitempty"`
	IsFailed     bool   `json:"is_failed,omitempty"`
	IsSuccessful bool   `json:"is_successful,omitempty"`
	StartTimeSec int64  `json:"start_time_sec,omitempty"`
	EndTimeSec   int64  `json:"end_time_sec,omitempty"`
	DurationSec  int64  `json:"duration_sec,omitempty"`
}

type TaskDetails struct {
	Task
	Log []string `json:"log,omitempty"`
}
