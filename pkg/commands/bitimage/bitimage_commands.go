package bitimage

import (
	"github.com/adcondev/poster/pkg/commands/shared"
)

// SelectBitImageMode stores bit image data in the print buffer with the specified mode.
//
// Format:
//
//	ASCII:   ESC * m nL nH d1...dk
//	Hex:     0x1B 0x2A m nL nH d1...dk
//	Decimal: 27 42 m nL nH d1...dk
//
// Range:
//
//	m = 0, 1, 32, 33
//	(nL + nH × 256) = 1–2400
//	d = 0–255
//
// Default:
//
//	None
//
// Parameters:
//
//	m: Bit image mode:
//	   0  -> 8-dot single-density (8 bits vertical, single-density horizontal)
//	   1  -> 8-dot double-density (8 bits vertical, double-density horizontal)
//	   32 -> 24-dot single-density (24 bits vertical, single-density horizontal)
//	   33 -> 24-dot double-density (24 bits vertical, double-density horizontal)
//	nL: Low byte of dot count in horizontal direction
//	nH: High byte of dot count in horizontal direction
//	data: Bit image data (column format) where 1 = print, 0 = no print
//
// Notes:
//   - The horizontal dot count is calculated as (nL + nH × 256)
//   - Data length (k) depends on the selected mode:
//   - For m = 0, 1: k = (nL + nH × 256)
//   - For m = 32, 33: k = (nL + nH × 256) × 3
//   - If the bit image data exceeds the number of dots to be printed on a line, the excess data is ignored
//   - The bit-image is not affected by most print modes (emphasized, double-strike, underline,
//     character size, white/black reverse printing, or 90° rotation), except for upside-down mode
//   - After printing a bit image, the printer processes normal data
//   - For multiple line bit images with aligned top/bottom parts, use unidirectional print mode (ESC U)
//   - Resolution by mode (dots per 25.4 mm / dots per inch):
//   - m = 0:  203/3 dpi vertical, 203/2 dpi horizontal
//   - m = 1:  203/3 dpi vertical, 203 dpi horizontal
//   - m = 32: 203 dpi vertical, 203/2 dpi horizontal
//   - m = 33: 203 dpi vertical, 203 dpi horizontal
//   - This command is typically used to print pictures or logos
//
// Errors:
//
//	Returns ErrBitImageMode if m is not a valid mode value (0, 1, 32, 33).
//	Returns ErrHorizontalDotsRange if the horizontal dot count is outside the valid range (1-2400).
//	Returns ErrDataLength if the provided data length doesn't match the expected length for the mode.
func (c *Commands) SelectBitImageMode(mode Mode, width uint16, data []byte) ([]byte, error) {
	if err := ValidateBitImageMode(mode); err != nil {
		return nil, err
	}
	if err := ValidateHorizontalDots(width); err != nil {
		return nil, err
	}

	expectedLength := CalculateDataLength(mode, width)
	if len(data) != expectedLength {
		return nil, ErrDataLength
	}

	nL, nH := shared.ToLittleEndian(width)

	cmd := make([]byte, 0, 5+len(data))
	cmd = append(cmd, shared.ESC, '*', byte(mode), nL, nH)
	cmd = append(cmd, data...)

	return cmd, nil
}

