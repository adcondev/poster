package bitimage

import (
	"github.com/adcondev/poster/pkg/commands/common"
)

// GetNVGraphicsCapacity transmits the entire capacity of the NV graphics area.
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
//	fn = 0, 48
//
// Default:
//
//	None
//
// Parameters:
//
//	fn: Function code (must be 0 or 48)
//
// Notes:
//   - Transmits the number of bytes in the NV graphics area
//   - ESC/POS Handshaking Protocol is not required for this function
//   - The printer response follows this format:
//   - Header: 1 byte (0x37, decimal 55)
//   - Identifier: 1 byte (0x30, decimal 48)
//   - Entire capacity: 1-8 bytes (ASCII digits 0-9)
//   - NUL: 1 byte (0x00, decimal 0)
//   - The entire capacity is the total byte count converted to ASCII digits
//   - Example: If capacity is 1200 bytes, the capacity field contains "1200" (4 bytes: 0x31,0x32,0x30,0x30)
//   - The entire capacity includes the area used for control information
//   - If the entire capacity is "0", the NV graphics function cannot be used
//
// Errors:
//
//	Returns ErrInvalidFunctionCode if fn is not a valid function code (0, 48).
func (c *NvGraphicsCommands) GetNVGraphicsCapacity(fn NVFunctionCode) ([]byte, error) {
	if err := ValidateNVCapacityFunctionCode(fn); err != nil {
		return nil, err
	}

	return []byte{common.GS, '(', 'L', 0x02, 0x00, 0x30, byte(fn)}, nil
}

// GetNVGraphicsRemainingCapacity transmits the remaining capacity (unused area) in the NV graphics area.
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
//	fn = 3, 51
//
// Default:
//
//	None
//
// Parameters:
//
//	fn: Function code (must be 3 or 51)
//
// Notes:
//   - Transmits the number of bytes of remaining memory (unused area) in the NV graphics area
//   - ESC/POS Handshaking Protocol is not required for this function
//   - The printer response follows this format:
//   - Header: 1 byte (0x37, decimal 55)
//   - Identifier: 1 byte (0x31, decimal 49)
//   - Unused capacity: 1-8 bytes (ASCII digits 0-9)
//   - NUL: 1 byte (0x00, decimal 0)
//   - The unused capacity is the total byte count for the unused area converted to ASCII digits
//   - Example: If available capacity is 120 bytes, the capacity field contains "120" (3 bytes: 0x31,0x32,0x30)
//   - The control information for NV graphics data is included in the capacity in use
//   - Do not use this function in conjunction with NV bit images (FS q)
//
// Errors:
//
//	Returns ErrInvalidFunctionCode if fn is not a valid function code (3, 51).
func (c *NvGraphicsCommands) GetNVGraphicsRemainingCapacity(fn NVFunctionCode) ([]byte, error) {
	if err := ValidateNVRemainingFunctionCode(fn); err != nil {
		return nil, err
	}

	return []byte{common.GS, '(', 'L', 0x02, 0x00, 0x30, byte(fn)}, nil
}

// GetNVGraphicsKeyCodeList transmits the list of key codes for defined NV graphics.
//
// Format:
//
//	ASCII:   GS ( L pL pH m fn d1 d2
//	Hex:     0x1D 0x28 0x4C 0x04 0x00 0x30 0x40 0x4B 0x43
//	Decimal: 29 40 76 4 0 48 64 75 67
//
// Range:
//
//	(pL + pH × 256) = 4
//	m = 48
//	fn = 64
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
//   - Transmits the defined NV graphics key code list
//   - ESC/POS Handshaking Protocol is required for this function
//   - The printer response when key codes are present follows this format:
//   - Header: 1 byte (0x37, decimal 55)
//   - Identifier: 1 byte (0x72, decimal 114)
//   - Identification status: 1 byte
//   - 0x40 (decimal 64): No more data (all key codes sent or ≤40 items)
//   - 0x41 (decimal 65): More data remains (>40 items)
//   - Data: 0-80 bytes containing key codes (ASCII 32-126)
//   - NUL: 1 byte (0x00, decimal 0)
//   - The printer response when no key codes are present:
//   - Header: 1 byte (0x37, decimal 55)
//   - Identifier: 1 byte (0x72, decimal 114)
//   - Identification status: 1 byte (0x40, decimal 64)
//   - NUL: 1 byte (0x00, decimal 0)
//   - Do not use this function in conjunction with NV bit images (FS q)
//
// Errors:
//
//	This function is safe and does not return errors.
func (c *NvGraphicsCommands) GetNVGraphicsKeyCodeList() []byte {
	return []byte{common.GS, '(', 'L', 0x04, 0x00, 0x30, 0x40, 'K', 'C'}
}

