package tools

import (
	"context"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"google.golang.org/protobuf/types/known/structpb"
)

// DebugEvaluateSchema returns the JSON Schema for the debug_evaluate tool.
func DebugEvaluateSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"session_id": map[string]any{
				"type":        "string",
				"description": "Session ID returned by debug_start",
			},
			"expression": map[string]any{
				"type":        "string",
				"description": "Expression to evaluate in the debug context",
			},
		},
		"required": []any{"session_id", "expression"},
	})
	return s
}

// DebugEvaluate returns a handler that informs the user evaluation requires an active DAP connection.
func DebugEvaluate() func(context.Context, *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		if err := helpers.ValidateRequired(req.Arguments, "session_id", "expression"); err != nil {
			return helpers.ErrorResult("validation_error", err.Error()), nil
		}
		return helpers.TextResult("Evaluation not supported without active DAP connection. Use your IDE debugger."), nil
	}
}
