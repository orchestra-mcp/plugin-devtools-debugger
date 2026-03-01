package tools

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"google.golang.org/protobuf/types/known/structpb"
)

// sessions tracks running debug processes: session_id -> PID.
var sessions = map[string]int{}

// DebugStartSchema returns the JSON Schema for the debug_start tool.
func DebugStartSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"program": map[string]any{
				"type":        "string",
				"description": "Path to the program to debug",
			},
			"runtime": map[string]any{
				"type":        "string",
				"description": "Runtime to use: go, node, or python (auto-detected from file extension if not given)",
				"enum":        []any{"go", "node", "python"},
			},
			"args": map[string]any{
				"type":        "string",
				"description": "Optional arguments to pass to the program",
			},
			"port": map[string]any{
				"type":        "number",
				"description": "Port for the debug adapter to listen on (default: 2345)",
			},
		},
		"required": []any{"program"},
	})
	return s
}

// DebugStart returns a handler that starts a debug session for a program.
func DebugStart() func(context.Context, *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		if err := helpers.ValidateRequired(req.Arguments, "program"); err != nil {
			return helpers.ErrorResult("validation_error", err.Error()), nil
		}
		program := helpers.GetString(req.Arguments, "program")
		runtime := helpers.GetString(req.Arguments, "runtime")
		args := helpers.GetString(req.Arguments, "args")
		port := helpers.GetInt(req.Arguments, "port")
		if port <= 0 {
			port = 2345
		}

		// Auto-detect runtime from extension if not provided.
		if runtime == "" {
			runtime = detectRuntime(program)
		}

		var cmd *exec.Cmd
		portStr := fmt.Sprintf("%d", port)
		switch runtime {
		case "node":
			cmdArgs := []string{"--inspect=" + portStr, program}
			if args != "" {
				cmdArgs = append(cmdArgs, strings.Fields(args)...)
			}
			cmd = exec.Command("node", cmdArgs...)
		case "python":
			cmdArgs := []string{"-m", "debugpy", "--listen", portStr, program}
			if args != "" {
				cmdArgs = append(cmdArgs, strings.Fields(args)...)
			}
			cmd = exec.Command("python", cmdArgs...)
		default: // go
			cmdArgs := []string{"debug", program, "--headless", "--api-version=2", "--listen=:" + portStr}
			if args != "" {
				cmdArgs = append(append(cmdArgs, "--"), strings.Fields(args)...)
			}
			cmd = exec.Command("dlv", cmdArgs...)
		}

		if err := cmd.Start(); err != nil {
			return helpers.ErrorResult("start_error", fmt.Sprintf("failed to start debug process: %v", err)), nil
		}

		pid := cmd.Process.Pid
		sessionID := fmt.Sprintf("dbg-%d", pid)
		sessions[sessionID] = pid

		// Detach from the process so it runs independently.
		go cmd.Wait() //nolint:errcheck

		return helpers.TextResult(fmt.Sprintf(
			"Debug session started. Connect debugger to :%d. PID: %d\nSession ID: %s",
			port, pid, sessionID,
		)), nil
	}
}

// detectRuntime infers the runtime from the program's file extension.
func detectRuntime(program string) string {
	switch filepath.Ext(program) {
	case ".js", ".mjs", ".cjs", ".ts":
		return "node"
	case ".py":
		return "python"
	default:
		return "go"
	}
}
