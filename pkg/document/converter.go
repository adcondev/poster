// Package document proporciona estructuras y funciones para construir documentos de impresión.
package document

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/adcondev/pos-printer/pkg/tables"
)

// TODO: Improve builder pattern, then add qr and table commands

// Builder ayuda a construir documentos programáticamente
type Builder struct {
	doc *Document
}

// NewBuilder crea un nuevo constructor de documentos
func NewBuilder() *Builder {
	return &Builder{
		doc: &Document{
			Version:  "1.0",
			Commands: []Command{},
		},
	}
}

// SetProfile configura el perfil de impresora
func (b *Builder) SetProfile(model string, width int, codeTable string) *Builder {
	b.doc.Profile = ProfileConfig{
		Model:      model,
		PaperWidth: width,
		CodeTable:  codeTable,
	}
	return b
}

// AddText generates a text command
func (b *Builder) AddText(content string, style *TextStyle, label *Label) *Builder {
	newLine := true
	cmd := TextCommand{
		Content: Content{
			Text: content,
		},
		NewLine: &newLine,
	}
	if style != nil {
		cmd.Content.Style = style
	}
	if label != nil {
		cmd.Label = label
	}

	data, err := json.Marshal(cmd)
	if err != nil {
		log.Printf("Error marshaling text command: %v", err)
	}
	b.doc.Commands = append(b.doc.Commands, Command{
		Type: "text",
		Data: data,
	})
	return b
}

// AddImage creates an image command
func (b *Builder) AddImage(base64Data string, width int, align string) *Builder {
	cmd := ImageCommand{
		Code:       base64Data,
		PixelWidth: width,
		Align:      align,
		Dithering:  "threshold",
	}

	data, err := json.Marshal(cmd)
	if err != nil {
		log.Printf("Error marshaling image command: %v", err)
	}
	b.doc.Commands = append(b.doc.Commands, Command{
		Type: "image",
		Data: data,
	})
	return b
}

// AddSeparator agrega un separador
func (b *Builder) AddSeparator(char string, length int) *Builder {
	cmd := SeparatorCommand{
		Char:   char,
		Length: length,
	}

	data, err := json.Marshal(cmd)
	if err != nil {
		log.Printf("Error marshaling separator command: %v", err)
	}
	b.doc.Commands = append(b.doc.Commands, Command{
		Type: "separator",
		Data: data,
	})
	return b
}

// AddFeed agrega avance de papel
func (b *Builder) AddFeed(lines int) *Builder {
	cmd := FeedCommand{Lines: lines}

	data, err := json.Marshal(cmd)
	if err != nil {
		log.Printf("Error marshaling feed command: %v", err)
	}
	b.doc.Commands = append(b.doc.Commands, Command{
		Type: "feed",
		Data: data,
	})
	return b
}

// AddCut agrega corte de papel
func (b *Builder) AddCut(mode string, feed int) *Builder {
	cmd := CutCommand{
		Mode: mode,
		Feed: feed,
	}

	data, err := json.Marshal(cmd)
	if err != nil {
		log.Printf("Error marshaling cut command: %v", err)
	}
	b.doc.Commands = append(b.doc.Commands, Command{
		Type: "cut",
		Data: data,
	})
	return b
}

// AddQR agrega un comando QR al documento
func (b *Builder) AddQR(data, text string, pixelWidth int, correction string, align, logo64 string, circle bool) *Builder {
	cmd := QRCommand{
		Data:        data,
		HumanText:   text,
		PixelWidth:  pixelWidth,
		Correction:  correction,
		Align:       align,
		Logo:        logo64,
		CircleShape: circle,
	}

	qrData, err := json.Marshal(cmd)
	if err != nil {
		log.Printf("Error marshaling QR command: %v", err)
		return b
	}

	b.doc.Commands = append(b.doc.Commands, Command{
		Type: "qr",
		Data: qrData,
	})
	return b
}

