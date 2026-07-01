package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"os"

	"github.com/fwfurtado/proxmox-mcp-server/internal/proxmox"
	"github.com/fwfurtado/proxmox-mcp-server/internal/tools"
	"github.com/fwfurtado/proxmox-mcp-server/internal/tools/readonly"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	TransportStdio          = "stdio"
	TransportStreamableHTTP = "streamable-http"
)

type Config struct {
	AllowWrite bool
	Transport  string
	HTTPAddr   string
	Proxmox    proxmox.Config
	Logger     *slog.Logger
}

func Run(ctx context.Context, config Config) error {
	logger := config.Logger
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(os.Stderr, nil))
	}

	logger.Info("starting proxmox MCP server", "allow_write", config.AllowWrite, "transport", config.Transport)

	if config.AllowWrite {
		return fmt.Errorf("read-write mode is not implemented yet; omit --allow-write to run read-only")
	}
	mode := tools.ModeReadOnly

	if err := validateTransport(config.Transport); err != nil {
		return err
	}

	logger.Info("connecting to Proxmox", "url", safeURL(config.Proxmox.URL), "auth", "api-token")
	proxmoxClient, err := proxmox.NewClient(ctx, config.Proxmox)
	if err != nil {
		return err
	}
	logger.Info("connected to Proxmox")

	server := sdkmcp.NewServer(&sdkmcp.Implementation{
		Name:    "proxmox-mcp-server",
		Version: "0.1.0",
	}, nil)

	allTools := []tools.Tool{
		readonly.NewListNodesTool(proxmoxClient),
		readonly.NewListVMsTool(proxmoxClient),
		readonly.NewGetVMTool(proxmoxClient),
	}

	for _, tool := range tools.ForMode(mode, allTools) {
		logger.Info("registering MCP tool", "name", tool.Name(), "mode", tool.Mode())
		if err := tool.Register(server); err != nil {
			return fmt.Errorf("register tool %q: %w", tool.Name(), err)
		}
	}

	if err := runTransport(ctx, server, config, logger); err != nil {
		return err
	}

	logger.Info("MCP server stopped")
	return nil
}

func runTransport(ctx context.Context, server *sdkmcp.Server, config Config, logger *slog.Logger) error {
	switch config.Transport {
	case TransportStdio:
		return runStdioTransport(ctx, server, logger)
	case TransportStreamableHTTP:
		return runStreamableHTTPTransport(ctx, server, config.HTTPAddr, logger)
	default:
		return fmt.Errorf("unsupported transport %q: expected %q or %q", config.Transport, TransportStdio, TransportStreamableHTTP)
	}
}

func validateTransport(transport string) error {
	switch transport {
	case TransportStdio, TransportStreamableHTTP:
		return nil
	default:
		return fmt.Errorf("unsupported transport %q: expected %q or %q", transport, TransportStdio, TransportStreamableHTTP)
	}
}

func safeURL(rawURL string) string {
	parsedURL, err := url.Parse(rawURL)
	if err != nil || parsedURL == nil {
		return rawURL
	}

	parsedURL.User = nil
	return parsedURL.String()
}
