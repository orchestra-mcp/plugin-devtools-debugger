package tools

import (
	"context"
	"os/exec"
	"strings"
	"testing"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"google.golang.org/protobuf/types/known/structpb"
)

// debuggerAvailable reports whether the binary required for the given runtime is on PATH.
func debuggerAvailable(runtime string) bool {
	bin := map[string]string{"go": "dlv", "node": "node", "python": "python"}[runtime]
	_, err := exec.LookPath(bin)
	return err == nil
}

// callTool invokes a tool handler with the given key/value argument pairs and returns the response.
func callTool(t *testing.T, handler func(context.Context, *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error), args map[string]any) *pluginv1.ToolResponse {
	t.Helper()
	s, err := structpb.NewStruct(args)
	if err != nil {
		t.Fatalf("callTool: structpb.NewStruct: %v", err)
	}
	resp, err := handler(context.Background(), &pluginv1.ToolRequest{Arguments: s})
	if err != nil {
		t.Fatalf("callTool: handler returned unexpected error: %v", err)
	}
	return resp
}

// isError returns true when the response represents a tool-level error.
func isError(resp *pluginv1.ToolResponse) bool {
	return resp != nil && !resp.Success
}

// errorCode extracts the error_code string from a failed response.
func errorCode(resp *pluginv1.ToolResponse) string {
	return resp.GetErrorCode()
}

// responseText extracts the plain text from a successful response.
// TextResult stores the text under the "text" key in Result.
func responseText(resp *pluginv1.ToolResponse) string {
	if resp == nil || resp.Result == nil {
		return ""
	}
	v, ok := resp.Result.Fields["text"]
	if !ok {
		return ""
	}
	return v.GetStringValue()
}

// ---------------------------------------------------------------------------
// debug_start
// ---------------------------------------------------------------------------

func TestDebugStart_MissingProgram(t *testing.T) {
	handler := DebugStart()
	resp := callTool(t, handler, map[string]any{})
	if !isError(resp) {
		t.Fatal("expected error response for missing program")
	}
	if errorCode(resp) != "validation_error" {
		t.Fatalf("expected error_code=validation_error, got %q", errorCode(resp))
	}
}

func TestDebugStart_GoNoDlv(t *testing.T) {
	if !debuggerAvailable("go") {
		t.Skip("dlv not found on PATH — skipping")
	}
	handler := DebugStart()
	resp := callTool(t, handler, map[string]any{
		"program": "/nonexistent/main.go",
		"runtime": "go",
	})
	if !isError(resp) {
		t.Fatal("expected error response when program path does not exist")
	}
	if errorCode(resp) != "start_error" {
		t.Fatalf("expected error_code=start_error, got %q", errorCode(resp))
	}
}

// ---------------------------------------------------------------------------
// debug_stop
// ---------------------------------------------------------------------------

func TestDebugStop_MissingSessionID(t *testing.T) {
	handler := DebugStop()
	resp := callTool(t, handler, map[string]any{})
	if !isError(resp) {
		t.Fatal("expected error response for missing session_id")
	}
	if errorCode(resp) != "validation_error" {
		t.Fatalf("expected error_code=validation_error, got %q", errorCode(resp))
	}
}

func TestDebugStop_UnknownSession(t *testing.T) {
	handler := DebugStop()
	resp := callTool(t, handler, map[string]any{
		"session_id": "dbg-unknown",
	})
	// debug_stop returns not_found for an unknown session ID.
	if !isError(resp) {
		t.Fatal("expected error response for unknown session")
	}
	if errorCode(resp) != "not_found" {
		t.Fatalf("expected error_code=not_found, got %q", errorCode(resp))
	}
}

// ---------------------------------------------------------------------------
// debug_set_breakpoint
// ---------------------------------------------------------------------------

func TestDebugSetBreakpoint_MissingArgs(t *testing.T) {
	handler := DebugSetBreakpoint()
	resp := callTool(t, handler, map[string]any{})
	if !isError(resp) {
		t.Fatal("expected error response for missing session_id")
	}
	if errorCode(resp) != "validation_error" {
		t.Fatalf("expected error_code=validation_error, got %q", errorCode(resp))
	}
}

