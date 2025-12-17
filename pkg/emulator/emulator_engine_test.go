package emulator_test

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"testing"

	"github.com/adcondev/poster/pkg/constants"
	"github.com/adcondev/poster/pkg/emulator"
)

// ============================================================================
// Engine Creation Tests
// ============================================================================

func TestNewEngine_Default(t *testing.T) {
	engine, err := emulator.NewDefaultEngine()
	if err != nil {
		t.Fatalf("NewDefaultEngine() error = %v", err)
	}
	if engine == nil {
		t.Fatal("NewDefaultEngine() returned nil")
	}

	state := engine.State()
	if state.PaperPxWidth != constants.PaperPxWidth80mm {
		t.Errorf("State().PaperPxWidth = %d, want %d", state.PaperPxWidth, constants.PaperPxWidth80mm)
	}
	if state.DPI != constants.DefaultDPI {
		t.Errorf("State().DPI = %d, want %d", state.DPI, constants.DefaultDPI)
	}
}

func TestNewEngine_58mm(t *testing.T) {
	engine, err := emulator.New58mmEngine()
	if err != nil {
		t.Fatalf("New58mmEngine() error = %v", err)
	}

	state := engine.State()
	if state.PaperPxWidth != constants.PaperPxWidth58mm {
		t.Errorf("State().PaperPxWidth = %d, want %d", state.PaperPxWidth, constants.PaperPxWidth58mm)
	}
}

func TestNewEngine_InvalidConfig(t *testing.T) {
	tests := []struct {
		name   string
		config emulator.Config
	}{
		{"zero width", emulator.Config{PaperPxWidth: 0}},
		{"negative width", emulator.Config{PaperPxWidth: -100}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := emulator.NewEngine(tt.config)
			if err == nil {
				t.Errorf("NewEngine() with %s should return error", tt.name)
			}
		})
	}
}

func TestNewEngine_DefaultDPI(t *testing.T) {
	// When DPI is 0 or negative, should default to 203
	config := emulator.Config{PaperPxWidth: 576, DPI: 0}
	engine, err := emulator.NewEngine(config)
	if err != nil {
		t.Fatalf("NewEngine() error = %v", err)
	}
	if engine.State().DPI != constants.DefaultDPI {
		t.Errorf("DPI = %d, want %d", engine.State().DPI, constants.DefaultDPI)
	}
}

// ============================================================================
// Text Rendering Tests
// ============================================================================

func TestPrintLine_AdvancesCursor(t *testing.T) {
	engine, _ := emulator.NewDefaultEngine()
	initialY := engine.State().CursorY

	engine.PrintLine("Test")

	if engine.State().CursorY <= initialY {
		t.Error("PrintLine() should advance CursorY")
	}
}

func TestPrintLine_EmptyString(t *testing.T) {
	engine, _ := emulator.NewDefaultEngine()
	initialY := engine.State().CursorY

	engine.PrintLine("")

	// Empty string should still advance line
	if engine.State().CursorY <= initialY {
		t.Error("PrintLine(\"\") should still advance cursor")
	}
}

func TestPrint_DoesNotAdvanceLine(t *testing.T) {
	engine, _ := emulator.NewDefaultEngine()
	initialY := engine.State().CursorY

	engine.Print("Test")

	if engine.State().CursorY != initialY {
		t.Error("Print() should not advance CursorY")
	}
}

func TestPrint_UpdatesCursorX(t *testing.T) {
	engine, _ := emulator.NewDefaultEngine()
	engine.AlignLeft()

	engine.Print("Test")

	if engine.State().CursorX <= 0 {
		t.Error("Print() should update CursorX")
	}
}

// ============================================================================
// Alignment Tests
// ============================================================================

func TestAlignment_AllModes(t *testing.T) {
	tests := []struct {
		name   string
		set    func(*emulator.Engine)
		expect string
	}{
		{"Left", func(e *emulator.Engine) { e.AlignLeft() }, constants.Left.String()},
		{"Center", func(e *emulator.Engine) { e.AlignCenter() }, constants.Center.String()},
		{"Right", func(e *emulator.Engine) { e.AlignRight() }, constants.Right.String()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine, _ := emulator.NewDefaultEngine()
			tt.set(engine)
			if engine.State().Align != tt.expect {
				t.Errorf("Align = %s, want %s", engine.State().Align, tt.expect)
			}
		})
	}
}