// AddTable adds a table command to the document
func (b *Builder) AddTable(definition tables.Definition, rows [][]string, showHeaders bool) *Builder {
	if len(definition.Columns) == 0 {
		log.Printf("Warning: table has no columns defined")
		return b
	}

	cmd := TableCommand{
		Definition:  definition,
		ShowHeaders: showHeaders,
		Rows:        rows,
		Options: &TableOptions{
			HeaderBold: true,
			WordWrap:   true,
		},
	}

	data, err := json.Marshal(cmd)
	if err != nil {
		log.Printf("Error marshaling table command: %v", err)
		return b
	}

	b.doc.Commands = append(b.doc.Commands, Command{
		Type: "table",
		Data: data,
	})
	return b
}

// AddBarcode agrega un comando de código de barras
func (b *Builder) AddBarcode(symbology, bcData string, width, height int, hriPos, hriFont, align string) *Builder {
	cmd := BarcodeCommand{
		Symbology: symbology,
		Data:      bcData,
	}

	// Solo establecer valores opcionales si son diferentes a los defaults
	if width > 0 {
		cmd.Width = &width
	}
	if height > 0 {
		cmd.Height = &height
	}
	if hriPos != "" {
		cmd.HRIPosition = &hriPos
	}
	if hriFont != "" {
		cmd.HRIFont = &hriFont
	}
	if align != "" {
		cmd.Align = &align
	}

	data, err := json.Marshal(cmd)
	if err != nil {
		log.Printf("Error marshaling barcode command: %v", err)
		return b
	}

	b.doc.Commands = append(b.doc.Commands, Command{
		Type: "barcode",
		Data: data,
	})
	return b
}

// AddRaw adds a raw command to the document
func (b *Builder) AddRaw(hexString string, comment string) *Builder {
	cmd := RawCommand{
		Hex:     hexString,
		Format:  "hex",
		Comment: comment,
	}

	data, err := json.Marshal(cmd)
	if err != nil {
		log.Printf("Error marshaling raw command: %v", err)
		return b
	}

	b.doc.Commands = append(b.doc.Commands, Command{
		Type: "raw",
		Data: data,
	})
	return b
}

// AddRawBytes adds raw bytes directly
func (b *Builder) AddRawBytes(bytes []byte, comment string) *Builder {
	hexStr := hex.EncodeToString(bytes)
	return b.AddRaw(hexStr, comment)
}

// ============================================================================
// Raw Command Convenience Methods (Escape Hatches)
// ============================================================================

// AddPulse opens cash drawer using standard pulse command
// This is a convenience wrapper around raw command for common operation
func (b *Builder) AddPulse() *Builder {
	// ESC p m t1 t2 (Pin 2/connector 1, 100ms on, 200ms off)
	return b.AddRawWithComment("1B 70 00 32 64", "Open cash drawer (100ms pulse)")
}

// AddBeep emits a beep sound (if buzzer is present)
func (b *Builder) AddBeep(times int) *Builder {
	if times <= 0 || times > 9 {
		times = 1
	}
	// ESC BEL (simple beep) - repeat for multiple beeps
	hexa := strings.Repeat("07 ", times)
	return b.AddRawWithComment(strings.TrimSpace(hexa), fmt.Sprintf("Beep %d times", times))
}

// AddRawWithComment adds a raw command with documentation
func (b *Builder) AddRawWithComment(hexString, comment string) *Builder {
	cmd := RawCommand{
		Hex:      hexString,
		Format:   "hex",
		Comment:  comment,
		SafeMode: false, // Convenience methods bypass safety by default
	}

	data, err := json.Marshal(cmd)
	if err != nil {
		log.Printf("Error marshaling raw command: %v", err)
		return b
	}

	b.doc.Commands = append(b.doc.Commands, Command{
		Type: "raw",
		Data: data,
	})
	return b
}

// AddRawSafe adds a raw command with safety checks enabled
func (b *Builder) AddRawSafe(hexString, comment string) *Builder {
	cmd := RawCommand{
		Hex:      hexString,
		Format:   "hex",
		Comment:  comment,
		SafeMode: true, // Explicitly enable safety
	}

	data, err := json.Marshal(cmd)
	if err != nil {
		log.Printf("Error marshaling raw command: %v", err)
		return b
	}

	b.doc.Commands = append(b.doc.Commands, Command{
		Type: "raw",
		Data: data,
	})
	return b
}

// Build construye el documento final
func (b *Builder) Build() *Document {
	return b.doc
}
