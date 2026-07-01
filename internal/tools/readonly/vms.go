package readonly

import (
	"context"
	"fmt"

	"github.com/fwfurtado/proxmox-mcp-server/internal/proxmox"
	"github.com/fwfurtado/proxmox-mcp-server/internal/tools"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

type ListVMsTool struct {
	client *proxmox.Client
}

type GetVMTool struct {
	client *proxmox.Client
}

type listVMsInput struct{}

type listVMsOutput struct {
	VMs   []*proxmox.VM `json:"vms" jsonschema:"Proxmox virtual machines"`
	Count int           `json:"count" jsonschema:"Number of virtual machines returned"`
}

type getVMInput struct {
	NodeName string `json:"node_name" jsonschema:"Name of the Proxmox node hosting the virtual machine"`
	VMID     int    `json:"vmid" jsonschema:"Virtual machine ID"`
}

type getVMOutput struct {
	VM *proxmox.VMDetails `json:"vm" jsonschema:"Proxmox virtual machine details"`
}

func NewListVMsTool(client *proxmox.Client) *ListVMsTool {
	return &ListVMsTool{client: client}
}

func (t *ListVMsTool) Name() string {
	return "list_vms"
}

func (t *ListVMsTool) Description() string {
	return "List Proxmox virtual machines across all cluster nodes."
}

func (t *ListVMsTool) Mode() tools.Mode {
	return tools.ModeReadOnly
}

func (t *ListVMsTool) Register(server *sdkmcp.Server) error {
	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        t.Name(),
		Description: t.Description(),
	}, t.listVMs)

	return nil
}

func (t *ListVMsTool) listVMs(ctx context.Context, _ *sdkmcp.CallToolRequest, _ listVMsInput) (*sdkmcp.CallToolResult, listVMsOutput, error) {
	vms, err := t.client.ListVMs(ctx)
	if err != nil {
		return nil, listVMsOutput{}, fmt.Errorf("list Proxmox VMs: %w", err)
	}

	return nil, listVMsOutput{
		VMs:   vms,
		Count: len(vms),
	}, nil
}

func NewGetVMTool(client *proxmox.Client) *GetVMTool {
	return &GetVMTool{client: client}
}

func (t *GetVMTool) Name() string {
	return "get_vm"
}

func (t *GetVMTool) Description() string {
	return "Get details for a Proxmox virtual machine by node name and VMID."
}

func (t *GetVMTool) Mode() tools.Mode {
	return tools.ModeReadOnly
}

func (t *GetVMTool) Register(server *sdkmcp.Server) error {
	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        t.Name(),
		Description: t.Description(),
	}, t.getVM)

	return nil
}

func (t *GetVMTool) getVM(ctx context.Context, _ *sdkmcp.CallToolRequest, input getVMInput) (*sdkmcp.CallToolResult, getVMOutput, error) {
	vm, err := t.client.GetVM(ctx, input.NodeName, input.VMID)
	if err != nil {
		return nil, getVMOutput{}, fmt.Errorf("get Proxmox VM: %w", err)
	}

	return nil, getVMOutput{VM: vm}, nil
}
