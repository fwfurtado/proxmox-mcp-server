package readwrite

import (
	"context"

	"github.com/fwfurtado/proxmox-mcp-server/internal/proxmox"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerRebootVMTool(server *sdkmcp.Server, client *proxmox.Client) {
	registerVMPowerTool(server, "reboot_vm", "Reboot a Proxmox virtual machine.", func(ctx context.Context, nodeName string, vmid int) (*proxmox.Task, error) {
		return client.RebootVM(ctx, nodeName, vmid)
	})
}
