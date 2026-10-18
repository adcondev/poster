package emulator_test

import (
	"bytes"
	"image/png"
	"testing"

	"github.com/adcondev/poster/pkg/commands/mechanismcontrol"
	"github.com/adcondev/poster/pkg/composer"
	"github.com/adcondev/poster/pkg/constants"
	"github.com/adcondev/poster/pkg/emulator"
)

// ============================================================================
// Text and Newline Parsing Tests
// ============================================================================

func TestParseESC_EmptyAndNonPrintable(t *testing.T) {
	eng, err := emulator.NewDefaultEngine()
	if err != nil {
		t.Fatalf("NewDefaultEngine() error = %v", err)
	}

	if err := eng.ParseESC(nil); err != nil {
		t.Errorf("ParseESC(nil) error = %v", err)
	}

	if err := eng.ParseESC([]byte{}); err != nil {
		t.Errorf("ParseESC([]byte{}) error = %v", err)
	}

	// Non-printable control characters should be safely skipped
	nonPrintable := []byte{0x00, 0x01, 0x04, 0x07, 0x11}
	if err := eng.ParseESC(nonPrintable); err != nil {
		t.Errorf("ParseESC(nonPrintable) error = %v", err)
	}
}

func TestParseESC_TextAndNewlines(t *testing.T) {
	eng, err := emulator.NewDefaultEngine()
	if err != nil {
		t.Fatalf("NewDefaultEngine() error = %v", err)
	}

	initialY := eng.State().CursorY

	// LF
	if err := eng.ParseESC([]byte("Line 1\nLine 2\n")); err != nil {
		t.Fatalf("ParseESC() error = %v", err)
	}
	yAfterLF := eng.State().CursorY
	if yAfterLF <= initialY {
		t.Errorf("CursorY after LF = %v, expected > %v", yAfterLF, initialY)
	}

	// CR LF
	eng.Reset()
	if err := eng.ParseESC([]byte("Line A\r\nLine B\r\n")); err != nil {
		t.Fatalf("ParseESC() error = %v", err)
	}
	yAfterCRLF := eng.State().CursorY
	if yAfterCRLF <= initialY {
		t.Errorf("CursorY after CRLF = %v, expected > %v", yAfterCRLF, initialY)
	}

	// CR alone
	eng.Reset()
	if err := eng.ParseESC([]byte("Line X\rLine Y\r")); err != nil {
		t.Fatalf("ParseESC() error = %v", err)
	}
	yAfterCR := eng.State().CursorY
	if yAfterCR <= initialY {
		t.Errorf("CursorY after CR = %v, expected > %v", yAfterCR, initialY)
	}
}

// ============================================================================
// ESC Command Tests
// ============================================================================

func TestParseESC_Initialize(t *testing.T) {
	eng, err := emulator.NewDefaultEngine()
	if err != nil {
		t.Fatalf("NewDefaultEngine() error = %v", err)
	}

	// Set bold, inverse, double size
	eng.SetBold(true)
	eng.SetInverse(true)
	eng.SetSize(2, 2)

	// Send ESC @ (0x1B 0x40)
	if err := eng.ParseESC([]byte{0x1B, '@'}); err != nil {
		t.Fatalf("ParseESC(ESC @) error = %v", err)
	}

	state := eng.State()
	if state.IsBold {
		t.Error("IsBold = true, expected false after reset")
	}
	if state.IsInverse {
		t.Error("IsInverse = true, expected false after reset")
	}
	if state.ScaleW != 1.0 || state.ScaleH != 1.0 {
		t.Errorf("Scale = (%v, %v), expected (1.0, 1.0) after reset", state.ScaleW, state.ScaleH)
	}
}

func TestParseESC_Bold(t *testing.T) {
	eng, err := emulator.NewDefaultEngine()
	if err != nil {
		t.Fatalf("NewDefaultEngine() error = %v", err)
	}

	// ESC E 1 -> bold on
	if err := eng.ParseESC([]byte{0x1B, 'E', 0x01}); err != nil {
		t.Fatalf("ParseESC(ESC E 1) error = %v", err)
	}
	if !eng.State().IsBold {
		t.Error("IsBold = false, want true")
	}

	// ESC E 0 -> bold off
	if err := eng.ParseESC([]byte{0x1B, 'E', 0x00}); err != nil {
		t.Fatalf("ParseESC(ESC E 0) error = %v", err)
	}
	if eng.State().IsBold {
		t.Error("IsBold = true, want false")
	}
}

