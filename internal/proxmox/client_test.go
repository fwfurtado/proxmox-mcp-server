package proxmox

import (
	"context"
	"testing"
	"time"

	proxmoxlib "github.com/luthermonson/go-proxmox"
)

func TestNewClientValidatesRequiredConfig(t *testing.T) {
	tests := []struct {
		name   string
		config Config
		want   string
	}{
		{name: "missing url", config: Config{}, want: "proxmox URL is required"},
		{name: "missing token id", config: Config{URL: "https://pve.example.com"}, want: "proxmox token ID is required"},
		{name: "missing token secret", config: Config{URL: "https://pve.example.com", TokenID: "root@pam!mcp"}, want: "proxmox token secret is required"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewClient(context.Background(), tt.config)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if err.Error() != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, err.Error())
			}
		})
	}
}

func TestVMSummary(t *testing.T) {
	vm := &proxmoxlib.VirtualMachine{
		Node:     "pve01",
		VMID:     proxmoxlib.StringOrUint64(101),
		Name:     "app-01",
		Status:   "running",
		CPUs:     4,
		CPU:      0.42,
		Mem:      2048,
		MaxMem:   4096,
		Disk:     1024,
		MaxDisk:  8192,
		Uptime:   3600,
		Template: proxmoxlib.IsTemplate(true),
		Tags:     "prod;linux",
	}

	got := vmSummary(vm)

	if got.Node != "pve01" || got.VMID != 101 || got.Name != "app-01" {
		t.Fatalf("unexpected VM identity: %+v", got)
	}
	if !got.Template {
		t.Fatal("expected template=true")
	}
	if got.Tags != "prod;linux" {
		t.Fatalf("expected tags to be preserved, got %q", got.Tags)
	}
}

func TestVMSummaryFromResource(t *testing.T) {
	resource := &proxmoxlib.ClusterResource{
		Type:     "qemu",
		Node:     "pve02",
		VMID:     202,
		Name:     "db-01",
		Status:   "stopped",
		MaxCPU:   8,
		CPU:      0.15,
		Mem:      4096,
		MaxMem:   8192,
		Disk:     2048,
		MaxDisk:  16384,
		Uptime:   7200,
		Template: 1,
		Tags:     "db;critical",
	}

	got := vmSummaryFromResource(resource)

	if got.Node != "pve02" || got.VMID != 202 || got.CPUs != 8 {
		t.Fatalf("unexpected VM summary: %+v", got)
	}
	if !got.Template {
		t.Fatal("expected template=true")
	}
}

func TestVMDetails(t *testing.T) {
	osType := "l26"
	vm := &proxmoxlib.VirtualMachine{
		Node:     "pve01",
		VMID:     proxmoxlib.StringOrUint64(101),
		Name:     "app-01",
		Status:   "running",
		Template: proxmoxlib.IsTemplate(false),
		VirtualMachineConfig: &proxmoxlib.VirtualMachineConfig{
			Name:        "app-01",
			Description: "application server",
			OSType:      &osType,
			Boot:        "order=scsi0",
			Agent:       "enabled=1",
			Tags:        "prod",
			OnBoot:      proxmoxlib.IntOrBool(true),
		},
	}

	got := vmDetails(vm)

	if got.Config == nil {
		t.Fatal("expected config to be populated")
	}
	if got.Config.OSType != "l26" {
		t.Fatalf("expected OSType l26, got %q", got.Config.OSType)
	}
	if !got.Config.OnBoot {
		t.Fatal("expected OnBoot=true")
	}
}

func TestVMDetailsWithoutConfig(t *testing.T) {
	vm := &proxmoxlib.VirtualMachine{
		Node: "pve01",
		VMID: proxmoxlib.StringOrUint64(101),
		Name: "app-01",
	}

	got := vmDetails(vm)

	if got.Config != nil {
		t.Fatalf("expected nil config, got %+v", got.Config)
	}
}

func TestVMSnapshotSummary(t *testing.T) {
	snapshot := &proxmoxlib.VirtualMachineSnapshot{
		Node:        "pve01",
		VMID:        101,
		Name:        "snap1",
		Description: "before update",
		Parent:      "base",
		Snaptime:    1234,
		Vmstate:     1,
		Snapstate:   "prepare",
	}

	got := vmSnapshotSummary(snapshot)

	if got.Kind != "vm" || got.Name != "snap1" || !got.VMState {
		t.Fatalf("unexpected VM snapshot summary: %+v", got)
	}
}

func TestContainerSummary(t *testing.T) {
	container := &proxmoxlib.Container{
		Node:    "pve01",
		VMID:    proxmoxlib.StringOrUint64(201),
		Name:    "ct-01",
		Status:  "running",
		CPUs:    2,
		MaxMem:  2048,
		MaxDisk: 8192,
		MaxSwap: 1024,
		Uptime:  1800,
		Tags:    "infra;lxc",
	}

	got := containerSummary(container)

	if got.Node != "pve01" || got.VMID != 201 || got.Name != "ct-01" {
		t.Fatalf("unexpected container identity: %+v", got)
	}
	if got.MaxSwap != 1024 {
		t.Fatalf("expected max swap 1024, got %d", got.MaxSwap)
	}
}

func TestContainerSummaryFromResource(t *testing.T) {
	resource := &proxmoxlib.ClusterResource{
		Type:    "lxc",
		Node:    "pve02",
		VMID:    202,
		Name:    "ct-02",
		Status:  "stopped",
		MaxCPU:  4,
		MaxMem:  4096,
		MaxDisk: 16384,
		Uptime:  3600,
		Tags:    "edge",
	}

	got := containerSummaryFromResource(resource)

	if got.Node != "pve02" || got.VMID != 202 || got.CPUs != 4 {
		t.Fatalf("unexpected container summary: %+v", got)
	}
}

