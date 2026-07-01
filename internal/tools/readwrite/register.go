package readwrite

import (
	"log/slog"

	"github.com/fwfurtado/proxmox-mcp-server/internal/proxmox"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

func RegisterTools(server *sdkmcp.Server, client *proxmox.Client) {
	logger := slog.Default()

	logger.Info("registering MCP tool", "name", "start_vm", "mode", "read-write")
	registerStartVMTool(server, client)
	logger.Info("registering MCP tool", "name", "stop_vm", "mode", "read-write")
	registerStopVMTool(server, client)
	logger.Info("registering MCP tool", "name", "shutdown_vm", "mode", "read-write")
	registerShutdownVMTool(server, client)
	logger.Info("registering MCP tool", "name", "reboot_vm", "mode", "read-write")
	registerRebootVMTool(server, client)
	logger.Info("registering MCP tool", "name", "reset_vm", "mode", "read-write")
	registerResetVMTool(server, client)

	logger.Info("registering MCP tool", "name", "start_container", "mode", "read-write")
	registerStartContainerTool(server, client)
	logger.Info("registering MCP tool", "name", "stop_container", "mode", "read-write")
	registerStopContainerTool(server, client)
	logger.Info("registering MCP tool", "name", "shutdown_container", "mode", "read-write")
	registerShutdownContainerTool(server, client)
	logger.Info("registering MCP tool", "name", "reboot_container", "mode", "read-write")
	registerRebootContainerTool(server, client)
}
