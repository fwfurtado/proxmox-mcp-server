package tools

import sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"

type Mode string

const (
	ModeReadWrite Mode = "read-write"
	ModeReadOnly  Mode = "read-only"
)

type Tool interface {
	Name() string
	Description() string
	Mode() Mode
	Register(server *sdkmcp.Server) error
}

func ForMode(mode Mode, all []Tool) []Tool {
	selected := make([]Tool, 0, len(all))

	for _, tool := range all {
		if mode == ModeReadOnly && tool.Mode() != ModeReadOnly {
			continue
		}

		selected = append(selected, tool)
	}

	return selected
}
