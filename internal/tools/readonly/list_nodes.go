package readonly

import (
	"context"
	"fmt"

	"github.com/fwfurtado/proxmox-mcp-server/internal/proxmox"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

type listNodesInput struct{}

type listNodesOutput struct {
	Nodes []*proxmox.Node `json:"nodes" jsonschema:"Proxmox cluster nodes"`
	Count int             `json:"count" jsonschema:"Number of nodes returned"`
}

func registerListNodesTool(server *sdkmcp.Server, client *proxmox.Client) {
	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "list_nodes",
		Description: "List Proxmox cluster nodes with their basic status.",
		Annotations: readOnlyToolAnnotations(),
	}, func(ctx context.Context, _ *sdkmcp.CallToolRequest, _ listNodesInput) (*sdkmcp.CallToolResult, listNodesOutput, error) {
		nodes, err := client.ListNodes(ctx)
		if err != nil {
			return nil, listNodesOutput{}, fmt.Errorf("list Proxmox nodes: %w", err)
		}

		return nil, listNodesOutput{
			Nodes: nodes,
			Count: len(nodes),
		}, nil
	})
}
