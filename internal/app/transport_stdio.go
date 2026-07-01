package app

import (
	"context"
	"log/slog"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

func runStdioTransport(ctx context.Context, server *sdkmcp.Server, logger *slog.Logger) error {
	logger.Info("MCP server ready", "transport", TransportStdio)

	if err := server.Run(ctx, &sdkmcp.StdioTransport{}); err != nil {
		logger.Error("MCP stdio transport stopped with error", "error", err)
		return err
	}

	return nil
}
