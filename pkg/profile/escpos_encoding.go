package profile

import (
	"fmt"
	"log"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"

	"github.com/adcondev/pos-printer/pkg/commands/character"
)

// TODO: Add new encodings as needed, many of them defined in ESCPOS don't have a direct mapping in Go

// TODO: Separate this concern to dedicated package, then add the state management concern to profile

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

	// Asian (basic mappings, may need special handling)
	character.Katakana: japanese.ShiftJIS,
	character.Hiragana: japanese.ShiftJIS,
	// Note: OnePassKanji would need special handling
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
		// FIXME: Side effect: logging directly avoids caller control. Return error or use configured logger.
		log.Printf("warning: unsupported encoding for code table %v, falling back to Windows-1252", codeTable)
		// FIXME: Silent fallback to Windows-1252 might lead to incorrect output without caller knowing.
		return charmap.Windows1252.NewEncoder()
	}
	// TODO: Optimization: .NewEncoder() is called on every invocation. Consider caching or reusing encoders if strict state is not required.
	return enc.NewEncoder()
}

// IsSupported checks if the specified code table is supported
func (e *Escpos) IsSupported(codeTable character.CodeTable) bool {
	_, ok := codeTableMap[codeTable]
	return ok
}
