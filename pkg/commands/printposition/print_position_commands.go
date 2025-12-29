package printposition

import (
	"fmt"

	"github.com/adcondev/poster/pkg/commands/shared"
)

// HorizontalTab moves the print position to the next horizontal tab position.
//
// Format:
//
//	ASCII:   HT
//	Hex:     0x09
//	Decimal: 9
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
//   - Moves the print position to the next horizontal tab position
//   - Ignored unless the next horizontal tab position has been set (ESC D)
//   - If the next tab position exceeds the print area, print position is set to [Print area width + 1]
//   - If processed when at [Print area width + 1], the printer executes print-buffer-full for the current line
//     and performs horizontal tab processing from the beginning of the next line.
//     In Page mode, printing is not executed but the print position is moved
//   - The printer will not move to the beginning of the line by executing this command
//   - When underline mode is on, the underline is not printed under the tab space skipped by this command
//
// Errors:
//
//	This function is safe and does not return errors.
func (c *Commands) HorizontalTab() []byte {
	return []byte{shared.HT}
}

// SetAbsolutePrintPosition sets the absolute print position.
//
// Format:
//
//	ASCII:   ESC $ nL nH
//	Hex:     0x1B 0x24 nL nH
//	Decimal: 27 36 nL nH
//
// Range:
//
//	(nL + nH × 256) = 0 – 65535
//
// Default:
//
//	None
//
// Parameters:
//
//	position: Absolute position as (nL + nH × 256) × (horizontal or vertical motion unit)
//	          from the left edge of the print area
//
// Notes:
//   - The printer ignores any setting that exceeds the print area
//   - In Standard mode the horizontal motion unit is used
//   - In Page mode the horizontal or vertical motion unit is used depending on
//     the print direction set by ESC T
//   - If the starting position is set to upper-left or lower-right using ESC T,
//     the horizontal motion unit is used; for upper-right or lower-left the vertical
//     motion unit is used
//   - If the motion unit changes after this command, the print position is not changed
//   - The printer will not move to the beginning of the line by executing this command
//   - When underline mode is on, the underline is not printed under the space skipped
//     by this command
//
// Errors:
//
//	This function is safe and does not return errors.
func (c *Commands) SetAbsolutePrintPosition(position uint16) []byte {
	nL, nH := shared.ToLittleEndian(position)
	return []byte{shared.ESC, '$', nL, nH}
}

// SetHorizontalTabPositions sets horizontal tab positions.
//
// Format:
//
//	ASCII:   ESC D n1 ... nk NUL
//	Hex:     0x1B 0x44 n1 ... nk 0x00
//	Decimal: 27 68 n1 ... nk 0
//
// Range:
//
//	n = 1–255
//	k = 0–32
//
// Default:
//
//	n = 8, 16, 24, 32, ..., 232, 240, 248 (every eight characters for the default font set)
//
// Parameters:
//
//	positions: Array of tab stop positions in ascending order (max 32 positions)
//	           Each value specifies the number of character widths from line start
//
// Notes:
//   - Sets horizontal tab positions. Each transmitted n value specifies the number
//     of character widths from the line start to the tab stop. Transmit the tab
//     stops in ascending order and terminate the list with NUL (0x00). Transmitting
//     ESC D NUL clears all horizontal tab positions
//   - The tab position is stored as [character width × n], where character width
//     includes right-side character spacing. Double-width characters count as
//     twice the width
//   - Character width and font/spacing/enlargement should be set before sending this command
//   - A maximum of 32 tab positions can be set; data beyond 32 is treated as normal data
//   - If a transmitted n is less than or equal to the previous value, tab-setting
//     is finished and subsequent bytes are processed as normal data
//   - Tab settings are preserved until ESC @ (initialize), printer reset, or power-off
//   - Changing the left margin will shift stored tab positions accordingly
//   - Horizontal tab positions that exceed the print area are allowed; they become
//     effective or not depending on the current print area width
//
// Errors:
//
//	Returns ErrTooManyTabPositions if more than 32 positions are provided.
//	Returns ErrTabPosition if positions are not in ascending order or contain invalid values.
func (c *Commands) SetHorizontalTabPositions(positions []byte) ([]byte, error) {
	// Check maximum number of positions
	if len(positions) > MaxTabPositions {
		return nil, ErrTooManyTabPositions
	}

	// Validate positions are in ascending order and within range
	prevPos := byte(0)
	for i, pos := range positions {
		if pos == 0 || pos > MaxTabValue {
			return nil, fmt.Errorf("%w: position %d at index %d", ErrTabPosition, pos, i)
		}
		if pos <= prevPos {
			return nil, fmt.Errorf("%w: position %d at index %d must be greater than %d", ErrTabPosition, pos, i, prevPos)
		}
		prevPos = pos
	}

	// Build command
	cmd := []byte{shared.ESC, 'D'}
	cmd = append(cmd, positions...)
	cmd = append(cmd, shared.NUL)
	return cmd, nil
}