func TestSetAlign_Normalization(t *testing.T) {
	tests := []struct {
		input  string
		expect string
	}{
		{"CENTER", constants.Center.String()},
		{"  center  ", constants.Center.String()},
		{"Center", constants.Center.String()},
		{"Right", constants.Right.String()},
		{"RIGHT", constants.Right.String()},
		{"  right  ", constants.Right.String()},
		{"left", constants.Left.String()},
		{"LEFT", constants.Left.String()},
		{"invalid", constants.Left.String()}, // Fallback
		{"", constants.Left.String()},        // Fallback
		{"unknown", constants.Left.String()}, // Fallback
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			engine, _ := emulator.NewDefaultEngine()
			engine.SetAlign(tt.input)
			if engine.State().Align != tt.expect {
				t.Errorf("SetAlign(%q) = %s, want %s", tt.input, engine.State().Align, tt.expect)
			}
		})
	}
}

// ============================================================================
// Font Tests
// ============================================================================

func TestSetFont_Normalization(t *testing.T) {
	tests := []struct {
		input  string
		expect string
	}{
		{"A", "A"},
		{"a", "A"},
		{"  A  ", "A"},
		{"B", "B"},
		{"b", "B"},
		{"  b  ", "B"},
		{"C", "A"},  // Invalid falls back to A
		{"", "A"},   // Empty falls back to A
		{"AB", "A"}, // Invalid falls back to A
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			engine, _ := emulator.NewDefaultEngine()
			engine.SetFont(tt.input)
			if engine.State().FontName != tt.expect {
				t.Errorf("SetFont(%q) = %s, want %s", tt.input, engine.State().FontName, tt.expect)
			}
		})
	}
}

// ============================================================================
// Style Tests
// ============================================================================

func TestSetSize_Clamping(t *testing.T) {
	tests := []struct {
		inputW, inputH   int
		expectW, expectH float64
	}{
		{1, 1, 1.0, 1.0},
		{8, 8, 8.0, 8.0},
		{0, 0, 1.0, 1.0},     // Clamped up
		{-1, -1, 1.0, 1.0},   // Clamped up
		{10, 10, 8.0, 8.0},   // Clamped down
		{100, 100, 8.0, 8.0}, // Clamped down
		{-1, 5, 1.0, 5.0},    // Mixed
		{5, -1, 5.0, 1.0},    // Mixed
		{3, 4, 3.0, 4.0},     // Normal
	}

	for _, tt := range tests {
		engine, _ := emulator.NewDefaultEngine()
		engine.SetSize(tt.inputW, tt.inputH)
		state := engine.State()

		if state.ScaleW != tt.expectW || state.ScaleH != tt.expectH {
			t.Errorf("SetSize(%d, %d) = (%.1f, %.1f), want (%.1f, %.1f)",
				tt.inputW, tt.inputH, state.ScaleW, state.ScaleH, tt.expectW, tt.expectH)
		}
	}
}

func TestSetUnderline_Clamping(t *testing.T) {
	tests := []struct {
		input  int
		expect int
	}{
		{0, 0},
		{1, 1},
		{2, 2},
		{-1, 0},
		{-100, 0},
		{3, 2},
		{5, 2},
		{100, 2},
	}

	for _, tt := range tests {
		engine, _ := emulator.NewDefaultEngine()
		engine.SetUnderline(tt.input)
		if engine.State().IsUnderline != tt.expect {
			t.Errorf("SetUnderline(%d) = %d, want %d", tt.input, engine.State().IsUnderline, tt.expect)
		}
	}
}

func TestSetBold(t *testing.T) {
	engine, _ := emulator.NewDefaultEngine()

	engine.SetBold(true)
	if !engine.State().IsBold {
		t.Error("SetBold(true) should set IsBold to true")
	}

	engine.SetBold(false)
	if engine.State().IsBold {
		t.Error("SetBold(false) should set IsBold to false")
	}
}

func TestSetInverse(t *testing.T) {
	engine, _ := emulator.NewDefaultEngine()

	engine.SetInverse(true)
	if !engine.State().IsInverse {
		t.Error("SetInverse(true) should set IsInverse to true")
	}

	engine.SetInverse(false)
	if engine.State().IsInverse {
		t.Error("SetInverse(false) should set IsInverse to false")
	}
}

// ============================================================================
// SetSize Auto-Cursor Adjustment Tests
// ============================================================================

