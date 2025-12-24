package barcode_test

import (
	"testing"

	"github.com/adcondev/poster/internal/testutils"
	"github.com/adcondev/poster/pkg/commands/barcode"
	"github.com/adcondev/poster/pkg/commands/shared"
)

var (
	// Common testutils barcodes
	Upca12     = []byte("123456789012")
	Upca11     = []byte("12345678901")
	Upce8      = []byte("12345678")
	Upce7      = []byte("1234567")
	Upce6      = []byte("123456")
	Ean1313    = []byte("1234567890123")
	Ean1312    = []byte("123456789012")
	Ean88      = []byte("12345678")
	Ean87      = []byte("1234567")
	Code39Noss = []byte("ABC-123")
	Code39Ss   = []byte("*ABC-123*")
	ITFEven    = []byte("123456")
	ITFMin     = []byte("12")
	ITFOdd     = []byte("12345")
	CodabarSs1 = []byte("A123456B")
	CodabarSs2 = []byte("C123456D")
	CodabarLss = []byte("a123456d")
)

// ============================================================================
// HRI Position Tests
// ============================================================================

func TestCommands_SelectHRICharacterPosition(t *testing.T) {
	cmd := barcode.NewCommands()
	prefix := testutils.BuildCommand(shared.GS, 'H')

	tests := []struct {
		name     string
		position barcode.HRIPosition
		want     []byte
		wantErr  error
	}{
		{
			name:     "not printed",
			position: barcode.HRINotPrinted,
			want:     append(prefix, 0x00),
			wantErr:  nil,
		},
		{
			name:     "above barcode",
			position: barcode.HRIAbove,
			want:     append(prefix, 0x01),
			wantErr:  nil,
		},
		{
			name:     "below barcode",
			position: barcode.HRIBelow,
			want:     append(prefix, 0x02),
			wantErr:  nil,
		},
		{
			name:     "both above and below",
			position: barcode.HRIBoth,
			want:     append(prefix, 0x03),
			wantErr:  nil,
		},
		{
			name:     "not printed ASCII",
			position: barcode.HRINotPrintedASCII,
			want:     append(prefix, '0'),
			wantErr:  nil,
		},
		{
			name:     "above ASCII",
			position: barcode.HRIAboveASCII,
			want:     append(prefix, '1'),
			wantErr:  nil,
		},
		{
			name:     "below ASCII",
			position: barcode.HRIBelowASCII,
			want:     append(prefix, '2'),
			wantErr:  nil,
		},
		{
			name:     "both ASCII",
			position: barcode.HRIBothASCII,
			want:     append(prefix, '3'),
			wantErr:  nil,
		},
		{
			name:     "invalid position",
			position: 99,
			want:     nil,
			wantErr:  barcode.ErrHRIPosition,
		},
		{
			name:     "invalid position 4",
			position: 4,
			want:     nil,
			wantErr:  barcode.ErrHRIPosition,
		},
		{
			name:     "invalid position 47",
			position: 47,
			want:     nil,
			wantErr:  barcode.ErrHRIPosition,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cmd.SelectHRICharacterPosition(tt.position)

			if !testutils.AssertErrorOccurred(t, err, tt.wantErr != nil, "SelectHRICharacterPosition") {
				return
			}
			if tt.wantErr != nil {
				testutils.AssertError(t, err, tt.wantErr)
				return
			}

			testutils.AssertBytes(t, got, tt.want, "SelectHRICharacterPosition(%v)", tt.position)
		})
	}
}

// ============================================================================
// HRI Font Tests
// ============================================================================

func TestCommands_SelectFontForHRI(t *testing.T) {
	cmd := barcode.NewCommands()
	prefix := testutils.BuildCommand(shared.GS, 'f')

	tests := []struct {
		name    string
		font    barcode.HRIFont
		want    []byte
		wantErr error
	}{
		{
			name:    "font A",
			font:    barcode.HRIFontA,
			want:    append(prefix, 0x00),
			wantErr: nil,
		},
		{
			name:    "font B",
			font:    barcode.HRIFontB,
			want:    append(prefix, 0x01),
			wantErr: nil,
		},
		{
			name:    "font C",
			font:    barcode.HRIFontC,
			want:    append(prefix, 0x02),
			wantErr: nil,
		},
		{
			name:    "font D",
			font:    barcode.HRIFontD,
			want:    append(prefix, 0x03),
			wantErr: nil,
		},
		{
			name:    "font E",
			font:    barcode.HRIFontE,
			want:    append(prefix, 0x04),
			wantErr: nil,
		},
		{
			name:    "font A ASCII",
			font:    barcode.HRIFontAASCII,
			want:    append(prefix, '0'),
			wantErr: nil,
		},
		{
			name:    "font B ASCII",
			font:    barcode.HRIFontBASCII,
			want:    append(prefix, '1'),
			wantErr: nil,
		},
		{
			name:    "font C ASCII",
			font:    barcode.HRIFontCASCII,
			want:    append(prefix, '2'),
			wantErr: nil,
		},
		{
			name:    "font D ASCII",
			font:    barcode.HRIFontDASCII,
			want:    append(prefix, '3'),
			wantErr: nil,
		},
		{
			name:    "font E ASCII",
			font:    barcode.HRIFontEASCII,
			want:    append(prefix, '4'),
			wantErr: nil,
		},
		{
			name:    "special font A",
			font:    barcode.HRISpecialFontA,
			want:    append(prefix, 97),
			wantErr: nil,
		},
		{
			name:    "special font B",
			font:    barcode.HRISpecialFontB,
			want:    append(prefix, 98),
			wantErr: nil,
		},
		{
			name:    "invalid font 99",
			font:    99,
			want:    nil,
			wantErr: barcode.ErrHRIFont,
		},
		{
			name:    "invalid font 5",
			font:    5,
			want:    nil,
			wantErr: barcode.ErrHRIFont,
		},
		{
			name:    "invalid font 96",
			font:    96,
			want:    nil,
			wantErr: barcode.ErrHRIFont,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cmd.SelectFontForHRI(tt.font)

			if !testutils.AssertErrorOccurred(t, err, tt.wantErr != nil, "SelectFontForHRI") {
				return
			}
			if tt.wantErr != nil {
				testutils.AssertError(t, err, tt.wantErr)
				return
			}

			testutils.AssertBytes(t, got, tt.want, "SelectFontForHRI(%v)", tt.font)
		})
	}
}

