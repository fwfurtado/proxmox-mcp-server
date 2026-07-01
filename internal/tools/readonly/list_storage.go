package readonly

import (
	"context"
	"fmt"

	"github.com/fwfurtado/proxmox-mcp-server/internal/proxmox"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

type listStorageInput struct{}

type listStorageOutput struct {
	Storage []*proxmox.Storage `json:"storage" jsonschema:"Proxmox storage entries"`
	Count   int                `json:"count" jsonschema:"Number of storage entries returned"`
}

func registerListStorageTool(server *sdkmcp.Server, client *proxmox.Client) {
	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "list_storage",
		Description: "List Proxmox storage across all cluster nodes.",
		Annotations: readOnlyToolAnnotations(),
	}, func(ctx context.Context, _ *sdkmcp.CallToolRequest, _ listStorageInput) (*sdkmcp.CallToolResult, listStorageOutput, error) {
		storage, err := client.ListStorage(ctx)
		if err != nil {
			return nil, listStorageOutput{}, fmt.Errorf("list Proxmox storage: %w", err)
		}

		return nil, listStorageOutput{
			Storage: storage,
			Count:   len(storage),
		}, nil
	})
}
