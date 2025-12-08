package mechanismcontrol

import (
	"github.com/adcondev/poster/pkg/commands/common"
)

// ReturnHome moves the print head to the standby position.
//
// Format:
//
//	ASCII:   ESC <
//	Hex:     0x1B 0x3C
//	Decimal: 27 60
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
//   - The standby position is different, depending on the printer model
//   - The command rechecks the standby position; therefore, the print position might be shifted before and after checking the standby position
//
// Errors:
//
//	This function is safe and does not return errors.
func (c *Commands) ReturnHome() []byte {
	return []byte{common.ESC, '<'}
}

// SetUnidirectionalPrintMode turns unidirectional print mode on or off.
//
// Format:
//
//	ASCII:   ESC U n
//	Hex:     0x1B 0x55 n
//	Decimal: 27 85 n
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
//	n: Controls unidirectional print mode (LSB is used):
//	   LSB = 0 -> Unidirectional print mode OFF (bidirectional mode ON)
//	   LSB = 1 -> Unidirectional print mode ON
//
// Notes:
//   - Mode can be set independently in Standard mode and in Page mode
//   - When used in Standard mode, the printer sets the mode for Standard mode
//   - When used in Page mode, the printer sets the mode for Page mode
//   - When unidirectional print mode is turned off, bidirectional print mode is automatically turned on
//   - In Page mode, the printer performs unidirectional printing for all data collectively printed using FF or ESC FF
//   - Unidirectional print mode is recommended when printing double-height characters, graphics, bit images, or two-dimensional codes to ensure that the top and bottom of the printing patterns are aligned
//   - Settings remain in effect until ESC @ is executed, the printer is reset, or power is turned off
//
// Errors:
//
//	This function is safe and does not return errors.
func (c *Commands) SetUnidirectionalPrintMode(n UnidirectionalMode) []byte {
	return []byte{common.ESC, 'U', byte(n)}
}

// PartialCut executes a partial cut of the roll paper (one point left uncut).
//
// Deprecated: This command is obsolete and will not be supported in future printer models.
// Use GS V or GS ( V instead for cutting paper. CutPaper GS V <Function A> gives the same result.
//
// Format:
//
//	ASCII:   ESC i
//	Hex:     0x1B 0x69
//	Decimal: 27 105
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
//   - The cutting shape depends on the specification of the mounted autocutter
//   - See GS V <Function A> for details
//
// Errors:
//
//	This function is safe and does not return errors.
func (c *Commands) PartialCut() []byte {
	return []byte{common.ESC, 'i'}
}

// PartialCutThreePoints executes a partial cut of the roll paper (three points left uncut).
//
// Deprecated: This command is obsolete and will not be supported in future printer models.
// Use GS V or GS ( V instead for cutting paper. CutPaper GS V <Function A> gives the same result.
//
// Format:
//
//	ASCII:   ESC m
//	Hex:     0x1B 0x6D
//	Decimal: 27 109
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
//   - The cutting shape depends on the specification of the mounted autocutter
//   - See GS V <Function A> for details
//
// Errors:
//
//	This function is safe and does not return errors.
func (c *Commands) PartialCutThreePoints() []byte {
	return []byte{common.ESC, 'm'}
}

// PaperCut executes a paper cut operation.
//
// Format:
//
//	ASCII:   GS ( V pL pH fn m
//	Hex:     0x1D 0x28 0x56 pL pH 0x30 m
//	Decimal: 29 40 86 pL pH 48 m
//
// Range:
//
//	(pL + pH × 256) = 2
//	fn = 48
//	m = 0, 1
//
// Default:
//
//	None
//
// Parameters:
//
//	m: Cut type:
//	   0 -> CutTypeFull cut
//	   1 -> CutTypePartial cut
//
// Notes:
//   - This command is not accompanied by a paper feed operation
//   - The pL, pH values are fixed at 2, 0 respectively
//   - The fn value is fixed at 48
//
// Errors:
//
//	Returns ErrCutType if m is not a valid cut type value (0, 1).
func (c *Commands) PaperCut(m CutType) ([]byte, error) {
	if err := ValidateCutType(m); err != nil {
		return nil, err
	}
	return []byte{common.GS, '(', 'V', 0x02, 0x00, 0x30, byte(m)}, nil
}

