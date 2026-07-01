package readwrite

import (
	"context"

	"github.com/fwfurtado/proxmox-mcp-server/internal/proxmox"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

type vmPowerInput struct {
	NodeName string `json:"node_name" jsonschema:"Name of the Proxmox node hosting the virtual machine"`
	VMID     int    `json:"vmid" jsonschema:"Virtual machine ID"`
}

type containerPowerInput struct {
	NodeName string `json:"node_name" jsonschema:"Name of the Proxmox node hosting the container"`
	VMID     int    `json:"vmid" jsonschema:"Container ID"`
}

type taskOutput struct {
	Task *proxmox.Task `json:"task" jsonschema:"Proxmox task started for this action"`
}

type vmActionFunc func(context.Context, string, int) (*proxmox.Task, error)
type containerActionFunc func(context.Context, string, int) (*proxmox.Task, error)

func registerVMPowerTool(server *sdkmcp.Server, name, description string, action vmActionFunc) {
	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        name,
		Description: description,
		Annotations: powerActionToolAnnotations(),
	}, func(ctx context.Context, _ *sdkmcp.CallToolRequest, input vmPowerInput) (*sdkmcp.CallToolResult, taskOutput, error) {
		task, err := action(ctx, input.NodeName, input.VMID)
		if err != nil {
			return nil, taskOutput{}, err
		}

		return nil, taskOutput{Task: task}, nil
	})
}

func registerContainerPowerTool(server *sdkmcp.Server, name, description string, action containerActionFunc) {
	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        name,
		Description: description,
		Annotations: powerActionToolAnnotations(),
	}, func(ctx context.Context, _ *sdkmcp.CallToolRequest, input containerPowerInput) (*sdkmcp.CallToolResult, taskOutput, error) {
		task, err := action(ctx, input.NodeName, input.VMID)
		if err != nil {
			return nil, taskOutput{}, err
		}

		return nil, taskOutput{Task: task}, nil
	})
}
