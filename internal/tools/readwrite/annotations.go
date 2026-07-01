package readwrite

import sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"

func powerActionToolAnnotations() *sdkmcp.ToolAnnotations {
	destructive := true
	openWorld := false

	return &sdkmcp.ToolAnnotations{
		ReadOnlyHint:    false,
		IdempotentHint:  false,
		DestructiveHint: &destructive,
		OpenWorldHint:   &openWorld,
	}
}