// ============================================================================
// Barcode Dimensions Tests
// ============================================================================

func TestCommands_SetBarcodeHeight(t *testing.T) {
	cmd := barcode.NewCommands()
	prefix := testutils.BuildCommand(shared.GS, 'h')

	tests := []struct {
		name    string
		height  barcode.Height
		want    []byte
		wantErr error
	}{
		{
			name:    "minimum height",
			height:  barcode.MinHeight,
			want:    append(prefix, 1),
			wantErr: nil,
		},
		{
			name:    "default height",
			height:  barcode.DefaultHeight,
			want:    append(prefix, 162),
			wantErr: nil,
		},
		{
			name:    "maximum height",
			height:  barcode.MaxHeight,
			want:    append(prefix, 255),
			wantErr: nil,
		},
		{
			name:    "typical height",
			height:  100,
			want:    append(prefix, 100),
			wantErr: nil,
		},
		{
			name:    "zero height (invalid)",
			height:  0,
			want:    nil,
			wantErr: barcode.ErrHeight,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cmd.SetBarcodeHeight(tt.height)

			if !testutils.AssertErrorOccurred(t, err, tt.wantErr != nil, "SetBarcodeHeight") {
				return
			}
			if tt.wantErr != nil {
				testutils.AssertError(t, err, tt.wantErr)
				return
			}

			testutils.AssertBytes(t, got, tt.want, "SetBarcodeHeight(%v)", tt.height)
		})
	}
}

func TestCommands_SetBarcodeWidth(t *testing.T) {
	cmd := barcode.NewCommands()
	prefix := testutils.BuildCommand(shared.GS, 'w')

	tests := []struct {
		name    string
		width   barcode.Width
		want    []byte
		wantErr error
	}{
		{
			name:    "minimum width",
			width:   barcode.MinWidth,
			want:    append(prefix, 2),
			wantErr: nil,
		},
		{
			name:    "default width",
			width:   barcode.DefaultWidth,
			want:    append(prefix, 3),
			wantErr: nil,
		},
		{
			name:    "width 4",
			width:   4,
			want:    append(prefix, 4),
			wantErr: nil,
		},
		{
			name:    "width 5",
			width:   5,
			want:    append(prefix, 5),
			wantErr: nil,
		},
		{
			name:    "maximum width",
			width:   barcode.MaxWidth,
			want:    append(prefix, 6),
			wantErr: nil,
		},
		{
			name:    "extended minimum width",
			width:   barcode.ExtendedMinWidth,
			want:    append(prefix, 68),
			wantErr: nil,
		},
		{
			name:    "extended width 70",
			width:   70,
			want:    append(prefix, 70),
			wantErr: nil,
		},
		{
			name:    "extended width 72",
			width:   72,
			want:    append(prefix, 72),
			wantErr: nil,
		},
		{
			name:    "extended maximum width",
			width:   barcode.ExtendedMaxWidth,
			want:    append(prefix, 76),
			wantErr: nil,
		},
		{
			name:    "invalid width (gap) 7",
			width:   7,
			want:    nil,
			wantErr: barcode.ErrWidth,
		},
		{
			name:    "invalid width (gap) 67",
			width:   67,
			want:    nil,
			wantErr: barcode.ErrWidth,
		},
		{
			name:    "invalid width (too low)",
			width:   1,
			want:    nil,
			wantErr: barcode.ErrWidth,
		},
		{
			name:    "invalid width (too high)",
			width:   77,
			want:    nil,
			wantErr: barcode.ErrWidth,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cmd.SetBarcodeWidth(tt.width)

			if !testutils.AssertErrorOccurred(t, err, tt.wantErr != nil, "SetBarcodeWidth") {
				return
			}
			if tt.wantErr != nil {
				testutils.AssertError(t, err, tt.wantErr)
				return
			}

			testutils.AssertBytes(t, got, tt.want, "SetBarcodeWidth(%v)", tt.width)
		})
	}
}

// ============================================================================
// Barcode Printing Tests - Function A
// ============================================================================

