package readwrite

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
		{name: "start_vm", register: registerStartVMTool},
		{name: "stop_vm", register: registerStopVMTool},
		{name: "shutdown_vm", register: registerShutdownVMTool},
		{name: "reboot_vm", register: registerRebootVMTool},
		{name: "reset_vm", register: registerResetVMTool},
		{name: "start_container", register: registerStartContainerTool},
		{name: "stop_container", register: registerStopContainerTool},
		{name: "shutdown_container", register: registerShutdownContainerTool},
		{name: "reboot_container", register: registerRebootContainerTool},
	}

	for _, tool := range tools {
		if !allowed.Allows(tool.name) {
			logger.Info("skipping MCP tool", "name", tool.name, "mode", "read-write", "reason", "not in allowlist")
			continue
		}

		logger.Info("registering MCP tool", "name", tool.name, "mode", "read-write")
		tool.register(server, client)
	}
}