func TestSetSize_AutoAdjustCursor_Enabled(t *testing.T) {
	// Default config has AutoAdjustCursorOnScale = true
	engine, err := emulator.NewDefaultEngine()
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	// Print some text first
	engine.PrintLine("Normal text")
	cursorAfterText := engine.State().CursorY

	// Scale up
	engine.SetSize(2, 2)
	cursorAfterScale := engine.State().CursorY

	// Cursor should have advanced
	if cursorAfterScale <= cursorAfterText {
		t.Errorf("Cursor should advance when scaling up with AutoAdjustCursorOnScale=true.  Before: %.1f, After:  %.1f",
			cursorAfterText, cursorAfterScale)
	}

	// The advance should be approximately (2-1) * FontAHeight = 24 pixels
	expectedAdvance := float64(constants.FontAHeight) * (2.0 - 1.0)
	actualAdvance := cursorAfterScale - cursorAfterText

	if actualAdvance != expectedAdvance {
		t.Errorf("Cursor advance = %.1f, want %.1f", actualAdvance, expectedAdvance)
	}
}

func TestSetSize_AutoAdjustCursor_Disabled(t *testing.T) {
	config := emulator.DefaultConfig()
	config.AutoAdjustCursorOnScale = false

	engine, err := emulator.NewEngine(config)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	// Print some text first
	engine.PrintLine("Normal text")
	cursorAfterText := engine.State().CursorY

	// Scale up
	engine.SetSize(2, 2)
	cursorAfterScale := engine.State().CursorY

	// Cursor should NOT have advanced
	if cursorAfterScale != cursorAfterText {
		t.Errorf("Cursor should not advance when AutoAdjustCursorOnScale=false. Before: %.1f, After: %.1f",
			cursorAfterText, cursorAfterScale)
	}
}

func TestSetSize_NoAdjustOnScaleDown(t *testing.T) {
	engine, err := emulator.NewDefaultEngine()
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	// Start with scaled text
	engine.SetSize(3, 3)
	engine.PrintLine("Big text")
	cursorAfterText := engine.State().CursorY

	// Scale down
	engine.SetSize(1, 1)
	cursorAfterScale := engine.State().CursorY

	// Cursor should NOT change when scaling down
	if cursorAfterScale != cursorAfterText {
		t.Errorf("Cursor should not change when scaling down. Before: %.1f, After: %.1f",
			cursorAfterText, cursorAfterScale)
	}
}

func TestSetSize_NoAdjustOnSameScale(t *testing.T) {
	engine, err := emulator.NewDefaultEngine()
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	engine.SetSize(2, 2)
	engine.PrintLine("Text")
	cursorAfterText := engine.State().CursorY

	// Set same size again
	engine.SetSize(2, 2)
	cursorAfterScale := engine.State().CursorY

	// Cursor should NOT change
	if cursorAfterScale != cursorAfterText {
		t.Errorf("Cursor should not change when setting same scale. Before: %.1f, After: %.1f",
			cursorAfterText, cursorAfterScale)
	}
}

func TestSetSize_MultipleScaleUps(t *testing.T) {
	engine, err := emulator.NewDefaultEngine()
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	engine.PrintLine("Text")
	cursor1 := engine.State().CursorY

	engine.SetSize(2, 2) // 1 -> 2
	cursor2 := engine.State().CursorY

	engine.SetSize(4, 4) // 2 -> 4
	cursor3 := engine.State().CursorY

	// Each scale up should advance cursor
	if cursor2 <= cursor1 {
		t.Error("First scale up should advance cursor")
	}
	if cursor3 <= cursor2 {
		t.Error("Second scale up should advance cursor")
	}

	// Verify the advances are proportional to scale changes
	advance1to2 := cursor2 - cursor1
	advance2to4 := cursor3 - cursor2

	// 1->2 advances by 1*24=24, 2->4 advances by 2*24=48
	expectedAdvance1 := float64(constants.FontAHeight) * 1.0
	expectedAdvance2 := float64(constants.FontAHeight) * 2.0

	if advance1to2 != expectedAdvance1 {
		t.Errorf("1->2 advance = %.1f, want %.1f", advance1to2, expectedAdvance1)
	}
	if advance2to4 != expectedAdvance2 {
		t.Errorf("2->4 advance = %.1f, want %.1f", advance2to4, expectedAdvance2)
	}
}

func TestSetSize_WidthOnlyChange(t *testing.T) {
	engine, err := emulator.NewDefaultEngine()
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	engine.PrintLine("Text")
	cursorBefore := engine.State().CursorY

	// Only change width, not height
	engine.SetSize(3, 1)
	cursorAfter := engine.State().CursorY

	// No height change means no cursor adjustment
	if cursorAfter != cursorBefore {
		t.Errorf("Cursor should not change when only width changes. Before: %.1f, After: %.1f",
			cursorBefore, cursorAfter)
	}
}