func TestParseESC_Underline(t *testing.T) {
	eng, err := emulator.NewDefaultEngine()
	if err != nil {
		t.Fatalf("NewDefaultEngine() error = %v", err)
	}

	tests := []struct {
		input []byte
		want  int
	}{
		{[]byte{0x1B, '-', 0x01}, 1},
		{[]byte{0x1B, '-', 0x02}, 2},
		{[]byte{0x1B, '-', 0x00}, 0},
		{[]byte{0x1B, '-', '1'}, 1},
		{[]byte{0x1B, '-', '2'}, 2},
		{[]byte{0x1B, '-', '0'}, 0},
	}

	for _, tt := range tests {
		if err := eng.ParseESC(tt.input); err != nil {
			t.Fatalf("ParseESC(%v) error = %v", tt.input, err)
		}
		if got := eng.State().IsUnderline; got != tt.want {
			t.Errorf("ParseESC(%v) IsUnderline = %d, want %d", tt.input, got, tt.want)
		}
	}
}

func TestParseESC_Font(t *testing.T) {
	eng, err := emulator.NewDefaultEngine()
	if err != nil {
		t.Fatalf("NewDefaultEngine() error = %v", err)
	}

	tests := []struct {
		input []byte
		want  string
	}{
		{[]byte{0x1B, 'M', 0x01}, "B"},
		{[]byte{0x1B, 'M', 0x00}, "A"},
		{[]byte{0x1B, 'M', '1'}, "B"},
		{[]byte{0x1B, 'M', '0'}, "A"},
	}

	for _, tt := range tests {
		if err := eng.ParseESC(tt.input); err != nil {
			t.Fatalf("ParseESC(%v) error = %v", tt.input, err)
		}
		if got := eng.State().FontName; got != tt.want {
			t.Errorf("ParseESC(%v) FontName = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestParseESC_Alignment(t *testing.T) {
	eng, err := emulator.NewDefaultEngine()
	if err != nil {
		t.Fatalf("NewDefaultEngine() error = %v", err)
	}

	tests := []struct {
		input []byte
		want  string
	}{
		{[]byte{0x1B, 'a', 0x00}, "left"},
		{[]byte{0x1B, 'a', 0x01}, "center"},
		{[]byte{0x1B, 'a', 0x02}, "right"},
		{[]byte{0x1B, 'a', '0'}, "left"},
		{[]byte{0x1B, 'a', '1'}, "center"},
		{[]byte{0x1B, 'a', '2'}, "right"},
		{[]byte{0x1B, 'a', 99}, "left"}, // fallback
	}

	for _, tt := range tests {
		if err := eng.ParseESC(tt.input); err != nil {
			t.Fatalf("ParseESC(%v) error = %v", tt.input, err)
		}
		if got := eng.State().Align; got != tt.want {
			t.Errorf("ParseESC(%v) Align = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestParseESC_FeedLines(t *testing.T) {
	eng, err := emulator.NewDefaultEngine()
	if err != nil {
		t.Fatalf("NewDefaultEngine() error = %v", err)
	}

	initialY := eng.State().CursorY
	if err := eng.ParseESC([]byte{0x1B, 'd', 0x03}); err != nil {
		t.Fatalf("ParseESC(ESC d 3) error = %v", err)
	}

	if got := eng.State().CursorY; got <= initialY {
		t.Errorf("CursorY = %v, expected advance after ESC d 3", got)
	}
}

func TestParseESC_LineSpacing(t *testing.T) {
	eng, err := emulator.NewDefaultEngine()
	if err != nil {
		t.Fatalf("NewDefaultEngine() error = %v", err)
	}

	// ESC 3 n: set line spacing
	if err := eng.ParseESC([]byte{0x1B, '3', 50}); err != nil {
		t.Fatalf("ParseESC(ESC 3 50) error = %v", err)
	}
	if got := eng.State().LineSpacing; got != 50.0 {
		t.Errorf("LineSpacing = %v, want 50.0", got)
	}

	// ESC 2: select default line spacing
	if err := eng.ParseESC([]byte{0x1B, '2'}); err != nil {
		t.Fatalf("ParseESC(ESC 2) error = %v", err)
	}
	if got := eng.State().LineSpacing; got != float64(constants.DefaultLineSpacing) {
		t.Errorf("LineSpacing = %v, want %v", got, float64(constants.DefaultLineSpacing))
	}
}

func TestParseESC_PrintMode(t *testing.T) {
	eng, err := emulator.NewDefaultEngine()
	if err != nil {
		t.Fatalf("NewDefaultEngine() error = %v", err)
	}

	// ESC ! with font B (0x01) | bold (0x08) | double-height (0x10) | double-width (0x20) | underline (0x80)
	mode := byte(0x01 | 0x08 | 0x10 | 0x20 | 0x80)
	if err := eng.ParseESC([]byte{0x1B, '!', mode}); err != nil {
		t.Fatalf("ParseESC(ESC ! mode) error = %v", err)
	}

	state := eng.State()
	if state.FontName != "B" {
		t.Errorf("FontName = %q, want \"B\"", state.FontName)
	}
	if !state.IsBold {
		t.Error("IsBold = false, want true")
	}
	if state.ScaleW != 2.0 || state.ScaleH != 2.0 {
		t.Errorf("Scale = (%v, %v), want (2.0, 2.0)", state.ScaleW, state.ScaleH)
	}
	if state.IsUnderline != 1 {
		t.Errorf("IsUnderline = %d, want 1", state.IsUnderline)
	}

	// Clear print mode
	if err := eng.ParseESC([]byte{0x1B, '!', 0x00}); err != nil {
		t.Fatalf("ParseESC(ESC ! 0) error = %v", err)
	}
	state = eng.State()
	if state.FontName != "A" {
		t.Errorf("FontName = %q, want \"A\"", state.FontName)
	}
	if state.IsBold {
		t.Error("IsBold = true, want false")
	}
	if state.ScaleW != 1.0 || state.ScaleH != 1.0 {
		t.Errorf("Scale = (%v, %v), want (1.0, 1.0)", state.ScaleW, state.ScaleH)
	}
	if state.IsUnderline != 0 {
		t.Errorf("IsUnderline = %d, want 0", state.IsUnderline)
	}
}

func TestParseESC_SkippedCommands(t *testing.T) {
	eng, err := emulator.NewDefaultEngine()
	if err != nil {
		t.Fatalf("NewDefaultEngine() error = %v", err)
	}

	// ESC t (code table), ESC p (pulse), ESC B (beep), and unknown ESC
	raw := []byte{
		0x1B, 't', 0x00,
		0x1B, 'p', 0x00, 0x19, 0x32,
		0x1B, 'B', 0x02, 0x03,
		0x1B, 0xFF,
	}

	if err := eng.ParseESC(raw); err != nil {
		t.Fatalf("ParseESC(skipped commands) error = %v", err)
	}
}

// ============================================================================
// GS Command Tests
// ============================================================================

func TestParseESC_GSCharacterSize(t *testing.T) {
	eng, err := emulator.NewDefaultEngine()
	if err != nil {
		t.Fatalf("NewDefaultEngine() error = %v", err)
	}

	tests := []struct {
		n     byte
		wantW float64
		wantH float64
	}{
		{0x00, 1.0, 1.0},
		{0x11, 2.0, 2.0},
		{0x77, 8.0, 8.0},
		{0x31, 4.0, 2.0}, // width multiplier 3+1=4, height multiplier 1+1=2
	}

	for _, tt := range tests {
		if err := eng.ParseESC([]byte{0x1D, '!', tt.n}); err != nil {
			t.Fatalf("ParseESC(GS ! %d) error = %v", tt.n, err)
		}
		state := eng.State()
		if state.ScaleW != tt.wantW || state.ScaleH != tt.wantH {
			t.Errorf("GS ! %d Scale = (%v, %v), want (%v, %v)",
				tt.n, state.ScaleW, state.ScaleH, tt.wantW, tt.wantH)
		}
	}
}

func TestParseESC_GSInverse(t *testing.T) {
	eng, err := emulator.NewDefaultEngine()
	if err != nil {
		t.Fatalf("NewDefaultEngine() error = %v", err)
	}

	// GS B 1: inverse on
	if err := eng.ParseESC([]byte{0x1D, 'B', 0x01}); err != nil {
		t.Fatalf("ParseESC(GS B 1) error = %v", err)
	}
	if !eng.State().IsInverse {
		t.Error("IsInverse = false, want true")
	}

	// GS B 0: inverse off
	if err := eng.ParseESC([]byte{0x1D, 'B', 0x00}); err != nil {
		t.Fatalf("ParseESC(GS B 0) error = %v", err)
	}
	if eng.State().IsInverse {
		t.Error("IsInverse = true, want false")
	}
}

func TestParseESC_GSCut(t *testing.T) {
	eng, err := emulator.NewDefaultEngine()
	if err != nil {
		t.Fatalf("NewDefaultEngine() error = %v", err)
	}

	// GS V 0 (full cut)
	if err := eng.ParseESC([]byte{0x1D, 'V', 0x00}); err != nil {
		t.Fatalf("ParseESC(GS V 0) error = %v", err)
	}

	// GS V 1 (partial cut)
	if err := eng.ParseESC([]byte{0x1D, 'V', 0x01}); err != nil {
		t.Fatalf("ParseESC(GS V 1) error = %v", err)
	}

	// GS V 65 n (feed n lines and cut)
	if err := eng.ParseESC([]byte{0x1D, 'V', 65, 3}); err != nil {
		t.Fatalf("ParseESC(GS V 65 3) error = %v", err)
	}

	// GS V 66 n (feed n lines and partial cut)
	if err := eng.ParseESC([]byte{0x1D, 'V', 66, 2}); err != nil {
		t.Fatalf("ParseESC(GS V 66 2) error = %v", err)
	}
}

func TestParseESC_GSRasterBitImage(t *testing.T) {
	eng, err := emulator.NewDefaultEngine()
	if err != nil {
		t.Fatalf("NewDefaultEngine() error = %v", err)
	}

	initialY := eng.State().CursorY

	// Construct a 16x8 monochrome raster image:
	// xL=2, xH=0 (2 bytes = 16 dots per row), yL=8, yH=0 (8 dots high)
	// total data bytes = 2 * 8 = 16 bytes
	xBytes := 2
	yDots := 8
	cmd := make([]byte, 0, 8+xBytes*yDots)
	cmd = append(cmd, 0x1D, 'v', '0', 0x00, byte(xBytes), 0x00, byte(yDots), 0x00)
	for i := 0; i < xBytes*yDots; i++ {
		cmd = append(cmd, 0xAA)
	}

	if err := eng.ParseESC(cmd); err != nil {
		t.Fatalf("ParseESC(GS v 0) error = %v", err)
	}

	if got := eng.State().CursorY; got <= initialY {
		t.Errorf("CursorY = %v, expected advance after raster image", got)
	}

	// Zero dimension raster bit image should be safely ignored
	zeroCmd := []byte{0x1D, 'v', '0', 0x00, 0x00, 0x00, 0x00, 0x00}
	if err := eng.ParseESC(zeroCmd); err != nil {
		t.Fatalf("ParseESC(zero raster) error = %v", err)
	}
}

func TestParseESC_GSBarcodeAnd2DSkip(t *testing.T) {
	eng, err := emulator.NewDefaultEngine()
	if err != nil {
		t.Fatalf("NewDefaultEngine() error = %v", err)
	}

	raw := []byte{
		// Barcode formatting parameters
		0x1D, 'h', 80, // height
		0x1D, 'w', 3, // width
		0x1D, 'H', 2, // HRI position
		0x1D, 'f', 0, // HRI font
		// Barcode type m <= 6 (null terminated)
		0x1D, 'k', 2, '1', '2', '3', '4', 0x00,
		// Barcode type m > 6 (explicit length prefix)
		0x1D, 'k', 67, 4, '9', '8', '7', '6',
		// 2D function GS ( k pL pH ...
		0x1D, '(', 'k', 0x04, 0x00, 0x31, 0x50, 0x30, 0x00,
		// Unknown GS command
		0x1D, 0xFE,
	}

	if err := eng.ParseESC(raw); err != nil {
		t.Fatalf("ParseESC(barcode and 2D skip) error = %v", err)
	}
}

// ============================================================================
// Truncated and Malformed Stream Tests
// ============================================================================

func TestParseESC_TruncatedInputs(t *testing.T) {
	eng, err := emulator.NewDefaultEngine()
	if err != nil {
		t.Fatalf("NewDefaultEngine() error = %v", err)
	}

	cases := map[string][]byte{
		"bare ESC":                           {0x1B},
		"bare GS":                            {0x1D},
		"ESC E truncated":                    {0x1B, 'E'},
		"ESC - truncated":                    {0x1B, '-'},
		"ESC M truncated":                    {0x1B, 'M'},
		"ESC a truncated":                    {0x1B, 'a'},
		"ESC d truncated":                    {0x1B, 'd'},
		"ESC 3 truncated":                    {0x1B, '3'},
		"ESC ! truncated":                    {0x1B, '!'},
		"ESC t truncated":                    {0x1B, 't'},
		"ESC p truncated 1":                  {0x1B, 'p'},
		"ESC p truncated 2":                  {0x1B, 'p', 0x00, 0x10},
		"ESC B truncated 1":                  {0x1B, 'B'},
		"ESC B truncated 2":                  {0x1B, 'B', 0x01},
		"GS ! truncated":                     {0x1D, '!'},
		"GS B truncated":                     {0x1D, 'B'},
		"GS V truncated 1":                   {0x1D, 'V'},
		"GS V 65 truncated":                  {0x1D, 'V', 65},
		"GS v truncated":                     {0x1D, 'v'},
		"GS v 0 truncated header":            {0x1D, 'v', '0', 0x00, 0x02},
		"GS v 0 truncated data":              {0x1D, 'v', '0', 0x00, 0x02, 0x00, 0x08, 0x00, 0xFF},
		"GS h truncated":                     {0x1D, 'h'},
		"GS k truncated":                     {0x1D, 'k'},
		"GS k type 0 missing null":           {0x1D, 'k', 0x00, '1', '2', '3'},
		"GS k type 67 truncated data":        {0x1D, 'k', 67, 10, '1'},
		"GS ( truncated header":              {0x1D, '(', 'k'},
		"GS ( truncated data":                {0x1D, '(', 'k', 0x20, 0x00, 0x01},
		"trailing ESC after text":            {'T', 'e', 'x', 't', 0x1B},
		"trailing GS after text":             {'T', 'e', 'x', 't', 0x1D},
		"trailing CR":                        {'T', 'e', 'x', 't', 0x0D},
	}

	for name, data := range cases {
		t.Run(name, func(t *testing.T) {
			if err := eng.ParseESC(data); err != nil {
				t.Fatalf("ParseESC(%s) unexpected error = %v", name, err)
			}
		})
	}
}

// ============================================================================
// RenderBytesToPNG Tests
// ============================================================================

func TestRenderBytesToPNG_ValidOutput(t *testing.T) {
	raw := []byte("Header Line\nItem 1           $10.00\nTotal            $10.00\n")
	pngBytes, err := emulator.RenderBytesToPNG(raw)
	if err != nil {
		t.Fatalf("RenderBytesToPNG() error = %v", err)
	}

	if len(pngBytes) == 0 {
		t.Fatal("RenderBytesToPNG() returned empty bytes")
	}

	img, err := png.Decode(bytes.NewReader(pngBytes))
	if err != nil {
		t.Fatalf("png.Decode() error = %v", err)
	}

	if got := img.Bounds().Dx(); got != constants.PaperPxWidth80mm {
		t.Errorf("Image width = %d, want %d", got, constants.PaperPxWidth80mm)
	}

	if got := img.Bounds().Dy(); got <= 0 {
		t.Errorf("Image height = %d, want > 0", got)
	}
}

func TestRenderBytesToPNG_EmptyCanvas(t *testing.T) {
	pngBytes, err := emulator.RenderBytesToPNG([]byte{})
	if err != nil {
		t.Fatalf("RenderBytesToPNG([]byte{}) error = %v", err)
	}

	img, err := png.Decode(bytes.NewReader(pngBytes))
	if err != nil {
		t.Fatalf("png.Decode() error = %v", err)
	}

	if got := img.Bounds().Dx(); got != constants.PaperPxWidth80mm {
		t.Errorf("Image width = %d, want %d", got, constants.PaperPxWidth80mm)
	}
}

// ============================================================================
// RenderBytesToPNGWithConfig Tests
// ============================================================================

func TestRenderBytesToPNGWithConfig_58mm(t *testing.T) {
	cfg := emulator.Config{
		PaperPxWidth: constants.PaperPxWidth58mm,
		DPI:          constants.DefaultDPI,
	}

	raw := []byte("58mm Thermal Receipt\nItem A           $5.00\n")
	pngBytes, err := emulator.RenderBytesToPNGWithConfig(raw, cfg)
	if err != nil {
		t.Fatalf("RenderBytesToPNGWithConfig() error = %v", err)
	}

	img, err := png.Decode(bytes.NewReader(pngBytes))
	if err != nil {
		t.Fatalf("png.Decode() error = %v", err)
	}

	if got := img.Bounds().Dx(); got != constants.PaperPxWidth58mm {
		t.Errorf("Image width = %d, want %d", got, constants.PaperPxWidth58mm)
	}
}

func TestRenderBytesToPNGWithConfig_InvalidConfig(t *testing.T) {
	cfg := emulator.Config{
		PaperPxWidth: -100,
	}

	_, err := emulator.RenderBytesToPNGWithConfig([]byte("Test"), cfg)
	if err == nil {
		t.Fatal("RenderBytesToPNGWithConfig() expected error for negative width, got nil")
	}
}

// ============================================================================
// Composer Integration Tests
// ============================================================================

func TestParseESC_ComposerIntegration(t *testing.T) {
	proto := composer.NewEscpos()

	var buf bytes.Buffer
	buf.Write(proto.InitializePrinter())
	buf.Write(proto.CenterAlign())

	line1, err := proto.PrintLn("SAMPLE RECEIPT")
	if err != nil {
		t.Fatalf("PrintLn error = %v", err)
	}
	buf.Write(line1)

	buf.Write(proto.LeftAlign())
	buf.Write(proto.DoubleSizeText())

	line2, err := proto.PrintLn("Large Coffee")
	if err != nil {
		t.Fatalf("PrintLn error = %v", err)
	}
	buf.Write(line2)

	buf.Write(proto.SingleSizeText())
	line3, err := proto.PrintLn("Total: $4.50")
	if err != nil {
		t.Fatalf("PrintLn error = %v", err)
	}
	buf.Write(line3)

	cutCmd, err := proto.MechanismControl.CutPaper(mechanismcontrol.CutModePartial)
	if err != nil {
		t.Fatalf("CutPaper error = %v", err)
	}
	buf.Write(cutCmd)

	pngBytes, err := emulator.RenderBytesToPNG(buf.Bytes())
	if err != nil {
		t.Fatalf("RenderBytesToPNG() error = %v", err)
	}

	img, err := png.Decode(bytes.NewReader(pngBytes))
	if err != nil {
		t.Fatalf("png.Decode() error = %v", err)
	}

	if got := img.Bounds().Dx(); got != constants.PaperPxWidth80mm {
		t.Errorf("Image width = %d, want %d", got, constants.PaperPxWidth80mm)
	}
	if got := img.Bounds().Dy(); got <= 0 {
		t.Errorf("Image height = %d, want > 0", got)
	}
}
