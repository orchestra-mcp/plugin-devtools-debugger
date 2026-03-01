// Command devtools-debugger is the entry point for the devtools.debugger
// plugin binary. It provides 8 MCP tools for managing debug sessions.
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/orchestra-mcp/plugin-devtools-debugger/internal"
	"github.com/orchestra-mcp/sdk-go/plugin"
)

func main() {
	builder := plugin.New("devtools.debugger").
		Version("0.1.0").
		Description("Debug session manager for Go, Node, and Python programs").
		Author("Orchestra").
		Binary("devtools-debugger")

	tp := &internal.DebuggerPlugin{}
	tp.RegisterTools(builder)

	p := builder.BuildWithTools()
	p.ParseFlags()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		cancel()
	}()

	if err := p.Run(ctx); err != nil {
		log.Fatalf("devtools.debugger: %v", err)
	}
}