func TestDebugSetBreakpoint_Valid(t *testing.T) {
	handler := DebugSetBreakpoint()
	resp := callTool(t, handler, map[string]any{
		"session_id": "dbg-1",
		"file":       "/src/main.go",
		"line":       float64(42),
	})
	if isError(resp) {
		t.Fatalf("expected success, got error_code=%q", errorCode(resp))
	}
	text := responseText(resp)
	if !strings.Contains(text, "main.go:42") {
		t.Fatalf("expected response to contain %q, got: %s", "main.go:42", text)
	}
}

// ---------------------------------------------------------------------------
// debug_remove_breakpoint
// ---------------------------------------------------------------------------

func TestDebugRemoveBreakpoint_MissingArgs(t *testing.T) {
	handler := DebugRemoveBreakpoint()
	resp := callTool(t, handler, map[string]any{})
	if !isError(resp) {
		t.Fatal("expected error response for missing session_id")
	}
	if errorCode(resp) != "validation_error" {
		t.Fatalf("expected error_code=validation_error, got %q", errorCode(resp))
	}
}

func TestDebugRemoveBreakpoint_Valid(t *testing.T) {
	handler := DebugRemoveBreakpoint()
	resp := callTool(t, handler, map[string]any{
		"session_id": "dbg-1",
		"file":       "/src/main.go",
		"line":       float64(42),
	})
	if isError(resp) {
		t.Fatalf("expected success, got error_code=%q", errorCode(resp))
	}
	text := responseText(resp)
	if !strings.Contains(text, "main.go:42") {
		t.Fatalf("expected response to contain %q, got: %s", "main.go:42", text)
	}
}

// ---------------------------------------------------------------------------
// debug_continue
// ---------------------------------------------------------------------------

func TestDebugContinue_MissingSession(t *testing.T) {
	handler := DebugContinue()
	resp := callTool(t, handler, map[string]any{})
	if !isError(resp) {
		t.Fatal("expected error response for missing session_id")
	}
	if errorCode(resp) != "validation_error" {
		t.Fatalf("expected error_code=validation_error, got %q", errorCode(resp))
	}
}

// ---------------------------------------------------------------------------
// debug_step_over
// ---------------------------------------------------------------------------

func TestDebugStepOver_MissingSession(t *testing.T) {
	handler := DebugStepOver()
	resp := callTool(t, handler, map[string]any{})
	if !isError(resp) {
		t.Fatal("expected error response for missing session_id")
	}
	if errorCode(resp) != "validation_error" {
		t.Fatalf("expected error_code=validation_error, got %q", errorCode(resp))
	}
}

// ---------------------------------------------------------------------------
// debug_step_into
// ---------------------------------------------------------------------------

func TestDebugStepInto_MissingSession(t *testing.T) {
	handler := DebugStepInto()
	resp := callTool(t, handler, map[string]any{})
	if !isError(resp) {
		t.Fatal("expected error response for missing session_id")
	}
	if errorCode(resp) != "validation_error" {
		t.Fatalf("expected error_code=validation_error, got %q", errorCode(resp))
	}
}

// ---------------------------------------------------------------------------
// debug_evaluate
// ---------------------------------------------------------------------------

func TestDebugEvaluate_MissingArgs(t *testing.T) {
	handler := DebugEvaluate()
	// Missing both session_id and expression.
	resp := callTool(t, handler, map[string]any{})
	if !isError(resp) {
		t.Fatal("expected error response for missing session_id/expression")
	}
	if errorCode(resp) != "validation_error" {
		t.Fatalf("expected error_code=validation_error, got %q", errorCode(resp))
	}
}

// ---------------------------------------------------------------------------
// debug_list_sessions
// ---------------------------------------------------------------------------

func TestDebugListSessions_NoArgs(t *testing.T) {
	// Ensure the package-level sessions map is empty for this test.
	for k := range sessions {
		delete(sessions, k)
	}

	handler := DebugListSessions()
	resp := callTool(t, handler, map[string]any{})
	if isError(resp) {
		t.Fatalf("expected success for empty sessions, got error_code=%q", errorCode(resp))
	}
	text := responseText(resp)
	if text == "" {
		t.Fatal("expected non-empty text response for debug_list_sessions with no sessions")
	}
}
