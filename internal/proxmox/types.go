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

type Snapshot struct {
	Kind        string `json:"kind"`
	Node        string `json:"node"`
	VMID        int    `json:"vmid"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Parent      string `json:"parent,omitempty"`
	SnapTimeSec int64  `json:"snap_time_sec,omitempty"`
	VMState     bool   `json:"vm_state,omitempty"`
	SnapState   string `json:"snap_state,omitempty"`
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

type Network struct {
	Node            string `json:"node"`
	Iface           string `json:"iface"`
	Type            string `json:"type,omitempty"`
	Active          bool   `json:"active"`
	Autostart       bool   `json:"autostart"`
	Address         string `json:"address,omitempty"`
	Address6        string `json:"address6,omitempty"`
	CIDR            string `json:"cidr,omitempty"`
	CIDR6           string `json:"cidr6,omitempty"`
	Gateway         string `json:"gateway,omitempty"`
	Gateway6        string `json:"gateway6,omitempty"`
	BridgePorts     string `json:"bridge_ports,omitempty"`
	BridgeVLANAware bool   `json:"bridge_vlan_aware,omitempty"`
	VLANID          string `json:"vlan_id,omitempty"`
	VLANRawDevice   string `json:"vlan_raw_device,omitempty"`
	MTU             string `json:"mtu,omitempty"`
	Method          string `json:"method,omitempty"`
	Method6         string `json:"method6,omitempty"`
	Comments        string `json:"comments,omitempty"`
}

type ClusterResource struct {
	ID         string  `json:"id,omitempty"`
	Type       string  `json:"type,omitempty"`
	Node       string  `json:"node,omitempty"`
	VMID       int     `json:"vmid,omitempty"`
	Name       string  `json:"name,omitempty"`
	Status     string  `json:"status,omitempty"`
	Pool       string  `json:"pool,omitempty"`
	Content    string  `json:"content,omitempty"`
	Storage    string  `json:"storage,omitempty"`
	CPU        float64 `json:"cpu,omitempty"`
	MaxCPU     int     `json:"max_cpu,omitempty"`
	Mem        uint64  `json:"mem_bytes,omitempty"`
	MaxMem     uint64  `json:"max_mem_bytes,omitempty"`
	Disk       uint64  `json:"disk_bytes,omitempty"`
	MaxDisk    uint64  `json:"max_disk_bytes,omitempty"`
	Uptime     uint64  `json:"uptime,omitempty"`
	Template   bool    `json:"template,omitempty"`
	Tags       string  `json:"tags,omitempty"`
	PluginType string  `json:"plugin_type,omitempty"`
}