// ============================================================================
// Feed and Paper Control Tests
// ============================================================================

func TestFeed_MultipleLines(t *testing.T) {
	engine, _ := emulator.NewDefaultEngine()
	initialY := engine.State().CursorY

	engine.Feed(5)
	afterFeed := engine.State().CursorY

	expectedAdvance := float64(5) * float64(constants.DefaultLineSpacing)
	actualAdvance := afterFeed - initialY

	// Allow some tolerance for line height calculations
	if actualAdvance < expectedAdvance*0.8 {
		t.Errorf("Feed(5) advanced %.1f, expected at least %.1f", actualAdvance, expectedAdvance*0.8)
	}
}

func TestFeed_ZeroOrNegative(t *testing.T) {
	engine, _ := emulator.NewDefaultEngine()
	initialY := engine.State().CursorY

	engine.Feed(0)
	afterFeed := engine.State().CursorY

	// Should still advance at least one line
	if afterFeed <= initialY {
		t.Error("Feed(0) should still advance cursor (default to 1 line)")
	}
}

func TestCut_UpdatesMaxY(t *testing.T) {
	engine, _ := emulator.NewDefaultEngine()
	engine.PrintLine("Before cut")
	beforeCut := engine.State().CursorY

	engine.Cut(true)

	if engine.State().CursorY <= beforeCut {
		t.Error("Cut() should advance cursor position")
	}
}

func TestCut_PartialVsFull(t *testing.T) {
	// Both should work without error and advance cursor
	for _, partial := range []bool{true, false} {
		engine, _ := emulator.NewDefaultEngine()
		initialY := engine.State().CursorY

		engine.Cut(partial)

		if engine.State().CursorY <= initialY {
			t.Errorf("Cut(partial=%v) should advance cursor", partial)
		}
	}
}

// ============================================================================
// Reset Tests
// ============================================================================

func TestReset_RestoresDefaults(t *testing.T) {
	engine, _ := emulator.NewDefaultEngine()

	// Modify everything
	engine.SetBold(true)
	engine.SetInverse(true)
	engine.SetUnderline(2)
	engine.SetSize(4, 4)
	engine.AlignCenter()
	engine.SetFont("B")
	engine.PrintLine("Modified state")

	engine.Reset()

	state := engine.State()
	if state.IsBold {
		t.Error("Reset() should clear IsBold")
	}
	if state.IsInverse {
		t.Error("Reset() should clear IsInverse")
	}
	if state.IsUnderline != 0 {
		t.Error("Reset() should clear IsUnderline")
	}
	if state.ScaleW != 1.0 || state.ScaleH != 1.0 {
		t.Error("Reset() should reset size to 1x1")
	}
	if state.Align != constants.Left.String() {
		t.Error("Reset() should reset alignment to left")
	}
	if state.FontName != "A" {
		t.Error("Reset() should reset font to A")
	}
}

// ============================================================================
// Output Tests
// ============================================================================

func TestWritePNG_ValidOutput(t *testing.T) {
	engine, _ := emulator.NewDefaultEngine()
	engine.PrintLine("Test Receipt")
	engine.Separator("-", 20)
	engine.Cut(false)

	var buf bytes.Buffer
	if err := engine.WritePNG(&buf); err != nil {
		t.Fatalf("WritePNG() error = %v", err)
	}

	if buf.Len() == 0 {
		t.Error("WritePNG() produced empty output")
	}

	// Verify it's a valid PNG
	img, err := png.Decode(bytes.NewReader(buf.Bytes()))
	if err != nil {
		t.Fatalf("WritePNG() output is not valid PNG: %v", err)
	}

	bounds := img.Bounds()
	if bounds.Dx() != constants.PaperPxWidth80mm {
		t.Errorf("PNG width = %d, want %d", bounds.Dx(), constants.PaperPxWidth80mm)
	}
}

func TestRender_EmptyCanvas(t *testing.T) {
	engine, _ := emulator.NewDefaultEngine()
	// Don't print anything

	img := engine.Render()
	if img == nil {
		t.Error("Render() should not return nil even for empty canvas")
	}
}

