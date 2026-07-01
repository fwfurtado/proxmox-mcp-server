package proxmox

import (
	"context"
	"fmt"
	"time"

	proxmoxlib "github.com/luthermonson/go-proxmox"
	"github.com/samber/lo"
)

type Config struct {
	URL         string
	TokenID     string
	TokenSecret string
	InsecureTLS bool
}

type Client struct {
	proxmox *proxmoxlib.Client
}

func NewClient(ctx context.Context, config Config) (*Client, error) {
	if config.URL == "" {
		return nil, fmt.Errorf("proxmox URL is required")
	}
	if config.TokenID == "" {
		return nil, fmt.Errorf("proxmox token ID is required")
	}
	if config.TokenSecret == "" {
		return nil, fmt.Errorf("proxmox token secret is required")
	}

	options := []proxmoxlib.Option{
		proxmoxlib.WithTimeout(30 * time.Second),
		proxmoxlib.WithAPIToken(config.TokenID, config.TokenSecret),
	}

	if config.InsecureTLS {
		options = append(options, proxmoxlib.WithInsecureSkipVerify())
	}

	client := proxmoxlib.NewClient(config.URL, options...)

	_, err := client.Version(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to Proxmox: %w", err)
	}

	return &Client{proxmox: client}, nil
}

func (c *Client) ListNodes(ctx context.Context) ([]*Node, error) {
	statuses, err := c.proxmox.Nodes(ctx)

	if err != nil {
		return nil, err
	}

	nodes := lo.Map(statuses, func(nodeStatus *proxmoxlib.NodeStatus, index int) *Node {
		return &Node{
			Node:      nodeStatus.Node,
			Status:    nodeStatus.Status,
			CPU:       nodeStatus.CPU,
			MaxCPU:    nodeStatus.MaxCPU,
			MemUsed:   nodeStatus.Mem,
			MaxMem:    nodeStatus.MaxMem,
			UptimeSec: nodeStatus.Uptime,
		}
	})

	return nodes, nil
}

func (c *Client) ListVMs(ctx context.Context) ([]*VM, error) {
	cluster, err := c.proxmox.Cluster(ctx)
	if err != nil {
		return nil, fmt.Errorf("get cluster: %w", err)
	}

	resources, err := cluster.Resources(ctx, "vm")
	if err != nil {
		return nil, fmt.Errorf("list cluster VM resources: %w", err)
	}

	vms := make([]*VM, 0, len(resources))
	for _, resource := range resources {
		if resource.Type != "qemu" {
			continue
		}

		vms = append(vms, vmSummaryFromResource(resource))
	}

	return vms, nil
}

func (c *Client) GetVM(ctx context.Context, nodeName string, vmid int) (*VMDetails, error) {
	if nodeName == "" {
		return nil, fmt.Errorf("node name is required")
	}
	if vmid <= 0 {
		return nil, fmt.Errorf("vmid must be greater than zero")
	}

	node, err := c.proxmox.Node(ctx, nodeName)
	if err != nil {
		return nil, fmt.Errorf("get node %q: %w", nodeName, err)
	}

	vm, err := node.VirtualMachine(ctx, vmid)
	if err != nil {
		return nil, fmt.Errorf("get VM %d on node %q: %w", vmid, nodeName, err)
	}

	return vmDetails(vm), nil
}

func vmSummary(vm *proxmoxlib.VirtualMachine) *VM {
	return &VM{
		Node:     vm.Node,
		VMID:     int(vm.VMID),
		Name:     vm.Name,
		Status:   vm.Status,
		CPUs:     vm.CPUs,
		CPU:      vm.CPU,
		Memory:   vm.Mem,
		MaxMem:   vm.MaxMem,
		Disk:     vm.Disk,
		MaxDisk:  vm.MaxDisk,
		Uptime:   vm.Uptime,
		Template: bool(vm.Template),
		Tags:     vm.Tags,
	}
}

func vmSummaryFromResource(resource *proxmoxlib.ClusterResource) *VM {
	return &VM{
		Node:     resource.Node,
		VMID:     int(resource.VMID),
		Name:     resource.Name,
		Status:   resource.Status,
		CPUs:     int(resource.MaxCPU),
		CPU:      resource.CPU,
		Memory:   resource.Mem,
		MaxMem:   resource.MaxMem,
		Disk:     resource.Disk,
		MaxDisk:  resource.MaxDisk,
		Uptime:   resource.Uptime,
		Template: resource.Template == 1,
		Tags:     resource.Tags,
	}
}

func vmDetails(vm *proxmoxlib.VirtualMachine) *VMDetails {
	details := &VMDetails{
		VM: *vmSummary(vm),
	}

	if vm.VirtualMachineConfig != nil {
		details.Config = &VMConfig{
			Name:        vm.VirtualMachineConfig.Name,
			Description: vm.VirtualMachineConfig.Description,
			Boot:        vm.VirtualMachineConfig.Boot,
			Agent:       vm.VirtualMachineConfig.Agent,
			Tags:        vm.VirtualMachineConfig.Tags,
			OnBoot:      bool(vm.VirtualMachineConfig.OnBoot),
		}

		if vm.VirtualMachineConfig.OSType != nil {
			details.Config.OSType = *vm.VirtualMachineConfig.OSType
		}
	}

	return details
}
