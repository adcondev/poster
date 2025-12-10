package test_test

import (
	"bytes"
	"testing"

	"github.com/adcondev/poster/pkg/commands/common"
	"github.com/adcondev/poster/pkg/commands/printposition"
)

func TestIntegration_PrintPosition_StandardModeWorkflow(t *testing.T) {
	cmd := printposition.NewCommands()

	t.Run("complete standard mode positioning", func(t *testing.T) {
		var buffer []byte

		// Set up print area
		buffer = append(buffer, cmd.SetLeftMargin(50)...)
		buffer = append(buffer, cmd.SetPrintAreaWidth(400)...)

		// Set justification
		justCmd, err := cmd.SelectJustification(printposition.Center)
		if err != nil {
			t.Fatalf("SelectJustification: %v", err)
		}
		buffer = append(buffer, justCmd...)

		// Set tab positions
		tabs := []byte{20, 40, 60, 80}
		tabCmd, err := cmd.SetHorizontalTabPositions(tabs)
		if err != nil {
			t.Fatalf("SetHorizontalTabPositions: %v", err)
		}
		buffer = append(buffer, tabCmd...)

		// Position commands
		buffer = append(buffer, cmd.SetAbsolutePrintPosition(100)...)
		buffer = append(buffer, cmd.HorizontalTab()...)
		buffer = append(buffer, cmd.SetRelativePrintPosition(50)...)

		// Move to beginning of line
		beginCmd, err := cmd.SetPrintPositionBeginningLine(printposition.Print)
		if err != nil {
			t.Fatalf("SetPrintPositionBeginningLine: %v", err)
		}
		buffer = append(buffer, beginCmd...)

		// Verify buffer has expected commands
		if len(buffer) == 0 {
			t.Error("Buffer should contain commands")
		}

		// Verify specific command sequences
		expectedStart := []byte{common.GS, 'L', 50, 0} // SetLeftMargin(50)
		if !bytes.Equal(buffer[:4], expectedStart) {
			t.Errorf("Buffer should start with SetLeftMargin command")
		}
	})

	t.Run("justification modes", func(t *testing.T) {
		var buffer []byte

		// Test all justification modes
		modes := []struct {
			name string
			mode printposition.Justification
		}{
			{"left", printposition.Left},
			{"center", printposition.Center},
			{"right", printposition.Right},
			{"left ASCII", printposition.LeftASCII},
			{"center ASCII", printposition.CenterASCII},
			{"right ASCII", printposition.RightASCII},
		}

		for _, m := range modes {
			justCmd, err := cmd.SelectJustification(m.mode)
			if err != nil {
				t.Errorf("SelectJustification(%s): %v", m.name, err)
				continue
			}
			buffer = append(buffer, justCmd...)
		}

		// Verify we have all justification commands
		expectedCmdCount := 6 * 3 // 6 modes × 3 bytes per command
		if len(buffer) != expectedCmdCount {
			t.Errorf("Buffer length = %d, want %d", len(buffer), expectedCmdCount)
		}
	})

	t.Run("tab position workflow", func(t *testing.T) {
		var buffer []byte

		// Clear tabs first
		clearCmd, err := cmd.SetHorizontalTabPositions([]byte{})
		if err != nil {
			t.Fatalf("Clear tabs: %v", err)
		}
		buffer = append(buffer, clearCmd...)

		// Set custom tabs
		customTabs := []byte{10, 25, 40, 55, 70}
		setCmd, err := cmd.SetHorizontalTabPositions(customTabs)
		if err != nil {
			t.Fatalf("Set custom tabs: %v", err)
		}
		buffer = append(buffer, setCmd...)

		// Use tabs
		buffer = append(buffer, cmd.HorizontalTab()...)
		buffer = append(buffer, cmd.HorizontalTab()...)
		buffer = append(buffer, cmd.HorizontalTab()...)

		// Verify tab usage
		htCount := bytes.Count(buffer, []byte{common.HT})
		if htCount != 3 {
			t.Errorf("HT count = %d, want 3", htCount)
		}
	})

	t.Run("relative and absolute positioning mix", func(t *testing.T) {
		var buffer []byte

		// Start at absolute position
		buffer = append(buffer, cmd.SetAbsolutePrintPosition(100)...)

		// Move relatively forward
		buffer = append(buffer, cmd.SetRelativePrintPosition(50)...)

		// Move relatively backward
		buffer = append(buffer, cmd.SetRelativePrintPosition(-30)...)

		// Reset to absolute position
		buffer = append(buffer, cmd.SetAbsolutePrintPosition(0)...)

		// Verify command sequence
		if len(buffer) != 4*4 { // 4 commands × 4 bytes each
			t.Errorf("Buffer length = %d, want 16", len(buffer))
		}

		// Check first absolute position command
		expected := []byte{common.ESC, '$', 100, 0}
		if !bytes.Equal(buffer[:4], expected) {
			t.Errorf("First command = %#v, want %#v", buffer[:4], expected)
		}
	})
}