// PrintNVBitImage prints a non-volatile bit image with the specified mode.
//
// Deprecated: This command is obsolete and may not be supported by future printer models.
// Use NV graphics functions (GS ( L / GS 8 L: <Function 51> and <Functions 64-69>) instead,
// which offer multiple logo data, key code control, color specification, raster format,
// and capacity monitoring.
//
// Format:
//
//	ASCII:   FS p n m
//	Hex:     0x1C 0x70 n m
//	Decimal: 28 112 n m
//
// Range:
//
//	n = 1–255
//	m = 0–3, 48–51
//
// Default:
//
//	None
//
// Parameters:
//
//	n: NV bit image number
//	m: Print mode:
//	   0 or 48 -> Normal (scale x1 horizontal, x1 vertical)
//	   1 or 49 -> Double-width (scale x2 horizontal, x1 vertical)
//	   2 or 50 -> Double-height (scale x1 horizontal, x2 vertical)
//	   3 or 51 -> Quadruple (scale x2 horizontal, x2 vertical)
//
// Notes:
//   - This command is not effective when the NV bit image specified by n has not been defined
//   - In Standard mode, this command is effective only when there is no data in the print buffer
//     and the printer is at the beginning of the line
//   - In Page mode, the NV bit image is only stored in the print buffer and is not printed
//   - If the NV bit image exceeds one line of print area, the printer does not print it
//   - This command is not affected by print modes except upside-down print mode
//   - The command executes paper feed for the amount needed regardless of paper feed settings
//   - After printing, the print position is moved to left side of the printable area
//   - When printing, using unidirectional print mode (ESC U) enables alignment of top and bottom parts
//   - The NV bit image is defined by FS q (also an obsolete command)
//   - NV bit image is printed in the default dot density defined by GS ( L <Function 49>
//
// Errors:
//
//	Returns ErrInvalidMode if m is not a valid mode value (0-3, 48-51).
//	Returns ErrInvalidBitImageNumber if n is 0.
func (c *Commands) PrintNVBitImage(n byte, mode PrintMode) ([]byte, error) {
	if err := ValidateNVBitImageNumber(n); err != nil {
		return nil, err
	}
	if err := ValidatePrintMode(mode); err != nil {
		return nil, err
	}

	return []byte{shared.FS, 'p', n, byte(mode)}, nil
}

// DefineNVBitImage defines non-volatile bit images in the NV graphics area.
//
// Deprecated: This command is obsolete and may not be supported by future printer models.
// Use NV graphics functions (GS ( L / GS 8 L: <Function 51> and <Functions 64-69>) instead,
// which offer multiple logo data, key code control, color specification, raster format,
// capacity monitoring, and continuous processing.
//
// Format:
//
//	ASCII:   FS q n [xL xH yL yH d1...dk]1 ... [xL xH yL yH d1...dk]n
//	Hex:     0x1C 0x71 n [xL xH yL yH d1...dk]1 ... [xL xH yL yH d1...dk]n
//	Decimal: 28 113 n [xL xH yL yH d1...dk]1 ... [xL xH yL yH d1...dk]n
//
// Range:
//
//	n = 1–255
//	(xL + xH × 256) = 1–1023
//	(yL + yH × 256) = 1–288
//	d = 0–255
//
// Default:
//
//	None
//
// Parameters:
//
//	n: Number of defined NV bit images
//	images: Slice of NVBitImageData structures containing:
//	  width: Horizontal size in bytes (xL + xH × 256)
//	  height: Vertical size in bytes (yL + yH × 256)
//	  data: Definition data (column format) where 1 = print, 0 = no print
//
// Notes:
//   - NV bit image means a bit image defined in non-volatile memory, effective until next definition
//   - In Standard mode, this command is only effective at the beginning of the line
//   - This command cannot be included in a macro definition
//   - Data (d) specifies a bit printed to 1 and not printed to 0
//   - All previously defined NV bit images are canceled
//   - NV graphics data and NV bit image data cannot be defined simultaneously
//   - After processing this command, the printer executes a software reset
//   - During command execution:
//   - Paper feed button is disabled
//   - Real-time commands are not processed
//   - ASB status cannot be transmitted
//   - Data is written to non-volatile memory; do not turn off power during execution
//   - Printer is BUSY when writing to non-volatile memory
//   - Excessive use may destroy non-volatile memory; limit combined usage of non-volatile
//     write commands to 10 times per day
//   - The NV graphics area size is 256 KB
//   - For each image, data length k = (width × height × 8)
//
// Errors:
//
//	Returns ErrInvalidImageCount if n is 0 or doesn't match the length of images slice.
//	Returns ErrInvalidImageDimensions if width or height values are outside valid ranges.
//	Returns ErrInvalidDataLength if the data length doesn't match the expected size.
func (c *Commands) DefineNVBitImage(n byte, images []NVBitImageData) ([]byte, error) {
	if n == 0 || int(n) != len(images) {
		return nil, ErrInvalidImageCount
	}

	cmd := []byte{shared.FS, 'q', n}

	for _, img := range images {
		if err := ValidateNVImageDimensions(img.Width, img.Height); err != nil {
			return nil, err
		}

		expectedLength := int(img.Width) * int(img.Height) * 8
		if len(img.Data) != expectedLength {
			return nil, ErrInvalidDataLength
		}

		xL, xH := shared.ToLittleEndian(img.Width)
		yL, yH := shared.ToLittleEndian(img.Height)

		cmd = append(cmd, xL, xH, yL, yH)
		cmd = append(cmd, img.Data...)
	}

	return cmd, nil
}