// PaperFeedAndCut executes paper cut after feeding paper by the specified amount.
//
// Format:
//
//	ASCII:   GS ( V pL pH fn m n
//	Hex:     0x1D 0x28 0x56 pL pH 0x31 m n
//	Decimal: 29 40 86 pL pH 49 m n
//
// Range:
//
//	(pL + pH × 256) = 3
//	fn = 49
//	m = 0, 1
//	n = 0–255
//
// Default:
//
//	None
//
// Parameters:
//
//	m: Cut type:
//	   0 -> CutTypeFull cut after paper feed
//	   1 -> CutTypePartial cut after paper feed
//	n: Paper feed amount before cutting: (cut position + [n × horizontal or vertical motion unit])
//
// Notes:
//   - Executes paper feeding to (cut position + [n × horizontal or vertical motion unit]) then executes the specified cut type
//   - Printers without an autocutter only feed the paper for the specified amount
//   - The pL, pH values are fixed at 3, 0 respectively
//   - The fn value is fixed at 49
//
// Errors:
//
//	Returns ErrCutType if m is not a valid cut type value (0, 1).
func (c *Commands) PaperFeedAndCut(m CutType, n byte) ([]byte, error) {
	if err := ValidateCutType(m); err != nil {
		return nil, err
	}
	return []byte{common.GS, '(', 'V', 0x03, 0x00, 0x31, byte(m), n}, nil
}

// ReservePaperCut reserves a paper cut operation to be executed during continuous printing.
//
// Format:
//
//	ASCII:   GS ( V pL pH fn m n
//	Hex:     0x1D 0x28 0x56 pL pH 0x33 m n
//	Decimal: 29 40 86 pL pH 51 m n
//
// Range:
//
//	(pL + pH × 256) = 3
//	fn = 51
//	m = 0, 1
//	n = 0–255
//
// Default:
//
//	None
//
// Parameters:
//
//	m: Cut type:
//	   0 -> CutTypeFull cut at specified position
//	   1 -> CutTypePartial cut at specified position
//	n: Paper feed amount: (cut position + [n × horizontal or vertical motion unit])
//
// Notes:
//   - The paper cut position reservation function is available only during continuous printing operation
//   - If data is interrupted for two seconds or more, the paper is automatically fed to the reserved cut position and then cut
//   - The preset cutting position is cleared if any of the following occurs:
//   - Execution of commands that perform a software reset
//   - Power off or service reset
//   - Paper feed by pressing paper feed button
//   - Execution of buffer clear commands
//   - The paper cut involves a stop of printing. This might affect the print quality
//   - The pL, pH values are fixed at 3, 0 respectively
//   - The fn value is fixed at 51
//
// Errors:
//
//	Returns ErrCutType if m is not a valid cut type value (0, 1).
func (c *Commands) ReservePaperCut(m CutType, n byte) ([]byte, error) {
	if err := ValidateCutType(m); err != nil {
		return nil, err
	}
	return []byte{common.GS, '(', 'V', 0x03, 0x00, 0x33, byte(m), n}, nil
}

// CutPaper executes paper cutting with the specified mode (Function A).
//
// Format:
//
//	ASCII:   GS V m
//	Hex:     0x1D 0x56 m
//	Decimal: 29 86 m
//
// Range:
//
//	m = 0, 1, 48, 49
//
// Default:
//
//	None
//
// Parameters:
//
//	m: Cutting shape:
//	   0 or 48 -> CutTypeFull cut
//	   1 or 49 -> CutTypePartial cut
//
// Notes:
//   - When Standard mode is selected, this command is enabled only when processed at the beginning of line
//   - This command can be used with an autocutter
//   - When using this command, note that there is a gap between the cutting position and the print position
//
// Errors:
//
//	Returns ErrCutMode if m is not a valid cutting mode value (0, 1, 48, 49).
func (c *Commands) CutPaper(m CutMode) ([]byte, error) {
	if err := ValidateCutMode(m); err != nil {
		return nil, err
	}
	return []byte{common.GS, 'V', byte(m)}, nil
}

