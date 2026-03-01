package tools

import (
	"context"
	"fmt"
	"sort"
	"strings"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"google.golang.org/protobuf/types/known/structpb"
)

// DebugListSessionsSchema returns the JSON Schema for the debug_list_sessions tool.
func DebugListSessionsSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type":       "object",
		"properties": map[string]any{},
	})
	return s
}

// DebugListSessions returns a handler that lists all active debug sessions.
func DebugListSessions() func(context.Context, *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		if len(sessions) == 0 {
			return helpers.TextResult("No active debug sessions."), nil
		}

		// Sort by session ID for consistent output.
		ids := make([]string, 0, len(sessions))
		for id := range sessions {
			ids = append(ids, id)
		}
		sort.Strings(ids)

		var b strings.Builder
		fmt.Fprintf(&b, "## Active Debug Sessions (%d)\n\n", len(sessions))
		fmt.Fprintf(&b, "| Session ID | PID |\n")
		fmt.Fprintf(&b, "|------------|-----|\n")
		for _, id := range ids {
			fmt.Fprintf(&b, "| %s | %d |\n", id, sessions[id])
		}
		return helpers.TextResult(b.String()), nil
	}
}
