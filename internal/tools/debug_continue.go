package tools

import (
	"context"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"google.golang.org/protobuf/types/known/structpb"
)

// DebugContinueSchema returns the JSON Schema for the debug_continue tool.
func DebugContinueSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"session_id": map[string]any{
				"type":        "string",
				"description": "Session ID returned by debug_start",
			},
		},
		"required": []any{"session_id"},
	})
	return s
}

// DebugContinue returns a handler that instructs the user to send continue via their DAP client.
func DebugContinue() func(context.Context, *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		if err := helpers.ValidateRequired(req.Arguments, "session_id"); err != nil {
			return helpers.ErrorResult("validation_error", err.Error()), nil
		}
		return helpers.TextResult("Send 'continue' command via your DAP client connected to the debug session."), nil
	}
}
