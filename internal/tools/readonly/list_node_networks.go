package readonly

import (
	"context"
	"fmt"

	"github.com/fwfurtado/proxmox-mcp-server/internal/proxmox"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

type listNodeNetworksInput struct {
	NodeName  string `json:"node_name" jsonschema:"Name of the Proxmox node"`
	IfaceType string `json:"iface_type,omitempty" jsonschema:"Optional Proxmox network interface type filter"`
}

type listNodeNetworksOutput struct {
	Networks []*proxmox.Network `json:"networks" jsonschema:"Proxmox node networks"`
	Count    int                `json:"count" jsonschema:"Number of network entries returned"`
}

func registerListNodeNetworksTool(server *sdkmcp.Server, client *proxmox.Client) {
	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "list_node_networks",
		Description: "List Proxmox network interfaces for a specific node.",
		Annotations: readOnlyToolAnnotations(),
	}, func(ctx context.Context, _ *sdkmcp.CallToolRequest, input listNodeNetworksInput) (*sdkmcp.CallToolResult, listNodeNetworksOutput, error) {
		networks, err := client.ListNodeNetworks(ctx, input.NodeName, input.IfaceType)
		if err != nil {
			return nil, listNodeNetworksOutput{}, fmt.Errorf("list Proxmox node networks: %w", err)
		}

		return nil, listNodeNetworksOutput{
			Networks: networks,
			Count:    len(networks),
		}, nil
	})
}