// DeleteAllNVGraphics deletes all NV graphics data from non-volatile memory.
//
// Format:
//
//	ASCII:   GS ( L pL pH m fn d1 d2 d3
//	Hex:     0x1D 0x28 0x4C 0x05 0x00 0x30 0x41 0x43 0x4C 0x52
//	Decimal: 29 40 76 5 0 48 65 67 76 82
//
// Range:
//
//	(pL + pH × 256) = 5
//	m = 48
//	fn = 65
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
//   - Deletes all NV graphics data that has been defined using Functions 67 or 68
//   - Deleted areas are designated as "Unused areas"
//   - All key codes are designated as undefined
//   - Use this function at the beginning of the line when the Standard mode is selected
//   - This command cannot be included in a macro definition
//   - When NV graphics data is being shared by multiple applications, executing this function
//     will delete ALL data being used by all applications. Use with caution.
//
// Errors:
//
//	This function is safe and does not return errors.
func (c *NvGraphicsCommands) DeleteAllNVGraphics() []byte {
	return []byte{common.GS, '(', 'L', 0x05, 0x00, 0x30, 0x41, 'C', 'L', 'R'}
}

// DeleteNVGraphicsByKeyCode deletes the specified NV graphics data identified by key codes.
//
// Format:
//
//	ASCII:   GS ( L pL pH m fn kc1 kc2
//	Hex:     0x1D 0x28 0x4C 0x04 0x00 0x30 0x42 kc1 kc2
//	Decimal: 29 40 76 4 0 48 66 kc1 kc2
//
// Range:
//
//	(pL + pH × 256) = 4
//	m = 48
//	fn = 66
//	kc1 = 32-126
//	kc2 = 32-126
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
//   - Deletes the NV graphics data defined by the specified key codes (kc1 and kc2)
//   - Deleted areas are designated as "Unused areas"
//   - Deleted key codes are designated as undefined
//   - Use this function at the beginning of the line when the Standard mode is selected
//   - This command cannot be included in a macro definition
//
// Errors:
//
//	Returns ErrInvalidKeyCode if either kc1 or kc2 is outside the valid range (32-126).
func (c *NvGraphicsCommands) DeleteNVGraphicsByKeyCode(kc1, kc2 byte) ([]byte, error) {
	if err := ValidateKeyCode(kc1); err != nil {
		return nil, err
	}
	if err := ValidateKeyCode(kc2); err != nil {
		return nil, err
	}

	return []byte{common.GS, '(', 'L', 0x04, 0x00, 0x30, 0x42, kc1, kc2}, nil
}

