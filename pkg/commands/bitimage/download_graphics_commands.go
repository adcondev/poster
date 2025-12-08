package bitimage

import (
	"github.com/adcondev/poster/pkg/commands/common"
)

// GetDownloadGraphicsRemainingCapacity transmits the remaining capacity (unused area) in the download graphics area.
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
//	fn = 4, 52
//
// Default:
//
//	None
//
// Parameters:
//
//	fn: Function code (must be 4 or 52)
//
// Notes:
//   - Transmits the number of bytes of remaining memory (unused area) in the download graphics area
//   - ESC/POS Handshaking Protocol is not required for this function
//   - The printer response follows this format:
//   - Header: 1 byte (0x37, decimal 55)
//   - Identifier: 1 byte (0x32, decimal 50)
//   - Unused capacity: 1-8 bytes (ASCII digits 0-9)
//   - NUL: 1 byte (0x00, decimal 0)
//   - The unused capacity is the total byte count for the unused area converted to ASCII digits
//   - Example: If available capacity is 120 bytes, the capacity field contains "120" (3 bytes: 0x31,0x32,0x30)
//   - The control information for download graphics data is included in the capacity in use
//   - Do not use this function in conjunction with download bit images (GS *) and download characters (ESC &)
//
// Errors:
//
//	Returns ErrInvalidFunctionCode if fn is not a valid function code (4, 52).
func (c *DownloadGraphicsCommands) GetDownloadGraphicsRemainingCapacity(fn DLFunctionCode) ([]byte, error) {
	if err := ValidateDLRemainingFunctionCode(fn); err != nil {
		return nil, err
	}

	return []byte{common.GS, '(', 'L', 0x02, 0x00, 0x30, byte(fn)}, nil
}

// GetDownloadGraphicsKeyCodeList transmits the list of key codes for defined download graphics.
//
// Format:
//
//	ASCII:   GS ( L pL pH m fn d1 d2
//	Hex:     0x1D 0x28 0x4C 0x04 0x00 0x30 0x50 0x4B 0x43
//	Decimal: 29 40 76 4 0 48 80 75 67
//
// Range:
//
//	(pL + pH × 256) = 4
//	m = 48
//	fn = 80
//	d1 = 75 (Character "K")
//	d2 = 67 (Character "C")
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
//   - Transmits the defined download graphics key code list
//   - ESC/POS Handshaking Protocol is required for this function
//   - The printer response when key codes are present follows this format:
//   - Header: 1 byte (0x37, decimal 55)
//   - Identifier: 1 byte (0x73, decimal 115)
//   - Identification status: 1 byte
//   - 0x40 (decimal 64): No more data (all key codes sent or ≤40 items)
//   - 0x41 (decimal 65): More data remains (>40 items)
//   - Data: 0-80 bytes containing key codes (ASCII 32-126)
//   - NUL: 1 byte (0x00, decimal 0)
//   - The printer response when no key codes are present:
//   - Header: 1 byte (0x37, decimal 55)
//   - Identifier: 1 byte (0x73, decimal 115)
//   - Identification status: 1 byte (0x40, decimal 64)
//   - NUL: 1 byte (0x00, decimal 0)
//   - Do not use this function in conjunction with downloaded bit images (GS *) or user-defined characters (ESC &)
//
// Errors:
//
//	This function is safe and does not return errors.
func (c *DownloadGraphicsCommands) GetDownloadGraphicsKeyCodeList() []byte {
	return []byte{common.GS, '(', 'L', 0x04, 0x00, 0x30, 0x50, 'K', 'C'}
}

