package character

import (
	"github.com/adcondev/poster/pkg/commands/shared"
)

// SetRightSideCharacterSpacing sets the right-side character spacing.
//
// Format:
//
//	ASCII:   ESC SP n
//	Hex:     0x1B 0x20 n
//	Decimal: 27 32 n
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
//	n: Sets the right-side character spacing to n × (horizontal or vertical motion unit)
//
// Notes:
//   - The character spacing set by this command is effective for alphanumeric,
//     Kana, and user-defined characters
//   - When characters are enlarged, the character spacing becomes n times the normal value
//   - In Standard mode the horizontal motion unit is used
//   - In Page mode the vertical or horizontal motion unit is used according
//     to the print direction set by ESC T
//   - When the starting position is set to the upper-left or lower-right of
//     the print area using ESC T, the horizontal motion unit is used
//   - When the starting position is set to the upper-right or lower-left of
//     the print area using ESC T, the vertical motion unit is used
//   - Character spacing can be set independently in Standard mode and in
//     Page mode; this command affects the spacing for the currently selected mode
//   - If the horizontal or vertical motion unit is changed after this
//     command is executed, the numeric character spacing value does not change
//   - The setting remains in effect until ESC @ is executed, the printer is
//     reset, or power is turned off
//   - This command is used to change the spacing between characters
//
// Errors:
//
//	This function is safe and does not return errors.
func (c *Commands) SetRightSideCharacterSpacing(n Spacing) []byte {
	return []byte{shared.ESC, shared.SP, byte(n)}
}

// SelectPrintModes selects character font and style bits (emphasized, double-height, double-width, underline) together.
//
// Format:
//
//	ASCII:   ESC ! n
//	Hex:     0x1B 0x21 n
//	Decimal: 27 33 n
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
//	n: Bit-mapped byte controlling font and styles:
//	   Bit 0 (0x01) - Font selection (0: Font 1, 1: Font 2)
//	   Bit 3 (0x08) - Emphasized mode (0: OFF, 1: ON)
//	   Bit 4 (0x10) - Double-height mode (0: OFF, 1: ON)
//	   Bit 5 (0x20) - Double-width mode (0: OFF, 1: ON)
//	   Bit 7 (0x80) - Underline mode (0: OFF, 1: ON)
//
// Notes:
//   - Configurations for Font 1 and Font 2 differ by model. If the desired
//     font type cannot be selected with this command, use ESC M
//   - Bits 0, 4, 5 and 7 affect 1-byte code characters. On some models,
//     bits 4, 5 and 7 also affect Korean characters
//   - Emphasized mode (bit 3) is effective for both 1-byte and multi-byte characters
//   - Settings remain in effect until ESC @ is executed, the printer is
//     reset, power is turned off, or one of these commands is executed:
//     Bit 0 (font): ESC M, Bit 3 (emphasized): ESC E,
//     Bit 4,5 (size): GS !, Bit 7 (underline): ESC -
//   - When some characters in a line are double-height, all characters on
//     the line are aligned at the baseline
//   - Double-width enlarges characters to the right from the left side of
//     the character. When both double-height and double-width are on,
//     characters become quadruple size
//   - In Standard mode double-height enlarges in the paper-feed direction
//     and double-width enlarges perpendicular to paper feed. Rotating
//     characters 90° clockwise swaps the relationship
//   - In Page mode double-height and double-width follow the character orientation
//   - Underline thickness is determined by ESC -, regardless of character
//     size. Underline color matches the printed character color (GS ( N <Function 48>)
//   - The following are not underlined: 90° clockwise-rotated characters,
//     white/black reverse characters, spaces set by HT, ESC $, and ESC \
//   - On printers with Automatic font replacement (GS ( E <Function 5> with
//     a = 111,112,113), the replacement font is selected by this command
//
// Errors:
//
//	This function is safe and does not return errors.
func (c *Commands) SelectPrintModes(n PrintMode) []byte {
	return []byte{shared.ESC, '!', byte(n)}
}

