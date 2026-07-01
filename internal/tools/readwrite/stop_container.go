package readwrite

import (
	"context"

	"github.com/fwfurtado/proxmox-mcp-server/internal/proxmox"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerStopContainerTool(server *sdkmcp.Server, client *proxmox.Client) {
	registerContainerPowerTool(server, "stop_container", "Force stop a Proxmox container.", func(ctx context.Context, nodeName string, vmid int) (*proxmox.Task, error) {
		return client.StopContainer(ctx, nodeName, vmid)
	})
}