func TestCommands_PrintBarcode_FunctionA(t *testing.T) {
	cmd := barcode.NewCommands()

	tests := []struct {
		name      string
		symbology barcode.Symbology
		data      []byte
		want      []byte
		wantErr   error
	}{
		{
			name:      "UPC-A 11 digits",
			symbology: barcode.UPCA,
			data:      Upca11,
			want:      append([]byte{shared.GS, 'k', 0}, append([]byte("12345678901"), shared.NUL)...),
			wantErr:   nil,
		},
		{
			name:      "UPC-A 12 digits with check",
			symbology: barcode.UPCA,
			data:      Upca12,
			want:      append([]byte{shared.GS, 'k', 0}, append([]byte("123456789012"), shared.NUL)...),
			wantErr:   nil,
		},
		{
			name:      "UPC-E 6 digits",
			symbology: barcode.UPCE,
			data:      Upce6,
			want:      append([]byte{shared.GS, 'k', 1}, append([]byte("123456"), shared.NUL)...),
			wantErr:   nil,
		},
		{
			name:      "UPC-E 7 digits",
			symbology: barcode.UPCE,
			data:      Upce7,
			want:      append([]byte{shared.GS, 'k', 1}, append([]byte("1234567"), shared.NUL)...),
			wantErr:   nil,
		},
		{
			name:      "UPC-E 8 digits",
			symbology: barcode.UPCE,
			data:      Upce8,
			want:      append([]byte{shared.GS, 'k', 1}, append([]byte("12345678"), shared.NUL)...),
			wantErr:   nil,
		},
		{
			name:      "EAN13 12 digits",
			symbology: barcode.JAN13,
			data:      Ean1312,
			want:      append([]byte{shared.GS, 'k', 2}, append([]byte("123456789012"), shared.NUL)...),
			wantErr:   nil,
		},
		{
			name:      "EAN13 13 digits with check",
			symbology: barcode.JAN13,
			data:      Ean1313,
			want:      append([]byte{shared.GS, 'k', 2}, append([]byte("1234567890123"), shared.NUL)...),
			wantErr:   nil,
		},
		{
			name:      "EAN8 7 digits",
			symbology: barcode.JAN8,
			data:      Ean87,
			want:      append([]byte{shared.GS, 'k', 3}, append([]byte("1234567"), shared.NUL)...),
			wantErr:   nil,
		},
		{
			name:      "EAN8 8 digits with check",
			symbology: barcode.JAN8,
			data:      Ean88,
			want:      append([]byte{shared.GS, 'k', 3}, append([]byte("12345678"), shared.NUL)...),
			wantErr:   nil,
		},
		{
			name:      "EAN8 8 digits with check",
			symbology: barcode.JAN8,
			data:      Ean88,
			want:      append([]byte{shared.GS, 'k', 3}, append([]byte("12345678"), shared.NUL)...),
			wantErr:   nil,
		},
		{
			name:      "CODE39 with start/stop",
			symbology: barcode.CODE39,
			data:      Code39Ss,
			want:      append([]byte{shared.GS, 'k', 4}, append([]byte("*ABC-123*"), shared.NUL)...),
			wantErr:   nil,
		},
		{
			name:      "CODE39 without start/stop",
			symbology: barcode.CODE39,
			data:      Code39Noss,
			want:      append([]byte{shared.GS, 'k', 4}, append([]byte("ABC-123"), shared.NUL)...),
			wantErr:   nil,
		},
		{
			name:      "ITF even digits",
			symbology: barcode.ITF,
			data:      ITFEven,
			want:      append([]byte{shared.GS, 'k', 5}, append([]byte("123456"), shared.NUL)...),
			wantErr:   nil,
		},
		{
			name:      "ITF minimum (2 digits)",
			symbology: barcode.ITF,
			data:      ITFMin,
			want:      append([]byte{shared.GS, 'k', 5}, append([]byte("12"), shared.NUL)...),
			wantErr:   nil,
		},
		{
			name:      "ITF odd digits (error)",
			symbology: barcode.ITF,
			data:      ITFOdd,
			want:      nil,
			wantErr:   barcode.ErrOddITFLength,
		},
		{
			name:      "CODABAR with start/stop A-B",
			symbology: barcode.CODABAR,
			data:      CodabarSs1,
			want:      append([]byte{shared.GS, 'k', 6}, append([]byte("A123456B"), shared.NUL)...),
			wantErr:   nil,
		},
		{
			name:      "CODABAR with start/stop C-D",
			symbology: barcode.CODABAR,
			data:      CodabarSs2,
			want:      append([]byte{shared.GS, 'k', 6}, append([]byte("C123456D"), shared.NUL)...),
			wantErr:   nil,
		},
		{
			name:      "CODABAR lowercase start/stop a-b",
			symbology: barcode.CODABAR,
			data:      CodabarLss,
			want:      append([]byte{shared.GS, 'k', 6}, append([]byte("a123456d"), shared.NUL)...),
			wantErr:   nil,
		},
		{
			name:      "empty data",
			symbology: barcode.UPCA,
			data:      []byte{},
			want:      nil,
			wantErr:   barcode.ErrDataTooShort,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cmd.PrintBarcode(tt.symbology, tt.data)

			if !testutils.AssertErrorOccurred(t, err, tt.wantErr != nil, "PrintBarcode") {
				return
			}
			if tt.wantErr != nil {
				testutils.AssertError(t, err, tt.wantErr)
				return
			}

			testutils.AssertBytes(t, got, tt.want, "PrintBarcode(%v, %q)", tt.symbology, tt.data)
		})
	}
}

// ============================================================================
// Barcode Printing Tests - Function B
// ============================================================================

