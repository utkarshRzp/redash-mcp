package redash

import (
	"fmt"
	"log/slog"

	"github.com/razorpay/mcp/redash/pkg/toolsets"
)

// NewToolSets creates a new ToolsetGroup with redash-specific toolsets
func NewToolSets(
	log *slog.Logger,
	enabledToolsets []string,
	readOnly bool,
) (*toolsets.ToolsetGroup, error) {
	group := toolsets.NewToolsetGroup(readOnly)

	// Add redash toolset
	redashToolset := NewRedashToolset(log, readOnly)
	group.AddToolset(redashToolset)

	// Enable specified toolsets or default to all
	if len(enabledToolsets) == 0 {
		enabledToolsets = []string{"redash"}
	}

	if err := group.EnableToolsets(enabledToolsets); err != nil {
		return nil, fmt.Errorf("failed to enable toolsets: %w", err)
	}

	return group, nil
}
