package readonly

import (
	"context"
	"fmt"

	"github.com/fwfurtado/proxmox-mcp-server/internal/proxmox"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

type listSnapshotsInput struct {
	Kind     string `json:"kind,omitempty" jsonschema:"Snapshot kind: vm or container; defaults to vm"`
	NodeName string `json:"node_name" jsonschema:"Name of the Proxmox node hosting the resource"`
	VMID     int    `json:"vmid" jsonschema:"Virtual machine or container ID"`
}

type listSnapshotsOutput struct {
	Snapshots []*proxmox.Snapshot `json:"snapshots" jsonschema:"Proxmox snapshots"`
	Count     int                 `json:"count" jsonschema:"Number of snapshots returned"`
}

func registerListSnapshotsTool(server *sdkmcp.Server, client *proxmox.Client) {
	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "list_snapshots",
		Description: "List Proxmox snapshots for a virtual machine or container.",
		Annotations: readOnlyToolAnnotations(),
	}, func(ctx context.Context, _ *sdkmcp.CallToolRequest, input listSnapshotsInput) (*sdkmcp.CallToolResult, listSnapshotsOutput, error) {
		snapshots, err := client.ListSnapshots(ctx, input.Kind, input.NodeName, input.VMID)
		if err != nil {
			return nil, listSnapshotsOutput{}, fmt.Errorf("list Proxmox snapshots: %w", err)
		}

		return nil, listSnapshotsOutput{
			Snapshots: snapshots,
			Count:     len(snapshots),
		}, nil
	})
}
