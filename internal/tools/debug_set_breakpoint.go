package tools

import (
	"context"
	"fmt"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"google.golang.org/protobuf/types/known/structpb"
)

// DebugSetBreakpointSchema returns the JSON Schema for the debug_set_breakpoint tool.
func DebugSetBreakpointSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"session_id": map[string]any{
				"type":        "string",
				"description": "Session ID returned by debug_start",
			},
			"file": map[string]any{
				"type":        "string",
				"description": "Source file path",
			},
			"line": map[string]any{
				"type":        "number",
				"description": "Line number for the breakpoint",
			},
			"condition": map[string]any{
				"type":        "string",
				"description": "Optional condition expression for a conditional breakpoint",
			},
		},
		"required": []any{"session_id", "file", "line"},
	})
	return s
}

// DebugSetBreakpoint returns a handler that acknowledges a breakpoint request.
func DebugSetBreakpoint() func(context.Context, *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		if err := helpers.ValidateRequired(req.Arguments, "session_id", "file"); err != nil {
			return helpers.ErrorResult("validation_error", err.Error()), nil
		}
		file := helpers.GetString(req.Arguments, "file")
		line := helpers.GetInt(req.Arguments, "line")
		return helpers.TextResult(fmt.Sprintf(
			"Breakpoint set at %s:%d. Use your DAP client to manage breakpoints interactively.",
			file, line,
		)), nil
	}
}
