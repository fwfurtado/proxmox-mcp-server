package readonly

import (
	"log/slog"

	"github.com/fwfurtado/proxmox-mcp-server/internal/proxmox"
	"github.com/fwfurtado/proxmox-mcp-server/internal/tools"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

type toolRegistration struct {
	name     string
	register func(*sdkmcp.Server, *proxmox.Client)
}

func RegisterTools(server *sdkmcp.Server, client *proxmox.Client, allowlist []string) {
	logger := slog.Default()
	allowed := tools.NewAllowlist(allowlist)

	tools := []toolRegistration{
		{name: "list_nodes", register: registerListNodesTool},
		{name: "list_vms", register: registerListVMsTool},
		{name: "get_vm", register: registerGetVMTool},
		{name: "get_vm_config", register: registerGetVMConfigTool},
		{name: "list_containers", register: registerListContainersTool},
		{name: "get_container", register: registerGetContainerTool},
		{name: "list_storage", register: registerListStorageTool},
		{name: "list_tasks", register: registerListTasksTool},
		{name: "get_task", register: registerGetTaskTool},
		{name: "list_snapshots", register: registerListSnapshotsTool},
		{name: "list_networks", register: registerListNetworksTool},
		{name: "list_node_networks", register: registerListNodeNetworksTool},
		{name: "list_cluster_resources", register: registerListClusterResourcesTool},
	}

	for _, tool := range tools {
		if !allowed.Allows(tool.name) {
			logger.Info("skipping MCP tool", "name", tool.name, "mode", "read-only", "reason", "not in allowlist")
			continue
		}

		logger.Info("registering MCP tool", "name", tool.name, "mode", "read-only")
		tool.register(server, client)
	}
}