func TestIntegration_PrintPosition_PageModeWorkflow(t *testing.T) {
	cmd := printposition.NewCommands()

	t.Run("complete page mode setup", func(t *testing.T) {
		var buffer []byte

		// Set print direction
		dirCmd, err := cmd.SelectPrintDirectionPageMode(printposition.LeftToRight)
		if err != nil {
			t.Fatalf("SelectPrintDirectionPageMode: %v", err)
		}
		buffer = append(buffer, dirCmd...)

		// Set print area
		areaCmd, err := cmd.SetPrintAreaPageMode(10, 20, 300, 400)
		if err != nil {
			t.Fatalf("SetPrintAreaPageMode: %v", err)
		}
		buffer = append(buffer, areaCmd...)

		// Set absolute positions
		buffer = append(buffer, cmd.SetAbsolutePrintPosition(50)...)
		buffer = append(buffer, cmd.SetAbsoluteVerticalPrintPosition(100)...)

		// Set relative positions
		buffer = append(buffer, cmd.SetRelativePrintPosition(25)...)
		buffer = append(buffer, cmd.SetRelativeVerticalPrintPosition(-50)...)

		// Verify buffer contains all commands
		if len(buffer) == 0 {
			t.Error("Buffer should contain commands")
		}

		// Verify print direction command
		expectedDir := []byte{common.ESC, 'T', byte(printposition.LeftToRight)}
		if !bytes.Equal(buffer[:3], expectedDir) {
			t.Errorf("First command = %#v, want %#v", buffer[:3], expectedDir)
		}
	})

	t.Run("all print directions", func(t *testing.T) {
		directions := []struct {
			name      string
			direction printposition.PrintDirection
		}{
			{"left to right", printposition.LeftToRight},
			{"bottom to top", printposition.BottomToTop},
			{"right to left", printposition.RightToLeft},
			{"top to bottom", printposition.TopToBottom},
			{"left to right ASCII", printposition.LeftToRightASCII},
			{"bottom to top ASCII", printposition.BottomToTopASCII},
			{"right to left ASCII", printposition.RightToLeftASCII},
			{"top to bottom ASCII", printposition.TopToBottomASCII},
		}

		for _, d := range directions {
			t.Run(d.name, func(t *testing.T) {
				cmd, err := cmd.SelectPrintDirectionPageMode(d.direction)
				if err != nil {
					t.Errorf("SelectPrintDirectionPageMode(%s): %v", d.name, err)
					return
				}

				expected := []byte{common.ESC, 'T', byte(d.direction)}
				if !bytes.Equal(cmd, expected) {
					t.Errorf("Command = %#v, want %#v", cmd, expected)
				}
			})
		}
	})

	t.Run("complex page mode positioning", func(t *testing.T) {
		var buffer []byte

		// Setup page mode area at specific location
		areaCmd, err := cmd.SetPrintAreaPageMode(50, 50, 200, 300)
		if err != nil {
			t.Fatalf("SetPrintAreaPageMode: %v", err)
		}
		buffer = append(buffer, areaCmd...)

		// Move to center of print area
		buffer = append(buffer, cmd.SetAbsolutePrintPosition(100)...)         // Horizontal center
		buffer = append(buffer, cmd.SetAbsoluteVerticalPrintPosition(150)...) // Vertical center

		// Move in a pattern
		moves := []struct {
			horizontal int16
			vertical   int16
		}{
			{20, 30},
			{-40, 0},
			{20, -30},
			{0, 60},
		}

		for _, move := range moves {
			if move.horizontal != 0 {
				buffer = append(buffer, cmd.SetRelativePrintPosition(move.horizontal)...)
			}
			if move.vertical != 0 {
				buffer = append(buffer, cmd.SetRelativeVerticalPrintPosition(move.vertical)...)
			}
		}

		// Verify we have the expected number of commands
		// 1 area + 2 absolute + 6 relative = 9 total position commands
		if len(buffer) < 10 {
			t.Error("Buffer should contain multiple positioning commands")
		}
	})

	t.Run("page mode with different starting positions", func(t *testing.T) {
		testCases := []struct {
			name      string
			direction printposition.PrintDirection
			x, y      uint16
			width     uint16
			height    uint16
		}{
			{"upper left start", printposition.LeftToRight, 0, 0, 100, 100},
			{"lower left start", printposition.BottomToTop, 0, 200, 100, 100},
			{"lower right start", printposition.RightToLeft, 200, 200, 100, 100},
			{"upper right start", printposition.TopToBottom, 200, 0, 100, 100},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				var buffer []byte

				// Set direction
				dirCmd, err := cmd.SelectPrintDirectionPageMode(tc.direction)
				if err != nil {
					t.Fatalf("SelectPrintDirectionPageMode: %v", err)
				}
				buffer = append(buffer, dirCmd...)

				// Set print area
				areaCmd, err := cmd.SetPrintAreaPageMode(tc.x, tc.y, tc.width, tc.height)
				if err != nil {
					t.Fatalf("SetPrintAreaPageMode: %v", err)
				}
				buffer = append(buffer, areaCmd...)

				// Verify commands were generated
				if len(buffer) == 0 {
					t.Error("Buffer should contain commands")
				}
			})
		}
	})
}