func TestCommands_PrintBarcode_FunctionB(t *testing.T) {
	cmd := barcode.NewCommands()

	tests := []struct {
		name      string
		symbology barcode.Symbology
		data      []byte
		want      []byte
		wantErr   error
	}{
		{
			name:      "UPC-A B 11 digits",
			symbology: barcode.UPCAB,
			data:      []byte("12345678901"),
			want:      append([]byte{shared.GS, 'k', 65, 11}, []byte("12345678901")...),
			wantErr:   nil,
		},
		{
			name:      "UPC-A B 12 digits",
			symbology: barcode.UPCAB,
			data:      []byte("123456789012"),
			want:      append([]byte{shared.GS, 'k', 65, 12}, []byte("123456789012")...),
			wantErr:   nil,
		},
		{
			name:      "UPC-E B 6 digits",
			symbology: barcode.UPCEB,
			data:      []byte("012345"),
			want:      append([]byte{shared.GS, 'k', 66, 6}, []byte("012345")...),
			wantErr:   nil,
		},
		{
			name:      "EAN13 B",
			symbology: barcode.EAN13,
			data:      []byte("1234567890123"),
			want:      append([]byte{shared.GS, 'k', 67, 13}, []byte("1234567890123")...),
			wantErr:   nil,
		},
		{
			name:      "EAN8 B",
			symbology: barcode.EAN8,
			data:      []byte("12345678"),
			want:      append([]byte{shared.GS, 'k', 68, 8}, []byte("12345678")...),
			wantErr:   nil,
		},
		{
			name:      "CODE39 B",
			symbology: barcode.CODE39B,
			data:      []byte("CODE39"),
			want:      append([]byte{shared.GS, 'k', 69, 6}, []byte("CODE39")...),
			wantErr:   nil,
		},
		{
			name:      "CODE39 B with special",
			symbology: barcode.CODE39B,
			data:      []byte("*ABC-123*"),
			want:      append([]byte{shared.GS, 'k', 69, 9}, []byte("*ABC-123*")...),
			wantErr:   nil,
		},
		{
			name:      "ITF B even",
			symbology: barcode.ITFB,
			data:      []byte("123456"),
			want:      append([]byte{shared.GS, 'k', 70, 6}, []byte("123456")...),
			wantErr:   nil,
		},
		{
			name:      "ITF B odd (error)",
			symbology: barcode.ITFB,
			data:      []byte("12345"),
			want:      nil,
			wantErr:   barcode.ErrOddITFLength,
		},
		{
			name:      "CODABAR B",
			symbology: barcode.CODABARB,
			data:      []byte("A1234B"),
			want:      append([]byte{shared.GS, 'k', 71, 6}, []byte("A1234B")...),
			wantErr:   nil,
		},
		{
			name:      "CODE93",
			symbology: barcode.CODE93,
			data:      []byte("CODE93"),
			want:      append([]byte{shared.GS, 'k', 72, 6}, []byte("CODE93")...),
			wantErr:   nil,
		},
		{
			name:      "CODE128 with code set",
			symbology: barcode.CODE128,
			data:      []byte{'{', 'B', 'H', 'e', 'l', 'l', 'o'},
			want:      append([]byte{shared.GS, 'k', 73, 7}, []byte{'{', 'B', 'H', 'e', 'l', 'l', 'o'}...),
			wantErr:   nil,
		},
		{
			name:      "CODE128 without code set (error)",
			symbology: barcode.CODE128,
			data:      []byte("Hello"),
			want:      nil,
			wantErr:   barcode.ErrCode128NoCodeSet,
		},
		{
			name:      "CODE128 single byte (error)",
			symbology: barcode.CODE128,
			data:      []byte("A"),
			want:      nil,
			wantErr:   barcode.ErrCode128NoCodeSet,
		},
		{
			name:      "GS1-128 with code set",
			symbology: barcode.GS1128,
			data:      []byte{'{', 'C', '0', '1', '2', '3', '4', '5'},
			want:      append([]byte{shared.GS, 'k', 74, 8}, []byte{'{', 'C', '0', '1', '2', '3', '4', '5'}...),
			wantErr:   nil,
		},
		{
			name:      "GS1-128 without code set (error)",
			symbology: barcode.GS1128,
			data:      []byte("12345"),
			want:      nil,
			wantErr:   barcode.ErrCode128NoCodeSet,
		},
		{
			name:      "CODE128 Auto",
			symbology: barcode.CODE128Auto,
			data:      []byte("Hello123"),
			want:      append([]byte{shared.GS, 'k', 79, 8}, []byte("Hello123")...),
			wantErr:   nil,
		},
		{
			name:      "GS1 DataBar Omnidirectional",
			symbology: barcode.GS1DataBarOmni,
			data:      []byte("0123456789012"),
			want:      append([]byte{shared.GS, 'k', 75, 13}, []byte("0123456789012")...),
			wantErr:   nil,
		},
		{
			name:      "GS1 DataBar Truncated",
			symbology: barcode.GS1DataBarTrunc,
			data:      []byte("0123456789012"),
			want:      append([]byte{shared.GS, 'k', 76, 13}, []byte("0123456789012")...),
			wantErr:   nil,
		},
		{
			name:      "GS1 DataBar Limited",
			symbology: barcode.GS1DataBarLim,
			data:      []byte("0123456789012"),
			want:      append([]byte{shared.GS, 'k', 77, 13}, []byte("0123456789012")...),
			wantErr:   nil,
		},
		{
			name:      "GS1 DataBar Expanded",
			symbology: barcode.GS1DataBarExp,
			data:      []byte("(01)12345678901234"),
			want:      append([]byte{shared.GS, 'k', 78, 18}, []byte("(01)12345678901234")...),
			wantErr:   nil,
		},
		{
			name:      "data too long",
			symbology: barcode.CODE39B,
			data:      testutils.RepeatByte(256, 'A'),
			want:      nil,
			wantErr:   barcode.ErrDataTooLong,
		},
		{
			name:      "maximum valid length",
			symbology: barcode.CODE39B,
			data:      testutils.RepeatByte(255, 'B'),
			want: append([]byte{shared.GS, 'k', 69, 255},
				testutils.RepeatByte(255, 'B')...),
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cmd.PrintBarcode(tt.symbology, tt.data)

			if !testutils.AssertErrorOccurred(t, err, tt.wantErr != nil, "PrintBarcode") {
				return
			}
			if tt.wantErr != nil {
				testutils.AssertError(t, err, tt.wantErr)
				return
			}

			testutils.AssertBytes(t, got, tt.want, "PrintBarcode(%v, %q)", tt.symbology, string(tt.data))
		})
	}
}

