package bitimage

import (
	"github.com/adcondev/poster/pkg/commands/shared"
)

// SetGraphicsDotDensity sets the reference dot density for graphics and bit image data.
//
// Format:
//
//	ASCII:   GS ( L pL pH m fn x y
//	Hex:     0x1D 0x28 0x4C 0x04 0x00 0x30 fn x y
//	Decimal: 29 40 76 4 0 48 fn x y
//
// Range:
//
//	(pL + pH × 256) = 4
//	m = 48
//	fn = 1, 49
//	x = 50, 51
//	y = 50 (when x = 50), 51 (when x = 51)
//
// Default:
//
//	x = 50, y = 50 (180 dpi × 180 dpi)
//
// Parameters:
//
//	fn: Function code (must be 1 or 49)
//	x: Horizontal density selector:
//	   50 -> 180 dpi horizontal
//	   51 -> 360 dpi horizontal
//	y: Vertical density selector:
//	   50 -> 180 dpi vertical (only valid when x = 50)
//	   51 -> 360 dpi vertical (only valid when x = 51)
//
// Notes:
//   - This function affects the processing of the following graphics commands:
//   - GS ( L / GS 8 L <Function 67>: Define NV graphics data (raster format)
//   - GS ( L / GS 8 L <Function 68>: Define NV graphics data (column format)
//   - GS ( L / GS 8 L <Function 83>: Define download graphics data (raster format)
//   - GS ( L / GS 8 L <Function 84>: Define download graphics data (column format)
//   - GS ( L / GS 8 L <Function 112>: Store graphics data in print buffer (raster format)
//   - GS ( L / GS 8 L <Function 113>: Store graphics data in print buffer (column format)
//   - Also affects the following bit image commands:
//   - FS q: Print NV bit image (obsolete)
//   - GS /: Print download bit image (obsolete)
//   - GS Q 0: Print variable vertical size bit image (obsolete)
//   - GS v 0: Print raster bit image (obsolete)
//   - The parameter combinations determine the dot density:
//   - x=50, y=50: 180 dpi × 180 dpi
//   - x=51, y=51: 360 dpi × 360 dpi
//
// Errors:
//
//	Returns ErrInvalidFunctionCode if fn is not a valid function code (1, 49).
//	Returns ErrInvalidDensityValue if x is not a valid value (50, 51).
//	Returns ErrInvalidDensityCombination if y does not match x (must be x=y=50 or x=y=51).
func (c *GraphicsCommands) SetGraphicsDotDensity(fn FunctionCode, x, y DotDensity) ([]byte, error) {
	if err := ValidateDensityFunctionCode(fn); err != nil {
		return nil, err
	}
	if err := ValidateDotDensity(x, y); err != nil {
		return nil, err
	}

	return []byte{shared.GS, '(', 'L', 0x04, 0x00, 0x30, byte(fn), byte(x), byte(y)}, nil
}

// PrintBufferedGraphics prints the graphics data that was previously stored in the print buffer.
//
// Format:
//
//	ASCII:   GS ( L pL pH m fn
//	Hex:     0x1D 0x28 0x4C 0x02 0x00 0x30 fn
//	Decimal: 29 40 76 2 0 48 fn
//
// Range:
//
//	(pL + pH × 256) = 2
//	m = 48
//	fn = 2, 50
//
// Default:
//
//	None
//
// Parameters:
//
//	fn: Function code (must be 2 or 50)
//
// Notes:
//   - The printer cannot print when there is no graphics data stored in the print buffer
//   - GS ( L / GS 8 L <Function 112> and GS ( L / GS 8 L <Function 113> are used to store
//     graphics data in the print buffer
//   - This function cannot be used when Page mode is enabled
//   - Feeds the paper the required distance when printing graphics data, regardless of
//     line feed pitch settings
//   - After printing the graphics, the print position is moved to the left side of the
//     printable area and the printer enters "Beginning of the line" status
//
// Errors:
//
//	Returns ErrInvalidFunctionCode if fn is not a valid function code (2, 50).
func (c *GraphicsCommands) PrintBufferedGraphics(fn FunctionCode) ([]byte, error) {
	if err := ValidatePrintFunctionCode(fn); err != nil {
		return nil, err
	}

	return []byte{shared.GS, '(', 'L', 0x02, 0x00, 0x30, byte(fn)}, nil
}

