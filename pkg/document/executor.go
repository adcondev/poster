package document

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/adcondev/pos-printer/internal/calculate"
	"github.com/adcondev/pos-printer/pkg/commands/character"
	"github.com/adcondev/pos-printer/pkg/service"
)

// Executor ejecuta documentos de impresión
type Executor struct {
	printer  *service.Printer
	handlers map[string]CommandHandler
}

// CommandHandler a command handler function
type CommandHandler func(printer *service.Printer, data json.RawMessage) error

// NewExecutor crea un nuevo ejecutor
func NewExecutor(printer *service.Printer) *Executor {
	if printer == nil {
		log.Panicf("printer cannot be nil")
	}

	e := &Executor{
		printer:  printer,
		handlers: make(map[string]CommandHandler),
	}

	// Registrar handlers básicos
	e.RegisterHandler("text", e.handleText)
	e.RegisterHandler("feed", e.handleFeed)
	e.RegisterHandler("cut", e.handleCut)

	// Registrar handlers avanzados
	e.RegisterHandler("image", e.handleImage)
	e.RegisterHandler("separator", e.handleSeparator)

	// Handlers para QR y tablas
	e.RegisterHandler("qr", e.handleQR)
	e.RegisterHandler("table", e.handleTable)

	// Handler para barcode
	e.RegisterHandler("barcode", e.handleBarcode)

	// TODO: Implement other commands

	return e
}

// RegisterHandler registers a new command handler
func (e *Executor) RegisterHandler(cmdType string, handler CommandHandler) {
	e.handlers[cmdType] = handler
}

// Execute ejecuta un documento completo
func (e *Executor) Execute(doc *Document) error {
	// Inicializar impresora
	if err := e.printer.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize printer: %w", err)
	}

	// Aplicar configuración del profile desde JSON
	if err := e.applyProfileFromDocument(doc); err != nil {
		log.Printf("Warning: failed to apply profile settings: %v", err)
	}

	// Execute commands
	for i, cmd := range doc.Commands {
		handler, exists := e.handlers[cmd.Type]
		if !exists {
			return fmt.Errorf("unknown command type at position %d: %s", i, cmd.Type)
		}

		if err := handler(e.printer, cmd.Data); err != nil {
			return fmt.Errorf("command %d (%s) failed: %w", i, cmd.Type, err)
		}
	}

	return nil
}

// setCodeTable configura la tabla de caracteres con fallback
func (e *Executor) setCodeTable(tableName string) error {
	// Mapa de nombres a constantes
	tables := map[string]character.CodeTable{
		"PC437":   character.PC437,
		"PC850":   character.PC850,
		"PC852":   character.PC852,
		"WPC1252": character.WPC1252,
		// TODO: Add Go Compatible code tables (check package profile in escpos_encoding.go)
	}

	table, ok := tables[tableName]
	if !ok {
		log.Printf("warning: unsupported code table %s, falling back to Windows-1252", tableName)
		return e.printer.SetCodeTable(character.WPC1252)
	}

	return e.printer.SetCodeTable(table)
}

// ExecuteJSON ejecuta un documento desde JSON
func (e *Executor) ExecuteJSON(data []byte) error {
	doc, err := ParseDocument(data)
	if err != nil {
		return err
	}
	return e.Execute(doc)
}

// applyProfileFromDocument aplica la configuración del profile desde el documento JSON
func (e *Executor) applyProfileFromDocument(doc *Document) error {
	profile := e.printer.Profile

	if doc == nil {
		return fmt.Errorf("document is nil")
	}

	if doc.Profile.Model == "" {
		return fmt.Errorf("profile model is required")
	}
	profile.Model = doc.Profile.Model
	log.Printf("Profile: Model set to %s from JSON", doc.Profile.Model)

	if err := e.applyProfileOrDefaults(doc.Profile); err != nil {
		return fmt.Errorf("failed to apply profile settings: %w", err)
	}

	profile.DotsPerLine = calculate.DotsPerLine(profile.PaperWidth, profile.DPI)
	log.Printf("Profile: DotsPerLine calculated as %d", profile.DotsPerLine)

	return nil
}

// applyProfileOrDefaults aplica configuraciones por defecto al profile
func (e *Executor) applyProfileOrDefaults(config ProfileConfig) error {
	profile := e.printer.Profile
	if config.PaperWidth == 0 {
		// Default paper width 80mm
		profile.PaperWidth = 80
		log.Printf("Profile: PaperWidth set to default 80mm")
	} else {
		profile.PaperWidth = float64(config.PaperWidth)
		log.Printf("Profile: PaperWidth set to %dmm from JSON", config.PaperWidth)
	}

	if config.CodeTable == "" {
		// Default code table WPC1252
		err := e.setCodeTable("WPC1252")
		if err != nil {
			log.Printf("failed to set default code table WPC1252: %v", err)
		}
		log.Printf("Profile: CodeTable set to default WPC1252")
	} else {
		err := e.setCodeTable(config.CodeTable)
		if err != nil {
			return fmt.Errorf("failed to set code table %s: %w", config.CodeTable, err)
		}
		log.Printf("Profile: CodeTable set to %s from JSON", config.CodeTable)
	}

	if config.DPI == 0 {
		// Default DPI 203
		profile.DPI = 203
		log.Printf("Profile: DPI set to default 203")
	} else {
		profile.DPI = config.DPI
		log.Printf("Profile: DPI set to %d from JSON", config.DPI)
	}

	if config.HasQR {
		profile.HasQR = true
		log.Printf("Profile: HasQR set to true from JSON")
	} else {
		profile.HasQR = false
		log.Printf("Profile: HasQR set to false from JSON")
	}

	return nil
}
