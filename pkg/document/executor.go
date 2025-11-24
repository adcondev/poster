package document

import (
	"encoding/json"
	"fmt"
	"log"

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
	if doc == nil {
		return fmt.Errorf("document is nil")
	}
	if doc.Profile.Model == "" && doc.Profile.PaperWidth == 0 && doc.Profile.CodeTable == "" && !doc.Profile.HasQR {
		// No hay configuración de profile para aplicar
		return nil
	}

	// Aplicar Model si se especifica
	if doc.Profile.Model != "" {
		e.printer.Profile.Model = doc.Profile.Model
		log.Printf("Profile: Model set to %s from JSON", doc.Profile.Model)
	}

	// Aplicar HasQR
	e.printer.Profile.HasQR = doc.Profile.HasQR
	log.Printf("Profile: HasQR set to %v from JSON", doc.Profile.HasQR)

	// Aplicar CodeTable si se especifica
	if doc.Profile.CodeTable != "" {
		if err := e.setCodeTable(doc.Profile.CodeTable); err != nil {
			log.Printf("Warning: failed to set code table %s: %v", doc.Profile.CodeTable, err)
		}
	}

	// Aplicar otros campos si vienen en el JSON
	if doc.Profile.PaperWidth > 0 {
		// Calcular DotsPerLine basado en PaperWidth y DPI
		if doc.Profile.DPI > 0 {
			e.printer.Profile.DPI = doc.Profile.DPI
		}
		e.printer.Profile.PaperWidth = float64(doc.Profile.PaperWidth)
		// Recalcular DotsPerLine: (PaperWidth_mm * DPI) / 25.4
		e.printer.Profile.DotsPerLine = int(float64(doc.Profile.PaperWidth) * float64(e.printer.Profile.DPI) / 25.4)
		log.Printf("Profile: Updated PaperWidth=%dmm, DotsPerLine=%d",
			doc.Profile.PaperWidth, e.printer.Profile.DotsPerLine)
	}

	// Model name (útil para debugging)
	if doc.Profile.Model != "" {
		e.printer.Profile.Model = doc.Profile.Model
		log.Printf("Profile: Model set to '%s'", doc.Profile.Model)
	}

	return nil
}