// DeleteAllDownloadGraphics deletes all download graphics data from the printer's memory.
//
// Format:
//
//	ASCII:   GS ( L pL pH m fn d1 d2 d3
//	Hex:     0x1D 0x28 0x4C 0x05 0x00 0x30 0x51 0x43 0x4C 0x52
//	Decimal: 29 40 76 5 0 48 81 67 76 82
//
// Range:
//
//	(pL + pH × 256) = 5
//	m = 48
//	fn = 81
//	d1 = 67 (Character "C")
//	d2 = 76 (Character "L")
//	d3 = 82 (Character "R")
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
//   - Deletes all download graphics data defined using:
//   - GS ( L / GS 8 L <Function 83>
//   - GS ( L / GS 8 L <Function 84>
//   - Deleted areas are designated as "Unused areas"
//   - All key codes are designated as undefined
//   - Use this function at the beginning of the line when the Standard mode is selected
//   - This command cannot be included in a macro definition
//   - When download graphics data is shared by multiple applications, executing this function
//     will delete ALL data being used by all applications. Use with caution.
//
// Errors:
//
//	This function is safe and does not return errors.
func (c *DownloadGraphicsCommands) DeleteAllDownloadGraphics() []byte {
	return []byte{common.GS, '(', 'L', 0x05, 0x00, 0x30, 0x51, 'C', 'L', 'R'}
}

// DeleteDownloadGraphicsByKeyCode deletes the specified download graphics data identified by key codes.
//
// Format:
//
//	ASCII:   GS ( L pL pH m fn kc1 kc2
//	Hex:     0x1D 0x28 0x4C 0x04 0x00 0x30 0x52 kc1 kc2
//	Decimal: 29 40 76 4 0 48 82 kc1 kc2
//
// Range:
//
//	(pL + pH × 256) = 4
//	m = 48
//	fn = 82
//	kc1 = 32–126
//	kc2 = 32–126
//
// Default:
//
//	None
//
// Parameters:
//
//	kc1: First byte of key code (ASCII printable character)
//	kc2: Second byte of key code (ASCII printable character)
//
// Notes:
//   - Deletes the download graphics data defined by the specified key codes (kc1 and kc2)
//   - Deleted areas are designated as "Unused areas"
//   - Deleted key codes are designated as undefined
//   - Use this function at the beginning of the line when the Standard mode is selected
//   - This command cannot be included in a macro definition
//
// Errors:
//
//	Returns ErrInvalidKeyCode if either kc1 or kc2 is outside the valid range (32–126).
func (c *DownloadGraphicsCommands) DeleteDownloadGraphicsByKeyCode(kc1, kc2 byte) ([]byte, error) {
	if err := ValidateKeyCode(kc1); err != nil {
		return nil, err
	}
	if err := ValidateKeyCode(kc2); err != nil {
		return nil, err
	}

	return []byte{common.GS, '(', 'L', 0x04, 0x00, 0x30, 0x52, kc1, kc2}, nil
}