// SelectPrintDirectionPageMode selects the print direction and starting position in Page mode.
//
// Format:
//
//	ASCII:   ESC T n
//	Hex:     0x1B 0x54 n
//	Decimal: 27 84 n
//
// Range:
//
//	n = 0–3, 48–51
//
// Default:
//
//	n = 0 (Left to right, starting position: upper left)
//
// Parameters:
//
//	direction: Selects print direction and starting position:
//	           0 or 48 -> Print direction: left to right;  Starting position: upper left
//	           1 or 49 -> Print direction: bottom to top;  Starting position: lower left
//	           2 or 50 -> Print direction: right to left;  Starting position: lower right
//	           3 or 51 -> Print direction: top to bottom;  Starting position: upper right
//
// Notes:
//   - In Page mode, selects print direction and starting position
//   - Effective only in Page mode; has no effect in Standard mode
//   - The meaning of horizontal/vertical motion units for other commands depends on the selected starting position (see command reference)
//   - Settings persist until ESC @ (initialize), printer reset, or power-off
//
// Errors:
//
//	Returns ErrPrintDirection if direction is not a valid value (0-3, 48-51).
func (c *Commands) SelectPrintDirectionPageMode(direction PrintDirection) ([]byte, error) {
	if err := ValidatePrintDirection(direction); err != nil {
		return nil, err
	}
	return []byte{shared.ESC, 'T', byte(direction)}, nil
}

// SetPrintAreaPageMode sets the print area and logical origin in Page mode.
//
// Format:
//
//	ASCII:   ESC W xL xH yL yH dxL dxH dyL dyH
//	Hex:     0x1B 0x57 xL xH yL yH dxL dxH dyL dyH
//	Decimal: 27 87 xL xH yL yH dxL dxH dyL dyH
//
// Range:
//
//	x, y = 0–65535
//	width, height = 1–65535
//
// Default:
//
//	Depends on printer model
//
// Parameters:
//
//	x: Horizontal logical origin = (xL + xH*256) × (horizontal motion unit)
//	y: Vertical logical origin = (yL + yH*256) × (vertical motion unit)
//	width: Print area width = (dxL + dxH*256) × (horizontal motion unit)
//	height: Print area height = (dyL + dyH*256) × (vertical motion unit)
//
// Notes:
//   - In Page mode, defines the logical origin (horizontal and vertical) and the
//     print area size. The transmitted parameters are interpreted as 16-bit
//     little-endian values: value = (low + high * 256), each measured in the
//     horizontal or vertical motion unit as appropriate
//   - This command only has effect in Page mode (ESC L) and is ignored in Standard mode
//   - Both print area width and height must be at least 1 (cannot be zero)
//   - Logical origins must lie within the printable area
//   - If origin + size exceeds the printable area, the size is reduced to fit
//   - Values are fixed even if motion units change later
//   - Settings persist until FF (in Page mode), ESC @ (initialize), reset, or power-off
//   - The absolute origin is the upper-left of the printable area
//   - For printers supporting GS ( P <Function 48>, the maximum and origin align with that printable area setting
//
// Errors:
//
//	Returns ErrPrintAreaWidthSize if width is 0.
//	Returns ErrPrintAreaHeightSize if height is 0.
func (c *Commands) SetPrintAreaPageMode(x, y, width, height uint16) ([]byte, error) {
	if err := ValidatePrintArea(width, height); err != nil {
		return nil, err
	}
	xL, xH := shared.ToLittleEndian(x)
	yL, yH := shared.ToLittleEndian(y)
	dxL, dxH := shared.ToLittleEndian(width)
	dyL, dyH := shared.ToLittleEndian(height)
	return []byte{shared.ESC, 'W', xL, xH, yL, yH, dxL, dxH, dyL, dyH}, nil
}

