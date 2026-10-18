package profile

import (
	"fmt"

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
	enc := e.getEncoding(e.CodeTable)
	result, err := enc.String(text)
	if err != nil {
		return "", fmt.Errorf("failed to encode string: %w", err)
	}
	return result, nil
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
