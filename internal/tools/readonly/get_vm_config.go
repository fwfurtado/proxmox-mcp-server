package readonly

import (
	"context"
	"fmt"

	"github.com/fwfurtado/proxmox-mcp-server/internal/proxmox"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

type getVMConfigInput struct {
	NodeName string `json:"node_name" jsonschema:"Name of the Proxmox node hosting the virtual machine"`
	VMID     int    `json:"vmid" jsonschema:"Virtual machine ID"`
}

type getVMConfigOutput struct {
	Config *proxmox.VMConfig `json:"config" jsonschema:"Proxmox virtual machine config"`
}

func registerGetVMConfigTool(server *sdkmcp.Server, client *proxmox.Client) {
	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "get_vm_config",
		Description: "Get config details for a Proxmox virtual machine by node name and VMID.",
		Annotations: readOnlyToolAnnotations(),
	}, func(ctx context.Context, _ *sdkmcp.CallToolRequest, input getVMConfigInput) (*sdkmcp.CallToolResult, getVMConfigOutput, error) {
		config, err := client.GetVMConfig(ctx, input.NodeName, input.VMID)
		if err != nil {
			return nil, getVMConfigOutput{}, fmt.Errorf("get Proxmox VM config: %w", err)
		}

		return nil, getVMConfigOutput{Config: config}, nil
	})
}