// SetRelativePrintPosition moves the print position relative to the current position.
//
// Format:
//
//	ASCII:   ESC \ nL nH
//	Hex:     0x1B 0x5C nL nH
//	Decimal: 27 92 nL nH
//
// Range:
//
//	(nL + nH × 256) = -32768 – 32767 (signed 16-bit value)
//
// Default:
//
//	None
//
// Parameters:
//
//	distance: Relative distance as signed 16-bit value in motion units.
//	          Positive moves to the right; negative moves to the left.
//
// Notes:
//   - Moves the print position by (nL + nH × 256) × (horizontal or vertical motion unit)
//     from the current position
//   - The printer ignores any setting that exceeds the print area
//   - In Standard mode the horizontal motion unit is used
//   - In Page mode the horizontal or vertical motion unit is used depending on the print direction set by ESC T
//   - If the starting position is upper-left or lower-right (ESC T), the horizontal motion unit is used.
//     If the starting position is upper-right or lower-left, the vertical motion unit is used
//   - Changing motion units after this command does not change the already-set print position
//   - Underline mode does not print under the space skipped by this command
//   - In JIS code, '\' corresponds to '¥'
//
// Errors:
//
//	This function is safe and does not return errors.
func (c *Commands) SetRelativePrintPosition(distance int16) []byte {
	// Convert signed int16 to unsigned bytes (little-endian)
	// intentional: preserve int16 two's-complement bit pattern for ESC \ command
	value := uint16(distance) //nolint:gosec
	nL, nH := shared.ToLittleEndian(value)
	return []byte{shared.ESC, '\\', nL, nH}
}

// SelectJustification selects text justification in Standard mode.
//
// Format:
//
//	ASCII:   ESC a n
//	Hex:     0x1B 0x61 n
//	Decimal: 27 97 n
//
// Range:
//
//	n = 0–2, 48–50
//
// Default:
//
//	n = 0 (Left)
//
// Parameters:
//
//	mode: Justification mode:
//	      0 or 48 -> Left justification
//	      1 or 49 -> CenterAlign
//	      2 or 50 -> Right justification
//
// Notes:
//   - In Standard mode, aligns all data in one line according to n
//   - Effective only in Standard mode and only when processed at the beginning of a line
//   - Has no effect in Page mode
//   - Justification is applied within the print area set by GS L and ESC W/GS W
//   - Affects characters, graphics, barcodes, 2D codes and space areas set by HT, ESC $, ESC \
//   - Setting persists until ESC @, reset, or power-off
//
// Errors:
//
//	Returns ErrJustification if mode is not a valid value (0-2, 48-50).
func (c *Commands) SelectJustification(mode Justification) ([]byte, error) {
	if err := ValidateJustification(mode); err != nil {
		return nil, err
	}
	return []byte{shared.ESC, 'a', byte(mode)}, nil
}

// SetAbsoluteVerticalPrintPosition sets the absolute vertical print position in Page mode.
//
// Format:
//
//	ASCII:   GS $ nL nH
//	Hex:     0x1D 0x24 nL nH
//	Decimal: 29 36 nL nH
//
// Range:
//
//	(nL + nH × 256) = 0 – 65535
//
// Default:
//
//	None
//
// Parameters:
//
//	position: Absolute vertical position from the starting position set by ESC T
//
// Notes:
//   - In Page mode, moves the vertical print position to (nL + nH × 256) × (vertical or horizontal motion unit)
//     from the starting position set by ESC T
//   - This command is enabled only in Page mode; it is ignored in Standard mode
//   - The printer ignores any setting that exceeds the print area set by ESC W
//   - The horizontal or vertical motion unit used depends on the print direction set by ESC T
//   - If the starting position is upper left or lower right, the vertical motion unit is used.
//     If the starting position is upper right or lower left, the horizontal motion unit is used
//   - Changing motion units after this command does not change the already-set print position
//
// Errors:
//
//	This function is safe and does not return errors.
func (c *Commands) SetAbsoluteVerticalPrintPosition(position uint16) []byte {
	nL := byte(position & 0xFF)
	nH := byte((position >> 8) & 0xFF)
	return []byte{shared.GS, '$', nL, nH}
}

// SetLeftMargin sets the left margin in Standard mode.
//
// Format:
//
//	ASCII:   GS L nL nH
//	Hex:     0x1D 0x4C nL nH
//	Decimal: 29 76 nL nH
//
// Range:
//
//	(nL + nH × 256) = 0 – 65535
//
// Default:
//
//	(nL + nH × 256) = 0
//
// Parameters:
//
//	margin: Left margin as (nL + nH × 256) × (horizontal motion unit)
//	        from the left edge of the printable area
//
// Notes:
//   - In Standard mode, sets the left margin from the left edge of the printable area
//   - Effective in Standard mode only when processed at the beginning of the line
//   - Has no effect while in Page mode; if issued in Page mode the value is stored
//     and enabled when returning to Standard mode
//   - If the setting exceeds the printable area, it is clamped to the printable-area maximum
//   - If this command and GS W would set the print area width to less than one character,
//     the print area width is extended to accommodate one character
//   - Uses the horizontal motion unit; changes to the motion unit after setting do not change the margin
//   - Setting persists until ESC @ (initialize), printer reset, or power-off
//   - The left margin is measured from the left edge of the printable area; changing the left margin moves that edge
//
// Errors:
//
//	This function is safe and does not return errors.
func (c *Commands) SetLeftMargin(margin uint16) []byte {
	nL := byte(margin & 0xFF)
	nH := byte((margin >> 8) & 0xFF)
	return []byte{shared.GS, 'L', nL, nH}
}

