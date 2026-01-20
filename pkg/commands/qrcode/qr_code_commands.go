package qrcode

import (
	"fmt"

	"github.com/adcondev/poster/pkg/commands/shared"
)

// SelectQRCodeModel selects the QR Code model to be used for encoding.
//
// Format:
//
//	ASCII:   GS ( k pL pH cn fn n1 n2
//	Hex:     0x1D 0x28 0x6B 0x04 0x00 0x31 0x41 n1 n2
//	Decimal: 29 40 107 4 0 49 65 n1 n2
//
// Range:
//
//	(pL + pH × 256) = 4
//	cn = 49
//	fn = 65
//	n1 = 49, 50, 51
//	n2 = 0
//
// Default:
//
//	n1 = 50, n2 = 0 (Model 2)
//
// Parameters:
//
//	n1: QR Code model selection:
//	   49 -> Model 1
//	   50 -> Model 2
//	   51 -> Micro QR Code
//	n2: Reserved (must be 0)
//
// Notes:
//   - Settings of this function affect the processing of GS ( k <Function 181> and GS ( k <Function 182>
//   - Settings remain effective until ESC @ is executed, the printer is reset, or the power is turned off
//   - Model 2 is the default and most commonly used QR Code model
//   - Model 1 is the original QR Code specification
//   - Micro QR Code is a smaller version for applications with limited space
//
// Errors:
//
//	Returns ErrQRModel if n1 is not a valid model value (49, 50, 51)
//	Returns ErrParameter if n2 is not 0
func (c *Commands) SelectQRCodeModel(n1 Model, n2 byte) ([]byte, error) {
	// Validate parameters
	if err := ValidateQRModel(n1); err != nil {
		return nil, err
	}
	if n2 != 0 {
		return nil, fmt.Errorf("%w: n2=%d", ErrParameter, n2)
	}

	// Build command
	return []byte{
		shared.GS, '(', 'k',
		0x04, 0x00, // pL, pH
		0x31,     // cn = 49
		65,       // fn = 65
		byte(n1), // model
		n2,       // reserved
	}, nil
}

// SetQRCodeModuleSize sets the size of the module (dot size) for QR Code symbols.
//
// Format:
//
//	ASCII:   GS ( k pL pH cn fn n
//	Hex:     0x1D 0x28 0x6B 0x03 0x00 0x31 0x43 n
//	Decimal: 29 40 107 3 0 49 67 n
//
// Range:
//
//	(pL + pH × 256) = 3
//	cn = 49
//	fn = 67
//	n = 1–16
//
// Default:
//
//	n = 3
//
// Parameters:
//
//	n: moduleSize of the module in dots (1–16)
//	   Each module is a square, so n represents both width and height
//
// Notes:
//   - Sets the size of the module for QR Code to n dots
//   - Settings of this function affect the processing of GS ( k <Function 181> and GS ( k <Function 182>
//   - Settings remain effective until ESC @ is executed, the printer is reset, or the power is turned off
//   - n = width of a module = height of a module (QR code modules are square)
//   - Larger module sizes produce larger, more easily scannable QR codes
//   - Smaller module sizes allow more data in limited print space but may be harder to scan
//
// Errors:
//
//	Returns ErrModuleSize if n is outside the valid range (1–16)
func (c *Commands) SetQRCodeModuleSize(n ModuleSize) ([]byte, error) {
	// Validate parameter
	if err := ValidateModuleSize(n); err != nil {
		return nil, err
	}

	// Build command
	return []byte{
		shared.GS, '(', 'k',
		0x03, 0x00, // pL, pH
		0x31,    // cn = 49
		67,      // fn = 67
		byte(n), // module size
	}, nil
}

// SetQRCodeErrorCorrectionLevel selects the error correction level for QR Code symbols.
//
// Format:
//
//	ASCII:   GS ( k pL pH cn fn n
//	Hex:     0x1D 0x28 0x6B 0x03 0x00 0x31 0x45 n
//	Decimal: 29 40 107 3 0 49 69 n
//
// Range:
//
//	(pL + pH × 256) = 3
//	cn = 49
//	fn = 69
//	n = 48–51
//
// Default:
//
//	n = 48 (Error correction level L)
//
// Parameters:
//
//	n: Error correction level:
//	   48 -> Level L (Recovery capacity ~7%)
//	   49 -> Level M (Recovery capacity ~15%)
//	   50 -> Level Q (Recovery capacity ~25%)
//	   51 -> Level H (Recovery capacity ~30%)
//
// Notes:
//   - Selects the error correction level for QR Code
//   - Settings of this function affect the processing of GS ( k <Function 181> and GS ( k <Function 182>
//   - QR Code employs Reed-Solomon error correction to generate a series of error correction codewords
//   - Settings remain effective until ESC @ is executed, the printer is reset, or the power is turned off
//   - Higher error correction levels allow the QR code to be read even if partially damaged or obscured
//   - Higher error correction levels increase the size of the QR code for the same data
//   - Level L is fastest to generate and produces smallest codes but has lowest error tolerance
//   - Level H provides highest error tolerance but produces largest codes
//
// Errors:
//
//	Returns ErrErrorCorrection if n is outside the valid range (48–51)
func (c *Commands) SetQRCodeErrorCorrectionLevel(n ErrorCorrection) ([]byte, error) {
	// Validate parameter
	if err := ValidateErrorCorrection(n); err != nil {
		return nil, err
	}

	// Build command
	return []byte{
		shared.GS, '(', 'k',
		0x03, 0x00, // pL, pH
		0x31,    // cn = 49
		69,      // fn = 69
		byte(n), // error correction level
	}, nil
}