// SetUnderlineMode sets underline mode on or off and selects underline thickness.
//
// Format:
//
//	ASCII:   ESC - n
//	Hex:     0x1B 0x2D n
//	Decimal: 27 45 n
//
// Range:
//
//	n = 0, 1, 2, 48, 49, 50
//
// Default:
//
//	n = 0
//
// Parameters:
//
//	n: Controls underline mode:
//	   0 or 48 -> Turns off underline mode
//	   1 or 49 -> Turns on underline mode (1-dot thick)
//	   2 or 50 -> Turns on underline mode (2-dots thick)
//
// Notes:
//   - The underline mode is effective for alphanumeric, Kana, and user
//     defined characters. On some models it is also effective for Korean characters
//   - The underline color matches the printed character color (see GS ( N <Function 48>)
//   - Changing the character size does not affect the current underline thickness
//   - When underline mode is turned off the underline thickness value is
//     retained but no underline is produced
//   - The printer does not underline 90° clockwise-rotated characters,
//     white/black reverse characters, or spaces produced by HT, ESC $, and ESC \
//   - The setting remains in effect until ESC ! is executed, ESC @ is
//     executed, the printer is reset, or power is turned off
//   - Some printer models support the 2-dot thick underline (n = 2 or 50)
//
// Errors:
//
//	Returns ErrUnderlineMode if n is not a valid underline mode value (0, 1, 2, 48, 49, 50).
func (c *Commands) SetUnderlineMode(n UnderlineMode) ([]byte, error) {
	// Validate allowed values
	if err := ValidateUnderlineMode(n); err != nil {
		return nil, err
	}
	return []byte{shared.ESC, '-', byte(n)}, nil
}

// SetEmphasizedMode turns emphasized (bold) mode on or off.
//
// Format:
//
//	ASCII:   ESC E n
//	Hex:     0x1B 0x45 n
//	Decimal: 27 69 n
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
//	n: Controls emphasized mode (LSB is used):
//	   LSB = 0 -> Emphasized mode OFF
//	   LSB = 1 -> Emphasized mode ON
//
// Notes:
//   - This mode is effective for alphanumeric, Kana, multilingual, and user-defined characters
//   - Settings of this command remain in effect until ESC ! is executed,
//     ESC @ is executed, the printer is reset, or power is turned off
//
// Errors:
//
//	This function is safe and does not return errors.
func (c *Commands) SetEmphasizedMode(n EmphasizedMode) []byte {
	return []byte{shared.ESC, 'E', byte(n)}
}

// SetDoubleStrikeMode turns double-strike mode on or off.
//
// Format:
//
//	ASCII:   ESC G n
//	Hex:     0x1B 0x47 n
//	Decimal: 27 71 n
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
//	n: Controls double-strike mode (LSB is used):
//	   LSB = 0 -> Double-strike mode OFF
//	   LSB = 1 -> Double-strike mode ON
//
// Notes:
//   - This mode is effective for alphanumeric, Kana, multilingual, and user-defined characters
//   - Settings of this command remain in effect until ESC ! is executed,
//     ESC @ is executed, the printer is reset, or power is turned off
//
// Errors:
//
//	This function is safe and does not return errors.
func (c *Commands) SetDoubleStrikeMode(n DoubleStrikeMode) []byte {
	return []byte{shared.ESC, 'G', byte(n)}
}

// SelectCharacterFont selects a character font.
//
// Format:
//
//	ASCII:   ESC M n
//	Hex:     0x1B 0x4D n
//	Decimal: 27 77 n
//
// Range:
//
//	n = 0–4, 48–52, 97, 98 (model-dependent)
//
// Default:
//
//	n = 0 or 1 (model-dependent)
//
// Parameters:
//
//	n: Selects character font:
//	   0 or 48 -> Font A
//	   1 or 49 -> Font B
//	   2 or 50 -> Font C
//	   3 or 51 -> Font D
//	   4 or 52 -> Font E
//	   97      -> Special font A
//	   98      -> Special font B
//
// Notes:
//   - The selected character font is effective for alphanumeric, Kana, and user-defined characters
//   - Configurations of Font A and Font B depend on the printer model
//   - Settings remain in effect until ESC ! is executed, ESC @ is executed,
//     the printer is reset, or the power is turned off
//   - On printers with the Automatic font replacement function, the
//     replaced font selected by GS ( E <Function 5> (a = 111, 112, 113)
//     is selected by this command
//
// Errors:
//
//	Returns ErrCharacterFont if n is not a valid font selector value.
func (c *Commands) SelectCharacterFont(n FontType) ([]byte, error) {
	// Validate allowed values
	if err := ValidateFontType(n); err != nil {
		return nil, err
	}
	return []byte{shared.ESC, 'M', byte(n)}, nil
}