// SetPrintPositionBeginningLine moves the print position to the beginning of the print line.
//
// Format:
//
//	ASCII:   GS T n
//	Hex:     0x1D 0x54 n
//	Decimal: 29 84 n
//
// Range:
//
//	n = 0, 1, 48, 49
//
// Default:
//
//	None
//
// Parameters:
//
//	mode: Controls how the print buffer is processed:
//	      0 or 48 -> Erase the data in the print buffer, then move the print position
//	      1 or 49 -> Print the data in the print buffer, then move the print position (starts a new line based on line spacing)
//
// Notes:
//   - In Standard mode, moves the print position to the beginning (left side) of the printable area
//     after performing the operation specified by n
//   - Effective only in Standard mode; ignored in Page mode
//   - Ignored if the print position is already at the beginning of the line
//   - If print position is not at the beginning of the line and n = 1 or 49, this behaves the same as LF
//   - Erase (n = 0 or 48) cancels the current print-buffered data but preserves other settings and buffer contents
//   - After execution the printer is in the "Beginning of the line" status
//   - Use this command immediately before other commands that require beginning-of-line to ensure they execute
//
// Errors:
//
//	Returns ErrBeginLineMode if mode is not a valid value (0, 1, 48, 49).
func (c *Commands) SetPrintPositionBeginningLine(mode BeginLine) ([]byte, error) {
	if err := ValidateBeginLineMode(mode); err != nil {
		return nil, err
	}
	return []byte{shared.GS, 'T', byte(mode)}, nil
}

// SetPrintAreaWidth sets the print area width in Standard mode.
//
// Format:
//
//	ASCII:   GS W nL nH
//	Hex:     0x1D 0x57 nL nH
//	Decimal: 29 87 nL nH
//
// Range:
//
//	(nL + nH × 256) = 0 – 65535
//
// Default:
//
//	Entire printable area (model-dependent; example: 576 for 80mm, 420 for 58mm paper)
//
// Parameters:
//
//	width: Print area width as (nL + nH × 256) × (horizontal motion unit)
//
// Notes:
//   - In Standard mode, sets the print area width
//   - This command is effective in Standard mode only when processed at the beginning of a line
//   - Has no effect in Page mode; if issued in Page mode the value is stored and enabled when returning to Standard mode
//   - If [left margin + print area width] exceeds the printable area, the print area width is clamped to [printable area - left margin]
//   - If this command together with GS L would set the print area width to less than one character, the width is extended to accommodate one character
//   - Uses the horizontal motion unit. Changing the motion unit after setting does not change the stored width
//   - Setting persists until ESC @, printer reset, or power-off
//
// Errors:
//
//	This function is safe and does not return errors.
func (c *Commands) SetPrintAreaWidth(width uint16) []byte {
	nL := byte(width & 0xFF)
	nH := byte((width >> 8) & 0xFF)
	return []byte{shared.GS, 'W', nL, nH}
}

// SetRelativeVerticalPrintPosition moves the vertical print position relative to the current position in Page mode.
//
// Format:
//
//	ASCII:   GS \ nL nH
//	Hex:     0x1D 0x5C nL nH
//	Decimal: 29 92 nL nH
//
// Range:
//
//	(nL + nH × 256) = -32768 – 32767 (signed 16-bit value)
//
// Default:
//
//	None
//
// Parameters:
//
//	distance: Relative vertical distance as signed 16-bit value in motion units.
//	          Positive moves downward; negative moves upward.
//
// Notes:
//   - In Page mode, moves the vertical print position by (nL + nH × 256) × (vertical or horizontal motion unit)
//     from the current position
//   - This command is enabled only in Page mode; it is ignored in Standard mode
//   - The printer ignores any setting that exceeds the print area set by ESC W
//   - The horizontal or vertical motion unit used depends on the print direction set by ESC T
//   - If starting position is upper-left or lower-right (ESC T), the vertical motion unit is used
//   - If starting position is upper-right or lower-left (ESC T), the horizontal motion unit is used
//   - Changes to the motion units after executing this command do not affect the already-set print position
//   - In JIS code, '\' corresponds to '¥'
//
// Errors:
//
//	This function is safe and does not return errors.
func (c *Commands) SetRelativeVerticalPrintPosition(distance int16) []byte {
	// Convert signed int16 to unsigned bytes (little-endian)
	// intentional: preserve int16 two's-complement bit pattern for ESC \ command
	value := uint16(distance) //nolint:gosec
	nL := byte(value & 0xFF)
	nH := byte((value >> 8) & 0xFF)
	return []byte{shared.GS, '\\', nL, nH}
}
