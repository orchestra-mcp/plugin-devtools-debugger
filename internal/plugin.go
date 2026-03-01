package internal

import (
	"github.com/orchestra-mcp/plugin-devtools-debugger/internal/tools"
	"github.com/orchestra-mcp/sdk-go/plugin"
)

// DebuggerPlugin registers all debug session tools with the plugin builder.
type DebuggerPlugin struct{}

// RegisterTools registers all 9 debugger tools on the given plugin builder.
func (dp *DebuggerPlugin) RegisterTools(builder *plugin.PluginBuilder) {
	builder.RegisterTool("debug_start",
		"Start a debug session for a Go, Node, or Python program",
		tools.DebugStartSchema(), tools.DebugStart())

	builder.RegisterTool("debug_stop",
		"Stop a running debug session by session ID",
		tools.DebugStopSchema(), tools.DebugStop())

	builder.RegisterTool("debug_set_breakpoint",
		"Set a breakpoint at a file and line (use your DAP client to manage interactively)",
		tools.DebugSetBreakpointSchema(), tools.DebugSetBreakpoint())

	builder.RegisterTool("debug_remove_breakpoint",
		"Remove a breakpoint at a file and line",
		tools.DebugRemoveBreakpointSchema(), tools.DebugRemoveBreakpoint())

	builder.RegisterTool("debug_continue",
		"Continue execution in a paused debug session",
		tools.DebugContinueSchema(), tools.DebugContinue())

	builder.RegisterTool("debug_step_over",
		"Step over the current line in a debug session",
		tools.DebugStepOverSchema(), tools.DebugStepOver())

	builder.RegisterTool("debug_step_into",
		"Step into the current function call in a debug session",
		tools.DebugStepIntoSchema(), tools.DebugStepInto())

	builder.RegisterTool("debug_evaluate",
		"Evaluate an expression in the context of a debug session",
		tools.DebugEvaluateSchema(), tools.DebugEvaluate())

	builder.RegisterTool("debug_list_sessions",
		"List all active debug sessions and their PIDs",
		tools.DebugListSessionsSchema(), tools.DebugListSessions())
}
