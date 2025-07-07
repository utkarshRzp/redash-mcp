package redash

import (
	"log/slog"

	"github.com/razorpay/mcp/redash/pkg/mcpgo"
	"github.com/razorpay/mcp/redash/pkg/toolsets"
)

// Server extends mcpgo.Server
type Server struct {
	log      *slog.Logger
	server   mcpgo.Server
	toolsets *toolsets.ToolsetGroup
}

// NewServer creates a new Server
func NewServer(
	log *slog.Logger,
	version string,
	enabledToolsets []string,
	readOnly bool,
) (*Server, error) {
	// Create default options
	opts := []mcpgo.ServerOption{
		mcpgo.WithLogging(),
		mcpgo.WithToolCapabilities(true),
	}

	// Create the mcpgo server
	server := mcpgo.NewServer(
		"redash-mcp-server",
		version,
		opts...,
	)

	//Initialize toolsets
	toolsets, err := NewToolSets(log, enabledToolsets, readOnly)
	if err != nil {
		return nil, err
	}

	// Create the server instance
	srv := &Server{
		log:      log,
		server:   server,
		toolsets: toolsets,
	}

	// Register all tools
	srv.RegisterTools()

	return srv, nil
}

// RegisterTools adds all available tools to the server
func (s *Server) RegisterTools() {
	s.toolsets.RegisterTools(s.server)
}

// GetMCPServer returns the underlying MCP server instance
func (s *Server) GetMCPServer() mcpgo.Server {
	return s.server
}
