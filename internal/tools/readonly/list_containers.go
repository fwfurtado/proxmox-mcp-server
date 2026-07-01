package readonly

import (
	"context"
	"fmt"

	"github.com/fwfurtado/proxmox-mcp-server/internal/proxmox"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

type listContainersInput struct{}

type listContainersOutput struct {
	Containers []*proxmox.Container `json:"containers" jsonschema:"Proxmox containers"`
	Count      int                  `json:"count" jsonschema:"Number of containers returned"`
}

func registerListContainersTool(server *sdkmcp.Server, client *proxmox.Client) {
	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "list_containers",
		Description: "List Proxmox containers across all cluster nodes.",
		Annotations: readOnlyToolAnnotations(),
	}, func(ctx context.Context, _ *sdkmcp.CallToolRequest, _ listContainersInput) (*sdkmcp.CallToolResult, listContainersOutput, error) {
		containers, err := client.ListContainers(ctx)
		if err != nil {
			return nil, listContainersOutput{}, fmt.Errorf("list Proxmox containers: %w", err)
		}

		return nil, listContainersOutput{
			Containers: containers,
			Count:      len(containers),
		}, nil
	})
}
