package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

func runStreamableHTTPTransport(ctx context.Context, server *sdkmcp.Server, addr string, logger *slog.Logger) error {
	handler := sdkmcp.NewStreamableHTTPHandler(func(*http.Request) *sdkmcp.Server {
		return server
	}, &sdkmcp.StreamableHTTPOptions{
		Logger: logger,
	})

	httpServer := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		logger.Info("MCP server ready", "transport", TransportStreamableHTTP, "addr", addr)
		errCh <- httpServer.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				logger.Warn("MCP streamable HTTP transport graceful shutdown timed out; forcing close", "error", err)
				if closeErr := httpServer.Close(); closeErr != nil && !errors.Is(closeErr, http.ErrServerClosed) {
					logger.Error("MCP streamable HTTP transport forced close failed", "error", closeErr)
					return closeErr
				}

				return nil
			}

			logger.Error("MCP streamable HTTP transport shutdown failed", "error", err)
			return err
		}

		return nil
	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}

		logger.Error("MCP streamable HTTP transport stopped with error", "error", err)
		return err
	}
}