func TestRenderWithInfo_Dimensions(t *testing.T) {
	engine, _ := emulator.NewDefaultEngine()
	engine.PrintLine("Line 1")
	engine.PrintLine("Line 2")

	result := engine.RenderWithInfo()

	if result.Width != constants.PaperPxWidth80mm {
		t.Errorf("RenderWithInfo().Width = %d, want %d", result.Width, constants.PaperPxWidth80mm)
	}
	if result.Height <= 0 {
		t.Error("RenderWithInfo().Height should be positive")
	}
	if result.Image == nil {
		t.Error("RenderWithInfo().Image should not be nil")
	}
}

// ============================================================================
// CharsPerLine Tests
// ============================================================================

func TestCharsPerLine_FontA(t *testing.T) {
	engine, _ := emulator.NewDefaultEngine()
	engine.SetFont("A")
	engine.SetSize(1, 1)

	chars := engine.CharsPerLine()

	// 576px / 12px = 48 chars for Font A on 80mm
	expected := constants.PaperPxWidth80mm / constants.FontAWidth
	if chars != expected {
		t.Errorf("CharsPerLine() with Font A = %d, want %d", chars, expected)
	}
}

func TestCharsPerLine_FontB(t *testing.T) {
	engine, _ := emulator.NewDefaultEngine()
	engine.SetFont("B")
	engine.SetSize(1, 1)

	chars := engine.CharsPerLine()

	// 576px / 9px = 64 chars for Font B on 80mm
	expected := constants.PaperPxWidth80mm / constants.FontBWidth
	if chars != expected {
		t.Errorf("CharsPerLine() with Font B = %d, want %d", chars, expected)
	}
}

func TestCharsPerLine_WithScaling(t *testing.T) {
	engine, _ := emulator.NewDefaultEngine()
	engine.SetSize(2, 2)

	chars := engine.CharsPerLine()

	// Should be half of normal
	expected := constants.PaperPxWidth80mm / (constants.FontAWidth * 2)
	if chars != expected {
		t.Errorf("CharsPerLine() with 2x scale = %d, want %d", chars, expected)
	}
}

func TestCharsPerLine_58mm(t *testing.T) {
	engine, _ := emulator.New58mmEngine()
	engine.SetFont("A")
	engine.SetSize(1, 1)

	chars := engine.CharsPerLine()

	// 384px / 12px = 32 chars for Font A on 58mm
	expected := constants.PaperPxWidth58mm / constants.FontAWidth
	if chars != expected {
		t.Errorf("CharsPerLine() on 58mm = %d, want %d", chars, expected)
	}
}

// ============================================================================
// Separator Tests
// ============================================================================

func TestSeparator_DefaultChar(t *testing.T) {
	engine, _ := emulator.NewDefaultEngine()
	initialY := engine.State().CursorY

	engine.Separator("", 20) // Empty char should default

	if engine.State().CursorY <= initialY {
		t.Error("Separator() should advance cursor")
	}
}

func TestSeparator_CustomChar(t *testing.T) {
	engine, _ := emulator.NewDefaultEngine()
	initialY := engine.State().CursorY

	engine.Separator("=", 32)

	if engine.State().CursorY <= initialY {
		t.Error("Separator() should advance cursor")
	}
}

// ============================================================================
// Constants Tests
// ============================================================================

func TestPaperWidthFromMM(t *testing.T) {
	tests := []struct {
		mm       int
		expected int
	}{
		{72, 576}, // 72mm * 8 = 576
		{48, 384}, // 48mm * 8 = 384
	}

	for _, tt := range tests {
		result := constants.PaperWidthFromMm(tt.mm)
		if result != tt.expected {
			t.Errorf("PaperWidthFromMM(%d) = %d, want %d", tt.mm, result, tt.expected)
		}
	}
}

func TestCharsPerLineForFont(t *testing.T) {
	tests := []struct {
		paperWidth int
		fontWidth  int
		expected   int
	}{
		{576, 12, 48}, // 80mm Font A
		{576, 9, 64},  // 80mm Font B
		{384, 12, 32}, // 58mm Font A
		{384, 9, 42},  // 58mm Font B
		{576, 0, 0},   // Zero font width
		{576, -1, 0},  // Negative font width
	}

	for _, tt := range tests {
		result := constants.CharsPerLineForFont(tt.paperWidth, tt.fontWidth)
		if result != tt.expected {
			t.Errorf("CharsPerLineForFont(%d, %d) = %d, want %d",
				tt.paperWidth, tt.fontWidth, result, tt.expected)
		}
	}
}

// ============================================================================
// Image Printing Methods Tests
// ============================================================================