// SelectInternationalCharacterSet selects an international character set.
//
// Format:
//
//	ASCII:   ESC R n
//	Hex:     0x1B 0x52 n
//	Decimal: 27 82 n
//
// Range:
//
//	n = 0–17 (shared)
//	n = 66–75, 82 (India-specific on some models)
//
// Default:
//
//	n = 0 (most models)
//	n = 8 (Japanese models)
//	n = 13 (Korean models)
//	n = 15 (Simplified Chinese models)
//
// Parameters:
//
//	n: Selects international character set:
//	   0   U.S.A.
//	   1   France
//	   2   Germany
//	   3   U.K.
//	   4   Denmark I
//	   5   Sweden
//	   6   Italy
//	   7   Spain I
//	   8   Japan
//	   9   Norway
//	   10  Denmark II
//	   11  Spain II
//	   12  Latin America
//	   13  Korea
//	   14  Slovenia / Croatia
//	   15  China
//	   16  Vietnam
//	   17  Arabia
//	   66  India (Devanagari)
//	   67  India (Bengali)
//	   68  India (Tamil)
//	   69  India (Telugu)
//	   70  India (Assamese)
//	   71  India (Oriya)
//	   72  India (Kannada)
//	   73  India (Malayalam)
//	   74  India (Gujarati)
//	   75  India (Punjabi)
//	   82  India (Marathi)
//
// Notes:
//   - The selected international character set remains in effect until ESC @
//     is executed, the printer is reset, or power is turned off
//   - Refer to the printer's Character Code Tables for model-specific
//     mappings and supported characters
//
// Errors:
//
//	Returns ErrCharacterSet if n is not a valid character set value.
func (c *Commands) SelectInternationalCharacterSet(n InternationalSet) ([]byte, error) {
	if err := ValidateInternationalSet(n); err != nil {
		return nil, err
	}
	return []byte{shared.ESC, 'R', byte(n)}, nil
}

// Set90DegreeClockwiseRotationMode turns 90° clockwise rotation mode on or off.
//
// Format:
//
//	ASCII:   ESC V n
//	Hex:     0x1B 0x56 n
//	Decimal: 27 86 n
//
// Range:
//
//	n = 0–2, 48–50
//
// Default:
//
//	n = 0
//
// Parameters:
//
//	n: Controls 90° clockwise rotation mode in Standard mode:
//	   0 or 48 -> Turns off 90° clockwise rotation mode
//	   1 or 49 -> Turns on 90° clockwise rotation mode (1-dot character spacing)
//	   2 or 50 -> Turns on 90° clockwise rotation mode (1.5-dot character spacing)
//
// Notes:
//   - This mode is effective for alphanumeric, Kana, multilingual, and user-defined characters
//   - When underline mode is turned on, the printer does not underline
//     90° clockwise-rotated characters
//   - When character orientation changes in 90° clockwise rotation mode,
//     the relationship between vertical and horizontal directions is reversed
//   - The 90° clockwise rotation mode has no effect in Page mode
//   - Some printer models support n = 2 (1.5-dot spacing); some models have
//     fonts for which 90° rotation is not effective
//   - Settings remain in effect until ESC @ is executed, the printer is
//     reset, or power is turned off
//
// Errors:
//
//	Returns ErrRotationMode if n is not a valid rotation mode value (0-2, 48-50).
func (c *Commands) Set90DegreeClockwiseRotationMode(n RotationMode) ([]byte, error) {
	// Validate allowed values
	if err := ValidateRotationMode(n); err != nil {
		return nil, err
	}
	return []byte{shared.ESC, 'V', byte(n)}, nil
}