// DefineDownloadGraphics defines the download graphics data (raster format) as a record specified by the key codes (kc1 and kc2) in the download graphics area.
//
// Format:
//
//	ASCII:   GS ( L pL pH m fn a kc1 kc2 b xL xH yL yH [c d1...dk]1...[c d1...dk]b
//	Hex:     0x1D 0x28 0x4C pL pH 0x30 0x53 a kc1 kc2 b xL xH yL yH [c d1...dk]1...[c d1...dk]b
//	Decimal: 29 40 76 pL pH 48 83 a kc1 kc2 b xL xH yL yH [c d1...dk]1...[c d1...dk]b
//
// Range:
//
//	(pL + pH × 256) = 12–65535
//	m = 48
//	fn = 83
//	kc1 = 32–126
//	kc2 = 32–126
//	(xL + xH × 256) = 1–8192
//	(yL + yH × 256) = 1–2304
//	d = 0–255
//	k = (int((xL + xH × 256) + 7) / 8) × (yL + yH × 256)
//	a = 48, 52
//	b = 1 (when a = 48)
//	b = 1–4 (when a = 52)
//	c = 49 (when a = 48)
//	c = 49–52 (when a = 52)
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
//	kc1: First byte of key code (ASCII printable character, 32–126)
//	kc2: Second byte of key code (ASCII printable character, 32–126)
//	numColors: Number of colors for the definition data (1–4, depending on tone)
//	width: Width of the graphics data in dots (1–8192)
//	height: Height of the graphics data in dots (1–2304)
//	data: Graphics data in raster format
//
// Notes:
//   - Download graphics data is defined in the printer's internal volatile memory (RAM).
//   - The defined data is valid until redefined, reset, or the printer is turned off.
//   - Definition data is stored as raster format, where 1 = print and 0 = no print.
//   - If the specified key code already exists in memory, it will be overwritten.
//   - This command cannot be included in a macro definition.
//   - Use GS ( L <Function 52> to check available memory capacity before defining data.
//   - When multiple colors are selected (b ≠ 1), specify one data group for each color.
//   - User-defined characters (ESC &) and downloaded bit images (GS *) are deleted when this command is used.
//
// Errors:
//
//	Returns ErrInvalidKeyCode if kc1 or kc2 is outside the valid range (32–126).
//	Returns ErrInvalidTone if tone is not 48 or 52.
//	Returns ErrInvalidNumColors if numColors is outside the valid range for the given tone.
//	Returns ErrInvalidDimensions if width or height is outside the valid range.
//	Returns ErrInvalidDataLength if the data length does not match the expected size.
func (c *DownloadGraphicsCommands) DefineDownloadGraphics(tone GraphicsTone, kc1, kc2 byte, width, height uint16,
	colorData []DLGraphicsColorData) ([]byte, error) {

	if err := ValidateGraphicsTone(tone); err != nil {
		return nil, err
	}
	if err := ValidateKeyCode(kc1); err != nil {
		return nil, err
	}
	if err := ValidateKeyCode(kc2); err != nil {
		return nil, err
	}
	if err := ValidateDLGraphicsDimensions(width, height); err != nil {
		return nil, err
	}
	if err := ValidateDLColorDataForTone(tone, colorData); err != nil {
		return nil, err
	}

	// Calculate expected data size and validate
	expectedDataSize := calculateDLRasterDataSize(width, height)
	for _, cd := range colorData {
		if len(cd.Data) != expectedDataSize {
			return nil, ErrInvalidDataLength
		}
	}

	// Calculate total command size
	totalSize := 9 + len(colorData)*(1+expectedDataSize) // 9 header bytes + color data
	if totalSize > MaxStandardCommandSize {
		return nil, ErrDataTooLarge
	}

	pL := byte((totalSize) & 0xFF)
	pH := byte(((totalSize) >> 8) & 0xFF)
	xL, xH := common.ToLittleEndian(width)
	yL, yH := common.ToLittleEndian(height)

	cmd := []byte{common.GS, '(', 'L', pL, pH, 0x30, 0x53, byte(tone), kc1, kc2, byte(len(colorData)), xL, xH, yL, yH}

	// Append color data groups
	for _, cd := range colorData {
		cmd = append(cmd, byte(cd.Color))
		cmd = append(cmd, cd.Data...)
	}

	return cmd, nil
}