func TestIntegration_PrintPosition_MixedModeTransitions(t *testing.T) {
	cmd := printposition.NewCommands()

	t.Run("standard to page mode transition", func(t *testing.T) {
		var buffer []byte

		// Standard mode setup
		buffer = append(buffer, cmd.SetLeftMargin(20)...)
		buffer = append(buffer, cmd.SetPrintAreaWidth(500)...)
		justCmd, _ := cmd.SelectJustification(printposition.Left)
		buffer = append(buffer, justCmd...)

		// Standard mode positioning
		buffer = append(buffer, cmd.SetAbsolutePrintPosition(100)...)
		buffer = append(buffer, cmd.SetRelativePrintPosition(50)...)

		// Page mode commands (would be after switching to page mode)
		dirCmd, _ := cmd.SelectPrintDirectionPageMode(printposition.LeftToRight)
		buffer = append(buffer, dirCmd...)

		areaCmd, _ := cmd.SetPrintAreaPageMode(0, 0, 400, 600)
		buffer = append(buffer, areaCmd...)

		buffer = append(buffer, cmd.SetAbsoluteVerticalPrintPosition(200)...)

		// Verify mixed commands are present
		hasStandardCmd := bytes.Contains(buffer, []byte{common.GS, 'L'}) // SetLeftMargin
		hasPageCmd := bytes.Contains(buffer, []byte{common.ESC, 'T'})    // SelectPrintDirectionPageMode

		if !hasStandardCmd {
			t.Error("Buffer should contain standard mode commands")
		}
		if !hasPageCmd {
			t.Error("Buffer should contain page mode commands")
		}
	})

	t.Run("beginning of line operations", func(t *testing.T) {
		var buffer []byte

		// Position somewhere
		buffer = append(buffer, cmd.SetAbsolutePrintPosition(200)...)

		// Move to beginning and erase
		eraseCmd, err := cmd.SetPrintPositionBeginningLine(printposition.Erase)
		if err != nil {
			t.Fatalf("SetPrintPositionBeginningLine(erase): %v", err)
		}
		buffer = append(buffer, eraseCmd...)

		// Position again
		buffer = append(buffer, cmd.SetAbsolutePrintPosition(150)...)

		// Move to beginning and print
		printCmd, err := cmd.SetPrintPositionBeginningLine(printposition.Print)
		if err != nil {
			t.Fatalf("SetPrintPositionBeginningLine(print): %v", err)
		}
		buffer = append(buffer, printCmd...)

		// Also utils ASCII variants
		eraseASCIICmd, err := cmd.SetPrintPositionBeginningLine(printposition.EraseASCII)
		if err != nil {
			t.Fatalf("SetPrintPositionBeginningLine(erase ASCII): %v", err)
		}
		buffer = append(buffer, eraseASCIICmd...)

		printASCIICmd, err := cmd.SetPrintPositionBeginningLine(printposition.PrintASCII)
		if err != nil {
			t.Fatalf("SetPrintPositionBeginningLine(print ASCII): %v", err)
		}
		buffer = append(buffer, printASCIICmd...)

		// Verify we have the expected commands
		beginLineCount := bytes.Count(buffer, []byte{common.GS, 'T'})
		if beginLineCount != 4 {
			t.Errorf("Begin line command count = %d, want 4", beginLineCount)
		}
	})
}