// DefineNVRasterGraphics defines NV graphics data in raster format in the printer's non-volatile memory
// using the standard command format (for data size ≤ 65,535 bytes).
//
// Format:
//
//	ASCII:   GS ( L pL pH m fn a kc1 kc2 b xL xH yL yH [c d1...dk]1...[c d1...dk]b
//	Hex:     0x1D 0x28 0x4C pL pH 0x30 0x43 a kc1 kc2 b xL xH yL yH [c d1...dk]1...[c d1...dk]b
//	Decimal: 29 40 76 pL pH 48 67 a kc1 kc2 b xL xH yL yH [c d1...dk]1...[c d1...dk]b
//
// Range:
//
//	(pL + pH × 256) = 12-65535
//	a = 48, 52
//	kc1 = 32-126
//	kc2 = 32-126
//	b = 1 (when a=48), 1-4 (when a=52)
//	(xL + xH × 256) = 1-8192
//	(yL + yH × 256) = 1-2304
//	c = 49-50 (when a=48), 49-52 (when a=52)
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
//	keyCode1: First byte of key code (ASCII printable character)
//	keyCode2: Second byte of key code (ASCII printable character)
//	width: Width of the graphics data in dots (1-8192)
//	height: Height of the graphics data in dots (1-2304)
//	colorData: Slice of NVGraphicsColorData structures, each containing:
//	  - Color: 49-50 (when tone=48), 49-52 (when tone=52)
//	  - Data: Graphics data in raster format (row by row) where 1 = print, 0 = no print
//
// Notes:
//   - NV graphics are image data stored in the printer's internal non-volatile memory
//   - NV graphics data definitions are valid until redefined or deleted
//   - Use either this function or GS ( L <Function 68> (column format) but not both
//   - If the specified key code already exists, the data will be overwritten
//   - Use this function at the beginning of the line when the Standard mode is selected
//   - This command cannot be included in a macro definition
//   - In cases where there is insufficient capacity, this function cannot be used
//   - Use GS ( L <Function 51> to confirm available capacity
//   - Data is stored using the dot density set by GS ( L <Function 49>
//   - For monochrome data (tone=48), only one color data group with color 1 or 2 is allowed
//   - For multiple tone data (tone=52), 1-4 color data groups with unique colors are allowed
//   - NV graphics data is printed using GS ( L <Function 69>
//   - This command cannot be used with NV bit image data (FS q)
//   - For optimal performance, use no more than 50 registered NV graphics
//   - ≤ 50 types: max execution time 60 seconds
//   - 100 types: execution time up to 120 seconds
//   - Data length k for each color = (int((width + 7)/8)) × height
//
// Errors:
//
//	Returns ErrInvalidTone if tone is not a valid value (48, 52).
//	Returns ErrInvalidKeyCode if keyCode1 or keyCode2 is outside the valid range (32-126).
//	Returns ErrInvalidWidth if width is outside the valid range (1-8192).
//	Returns ErrInvalidHeight if height is outside the valid range (1-2304).
//	Returns ErrInvalidColorCount if the number of color data groups doesn't match requirements.
//	Returns ErrInvalidColor if any color value is outside the valid range.
//	Returns ErrInvalidDataLength if any data length doesn't match the calculated size.
//	Returns ErrDataTooLarge if the total command size exceeds 65,535 bytes.
func (c *NvGraphicsCommands) DefineNVRasterGraphics(tone GraphicsTone, kc1, kc2 byte, width, height uint16,
	colorData []NVGraphicsColorData) ([]byte, error) {

	if err := ValidateGraphicsTone(tone); err != nil {
		return nil, err
	}
	if err := ValidateKeyCode(kc1); err != nil {
		return nil, err
	}
	if err := ValidateKeyCode(kc2); err != nil {
		return nil, err
	}
	if err := ValidateNVGraphicsDimensions(width, height); err != nil {
		return nil, err
	}
	if err := validateColorDataForTone(tone, colorData); err != nil {
		return nil, err
	}

	// Calculate expected data size and validate
	expectedDataSize := calculateNVRasterDataSize(width, height)
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

	cmd := []byte{common.GS, '(', 'L', pL, pH, 0x30, 0x43, byte(tone), kc1, kc2, byte(len(colorData)), xL, xH, yL, yH}

	// Append color data groups
	for _, cd := range colorData {
		cmd = append(cmd, byte(cd.Color))
		cmd = append(cmd, cd.Data...)
	}

	return cmd, nil
}