// DefineDownloadGraphicsLarge defines the download graphics data (raster format) as a record specified by the key codes (kc1 and kc2) in the download graphics area.
//
// Format:
//
//	ASCII:   GS 8 L p1 p2 p3 p4 m fn a kc1 kc2 b xL xH yL yH [c d1...dk]1...[c d1...dk]b
//	Hex:     0x1D 0x38 0x4C p1 p2 p3 p4 0x30 0x53 a kc1 kc2 b xL xH yL yH [c d1...dk]1...[c d1...dk]b
//	Decimal: 29 56 76 p1 p2 p3 p4 48 83 a kc1 kc2 b xL xH yL yH [c d1...dk]1...[c d1...dk]b
//
// Range:
//
//	(p1 + p2 × 256 + p3 × 65536 + p4 × 16777216) = 12–4294967295
//	m = 48
//	fn = 83
//	kc1 = 32–126
//	kc2 = 32–126
//	(xL + xH × 256) = 1–8192
//	(yL + yH × 256) = 1–2304
//	d = 0–255
//	k = (int((xL + xH × 256) + 7) / 8) × (yL + yH × 256)
//	a = 48, 52
//	b = 1 (when a = 48)
//	b = 1–4 (when a = 52)
//	c = 49 (when a = 48)
//	c = 49–52 (when a = 52)
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
//	kc1: First byte of key code (ASCII printable character, 32–126)
//	kc2: Second byte of key code (ASCII printable character, 32–126)
//	numColors: Number of colors for the definition data (1–4, depending on tone)
//	width: Width of the graphics data in dots (1–8192)
//	height: Height of the graphics data in dots (1–2304)
//	data: Graphics data in raster format
//
// Notes:
//   - Download graphics data is defined in the printer's internal volatile memory (RAM).
//   - The defined data is valid until redefined, reset, or the printer is turned off.
//   - Definition data is stored as raster format, where 1 = print and 0 = no print.
//   - If the specified key code already exists in memory, it will be overwritten.
//   - This command cannot be included in a macro definition.
//   - Use GS ( L <Function 52> to check available memory capacity before defining data.
//   - When multiple colors are selected (b ≠ 1), specify one data group for each color.
//   - User-defined characters (ESC &) and downloaded bit images (GS *) are deleted when this command is used.
//
// Errors:
//
//	Returns ErrInvalidKeyCode if kc1 or kc2 is outside the valid range (32–126).
//	Returns ErrInvalidTone if tone is not 48 or 52.
//	Returns ErrInvalidNumColors if numColors is outside the valid range for the given tone.
//	Returns ErrInvalidDimensions if width or height is outside the valid range.
//	Returns ErrInvalidDataLength if the data length does not match the expected size.
func (c *DownloadGraphicsCommands) DefineDownloadGraphicsLarge(tone GraphicsTone, kc1, kc2 byte, width, height uint16,
	colorData []DLGraphicsColorData) ([]byte, error) {

	if err := ValidateGraphicsTone(tone); err != nil {
		return nil, err
	}
	if err := ValidateKeyCode(kc1); err != nil {
		return nil, err
	}
	if err := ValidateKeyCode(kc2); err != nil {
		return nil, err
	}
	if err := ValidateDLGraphicsDimensions(width, height); err != nil {
		return nil, err
	}
	if err := ValidateDLColorDataForTone(tone, colorData); err != nil {
		return nil, err
	}

	// Calculate expected data size and validate
	expectedDataSize := calculateDLRasterDataSize(width, height)
	for _, cd := range colorData {
		if len(cd.Data) != expectedDataSize {
			return nil, ErrInvalidDataLength
		}
	}

	// Calculate total command size
	totalSize := uint32(9 + len(colorData)*(1+expectedDataSize)) //nolint:gosec
	if totalSize > MaxExtendedCommandSize {
		return nil, ErrDataTooLarge
	}

	p1, p2, p3, p4 := common.ToLittleEndian32(totalSize)
	xL, xH := common.ToLittleEndian(width)
	yL, yH := common.ToLittleEndian(height)

	cmd := []byte{common.GS, '8', 'L', p1, p2, p3, p4, 0x30, 0x53, byte(tone), kc1, kc2, byte(len(colorData)), xL, xH, yL, yH}

	// Append color data groups
	for _, cd := range colorData {
		cmd = append(cmd, byte(cd.Color))
		cmd = append(cmd, cd.Data...)
	}

	return cmd, nil
}

