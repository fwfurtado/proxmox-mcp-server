package readonly

import sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"

func readOnlyToolAnnotations() *sdkmcp.ToolAnnotations {
	destructive := false
	openWorld := false

	return &sdkmcp.ToolAnnotations{
		ReadOnlyHint:    true,
		IdempotentHint:  true,
		DestructiveHint: &destructive,
		OpenWorldHint:   &openWorld,
	}
}