func TestIntegration_PrintPosition_EdgeCases(t *testing.T) {
	cmd := printposition.NewCommands()

	t.Run("maximum values", func(t *testing.T) {
		var buffer []byte

		// Maximum positions
		buffer = append(buffer, cmd.SetAbsolutePrintPosition(65535)...)
		buffer = append(buffer, cmd.SetAbsoluteVerticalPrintPosition(65535)...)

		// Maximum margins and widths
		buffer = append(buffer, cmd.SetLeftMargin(65535)...)
		buffer = append(buffer, cmd.SetPrintAreaWidth(65535)...)

		// Maximum page mode area
		areaCmd, err := cmd.SetPrintAreaPageMode(65535, 65535, 65535, 65535)
		if err != nil {
			t.Fatalf("SetPrintAreaPageMode with max values: %v", err)
		}
		buffer = append(buffer, areaCmd...)

		// Maximum relative movements
		buffer = append(buffer, cmd.SetRelativePrintPosition(32767)...)
		buffer = append(buffer, cmd.SetRelativeVerticalPrintPosition(32767)...)

		// Minimum (negative maximum) relative movements
		buffer = append(buffer, cmd.SetRelativePrintPosition(-32768)...)
		buffer = append(buffer, cmd.SetRelativeVerticalPrintPosition(-32768)...)

		// Verify all commands were generated
		if len(buffer) == 0 {
			t.Error("Buffer should contain commands")
		}
	})

	t.Run("minimum values", func(t *testing.T) {
		var buffer []byte

		// Minimum positions (zero)
		buffer = append(buffer, cmd.SetAbsolutePrintPosition(0)...)
		buffer = append(buffer, cmd.SetAbsoluteVerticalPrintPosition(0)...)

		// Minimum margins and widths
		buffer = append(buffer, cmd.SetLeftMargin(0)...)
		buffer = append(buffer, cmd.SetPrintAreaWidth(1)...) // Minimum non-zero

		// Minimum page mode area
		areaCmd, err := cmd.SetPrintAreaPageMode(0, 0, 1, 1)
		if err != nil {
			t.Fatalf("SetPrintAreaPageMode with min values: %v", err)
		}
		buffer = append(buffer, areaCmd...)

		// Zero relative movements
		buffer = append(buffer, cmd.SetRelativePrintPosition(0)...)
		buffer = append(buffer, cmd.SetRelativeVerticalPrintPosition(0)...)

		// Verify all commands were generated
		if len(buffer) == 0 {
			t.Error("Buffer should contain commands")
		}
	})

	t.Run("tab position edge cases", func(t *testing.T) {
		// Maximum number of tabs (32)
		maxTabs := make([]byte, 32)
		for i := range maxTabs {
			maxTabs[i] = byte((i + 1) * 8) // Every 8 units
		}
		maxTabs[31] = 255 // Last tab at maximum position

		maxCmd, err := cmd.SetHorizontalTabPositions(maxTabs)
		if err != nil {
			t.Fatalf("SetHorizontalTabPositions with 32 tabs: %v", err)
		}

		// Clear all tabs
		clearCmd, err := cmd.SetHorizontalTabPositions([]byte{})
		if err != nil {
			t.Fatalf("SetHorizontalTabPositions clear: %v", err)
		}

		// Single tab at maximum position
		singleMaxTab, err := cmd.SetHorizontalTabPositions([]byte{255})
		if err != nil {
			t.Fatalf("SetHorizontalTabPositions with max position: %v", err)
		}

		// Verify commands
		if len(maxCmd) != 32+3 { // 32 tabs + ESC D + NUL
			t.Errorf("Max tabs command length = %d, want %d", len(maxCmd), 35)
		}
		if len(clearCmd) != 3 { // ESC D NUL
			t.Errorf("Clear tabs command length = %d, want 3", len(clearCmd))
		}
		if len(singleMaxTab) != 4 { // ESC D 255 NUL
			t.Errorf("Single max tab command length = %d, want 4", len(singleMaxTab))
		}
	})
}

