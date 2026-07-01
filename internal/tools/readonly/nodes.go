package readonly

import (
	"context"
	"fmt"

	"github.com/fwfurtado/proxmox-mcp-server/internal/proxmox"
	"github.com/fwfurtado/proxmox-mcp-server/internal/tools"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

type ListNodesTool struct {
	client *proxmox.Client
}

type listNodesInput struct{}

type listNodesOutput struct {
	Nodes []*proxmox.Node `json:"nodes" jsonschema:"Proxmox cluster nodes"`
	Count int             `json:"count" jsonschema:"Number of nodes returned"`
}

func NewListNodesTool(client *proxmox.Client) *ListNodesTool {
	return &ListNodesTool{client: client}
}

func (t *ListNodesTool) Name() string {
	return "list_nodes"
}

func (t *ListNodesTool) Description() string {
	return "List Proxmox cluster nodes with their basic status."
}

func (t *ListNodesTool) Mode() tools.Mode {
	return tools.ModeReadOnly
}

func (t *ListNodesTool) Register(server *sdkmcp.Server) error {
	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        t.Name(),
		Description: t.Description(),
	}, t.listNodes)

	return nil
}

func (t *ListNodesTool) listNodes(ctx context.Context, _ *sdkmcp.CallToolRequest, _ listNodesInput) (*sdkmcp.CallToolResult, listNodesOutput, error) {
	nodes, err := t.client.ListNodes(ctx)
	if err != nil {
		return nil, listNodesOutput{}, fmt.Errorf("list Proxmox nodes: %w", err)
	}

	return nil, listNodesOutput{
		Nodes: nodes,
		Count: len(nodes),
	}, nil
}
