package readonly

import (
	"context"
	"fmt"

	"github.com/fwfurtado/proxmox-mcp-server/internal/proxmox"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

type getVMInput struct {
	NodeName string `json:"node_name" jsonschema:"Name of the Proxmox node hosting the virtual machine"`
	VMID     int    `json:"vmid" jsonschema:"Virtual machine ID"`
}

type getVMOutput struct {
	VM *proxmox.VMDetails `json:"vm" jsonschema:"Proxmox virtual machine details"`
}

func registerGetVMTool(server *sdkmcp.Server, client *proxmox.Client) {
	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "get_vm",
		Description: "Get details for a Proxmox virtual machine by node name and VMID.",
		Annotations: readOnlyToolAnnotations(),
	}, func(ctx context.Context, _ *sdkmcp.CallToolRequest, input getVMInput) (*sdkmcp.CallToolResult, getVMOutput, error) {
		vm, err := client.GetVM(ctx, input.NodeName, input.VMID)
		if err != nil {
			return nil, getVMOutput{}, fmt.Errorf("get Proxmox VM: %w", err)
		}

		return nil, getVMOutput{VM: vm}, nil
	})
}
