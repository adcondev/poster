package print

import (
	"github.com/adcondev/poster/pkg/commands/shared"
)

// FormFeed executes form feed operation (behavior varies by mode).
//
// Format:
//
//	ASCII:   FF
//	Hex:     0x0C
//	Decimal: 12
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
//   - In Standard mode: Indicates "Printing is completed" for the current job,
//     feeds to cutting position if reserved by GS V commands
//   - In Page mode: Prints all data in the print buffer collectively,
//     switches from Page mode to Standard mode, clears the print area after printing,
//     resets ESC W values to defaults
//   - Mode-specific behavior; check current mode before use
//   - ESC T value is maintained when returning from Page mode
//
// Errors:
//
//	This function is safe and does not return errors.
func (c *Commands) FormFeed() []byte {
	return []byte{FF}
}

// PrintAndCarriageReturn prints the data in the print buffer and performs a carriage return.
//
// Format:
//
//	ASCII:   CR
//	Hex:     0x0D
//	Decimal: 13
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
//   - Executes one of the following operations depending on the print head type
//     and auto line feed setting:
//     When auto line feed is enabled: executes printing and one line feed (same as LF)
//     When auto line feed is disabled: behavior depends on print head type
//   - In Standard mode: prints data and moves to beginning of line
//   - In Page mode: only moves print position without printing
//   - With a serial interface, the command performs as if auto line feed is disabled
//   - Auto line feed can be configured by DIP switch or memory switch (GS ( E <Function 3>)
//   - After printing, the print position is moved to the left side of the printable area
//
// Errors:
//
//	This function is safe and does not return errors.
func (c *Commands) PrintAndCarriageReturn() []byte {
	return []byte{CR}
}

// PrintAndLineFeed prints the data in the print buffer and feeds one line.
//
// Format:
//
//	ASCII:   LF
//	Hex:     0x0A
//	Decimal: 10
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
//   - Prints the data in the print buffer and feeds one line, based on the current line spacing
//   - The amount of paper fed per line is based on the value set using the
//     line spacing command (ESC 2 or ESC 3)
//   - After printing, the print position is moved to the left side of the
//     printable area and the printer enters the "Beginning of the line" status
//   - When this command is processed in Page mode, only the print position
//     moves and the printer does not perform actual printing
//
// Errors:
//
//	This function is safe and does not return errors.
func (c *Commands) PrintAndLineFeed() []byte {
	return []byte{LF}
}

// PrintAndFeedPaper prints the data in the print buffer and feeds the paper.
//
// Format:
//
//	ASCII:   ESC J n
//	Hex:     0x1B 0x4A n
//	Decimal: 27 74 n
//
// Range:
//
//	n = 0–255
//
// Default:
//
//	None
//
// Parameters:
//
//	n: Feeds the paper [n × (vertical or horizontal motion unit)]
//
// Notes:
//   - Prints the data in the print buffer and feeds the paper
//   - In Standard mode the vertical motion unit is used
//   - In Page mode the vertical or horizontal motion unit is used according
//     to the print direction set by ESC T
//   - Maximum paper feed amount depends on the printer model
//   - After printing, the print position moves to the beginning of the line
//   - In Page mode, only the print position moves without actual printing
//
// Errors:
//
//	This function is safe and does not return errors.
func (c *Commands) PrintAndFeedPaper(n byte) []byte {
	return []byte{shared.ESC, 'J', n}
}

// PrintAndFeedLines prints the data in the print buffer and feeds n lines.
//
// Format:
//
//	ASCII:   ESC d n
//	Hex:     0x1B 0x64 n
//	Decimal: 27 100 n
//
// Range:
//
//	n = 0–255
//
// Default:
//
//	None
//
// Parameters:
//
//	n: Number of lines to feed
//
// Notes:
//   - Prints the data in the print buffer and feeds n lines
//   - Paper feed per line based on line spacing (ESC 2 or ESC 3)
//   - Maximum feed depends on printer model
//   - After printing, print position moves to beginning of line
//   - In Page mode, only print position moves without actual printing
//
// Errors:
//
//	This function is safe and does not return errors.
func (c *Commands) PrintAndFeedLines(n byte) []byte {
	return []byte{shared.ESC, 'd', n}
}

// PrintAndReverseFeed prints the data in the print buffer and feeds paper in reverse.
//
// Format:
//
//	ASCII:   ESC K n
//	Hex:     0x1B 0x4B n
//	Decimal: 27 75 n
//
// Range:
//
//	n = 0–48
//
// Default:
//
//	None
//
// Parameters:
//
//	n: Feeds the paper n × (vertical or horizontal motion unit) in reverse direction
//
// Notes:
//   - Prints the data in the print buffer and feeds the paper in reverse
//   - Motion unit used depends on mode and ESC T setting
//   - If n exceeds model maximum, reverse feed is not executed but printing occurs
//   - Some printers perform small forward feed after reverse feed
//   - In Page mode, only print position moves without actual printing
//
// Errors:
//
//	Returns ErrReverseUnits if n exceeds MaxReverseMotionUnits (48).
func (c *Commands) PrintAndReverseFeed(n byte) ([]byte, error) {
	if err := ValidateReverseFeedUnits(n); err != nil {
		return nil, err
	}
	return []byte{shared.ESC, 'K', n}, nil
}

// PrintAndReverseFeedLines prints the data in the print buffer and feeds n lines in reverse.
//
// Format:
//
//	ASCII:   ESC e n
//	Hex:     0x1B 0x65 n
//	Decimal: 27 101 n
//
// Range:
//
//	n = 0–2
//
// Default:
//
//	None
//
// Parameters:
//
//	n: Number of lines to feed in reverse direction
//
// Notes:
//   - Prints the data in the print buffer and feeds n lines in reverse
//   - Paper feed per line based on line spacing (ESC 2 or ESC 3)
//   - If n exceeds model maximum, reverse feed is not executed but printing occurs
//   - Some printers perform small forward feed after reverse feed
//   - In Page mode, only print position moves without actual printing
//
// Errors:
//
//	Returns ErrReverseLines if n exceeds MaxReverseFeedLines (2).
func (c *Commands) PrintAndReverseFeedLines(n byte) ([]byte, error) {
	if err := ValidateReverseFeedLines(n); err != nil {
		return nil, err
	}
	return []byte{shared.ESC, 'e', n}, nil
}

// PrintDataInPageMode prints the data in the print buffer collectively in Page mode.
//
// Format:
//
//	ASCII:   ESC FF
//	Hex:     0x1B 0x0C
//	Decimal: 27 12
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
//   - In Page mode, prints all data currently buffered in the print area collectively
//   - Enabled only in Page mode (selected by ESC L)
//   - After printing, buffer data, print position, and settings remain intact
//   - Commonly used for reprinting the same Page-mode data multiple times
//   - Returns to Standard mode when FF, ESC S, or ESC @ is issued
//
// Errors:
//
//	This function is safe and does not return errors.
func (c *Commands) PrintDataInPageMode() []byte {
	return []byte{shared.ESC, FF}
}

// CancelData deletes all print data in the current print area (Page mode only).
//
// Format:
//
//	ASCII:   CAN
//	Hex:     0x18
//	Decimal: 24
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
//   - In Page mode, deletes all the print data in the current print area
//   - Enabled only in Page mode (selected by ESC L)
//   - Also deletes overlapping data from previously specified print areas
//   - Has no effect in Standard mode
//
// Errors:
//
//	This function is safe and does not return errors.
func (c *Commands) CancelData() []byte {
	return []byte{CAN}
}
