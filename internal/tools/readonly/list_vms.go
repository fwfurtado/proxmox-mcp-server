package readonly

import (
	"context"
	"fmt"

	"github.com/fwfurtado/proxmox-mcp-server/internal/proxmox"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

type listVMsInput struct{}

type listVMsOutput struct {
	VMs   []*proxmox.VM `json:"vms" jsonschema:"Proxmox virtual machines"`
	Count int           `json:"count" jsonschema:"Number of virtual machines returned"`
}

func registerListVMsTool(server *sdkmcp.Server, client *proxmox.Client) {
	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "list_vms",
		Description: "List Proxmox virtual machines across all cluster nodes.",
		Annotations: readOnlyToolAnnotations(),
	}, func(ctx context.Context, _ *sdkmcp.CallToolRequest, _ listVMsInput) (*sdkmcp.CallToolResult, listVMsOutput, error) {
		vms, err := client.ListVMs(ctx)
		if err != nil {
			return nil, listVMsOutput{}, fmt.Errorf("list Proxmox VMs: %w", err)
		}

		return nil, listVMsOutput{
			VMs:   vms,
			Count: len(vms),
		}, nil
	})
}
