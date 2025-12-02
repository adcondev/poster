package builder

import (
	"encoding/json"
	"log"

	"github.com/adcondev/pos-printer/pkg/document/schema"
)

// DocumentBuilder constructs print documents using a fluent API
type DocumentBuilder struct {
	version  string
	profile  schema.ProfileConfig
	debugLog bool
	commands []schema.Command
}

// NewDocument creates a new document builder
func NewDocument() *DocumentBuilder {
	return &DocumentBuilder{
		version:  "1.0",
		commands: []schema.Command{},
	}
}

// SetProfile configures the printer profile
func (b *DocumentBuilder) SetProfile(model string, paperWidth int, codeTable string) *DocumentBuilder {
	b.profile = schema.ProfileConfig{
		Model:      model,
		PaperWidth: paperWidth,
		CodeTable:  codeTable,
	}
	return b
}

// SetDPI sets the printer DPI
func (b *DocumentBuilder) SetDPI(dpi int) *DocumentBuilder {
	b.profile.DPI = dpi
	return b
}

// SetHasQR indicates native QR support
func (b *DocumentBuilder) SetHasQR(hasQR bool) *DocumentBuilder {
	b.profile.HasQR = hasQR
	return b
}

// EnableDebug enables debug logging
func (b *DocumentBuilder) EnableDebug() *DocumentBuilder {
	b.debugLog = true
	return b
}

// Build creates the final Document
func (b *DocumentBuilder) Build() *schema.Document {
	return &schema.Document{
		Version:  b.version,
		Profile:  b.profile,
		DebugLog: b.debugLog,
		Commands: b.commands,
	}
}

// ToJSON converts the document to JSON bytes
func (b *DocumentBuilder) ToJSON() ([]byte, error) {
	return json.MarshalIndent(b.Build(), "", "  ")
}

// addCommand is an internal helper to add commands
func (b *DocumentBuilder) addCommand(cmdType string, data interface{}) *DocumentBuilder {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshaling %s command: %v", cmdType, err)
		return b
	}
	b.commands = append(b.commands, schema.Command{
		Type: cmdType,
		Data: jsonData,
	})
	return b
}

// ============================================================================
// Builder-returning methods
// ============================================================================

// Text starts building a text command
func (b *DocumentBuilder) Text(content string) *TextBuilder {
	return newTextBuilder(b, content)
}

// Table starts building a table command
func (b *DocumentBuilder) Table() *TableBuilder {
	return newTableBuilder(b)
}

// QR starts building a QR code command
func (b *DocumentBuilder) QR(data string) *QRBuilder {
	return newQRBuilder(b, data)
}

// Barcode starts building a barcode command
func (b *DocumentBuilder) Barcode(symbology, data string) *BarcodeBuilder {
	return newBarcodeBuilder(b, symbology, data)
}

// Image starts building an image command
func (b *DocumentBuilder) Image(base64Data string) *ImageBuilder {
	return newImageBuilder(b, base64Data)
}

// Raw starts building a raw command
func (b *DocumentBuilder) Raw(hexData string) *RawBuilder {
	return newRawBuilder(b, hexData)
}

// ============================================================================
// Internal command structures
// ============================================================================

type feedCommand struct {
	Lines int `json:"lines"`
}

type cutCommand struct {
	Mode string `json:"mode,omitempty"`
	Feed int    `json:"feed,omitempty"`
}

type separatorCommand struct {
	Char   string `json:"char,omitempty"`
	Length int    `json:"length,omitempty"`
}
