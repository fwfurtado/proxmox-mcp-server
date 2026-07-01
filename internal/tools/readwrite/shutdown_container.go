package readwrite

import (
	"context"

	"github.com/fwfurtado/proxmox-mcp-server/internal/proxmox"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

type shutdownContainerInput struct {
	NodeName string `json:"node_name" jsonschema:"Name of the Proxmox node hosting the container"`
	VMID     int    `json:"vmid" jsonschema:"Container ID"`
	Force    bool   `json:"force,omitempty" jsonschema:"Force stop if graceful shutdown fails"`
	Timeout  int    `json:"timeout,omitempty" jsonschema:"Shutdown timeout in seconds; defaults to Proxmox behavior when zero"`
}

func registerShutdownContainerTool(server *sdkmcp.Server, client *proxmox.Client) {
	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "shutdown_container",
		Description: "Gracefully shut down a Proxmox container.",
		Annotations: powerActionToolAnnotations(),
	}, func(ctx context.Context, _ *sdkmcp.CallToolRequest, input shutdownContainerInput) (*sdkmcp.CallToolResult, taskOutput, error) {
		task, err := client.ShutdownContainer(ctx, input.NodeName, input.VMID, input.Force, input.Timeout)
		if err != nil {
			return nil, taskOutput{}, err
		}

		return nil, taskOutput{Task: task}, nil
	})
}