// StoreRasterGraphicsInBuffer stores graphics data in raster format in the print buffer
// using the standard command format (for data size ≤ 65,535 bytes).
//
// Format:
//
//	ASCII:   GS ( L pL pH m fn a bx by c xL xH yL yH d1...dk
//	Hex:     0x1D 0x28 0x4C pL pH 0x30 0x70 a bx by c xL xH yL yH d1...dk
//	Decimal: 29 40 76 pL pH 48 112 a bx by c xL xH yL yH d1...dk
//
// Range:
//
//	(pL + pH × 256) = 11-65535
//	a = 48, 52
//	bx = 1, 2
//	by = 1, 2
//	c = 49-52
//	(xL + xH × 256) = 1-2400
//	(yL + yH × 256) = 1-2400 [when a=48, by=1]
//	(yL + yH × 256) = 1-1200 [when a=48, by=2]
//	(yL + yH × 256) = 1-600 [when a=52, by=1]
//	(yL + yH × 256) = 1-300 [when a=52, by=2]
//	d = 0-255
//
// Default:
//
//	None
//
// Parameters:
//
//	tone: Graphics tone:
//	      48 -> Monochrome (digital)
//	      52 -> Multiple tone
//	horizontalScale: Horizontal scaling factor:
//	      1 -> Normal width
//	      2 -> Double width
//	verticalScale: Vertical scaling factor:
//	      1 -> Normal height
//	      2 -> Double height
//	color: Color selection:
//	      49 (1) -> Color 1
//	      50 (2) -> Color 2
//	      51 (3) -> Color 3
//	      52 (4) -> Color 4
//	width: Width of the graphics data in dots (1-2400)
//	height: Height of the graphics data in dots (range depends on tone and verticalScale)
//	data: Graphics data in raster format (row by row) where 1 = print, 0 = no print
//
// Notes:
//   - This function stores graphics data in raster format (row by row)
//   - Use this function when data size is ≤ 65,535 bytes; use StoreRasterGraphicsInBufferLarge otherwise
//   - Use either this function or GS ( L <Function 113> (column format) but not both
//   - Use this function when the printer is at the beginning of the line or when there is
//     no data besides graphics data in the print buffer during Standard mode
//   - Multiple tone printing is not supported in Page mode
//   - This command cannot be included in a macro definition
//   - Graphics data that exceeds the print area for one line will not be printed
//   - Settings for text effects and font size do not affect graphics printing
//   - Print position does not change before and after this function
//   - Overprinting of data in multiple colors can be performed by changing the color (c)
//     and running this function again, but overprinting the same color is not possible
//   - Use GS ( L <Function 50> to print the graphics after storing in the buffer
//   - Graphics data is defined using the dot density set by GS ( L <Function 49>
//   - During processing of this function, real-time commands are not available
//   - Data length k = (int((width + 7)/8)) × height
//
// Errors:
//
//	Returns ErrInvalidTone if tone is not a valid value (48, 52).
//	Returns ErrInvalidScale if horizontalScale or verticalScale is not 1 or 2.
//	Returns ErrInvalidColor if color is not a valid value (49-52).
//	Returns ErrInvalidWidth if width is outside the valid range (1-2400).
//	Returns ErrInvalidHeight if height exceeds limits based on tone and verticalScale.
//	Returns ErrInvalidDataLength if data length doesn't match calculated size.
//	Returns ErrDataTooLarge if the total command size exceeds 65,535 bytes.
func (c *GraphicsCommands) StoreRasterGraphicsInBuffer(tone GraphicsTone, horizontalScale, verticalScale GraphicsScale,
	color GraphicsColor, width, height uint16, data []byte) ([]byte, error) {
	if err := ValidateGraphicsTone(tone); err != nil {
		return nil, err
	}
	if err := ValidateGraphicsScale(horizontalScale); err != nil {
		return nil, err
	}
	if err := ValidateGraphicsScale(verticalScale); err != nil {
		return nil, err
	}
	if err := ValidateGraphicsColor(color); err != nil {
		return nil, err
	}
	if err := ValidateRasterDimensions(width, height, tone, verticalScale); err != nil {
		return nil, err
	}

	expectedLength := calculateRasterDataSize(width, height)
	if len(data) != expectedLength {
		return nil, ErrInvalidDataLength
	}

	// Check total command size
	totalSize := 11 + len(data)
	if totalSize > MaxStandardCommandSize {
		return nil, ErrDataTooLarge
	}

	pL := byte((totalSize) & 0xFF)
	pH := byte(((totalSize) >> 8) & 0xFF)
	xL, xH := shared.ToLittleEndian(width)
	yL, yH := shared.ToLittleEndian(height)

	cmd := []byte{shared.GS, '(', 'L', pL, pH, 0x30, 0x70, byte(tone),
		byte(horizontalScale), byte(verticalScale), byte(color), xL, xH, yL, yH}
	cmd = append(cmd, data...)

	return cmd, nil
}

