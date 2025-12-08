package character

import (
	"github.com/adcondev/poster/pkg/commands/common"
)

// SelectCharacterEncodeSystem selects the character encoding system.
//
// Format:
//
//	ASCII:   FS ( C pL pH fn m
//	Hex:     0x1C 0x28 0x43 0x02 0x00 0x30 m
//	Decimal: 28 40 67 2 0 48 m
//
// Range:
//
//	pL = 0x02, pH = 0x00
//	m = 1, 2, 49, 50
//
// Default:
//
//	m = 1 (1-byte encoding)
//
// Parameters:
//
//	m: Selects the character encoding system:
//	   1 or 49 -> 1-byte (legacy) encoding (model-dependent legacy code pages)
//	   2 or 50 -> UTF-8 (Unicode)
//
// Notes:
//   - When UTF-8 is selected, ESC t (code table selection) is ignored
//   - Settings persist until ESC @ (initialize), printer reset, or power off
//   - Availability of specific legacy encodings is model-dependent
//
// Errors:
//
//	Returns ErrEncoding if m is not a valid encoding system value (1, 2, 49, 50).
func (c *CodeConversionCommands) SelectCharacterEncodeSystem(m EncodeSystem) ([]byte, error) {
	// Validate allowed values
	if err := ValidateEncodeSystem(m); err != nil {
		return nil, err
	}
	return []byte{common.FS, '(', 'C', 0x02, 0x00, 0x30, byte(m)}, nil
}

// SetFontPriority sets the font priority.
//
// Format:
//
//	ASCII:   FS ( C pL pH fn m a
//	Hex:     0x1C 0x28 0x43 0x03 0x00 0x3C m a
//	Decimal: 28 40 67 3 0 60 m a
//
// Range:
//
//	pL = 0x03, pH = 0x00
//	m = 0–1
//	a = 0, 11, 20, 30, 41
//
// Default:
//
//	m = 0, a = 0
//
// Parameters:
//
//	m: Priority rank (0 = 1st priority, 1 = 2nd priority)
//	a: Font type:
//	   0  -> AnkSansSerif font (Sans serif)
//	   11 -> Japanese font (Gothic)
//	   20 -> Simplified Chinese font (Mincho)
//	   30 -> Traditional Chinese font (Mincho)
//	   41 -> Korean font (Gothic)
//
// Notes:
//   - Assigns a font style to a priority slot (1st or 2nd)
//   - If the style already exists in the priority list, promotion/demotion
//     is handled so that the newly specified font becomes the selected priority
//   - Settings persist until ESC @ (initialize), printer reset, or power-off
//
// Errors:
//
//	Returns ErrFontPriority if m is greater than 1.
//	Returns ErrFontType if a is not a valid font type value.
func (c *CodeConversionCommands) SetFontPriority(m FontPriority, a FontFunction) ([]byte, error) {
	// Validate allowed values
	if err := ValidateFontPriority(m); err != nil {
		return nil, err
	}
	if err := ValidateFontFunction(a); err != nil {
		return nil, err
	}

	return []byte{common.FS, '(', 'C', 0x03, 0x00, 0x3C, byte(m), byte(a)}, nil
}
