package proxmox

import (
	"context"
	"testing"

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
