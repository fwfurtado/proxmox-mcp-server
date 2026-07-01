package readwrite

import (
	"context"

	"github.com/fwfurtado/proxmox-mcp-server/internal/proxmox"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerStartContainerTool(server *sdkmcp.Server, client *proxmox.Client) {
	registerContainerPowerTool(server, "start_container", "Start a Proxmox container.", func(ctx context.Context, nodeName string, vmid int) (*proxmox.Task, error) {
		return client.StartContainer(ctx, nodeName, vmid)
	})
}