// DefineDownloadGraphicsColumn defines the download graphics data (column format) as a record specified by the key codes (kc1 and kc2) in the download graphics area.
//
// Format:
//
//	ASCII:   GS ( L pL pH m fn a kc1 kc2 b xL xH yL yH [c d1...dk]1...[c d1...dk]b
//	Hex:     0x1D 0x28 0x4C pL pH 0x30 0x54 0x30 kc1 kc2 b xL xH yL yH [c d1...dk]1...[c d1...dk]b
//	Decimal: 29 40 76 pL pH 48 84 48 kc1 kc2 b xL xH yL yH [c d1...dk]1...[c d1...dk]b
//
// Range:
//
//	(pL + pH × 256) = 12–65535
//	m = 48
//	fn = 84
//	a = 48
//	kc1 = 32–126
//	kc2 = 32–126
//	(xL + xH × 256) = 1–8192
//	(yL + yH × 256) = 1–2304
//	d = 0–255
//	k = (xL + xH × 256) × (int((yL + yH × 256) + 7) / 8)
//	b = 1, 2 (when c = 49, 50)
//	b = 1 (when c = 51)
//	c = 49–51
//
// Default:
//
//	None
//
// Parameters:
//
//	kc1: First byte of key code (ASCII printable character, 32–126)
//	kc2: Second byte of key code (ASCII printable character, 32–126)
//	numColors: Number of colors for the definition data
//	width: Width of the graphics data in dots (1–8192)
//	height: Height of the graphics data in dots (1–2304)
//	color: Color specification for each group:
//	      49 -> Color 1
//	      50 -> Color 2
//	      51 -> Color 3
//	data: Graphics data in column format
//
// Notes:
//   - Download graphics data is defined in the printer's internal volatile memory (RAM).
//   - The defined data is valid until redefined, reset, or the printer is turned off.
//   - Definition data is stored as column format, where 1 = print and 0 = no print.
//   - If the specified key code already exists in memory, it will be overwritten.
//   - This command cannot be included in a macro definition.
//   - Use GS ( L <Function 52> to check available memory capacity before defining data.
//   - When multiple colors are selected (b ≠ 1), specify one data group for each color.
//   - User-defined characters (ESC &) and downloaded bit images (GS *) are deleted when this command is used.
//
// Errors:
//
//	Returns ErrInvalidKeyCode if kc1 or kc2 is outside the valid range (32–126).
//	Returns ErrInvalidNumColors if numColors is outside the valid range for the given color specification.
//	Returns ErrInvalidDimensions if width or height is outside the valid range.
//	Returns ErrInvalidDataLength if the data length does not match the expected size.
func (c *DownloadGraphicsCommands) DefineDownloadGraphicsColumn(kc1, kc2 byte, width, height uint16,
	colorData []DLGraphicsColorData) ([]byte, error) {

	if err := ValidateKeyCode(kc1); err != nil {
		return nil, err
	}
	if err := ValidateKeyCode(kc2); err != nil {
		return nil, err
	}
	if err := ValidateDLGraphicsDimensions(width, height); err != nil {
		return nil, err
	}
	if err := ValidateDLColumnColorData(colorData); err != nil {
		return nil, err
	}

	// Calculate expected data size and validate
	expectedDataSize := calculateDLColumnDataSize(width, height)
	for _, cd := range colorData {
		if len(cd.Data) != expectedDataSize {
			return nil, ErrInvalidDataLength
		}
	}

	// Calculate total command size
	totalSize := 10 + len(colorData)*(1+expectedDataSize) // 10 header bytes + color data
	if totalSize > MaxStandardCommandSize {
		return nil, ErrDataTooLarge
	}

	pL := byte((totalSize) & 0xFF)
	pH := byte(((totalSize) >> 8) & 0xFF)
	xL, xH := common.ToLittleEndian(width)
	yL, yH := common.ToLittleEndian(height)

	cmd := []byte{common.GS, '(', 'L', pL, pH, 0x30, 0x54, 0x30, kc1, kc2, byte(len(colorData)), xL, xH, yL, yH}

	// Append color data groups
	for _, cd := range colorData {
		cmd = append(cmd, byte(cd.Color))
		cmd = append(cmd, cd.Data...)
	}

	return cmd, nil
}

