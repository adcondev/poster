package barcode

import (
	"fmt"

	"github.com/adcondev/poster/pkg/commands/common"
)

// ============================================================================
// Barcode Command Implementations
// ============================================================================

// SelectHRICharacterPosition selects the print position of HRI (Human Readable Interpretation) characters.
//
// Format:
//
//	ASCII:   GS H n
//	Hex:     0x1D 0x48 n
//	Decimal: 29 72 n
//
// Range:
//
//	n = 0–3, 48–51
//
// Default:
//
//	n = 0 (Not printed)
//
// Parameters:
//
//	n: Selects the print position of HRI characters when printing a barcode:
//	   0 or 48 -> Not printed
//	   1 or 49 -> Above the barcode
//	   2 or 50 -> Below the barcode
//	   3 or 51 -> Both above and below the barcode
//
// Notes:
//   - HRI characters are printed using the font specified by GS f
//   - The setting persists until ESC @ (initialize), printer reset, or power-off
//
// Errors:
//
//	Returns ErrHRIPosition if n is outside the valid range (0-3, 48-51).
func (c *Commands) SelectHRICharacterPosition(n HRIPosition) ([]byte, error) {
	// Validate allowed values
	if err := ValidateHRIPosition(n); err != nil {
		return nil, err
	}
	return []byte{common.GS, 'H', byte(n)}, nil
}

// SelectFontForHRI selects the font used to print HRI (Human Readable Interpretation) characters.
//
// Format:
//
//	ASCII:   GS f n
//	Hex:     0x1D 0x66 n
//	Decimal: 29 102 n
//
// Range:
//
//	n = 0–4, 48–52, 97, 98 (model-dependent)
//
// Default:
//
//	n = 0
//
// Parameters:
//
//	n: Selects the font for HRI characters printed with barcodes:
//	   0 or 48  -> Font A
//	   1 or 49  -> Font B
//	   2 or 50  -> Font C
//	   3 or 51  -> Font D
//	   4 or 52  -> Font E
//	   97       -> Special font A (model dependent)
//	   98       -> Special font B (model dependent)
//
// Notes:
//   - The chosen font applies only to HRI characters
//   - HRI characters are printed at the position set by GS H
//   - Built-in font availability and metrics vary by model
//
// Errors:
//
//	Returns ErrHRIFont if n is not a valid font selector value.
func (c *Commands) SelectFontForHRI(n HRIFont) ([]byte, error) {
	// Validate allowed values
	if err := ValidateHRIFont(n); err != nil {
		return nil, err
	}
	return []byte{common.GS, 'f', byte(n)}, nil
}

// SetBarcodeHeight sets the barcode height.
//
// Format:
//
//	ASCII:   GS h n
//	Hex:     0x1D 0x68 n
//	Decimal: 29 104 n
//
// Range:
//
//	n = 1–255
//
// Default:
//
//	n = model dependent (example: 162)
//
// Parameters:
//
//	n: Sets the height of a barcode to n dots
//
// Notes:
//   - The units for n depend on the printer model
//   - This setting remains effective until ESC @ (initialize), printer reset, or power-off
//
// Errors:
//
//	Returns an error if height is outside the valid range (MinHeight to MaxHeight).
func (c *Commands) SetBarcodeHeight(height Height) ([]byte, error) {
	// Validate range using constants
	if height < MinHeight || height > MaxHeight {
		return nil, fmt.Errorf("%w: %d (allowed %d-%d)", ErrHeight, height, MinHeight, MaxHeight)
	}
	return []byte{common.GS, 'h', byte(height)}, nil
}

// SetBarcodeWidth sets the horizontal module width for barcodes.
//
// Format:
//
//	ASCII:   GS w n
//	Hex:     0x1D 0x77 n
//	Decimal: 29 119 n
//
// Range:
//
//	n = 2–6 (typical numeric values)
//	n = 68–76 (model-dependent alternate values)
//
// Default:
//
//	n = 3 (model-dependent)
//
// Parameters:
//
//	n: Sets the barcode module width (horizontal size)
//
// Notes:
//   - Units and exact effect depend on the printer model
//   - The setting remains effective until ESC @, printer reset, or power-off
//   - This affects the module width for various barcode types (see printer spec)
//   - The command does not validate model-specific allowed values; caller must
//     supply a value supported by the target printer
//
// Errors:
//
//	Returns an error if width is outside both standard (MinWidth-MaxWidth) and
//	extended (ExtendedMinWidth-ExtendedMaxWidth) ranges.
func (c *Commands) SetBarcodeWidth(width Width) ([]byte, error) {
	// Validate standard and extended ranges
	if (width >= MinWidth && width <= MaxWidth) ||
		(width >= ExtendedMinWidth && width <= ExtendedMaxWidth) {
		return []byte{common.GS, 'w', byte(width)}, nil
	}
	return nil, fmt.Errorf("%w: %d (allowed %d-%d or %d-%d)", ErrWidth, width, MinWidth, MaxWidth, ExtendedMinWidth, ExtendedMaxWidth)
}

