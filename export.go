package devtoolsdebugger

import (
	"github.com/orchestra-mcp/plugin-devtools-debugger/internal"
	"github.com/orchestra-mcp/sdk-go/plugin"
)

// Register adds all debugger tools to the builder.
func Register(builder *plugin.PluginBuilder) {
	dp := &internal.DebuggerPlugin{}
	dp.RegisterTools(builder)
}
