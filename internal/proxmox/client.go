package proxmox

import (
	"context"
	"fmt"
	"sort"
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

func (c *Client) GetVMConfig(ctx context.Context, nodeName string, vmid int) (*VMConfig, error) {
	vm, err := c.GetVM(ctx, nodeName, vmid)
	if err != nil {
		return nil, err
	}
	if vm.Config == nil {
		return nil, fmt.Errorf("VM %d on node %q has no config available", vmid, nodeName)
	}

	return vm.Config, nil
}

func (c *Client) ListContainers(ctx context.Context) ([]*Container, error) {
	cluster, err := c.proxmox.Cluster(ctx)
	if err != nil {
		return nil, fmt.Errorf("get cluster: %w", err)
	}

	resources, err := cluster.Resources(ctx, "vm")
	if err != nil {
		return nil, fmt.Errorf("list cluster container resources: %w", err)
	}

	containers := make([]*Container, 0, len(resources))
	for _, resource := range resources {
		if resource.Type != "lxc" {
			continue
		}

		containers = append(containers, containerSummaryFromResource(resource))
	}

	return containers, nil
}

func (c *Client) GetContainer(ctx context.Context, nodeName string, vmid int) (*ContainerDetails, error) {
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

	container, err := node.Container(ctx, vmid)
	if err != nil {
		return nil, fmt.Errorf("get container %d on node %q: %w", vmid, nodeName, err)
	}

	return containerDetails(container), nil
}

func (c *Client) ListSnapshots(ctx context.Context, kind, nodeName string, vmid int) ([]*Snapshot, error) {
	if nodeName == "" {
		return nil, fmt.Errorf("node name is required")
	}
	if vmid <= 0 {
		return nil, fmt.Errorf("vmid must be greater than zero")
	}
	if kind == "" {
		kind = "vm"
	}

	node, err := c.proxmox.Node(ctx, nodeName)
	if err != nil {
		return nil, fmt.Errorf("get node %q: %w", nodeName, err)
	}

	switch kind {
	case "vm":
		vm, err := node.VirtualMachine(ctx, vmid)
		if err != nil {
			return nil, fmt.Errorf("get VM %d on node %q for snapshots: %w", vmid, nodeName, err)
		}

		snapshots, err := vm.Snapshots(ctx)
		if err != nil {
			return nil, fmt.Errorf("list VM %d snapshots on node %q: %w", vmid, nodeName, err)
		}

		return lo.Map(snapshots, func(snapshot *proxmoxlib.VirtualMachineSnapshot, _ int) *Snapshot {
			return vmSnapshotSummary(snapshot)
		}), nil
	case "container":
		container, err := node.Container(ctx, vmid)
		if err != nil {
			return nil, fmt.Errorf("get container %d on node %q for snapshots: %w", vmid, nodeName, err)
		}

		snapshots, err := container.Snapshots(ctx)
		if err != nil {
			return nil, fmt.Errorf("list container %d snapshots on node %q: %w", vmid, nodeName, err)
		}

		return lo.Map(snapshots, func(snapshot *proxmoxlib.ContainerSnapshot, _ int) *Snapshot {
			return containerSnapshotSummary(snapshot)
		}), nil
	default:
		return nil, fmt.Errorf("unsupported snapshot kind %q: expected \"vm\" or \"container\"", kind)
	}
}

func (c *Client) ListStorage(ctx context.Context) ([]*Storage, error) {
	nodeStatuses, err := c.proxmox.Nodes(ctx)
	if err != nil {
		return nil, fmt.Errorf("list nodes for storage: %w", err)
	}

	storages := make([]*Storage, 0)
	for _, nodeStatus := range nodeStatuses {
		node, err := c.proxmox.Node(ctx, nodeStatus.Node)
		if err != nil {
			return nil, fmt.Errorf("get node %q for storage: %w", nodeStatus.Node, err)
		}

		nodeStorages, err := node.Storages(ctx)
		if err != nil {
			return nil, fmt.Errorf("list storage on node %q: %w", nodeStatus.Node, err)
		}

		for _, storage := range nodeStorages {
			storages = append(storages, storageSummary(storage))
		}
	}

	return storages, nil
}

func (c *Client) ListNetworks(ctx context.Context, ifaceType string) ([]*Network, error) {
	nodeStatuses, err := c.proxmox.Nodes(ctx)
	if err != nil {
		return nil, fmt.Errorf("list nodes for networks: %w", err)
	}

	networks := make([]*Network, 0)
	for _, nodeStatus := range nodeStatuses {
		nodeNetworks, err := c.ListNodeNetworks(ctx, nodeStatus.Node, ifaceType)
		if err != nil {
			return nil, err
		}
		networks = append(networks, nodeNetworks...)
	}

	return networks, nil
}

func (c *Client) ListNodeNetworks(ctx context.Context, nodeName, ifaceType string) ([]*Network, error) {
	if nodeName == "" {
		return nil, fmt.Errorf("node name is required")
	}

	node, err := c.proxmox.Node(ctx, nodeName)
	if err != nil {
		return nil, fmt.Errorf("get node %q: %w", nodeName, err)
	}

	var rawNetworks proxmoxlib.NodeNetworks
	if ifaceType != "" {
		rawNetworks, err = node.Networks(ctx, ifaceType)
	} else {
		rawNetworks, err = node.Networks(ctx)
	}
	if err != nil {
		return nil, fmt.Errorf("list networks on node %q: %w", nodeName, err)
	}

	return lo.Map(rawNetworks, func(network *proxmoxlib.NodeNetwork, _ int) *Network {
		return networkSummary(network)
	}), nil
}

func (c *Client) ListClusterResources(ctx context.Context, filter string) ([]*ClusterResource, error) {
	cluster, err := c.proxmox.Cluster(ctx)
	if err != nil {
		return nil, fmt.Errorf("get cluster: %w", err)
	}

	var resources proxmoxlib.ClusterResources
	if filter != "" {
		resources, err = cluster.Resources(ctx, filter)
	} else {
		resources, err = cluster.Resources(ctx)
	}
	if err != nil {
		return nil, fmt.Errorf("list cluster resources: %w", err)
	}

	return lo.Map(resources, func(resource *proxmoxlib.ClusterResource, _ int) *ClusterResource {
		return clusterResourceSummary(resource)
	}), nil
}

func (c *Client) ListTasks(ctx context.Context) ([]*Task, error) {
	cluster, err := c.proxmox.Cluster(ctx)
	if err != nil {
		return nil, fmt.Errorf("get cluster: %w", err)
	}

	tasks, err := cluster.Tasks(ctx)
	if err != nil {
		return nil, fmt.Errorf("list cluster tasks: %w", err)
	}

	return lo.Map(tasks, func(task *proxmoxlib.Task, _ int) *Task {
		return taskSummary(task)
	}), nil
}

func (c *Client) ListNodeTasks(ctx context.Context, nodeName string, limit int) ([]*Task, error) {
	if nodeName == "" {
		return nil, fmt.Errorf("node name is required")
	}

	node, err := c.proxmox.Node(ctx, nodeName)
	if err != nil {
		return nil, fmt.Errorf("get node %q: %w", nodeName, err)
	}

	opts := &proxmoxlib.NodeTasksOptions{}
	if limit > 0 {
		opts.Limit = limit
	}

	tasks, err := node.Tasks(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("list tasks on node %q: %w", nodeName, err)
	}

	return lo.Map(tasks, func(task *proxmoxlib.Task, _ int) *Task {
		return taskSummary(task)
	}), nil
}

func (c *Client) GetTask(ctx context.Context, upid string, logStart, logLimit int) (*TaskDetails, error) {
	if upid == "" {
		return nil, fmt.Errorf("upid is required")
	}
	if logStart < 0 {
		return nil, fmt.Errorf("log_start must be greater than or equal to zero")
	}
	if logLimit <= 0 {
		logLimit = 50
	}

	task := proxmoxlib.NewTask(proxmoxlib.UPID(upid), c.proxmox)
	if task == nil || task.Node == "" {
		return nil, fmt.Errorf("invalid upid %q", upid)
	}

	if err := task.Ping(ctx); err != nil {
		return nil, fmt.Errorf("get task %q status: %w", upid, err)
	}

	logLines, err := task.Log(ctx, logStart, logLimit)
	if err != nil {
		return nil, fmt.Errorf("get task %q log: %w", upid, err)
	}

	return &TaskDetails{
		Task: *taskSummary(task),
		Log:  taskLogLines(logLines),
	}, nil
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

func vmSnapshotSummary(snapshot *proxmoxlib.VirtualMachineSnapshot) *Snapshot {
	return &Snapshot{
		Kind:        "vm",
		Node:        snapshot.Node,
		VMID:        snapshot.VMID,
		Name:        snapshot.Name,
		Description: snapshot.Description,
		Parent:      snapshot.Parent,
		SnapTimeSec: snapshot.Snaptime,
		VMState:     snapshot.Vmstate == 1,
		SnapState:   snapshot.Snapstate,
	}
}

func containerSummary(container *proxmoxlib.Container) *Container {
	return &Container{
		Node:    container.Node,
		VMID:    int(container.VMID),
		Name:    container.Name,
		Status:  container.Status,
		CPUs:    container.CPUs,
		MaxMem:  container.MaxMem,
		MaxDisk: container.MaxDisk,
		MaxSwap: container.MaxSwap,
		Uptime:  container.Uptime,
		Tags:    container.Tags,
	}
}

func containerSummaryFromResource(resource *proxmoxlib.ClusterResource) *Container {
	return &Container{
		Node:    resource.Node,
		VMID:    int(resource.VMID),
		Name:    resource.Name,
		Status:  resource.Status,
		CPUs:    int(resource.MaxCPU),
		MaxMem:  resource.MaxMem,
		MaxDisk: resource.MaxDisk,
		Uptime:  resource.Uptime,
		Tags:    resource.Tags,
	}
}

func containerDetails(container *proxmoxlib.Container) *ContainerDetails {
	details := &ContainerDetails{
		Container: *containerSummary(container),
	}

	if container.ContainerConfig != nil {
		details.Config = &ContainerConfig{
			Hostname:    container.ContainerConfig.Hostname,
			Description: container.ContainerConfig.Description,
			OSType:      container.ContainerConfig.OSType,
			OnBoot:      bool(container.ContainerConfig.OnBoot),
			Tags:        container.ContainerConfig.Tags,
			RootFS:      container.ContainerConfig.RootFS,
		}
	}

	return details
}

func containerSnapshotSummary(snapshot *proxmoxlib.ContainerSnapshot) *Snapshot {
	return &Snapshot{
		Kind:        "container",
		Node:        snapshot.Node,
		VMID:        snapshot.VMID,
		Name:        snapshot.Name,
		Description: snapshot.Description,
		Parent:      snapshot.Parent,
		SnapTimeSec: snapshot.SnapshotCreationTime,
	}
}

func storageSummary(storage *proxmoxlib.Storage) *Storage {
	return &Storage{
		Node:         storage.Node,
		Name:         storage.Name,
		Type:         storage.Type,
		Content:      storage.Content,
		Active:       storage.Active == 1,
		Enabled:      storage.Enabled == 1,
		Shared:       storage.Shared == 1,
		UsedFraction: storage.UsedFraction,
		Avail:        storage.Avail,
		Used:         storage.Used,
		Total:        storage.Total,
	}
}

func networkSummary(network *proxmoxlib.NodeNetwork) *Network {
	return &Network{
		Node:            network.Node,
		Iface:           network.Iface,
		Type:            network.Type,
		Active:          int(network.Active) == 1,
		Autostart:       network.Autostart == 1,
		Address:         network.Address,
		Address6:        network.Address6,
		CIDR:            network.CIDR,
		CIDR6:           network.CIDR6,
		Gateway:         network.Gateway,
		Gateway6:        network.Gateway6,
		BridgePorts:     network.BridgePorts,
		BridgeVLANAware: network.BridgeVLANAware == 1,
		VLANID:          network.VLANID,
		VLANRawDevice:   network.VLANRawDevice,
		MTU:             network.MTU,
		Method:          network.Method,
		Method6:         network.Method6,
		Comments:        network.Comments,
	}
}

func clusterResourceSummary(resource *proxmoxlib.ClusterResource) *ClusterResource {
	return &ClusterResource{
		ID:         resource.ID,
		Type:       resource.Type,
		Node:       resource.Node,
		VMID:       int(resource.VMID),
		Name:       resource.Name,
		Status:     resource.Status,
		Pool:       resource.Pool,
		Content:    resource.Content,
		Storage:    resource.Storage,
		CPU:        resource.CPU,
		MaxCPU:     int(resource.MaxCPU),
		Mem:        resource.Mem,
		MaxMem:     resource.MaxMem,
		Disk:       resource.Disk,
		MaxDisk:    resource.MaxDisk,
		Uptime:     resource.Uptime,
		Template:   resource.Template == 1,
		Tags:       resource.Tags,
		PluginType: resource.PluginType,
	}
}

func taskSummary(task *proxmoxlib.Task) *Task {
	summary := &Task{
		UPID:         string(task.UPID),
		ID:           task.ID,
		Type:         task.Type,
		User:         task.User,
		Status:       task.Status,
		Node:         task.Node,
		PID:          task.PID,
		PStart:       task.PStart,
		Saved:        task.Saved,
		ExitStatus:   task.ExitStatus,
		IsCompleted:  task.IsCompleted,
		IsRunning:    task.IsRunning,
		IsFailed:     task.IsFailed,
		IsSuccessful: task.IsSuccessful,
	}

	if !task.StartTime.IsZero() {
		summary.StartTimeSec = task.StartTime.Unix()
	}
	if !task.EndTime.IsZero() {
		summary.EndTimeSec = task.EndTime.Unix()
	}
	if task.Duration > 0 {
		summary.DurationSec = int64(task.Duration / time.Second)
	}

	return summary
}

func taskLogLines(log proxmoxlib.Log) []string {
	lines := make([]int, 0, len(log))
	for line := range log {
		lines = append(lines, line)
	}
	sort.Ints(lines)

	result := make([]string, 0, len(lines))
	for _, line := range lines {
		result = append(result, log[line])
	}

	return result
}
