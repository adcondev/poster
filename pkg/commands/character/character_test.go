package character_test

import (
	"testing"

	"github.com/adcondev/poster/internal/testutils"
	"github.com/adcondev/poster/pkg/commands/character"
	"github.com/adcondev/poster/pkg/commands/shared"
)

// ============================================================================
// Utility Functions Tests
// ============================================================================

// Test_BuildCharacterSize tests the NewSize utility function which is used for SelectCharacterSize
func Test_BuildCharacterSize(t *testing.T) {
	// Setup
	tests := []struct {
		name    string
		width   byte
		height  byte
		want    character.Size
		wantErr error
	}{
		{
			name:    "normal size",
			width:   1,
			height:  1,
			want:    character.Size1x1,
			wantErr: nil,
		},
		{
			name:    "double width",
			width:   2,
			height:  1,
			want:    character.Size2x1,
			wantErr: nil,
		},
		{
			name:    "double height",
			width:   1,
			height:  2,
			want:    character.Size1x2,
			wantErr: nil,
		},
		{
			name:    "double size",
			width:   2,
			height:  2,
			want:    character.Size2x2,
			wantErr: nil,
		},
		{
			name:    "maximum size",
			width:   8,
			height:  8,
			want:    0x77,
			wantErr: nil,
		},
		{
			name:    "invalid width",
			width:   9,
			height:  1,
			want:    0,
			wantErr: character.ErrCharacterWidth,
		},
		{
			name:    "invalid height",
			width:   1,
			height:  9,
			want:    0,
			wantErr: character.ErrCharacterHeight,
		},
		{
			name:    "zero width",
			width:   0,
			height:  1,
			want:    0,
			wantErr: character.ErrCharacterWidth,
		},
		{
			name:    "zero height",
			width:   1,
			height:  0,
			want:    0,
			wantErr: character.ErrCharacterHeight,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Execute
			got, err := character.NewSize(tt.width, tt.height)

			// Verify
			if !testutils.AssertErrorOccurred(t, err, tt.wantErr != nil, "NewSize") {
				return
			}
			if tt.wantErr != nil {
				testutils.AssertError(t, err, tt.wantErr)
				return
			}

			// Verify
			if got != tt.want {
				t.Errorf("NewSize(%d, %d) = %v, want %v",
					tt.width, tt.height, got, tt.want)
			}
		})
	}
}

// ============================================================================
// Character Commands Tests
// ============================================================================

func TestCommands_SetRightSideCharacterSpacing(t *testing.T) {
	// Setup
	cmd := character.NewCommands()
	prefix := testutils.BuildCommand(shared.ESC, shared.SP)

	tests := []struct {
		name    string
		spacing character.Spacing
		want    []byte
	}{
		{
			name:    "no spacing",
			spacing: 0,
			want:    append(prefix, 0),
		},
		{
			name:    "normal spacing",
			spacing: 5,
			want:    append(prefix, 5),
		},
		{
			name:    "maximum spacing",
			spacing: 255,
			want:    append(prefix, 255),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Execute
			got := cmd.SetRightSideCharacterSpacing(tt.spacing)

			// Verify
			testutils.AssertBytes(t, got, tt.want, "SetRightSideCharacterSpacing(%d)", tt.spacing)
		})
	}
}

func TestCommands_SelectPrintModes(t *testing.T) {
	// Setup
	cmd := character.NewCommands()
	prefix := testutils.BuildCommand(shared.ESC, '!')

	tests := []struct {
		name string
		mode character.PrintMode
		want []byte
	}{
		{
			name: "normal mode",
			mode: character.FontAPm | character.EmphasizedOffPm | character.DoubleHeightOffPm | character.DoubleWidthOffPm | character.UnderlineOffPm,
			want: append(prefix, 0),
		},
		{
			name: "font B",
			mode: character.FontBPm,
			want: append(prefix, 0x01),
		},
		{
			name: "emphasized",
			mode: character.EmphasizedOnPm,
			want: append(prefix, 0x08),
		},
		{
			name: "double height",
			mode: character.DoubleHeightOnPm,
			want: append(prefix, 0x10),
		},
		{
			name: "double width",
			mode: character.DoubleWidthOnPm,
			want: append(prefix, 0x20),
		},
		{
			name: "underline",
			mode: character.UnderlineOnPm,
			want: append(prefix, 0x80),
		},
		{
			name: "combined modes",
			mode: character.FontBPm | character.EmphasizedOnPm | character.DoubleHeightOnPm | character.DoubleWidthOnPm | character.UnderlineOnPm,
			want: append(prefix, 0xB9),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Execute
			got := cmd.SelectPrintModes(tt.mode)

			// Verify
			testutils.AssertBytes(t, got, tt.want, "SelectPrintModes(%#x)", tt.mode)
		})
	}
}

