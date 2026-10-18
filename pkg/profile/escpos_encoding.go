package profile

import (
	"fmt"
	"strings"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"

	"github.com/adcondev/poster/pkg/commands/character"
)

// codeTableMap maps ESC/POS code tables to Go encodings
var codeTableMap = map[character.CodeTable]encoding.Encoding{
	// Western European and Americas
	character.PC437:     charmap.CodePage437,
	character.PC850:     charmap.CodePage850,
	character.PC852:     charmap.CodePage852,
	character.PC858:     charmap.CodePage858,
	character.PC860:     charmap.CodePage860,
	character.PC863:     charmap.CodePage863,
	character.PC865:     charmap.CodePage865,
	character.WPC1252:   charmap.Windows1252,
	character.ISO885915: charmap.ISO8859_15,

	// Cyrillic
	character.PC855:   charmap.CodePage855,
	character.PC866:   charmap.CodePage866,
	character.WPC1251: charmap.Windows1251,

	// Greek
	character.ISO88597: charmap.ISO8859_7,
	character.WPC1253:  charmap.Windows1253,

	// Turkish
	character.WPC1254: charmap.Windows1254,

	// Baltic
	character.WPC1257: charmap.Windows1257,

	// Hebrew/Arabic
	character.PC862:   charmap.CodePage862,
	character.WPC1255: charmap.Windows1255,
	character.WPC1256: charmap.Windows1256,

	// Central European
	character.WPC1250:  charmap.Windows1250,
	character.ISO88592: charmap.ISO8859_2,

	// Vietnamese
	character.WPC1258: charmap.Windows1258,

	// Asian (basic mappings)
	character.Katakana: japanese.ShiftJIS,
	character.Hiragana: japanese.ShiftJIS,
}

// EncodeString encodes a string using the specified code table
func (e *Escpos) EncodeString(text string) (string, error) {
	text = NormalizeTypography(text)
	enc := e.getEncoding(e.CodeTable)
	result, err := enc.String(text)
	if err != nil {
		return "", fmt.Errorf("failed to encode string: %w", err)
	}
	return result, nil
}

// ParseCodeTable normalizes a code table name or alias and returns the corresponding character.CodeTable.
// It removes hyphens, underscores, and spaces and converts to uppercase before matching.
func ParseCodeTable(name string) (character.CodeTable, bool) {
	clean := strings.ToUpper(strings.TrimSpace(name))
	clean = strings.ReplaceAll(clean, "-", "")
	clean = strings.ReplaceAll(clean, "_", "")
	clean = strings.ReplaceAll(clean, " ", "")

	switch clean {
	case "PC437", "CP437", "437":
		return character.PC437, true
	case "PC850", "CP850", "850":
		return character.PC850, true
	case "PC852", "CP852", "852":
		return character.PC852, true
	case "PC855", "CP855", "855":
		return character.PC855, true
	case "PC858", "CP858", "858":
		return character.PC858, true
	case "PC860", "CP860", "860":
		return character.PC860, true
	case "PC862", "CP862", "862":
		return character.PC862, true
	case "PC863", "CP863", "863":
		return character.PC863, true
	case "PC865", "CP865", "865":
		return character.PC865, true
	case "PC866", "CP866", "866":
		return character.PC866, true
	case "WPC1250", "CP1250", "WINDOWS1250", "1250":
		return character.WPC1250, true
	case "WPC1251", "CP1251", "WINDOWS1251", "1251":
		return character.WPC1251, true
	case "WPC1252", "CP1252", "WINDOWS1252", "1252":
		return character.WPC1252, true
	case "WPC1253", "CP1253", "WINDOWS1253", "1253":
		return character.WPC1253, true
	case "WPC1254", "CP1254", "WINDOWS1254", "1254":
		return character.WPC1254, true
	case "WPC1255", "CP1255", "WINDOWS1255", "1255":
		return character.WPC1255, true
	case "WPC1256", "CP1256", "WINDOWS1256", "1256":
		return character.WPC1256, true
	case "WPC1257", "CP1257", "WINDOWS1257", "1257":
		return character.WPC1257, true
	case "WPC1258", "CP1258", "WINDOWS1258", "1258":
		return character.WPC1258, true
	case "ISO88592", "88592":
		return character.ISO88592, true
	case "ISO88597", "88597":
		return character.ISO88597, true
	case "ISO885915", "885915":
		return character.ISO885915, true
	case "KATAKANA", "SHIFTJIS", "SJIS":
		return character.Katakana, true
	case "HIRAGANA":
		return character.Hiragana, true
	default:
		return 0, false
	}
}

// getEncoding returns the encoding.Encoder for the specified code table
func (e *Escpos) getEncoding(codeTable character.CodeTable) *encoding.Encoder {
	enc, ok := codeTableMap[codeTable]
	if !ok {
		return charmap.Windows1252.NewEncoder()
	}
	return enc.NewEncoder()
}

// IsSupported checks if the specified code table is supported
func (e *Escpos) IsSupported(codeTable character.CodeTable) bool {
	_, ok := codeTableMap[codeTable]
	return ok
}