func TestPrintImage_Basic(t *testing.T) {
	engine, err := emulator.NewDefaultEngine()
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	img := createTestImage(100, 50)
	initialY := engine.State().CursorY

	err = engine.PrintImage(img)
	if err != nil {
		t.Fatalf("PrintImage() error = %v", err)
	}

	// Cursor should advance
	if engine.State().CursorY <= initialY {
		t.Error("PrintImage() should advance CursorY")
	}
}

func TestPrintImage_NilImage(t *testing.T) {
	engine, err := emulator.NewDefaultEngine()
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	err = engine.PrintImage(nil)
	if err == nil {
		t.Error("PrintImage(nil) should return error")
	}
}

func TestPrintImage_CursorAdvanceMatchesHeight(t *testing.T) {
	engine, err := emulator.NewDefaultEngine()
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	imgHeight := 75
	img := createTestImage(100, imgHeight)
	initialY := engine.State().CursorY

	// Use options that don't resize
	opts := emulator.DefaultImageOptions()
	opts.PixelWidth = 100
	opts.PreserveAspect = true

	err = engine.PrintImageWithOptions(img, opts)
	if err != nil {
		t.Fatalf("PrintImageWithOptions() error = %v", err)
	}

	expectedY := initialY + float64(imgHeight)
	if engine.State().CursorY != expectedY {
		t.Errorf("CursorY = %.1f, want %.1f", engine.State().CursorY, expectedY)
	}
}

func TestPrintImageWithOptions_NilOptions(t *testing.T) {
	engine, err := emulator.NewDefaultEngine()
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	img := createTestImage(100, 50)

	// Should use defaults and not panic
	err = engine.PrintImageWithOptions(img, nil)
	if err != nil {
		t.Errorf("PrintImageWithOptions(img, nil) should not error:  %v", err)
	}
}

func TestPrintImageWithOptions_CustomWidth(t *testing.T) {
	engine, err := emulator.NewDefaultEngine()
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	img := createTestImage(200, 100)
	initialY := engine.State().CursorY

	opts := emulator.DefaultImageOptions()
	opts.PixelWidth = 100 // Scale down to half
	opts.PreserveAspect = true

	err = engine.PrintImageWithOptions(img, opts)
	if err != nil {
		t.Fatalf("PrintImageWithOptions() error = %v", err)
	}

	// With 2: 1 aspect ratio scaled to 100px wide, height should be 50
	expectedAdvance := 50.0
	actualAdvance := engine.State().CursorY - initialY

	if actualAdvance != expectedAdvance {
		t.Errorf("Cursor advance = %.1f, want %.1f", actualAdvance, expectedAdvance)
	}
}

func TestPrintImageWithOptions_SimulateThermal(t *testing.T) {
	engine, err := emulator.NewDefaultEngine()
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	// Create gray image
	img := createGrayTestImage(100, 50, 128)

	opts := emulator.DefaultImageOptions()
	opts.SimulateThermal = true
	opts.PixelWidth = 100

	err = engine.PrintImageWithOptions(img, opts)
	if err != nil {
		t.Fatalf("PrintImageWithOptions() with SimulateThermal error = %v", err)
	}

	// Verify output is valid
	result := engine.Render()
	if result == nil {
		t.Error("Render() should not return nil after thermal preview")
	}
}

func TestPrintImageAligned_AllAlignments(t *testing.T) {
	alignments := []string{
		constants.Left.String(),
		constants.Center.String(),
		constants.Right.String(),
	}

	for _, align := range alignments {
		t.Run(align, func(t *testing.T) {
			engine, err := emulator.NewDefaultEngine()
			if err != nil {
				t.Fatalf("Failed to create engine: %v", err)
			}

			img := createTestImage(100, 50)

			err = engine.PrintImageAligned(img, 100, align)
			if err != nil {
				t.Errorf("PrintImageAligned(img, 100, %q) error = %v", align, err)
			}
		})
	}
}

func TestPrintImageAligned_InvalidAlignment(t *testing.T) {
	engine, err := emulator.NewDefaultEngine()
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	img := createTestImage(100, 50)

	// Invalid alignment should fall back to left (not error)
	err = engine.PrintImageAligned(img, 100, "invalid")
	if err != nil {
		t.Errorf("PrintImageAligned with invalid alignment should not error: %v", err)
	}
}

