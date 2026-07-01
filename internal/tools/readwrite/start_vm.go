package readwrite

import (
	"context"

	"github.com/fwfurtado/proxmox-mcp-server/internal/proxmox"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerStartVMTool(server *sdkmcp.Server, client *proxmox.Client) {
	registerVMPowerTool(server, "start_vm", "Start a Proxmox virtual machine.", func(ctx context.Context, nodeName string, vmid int) (*proxmox.Task, error) {
		return client.StartVM(ctx, nodeName, vmid)
	})
}