// StoreRasterGraphicsInBufferLarge stores graphics data in raster format in the print buffer
// using the extended command format (for data size > 65,535 bytes).
//
// Format:
//
//	ASCII:   GS 8 L p1 p2 p3 p4 m fn a bx by c xL xH yL yH d1...dk
//	Hex:     0x1D 0x38 0x4C p1 p2 p3 p4 0x30 0x70 a bx by c xL xH yL yH d1...dk
//	Decimal: 29 56 76 p1 p2 p3 p4 48 112 a bx by c xL xH yL yH d1...dk
//
// Range:
//
//	(p1 + p2 × 256 + p3 × 65536 + p4 × 16777216) = 11-4294967295
//	a = 48, 52
//	bx = 1, 2
//	by = 1, 2
//	c = 49-52
//	(xL + xH × 256) = 1-2400
//	(yL + yH × 256) = 1-2400 [when a=48, by=1]
//	(yL + yH × 256) = 1-1200 [when a=48, by=2]
//	(yL + yH × 256) = 1-600 [when a=52, by=1]
//	(yL + yH × 256) = 1-300 [when a=52, by=2]
//	d = 0-255
//
// Default:
//
//	None
//
// Parameters:
//
//	tone: Graphics tone:
//	      48 -> Monochrome (digital)
//	      52 -> Multiple tone
//	horizontalScale: Horizontal scaling factor:
//	      1 -> Normal width
//	      2 -> Double width
//	verticalScale: Vertical scaling factor:
//	      1 -> Normal height
//	      2 -> Double height
//	color: Color selection:
//	      49 (1) -> Color 1
//	      50 (2) -> Color 2
//	      51 (3) -> Color 3
//	      52 (4) -> Color 4
//	width: Width of the graphics data in dots (1-2400)
//	height: Height of the graphics data in dots (range depends on tone and verticalScale)
//	data: Graphics data in raster format (row by row) where 1 = print, 0 = no print
//
// Notes:
//   - This function stores graphics data in raster format (row by row)
//   - Use this function when data size is > 65,535 bytes; use StoreRasterGraphicsInBuffer otherwise
//   - Use either this function or GS 8 L <Function 113> (column format) but not both
//   - Use this function when the printer is at the beginning of the line or when there is
//     no data besides graphics data in the print buffer during Standard mode
//   - Multiple tone printing is not supported in Page mode
//   - This command cannot be included in a macro definition
//   - Graphics data that exceeds the print area for one line will not be printed
//   - Settings for text effects and font size do not affect graphics printing
//   - Print position does not change before and after this function
//   - Overprinting of data in multiple colors can be performed by changing the color (c)
//     and running this function again, but overprinting the same color is not possible
//   - Use GS ( L <Function 50> to print the graphics after storing in the buffer
//   - Graphics data is defined using the dot density set by GS ( L <Function 49>
//   - During processing of this function, real-time commands are not available
//   - Data length k = (int((width + 7)/8)) × height
//
// Errors:
//
//	Returns ErrInvalidTone if tone is not a valid value (48, 52).
//	Returns ErrInvalidScale if horizontalScale or verticalScale is not 1 or 2.
//	Returns ErrInvalidColor if color is not a valid value (49-52).
//	Returns ErrInvalidWidth if width is outside the valid range (1-2400).
//	Returns ErrInvalidHeight if height exceeds limits based on tone and verticalScale.
//	Returns ErrInvalidDataLength if data length doesn't match calculated size.
//	Returns ErrDataTooLarge if the total command size exceeds 4,294,967,295 bytes.
func (c *GraphicsCommands) StoreRasterGraphicsInBufferLarge(tone GraphicsTone, horizontalScale, verticalScale GraphicsScale,
	color GraphicsColor, width, height uint16, data []byte) ([]byte, error) {
	if err := ValidateGraphicsTone(tone); err != nil {
		return nil, err
	}
	if err := ValidateGraphicsScale(horizontalScale); err != nil {
		return nil, err
	}
	if err := ValidateGraphicsScale(verticalScale); err != nil {
		return nil, err
	}
	if err := ValidateGraphicsColor(color); err != nil {
		return nil, err
	}
	if err := ValidateRasterDimensions(width, height, tone, verticalScale); err != nil {
		return nil, err
	}

	expectedLength := calculateRasterDataSize(width, height)
	if len(data) != expectedLength {
		return nil, ErrInvalidDataLength
	}

	// Check total command size
	totalSize := uint32(11 + len(data)) //nolint:gosec
	if totalSize > MaxExtendedCommandSize {
		return nil, ErrDataTooLarge
	}

	p1, p2, p3, p4 := shared.ToLittleEndian32(totalSize)
	xL, xH := shared.ToLittleEndian(width)
	yL, yH := shared.ToLittleEndian(height)

	cmd := []byte{shared.GS, '8', 'L', p1, p2, p3, p4, 0x30, 0x70, byte(tone),
		byte(horizontalScale), byte(verticalScale), byte(color), xL, xH, yL, yH}
	cmd = append(cmd, data...)

	return cmd, nil
}