// DefineDownloadGraphicsColumnLarge defines the download graphics data (column format) as a record specified by the key codes (kc1 and kc2) in the download graphics area.
//
// Format:
//
//	ASCII:   GS 8 L p1 p2 p3 p4 m fn a kc1 kc2 b xL xH yL yH [c d1...dk]1...[c d1...dk]b
//	Hex:     0x1D 0x38 0x4C p1 p2 p3 p4 0x30 0x54 0x30 kc1 kc2 b xL xH yL yH [c d1...dk]1...[c d1...dk]b
//	Decimal: 29 56 76 p1 p2 p3 p4 48 84 48 kc1 kc2 b xL xH yL yH [c d1...dk]1...[c d1...dk]b
//
// Range:
//
//	(p1 + p2 × 256 + p3 × 65536 + p4 × 16777216) = 12–4294967295
//	m = 48
//	fn = 84
//	a = 48
//	kc1 = 32–126
//	kc2 = 32–126
//	(xL + xH × 256) = 1–8192
//	(yL + yH × 256) = 1–2304
//	d = 0–255
//	k = (xL + xH × 256) × (int((yL + yH × 256) + 7) / 8)
//	b = 1, 2 (when c = 49, 50)
//	b = 1 (when c = 51)
//	c = 49–51
//
// Default:
//
//	None
//
// Parameters:
//
//	kc1: First byte of key code (ASCII printable character, 32–126)
//	kc2: Second byte of key code (ASCII printable character, 32–126)
//	numColors: Number of colors for the definition data
//	width: Width of the graphics data in dots (1–8192)
//	height: Height of the graphics data in dots (1–2304)
//	color: Color specification for each group:
//	      49 -> Color 1
//	      50 -> Color 2
//	      51 -> Color 3
//	data: Graphics data in column format
//
// Notes:
//   - Download graphics data is defined in the printer's internal volatile memory (RAM).
//   - The defined data is valid until redefined, reset, or the printer is turned off.
//   - Definition data is stored as column format, where 1 = print and 0 = no print.
//   - If the specified key code already exists in memory, it will be overwritten.
//   - This command cannot be included in a macro definition.
//   - Use GS ( L <Function 52> to check available memory capacity before defining data.
//   - When multiple colors are selected (b ≠ 1), specify one data group for each color.
//   - User-defined characters (ESC &) and downloaded bit images (GS *) are deleted when this command is used.
//
// Errors:
//
//	Returns ErrInvalidKeyCode if kc1 or kc2 is outside the valid range (32–126).
//	Returns ErrInvalidNumColors if numColors is outside the valid range for the given color specification.
//	Returns ErrInvalidDimensions if width or height is outside the valid range.
//	Returns ErrInvalidDataLength if the data length does not match the expected size.
func (c *DownloadGraphicsCommands) DefineDownloadGraphicsColumnLarge(kc1, kc2 byte, width, height uint16,
	colorData []DLGraphicsColorData) ([]byte, error) {

	if err := ValidateKeyCode(kc1); err != nil {
		return nil, err
	}
	if err := ValidateKeyCode(kc2); err != nil {
		return nil, err
	}
	if err := ValidateDLGraphicsDimensions(width, height); err != nil {
		return nil, err
	}
	if err := ValidateDLColumnColorData(colorData); err != nil {
		return nil, err
	}

	// Calculate expected data size and validate
	expectedDataSize := calculateDLColumnDataSize(width, height)
	for _, cd := range colorData {
		if len(cd.Data) != expectedDataSize {
			return nil, ErrInvalidDataLength
		}
	}

	// Calculate total command size
	totalSize := uint32(10 + len(colorData)*(1+expectedDataSize)) //nolint:gosec
	if totalSize > MaxExtendedCommandSize {
		return nil, ErrDataTooLarge
	}

	p1, p2, p3, p4 := common.ToLittleEndian32(totalSize)
	xL, xH := common.ToLittleEndian(width)
	yL, yH := common.ToLittleEndian(height)

	cmd := []byte{common.GS, '8', 'L', p1, p2, p3, p4, 0x30, 0x54, 0x30, kc1, kc2, byte(len(colorData)), xL, xH, yL, yH}

	// Append color data groups
	for _, cd := range colorData {
		cmd = append(cmd, byte(cd.Color))
		cmd = append(cmd, cd.Data...)
	}

	return cmd, nil
}