func TestContainerDetails(t *testing.T) {
	container := &proxmoxlib.Container{
		Node: "pve01",
		VMID: proxmoxlib.StringOrUint64(201),
		Name: "ct-01",
		ContainerConfig: &proxmoxlib.ContainerConfig{
			Hostname:    "ct-01",
			Description: "container server",
			OSType:      "debian",
			OnBoot:      proxmoxlib.IntOrBool(true),
			Tags:        "prod",
			RootFS:      "local-lvm:subvol-201-disk-0",
		},
	}

	got := containerDetails(container)

	if got.Config == nil {
		t.Fatal("expected config to be populated")
	}
	if got.Config.OSType != "debian" || !got.Config.OnBoot {
		t.Fatalf("unexpected container config: %+v", got.Config)
	}
}

func TestContainerSnapshotSummary(t *testing.T) {
	snapshot := &proxmoxlib.ContainerSnapshot{
		Node:                 "pve01",
		VMID:                 201,
		Name:                 "snap1",
		Description:          "before upgrade",
		Parent:               "base",
		SnapshotCreationTime: 5678,
	}

	got := containerSnapshotSummary(snapshot)

	if got.Kind != "container" || got.Name != "snap1" || got.SnapTimeSec != 5678 {
		t.Fatalf("unexpected container snapshot summary: %+v", got)
	}
}

func TestStorageSummary(t *testing.T) {
	storage := &proxmoxlib.Storage{
		Node:         "pve01",
		Name:         "local-lvm",
		Type:         "lvmthin",
		Content:      "images,rootdir",
		Active:       1,
		Enabled:      1,
		Shared:       0,
		UsedFraction: 0.25,
		Avail:        100,
		Used:         50,
		Total:        150,
	}

	got := storageSummary(storage)

	if got.Name != "local-lvm" || !got.Active || !got.Enabled || got.Shared {
		t.Fatalf("unexpected storage summary: %+v", got)
	}
}

func TestNetworkSummary(t *testing.T) {
	network := &proxmoxlib.NodeNetwork{
		Node:            "pve01",
		Iface:           "vmbr0",
		Type:            "bridge",
		Active:          proxmoxlib.StringOrInt(1),
		Autostart:       1,
		Address:         "10.0.0.2",
		BridgePorts:     "eno1",
		BridgeVLANAware: 1,
		VLANID:          "100",
		Method:          "static",
	}

	got := networkSummary(network)

	if got.Node != "pve01" || got.Iface != "vmbr0" || !got.Active || !got.Autostart || !got.BridgeVLANAware {
		t.Fatalf("unexpected network summary: %+v", got)
	}
}

func TestClusterResourceSummary(t *testing.T) {
	resource := &proxmoxlib.ClusterResource{
		ID:         "qemu/101",
		Type:       "qemu",
		Node:       "pve01",
		VMID:       101,
		Name:       "app-01",
		Status:     "running",
		Pool:       "prod",
		CPU:        0.2,
		MaxCPU:     4,
		Mem:        2048,
		MaxMem:     4096,
		Disk:       1024,
		MaxDisk:    8192,
		Uptime:     100,
		Template:   1,
		Tags:       "prod",
		PluginType: "qemu",
	}

	got := clusterResourceSummary(resource)

	if got.ID != "qemu/101" || got.VMID != 101 || !got.Template {
		t.Fatalf("unexpected cluster resource summary: %+v", got)
	}
}

func TestTaskSummary(t *testing.T) {
	start := time.Unix(1000, 0)
	end := time.Unix(1060, 0)
	task := &proxmoxlib.Task{
		UPID:         "UPID:pve01:00000001:00000002:00000003:qmstart:101:root@pam:",
		ID:           "101",
		Type:         "qmstart",
		User:         "root@pam",
		Status:       "stopped",
		Node:         "pve01",
		ExitStatus:   "OK",
		IsCompleted:  true,
		IsSuccessful: true,
		StartTime:    start,
		EndTime:      end,
		Duration:     end.Sub(start),
	}

	got := taskSummary(task)

	if got.UPID != string(task.UPID) || got.DurationSec != 60 {
		t.Fatalf("unexpected task summary: %+v", got)
	}
	if got.StartTimeSec != 1000 || got.EndTimeSec != 1060 {
		t.Fatalf("unexpected task timestamps: %+v", got)
	}
}

func TestTaskLogLines(t *testing.T) {
	log := proxmoxlib.Log{
		2: "third",
		0: "first",
		1: "second",
	}

	got := taskLogLines(log)

	want := []string{"first", "second", "third"}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("expected %v, got %v", want, got)
		}
	}
}

func TestGetTaskValidatesInput(t *testing.T) {
	client := &Client{}

	_, err := client.GetTask(context.Background(), "", 0, 50)
	if err == nil || err.Error() != "upid is required" {
		t.Fatalf("expected upid validation error, got %v", err)
	}

	_, err = client.GetTask(context.Background(), "UPID:pve", -1, 50)
	if err == nil || err.Error() != "log_start must be greater than or equal to zero" {
		t.Fatalf("expected log_start validation error, got %v", err)
	}
}

func TestGetVMConfigValidatesInput(t *testing.T) {
	client := &Client{}

	_, err := client.GetVMConfig(context.Background(), "", 101)
	if err == nil || err.Error() != "node name is required" {
		t.Fatalf("expected node name validation error, got %v", err)
	}
}