// FeedAndCutPaper feeds paper to the specified position and executes paper cutting (Function B).
//
// Format:
//
//	ASCII:   GS V m n
//	Hex:     0x1D 0x56 m n
//	Decimal: 29 86 m n
//
// Range:
//
//	m = 65, 66
//	n = 0–255
//
// Default:
//
//	None
//
// Parameters:
//
//	m: Cutting shape:
//	   65 -> CutTypeFull cut
//	   66 -> CutTypePartial cut
//	n: FeedLines amount before cutting: [cutting position + (n × vertical motion unit)]
//
// Notes:
//   - When Standard mode is selected, this command is enabled only when processed at the beginning of line
//   - This command can be used with an autocutter
//   - Printers without an autocutter only feed the paper for specified amount
//
// Errors:
//
//	Returns ErrFeedCutMode if m is not a valid value (65, 66).
func (c *Commands) FeedAndCutPaper(m FeedCut, n byte) ([]byte, error) {
	if err := ValidateFeedCutMode(m); err != nil {
		return nil, err
	}
	return []byte{common.GS, 'V', byte(m), n}, nil
}

// SetCutPosition sets the paper cutting position and executes paper cut when it reaches the autocutter position (Function C).
//
// Format:
//
//	ASCII:   GS V m n
//	Hex:     0x1D 0x56 m n
//	Decimal: 29 86 m n
//
// Range:
//
//	m = 97, 98
//	n = 0–255
//
// Default:
//
//	None
//
// Parameters:
//
//	m: Cutting shape:
//	   97 -> CutTypeFull cut
//	   98 -> CutTypePartial cut
//	n: Position offset: [cutting position + (n × vertical motion unit)]
//
// Notes:
//   - When Standard mode is selected, this command is enabled only when processed at the beginning of line
//   - This command can be used with an autocutter
//   - The preset cutting position is cleared if any of the following occurs:
//   - Execution of ESC @
//   - Execution of commands that perform a software reset
//   - Power off or service reset
//   - Paper feed by pressing paper feed button
//   - Execution of buffer clear commands
//   - The paper cut involves a stop of printing. This might affect the print quality
//   - This function makes it possible to save paper by reducing the top margin for continuous printing
//
// Errors:
//
//	Returns ErrPositionCutMode if m is not a valid value (97, 98).
func (c *Commands) SetCutPosition(m PositionCut, n byte) ([]byte, error) {
	if err := ValidatePositionCutMode(m); err != nil {
		return nil, err
	}
	return []byte{common.GS, 'V', byte(m), n}, nil
}

// FeedCutAndReturnPaper feeds paper to the specified position, executes paper cutting, and feeds paper to the print starting position (Function D).
//
// Format:
//
//	ASCII:   GS V m n
//	Hex:     0x1D 0x56 m n
//	Decimal: 29 86 m n
//
// Range:
//
//	m = 103, 104
//	n = 0–255
//
// Default:
//
//	None
//
// Parameters:
//
//	m: Cutting shape:
//	   103 -> CutTypeFull cut
//	   104 -> CutTypePartial cut
//	n: FeedLines amount before cutting: [cutting position + (n × vertical motion unit)]
//
// Notes:
//   - When Standard mode is selected, this command is enabled only when processed at the beginning of line
//   - This command can be used with an autocutter
//   - This function makes it possible to save paper by reducing the top margin for continuous printing
//   - The top margin of the print start position is different depending on the printer model
//   - After cutting, the paper is moved to the print start position by reverse feeding
//
// Errors:
//
//	Returns ErrFeedCutReturnMode if m is not a valid value (103, 104).
func (c *Commands) FeedCutAndReturnPaper(m FeedCutReturn, n byte) ([]byte, error) {
	if err := ValidateFeedCutReturnMode(m); err != nil {
		return nil, err
	}
	return []byte{common.GS, 'V', byte(m), n}, nil
}