func TestCommands_SetUnderlineMode(t *testing.T) {
	// Setup
	cmd := character.NewCommands()
	prefix := testutils.BuildCommand(shared.ESC, '-')

	tests := []struct {
		name    string
		mode    character.UnderlineMode
		want    []byte
		wantErr error
	}{
		{
			name:    "underline off",
			mode:    character.NoDot,
			want:    append(prefix, 0x00),
			wantErr: nil,
		},
		{
			name:    "underline 1 dot",
			mode:    character.OneDot,
			want:    append(prefix, 0x01),
			wantErr: nil,
		},
		{
			name:    "underline 2 dots",
			mode:    character.TwoDot,
			want:    append(prefix, 0x02),
			wantErr: nil,
		},
		{
			name:    "underline off ASCII",
			mode:    character.NoDotASCII,
			want:    append(prefix, '0'),
			wantErr: nil,
		},
		{
			name:    "underline 1 dot ASCII",
			mode:    character.OneDotASCII,
			want:    append(prefix, '1'),
			wantErr: nil,
		},
		{
			name:    "underline 2 dots ASCII",
			mode:    character.TwoDotASCII,
			want:    append(prefix, '2'),
			wantErr: nil,
		},
		{
			name:    "invalid mode",
			mode:    3,
			want:    nil,
			wantErr: character.ErrUnderlineMode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Execute
			got, err := cmd.SetUnderlineMode(tt.mode)

			// Verify error
			if !testutils.AssertErrorOccurred(t, err, tt.wantErr != nil, "SetUnderlineMode") {
				return
			}
			if tt.wantErr != nil {
				testutils.AssertError(t, err, tt.wantErr)
				return
			}

			// Verify result
			testutils.AssertBytes(t, got, tt.want, "SetUnderlineMode(%v)", tt.mode)
		})
	}
}

func TestCommands_SetEmphasizedMode(t *testing.T) {
	// Setup
	cmd := character.NewCommands()
	prefix := testutils.BuildCommand(shared.ESC, 'E')

	tests := []struct {
		name string
		mode character.EmphasizedMode
		want []byte
	}{
		{
			name: "emphasized off",
			mode: character.OffEm,
			want: append(prefix, 0x00),
		},
		{
			name: "emphasized on",
			mode: character.OnEm,
			want: append(prefix, 0x01),
		},
		{
			name: "any even number (LSB=0)",
			mode: 0xFE,
			want: append(prefix, 0xFE),
		},
		{
			name: "any odd number (LSB=1)",
			mode: 0xFF,
			want: append(prefix, 0xFF),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Execute
			got := cmd.SetEmphasizedMode(tt.mode)

			// Verify
			testutils.AssertBytes(t, got, tt.want, "SetEmphasizedMode(%d)", tt.mode)
		})
	}
}

func TestCommands_SetDoubleStrikeMode(t *testing.T) {
	// Setup
	cmd := character.NewCommands()
	prefix := testutils.BuildCommand(shared.ESC, 'G')

	tests := []struct {
		name string
		mode character.DoubleStrikeMode
		want []byte
	}{
		{
			name: "double-strike off",
			mode: character.OffDsm,
			want: append(prefix, 0x00),
		},
		{
			name: "double-strike on",
			mode: character.OnDsm,
			want: append(prefix, 0x01),
		},
		{
			name: "any even number (LSB=0)",
			mode: 0xFE,
			want: append(prefix, 0xFE),
		},
		{
			name: "any odd number (LSB=1)",
			mode: 0xFF,
			want: append(prefix, 0xFF),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Execute
			got := cmd.SetDoubleStrikeMode(tt.mode)

			// Verify
			testutils.AssertBytes(t, got, tt.want, "SetDoubleStrikeMode(%d)", tt.mode)
		})
	}
}