// SelectPrintColor selects the print color.
//
// Format:
//
//	ASCII:   ESC r n
//	Hex:     0x1B 0x72 n
//	Decimal: 27 114 n
//
// Range:
//
//	n = 0, 1, 48, 49
//
// Default:
//
//	n = 0
//
// Parameters:
//
//	n: Selects print color:
//	   0 or 48 -> Black
//	   1 or 49 -> Red
//
// Notes:
//   - In Standard mode this command is enabled only when processed at the
//     Beginning of the line
//   - In Page mode the color setting is applied to all data collectively
//     printed by FF (in Page mode) or ESC FF
//   - The setting remains in effect until ESC @ is executed, the printer is
//     reset, or power is turned off
//   - For printers that support two-color printing, GS ( N and GS ( L / GS 8 L
//     are available to define and control character/background/graphics
//     color layers. Use model-specific GS ( N / GS ( L / GS 8 L commands when
//     available for more advanced two-color workflows
//
// Errors:
//
//	Returns ErrPrintColor if n is not a valid color value (0, 1, 48, 49).
func (c *Commands) SelectPrintColor(n PrintColor) ([]byte, error) {
	// Validate allowed values
	if err := ValidatePrintColor(n); err != nil {
		return nil, err
	}
	return []byte{shared.ESC, 'r', byte(n)}, nil
}

// SelectCharacterCodeTable selects a character code table page.
//
// Format:
//
//	ASCII:   ESC t n
//	Hex:     0x1B 0x74 n
//	Decimal: 27 116 n
//
// Range:
//
//	n = 0–8, 11–26, 30–53, 66–75, 82, 254, 255 (model-dependent)
//
// Default:
//
//	n = 0
//
// Parameters:
//
//	n: Selects a page from the character code table:
//	   0   Page 0  [PC437: USA, Standard Europe]
//	   1   Page 1  [Katakana]
//	   2   Page 2  [PC850: Multilingual]
//	   3   Page 3  [PC860: Portuguese]
//	   4   Page 4  [PC863: Canadian-French]
//	   5   Page 5  [PC865: Nordic]
//	   6   Page 6  [Hiragana]
//	   7   Page 7  [One-pass printing Kanji characters]
//	   8   Page 8  [One-pass printing Kanji characters]
//	   11  Page 11 [PC851: Greek]
//	   12  Page 12 [PC853: Turkish]
//	   13  Page 13 [PC857: Turkish]
//	   14  Page 14 [PC737: Greek]
//	   15  Page 15 [ISO8859-7: Greek]
//	   16  Page 16 [WPC1252]
//	   17  Page 17 [PC866: Cyrillic #2]
//	   18  Page 18 [PC852: Latin 2]
//	   19  Page 19 [PC858: Euro]
//	   20-26 Thai Character Codes
//	   30-31 TCVN-3: Vietnamese
//	   32-53 Various international code pages
//	   66-75 India language pages
//	   82  Page 82 [Marathi]
//	   254, 255 Special pages
//
// Notes:
//   - The alphanumeric range (ASCII 0x20–0x7F / decimal 32–127) is the same
//     across pages; differences appear in the extended range (0x80–0xFF)
//   - The selected code table remains in effect until ESC @ is executed, the
//     printer is reset, or power is turned off
//   - Consult your printer's Character Code Tables for exact glyph mappings
//     per page and model-specific supported pages
//
// Errors:
//
//	Returns ErrCodeTablePage if n is not a valid code table page number.
func (c *Commands) SelectCharacterCodeTable(n CodeTable) ([]byte, error) {
	if err := ValidateCodeTable(n); err != nil {
		return nil, err
	}
	return []byte{shared.ESC, 't', byte(n)}, nil
}

// SetUpsideDownMode turns upside-down (180° rotated) print mode on or off.
//
// Format:
//
//	ASCII:   ESC { n
//	Hex:     0x1B 0x7B n
//	Decimal: 27 123 n
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
//	n: Controls upside-down print mode in Standard mode (LSB is used):
//	   LSB = 0 -> Upside-down print mode OFF
//	   LSB = 1 -> Upside-down print mode ON
//
// Notes:
//   - In Standard mode this command is only valid when processed at the
//     beginning of a line
//   - Upside-down mode is effective for all Standard-mode data except certain
//     graphics and obsolete raster/variable-size image commands (see model documentation)
//   - The mode has no effect in Page mode
//   - When turned on, characters are printed rotated 180° from right to
//     left. The line printing order is not reversed, so take care with the
//     order of transmitted data
//   - Settings remain in effect until ESC @ is executed, the printer is
//     reset, or power is turned off
//
// Errors:
//
//	This function is safe and does not return errors.
func (c *Commands) SetUpsideDownMode(n UpsideDownMode) []byte {
	return []byte{shared.ESC, '{', byte(n)}
}