func TestIntegration_PrintPosition_ErrorConditions(t *testing.T) {
	cmd := printposition.NewCommands()

	t.Run("invalid parameters", func(t *testing.T) {
		// Invalid justification
		_, err := cmd.SelectJustification(99)
		if err == nil {
			t.Error("SelectJustification(99) should return error")
		}

		// Invalid print direction
		_, err = cmd.SelectPrintDirectionPageMode(99)
		if err == nil {
			t.Error("SelectPrintDirectionPageMode(99) should return error")
		}

		// Invalid begin line mode
		_, err = cmd.SetPrintPositionBeginningLine(99)
		if err == nil {
			t.Error("SetPrintPositionBeginningLine(99) should return error")
		}

		// Zero width in page mode area
		_, err = cmd.SetPrintAreaPageMode(0, 0, 0, 100)
		if err == nil {
			t.Error("SetPrintAreaPageMode with zero width should return error")
		}

		// Zero height in page mode area
		_, err = cmd.SetPrintAreaPageMode(0, 0, 100, 0)
		if err == nil {
			t.Error("SetPrintAreaPageMode with zero height should return error")
		}

		// Too many tab positions
		tooManyTabs := make([]byte, 33)
		for i := range tooManyTabs {
			tooManyTabs[i] = byte(i + 1)
		}
		_, err = cmd.SetHorizontalTabPositions(tooManyTabs)
		if err == nil {
			t.Error("SetHorizontalTabPositions with 33 tabs should return error")
		}

		// Invalid tab positions (not ascending)
		_, err = cmd.SetHorizontalTabPositions([]byte{20, 10, 30})
		if err == nil {
			t.Error("SetHorizontalTabPositions with non-ascending values should return error")
		}

		// Invalid tab position (zero)
		_, err = cmd.SetHorizontalTabPositions([]byte{0, 10, 20})
		if err == nil {
			t.Error("SetHorizontalTabPositions with zero value should return error")
		}
	})
}