// ============================================================================
// PrintBarcodeWithCodeSet Tests
// ============================================================================

func TestCommands_PrintBarcodeWithCodeSet(t *testing.T) {
	cmd := barcode.NewCommands()

	tests := []struct {
		name      string
		symbology barcode.Symbology
		codeSet   barcode.Code128Set
		data      []byte
		want      []byte
		wantErr   error
	}{
		{
			name:      "CODE128 Set A uppercase",
			symbology: barcode.CODE128,
			codeSet:   barcode.Code128SetA,
			data:      []byte("HELLO"),
			want:      append([]byte{shared.GS, 'k', 73, 7}, []byte{'{', 'A', 'H', 'E', 'L', 'L', 'O'}...),
			wantErr:   nil,
		},
		{
			name:      "CODE128 Set A control chars",
			symbology: barcode.CODE128,
			codeSet:   barcode.Code128SetA,
			data:      []byte("ABC\x00DEF"),
			want:      append([]byte{shared.GS, 'k', 73, 9}, []byte{'{', 'A', 'A', 'B', 'C', 0x00, 'D', 'E', 'F'}...),
			wantErr:   nil,
		},
		{
			name:      "CODE128 Set B mixed case",
			symbology: barcode.CODE128,
			codeSet:   barcode.Code128SetB,
			data:      []byte("Hello"),
			want:      append([]byte{shared.GS, 'k', 73, 7}, []byte{'{', 'B', 'H', 'e', 'l', 'l', 'o'}...),
			wantErr:   nil,
		},
		{
			name:      "CODE128 Set B with numbers",
			symbology: barcode.CODE128,
			codeSet:   barcode.Code128SetB,
			data:      []byte("Item123"),
			want:      append([]byte{shared.GS, 'k', 73, 9}, []byte{'{', 'B', 'I', 't', 'e', 'm', '1', '2', '3'}...),
			wantErr:   nil,
		},
		{
			name:      "CODE128 Set C numeric pairs",
			symbology: barcode.CODE128,
			codeSet:   barcode.Code128SetC,
			data:      []byte("123456"),
			want:      append([]byte{shared.GS, 'k', 73, 8}, []byte{'{', 'C', '1', '2', '3', '4', '5', '6'}...),
			wantErr:   nil,
		},
		{
			name:      "CODE128 Set C long numeric",
			symbology: barcode.CODE128,
			codeSet:   barcode.Code128SetC,
			data:      []byte("00112233445566778899"),
			want: append([]byte{shared.GS, 'k', 73, 22},
				append([]byte{'{', 'C'}, []byte("00112233445566778899")...)...),
			wantErr: nil,
		},
		{
			name:      "GS1-128 Set A",
			symbology: barcode.GS1128,
			codeSet:   barcode.Code128SetA,
			data:      []byte("BATCH001"),
			want:      append([]byte{shared.GS, 'k', 74, 10}, []byte{'{', 'A', 'B', 'A', 'T', 'C', 'H', '0', '0', '1'}...),
			wantErr:   nil,
		},
		{
			name:      "GS1-128 Set B with AI",
			symbology: barcode.GS1128,
			codeSet:   barcode.Code128SetB,
			data:      []byte("(01)12345678901231"),
			want:      append([]byte{shared.GS, 'k', 74, 20}, append([]byte{'{', 'B'}, []byte("(01)12345678901231")...)...),
			wantErr:   nil,
		},
		{
			name:      "GS1-128 Set C numeric",
			symbology: barcode.GS1128,
			codeSet:   barcode.Code128SetC,
			data:      []byte("01234567890123"),
			want: append([]byte{shared.GS, 'k', 74, 16},
				append([]byte{'{', 'C'}, []byte("01234567890123")...)...),
			wantErr: nil,
		},
		{
			name:      "invalid symbology CODE39",
			symbology: barcode.CODE39B,
			codeSet:   barcode.Code128SetA,
			data:      []byte("DATA"),
			want:      nil,
			wantErr:   barcode.ErrSymbology,
		},
		{
			name:      "invalid symbology UPCA",
			symbology: barcode.UPCA,
			codeSet:   barcode.Code128SetB,
			data:      []byte("12345678901"),
			want:      nil,
			wantErr:   barcode.ErrSymbology,
		},
		{
			name:      "invalid code set 64",
			symbology: barcode.CODE128,
			codeSet:   64,
			data:      []byte("DATA"),
			want:      nil,
			wantErr:   barcode.ErrCode128Set,
		},
		{
			name:      "invalid code set 68",
			symbology: barcode.CODE128,
			codeSet:   68,
			data:      []byte("DATA"),
			want:      nil,
			wantErr:   barcode.ErrCode128Set,
		},
		{
			name:      "invalid code set 99",
			symbology: barcode.CODE128,
			codeSet:   99,
			data:      []byte("DATA"),
			want:      nil,
			wantErr:   barcode.ErrCode128Set,
		},
		{
			name:      "data too long",
			symbology: barcode.CODE128,
			codeSet:   barcode.Code128SetB,
			data:      testutils.RepeatByte(254, 'A'),
			want:      nil,
			wantErr:   barcode.ErrDataTooLong,
		},
		{
			name:      "maximum valid length CODE128",
			symbology: barcode.CODE128,
			codeSet:   barcode.Code128SetB,
			data:      testutils.RepeatByte(253, 'X'), // 253 + 2 for prefix = 255
			want: append([]byte{shared.GS, 'k', 73, 255},
				append([]byte{'{', 'B'}, testutils.RepeatByte(253, 'X')...)...),
			wantErr: nil,
		},
		{
			name:      "empty data",
			symbology: barcode.CODE128,
			codeSet:   barcode.Code128SetA,
			data:      []byte{},
			want: append([]byte{shared.GS, 'k', 73, 2},
				[]byte{'{', 'A'}...),
			wantErr: nil,
		},
		{
			name:      "CODE128 Set B with escaping {",
			symbology: barcode.CODE128,
			codeSet:   barcode.Code128SetB,
			data:      []byte("USER{ID}"),
			// Length: 2 (prefijo) + 4 (USER) + 2 ({{) + 2 (ID) + 1 (}) = 11 bytes
			want:    append([]byte{shared.GS, 'k', 73, 11}, []byte{'{', 'B', 'U', 'S', 'E', 'R', '{', '{', 'I', 'D', '}'}...),
			wantErr: nil,
		},
		{
			name:      "CODE128 Set B multiple braces",
			symbology: barcode.CODE128,
			codeSet:   barcode.Code128SetB,
			data:      []byte("A{{B"),
			// Length: 2 (prefijo) + 1 (A) + 2 ({{) + 2 ({{) + 1 (B) = 8 bytes
			want:    append([]byte{shared.GS, 'k', 73, 8}, []byte{'{', 'B', 'A', '{', '{', '{', '{', 'B'}...),
			wantErr: nil,
		},
		{
			name:      "CODE128 Set A no escaping",
			symbology: barcode.CODE128,
			codeSet:   barcode.Code128SetA,
			data:      []byte("TEST{"),
			want:      append([]byte{shared.GS, 'k', 73, 7}, []byte{'{', 'A', 'T', 'E', 'S', 'T', '{'}...),
			wantErr:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cmd.PrintBarcodeWithCodeSet(tt.symbology, tt.codeSet, tt.data)

			if !testutils.AssertErrorOccurred(t, err, tt.wantErr != nil, "PrintBarcodeWithCodeSet") {
				return
			}
			if tt.wantErr != nil {
				testutils.AssertError(t, err, tt.wantErr)
				return
			}

			testutils.AssertBytes(t, got, tt.want, "PrintBarcodeWithCodeSet(%v, %v, %q)",
				tt.symbology, tt.codeSet, string(tt.data))
		})
	}
}