// PrintBarcode builds the GS k command byte sequence to print a barcode.
//
// Format:
//
//	Function A (m = 0–6):
//	  GS k m d1...dk NUL
//	Function B (m = 65–79):
//	  GS k m n d1...dn
//
// Range:
//
//	Function A: m = 0–6, data terminated with NUL
//	Function B: m = 65–79, n = 1–255 (data length)
//
// Default:
//
//	None
//
// Parameters:
//
//	symbology: Barcode type selector (m value):
//	  Function A (classic forms):
//	    0  UPC-A     (k = 11 or 12 digits, numeric)
//	    1  UPC-E     (k = 6–8, 11, 12, numeric; k=7/8/11/12 must start with '0')
//	    2  JAN13/EAN13 (k = 12 or 13 digits, numeric)
//	    3  JAN8/EAN8   (k = 7 or 8 digits, numeric)
//	    4  CODE39    (k >= 1, 0–9 A–Z space $ % * + - . /, Start/stop '*' auto if omitted)
//	    5  ITF       (k >= 2 even, numeric; odd final digit ignored)
//	    6  CODABAR   (k >= 2, Start/stop A–D/a–d must be present; not auto-added)
//	  Function B (extended forms):
//	    65 UPC-A     (n = 11 or 12)
//	    66 UPC-E     (n = 6–8, 11, 12)
//	    67 EAN13     (n = 12 or 13)
//	    68 EAN8      (n = 7 or 8)
//	    69 CODE39    (1–255)
//	    70 ITF       (2–254 even)
//	    71 CODABAR   (2–255)
//	    72 CODE93    (1–255, start/stop + 2 check chars auto)
//	    73 CODE128   (2–255, d1= '{' (0x7B), d2= 65–67 => Set A/B/C; check digit auto)
//	    74 GS1-128   (2–255, FNC1, check digits auto; special SP,(,),* rules)
//	    75 GS1 DataBar Omnidirectional (n=13 digits; AI(01), check digit auto)
//	    76 GS1 DataBar Truncated (n=13)
//	    77 GS1 DataBar Limited (n=13; first digit constraint)
//	    78 GS1 DataBar Expanded (2–255; uses '{'+code for FNC1 / '(' / ')')
//	    79 CODE128 Auto (1–255; 0–255 byte data)
//	data: Barcode data bytes
//
// Notes:
//   - This function DOES NOT validate symbology-specific content or lengths;
//     caller must supply conforming data
//   - After printing, printer returns to "beginning of line" state
//   - Not affected by most text print modes (except upside-down)
//   - In Page mode, data is buffered (rendering per Page mode rules)
//   - Width exceeding print area is ignored/clipped by service
//
// Errors:
//
//	Returns ErrDataTooShort if data is empty.
//	Returns ErrSymbology if symbology is not recognized.
func (c *Commands) PrintBarcode(symbology Symbology, data []byte) ([]byte, error) {
	// Validate that data exists
	if len(data) == 0 {
		return nil, ErrDataTooShort
	}

	// TODO: Check if builders are better in composer package

	// Select construction based on symbology
	if symbology <= CODABAR {
		// Function A (NUL-terminated)
		return c.buildFunctionA(symbology, data)
	} else if symbology >= UPCAB && symbology <= CODE128Auto {
		// Function B (length-prefixed)
		return c.buildFunctionB(symbology, data)
	}

	return nil, ErrSymbology
}

// PrintBarcodeWithCodeSet prints a CODE128 or GS1-128 barcode with explicit code set.
//
// Format:
//
//	GS k m n '{' codeSet data...
//
// Range:
//
//	m = 73 (CODE128) or 74 (GS1-128)
//	codeSet = 65–67 (Code set A/B/C)
//	n = 2–255 (total data length including prefix)
//
// Default:
//
//	None
//
// Parameters:
//
//	symbology: Must be CODE128 (m=73) or GS1-128 (m=74)
//	codeSet: Code set selector (65=A, 66=B, 67=C)
//	data: Barcode data (without the '{' and code set prefix)
//
// Notes:
//   - Specialized method for CODE128 and GS1-128 that require code set specification
//   - The first two bytes of the barcode data must be: d1 = '{' (0x7B), d2 = 65-67
//   - Use this method when you need explicit control over CODE128 code sets
//   - For automatic code set selection, use PrintBarcode with CODE128Auto (m=79)
//
// Errors:
//
//	Returns an error if symbology is not CODE128 or GS1-128.
//	Returns ErrCode128Set if codeSet is invalid.
func (c *Commands) PrintBarcodeWithCodeSet(symbology Symbology, codeSet Code128Set, data []byte) ([]byte, error) {
	// Validate that symbology supports code sets
	if symbology != CODE128 && symbology != GS1128 {
		return nil, fmt.Errorf("%w: symbology %d does not support code sets", ErrSymbology, symbology)
	}

	// Validate code set
	if codeSet < Code128SetA || codeSet > Code128SetC {
		return nil, ErrCode128Set
	}

	if codeSet == Code128SetB {
		data = escCode128Data(data)
	}

	// Build data with code set prefix
	prefixedData := make([]byte, 0, len(data)+2)
	prefixedData = append(prefixedData, '{', byte(codeSet))
	prefixedData = append(prefixedData, data...)

	return c.buildFunctionB(symbology, prefixedData)
}

// Si usas CODE128 manual, considera sanitizar el input
func escCode128Data(data []byte) []byte {
	result := make([]byte, 0, len(data)*2)
	for _, b := range data {
		result = append(result, b)
		if b == '{' { // 0x7B
			result = append(result, '{') // EPSON requiere duplicarlo para imprimirlo como literal
		}
	}
	return result
}