// DefineDownloadedBitImage defines a bit image in the downloaded graphic area.
//
// Deprecated: This command is obsolete and may not be supported by future printer models.
// Use download graphics functions (GS ( L / GS 8 L: <Function 52> and <Functions 80-85>) instead,
// which offer multiple logo data, key code control, color selection, raster format,
// capacity monitoring, and other features.
//
// Format:
//
//	ASCII:   GS * x y d1...d(x × y × 8)
//	Hex:     0x1D 0x2A x y d1...d(x × y × 8)
//	Decimal: 29 42 x y d1...d(x × y × 8)
//
// Range:
//
//	x = 1–255
//	y = 1–48
//	x × y = 1–1536
//	d = 0–255
//
// Default:
//
//	None
//
// Parameters:
//
//	x: Number of bytes in horizontal direction
//	y: Number of bytes in vertical direction
//	data: Bit image data (column format) where 1 = print, 0 = no print
//
// Notes:
//   - The downloaded bit image is only available in current session until another definition is made,
//     ESC @ is executed, the printer is reset, or the power is turned off
//   - On some printer models, a downloaded bit image and a user-defined character (ESC &)
//     cannot be defined simultaneously
//   - When this command is executed, any previously defined user-defined characters are cleared
//   - When ESC & is executed, the downloaded bit image data is cleared
//   - The downloaded bit image is printed by GS / (also an obsolete command)
//   - Data is defined in column format with each bit representing a printed (1) or not printed (0) dot
//
// Errors:
//
//	Returns ErrInvalidHorizontalSize if x is 0.
//	Returns ErrInvalidVerticalSize if y is 0 or greater than 48.
//	Returns ErrInvalidDimensions if x × y exceeds 1536.
//	Returns ErrInvalidDataLength if the data length doesn't match x × y × 8.
func (c *Commands) DefineDownloadedBitImage(x, y byte, data []byte) ([]byte, error) {
	if err := ValidateDownloadedImageDimensions(x, y); err != nil {
		return nil, err
	}

	expectedLength := int(x) * int(y) * 8
	if len(data) != expectedLength {
		return nil, ErrInvalidDataLength
	}

	cmd := make([]byte, 0, 4+len(data))
	cmd = append(cmd, shared.GS, '*', x, y)
	cmd = append(cmd, data...)

	return cmd, nil
}