// DefineNVRasterGraphicsLarge defines NV graphics data in raster format in the printer's non-volatile memory
// using the extended command format (for data size > 65,535 bytes).
//
// Format:
//
//	ASCII:   GS 8 L p1 p2 p3 p4 m fn a kc1 kc2 b xL xH yL yH [c d1...dk]1...[c d1...dk]b
//	Hex:     0x1D 0x38 0x4C p1 p2 p3 p4 0x30 0x43 a kc1 kc2 b xL xH yL yH [c d1...dk]1...[c d1...dk]b
//	Decimal: 29 56 76 p1 p2 p3 p4 48 67 a kc1 kc2 b xL xH yL yH [c d1...dk]1...[c d1...dk]b
//
// Range:
//
//	(p1 + p2 × 256 + p3 × 65536 + p4 × 16777216) = 12-4294967295
//	a = 48, 52
//	kc1 = 32-126
//	kc2 = 32-126
//	b = 1 (when a=48), 1-4 (when a=52)
//	(xL + xH × 256) = 1-8192
//	(yL + yH × 256) = 1-2304
//	c = 49-50 (when a=48), 49-52 (when a=52)
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
//	keyCode1: First byte of key code (ASCII printable character)
//	keyCode2: Second byte of key code (ASCII printable character)
//	width: Width of the graphics data in dots (1-8192)
//	height: Height of the graphics data in dots (1-2304)
//	colorData: Slice of NVGraphicsColorData structures, each containing:
//	  - Color: 49-50 (when tone=48), 49-52 (when tone=52)
//	  - Data: Graphics data in raster format (row by row) where 1 = print, 0 = no print
//
// Notes:
//   - NV graphics are image data stored in the printer's internal non-volatile memory
//   - Use this function when data size is > 65,535 bytes; use DefineNVRasterGraphics otherwise
//   - NV graphics data definitions are valid until redefined or deleted
//   - Use either this function or GS 8 L <Function 68> (column format) but not both
//   - If the specified key code already exists, the data will be overwritten
//   - Use this function at the beginning of the line when the Standard mode is selected
//   - This command cannot be included in a macro definition
//   - In cases where there is insufficient capacity, this function cannot be used
//   - Use GS ( L <Function 51> to confirm available capacity
//   - Data is stored using the dot density set by GS ( L <Function 49>
//   - For monochrome data (tone=48), only one color data group with color 1 or 2 is allowed
//   - For multiple tone data (tone=52), 1-4 color data groups with unique colors are allowed
//   - NV graphics data is printed using GS ( L <Function 69>
//   - This command cannot be used with NV bit image data (FS q)
//   - For optimal performance, use no more than 50 registered NV graphics
//   - ≤ 50 types: max execution time 60 seconds
//   - 100 types: execution time up to 120 seconds
//   - Data length k for each color = (int((width + 7)/8)) × height
//
// Errors:
//
//	Returns ErrInvalidTone if tone is not a valid value (48, 52).
//	Returns ErrInvalidKeyCode if keyCode1 or keyCode2 is outside the valid range (32-126).
//	Returns ErrInvalidWidth if width is outside the valid range (1-8192).
//	Returns ErrInvalidHeight if height is outside the valid range (1-2304).
//	Returns ErrInvalidColorCount if the number of color data groups doesn't match requirements.
//	Returns ErrInvalidColor if any color value is outside the valid range.
//	Returns ErrInvalidDataLength if any data length doesn't match the calculated size.
//	Returns ErrDataTooLarge if the total command size exceeds 4,294,967,295 bytes.
func (c *NvGraphicsCommands) DefineNVRasterGraphicsLarge(tone GraphicsTone, kc1, kc2 byte, width, height uint16,
	colorData []NVGraphicsColorData) ([]byte, error) {

	if err := ValidateGraphicsTone(tone); err != nil {
		return nil, err
	}
	if err := ValidateKeyCode(kc1); err != nil {
		return nil, err
	}
	if err := ValidateKeyCode(kc2); err != nil {
		return nil, err
	}
	if err := ValidateNVGraphicsDimensions(width, height); err != nil {
		return nil, err
	}
	if err := validateColorDataForTone(tone, colorData); err != nil {
		return nil, err
	}

	// Calculate expected data size and validate
	expectedDataSize := calculateNVRasterDataSize(width, height)
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

	cmd := []byte{common.GS, '8', 'L', p1, p2, p3, p4, 0x30, 0x43, byte(tone), kc1, kc2, byte(len(colorData)), xL, xH, yL, yH}

	// Append color data groups
	for _, cd := range colorData {
		cmd = append(cmd, byte(cd.Color))
		cmd = append(cmd, cd.Data...)
	}

	return cmd, nil
}

