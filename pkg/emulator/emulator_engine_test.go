// engine_test.go - Tests expandidos y estandarizados
package emulator_test

import (
	"bytes"
	"image/png"
	"testing"

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
	if state.PaperWidthPx != emulator.PaperWidth80mm {
		t.Errorf("State().PaperWidthPx = %d, want %d", state.PaperWidthPx, emulator.PaperWidth80mm)
	}
}

func TestNewEngine_58mm(t *testing.T) {
	engine, err := emulator.New58mmEngine()
	if err != nil {
		t.Fatalf("New58mmEngine() error = %v", err)
	}

	state := engine.State()
	if state.PaperWidthPx != emulator.PaperWidth58mm {
		t.Errorf("State().PaperWidthPx = %d, want %d", state.PaperWidthPx, emulator.PaperWidth58mm)
	}
}

func TestNewEngine_InvalidConfig(t *testing.T) {
	_, err := emulator.NewEngine(emulator.Config{PaperWidthPx: 0})
	if err == nil {
		t.Error("NewEngine() with zero width should return error")
	}

	_, err = emulator.NewEngine(emulator.Config{PaperWidthPx: -100})
	if err == nil {
		t.Error("NewEngine() with negative width should return error")
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

// ============================================================================
// Alignment Tests
// ============================================================================

func TestAlignment_AllModes(t *testing.T) {
	tests := []struct {
		name   string
		set    func(*emulator.Engine)
		expect string
	}{
		{"Left", func(e *emulator.Engine) { e.AlignLeft() }, emulator.AlignLeft},
		{"Center", func(e *emulator.Engine) { e.AlignCenter() }, emulator.AlignCenter},
		{"Right", func(e *emulator.Engine) { e.AlignRight() }, emulator.AlignRight},
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
		{"CENTER", emulator.AlignCenter},
		{"  center  ", emulator.AlignCenter},
		{"Right", emulator.AlignRight},
		{"invalid", emulator.AlignLeft}, // Fallback
		{"", emulator.AlignLeft},
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
// Style Tests
// ============================================================================

func TestSetSize_Clamping(t *testing.T) {
	tests := []struct {
		inputW, inputH   int
		expectW, expectH float64
	}{
		{1, 1, 1.0, 1.0},
		{8, 8, 8.0, 8.0},
		{0, 0, 1.0, 1.0},   // Clamped up
		{10, 10, 8.0, 8.0}, // Clamped down
		{-1, 5, 1.0, 5.0},
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
		{5, 2},
	}

	for _, tt := range tests {
		engine, _ := emulator.NewDefaultEngine()
		engine.SetUnderline(tt.input)
		if engine.State().IsUnderline != tt.expect {
			t.Errorf("SetUnderline(%d) = %d, want %d", tt.input, engine.State().IsUnderline, tt.expect)
		}
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

	expectedAdvance := float64(5) * float64(emulator.DefaultLineSpacing)
	actualAdvance := afterFeed - initialY

	// Allow some tolerance for line height calculations
	if actualAdvance < expectedAdvance*0.8 {
		t.Errorf("Feed(5) advanced %.1f, expected at least %.1f", actualAdvance, expectedAdvance*0.8)
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
	if state.Align != emulator.AlignLeft {
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
	if bounds.Dx() != emulator.PaperWidth80mm {
		t.Errorf("PNG width = %d, want %d", bounds.Dx(), emulator.PaperWidth80mm)
	}
}

func TestRenderWithInfo_Dimensions(t *testing.T) {
	engine, _ := emulator.NewDefaultEngine()
	engine.PrintLine("Line 1")
	engine.PrintLine("Line 2")

	result := engine.RenderWithInfo()

	if result.Width != emulator.PaperWidth80mm {
		t.Errorf("RenderWithInfo().Width = %d, want %d", result.Width, emulator.PaperWidth80mm)
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
	expected := emulator.PaperWidth80mm / emulator.FontAWidth
	if chars != expected {
		t.Errorf("CharsPerLine() with Font A = %d, want %d", chars, expected)
	}
}

func TestCharsPerLine_WithScaling(t *testing.T) {
	engine, _ := emulator.NewDefaultEngine()
	engine.SetSize(2, 2)

	chars := engine.CharsPerLine()

	// Should be half of normal
	expected := emulator.PaperWidth80mm / (emulator.FontAWidth * 2)
	if chars != expected {
		t.Errorf("CharsPerLine() with 2x scale = %d, want %d", chars, expected)
	}
}
