package readwrite

import (
	"context"

	"github.com/fwfurtado/proxmox-mcp-server/internal/proxmox"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerResetVMTool(server *sdkmcp.Server, client *proxmox.Client) {
	registerVMPowerTool(server, "reset_vm", "Hard reset a Proxmox virtual machine.", func(ctx context.Context, nodeName string, vmid int) (*proxmox.Task, error) {
		return client.ResetVM(ctx, nodeName, vmid)
	})
}