// StoreQRCodeData stores the data in the QR Code symbol storage area for later encoding and printing.
//
// Format:
//
//	ASCII:   GS ( k pL pH cn fn m d1...dk
//	Hex:     0x1D 0x28 0x6B pL pH 0x31 0x50 0x30 d1...dk
//	Decimal: 29 40 107 pL pH 49 80 48 d1...dk
//
// Range:
//
//	(pL + pH × 256) = 4–7092
//	cn = 49
//	fn = 80
//	m = 48
//	d = 0–255
//	k = (pL + pH × 256) − 3
//
// Default:
//
//	None
//
// Parameters:
//
//	data: QR Code symbol data to store (d1...dk)
//	      The data can include the following categories:
//	      - Numerical Dithering: "0"–"9"
//	      - Alphanumeric Dithering: "0"–"9", "A"–"Z", SP, $, %, *, +, -, ., /, :
//	      - Kanji Dithering: Shift JIS values (shifted from JIS X 0208)
//	      - 8-Bit Byte Dithering: 0x00–0xFF
//
// Notes:
//   - Stores the QR Code symbol data in the symbol storage area
//   - The stored symbol data is encoded by GS ( k <Function 181> and GS ( k <Function 182>
//   - After encoding/printing, the symbol data in the storage area is retained
//   - k bytes of d1...dk are processed as symbol data
//   - Data should only include valid QR Code encodable characters (see Parameters section)
//   - Settings remain effective until one of the following occurs:
//   - GS ( k <Function 80/180/280/380/480> is executed (stores new data)
//   - ESC @ is executed
//   - The printer is reset or power is turned off
//   - Maximum data size is 7089 bytes (7092 - 3 header bytes)
//
// Errors:
//
//	Returns ErrDataLength if data length is less than 1 or greater than 7089 bytes
func (c *Commands) StoreQRCodeData(data []byte) ([]byte, error) {
	// Validate data length
	if err := ValidateDataLength(data); err != nil {
		return nil, err
	}

	// Calculate pL and pH for command length without converting to uint16
	// Total length = 3 (cn + fn + m) + data length
	totalLen := 3 + len(data) // int arithmetic avoids int->uint16 conversion
	pL := byte(totalLen & 0xFF)
	pH := byte((totalLen >> 8) & 0xFF)

	// Build command header
	cmd := make([]byte, 0, 8+len(data))
	cmd = append(cmd, []byte{
		shared.GS, '(', 'k',
		pL, pH, // length bytes
		0x31, // cn = 49
		80,   // fn = 80
		0x30, // m = 48
	}...)

	// Append data
	cmd = append(cmd, data...)

	return cmd, nil
}