func TestCommands_SelectCharacterFont(t *testing.T) {
	// Setup
	cmd := character.NewCommands()
	prefix := testutils.BuildCommand(shared.ESC, 'M')

	tests := []struct {
		name    string
		font    character.FontType
		want    []byte
		wantErr error
	}{
		{
			name:    "font A",
			font:    character.FontA,
			want:    append(prefix, 0x00),
			wantErr: nil,
		},
		{
			name:    "font B",
			font:    character.FontB,
			want:    append(prefix, 0x01),
			wantErr: nil,
		},
		{
			name:    "font C",
			font:    character.FontC,
			want:    append(prefix, 0x02),
			wantErr: nil,
		},
		{
			name:    "font A ASCII",
			font:    character.FontAAscii,
			want:    append(prefix, '0'),
			wantErr: nil,
		},
		{
			name:    "special font A",
			font:    character.SpecialFontA,
			want:    append(prefix, 0x61),
			wantErr: nil,
		},
		{
			name:    "special font B",
			font:    character.SpecialFontB,
			want:    append(prefix, 0x62),
			wantErr: nil,
		},
		{
			name:    "invalid font",
			font:    99,
			want:    nil,
			wantErr: character.ErrCharacterFont,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Execute
			got, err := cmd.SelectCharacterFont(tt.font)

			// Verify error
			if !testutils.AssertErrorOccurred(t, err, tt.wantErr != nil, "SelectCharacterFont") {
				return
			}
			if tt.wantErr != nil {
				testutils.AssertError(t, err, tt.wantErr)
				return
			}

			// Verify result
			testutils.AssertBytes(t, got, tt.want, "SelectCharacterFont(%v)", tt.font)
		})
	}
}

func TestCommands_SelectInternationalCharacterSet(t *testing.T) {
	// Setup
	cmd := character.NewCommands()
	prefix := testutils.BuildCommand(shared.ESC, 'R')

	tests := []struct {
		name    string
		charset character.InternationalSet
		want    []byte
		wantErr error
	}{
		{
			name:    "USA",
			charset: character.USA,
			want:    append(prefix, 0x00),
			wantErr: nil,
		},
		{
			name:    "Japan",
			charset: character.Japan,
			want:    append(prefix, 0x08),
			wantErr: nil,
		},
		{
			name:    "India Devanagari",
			charset: character.IndiaDevanagari,
			want:    append(prefix, 0x42),
			wantErr: nil,
		},
		{
			name:    "India Marathi",
			charset: character.IndiaMarathi,
			want:    append(prefix, 0x52),
			wantErr: nil,
		},
		{
			name:    "invalid charset",
			charset: 200,
			want:    nil,
			wantErr: character.ErrCharacterSet,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Execute
			got, err := cmd.SelectInternationalCharacterSet(tt.charset)

			// Verify error
			if !testutils.AssertErrorOccurred(t, err, tt.wantErr != nil, "SelectInternationalCharacterSet") {
				return
			}
			if tt.wantErr != nil {
				testutils.AssertError(t, err, tt.wantErr)
				return
			}

			// Verify result
			testutils.AssertBytes(t, got, tt.want, "SelectInternationalCharacterSet(%v)", tt.charset)
		})
	}
}

