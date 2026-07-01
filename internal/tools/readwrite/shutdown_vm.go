package readwrite

import (
	"context"

	"github.com/fwfurtado/proxmox-mcp-server/internal/proxmox"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerShutdownVMTool(server *sdkmcp.Server, client *proxmox.Client) {
	registerVMPowerTool(server, "shutdown_vm", "Gracefully shut down a Proxmox virtual machine.", func(ctx context.Context, nodeName string, vmid int) (*proxmox.Task, error) {
		return client.ShutdownVM(ctx, nodeName, vmid)
	})
}