// PrintQRCode encodes and prints the QR Code symbol data stored in the symbol storage area.
//
// Format:
//
//	ASCII:   GS ( k pL pH cn fn m
//	Hex:     0x1D 0x28 0x6B 0x03 0x00 0x31 0x51 m
//	Decimal: 29 40 107 3 0 49 81 m
//
// Range:
//
//	(pL + pH × 256) = 3
//	cn = 49
//	fn = 81
//	m = 48
//
// Default:
//
//	None
//
// Parameters:
//
//	None - All parameters are fixed for this function
//
// Notes:
//   - Encodes and prints the QR Code symbol data stored via GS ( k <Function 180>
//   - In Standard mode, use this function when the printer is "at the beginning of a line" or
//     "there is no data in the print buffer"
//   - Symbol size that exceeds the print area cannot be printed
//   - Printing fails if there are errors in the symbol storage area data:
//   - No data exists (GS ( k <Function 180> was not executed)
//   - Data exceeds what is allowed by the specified model and data compaction mode
//   - Automatically selects the best data compaction mode (in order of compaction rate):
//     1. Numerical mode
//     2. Alphanumeric mode
//     3. Kanji mode
//     4. 8-Bit Byte Dithering
//   - The following elements are added automatically during encoding:
//   - Position Detection Patterns
//   - Separators for Position Detection Patterns
//   - Timing Patterns
//   - Format Information
//   - Version Information
//   - Error Correction codewords (Reed-Solomon algorithm)
//   - Pad codeword
//   - Number of bits in Character Count Indicator
//   - Dithering Indicator
//   - Terminator
//   - Alignment Patterns (when model 2 is selected)
//   - Extension Patterns (when model 1 is selected)
//   - Printing is not affected by print modes (emphasized, double-strike, underline, white/black reverse,
//     90° clockwise-rotated), except for character size and upside-down print mode
//   - In Standard mode: executes paper feeding for the symbol, moves print position to left side of
//     printable area, and sets printer status to "Beginning of the line"
//   - In Page mode: stores symbol data in print buffer without printing, moves print position to
//     the next dot after the last data of the symbol
//   - The quiet zone is NOT included in the printing data - ensure adequate quiet zone space
//   - Symbol is printed at appropriate speed regardless of printer's maximum speed or GS ( K <Function 48> settings
//
// Errors:
//
//	This function is safe and does not return errors
func (c *Commands) PrintQRCode() []byte {
	// Build command
	return []byte{
		shared.GS, '(', 'k',
		0x03, 0x00, // pL, pH
		0x31, // cn = 49
		81,   // fn = 81
		0x30, // m = 48
	}
}

// GetQRCodeSize transmits the size information of the encoded QR Code symbol data in the symbol storage area.
//
// Format:
//
//	ASCII:   GS ( k pL pH cn fn m
//	Hex:     0x1D 0x28 0x6B 0x03 0x00 0x31 0x52 m
//	Decimal: 29 40 107 3 0 49 82 m
//
// Range:
//
//	(pL + pH × 256) = 3
//	cn = 49
//	fn = 82
//	m = 48
//
// Default:
//
//	None
//
// Parameters:
//
//	None - All parameters are fixed for this function
//
// Notes:
//   - Transmits the size information for the encoded QR Code symbol data stored via GS ( k <Function 180>
//   - In Standard mode, use this function when the printer is "at the beginning of a line" or
//     "there is no data in the print buffer"
//   - The printer response follows this format:
//   - Name: 1 byte (0x37, decimal 55)
//   - Identifier: 1 byte (0x36, decimal 54)
//   - Horizontal size: 1-5 bytes (ASCII digits 0-9)
//   - Separator: 1 byte (0x1F, decimal 31)
//   - Vertical size: 1-5 bytes (ASCII digits 0-9)
//   - Separator: 1 byte (0x1F, decimal 31)
//   - Fixed value: 1 byte (0x31, decimal 49)
//   - Separator: 1 byte (0x1F, decimal 31)
//   - Other information: 1 byte
//   - 0x30 (decimal 48): Printing is possible
//   - 0x31 (decimal 49): Printing is impossible
//   - NUL: 1 byte (0x00, decimal 0)
//   - Horizontal size and vertical size indicate the number of dots of the symbol
//   - ModuleSize values are converted to ASCII digits starting from the high order
//   - Example: If horizontal size is 120 dots, the field contains "120" (3 bytes: 0x31,0x32,0x30)
//   - ModuleSize information indicates the size of symbol that would be printed by GS ( k <Function 181>
//   - The quiet zone is NOT included in the size information
//   - This function does NOT print - it only transmits size information
//   - If "Other information" indicates "Printing is impossible" (49), possible causes and solutions:
//   - Cause: Data in print buffer (Standard mode)
//     Solution: Execute GS T or print commands (LF, CR, ESC J) to clear buffer
//   - Cause: Symbol is bigger than current print area
//     Solution: Expand print area (GS W, ESC W, ESC $), reduce module size (GS ( k <Function 167>),
//     or lower error correction level (GS ( k <Function 169>)
//   - Cause: Data in symbol storage area is too large
//     Solution: Send correct data (GS ( k <Function 180>), select other model (GS ( k <Function 165>),
//     or lower error correction level (GS ( k <Function 169>)
//   - Cause: No data in symbol storage area
//     Solution: Send data to storage area (GS ( k <Function 180>)
//
// Errors:
//
//	This function is safe and does not return errors
func (c *Commands) GetQRCodeSize() []byte {
	// Build command
	return []byte{
		shared.GS, '(', 'k',
		0x03, 0x00, // pL, pH
		0x31, // cn = 49
		82,   // fn = 82
		0x30, // m = 48
	}
}