// PrintDownloadedBitImage prints a downloaded bit image with the specified scaling mode.
//
// Deprecated: This command is obsolete and may not be supported by future printer models.
// Use download graphics functions (GS ( L / GS 8 L: <Function 52> and <Functions 80-85>) instead,
// which offer multiple logo data, key code control, color selection, raster format,
// and capacity monitoring.
//
// Format:
//
//	ASCII:   GS / m
//	Hex:     0x1D 0x2F m
//	Decimal: 29 47 m
//
// Range:
//
//	m = 0–3, 48–51
//
// Default:
//
//	None
//
// Parameters:
//
//	m: Print mode:
//	   0 or 48 -> Normal (scale x1 horizontal, x1 vertical)
//	   1 or 49 -> Double-width (scale x2 horizontal, x1 vertical)
//	   2 or 50 -> Double-height (scale x1 horizontal, x2 vertical)
//	   3 or 51 -> Quadruple (scale x2 horizontal, x2 vertical)
//
// Notes:
//   - This command is ignored if a downloaded bit image has not been defined
//   - In Standard mode, this command is enabled only when there is no data in the print buffer
//     and the printer is at the beginning of the line
//   - In Page mode, this command develops the bit image in the print buffer but does not print it
//   - If a downloaded bit image exceeds one line, the excess data is not printed
//   - Downloaded bit image is printed in the default dot density defined by GS ( L <Function 49>
//   - The scales for width and height are specified by m. In Page mode with rotation,
//     the printer applies print area and dot density accordingly
//   - The downloaded bit image is not affected by most print modes, except upside-down print mode
//   - This command feeds paper as required, regardless of line spacing set by ESC 2 or ESC 3
//   - After printing, the print position moves to the left side of the printable area
//   - When printing, using unidirectional print mode (ESC U) enables alignment of top and bottom parts
//   - The downloaded bit image is defined by GS * (also an obsolete command)
//
// Errors:
//
//	Returns ErrInvalidMode if m is not a valid mode value (0-3, 48-51).
func (c *Commands) PrintDownloadedBitImage(mode PrintMode) ([]byte, error) {
	if err := ValidatePrintMode(mode); err != nil {
		return nil, err
	}

	return []byte{shared.GS, '/', byte(mode)}, nil
}

// PrintVariableVerticalSizeBitImage prints a bit image with variable vertical size.
//
// Deprecated: This command is obsolete and may not be supported by future printer models.
// Use graphics functions (GS ( L <Function 50> and GS ( L / GS 8 L <Function 113>) instead,
// which allow color selection and size specification in dot units.
//
// Format:
//
//	ASCII:   GS Q 0 m xL xH yL yH d1...dk
//	Hex:     0x1D 0x51 0x30 m xL xH yL yH d1...dk
//	Decimal: 29 81 48 m xL xH yL yH d1...dk
//
// Range:
//
//	m = 0–3, 48–51
//	(xL + xH × 256) = 1–4256
//	(yL + yH × 256) = 1–16
//	d = 0–255
//
// Default:
//
//	None
//
// Parameters:
//
//	m: Print mode:
//	   0 or 48 -> Normal (scale x1 horizontal, x1 vertical)
//	   1 or 49 -> Double-width (scale x2 horizontal, x1 vertical)
//	   2 or 50 -> Double-height (scale x1 horizontal, x2 vertical)
//	   3 or 51 -> Quadruple (scale x2 horizontal, x2 vertical)
//	width: Horizontal size in dots (xL + xH × 256)
//	height: Vertical size in bytes (yL + yH × 256)
//	data: Bit image data (column format) where 1 = print, 0 = no print
//
// Notes:
//   - In Standard mode, this command is effective only when there is no data in the print buffer
//     and the printer is at the beginning of the line
//   - Data (d) specifies a bit printed to 1 and not printed to 0
//   - If a bit image exceeds the print area for a line, the excess image data is ignored
//   - The scales for width and height are specified by m. In Page mode with rotation,
//     the printer applies print area and dot density accordingly
//   - Character size and all print modes do not affect printing of the bit image
//   - This command feeds paper as required, regardless of paper feed setting commands
//   - This command cannot be included in a macro definition
//   - After printing, the print position moves to the left side of the printable area
//   - Total data length k = (width × height)
//
// Errors:
//
//	Returns ErrInvalidMode if m is not a valid mode value (0-3, 48-51).
//	Returns ErrInvalidImageWidth if width is outside the valid range (1-4256).
//	Returns ErrInvalidImageHeight if height is outside the valid range (1-16).
//	Returns ErrInvalidDataLength if the data length doesn't match width × height.
func (c *Commands) PrintVariableVerticalSizeBitImage(mode PrintMode, width, height uint16, data []byte) ([]byte, error) {
	if err := ValidatePrintMode(mode); err != nil {
		return nil, err
	}
	if err := ValidateVariableImageDimensions(width, height, mode); err != nil {
		return nil, err
	}

	expectedLength := int(width) * int(height)
	if len(data) != expectedLength {
		return nil, ErrInvalidDataLength
	}

	xL, xH := shared.ToLittleEndian(width)
	yL, yH := shared.ToLittleEndian(height)

	cmd := make([]byte, 0, 8+len(data))
	cmd = append(cmd, shared.GS, 'Q', 0x30, byte(mode), xL, xH, yL, yH)
	cmd = append(cmd, data...)

	return cmd, nil
}