func TestPrintImageAligned_ZeroWidth(t *testing.T) {
	engine, err := emulator.NewDefaultEngine()
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	img := createTestImage(100, 50)
	initialY := engine.State().CursorY

	// Zero width should use original or default
	err = engine.PrintImageAligned(img, 0, constants.Center.String())
	if err != nil {
		t.Errorf("PrintImageAligned with zero width should not error: %v", err)
	}

	// Should still render
	if engine.State().CursorY <= initialY {
		t.Error("Image should still render with zero width")
	}
}

func TestPrintImageThermalPreview_Basic(t *testing.T) {
	engine, err := emulator.NewDefaultEngine()
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	img := createTestImage(200, 100)
	initialY := engine.State().CursorY

	err = engine.PrintImageThermalPreview(img, 150)
	if err != nil {
		t.Fatalf("PrintImageThermalPreview() error = %v", err)
	}

	if engine.State().CursorY <= initialY {
		t.Error("PrintImageThermalPreview() should advance cursor")
	}
}

func TestPrintImageThermalPreview_OutputIsBW(t *testing.T) {
	engine, err := emulator.NewDefaultEngine()
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	// Create colored image
	img := createColorTestImage(100, 50)

	err = engine.PrintImageThermalPreview(img, 100)
	if err != nil {
		t.Fatalf("PrintImageThermalPreview() error = %v", err)
	}

	result := engine.Render()

	// Check that pixels in the image area are pure B&W
	allBW := true
	bounds := result.Bounds()
	for y := 0; y < bounds.Dy() && allBW; y++ {
		for x := 0; x < bounds.Dx() && allBW; x++ {
			r, g, b, _ := result.At(x, y).RGBA()
			r8, g8, b8 := r>>8, g>>8, b>>8

			isWhite := r8 == 255 && g8 == 255 && b8 == 255
			isBlack := r8 == 0 && g8 == 0 && b8 == 0

			if !isWhite && !isBlack {
				allBW = false
			}
		}
	}

	if !allBW {
		t.Error("PrintImageThermalPreview should produce only black and white pixels")
	}
}

func TestPrintImage_MultipleImages(t *testing.T) {
	engine, err := emulator.NewDefaultEngine()
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	img1 := createTestImage(100, 50)
	img2 := createTestImage(100, 75)

	err = engine.PrintImage(img1)
	if err != nil {
		t.Fatalf("First PrintImage() error = %v", err)
	}
	cursor1 := engine.State().CursorY

	err = engine.PrintImage(img2)
	if err != nil {
		t.Fatalf("Second PrintImage() error = %v", err)
	}
	cursor2 := engine.State().CursorY

	if cursor2 <= cursor1 {
		t.Error("Second image should be below first image")
	}
}

func TestPrintImage_WithTextBefore(t *testing.T) {
	engine, err := emulator.NewDefaultEngine()
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	engine.PrintLine("Header text")
	cursorAfterText := engine.State().CursorY

	img := createTestImage(100, 50)
	err = engine.PrintImage(img)
	if err != nil {
		t.Fatalf("PrintImage() error = %v", err)
	}

	if engine.State().CursorY <= cursorAfterText {
		t.Error("Image should be below text")
	}
}

func TestPrintImage_WithTextAfter(t *testing.T) {
	engine, err := emulator.NewDefaultEngine()
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	img := createTestImage(100, 50)
	err = engine.PrintImage(img)
	if err != nil {
		t.Fatalf("PrintImage() error = %v", err)
	}
	cursorAfterImage := engine.State().CursorY

	engine.PrintLine("Footer text")

	if engine.State().CursorY <= cursorAfterImage {
		t.Error("Text should be below image")
	}
}

func TestPrintImage_WidthCappedToPaper_80mm(t *testing.T) {
	engine, err := emulator.NewDefaultEngine() // 576px
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	// Image wider than paper
	img := createTestImage(800, 100)

	opts := emulator.DefaultImageOptions()
	opts.PixelWidth = 800 // Request more than paper width

	err = engine.PrintImageWithOptions(img, opts)
	if err != nil {
		t.Fatalf("PrintImageWithOptions() error = %v", err)
	}

	result := engine.Render()
	if result.Bounds().Dx() != constants.PaperPxWidth80mm {
		t.Errorf("Output width = %d, want %d", result.Bounds().Dx(), constants.PaperPxWidth80mm)
	}
}

func TestPrintImage_WidthCappedToPaper_58mm(t *testing.T) {
	engine, err := emulator.New58mmEngine() // 384px
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	img := createTestImage(500, 100)

	opts := emulator.DefaultImageOptions()
	opts.PixelWidth = 500

	err = engine.PrintImageWithOptions(img, opts)
	if err != nil {
		t.Fatalf("PrintImageWithOptions() error = %v", err)
	}

	result := engine.Render()
	if result.Bounds().Dx() != constants.PaperPxWidth58mm {
		t.Errorf("Output width = %d, want %d", result.Bounds().Dx(), constants.PaperPxWidth58mm)
	}
}

