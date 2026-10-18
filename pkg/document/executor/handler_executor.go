package executor

import (
	"encoding/json"
	"fmt"

	"github.com/adcondev/poster/internal/calculate"
	"github.com/adcondev/poster/pkg/commands/character"
	"github.com/adcondev/poster/pkg/constants"
	"github.com/adcondev/poster/pkg/document/schema"
	"github.com/adcondev/poster/pkg/profile"
	"github.com/adcondev/poster/pkg/service"
)

// Executor executes JSON ticket documents against a PrinterActions target.
type Executor struct {
	printer  service.PrinterActions
	registry *HandlerRegistry
	handlers map[string]CommandHandler
}

// CommandHandler processes a single document command payload.
type CommandHandler func(printer service.PrinterActions, data json.RawMessage) error

// NewExecutor creates a new Executor bound to the provided PrinterActions implementation.
func NewExecutor(printer service.PrinterActions) *Executor {
	reg := NewRegistry()
	e := &Executor{
		printer:  printer,
		registry: reg,
		handlers: reg.handlers,
	}

	// Register basic handlers
	e.registerHandler("text", e.handleText)
	e.registerHandler("feed", e.handleFeed)
	e.registerHandler("cut", e.handleCut)
	e.registerHandler("separator", e.handleSeparator)

	// Register hardware control handlers
	e.registerHandler("pulse", e.handlePulse)
	e.registerHandler("beep", e.handleBeep)

	// Register 2D & barcode handlers
	e.registerHandler("image", e.handleImage)
	e.registerHandler("qr", e.handleQR)
	e.registerHandler("barcode", e.handleBarcode)

	// Register layout & raw handlers
	e.registerHandler("table", e.handleTable)
	e.registerHandler("raw", e.handleRaw)

	return e
}

// RegisterHandler registers or overrides a command handler for a custom or built-in command type.
func (e *Executor) RegisterHandler(cmdType string, handler CommandHandler) {
	e.registerHandler(cmdType, handler)
}

// registerHandler registers a new command handler in the underlying registry.
func (e *Executor) registerHandler(cmdType string, handler CommandHandler) {
	if e.registry == nil {
		e.registry = NewRegistry()
		e.handlers = e.registry.handlers
	}
	e.registry.Register(cmdType, handler)
}

// Execute validates and executes a complete schema.Document.
func (e *Executor) Execute(doc *schema.Document) error {
	if e == nil || e.printer == nil {
		return fmt.Errorf("printer cannot be nil")
	}
	if doc == nil {
		return fmt.Errorf("document is nil")
	}
	if doc.Version != "" && doc.Version != constants.DefaultVersion {
		return fmt.Errorf("unsupported document version: %s (expected %s)", doc.Version, constants.DefaultVersion)
	}
	if err := doc.Validate(); err != nil {
		return fmt.Errorf("invalid document: %w", err)
	}

	if err := e.printer.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize printer: %w", err)
	}

	if err := e.applyProfileFromDocument(doc); err != nil {
		return fmt.Errorf("failed to apply profile settings: %w", err)
	}

	for i, cmd := range doc.Commands {
		handler, exists := e.registry.Get(cmd.Type)
		if !exists {
			return fmt.Errorf("unknown command type at position %d: %s", i, cmd.Type)
		}

		if err := handler(e.printer, cmd.Data); err != nil {
			return fmt.Errorf("command %d (%s) failed: %w", i, cmd.Type, err)
		}
	}

	return nil
}

// setCodeTable configures the character code table with fallback to Windows-1252.
func (e *Executor) setCodeTable(tableName string) error {
	table, ok := profile.ParseCodeTable(tableName)
	if !ok {
		return e.printer.SetCodeTable(character.WPC1252)
	}

	return e.printer.SetCodeTable(table)
}

// ExecuteJSON parses and executes a document from raw JSON bytes.
func (e *Executor) ExecuteJSON(data []byte) error {
	doc, err := schema.ParseDocument(data)
	if err != nil {
		return err
	}
	return e.Execute(doc)
}

// applyProfileFromDocument applies profile configuration from the parsed JSON document.
func (e *Executor) applyProfileFromDocument(doc *schema.Document) error {
	if doc == nil {
		return fmt.Errorf("document is nil")
	}
	if doc.Profile.Model == "" {
		return fmt.Errorf("profile model is required")
	}

	prof := e.printer.GetProfile()
	if prof == nil {
		return fmt.Errorf("printer profile is nil")
	}
	prof.Model = doc.Profile.Model

	if err := e.applyProfileOrDefaults(doc.Profile); err != nil {
		return fmt.Errorf("failed to apply profile settings: %w", err)
	}

	switch {
	case prof.PaperWidth == constants.Paper80mm && prof.DPI == constants.DefaultDPI:
		prof.DotsPerLine = constants.PaperPxWidth80mm
	case prof.PaperWidth == constants.Paper58mm && prof.DPI == constants.DefaultDPI:
		prof.DotsPerLine = constants.PaperPxWidth58mm
	default:
		prof.DotsPerLine = calculate.DotsPerLine(prof.PaperWidth-8, prof.DPI)
	}

	return nil
}

// applyProfileOrDefaults applies default profile settings when omitted in JSON.
func (e *Executor) applyProfileOrDefaults(config schema.ProfileConfig) error {
	prof := e.printer.GetProfile()
	if config.PaperWidth > 0 {
		prof.PaperWidth = float64(config.PaperWidth)
	} else {
		if preset, ok := profile.ProfileFromModel(config.Model); ok && preset.PaperWidth > 0 {
			prof.PaperWidth = preset.PaperWidth
		} else if prof.PaperWidth <= 0 {
			prof.PaperWidth = constants.Paper80mm
		}
	}

	if config.CodeTable == "" {
		if err := e.setCodeTable("WPC1252"); err != nil {
			return fmt.Errorf("failed to set default code table WPC1252: %w", err)
		}
	} else {
		if err := e.setCodeTable(config.CodeTable); err != nil {
			return fmt.Errorf("failed to set code table %s: %w", config.CodeTable, err)
		}
	}

	if config.DPI == 0 {
		prof.DPI = constants.DefaultDPI
	} else {
		prof.DPI = config.DPI
	}

	prof.HasQR = config.HasQR

	return nil
}
