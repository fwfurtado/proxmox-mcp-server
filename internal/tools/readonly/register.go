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

	logger.Info("registering MCP tool", "name", "get_vm_config", "mode", "read-only")
	registerGetVMConfigTool(server, client)

	logger.Info("registering MCP tool", "name", "list_containers", "mode", "read-only")
	registerListContainersTool(server, client)

	logger.Info("registering MCP tool", "name", "get_container", "mode", "read-only")
	registerGetContainerTool(server, client)

	logger.Info("registering MCP tool", "name", "list_storage", "mode", "read-only")
	registerListStorageTool(server, client)

	logger.Info("registering MCP tool", "name", "list_tasks", "mode", "read-only")
	registerListTasksTool(server, client)

	logger.Info("registering MCP tool", "name", "get_task", "mode", "read-only")
	registerGetTaskTool(server, client)

	logger.Info("registering MCP tool", "name", "list_snapshots", "mode", "read-only")
	registerListSnapshotsTool(server, client)

	logger.Info("registering MCP tool", "name", "list_networks", "mode", "read-only")
	registerListNetworksTool(server, client)

	logger.Info("registering MCP tool", "name", "list_node_networks", "mode", "read-only")
	registerListNodeNetworksTool(server, client)

	logger.Info("registering MCP tool", "name", "list_cluster_resources", "mode", "read-only")
	registerListClusterResourcesTool(server, client)
}