func TestCommands_Set90DegreeClockwiseRotationMode(t *testing.T) {
	// Setup
	cmd := character.NewCommands()
	prefix := testutils.BuildCommand(shared.ESC, 'V')

	tests := []struct {
		name    string
		mode    character.RotationMode
		want    []byte
		wantErr error
	}{
		{
			name:    "no rotation",
			mode:    character.NoRotation,
			want:    append(prefix, 0x00),
			wantErr: nil,
		},
		{
			name:    "90 degree 1 dot spacing",
			mode:    character.On90Dot1,
			want:    append(prefix, 0x01),
			wantErr: nil,
		},
		{
			name:    "90 degree 1.5 dot spacing",
			mode:    character.On90Dot15,
			want:    append(prefix, 0x02),
			wantErr: nil,
		},
		{
			name:    "no rotation ASCII",
			mode:    character.NoRotationASCII,
			want:    append(prefix, '0'),
			wantErr: nil,
		},
		{
			name:    "invalid rotation",
			mode:    99,
			want:    nil,
			wantErr: character.ErrRotationMode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Execute
			got, err := cmd.Set90DegreeClockwiseRotationMode(tt.mode)

			// Verify error
			if !testutils.AssertErrorOccurred(t, err, tt.wantErr != nil, "Set90DegreeClockwiseRotationMode") {
				return
			}
			if tt.wantErr != nil {
				testutils.AssertError(t, err, tt.wantErr)
				return
			}

			// Verify result
			testutils.AssertBytes(t, got, tt.want, "Set90DegreeClockwiseRotationMode(%v)", tt.mode)
		})
	}
}

func TestCommands_SelectPrintColor(t *testing.T) {
	// Setup
	cmd := character.NewCommands()
	prefix := testutils.BuildCommand(shared.ESC, 'r')

	tests := []struct {
		name    string
		color   character.PrintColor
		want    []byte
		wantErr error
	}{
		{
			name:    "black",
			color:   character.Black,
			want:    append(prefix, 0x00),
			wantErr: nil,
		},
		{
			name:    "red",
			color:   character.Red,
			want:    append(prefix, 0x01),
			wantErr: nil,
		},
		{
			name:    "black ASCII",
			color:   character.BlackASCII,
			want:    append(prefix, '0'),
			wantErr: nil,
		},
		{
			name:    "red ASCII",
			color:   character.RedASCII,
			want:    append(prefix, '1'),
			wantErr: nil,
		},
		{
			name:    "invalid color",
			color:   99,
			want:    nil,
			wantErr: character.ErrPrintColor,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Execute
			got, err := cmd.SelectPrintColor(tt.color)

			// Verify error
			if !testutils.AssertErrorOccurred(t, err, tt.wantErr != nil, "SelectPrintColor") {
				return
			}
			if tt.wantErr != nil {
				testutils.AssertError(t, err, tt.wantErr)
				return
			}

			// Verify result
			testutils.AssertBytes(t, got, tt.want, "SelectPrintColor(%v)", tt.color)
		})
	}
}

func TestCommands_SelectCharacterCodeTable(t *testing.T) {
	// Setup
	cmd := character.NewCommands()
	prefix := testutils.BuildCommand(shared.ESC, 't')

	tests := []struct {
		name    string
		page    character.CodeTable
		want    []byte
		wantErr error
	}{
		{
			name:    "PC437",
			page:    character.PC437,
			want:    append(prefix, 0x00),
			wantErr: nil,
		},
		{
			name:    "Katakana",
			page:    character.Katakana,
			want:    append(prefix, 0x01),
			wantErr: nil,
		},
		{
			name:    "Devanagari",
			page:    character.Devanagari,
			want:    append(prefix, 0x42),
			wantErr: nil,
		},
		{
			name:    "Special254",
			page:    character.Special254,
			want:    append(prefix, 0xFE),
			wantErr: nil,
		},
		{
			name:    "invalid page",
			page:    100,
			want:    nil,
			wantErr: character.ErrCodeTablePage,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Execute
			got, err := cmd.SelectCharacterCodeTable(tt.page)

			// Verify error
			if !testutils.AssertErrorOccurred(t, err, tt.wantErr != nil, "SelectCharacterCodeTable") {
				return
			}
			if tt.wantErr != nil {
				testutils.AssertError(t, err, tt.wantErr)
				return
			}

			// Verify result
			testutils.AssertBytes(t, got, tt.want, "SelectCharacterCodeTable(%v)", tt.page)
		})
	}
}