func TestCommands_EdgeCases(t *testing.T) {
	cmd := barcode.NewCommands()

	t.Run("maximum data at boundaries", func(t *testing.T) {
		// Test exactly 255 bytes
		maxData := testutils.RepeatByte(255, 'X')
		_, err := cmd.PrintBarcode(barcode.CODE39B, maxData)
		testutils.AssertError(t, err, nil)

		// Test 254 bytes (should work)
		almostMaxData := testutils.RepeatByte(254, 'Y')
		_, err = cmd.PrintBarcode(barcode.CODE39B, almostMaxData)
		testutils.AssertError(t, err, nil)
	})

	t.Run("special characters in different positions", func(t *testing.T) {
		// Test null byte in middle of data
		dataWithNull := []byte("ABC\x00DEF")
		_, err := cmd.PrintBarcodeWithCodeSet(
			barcode.CODE128,
			barcode.Code128SetA,
			dataWithNull,
		)
		testutils.AssertError(t, err, nil)
	})

	t.Run("CODE128 Set B escaping exceeds max length", func(t *testing.T) {
		// Creamos datos de 253 bytes (límite teórico: 255 - 2 prefijo = 253)
		// Si uno de ellos es '{', se duplicará y excederá el límite.
		data := testutils.RepeatByte(252, 'A')
		data = append(data, '{') // El byte 253 es una llave

		_, err := cmd.PrintBarcodeWithCodeSet(
			barcode.CODE128,
			barcode.Code128SetB,
			data,
		)
		// Esperamos error porque 252 + 2('{') + 2(prefijo) = 256 bytes total
		testutils.AssertError(t, err, barcode.ErrDataTooLong)
	})
}

// ============================================================================
// Helper Functions Tests
// ============================================================================