// [Range]
// n = 0xxx0xxxb (n = 0 – 7, 16 – 23, 32 – 39, 48 – 55, 64 – 71, 80 – 87, 96 – 103, 112 – 119)
// (Enlargement in vertical direction: 1–8, Enlargement in horizontal direction: 1–8)

// SelectCharacterSize selects character size (width and height magnification).
//
// Format:
//
//	ASCII:   GS ! n
//	Hex:     0x1D 0x21 n
//	Decimal: 29 33 n
//
// Range:
//
//	n = 0xxx0xxxb (bit pattern)
//	Width magnification: 1–8
//	Height magnification: 1–8
//	Valid n examples: 0–7, 16–23, 32–39, 48–55, 64–71, 80–87, 96–103, 112–119
//
// Default:
//
//	n = 0 (normal size)
//
// Parameters:
//
//	n: Character size encoded as:
//	   Bits 0–2: Height magnification (value 0..7 -> x1..x8 where stored value = height-1)
//	   Bits 4–6: Width magnification (value 0..7 -> x1..x8 where stored value = width-1)
//	   Formula: n = ((width-1) << 4) | (height-1)
//
// Notes:
//   - The character size set by this command is effective for alphanumeric,
//     Kana, multilingual, and user-defined characters
//   - When characters on a line have different heights, they are aligned at the baseline
//   - Width enlargement extends characters to the right from the left side
//   - ESC ! can also toggle double-width and double-height modes
//   - In Standard mode double-height enlarges in the paper-feed direction
//     and double-width enlarges perpendicular to the paper feed. In 90°
//     rotated mode the relationship is reversed. In Page mode the size
//     follows the character orientation
//   - The setting for alphanumeric and Katakana remains until ESC !,
//     ESC @, reset, or power-off. For Kanji/multilingual chars the setting
//     remains until FS !, FS W, ESC @, reset, or power-off
//
// Errors:
//
//	This function is safe and does not return errors.
func (c *Commands) SelectCharacterSize(n Size) []byte {
	return []byte{shared.GS, '!', byte(n)}
}

// SetWhiteBlackReverseMode turns white/black reverse print mode on or off.
//
// Format:
//
//	ASCII:   GS B n
//	Hex:     0x1D 0x42 n
//	Decimal: 29 66 n
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
//	n: Controls white/black reverse print mode (LSB is used):
//	   LSB = 0 -> Reverse mode OFF
//	   LSB = 1 -> Reverse mode ON
//
// Notes:
//   - The white/black reverse print is effective for both single-byte and multi-byte code characters
//   - When reverse mode is turned on, characters are printed in white on a black background
//   - Reverse mode affects right-side character spacing set by ESC SP and
//     left/right spacing of multi-byte characters set by FS S
//   - Reverse mode does not affect line spacing or spaces skipped by HT, ESC $, or ESC \
//   - When underline mode is turned on, the printer does not underline
//     white/black reversed characters
//   - The setting remains in effect until ESC @ is executed, the printer is
//     reset, or the power is turned off
//
// Errors:
//
//	This function is safe and does not return errors.
func (c *Commands) SetWhiteBlackReverseMode(n ReverseMode) []byte {
	return []byte{shared.GS, 'B', byte(n)}
}

// SetSmoothingMode turns smoothing mode on or off.
//
// Format:
//
//	ASCII:   GS b n
//	Hex:     0x1D 0x62 n
//	Decimal: 29 98 n
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
//	n: Controls smoothing mode (LSB is used):
//	   LSB = 0 -> Smoothing mode OFF
//	   LSB = 1 -> Smoothing mode ON
//
// Notes:
//   - The smoothing mode is effective for quadruple-size or larger characters
//     (alphanumeric, Kana, multilingual, and user-defined characters)
//   - The setting remains in effect until ESC @ is executed, the printer is
//     reset, or the power is turned off
//
// Errors:
//
//	This function is safe and does not return errors.
func (c *Commands) SetSmoothingMode(n SmoothingMode) []byte {
	return []byte{shared.GS, 'b', byte(n)}
}