func TestPrintImage_TransparentImage(t *testing.T) {
	engine, err := emulator.NewDefaultEngine()
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	// Create transparent image
	img := createTransparentTestImage(100, 50)

	err = engine.PrintImage(img)
	if err != nil {
		t.Fatalf("PrintImage() with transparent image error = %v", err)
	}

	// Should render without error
	result := engine.Render()
	if result == nil {
		t.Error("Render() should not return nil")
	}
}

func TestPrintImage_VerySmallImage(t *testing.T) {
	engine, err := emulator.NewDefaultEngine()
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	// 1x1 pixel image
	img := createTestImage(1, 1)

	err = engine.PrintImage(img)
	if err != nil {
		t.Fatalf("PrintImage() with 1x1 image error = %v", err)
	}
}

func TestPrintImage_CanvasGrows(t *testing.T) {
	engine, err := emulator.NewDefaultEngine()
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	// Print text to fill canvas
	for i := 0; i < 40; i++ {
		engine.PrintLine("Line of text")
	}

	// Now add a tall image
	img := createTestImage(200, 800)

	err = engine.PrintImage(img)
	if err != nil {
		t.Fatalf("PrintImage() error = %v", err)
	}

	result := engine.Render()
	if result.Bounds().Dy() < 800 {
		t.Error("Canvas should grow to accommodate tall image")
	}
}

func TestPrintImage_Reset(t *testing.T) {
	engine, err := emulator.NewDefaultEngine()
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	img := createTestImage(100, 50)
	err = engine.PrintImage(img)
	if err != nil {
		t.Fatalf("PrintImage() error = %v", err)
	}

	engine.Reset()

	// After reset, should be able to print image again
	err = engine.PrintImage(img)
	if err != nil {
		t.Fatalf("PrintImage() after Reset() error = %v", err)
	}
}

// ============================================================================
// Image Test Helpers
// ============================================================================

func createTestImage(width, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	white := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, white)
		}
	}
	return img
}

func createGrayTestImage(width, height int, gray uint8) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	c := color.RGBA{R: gray, G: gray, B: gray, A: 255}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, c)
		}
	}
	return img
}

func createColorTestImage(width, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Create rainbow pattern
			r := uint8((x * 255) / width)  //nolint:gosec
			g := uint8((y * 255) / height) //nolint:gosec
			b := uint8(128)
			img.Set(x, y, color.RGBA{R: r, G: g, B: b, A: 255})
		}
	}
	return img
}

func createTransparentTestImage(width, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	// Leave pixels as default (transparent black)
	return img
}

// ============================================================================
// Image Printing Benchmark Tests
// ============================================================================

func BenchmarkPrintImage_Small(b *testing.B) {
	engine, _ := emulator.NewDefaultEngine()
	img := createTestImage(100, 100)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		engine.Reset()
		_ = engine.PrintImage(img)
	}
}

func BenchmarkPrintImage_Medium(b *testing.B) {
	engine, _ := emulator.NewDefaultEngine()
	img := createTestImage(384, 300)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		engine.Reset()
		_ = engine.PrintImage(img)
	}
}

func BenchmarkPrintImage_ThermalPreview(b *testing.B) {
	engine, _ := emulator.NewDefaultEngine()
	img := createTestImage(384, 300)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		engine.Reset()
		_ = engine.PrintImageThermalPreview(img, 384)
	}
}

// ============================================================================
// Benchmark Tests
// ============================================================================

func BenchmarkPrintLine(b *testing.B) {
	engine, _ := emulator.NewDefaultEngine()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		engine.PrintLine("Benchmark test line")
	}
}

func BenchmarkRender(b *testing.B) {
	engine, _ := emulator.NewDefaultEngine()
	for i := 0; i < 20; i++ {
		engine.PrintLine("Line of text for benchmarking")
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = engine.Render()
	}
}

func BenchmarkWritePNG(b *testing.B) {
	engine, _ := emulator.NewDefaultEngine()
	for i := 0; i < 20; i++ {
		engine.PrintLine("Line of text for benchmarking")
	}
	var buf bytes.Buffer
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		buf.Reset()
		_ = engine.WritePNG(&buf)
	}
}
