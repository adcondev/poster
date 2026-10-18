package constants_test

import (
	"testing"

	"github.com/adcondev/poster/pkg/constants"
)

func TestConstants_StringMethods(t *testing.T) {
	// Alignment
	alignments := []struct {
		val  constants.Alignment
		want string
	}{
		{constants.Left, "left"},
		{constants.Center, "center"},
		{constants.Right, "right"},
	}
	for _, tt := range alignments {
		if got := tt.val.String(); got != tt.want {
			t.Errorf("Alignment.String() = %q, want %q", got, tt.want)
		}
	}

	// Dithering
	ditherings := []struct {
		val  constants.Dithering
		want string
	}{
		{constants.Threshold, "threshold"},
		{constants.Atkinson, "atkinson"},
		{constants.FloydSteinberg, "floyd_steinberg"},
		{constants.Ordered, "ordered"},
	}
	for _, tt := range ditherings {
		if got := tt.val.String(); got != tt.want {
			t.Errorf("Dithering.String() = %q, want %q", got, tt.want)
		}
	}

	// Scaling
	scalings := []struct {
		val  constants.Scaling
		want string
	}{
		{constants.NearestNeighbor, "nearest_neighbor"},
		{constants.Bilinear, "bilinear"},
		{constants.Bicubic, "bicubic"},
	}
	for _, tt := range scalings {
		if got := tt.val.String(); got != tt.want {
			t.Errorf("Scaling.String() = %q, want %q", got, tt.want)
		}
	}

	// Size
	sizes := []struct {
		val  constants.Size
		want string
	}{
		{constants.Normal, "1x1"},
		{constants.Double, "2x2"},
		{constants.Triple, "3x3"},
		{constants.Quad, "4x4"},
		{constants.Penta, "5x5"},
		{constants.Hexa, "6x6"},
		{constants.Hepta, "7x7"},
		{constants.Octa, "8x8"},
	}
	for _, tt := range sizes {
		if got := tt.val.String(); got != tt.want {
			t.Errorf("Size.String() = %q, want %q", got, tt.want)
		}
	}

	// Underline
	underlines := []struct {
		val  constants.Underline
		want string
	}{
		{constants.NoDot, "0pt"},
		{constants.OneDot, "1pt"},
		{constants.TwoDot, "2pt"},
	}
	for _, tt := range underlines {
		if got := tt.val.String(); got != tt.want {
			t.Errorf("Underline.String() = %q, want %q", got, tt.want)
		}
	}

	// Font
	fonts := []struct {
		val  constants.Font
		want string
	}{
		{constants.A, "a"},
		{constants.B, "b"},
		{constants.C, "c"},
		{constants.D, "d"},
	}
	for _, tt := range fonts {
		if got := tt.val.String(); got != tt.want {
			t.Errorf("Font.String() = %q, want %q", got, tt.want)
		}
	}

	// Mode (Cut)
	modes := []struct {
		val  constants.Mode
		want string
	}{
		{constants.Partial, "partial"},
		{constants.Full, "full"},
	}
	for _, tt := range modes {
		if got := tt.val.String(); got != tt.want {
			t.Errorf("Mode.String() = %q, want %q", got, tt.want)
		}
	}

	// ErrorCorrection
	errorCorrections := []struct {
		val  constants.ErrorCorrection
		want string
	}{
		{constants.L, "L"},
		{constants.M, "M"},
		{constants.Q, "Q"},
		{constants.H, "H"},
	}
	for _, tt := range errorCorrections {
		if got := tt.val.String(); got != tt.want {
			t.Errorf("ErrorCorrection.String() = %q, want %q", got, tt.want)
		}
	}

	// Symbology
	symbologies := []struct {
		val  constants.Symbology
		want string
	}{
		{constants.CODE39, "code39"},
		{constants.CODE128, "code128"},
		{constants.EAN13, "ean13"},
		{constants.UPCA, "upca"},
		{constants.ITF, "itf"},
		{constants.CODABAR, "codabar"},
		{constants.EAN8, "ean8"},
		{constants.UPCE, "upce"},
		{constants.JAN13, "jan13"},
		{constants.JAN8, "jan8"},
		{constants.CODE93, "code93"},
		{constants.GS1128, "gs1_128"},
		{constants.CODE128Auto, "code128_auto"},
	}
	for _, tt := range symbologies {
		if got := tt.val.String(); got != tt.want {
			t.Errorf("Symbology.String() = %q, want %q", got, tt.want)
		}
	}

	// HriPosition
	hriPositions := []struct {
		val  constants.HriPosition
		want string
	}{
		{constants.Below, "below"},
		{constants.Above, "above"},
		{constants.Both, "both"},
		{constants.None, "none"},
	}
	for _, tt := range hriPositions {
		if got := tt.val.String(); got != tt.want {
			t.Errorf("HriPosition.String() = %q, want %q", got, tt.want)
		}
	}

	// Format
	formats := []struct {
		val  constants.Format
		want string
	}{
		{constants.Hex, "hex"},
		{constants.Base64, "base64"},
	}
	for _, tt := range formats {
		if got := tt.val.String(); got != tt.want {
			t.Errorf("Format.String() = %q, want %q", got, tt.want)
		}
	}
}

func TestConversions(t *testing.T) {
	// PaperWidthFromMm
	if got := constants.PaperWidthFromMm(80); got != 80*constants.DefaultDotsPerMm {
		t.Errorf("PaperWidthFromMm(80) = %d, want %d", got, 80*constants.DefaultDotsPerMm)
	}

	// CharsPerLineForFont
	if got := constants.CharsPerLineForFont(576, 12); got != 48 {
		t.Errorf("CharsPerLineForFont(576, 12) = %d, want 48", got)
	}
	if got := constants.CharsPerLineForFont(576, 0); got != 0 {
		t.Errorf("CharsPerLineForFont(576, 0) = %d, want 0", got)
	}
	if got := constants.CharsPerLineForFont(576, -5); got != 0 {
		t.Errorf("CharsPerLineForFont(576, -5) = %d, want 0", got)
	}

	// MaxCharsForPaper
	if got := constants.MaxCharsForPaper(384, 12); got != 32 {
		t.Errorf("MaxCharsForPaper(384, 12) = %d, want 32", got)
	}
	if got := constants.MaxCharsForPaper(384, 0); got != 0 {
		t.Errorf("MaxCharsForPaper(384, 0) = %d, want 0", got)
	}
	if got := constants.MaxCharsForPaper(384, -2); got != 0 {
		t.Errorf("MaxCharsForPaper(384, -2) = %d, want 0", got)
	}

	// MaxCharsForPaperFontA
	if got := constants.MaxCharsForPaperFontA(384); got != 384/constants.FontAWidth {
		t.Errorf("MaxCharsForPaperFontA(384) = %d, want %d", got, 384/constants.FontAWidth)
	}

	// MaxCharsForPaperFontB
	if got := constants.MaxCharsForPaperFontB(384); got != 384/constants.FontBWidth {
		t.Errorf("MaxCharsForPaperFontB(384) = %d, want %d", got, 384/constants.FontBWidth)
	}
}