// StoreColumnGraphicsInBuffer stores graphics data in column format in the print buffer
// using the standard command format (for data size ≤ 65,535 bytes).
//
// Format:
//
//	ASCII:   GS ( L pL pH m fn a bx by c xL xH yL yH d1...dk
//	Hex:     0x1D 0x28 0x4C pL pH 0x30 0x71 0x30 bx by c xL xH yL yH d1...dk
//	Decimal: 29 40 76 pL pH 48 113 48 bx by c xL xH yL yH d1...dk
//
// Range:
//
//	(pL + pH × 256) = 11-65535
//	a = 48 (fixed)
//	bx = 1, 2
//	by = 1, 2
//	c = 49-51
//	(xL + xH × 256) = 1-2048
//	(yL + yH × 256) = 1-128
//	d = 0-255
//
// Default:
//
//	None
//
// Parameters:
//
//	horizontalScale: Horizontal scaling factor:
//	      1 -> Normal width
//	      2 -> Double width
//	verticalScale: Vertical scaling factor:
//	      1 -> Normal height
//	      2 -> Double height
//	color: Color selection:
//	      49 (1) -> Color 1
//	      50 (2) -> Color 2
//	      51 (3) -> Color 3
//	width: Width of the graphics data in dots (1-2048)
//	height: Height of the graphics data in dots (1-128)
//	data: Graphics data in column format (column by column) where 1 = print, 0 = no print
//
// Notes:
//   - This function stores graphics data in column format (column by column)
//   - Use this function when data size is ≤ 65,535 bytes; use StoreColumnGraphicsInBufferLarge otherwise
//   - Use either this function or GS ( L <Function 112> (raster format) but not both
//   - Use this function when the printer is at the beginning of the line or when there is
//     no data besides graphics data in the print buffer during Standard mode
//   - This command cannot be included in a macro definition
//   - Graphics data that exceeds the print area for one line will not be printed
//   - Settings for text effects and font size do not affect graphics printing
//   - Print position does not change before and after this function
//   - Overprinting of data in multiple colors can be performed by changing the color (c)
//     and running this function again, but overprinting the same color is not possible
//   - Use GS ( L <Function 50> to print the graphics after storing in the buffer
//   - Graphics data is defined using the dot density set by GS ( L <Function 49>
//   - During processing of this function, real-time commands are not available
//   - Data length k = (width) × (int((height + 7)/8))
//
// Errors:
//
//	Returns ErrInvalidScale if horizontalScale or verticalScale is not 1 or 2.
//	Returns ErrInvalidColor if color is not a valid value (49-51).
//	Returns ErrInvalidWidth if width is outside the valid range (1-2048).
//	Returns ErrInvalidHeight if height is outside the valid range (1-128).
//	Returns ErrInvalidDataLength if data length doesn't match calculated size.
//	Returns ErrDataTooLarge if the total command size exceeds 65,535 bytes.
func (c *GraphicsCommands) StoreColumnGraphicsInBuffer(horizontalScale, verticalScale GraphicsScale,
	color GraphicsColor, width, height uint16, data []byte) ([]byte, error) {
	if err := ValidateGraphicsScale(horizontalScale); err != nil {
		return nil, err
	}
	if err := ValidateGraphicsScale(verticalScale); err != nil {
		return nil, err
	}
	// Note: Column format only supports colors 1-3
	if color < Color1 || color > Color3 {
		return nil, ErrInvalidColor
	}
	if err := ValidateColumnDimensions(width, height); err != nil {
		return nil, err
	}

	expectedLength := calculateColumnDataSize(width, height)
	if len(data) != expectedLength {
		return nil, ErrInvalidDataLength
	}

	// Check total command size
	totalSize := 11 + len(data)
	if totalSize > MaxStandardCommandSize {
		return nil, ErrDataTooLarge
	}

	pL := byte((totalSize) & 0xFF)
	pH := byte(((totalSize) >> 8) & 0xFF)
	xL, xH := shared.ToLittleEndian(width)
	yL, yH := shared.ToLittleEndian(height)

	cmd := []byte{shared.GS, '(', 'L', pL, pH, 0x30, 0x71, 0x30,
		byte(horizontalScale), byte(verticalScale), byte(color), xL, xH, yL, yH}
	cmd = append(cmd, data...)

	return cmd, nil
}

