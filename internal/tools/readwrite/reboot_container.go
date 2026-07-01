package readwrite

import (
	"context"

	"github.com/fwfurtado/proxmox-mcp-server/internal/proxmox"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerRebootContainerTool(server *sdkmcp.Server, client *proxmox.Client) {
	registerContainerPowerTool(server, "reboot_container", "Reboot a Proxmox container.", func(ctx context.Context, nodeName string, vmid int) (*proxmox.Task, error) {
		return client.RebootContainer(ctx, nodeName, vmid)
	})
}