// DefineNVColumnGraphics defines NV graphics data in column format in the printer's non-volatile memory
// using the standard command format (for data size ≤ 65,535 bytes).
//
// Format:
//
//	ASCII:   GS ( L pL pH m fn a kc1 kc2 b xL xH yL yH [c d1...dk]1...[c d1...dk]b
//	Hex:     0x1D 0x28 0x4C pL pH 0x30 0x44 0x30 kc1 kc2 b xL xH yL yH [c d1...dk]1...[c d1...dk]b
//	Decimal: 29 40 76 pL pH 48 68 48 kc1 kc2 b xL xH yL yH [c d1...dk]1...[c d1...dk]b
//
// Range:
//
//	(pL + pH × 256) = 12-65535
//	a = 48 (fixed)
//	kc1 = 32-126
//	kc2 = 32-126
//	b = 1, 2 (when c=49,50), 1 (when c=51)
//	(xL + xH × 256) = 1-8192
//	(yL + yH × 256) = 1-2304
//	c = 49-51
//	d = 0-255
//
// Default:
//
//	None
//
// Parameters:
//
//	keyCode1: First byte of key code (ASCII printable character)
//	keyCode2: Second byte of key code (ASCII printable character)
//	width: Width of the graphics data in dots (1-8192)
//	height: Height of the graphics data in dots (1-2304)
//	colorData: Slice of NVGraphicsColorData structures, each containing:
//	  - Color: 49-51 (colors 1-3)
//	  - Data: Graphics data in column format (column by column) where 1 = print, 0 = no print
//
// Notes:
//   - NV graphics are image data stored in the printer's internal non-volatile memory
//   - NV graphics data definitions are valid until redefined or deleted
//   - Use either this function or GS ( L <Function 67> (raster format) but not both
//   - If the specified key code already exists, the data will be overwritten
//   - Use this function at the beginning of the line when the Standard mode is selected
//   - This command cannot be included in a macro definition
//   - In cases where there is insufficient capacity, this function cannot be used
//   - Use GS ( L <Function 51> to confirm available capacity
//   - Data is stored using the dot density set by GS ( L <Function 49>
//   - Color restrictions:
//   - For colors 1 and 2 (49, 50): can use 1 or 2 color data groups
//   - For color 3 (51): only 1 color data group allowed
//   - NV graphics data is printed using GS ( L <Function 69>
//   - This command cannot be used with NV bit image data (FS q)
//   - Data length k for each color = (width) × (int((height + 7)/8))
//
// Errors:
//
//	Returns ErrInvalidKeyCode if keyCode1 or keyCode2 is outside the valid range (32-126).
//	Returns ErrInvalidWidth if width is outside the valid range (1-8192).
//	Returns ErrInvalidHeight if height is outside the valid range (1-2304).
//	Returns ErrInvalidColorCount if the number of color data groups doesn't match requirements.
//	Returns ErrInvalidColor if any color value is outside the valid range or has invalid restrictions.
//	Returns ErrInvalidDataLength if any data length doesn't match the calculated size.
//	Returns ErrDataTooLarge if the total command size exceeds 65,535 bytes.
func (c *NvGraphicsCommands) DefineNVColumnGraphics(kc1, kc2 byte, width, height uint16,
	colorData []NVGraphicsColorData) ([]byte, error) {

	if err := ValidateKeyCode(kc1); err != nil {
		return nil, err
	}
	if err := ValidateKeyCode(kc2); err != nil {
		return nil, err
	}
	if err := ValidateNVGraphicsDimensions(width, height); err != nil {
		return nil, err
	}
	if err := validateColumnColorData(colorData); err != nil {
		return nil, err
	}

	// Calculate expected data size and validate
	expectedDataSize := calculateNVColumnDataSize(width, height)
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

	cmd := []byte{common.GS, '(', 'L', pL, pH, 0x30, 0x44, 0x30, kc1, kc2, byte(len(colorData)), xL, xH, yL, yH}

	// Append color data groups
	for _, cd := range colorData {
		cmd = append(cmd, byte(cd.Color))
		cmd = append(cmd, cd.Data...)
	}

	return cmd, nil
}