// PrintRasterBitImage prints a raster format bit image with the specified scaling mode.
//
// Deprecated: This command is obsolete and may not be supported by future printer models.
// Use graphics functions (GS ( L <Function 50> and GS ( L / GS 8 L <Function 112>) instead,
// which allow color selection and size specification in dot units.
//
// Format:
//
//	ASCII:   GS v 0 m xL xH yL yH d1...dk
//	Hex:     0x1D 0x76 0x30 m xL xH yL yH d1...dk
//	Decimal: 29 118 48 m xL xH yL yH d1...dk
//
// Range:
//
//	m = 0–3, 48–51
//	(xL + xH × 256) = 1–65535
//	(yL + yH × 256) = 1–2303
//	d = 0–255
//
// Default:
//
//	None
//
// Parameters:
//
//	m: Print mode:
//	   0 or 48 -> Normal (scale x1 horizontal, x1 vertical)
//	   1 or 49 -> Double-width (scale x2 horizontal, x1 vertical)
//	   2 or 50 -> Double-height (scale x1 horizontal, x2 vertical)
//	   3 or 51 -> Quadruple (scale x2 horizontal, x2 vertical)
//	width: Horizontal size in bytes (xL + xH × 256)
//	height: Vertical size in dots (yL + yH × 256)
//	data: Bit image data (raster format) where 1 = print, 0 = no print
//
// Notes:
//   - In Standard mode, this command is effective only when there is no data in the print buffer
//     and the printer is at the beginning of the line
//   - In Page mode, the bit image is only stored in the print buffer and is not printed
//   - Data (d) specifies a bit printed to 1 and not printed to 0
//   - If a raster bit image exceeds one line, the excess data is not printed
//   - The scales for width and height are specified by m. In Page mode with rotation,
//     the printer applies print area and dot density accordingly
//   - The raster bit image is not affected by print modes (emphasized, double-strike, underline,
//     character size, white/black reverse, upside-down, or 90° rotation)
//   - This command feeds paper as required, regardless of line spacing set by ESC 2 or ESC 3
//   - This command cannot be included in a macro definition
//   - After printing, the print position moves to the left side of the printable area
//   - The bit image is printed in the default dot density defined by GS ( L <Function 49>
//   - Total data length k = (width × height)
//
// Errors:
//
//	Returns ErrInvalidMode if m is not a valid mode value (0-3, 48-51).
//	Returns ErrInvalidImageWidth if width is outside the valid range (1-65535).
//	Returns ErrInvalidImageHeight if height is outside the valid range (1-2303).
//	Returns ErrInvalidDataLength if the data length doesn't match width × height.
func (c *Commands) PrintRasterBitImage(mode PrintMode, width, height uint16, data []byte) ([]byte, error) {
	if err := ValidatePrintMode(mode); err != nil {
		return nil, err
	}
	if err := ValidateRasterImageDimensions(width, height, mode); err != nil {
		return nil, err
	}

	expectedLength := int(width) * int(height)
	if len(data) != expectedLength {
		return nil, ErrInvalidDataLength
	}

	xL, xH := shared.ToLittleEndian(width)
	yL, yH := shared.ToLittleEndian(height)

	cmd := make([]byte, 0, 8+len(data))
	cmd = append(cmd, shared.GS, 'v', 0x30, byte(mode), xL, xH, yL, yH)
	cmd = append(cmd, data...)

	return cmd, nil
}
