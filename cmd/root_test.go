package cmd

import (
	"testing"

	"github.com/fwfurtado/proxmox-mcp-server/internal/proxmox"
	"github.com/spf13/cobra"
)

func TestBoolFlagFromEnv(t *testing.T) {
	t.Run("uses flag value when flag changed", func(t *testing.T) {
		cmd := &cobra.Command{}
		cmd.Flags().Bool("allow-write", false, "")
		if err := cmd.Flags().Set("allow-write", "true"); err != nil {
			t.Fatalf("set flag: %v", err)
		}
		t.Setenv("MCP_ALLOW_WRITE", "false")

		got := boolFlagFromEnv(cmd, "allow-write", true, "MCP_ALLOW_WRITE")
		if !got {
			t.Fatal("expected true from changed flag")
		}
	})

	t.Run("uses env value when flag unchanged", func(t *testing.T) {
		cmd := &cobra.Command{}
		cmd.Flags().Bool("allow-write", false, "")
		t.Setenv("MCP_ALLOW_WRITE", "true")

		got := boolFlagFromEnv(cmd, "allow-write", false, "MCP_ALLOW_WRITE")
		if !got {
			t.Fatal("expected true from env")
		}
	})

	t.Run("falls back to input value on invalid env", func(t *testing.T) {
		cmd := &cobra.Command{}
		cmd.Flags().Bool("allow-write", false, "")
		t.Setenv("MCP_ALLOW_WRITE", "invalid")

		got := boolFlagFromEnv(cmd, "allow-write", false, "MCP_ALLOW_WRITE")
		if got {
			t.Fatal("expected false fallback value")
		}
	})
}

func TestStringFlagFromEnv(t *testing.T) {
	t.Run("uses flag value when flag changed", func(t *testing.T) {
		cmd := &cobra.Command{}
		cmd.Flags().String("transport", "", "")
		if err := cmd.Flags().Set("transport", "streamable-http"); err != nil {
			t.Fatalf("set flag: %v", err)
		}
		t.Setenv("MCP_TRANSPORT", "stdio")

		got := stringFlagFromEnv(cmd, "transport", "streamable-http", "MCP_TRANSPORT")
		if got != "streamable-http" {
			t.Fatalf("expected flag value, got %q", got)
		}
	})

	t.Run("uses env value when flag unchanged", func(t *testing.T) {
		cmd := &cobra.Command{}
		cmd.Flags().String("transport", "", "")
		t.Setenv("MCP_TRANSPORT", "streamable-http")

		got := stringFlagFromEnv(cmd, "transport", "stdio", "MCP_TRANSPORT")
		if got != "streamable-http" {
			t.Fatalf("expected env value, got %q", got)
		}
	})
}

func TestProxmoxConfigFromEnv(t *testing.T) {
	t.Setenv("PROXMOX_URL", "https://pve.example.com")
	t.Setenv("PROXMOX_TOKEN_ID", "root@pam!mcp")
	t.Setenv("PROXMOX_TOKEN_SECRET", "secret")
	t.Setenv("PROXMOX_INSECURE_TLS", "true")

	got := proxmoxConfigFromEnv(proxmox.Config{})

	if got.URL != "https://pve.example.com" {
		t.Fatalf("expected URL from env, got %q", got.URL)
	}
	if got.TokenID != "root@pam!mcp" {
		t.Fatalf("expected TokenID from env, got %q", got.TokenID)
	}
	if got.TokenSecret != "secret" {
		t.Fatalf("expected TokenSecret from env, got %q", got.TokenSecret)
	}
	if !got.InsecureTLS {
		t.Fatal("expected InsecureTLS from env")
	}
}

func TestProxmoxConfigFromEnvPreservesExplicitValues(t *testing.T) {
	t.Setenv("PROXMOX_URL", "https://ignored.example.com")
	t.Setenv("PROXMOX_TOKEN_ID", "ignored")
	t.Setenv("PROXMOX_TOKEN_SECRET", "ignored")
	t.Setenv("PROXMOX_INSECURE_TLS", "true")

	input := proxmox.Config{
		URL:         "https://configured.example.com",
		TokenID:     "configured-id",
		TokenSecret: "configured-secret",
		InsecureTLS: false,
	}

	got := proxmoxConfigFromEnv(input)

	if got.URL != input.URL {
		t.Fatalf("expected explicit URL %q, got %q", input.URL, got.URL)
	}
	if got.TokenID != input.TokenID {
		t.Fatalf("expected explicit TokenID %q, got %q", input.TokenID, got.TokenID)
	}
	if got.TokenSecret != input.TokenSecret {
		t.Fatalf("expected explicit TokenSecret %q, got %q", input.TokenSecret, got.TokenSecret)
	}
	if !got.InsecureTLS {
		t.Fatal("expected InsecureTLS to be loaded from env when explicit value is false")
	}
}