// DefineNVColumnGraphicsLarge defines NV graphics data in column format in the printer's non-volatile memory
// using the extended command format (for data size > 65,535 bytes).
//
// Format:
//
//	ASCII:   GS 8 L p1 p2 p3 p4 m fn a kc1 kc2 b xL xH yL yH [c d1...dk]1...[c d1...dk]b
//	Hex:     0x1D 0x38 0x4C p1 p2 p3 p4 0x30 0x44 0x30 kc1 kc2 b xL xH yL yH [c d1...dk]1...[c d1...dk]b
//	Decimal: 29 56 76 p1 p2 p3 p4 48 68 48 kc1 kc2 b xL xH yL yH [c d1...dk]1...[c d1...dk]b
//
// Range:
//
//	(p1 + p2 × 256 + p3 × 65536 + p4 × 16777216) = 12-4294967295
//	a = 48 (fixed)
//	kc1 = 32-126
//	kc2 = 32-126
//	b = 1, 2 (when c=49,50), 1 (when c=51)
//	(xL + xH × 256) = 1-8192
//	(yL + yH × 256) = 1-2304
//	c = 49-51
//	d = 0-255
//
// Default:
//
//	None
//
// Parameters:
//
//	keyCode1: First byte of key code (ASCII printable character)
//	keyCode2: Second byte of key code (ASCII printable character)
//	width: Width of the graphics data in dots (1-8192)
//	height: Height of the graphics data in dots (1-2304)
//	colorData: Slice of NVGraphicsColorData structures, each containing:
//	  - Color: 49-51 (colors 1-3)
//	  - Data: Graphics data in column format (column by column) where 1 = print, 0 = no print
//
// Notes:
//   - NV graphics are image data stored in the printer's internal non-volatile memory
//   - Use this function when data size is > 65,535 bytes; use DefineNVColumnGraphics otherwise
//   - NV graphics data definitions are valid until redefined or deleted
//   - Use either this function or GS 8 L <Function 67> (raster format) but not both
//   - If the specified key code already exists, the data will be overwritten
//   - Use this function at the beginning of the line when the Standard mode is selected
//   - This command cannot be included in a macro definition
//   - In cases where there is insufficient capacity, this function cannot be used
//   - Use GS ( L <Function 51> to confirm available capacity
//   - Data is stored using the dot density set by GS ( L <Function 49>
//   - Color restrictions:
//   - For colors 1 and 2 (49, 50): can use 1 or 2 color data groups
//   - For color 3 (51): only 1 color data group allowed
//   - NV graphics data is printed using GS ( L <Function 69>
//   - This command cannot be used with NV bit image data (FS q)
//   - Data length k for each color = (width) × (int((height + 7)/8))
//
// Errors:
//
//	Returns ErrInvalidKeyCode if keyCode1 or keyCode2 is outside the valid range (32-126).
//	Returns ErrInvalidWidth if width is outside the valid range (1-8192).
//	Returns ErrInvalidHeight if height is outside the valid range (1-2304).
//	Returns ErrInvalidColorCount if the number of color data groups doesn't match requirements.
//	Returns ErrInvalidColor if any color value is outside the valid range or has invalid restrictions.
//	Returns ErrInvalidDataLength if any data length doesn't match the calculated size.
//	Returns ErrDataTooLarge if the total command size exceeds 4,294,967,295 bytes.
func (c *NvGraphicsCommands) DefineNVColumnGraphicsLarge(kc1, kc2 byte, width, height uint16,
	colorData []NVGraphicsColorData) ([]byte, error) {

	if err := ValidateKeyCode(kc1); err != nil {
		return nil, err
	}
	if err := ValidateKeyCode(kc2); err != nil {
		return nil, err
	}
	if err := ValidateNVGraphicsDimensions(width, height); err != nil {
		return nil, err
	}
	if err := validateColumnColorData(colorData); err != nil {
		return nil, err
	}

	// Calculate expected data size and validate
	expectedDataSize := calculateNVColumnDataSize(width, height)
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

	cmd := []byte{common.GS, '8', 'L', p1, p2, p3, p4, 0x30, 0x44, 0x30, kc1, kc2, byte(len(colorData)), xL, xH, yL, yH}

	// Append color data groups
	for _, cd := range colorData {
		cmd = append(cmd, byte(cd.Color))
		cmd = append(cmd, cd.Data...)
	}

	return cmd, nil
}