func TestIntegration_PrintPosition_RealWorldScenarios(t *testing.T) {
	cmd := printposition.NewCommands()

	t.Run("receipt header alignment", func(t *testing.T) {
		var buffer []byte

		// Center company name
		centerCmd, _ := cmd.SelectJustification(printposition.Center)
		buffer = append(buffer, centerCmd...)
		// ... print company name ...

		// Left align address
		leftCmd, _ := cmd.SelectJustification(printposition.Left)
		buffer = append(buffer, leftCmd...)
		// ... print address ...

		// Right align date/time
		rightCmd, _ := cmd.SelectJustification(printposition.Right)
		buffer = append(buffer, rightCmd...)
		// ... print date/time ...

		// Reset to left
		resetCmd, _ := cmd.SelectJustification(printposition.Left)
		buffer = append(buffer, resetCmd...)

		// Verify justification changes
		if len(buffer) != 4*3 { // 4 commands × 3 bytes each
			t.Errorf("Buffer length = %d, want 12", len(buffer))
		}
	})

	t.Run("tabulated item list", func(t *testing.T) {
		var buffer []byte

		// Set up tabs for: Qty | Description | Price | Total
		// Assuming 58mm paper (420 dots)
		tabs := []byte{5, 10, 35, 42} // Character positions
		tabCmd, _ := cmd.SetHorizontalTabPositions(tabs)
		buffer = append(buffer, tabCmd...)

		// Print header row with tabs
		for i := 0; i < 3; i++ {
			buffer = append(buffer, cmd.HorizontalTab()...)
		}

		// Print item rows
		for item := 0; item < 5; item++ {
			// Reset to beginning
			beginCmd, _ := cmd.SetPrintPositionBeginningLine(printposition.Print)
			buffer = append(buffer, beginCmd...)

			// Tab to each column
			for col := 0; col < 3; col++ {
				buffer = append(buffer, cmd.HorizontalTab()...)
			}
		}

		// Count total tabs used
		tabCount := bytes.Count(buffer, []byte{common.HT})
		expectedTabs := 3 + (5 * 3) // Header + 5 items
		if tabCount != expectedTabs {
			t.Errorf("Tab count = %d, want %d", tabCount, expectedTabs)
		}
	})

	t.Run("barcode positioning", func(t *testing.T) {
		var buffer []byte

		// Set narrow print area for product label
		buffer = append(buffer, cmd.SetLeftMargin(50)...)
		buffer = append(buffer, cmd.SetPrintAreaWidth(300)...)

		// Center barcode
		centerCmd, _ := cmd.SelectJustification(printposition.Center)
		buffer = append(buffer, centerCmd...)

		// Position for barcode
		buffer = append(buffer, cmd.SetAbsolutePrintPosition(150)...)
		// ... print barcode ...

		// Reset for text below
		leftCmd, _ := cmd.SelectJustification(printposition.Left)
		buffer = append(buffer, leftCmd...)
		beginCmd, _ := cmd.SetPrintPositionBeginningLine(printposition.Print)
		buffer = append(buffer, beginCmd...)

		// Verify positioning sequence
		if len(buffer) < 5 {
			t.Error("Buffer should contain multiple positioning commands")
		}
	})

	t.Run("multi-column layout in page mode", func(t *testing.T) {
		var buffer []byte

		// Set up page mode for two-column layout
		dirCmd, _ := cmd.SelectPrintDirectionPageMode(printposition.LeftToRight)
		buffer = append(buffer, dirCmd...)

		// Left column
		leftColCmd, _ := cmd.SetPrintAreaPageMode(0, 0, 200, 400)
		buffer = append(buffer, leftColCmd...)
		// ... print left column content ...

		// Right column
		rightColCmd, _ := cmd.SetPrintAreaPageMode(220, 0, 200, 400)
		buffer = append(buffer, rightColCmd...)
		// ... print right column content ...

		// Verify we have both column definitions
		if len(buffer) < 2*10 { // At least 2 SetPrintAreaPageMode commands
			t.Error("Buffer should contain column definitions")
		}
	})
}
