package tools

import (
	"context"
	"fmt"
	"os"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"google.golang.org/protobuf/types/known/structpb"
)

// DebugStopSchema returns the JSON Schema for the debug_stop tool.
func DebugStopSchema() *structpb.Struct {
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

// DebugStop returns a handler that stops a running debug session.
func DebugStop() func(context.Context, *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		if err := helpers.ValidateRequired(req.Arguments, "session_id"); err != nil {
			return helpers.ErrorResult("validation_error", err.Error()), nil
		}
		sessionID := helpers.GetString(req.Arguments, "session_id")

		pid, ok := sessions[sessionID]
		if !ok {
			return helpers.ErrorResult("not_found", fmt.Sprintf("no debug session found with ID %q", sessionID)), nil
		}

		proc, err := os.FindProcess(pid)
		if err != nil {
			delete(sessions, sessionID)
			return helpers.ErrorResult("process_error", fmt.Sprintf("failed to find process %d: %v", pid, err)), nil
		}

		if err := proc.Kill(); err != nil {
			delete(sessions, sessionID)
			return helpers.ErrorResult("kill_error", fmt.Sprintf("failed to kill process %d: %v", pid, err)), nil
		}

		delete(sessions, sessionID)
		return helpers.TextResult(fmt.Sprintf("Debug session %s stopped", sessionID)), nil
	}
}
