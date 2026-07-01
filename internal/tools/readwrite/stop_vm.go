package readwrite

import (
	"context"

	"github.com/fwfurtado/proxmox-mcp-server/internal/proxmox"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerStopVMTool(server *sdkmcp.Server, client *proxmox.Client) {
	registerVMPowerTool(server, "stop_vm", "Force stop a Proxmox virtual machine.", func(ctx context.Context, nodeName string, vmid int) (*proxmox.Task, error) {
		return client.StopVM(ctx, nodeName, vmid)
	})
}
