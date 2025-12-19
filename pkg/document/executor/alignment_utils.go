package executor

import (
	"strings"

	"github.com/adcondev/poster/pkg/constants"
	"github.com/adcondev/poster/pkg/service"
)

// ============================================================================
// Alignment Utilities
// ============================================================================

// AlignScope executes a function with the specified alignment, then resets to left.
// This eliminates repetitive align/reset patterns across handlers.
func AlignScope(printer *service.Printer, align string, fn func() error) error {
	if err := ApplyAlignment(printer, align); err != nil {
		return err
	}

	if err := fn(); err != nil {
		return err
	}

	return printer.AlignLeft()
}

// ApplyAlignment applies the specified alignment to the printer.
func ApplyAlignment(printer *service.Printer, align string) error {
	switch strings.ToLower(align) {
	case constants.Center.String():
		return printer.AlignCenter()
	case constants.Right.String():
		return printer.AlignRight()
	case constants.Left.String():
		return printer.AlignLeft()
	default:
		return printer.AlignLeft()
	}
}

// ResolveAlignment returns the alignment string, using default if empty.
func ResolveAlignment(align string, defaultAlign constants.Alignment) string {
	if align == "" {
		return defaultAlign.String()
	}
	return strings.ToLower(align)
}
