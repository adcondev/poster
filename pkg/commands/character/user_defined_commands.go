package character

import (
	"fmt"

	"github.com/adcondev/poster/pkg/commands/common"
)

// SelectUserDefinedCharacterSet selects or cancels the user-defined character set.
//
// Format:
//
//	ASCII:   ESC % n
//	Hex:     0x1B 0x25 n
//	Decimal: 27 37 n
//
// Range:
//
//	n = 0–255
//
// Default:
//
//	n = 0
//
// Parameters:
//
//	n: Controls user-defined character set (LSB is used):
//	   LSB = 0 -> User-defined character set canceled
//	   LSB = 1 -> User-defined character set selected
//
// Notes:
//   - When the user-defined character set is canceled the resident (built-in)
//     character set is automatically selected
//   - The setting remains in effect until ESC @ is executed, the printer is
//     reset, or power is turned off
//   - This command affects alphanumeric, Kana, multilingual, and user-defined
//     characters as applicable per model
//
// Errors:
//
//	This function is safe and does not return errors.
func (udc *UserDefinedCommands) SelectUserDefinedCharacterSet(n byte) []byte {
	return []byte{0x1B, '%', n}
}

// DefineUserDefinedCharacters defines user-defined glyph patterns for character codes.
//
// Format:
//
//	ASCII:   ESC & y c1 c2 [x1 d1...d(y*x1)]...[xk d1...d(y*xk)]
//	Hex:     0x1B 0x26 y c1 c2 [x1 data...]...
//	Decimal: 27 38 y c1 c2 [x1 data...]...
//
// Range:
//
//	y = 1–255 (model-dependent, typically 3 for 12x24 or 9x17 fonts)
//	c1, c2 = 32–126 (typically, model-dependent)
//
// Default:
//
//	None
//
// Parameters:
//
//	y: Number of bytes in the vertical direction for each column
//	c1: First character code to define (inclusive)
//	c2: Last character code to define (inclusive)
//	definitions: Slice of per-character definitions in order for codes c1...c2
//
// Notes:
//   - Typical y values are model/font dependent (e.g., 3 for 12x24 or 9x17 fonts)
//   - c1 and c2 typically in the printable range (32..126) depending on model
//   - For each character i between c1 and c2, definitions[i - int(c1)] must contain
//     the width byte and exactly y*width data bytes
//   - Each definition is encoded as: 1 byte width xi, followed by y*xi bytes
//     of column data (column-major)
//   - Each column is described by y bytes (little-endian vertical bit order)
//   - Existing user-defined characters for the specified codes are replaced
//   - Definitions persist until cleared (ESC ?, ESC @), reset, or power-off
//   - To use defined glyphs, send ESC % 1 (select user-defined character set)
//
// Errors:
//
//	Returns ErrYValue if y is 0.
//	Returns ErrCharacterCode if c1 is outside the valid range.
//	Returns ErrCodeRange if c2 is less than c1 or outside the valid range.
//	Returns ErrDefinition if the number of definitions doesn't match the code range.
//	Returns ErrDataLength if a definition's data length doesn't match y*width.
func (udc *UserDefinedCommands) DefineUserDefinedCharacters(y, c1, c2 byte, definitions []UserDefinedChar) ([]byte, error) {
	// Validation
	if y == 0 {
		return nil, ErrYValue
	}
	if err := ValidateCodeRange(c1, c2); err != nil {
		return nil, err
	}
	expected := int(c2 - c1 + 1)
	if len(definitions) != expected {
		return nil, fmt.Errorf("%w: got %d, expected %d", ErrDefinition, len(definitions), expected)
	}

	// Build command
	seq := []byte{common.ESC, '&', y, c1, c2}
	bytesPerCol := int(y)

	for idx, def := range definitions {
		// Width validation is printer & font dependent
		if def.Width == 0 {
			// Zero width is allowed (blank char)
			seq = append(seq, def.Width)
			continue
		}
		expectedDataLen := bytesPerCol * int(def.Width)
		if len(def.Data) != expectedDataLen {
			return nil, fmt.Errorf("%w: char %d has %d bytes, expected %d",
				ErrDataLength, int(c1)+idx, len(def.Data), expectedDataLen)
		}
		seq = append(seq, def.Width)
		seq = append(seq, def.Data...)
	}

	return seq, nil
}

// CancelUserDefinedCharacter deletes (cancels) a user-defined character.
//
// Format:
//
//	ASCII:   ESC ? n
//	Hex:     0x1B 0x3F n
//	Decimal: 27 63 n
//
// Range:
//
//	n = 32–126
//
// Default:
//
//	None
//
// Parameters:
//
//	n: Character code to delete
//
// Notes:
//   - Deletes the user-defined character pattern specified by character code n
//   - After cancellation, the resident (built-in) character for that code is printed
//   - This command can cancel user-defined characters per font. Select the font
//     with ESC ! or ESC M before issuing this command if needed
//   - Settings take effect immediately; the deleted definition remains cleared
//     until redefined, ESC @ (initialize), power-off, or reset
//
// Errors:
//
//	Returns ErrCharacterCode if n is outside the valid range (32-126).
func (udc *UserDefinedCommands) CancelUserDefinedCharacter(n byte) ([]byte, error) {
	if err := ValidateCharacterCode(n); err != nil {
		return nil, err
	}
	return []byte{common.ESC, '?', n}, nil
}
