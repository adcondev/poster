package linespacing

import (
	"github.com/adcondev/poster/pkg/commands/shared"
)

// SetLineSpacing sets the line spacing to spacing × (vertical or horizontal motion unit).
//
// Format:
//
//	ASCII:   ESC 3 n
//	Hex:     0x1B 0x33 n
//	Decimal: 27 51 n
//
// Range:
//
//	n = 0–255
//
// Default:
//
//	The amount of line spacing corresponding to the "default line spacing"
//	(equivalent to a spacing between 30 and 80 dots)
//
// Parameters:
//
//	n: Sets the line spacing to n × (vertical or horizontal motion unit)
//
// Notes:
//   - Maximum line spacing is 1016 mm (40 inches); if exceeded, printer uses maximum
//   - In Standard mode the vertical motion unit is used
//   - In Page mode the motion unit depends on ESC T setting
//   - Line spacing can be set independently in Standard and Page modes
//   - Motion unit changes after setting don't affect the numeric spacing
//   - Remains in effect until ESC 2, ESC @, reset, or power off
//
// Errors:
//
//	This function is safe and does not return errors.
func (lsc *Commands) SetLineSpacing(n Spacing) []byte {
	return []byte{shared.ESC, '3', byte(n)}
}

// SelectDefaultLineSpacing sets the line spacing to the printer's default.
//
// Format:
//
//	ASCII:   ESC 2
//	Hex:     0x1B 0x32
//	Decimal: 27 50
//
// Range:
//
//	Not applicable
//
// Default:
//
//	Not applicable
//
// Parameters:
//
//	None
//
// Notes:
//   - Sets the line spacing to the default line spacing
//   - Line spacing can be set independently in Standard and Page modes
//   - Remains in effect until ESC 3, ESC @, reset, or power off
//
// Errors:
//
//	This function is safe and does not return errors.
func (lsc *Commands) SelectDefaultLineSpacing() []byte {
	return []byte{shared.ESC, '2'}
}
