package readonly

import (
	"context"
	"fmt"

	"github.com/fwfurtado/proxmox-mcp-server/internal/proxmox"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

type getTaskInput struct {
	UPID     string `json:"upid" jsonschema:"Proxmox task UPID"`
	LogStart int    `json:"log_start,omitempty" jsonschema:"Zero-based log line offset to start from"`
	LogLimit int    `json:"log_limit,omitempty" jsonschema:"Maximum number of task log lines to return; default 50"`
}

type getTaskOutput struct {
	Task *proxmox.TaskDetails `json:"task" jsonschema:"Proxmox task details with log lines"`
}

func registerGetTaskTool(server *sdkmcp.Server, client *proxmox.Client) {
	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "get_task",
		Description: "Get details and log lines for a Proxmox task by UPID.",
		Annotations: readOnlyToolAnnotations(),
	}, func(ctx context.Context, _ *sdkmcp.CallToolRequest, input getTaskInput) (*sdkmcp.CallToolResult, getTaskOutput, error) {
		task, err := client.GetTask(ctx, input.UPID, input.LogStart, input.LogLimit)
		if err != nil {
			return nil, getTaskOutput{}, fmt.Errorf("get Proxmox task: %w", err)
		}

		return nil, getTaskOutput{Task: task}, nil
	})
}
