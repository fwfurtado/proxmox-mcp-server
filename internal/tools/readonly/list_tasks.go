package readonly

import (
	"context"
	"fmt"

	"github.com/fwfurtado/proxmox-mcp-server/internal/proxmox"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

type listTasksInput struct {
	NodeName string `json:"node_name,omitempty" jsonschema:"Optional Proxmox node name to restrict task listing"`
	Limit    int    `json:"limit,omitempty" jsonschema:"Optional maximum number of tasks to return for node-scoped listing"`
}

type listTasksOutput struct {
	Tasks []*proxmox.Task `json:"tasks" jsonschema:"Proxmox tasks"`
	Count int             `json:"count" jsonschema:"Number of tasks returned"`
}

func registerListTasksTool(server *sdkmcp.Server, client *proxmox.Client) {
	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "list_tasks",
		Description: "List Proxmox tasks across the cluster or for a specific node.",
		Annotations: readOnlyToolAnnotations(),
	}, func(ctx context.Context, _ *sdkmcp.CallToolRequest, input listTasksInput) (*sdkmcp.CallToolResult, listTasksOutput, error) {
		var (
			tasks []*proxmox.Task
			err   error
		)

		if input.NodeName != "" {
			tasks, err = client.ListNodeTasks(ctx, input.NodeName, input.Limit)
		} else {
			tasks, err = client.ListTasks(ctx)
		}
		if err != nil {
			return nil, listTasksOutput{}, fmt.Errorf("list Proxmox tasks: %w", err)
		}

		return nil, listTasksOutput{
			Tasks: tasks,
			Count: len(tasks),
		}, nil
	})
}
