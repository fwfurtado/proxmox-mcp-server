package readonly

import (
	"context"
	"fmt"

	"github.com/fwfurtado/proxmox-mcp-server/internal/proxmox"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

type getContainerInput struct {
	NodeName string `json:"node_name" jsonschema:"Name of the Proxmox node hosting the container"`
	VMID     int    `json:"vmid" jsonschema:"Container ID"`
}

type getContainerOutput struct {
	Container *proxmox.ContainerDetails `json:"container" jsonschema:"Proxmox container details"`
}

func registerGetContainerTool(server *sdkmcp.Server, client *proxmox.Client) {
	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "get_container",
		Description: "Get details for a Proxmox container by node name and VMID.",
		Annotations: readOnlyToolAnnotations(),
	}, func(ctx context.Context, _ *sdkmcp.CallToolRequest, input getContainerInput) (*sdkmcp.CallToolResult, getContainerOutput, error) {
		container, err := client.GetContainer(ctx, input.NodeName, input.VMID)
		if err != nil {
			return nil, getContainerOutput{}, fmt.Errorf("get Proxmox container: %w", err)
		}

		return nil, getContainerOutput{Container: container}, nil
	})
}
