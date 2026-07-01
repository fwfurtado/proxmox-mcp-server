/*
Copyright © 2026 Mimi

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"os"
	"strconv"

	"github.com/fwfurtado/proxmox-mcp-server/internal/app"
	"github.com/fwfurtado/proxmox-mcp-server/internal/proxmox"
	"github.com/spf13/cobra"
)

var readOnly bool
var transport string
var httpAddr string
var proxmoxConfig proxmox.Config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "proxmox-mcp-server",
	Short: "A proxmox mcp server",
	Long:  `A proxmox mcp server`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runServer(cmd)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&readOnly, "read-only", false, "start the MCP server with only read-only tools enabled")
	rootCmd.PersistentFlags().StringVar(&transport, "transport", app.TransportStdio, "MCP transport: stdio or streamable-http (env: MCP_TRANSPORT)")
	rootCmd.PersistentFlags().StringVar(&httpAddr, "http-addr", ":8080", "listen address for streamable-http transport (env: MCP_HTTP_ADDR)")
	rootCmd.PersistentFlags().StringVar(&proxmoxConfig.URL, "proxmox-url", "", "Proxmox API URL (env: PROXMOX_URL)")
	rootCmd.PersistentFlags().StringVar(&proxmoxConfig.TokenID, "proxmox-token-id", "", "Proxmox API token ID, for example root@pam!mcp (env: PROXMOX_TOKEN_ID)")
	rootCmd.PersistentFlags().StringVar(&proxmoxConfig.TokenSecret, "proxmox-token-secret", "", "Proxmox API token secret (env: PROXMOX_TOKEN_SECRET)")
	rootCmd.PersistentFlags().BoolVar(&proxmoxConfig.InsecureTLS, "proxmox-insecure-tls", false, "skip Proxmox TLS certificate verification (env: PROXMOX_INSECURE_TLS)")
}

func runServer(cmd *cobra.Command) error {
	return app.Run(cmd.Context(), app.Config{
		ReadOnly:  readOnly,
		Transport: stringFlagFromEnv(cmd, "transport", transport, "MCP_TRANSPORT"),
		HTTPAddr:  stringFlagFromEnv(cmd, "http-addr", httpAddr, "MCP_HTTP_ADDR"),
		Proxmox:   proxmoxConfigFromEnv(proxmoxConfig),
	})
}

func proxmoxConfigFromEnv(config proxmox.Config) proxmox.Config {
	if config.URL == "" {
		config.URL = os.Getenv("PROXMOX_URL")
	}

	if config.TokenID == "" {
		config.TokenID = os.Getenv("PROXMOX_TOKEN_ID")
	}
	if config.TokenSecret == "" {
		config.TokenSecret = os.Getenv("PROXMOX_TOKEN_SECRET")
	}
	if !config.InsecureTLS {
		config.InsecureTLS, _ = strconv.ParseBool(os.Getenv("PROXMOX_INSECURE_TLS"))
	}

	return config
}

func stringFlagFromEnv(cmd *cobra.Command, flagName, value, envName string) string {
	if cmd.Flags().Changed(flagName) {
		return value
	}
	if envValue := os.Getenv(envName); envValue != "" {
		return envValue
	}

	return value
}
