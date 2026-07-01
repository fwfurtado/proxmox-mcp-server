package readonly

import (
	"log/slog"

	"github.com/fwfurtado/proxmox-mcp-server/internal/proxmox"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

func RegisterTools(server *sdkmcp.Server, client *proxmox.Client) {
	logger := slog.Default()

	logger.Info("registering MCP tool", "name", "list_nodes", "mode", "read-only")
	registerListNodesTool(server, client)

	logger.Info("registering MCP tool", "name", "list_vms", "mode", "read-only")
	registerListVMsTool(server, client)

	logger.Info("registering MCP tool", "name", "get_vm", "mode", "read-only")
	registerGetVMTool(server, client)
}