// PrintDownloadGraphics prints the specified download graphics data identified by key codes.
//
// Format:
//
//	ASCII:   GS ( L pL pH m fn kc1 kc2 x y
//	Hex:     0x1D 0x28 0x4C 0x06 0x00 0x30 0x55 kc1 kc2 x y
//	Decimal: 29 40 76 6 0 48 85 kc1 kc2 x y
//
// Range:
//
//	(pL + pH × 256) = 6
//	m = 48
//	fn = 85
//	kc1 = 32–126
//	kc2 = 32–126
//	x = 1, 2
//	y = 1, 2
//
// Default:
//
//	None
//
// Parameters:
//
//	kc1: First byte of key code (ASCII printable character)
//	kc2: Second byte of key code (ASCII printable character)
//	x: Horizontal scaling factor:
//	      1 -> Normal width
//	      2 -> Double width
//	y: Vertical scaling factor:
//	      1 -> Normal height
//	      2 -> Double height
//
// Notes:
//   - Prints the download graphics data defined by the specified key codes (kc1 and kc2)
//   - The printer only prints if download graphics data corresponding to the specified key code is present
//   - Use this function at the beginning of the line or when the print buffer has no data in Standard mode
//   - In Page mode, storing image data to the buffer does not initiate printing
//   - Download graphics data exceeding the print area for one line will not be printed
//   - Scaling for width and height is specified by (x, y)
//   - In Page mode with 90° or 270° rotation, scaling applies to the respective orientation
//   - Text effects (bold, underline, orientation, etc., except upside-down) and font size do not affect the printing
//   - This command feeds paper the required distance, ignoring the line feed pitch settings
//   - After printing, the print position is moved to the left side of the printable area, and the printer enters "Beginning of the line" status
//   - Use ESC U to specify unidirectional print for better vertical alignment of printed text
//   - The model information (dot density, maximum print area, etc.) is the same as Function 69
//
// Errors:
//
//	Returns ErrInvalidKeyCode if either kc1 or kc2 is outside the valid range (32–126).
//	Returns ErrInvalidScaleFactor if x or y is outside the valid range (1, 2).
func (c *DownloadGraphicsCommands) PrintDownloadGraphics(kc1, kc2 byte, horizontalScale, verticalScale GraphicsScale) ([]byte, error) {
	if err := ValidateKeyCode(kc1); err != nil {
		return nil, err
	}
	if err := ValidateKeyCode(kc2); err != nil {
		return nil, err
	}
	if err := ValidateGraphicsScale(horizontalScale); err != nil {
		return nil, err
	}
	if err := ValidateGraphicsScale(verticalScale); err != nil {
		return nil, err
	}

	return []byte{common.GS, '(', 'L', 0x06, 0x00, 0x30, 0x55, kc1, kc2, byte(horizontalScale), byte(verticalScale)}, nil
}

// DefineBMPDownloadGraphics defines Windows BMP download graphics data in the printer's memory.
//
// Format:
//
//	ASCII:   GS D m fn a kc1 kc2 b c d1...dk
//	Hex:     0x1D 0x44 0x30 0x53 0x30 kc1 kc2 b 0x31 d1...dk
//	Decimal: 29 68 48 83 48 kc1 kc2 b 49 d1...dk
//
// Range:
//
//	m = 48
//	fn = 83
//	a = 48
//	kc1 = 32–126 (ASCII printable character)
//	kc2 = 32–126 (ASCII printable character)
//	b = 48, 52
//	c = 49
//	d = 0–255
//	k: Depends on the BMP file size
//
// Default:
//
//	None
//
// Parameters:
//
//	kc1: First byte of key code (ASCII printable character)
//	kc2: Second byte of key code (ASCII printable character)
//	tone: Tone of data to define:
//	      48 -> Monochrome (digital)
//	      52 -> Multiple tone
//	color: Color of the graphics data (fixed at 49 -> Color 1)
//	data: BMP definition data in raster format
//
// Notes:
//   - Converts Windows BMP data to the specified tone and defines download graphics data (raster format)
//   - The graphics data is associated with the specified key codes (kc1, kc2)
//   - If the number of definition data exceeds the remaining capacity of the download graphics memory, this function is ignored
//   - The [data value (k) + control information data value (16 bytes)] area of the download graphics data domain is used when this function is executed
//   - Tone (b parameter):
//   - 48: Monochrome (digital)
//   - 52: Multiple tone
//   - Color (c parameter) is fixed at 49 (Color 1)
//
// Errors:
//
//	Returns ErrInvalidKeyCode if either kc1 or kc2 is outside the valid range (32–126).
//	Returns ErrInvalidTone if tone is not a valid value (48, 52).
//	Returns ErrInvalidColor if color is not 49.
func (c *DownloadGraphicsCommands) DefineBMPDownloadGraphics(kc1, kc2 byte, tone GraphicsTone, bmpData []byte) ([]byte, error) {
	if err := ValidateKeyCode(kc1); err != nil {
		return nil, err
	}
	if err := ValidateKeyCode(kc2); err != nil {
		return nil, err
	}
	if err := ValidateGraphicsTone(tone); err != nil {
		return nil, err
	}
	if err := ValidateBMPData(bmpData); err != nil {
		return nil, err
	}

	cmd := []byte{common.GS, 'D', 0x30, 0x53, 0x30, kc1, kc2, byte(tone), 0x31}
	cmd = append(cmd, bmpData...)

	return cmd, nil
}