// PrintNVGraphics prints NV graphics data defined by the specified key codes.
//
// Format:
//
//	ASCII:   GS ( L pL pH m fn kc1 kc2 x y
//	Hex:     0x1D 0x28 0x4C 0x06 0x00 0x30 0x45 kc1 kc2 x y
//	Decimal: 29 40 76 6 0 48 69 kc1 kc2 x y
//
// Range:
//
//	(pL + pH × 256) = 6
//	m = 48
//	fn = 69
//	kc1 = 32-126
//	kc2 = 32-126
//	x = 1, 2
//	y = 1, 2
//
// Default:
//
//	None
//
// Parameters:
//
//	keyCode1: First byte of key code (ASCII printable character)
//	keyCode2: Second byte of key code (ASCII printable character)
//	horizontalScale: Horizontal scaling factor:
//	      1 -> Normal width
//	      2 -> Double width
//	verticalScale: Vertical scaling factor:
//	      1 -> Normal height
//	      2 -> Double height
//
// Notes:
//   - This function prints NV graphics data defined using GS ( L / GS 8 L <Function 67> or <Function 68>
//   - The printer only prints when NV graphics data corresponding to the specified key code exists
//   - Use this function at the beginning of the line or when there is no data in the print buffer
//     during Standard mode
//   - In Page mode, the data is only stored in the print buffer and not printed
//   - NV graphics data that exceeds the print area for one line will not be printed
//   - Settings for text effects and font size do not affect graphics printing (except upside-down)
//   - The printer feeds paper as required, regardless of line feed pitch settings
//   - After printing the graphics, the print position moves to the left side of the printable area
//   - When printing NV graphic files, use ESC U to specify unidirectional print for vertical alignment
//   - Dot density (at normal scale):
//   - Horizontal: 203 dpi
//   - Vertical: 203 dpi
//   - Maximum print area (at normal scale) depends on paper width:
//   - 40mm paper: 280 dots horizontal
//   - 58mm paper: 420 dots horizontal
//   - 80mm paper: 576 dots horizontal
//   - All paper: 1476 dots vertical
//
// Errors:
//
//	Returns ErrInvalidKeyCode if keyCode1 or keyCode2 is outside the valid range (32-126).
//	Returns ErrInvalidScale if horizontalScale or verticalScale is not 1 or 2.
func (c *NvGraphicsCommands) PrintNVGraphics(kc1, kc2 byte, horizontalScale, verticalScale GraphicsScale) ([]byte, error) {
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

	return []byte{common.GS, '(', 'L', 0x06, 0x00, 0x30, 0x45, kc1, kc2, byte(horizontalScale), byte(verticalScale)}, nil
}

// DefineWindowsBMPNVGraphics converts a Windows BMP file to NV graphics data and stores it
// in the printer's non-volatile memory.
//
// Format:
//
//	ASCII:   GS D m fn a kc1 kc2 b c d1...dk
//	Hex:     0x1D 0x44 0x30 0x43 0x30 kc1 kc2 b 0x31 d1...dk
//	Decimal: 29 68 48 67 48 kc1 kc2 b 49 d1...dk
//
// Range:
//
//	m = 48 (fixed)
//	fn = 67 (fixed)
//	a = 48 (fixed)
//	kc1 = 32-126
//	kc2 = 32-126
//	b = 48, 52
//	c = 49 (fixed)
//	d = 0-255
//
// Default:
//
//	None
//
// Parameters:
//
//	keyCode1: First byte of key code (ASCII printable character)
//	keyCode2: Second byte of key code (ASCII printable character)
//	tone: Graphics tone:
//	      48 -> Monochrome (digital)
//	      52 -> Multiple tone
//	bmpData: Windows BMP format image data
//
// Notes:
//   - Converts Windows BMP data to the specified tone and defines it as NV graphics data
//     (raster format) associated with the specified key codes
//   - For optimal performance, keep the number of NV graphics registrations below 50:
//   - ≤50 items: max execution time 60 seconds
//   - 100 items: execution time up to 120 seconds
//   - The NV graphics memory usage is [BMP data size + 24 bytes of control information]
//   - NV graphics and NV bit image (FS q) cannot be defined simultaneously
//   - When this function is executed, all NV bit images are deleted
//   - If the data size exceeds the remaining NV graphics memory capacity, this function is ignored
//   - Input data must be in standard Windows BMP file format
//
// Errors:
//
//	Returns ErrInvalidKeyCode if keyCode1 or keyCode2 is outside the valid range (32-126).
//	Returns ErrInvalidTone if tone is not a valid value (48, 52).
//	Returns ErrInvalidBMPFormat if the BMP data is not in a valid format.
func (c *NvGraphicsCommands) DefineWindowsBMPNVGraphics(kc1, kc2 byte, tone GraphicsTone, bmpData []byte) ([]byte, error) {
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

	cmd := []byte{common.GS, 'D', 0x30, 0x43, 0x30, kc1, kc2, byte(tone), 0x31}
	cmd = append(cmd, bmpData...)

	return cmd, nil
}