func TestValidateNumericData(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		want bool
	}{
		{
			name: "all numeric",
			data: []byte("0123456789"),
			want: true,
		},
		{
			name: "single digit",
			data: []byte("5"),
			want: true,
		},
		{
			name: "contains letter",
			data: []byte("12345A"),
			want: false,
		},
		{
			name: "contains space",
			data: []byte("123 456"),
			want: false,
		},
		{
			name: "empty",
			data: []byte{},
			want: true,
		},
		{
			name: "special characters",
			data: []byte("123-456"),
			want: false,
		},
		{
			name: "decimal point",
			data: []byte("123.456"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := barcode.ValidateNumericData(tt.data)
			if got != tt.want {
				t.Errorf("ValidateNumericData(%q) = %v, want %v", tt.data, got, tt.want)
			}

			// Also verify with assertion helper for consistency
			if tt.want {
				testutils.AssertNumeric(t, tt.data, "ValidateNumericData testutils: %s", tt.name)
			}
		})
	}
}

func TestValidateCode39Data(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		want bool
	}{
		{
			name: "valid alphanumeric",
			data: []byte("CODE39"),
			want: true,
		},
		{
			name: "valid with digits",
			data: []byte("ABC123"),
			want: true,
		},
		{
			name: "valid with start/stop",
			data: []byte("*CODE39*"),
			want: true,
		},
		{
			name: "valid with hyphen",
			data: []byte("CODE-39"),
			want: true,
		},
		{
			name: "valid with dot",
			data: []byte("CODE.39"),
			want: true,
		},
		{
			name: "valid with slash",
			data: []byte("CODE/39"),
			want: true,
		},
		{
			name: "valid with plus",
			data: []byte("CODE+39"),
			want: true,
		},
		{
			name: "valid with percent",
			data: []byte("CODE%39"),
			want: true,
		},
		{
			name: "valid with dollar",
			data: []byte("CODE$39"),
			want: true,
		},
		{
			name: "valid with space",
			data: []byte("CODE 39"),
			want: true,
		},
		{
			name: "valid all special",
			data: []byte(" $%*+-./ "),
			want: true,
		},
		{
			name: "invalid lowercase",
			data: []byte("code39"),
			want: false,
		},
		{
			name: "invalid mixed case",
			data: []byte("Code39"),
			want: false,
		},
		{
			name: "invalid character @",
			data: []byte("CODE@39"),
			want: false,
		},
		{
			name: "invalid character #",
			data: []byte("CODE#39"),
			want: false,
		},
		{
			name: "invalid character &",
			data: []byte("CODE&39"),
			want: false,
		},
		{
			name: "empty",
			data: []byte{},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := barcode.ValidateCode39Data(tt.data)
			if got != tt.want {
				t.Errorf("ValidateCode39Data(%q) = %v, want %v", tt.data, got, tt.want)
			}

			// Verify with assertion helpers
			if tt.want && len(tt.data) > 0 {
				// CODE39 uses uppercase and special chars
				allowed := []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ $%*+-./")
				testutils.AssertContainsOnly(t, tt.data, allowed, "CODE39 validation: %s", tt.name)
			}
		})
	}
}

func TestValidateCodabarData(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		want bool
	}{
		{
			name: "valid uppercase A-B",
			data: []byte("A12345B"),
			want: true,
		},
		{
			name: "valid uppercase C-D",
			data: []byte("C67890D"),
			want: true,
		},
		{
			name: "valid uppercase A-D",
			data: []byte("A54321D"),
			want: true,
		},
		{
			name: "valid lowercase a-b",
			data: []byte("a12345b"),
			want: true,
		},
		{
			name: "valid lowercase c-d",
			data: []byte("c67890d"),
			want: true,
		},
		{
			name: "valid mixed case A-d",
			data: []byte("A12345d"),
			want: true,
		},
		{
			name: "valid mixed case c-B",
			data: []byte("c12345B"),
			want: true,
		},
		{
			name: "valid minimum length",
			data: []byte("AB"),
			want: true,
		},
		{
			name: "invalid start E",
			data: []byte("E12345B"),
			want: false,
		},
		{
			name: "invalid start 1",
			data: []byte("112345B"),
			want: false,
		},
		{
			name: "invalid stop E",
			data: []byte("A12345E"),
			want: false,
		},
		{
			name: "invalid stop 5",
			data: []byte("A12345"),
			want: false,
		},
		{
			name: "missing start/stop",
			data: []byte("12345"),
			want: false,
		},
		{
			name: "too short single char",
			data: []byte("A"),
			want: false,
		},
		{
			name: "empty",
			data: []byte{},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := barcode.ValidateCodabarData(tt.data)
			if got != tt.want {
				t.Errorf("ValidateCodabarData(%q) = %v, want %v", tt.data, got, tt.want)
			}
		})
	}
}

// ============================================================================
// Error Cases Tests
// ============================================================================

func TestCommands_InvalidSymbology(t *testing.T) {
	cmd := barcode.NewCommands()

	invalidSymbologies := []barcode.Symbology{
		7,   // Gap in Function A
		64,  // Gap before Function B
		80,  // Beyond CODE128Auto
		100, // Way beyond valid range
		200,
		255,
	}

	for _, symbology := range invalidSymbologies {
		t.Run(string(rune(symbology)), func(t *testing.T) {
			_, err := cmd.PrintBarcode(symbology, []byte("testutils"))
			if err == nil {
				t.Errorf("PrintBarcode with symbology %d should return error", symbology)
			}
			testutils.AssertError(t, err, barcode.ErrSymbology)
		})
	}
}