func TestCommands_SetUpsideDownMode(t *testing.T) {
	// Setup
	cmd := character.NewCommands()
	prefix := testutils.BuildCommand(shared.ESC, '{')

	tests := []struct {
		name string
		mode character.UpsideDownMode
		want []byte
	}{
		{
			name: "upside-down off",
			mode: character.OffUdm,
			want: append(prefix, 0x00),
		},
		{
			name: "upside-down on",
			mode: character.OnUdm,
			want: append(prefix, 0x01),
		},
		{
			name: "any even number (LSB=0)",
			mode: 0xFE,
			want: append(prefix, 0xFE),
		},
		{
			name: "any odd number (LSB=1)",
			mode: 0xFF,
			want: append(prefix, 0xFF),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Execute
			got := cmd.SetUpsideDownMode(tt.mode)

			// Verify
			testutils.AssertBytes(t, got, tt.want, "SetUpsideDownMode(%d)", tt.mode)
		})
	}
}

func TestCommands_SelectCharacterSize(t *testing.T) {
	// Setup
	cmd := character.NewCommands()
	prefix := testutils.BuildCommand(shared.GS, '!')

	tests := []struct {
		name string
		size character.Size
		want []byte
	}{
		{
			name: "normal size (1x1)",
			size: character.Size1x1,
			want: append(prefix, 0x00),
		},
		{
			name: "double width (2x1)",
			size: character.Size2x1,
			want: append(prefix, 0x10),
		},
		{
			name: "double height (1x2)",
			size: character.Size1x2,
			want: append(prefix, 0x01),
		},
		{
			name: "double size (2x2)",
			size: character.Size2x2,
			want: append(prefix, 0x11),
		},
		{
			name: "quadruple width (4x1)",
			size: character.Size4x1,
			want: append(prefix, 0x30),
		},
		{
			name: "maximum size (8x8)",
			size: 0x77,
			want: append(prefix, 0x77),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Execute
			got := cmd.SelectCharacterSize(tt.size)

			// Verify
			testutils.AssertBytes(t, got, tt.want, "SelectCharacterSize(%#x)", tt.size)
		})
	}
}

func TestCommands_SetWhiteBlackReverseMode(t *testing.T) {
	// Setup
	cmd := character.NewCommands()
	prefix := testutils.BuildCommand(shared.GS, 'B')

	tests := []struct {
		name string
		mode character.ReverseMode
		want []byte
	}{
		{
			name: "reverse off",
			mode: character.OffRm,
			want: append(prefix, 0x00),
		},
		{
			name: "reverse on",
			mode: character.OnRm,
			want: append(prefix, 0x01),
		},
		{
			name: "reverse off ASCII",
			mode: character.OffRmASCII,
			want: append(prefix, '0'),
		},
		{
			name: "reverse on ASCII",
			mode: character.OnRmASCII,
			want: append(prefix, '1'),
		},
		{
			name: "any even number (LSB=0)",
			mode: 0xFE,
			want: append(prefix, 0xFE),
		},
		{
			name: "any odd number (LSB=1)",
			mode: 0xFF,
			want: append(prefix, 0xFF),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Execute
			got := cmd.SetWhiteBlackReverseMode(tt.mode)

			// Verify
			testutils.AssertBytes(t, got, tt.want, "SetWhiteBlackReverseMode(%d)", tt.mode)
		})
	}
}

func TestCommands_SetSmoothingMode(t *testing.T) {
	// Setup
	cmd := character.NewCommands()
	prefix := testutils.BuildCommand(shared.GS, 'b')

	tests := []struct {
		name string
		mode character.SmoothingMode
		want []byte
	}{
		{
			name: "smoothing off",
			mode: character.OffSm,
			want: append(prefix, 0x00),
		},
		{
			name: "smoothing on",
			mode: character.OnSm,
			want: append(prefix, 0x01),
		},
		{
			name: "smoothing off ASCII",
			mode: character.OffSmASCII,
			want: append(prefix, '0'),
		},
		{
			name: "smoothing on ASCII",
			mode: character.OnSmASCII,
			want: append(prefix, '1'),
		},
		{
			name: "any even number (LSB=0)",
			mode: 0xFE,
			want: append(prefix, 0xFE),
		},
		{
			name: "any odd number (LSB=1)",
			mode: 0xFF,
			want: append(prefix, 0xFF),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Execute
			got := cmd.SetSmoothingMode(tt.mode)

			// Verify
			testutils.AssertBytes(t, got, tt.want, "SetSmoothingMode(%d)", tt.mode)
		})
	}
}