// StoreColumnGraphicsInBufferLarge stores graphics data in column format in the print buffer
// using the extended command format (for data size > 65,535 bytes).
//
// Format:
//
//	ASCII:   GS 8 L p1 p2 p3 p4 m fn a bx by c xL xH yL yH d1...dk
//	Hex:     0x1D 0x38 0x4C p1 p2 p3 p4 0x30 0x71 0x30 bx by c xL xH yL yH d1...dk
//	Decimal: 29 56 76 p1 p2 p3 p4 48 113 48 bx by c xL xH yL yH d1...dk
//
// Range:
//
//	(p1 + p2 × 256 + p3 × 65536 + p4 × 16777216) = 11-4294967295
//	a = 48 (fixed)
//	bx = 1, 2
//	by = 1, 2
//	c = 49-51
//	(xL + xH × 256) = 1-2048
//	(yL + yH × 256) = 1-128
//	d = 0-255
//
// Default:
//
//	None
//
// Parameters:
//
//	horizontalScale: Horizontal scaling factor:
//	      1 -> Normal width
//	      2 -> Double width
//	verticalScale: Vertical scaling factor:
//	      1 -> Normal height
//	      2 -> Double height
//	color: Color selection:
//	      49 (1) -> Color 1
//	      50 (2) -> Color 2
//	      51 (3) -> Color 3
//	width: Width of the graphics data in dots (1-2048)
//	height: Height of the graphics data in dots (1-128)
//	data: Graphics data in column format (column by column) where 1 = print, 0 = no print
//
// Notes:
//   - This function stores graphics data in column format (column by column)
//   - Use this function when data size is > 65,535 bytes; use StoreColumnGraphicsInBuffer otherwise
//   - Use either this function or GS 8 L <Function 112> (raster format) but not both
//   - Use this function when the printer is at the beginning of the line or when there is
//     no data besides graphics data in the print buffer during Standard mode
//   - This command cannot be included in a macro definition
//   - Graphics data that exceeds the print area for one line will not be printed
//   - Settings for text effects and font size do not affect graphics printing
//   - Print position does not change before and after this function
//   - Overprinting of data in multiple colors can be performed by changing the color (c)
//     and running this function again, but overprinting the same color is not possible
//   - Use GS ( L <Function 50> to print the graphics after storing in the buffer
//   - Graphics data is defined using the dot density set by GS ( L <Function 49>
//   - During processing of this function, real-time commands are not available
//   - Data length k = (width) × (int((height + 7)/8))
//
// Errors:
//
//	Returns ErrInvalidScale if horizontalScale or verticalScale is not 1 or 2.
//	Returns ErrInvalidColor if color is not a valid value (49-51).
//	Returns ErrInvalidWidth if width is outside the valid range (1-2048).
//	Returns ErrInvalidHeight if height is outside the valid range (1-128).
//	Returns ErrInvalidDataLength if data length doesn't match calculated size.
//	Returns ErrDataTooLarge if the total command size exceeds 4,294,967,295 bytes.
func (c *GraphicsCommands) StoreColumnGraphicsInBufferLarge(horizontalScale, verticalScale GraphicsScale,
	color GraphicsColor, width, height uint16, data []byte) ([]byte, error) {
	if err := ValidateGraphicsScale(horizontalScale); err != nil {
		return nil, err
	}
	if err := ValidateGraphicsScale(verticalScale); err != nil {
		return nil, err
	}
	// Note: Column format only supports colors 1-3
	if color < Color1 || color > Color3 {
		return nil, ErrInvalidColor
	}
	if err := ValidateColumnDimensions(width, height); err != nil {
		return nil, err
	}

	expectedLength := calculateColumnDataSize(width, height)
	if len(data) != expectedLength {
		return nil, ErrInvalidDataLength
	}

	// Check total command size
	totalSize := uint32(11 + len(data)) //nolint:gosec
	if totalSize > MaxExtendedCommandSize {
		return nil, ErrDataTooLarge
	}

	p1, p2, p3, p4 := shared.ToLittleEndian32(totalSize)
	xL, xH := shared.ToLittleEndian(width)
	yL, yH := shared.ToLittleEndian(height)

	cmd := []byte{shared.GS, '8', 'L', p1, p2, p3, p4, 0x30, 0x71, 0x30,
		byte(horizontalScale), byte(verticalScale), byte(color), xL, xH, yL, yH}
	cmd = append(cmd, data...)

	return cmd, nil
}
