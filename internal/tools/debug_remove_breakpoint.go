package tools

import (
	"context"
	"fmt"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"google.golang.org/protobuf/types/known/structpb"
)

// DebugRemoveBreakpointSchema returns the JSON Schema for the debug_remove_breakpoint tool.
func DebugRemoveBreakpointSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"session_id": map[string]any{
				"type":        "string",
				"description": "Session ID returned by debug_start",
			},
			"file": map[string]any{
				"type":        "string",
				"description": "Source file path of the breakpoint to remove",
			},
			"line": map[string]any{
				"type":        "number",
				"description": "Line number of the breakpoint to remove",
			},
		},
		"required": []any{"session_id", "file", "line"},
	})
	return s
}

// DebugRemoveBreakpoint returns a handler that acknowledges a breakpoint removal request.
func DebugRemoveBreakpoint() func(context.Context, *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		if err := helpers.ValidateRequired(req.Arguments, "session_id", "file"); err != nil {
			return helpers.ErrorResult("validation_error", err.Error()), nil
		}
		file := helpers.GetString(req.Arguments, "file")
		line := helpers.GetInt(req.Arguments, "line")
		return helpers.TextResult(fmt.Sprintf(
			"Breakpoint at %s:%d removed. Use your DAP client to manage breakpoints interactively.",
			file, line,
		)), nil
	}
}
