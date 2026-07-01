package readonly

import (
	"context"
	"fmt"

	"github.com/fwfurtado/proxmox-mcp-server/internal/proxmox"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

type listClusterResourcesInput struct {
	Type string `json:"type,omitempty" jsonschema:"Optional Proxmox cluster resource type filter, for example vm, node, storage, or sdn"`
}

type listClusterResourcesOutput struct {
	Resources []*proxmox.ClusterResource `json:"resources" jsonschema:"Proxmox cluster resources"`
	Count     int                        `json:"count" jsonschema:"Number of cluster resources returned"`
}

func registerListClusterResourcesTool(server *sdkmcp.Server, client *proxmox.Client) {
	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "list_cluster_resources",
		Description: "List raw Proxmox cluster resources, optionally filtered by type.",
		Annotations: readOnlyToolAnnotations(),
	}, func(ctx context.Context, _ *sdkmcp.CallToolRequest, input listClusterResourcesInput) (*sdkmcp.CallToolResult, listClusterResourcesOutput, error) {
		resources, err := client.ListClusterResources(ctx, input.Type)
		if err != nil {
			return nil, listClusterResourcesOutput{}, fmt.Errorf("list Proxmox cluster resources: %w", err)
		}

		return nil, listClusterResourcesOutput{
			Resources: resources,
			Count:     len(resources),
		}, nil
	})
}
