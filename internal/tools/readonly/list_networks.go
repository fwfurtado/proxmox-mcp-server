package readonly

import (
	"context"
	"fmt"

	"github.com/fwfurtado/proxmox-mcp-server/internal/proxmox"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

type listNetworksInput struct {
	IfaceType string `json:"iface_type,omitempty" jsonschema:"Optional Proxmox network interface type filter"`
}

type listNetworksOutput struct {
	Networks []*proxmox.Network `json:"networks" jsonschema:"Proxmox node networks across the cluster"`
	Count    int                `json:"count" jsonschema:"Number of network entries returned"`
}

func registerListNetworksTool(server *sdkmcp.Server, client *proxmox.Client) {
	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "list_networks",
		Description: "List Proxmox network interfaces across all cluster nodes.",
		Annotations: readOnlyToolAnnotations(),
	}, func(ctx context.Context, _ *sdkmcp.CallToolRequest, input listNetworksInput) (*sdkmcp.CallToolResult, listNetworksOutput, error) {
		networks, err := client.ListNetworks(ctx, input.IfaceType)
		if err != nil {
			return nil, listNetworksOutput{}, fmt.Errorf("list Proxmox networks: %w", err)
		}

		return nil, listNetworksOutput{
			Networks: networks,
			Count:    len(networks),
		}, nil
	})
}