func TestCommands_SpecialSymbology(t *testing.T) {
	t.Run("UPC-A requires numeric data", func(t *testing.T) {
		validUPC := []byte("12345678901")
		invalidUPC := []byte("12345A78901")

		// Use the assertion helper
		testutils.AssertNumeric(t, validUPC, "UPC-A data should be numeric")

		// This will fail with a descriptive error
		// testutils.AssertNumeric(t, invalidUPC, "UPC-A validation testutils")

		// Verify barcode package agrees
		if !barcode.ValidateNumericData(validUPC) {
			t.Error("ValidateNumericData should accept valid UPC")
		}
		if barcode.ValidateNumericData(invalidUPC) {
			t.Error("ValidateNumericData should reject invalid UPC")
		}
	})

	t.Run("CODE39 requires uppercase and special chars", func(t *testing.T) {
		validCode39 := []byte("ABC-123")

		// Check it's uppercase where alphabetic
		testutils.AssertUppercase(t, []byte("ABC"), "CODE39 letters should be uppercase")

		// Check allowed character set
		allowed := []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ $%*+-./")
		testutils.AssertContainsOnly(t, validCode39, allowed, "CODE39 character validation")

		// Verify with barcode package
		if !barcode.ValidateCode39Data(validCode39) {
			t.Error("ValidateCode39Data should accept valid data")
		}
	})

	t.Run("ITF requires even length numeric", func(t *testing.T) {
		evenData := []byte("123456")
		oddData := []byte("12345")

		testutils.AssertNumeric(t, evenData, "ITF data should be numeric")
		testutils.AssertEvenLength(t, evenData, "ITF data should have even length")

		// This should fail
		if testutils.IsEvenLength(oddData) {
			t.Error("Odd data should not pass even length check")
		}
	})

	t.Run("Function A barcodes have null terminator", func(t *testing.T) {
		cmd := barcode.NewCommands()
		result, _ := cmd.PrintBarcode(barcode.UPCA, []byte("12345678901"))

		testutils.AssertHasNullTerminator(t, result, "Function A barcodes should end with NUL")
	})

	t.Run("barcode data length validation", func(t *testing.T) {
		data := []byte("TEST123")
		testutils.AssertValidLength(t, data, 1, 255, "Barcode data length check")

		// Test printable ASCII for CODE39
		testutils.AssertPrintableASCII(t, data, "CODE39 should use printable ASCII")
	})

	t.Run("CODABAR start/stop validation", func(t *testing.T) {
		validCodabar := []byte("A12345B")

		// Check start/stop characters are in allowed range
		startStop := []byte{'A', 'B', 'C', 'D', 'a', 'b', 'c', 'd'}
		if !testutils.ContainsAny([]byte{validCodabar[0]}, startStop) {
			t.Error("Invalid CODABAR start character")
		}
		if !testutils.ContainsAny([]byte{validCodabar[len(validCodabar)-1]}, startStop) {
			t.Error("Invalid CODABAR stop character")
		}

		// Verify with barcode package
		if !barcode.ValidateCodabarData(validCodabar) {
			t.Error("ValidateCodabarData should accept valid data")
		}
	})
}

// ============================================================================
// Boundary Tests
// ============================================================================

func TestCommands_BoundaryValues(t *testing.T) {
	cmd := barcode.NewCommands()

	t.Run("height boundaries", func(t *testing.T) {
		// Valid boundaries
		validHeights := []barcode.Height{1, 2, 100, 254, 255}
		for _, height := range validHeights {
			_, err := cmd.SetBarcodeHeight(height)
			if err != nil {
				t.Errorf("SetBarcodeHeight(%d) should be valid: %v", height, err)
			}
		}

		// Invalid boundary
		_, err := cmd.SetBarcodeHeight(0)
		if err == nil {
			t.Error("SetBarcodeHeight(0) should return error")
		}
	})

	t.Run("width boundaries", func(t *testing.T) {
		// Valid boundaries
		validWidths := []barcode.Width{2, 3, 4, 5, 6, 68, 69, 75, 76}
		for _, width := range validWidths {
			_, err := cmd.SetBarcodeWidth(width)
			if err != nil {
				t.Errorf("SetBarcodeWidth(%d) should be valid: %v", width, err)
			}
		}

		// Invalid boundaries
		invalidWidths := []barcode.Width{0, 1, 7, 67, 77, 100, 255}
		for _, width := range invalidWidths {
			_, err := cmd.SetBarcodeWidth(width)
			if err == nil {
				t.Errorf("SetBarcodeWidth(%d) should return error", width)
			}
		}
	})

	t.Run("data length boundaries", func(t *testing.T) {
		// Maximum valid length for Function B
		maxData := testutils.RepeatByte(255, 'X')
		_, err := cmd.PrintBarcode(barcode.CODE39B, maxData)
		if err != nil {
			t.Errorf("PrintBarcode with 255 bytes should be valid: %v", err)
		}

		// Over maximum
		overData := testutils.RepeatByte(256, 'Y')
		_, err = cmd.PrintBarcode(barcode.CODE39B, overData)
		if err == nil {
			t.Error("PrintBarcode with 256 bytes should return error")
		}
	})
}
