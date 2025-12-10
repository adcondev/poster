package printposition_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/adcondev/poster/pkg/commands/common"
	"github.com/adcondev/poster/pkg/commands/printposition"
)

// ============================================================================
// Commands Tests
// ============================================================================

// Naming Convention: Test{Struct}_{Method}_{Optional Scenario}

func TestCommands_HorizontalTab(t *testing.T) {
	cmd := printposition.NewCommands()
	got := cmd.HorizontalTab()
	want := []byte{common.HT}

	if !bytes.Equal(got, want) {
		t.Errorf("HorizontalTab() = %#v, want %#v", got, want)
	}
}

func TestCommands_SetAbsolutePrintPosition(t *testing.T) {
	cmd := printposition.NewCommands()

	tests := []struct {
		name     string
		position uint16
		want     []byte
	}{
		{
			name:     "zero position",
			position: 0,
			want:     []byte{common.ESC, '$', 0x00, 0x00},
		},
		{
			name:     "position 100",
			position: 100,
			want:     []byte{common.ESC, '$', 0x64, 0x00},
		},
		{
			name:     "position 256",
			position: 256,
			want:     []byte{common.ESC, '$', 0x00, 0x01},
		},
		{
			name:     "maximum position",
			position: 65535,
			want:     []byte{common.ESC, '$', 0xFF, 0xFF},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cmd.SetAbsolutePrintPosition(tt.position)
			if !bytes.Equal(got, tt.want) {
				t.Errorf("SetAbsolutePrintPosition(%d) = %#v, want %#v",
					tt.position, got, tt.want)
			}
		})
	}
}

func TestCommands_SetRelativePrintPosition(t *testing.T) {
	cmd := printposition.NewCommands()

	tests := []struct {
		name     string
		distance int16
		want     []byte
	}{
		{
			name:     "zero distance",
			distance: 0,
			want:     []byte{common.ESC, '\\', 0x00, 0x00},
		},
		{
			name:     "positive distance",
			distance: 100,
			want:     []byte{common.ESC, '\\', 0x64, 0x00},
		},
		{
			name:     "negative distance",
			distance: -100,
			want:     []byte{common.ESC, '\\', 0x9C, 0xFF}, // Two's complement of -100
		},
		{
			name:     "maximum positive",
			distance: 32767,
			want:     []byte{common.ESC, '\\', 0xFF, 0x7F},
		},
		{
			name:     "maximum negative",
			distance: -32768,
			want:     []byte{common.ESC, '\\', 0x00, 0x80},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cmd.SetRelativePrintPosition(tt.distance)
			if !bytes.Equal(got, tt.want) {
				t.Errorf("SetRelativePrintPosition(%d) = %#v, want %#v",
					tt.distance, got, tt.want)
			}
		})
	}
}

func TestCommands_SetHorizontalTabPositions(t *testing.T) {
	cmd := printposition.NewCommands()

	tests := []struct {
		name      string
		positions []byte
		want      []byte
		wantErr   bool
	}{
		{
			name:      "single tab",
			positions: []byte{8},
			want:      []byte{common.ESC, 'D', 8, common.NUL},
			wantErr:   false,
		},
		{
			name:      "multiple tabs",
			positions: []byte{8, 16, 24, 32},
			want:      []byte{common.ESC, 'D', 8, 16, 24, 32, common.NUL},
			wantErr:   false,
		},
		{
			name:      "clear all tabs",
			positions: []byte{},
			want:      []byte{common.ESC, 'D', common.NUL},
			wantErr:   false,
		},
		{
			name: "maximum tabs (32)",
			positions: func() []byte {
				tabs := make([]byte, 32)
				for i := range tabs {
					tabs[i] = byte(i + 1)
				}
				return tabs
			}(),
			// FIXME: change anonymous func to utils helpers
			want: func() []byte {
				cmd := []byte{common.ESC, 'D'}
				for i := 1; i <= 32; i++ {
					cmd = append(cmd, byte(i))
				}
				cmd = append(cmd, common.NUL)
				return cmd
			}(),
			wantErr: false,
		},
		{
			name: "too many tabs",
			// FIXME: change anonymous func to utils helpers
			positions: func() []byte {
				tabs := make([]byte, 33)
				for i := range tabs {
					tabs[i] = byte(i + 1)
				}
				return tabs
			}(),
			want:    nil,
			wantErr: true,
		},
		{
			name:      "invalid tab (zero)",
			positions: []byte{0, 8, 16},
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "not ascending order",
			positions: []byte{16, 8, 24},
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "duplicate position",
			positions: []byte{8, 16, 16, 24},
			want:      nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cmd.SetHorizontalTabPositions(tt.positions)

			// Standardized error checking
			if (err != nil) != tt.wantErr {
				t.Errorf("SetHorizontalTabPositions() error = %v, wantErr %v",
					err, tt.wantErr)
				return
			}

			// Check specific error type if expecting error
			if tt.wantErr && err != nil {
				var baseErr error
				switch tt.name {
				case "too many tabs":
					baseErr = printposition.ErrTooManyTabPositions
				case "invalid tab (zero)", "not ascending order", "duplicate position":
					baseErr = printposition.ErrTabPosition
				default:
					baseErr = nil
				}

				if baseErr != nil && !errors.Is(err, baseErr) {
					t.Errorf("SetHorizontalTabPositions() error = %v, want error containing %v",
						err, baseErr)
				}
				return
			}

			// Check result if no error expected
			if !tt.wantErr && !bytes.Equal(got, tt.want) {
				t.Errorf("SetHorizontalTabPositions() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestCommands_SelectJustification(t *testing.T) {
	cmd := printposition.NewCommands()

	tests := []struct {
		name    string
		mode    printposition.Justification
		want    []byte
		wantErr bool
	}{
		{
			name:    "left justification",
			mode:    0,
			want:    []byte{common.ESC, 'a', 0},
			wantErr: false,
		},
		{
			name:    "center justification",
			mode:    1,
			want:    []byte{common.ESC, 'a', 1},
			wantErr: false,
		},
		{
			name:    "right justification",
			mode:    2,
			want:    []byte{common.ESC, 'a', 2},
			wantErr: false,
		},
		{
			name:    "left justification ASCII",
			mode:    '0',
			want:    []byte{common.ESC, 'a', '0'},
			wantErr: false,
		},
		{
			name:    "center justification ASCII",
			mode:    '1',
			want:    []byte{common.ESC, 'a', '1'},
			wantErr: false,
		},
		{
			name:    "right justification ASCII",
			mode:    '2',
			want:    []byte{common.ESC, 'a', '2'},
			wantErr: false,
		},
		{
			name:    "invalid justification",
			mode:    3,
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cmd.SelectJustification(tt.mode)

			// Standardized error checking
			if (err != nil) != tt.wantErr {
				t.Errorf("SelectJustification(%v) error = %v, wantErr %v",
					tt.mode, err, tt.wantErr)
				return
			}

			// Check specific error type if expecting error
			if tt.wantErr && err != nil {
				if !errors.Is(err, printposition.ErrJustification) {
					t.Errorf("SelectJustification(%v) error = %v, want %v",
						tt.mode, err, printposition.ErrJustification)
				}
				return
			}

			// Check result if no error expected
			if !tt.wantErr && !bytes.Equal(got, tt.want) {
				t.Errorf("SelectJustification(%v) = %#v, want %#v",
					tt.mode, got, tt.want)
			}
		})
	}
}

func TestCommands_SetLeftMargin(t *testing.T) {
	cmd := printposition.NewCommands()

	tests := []struct {
		name   string
		margin uint16
		want   []byte
	}{
		{
			name:   "zero margin",
			margin: 0,
			want:   []byte{common.GS, 'L', 0x00, 0x00},
		},
		{
			name:   "margin 100",
			margin: 100,
			want:   []byte{common.GS, 'L', 0x64, 0x00},
		},
		{
			name:   "margin 256",
			margin: 256,
			want:   []byte{common.GS, 'L', 0x00, 0x01},
		},
		{
			name:   "maximum margin",
			margin: 65535,
			want:   []byte{common.GS, 'L', 0xFF, 0xFF},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cmd.SetLeftMargin(tt.margin)
			if !bytes.Equal(got, tt.want) {
				t.Errorf("SetLeftMargin(%d) = %#v, want %#v",
					tt.margin, got, tt.want)
			}
		})
	}
}

func TestCommands_SetPrintAreaWidth(t *testing.T) {
	cmd := printposition.NewCommands()

	tests := []struct {
		name  string
		width uint16
		want  []byte
	}{
		{
			name:  "minimum width",
			width: 1,
			want:  []byte{common.GS, 'W', 0x01, 0x00},
		},
		{
			name:  "typical 58mm width",
			width: 420,
			want:  []byte{common.GS, 'W', 0xA4, 0x01},
		},
		{
			name:  "typical 80mm width",
			width: 576,
			want:  []byte{common.GS, 'W', 0x40, 0x02},
		},
		{
			name:  "maximum width",
			width: 65535,
			want:  []byte{common.GS, 'W', 0xFF, 0xFF},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cmd.SetPrintAreaWidth(tt.width)
			if !bytes.Equal(got, tt.want) {
				t.Errorf("SetPrintAreaWidth(%d) = %#v, want %#v",
					tt.width, got, tt.want)
			}
		})
	}
}

func TestCommands_SetPrintPositionBeginningLine(t *testing.T) {
	cmd := printposition.NewCommands()

	tests := []struct {
		name    string
		mode    printposition.BeginLine
		want    []byte
		wantErr bool
	}{
		{
			name:    "erase buffer",
			mode:    0,
			want:    []byte{common.GS, 'T', 0},
			wantErr: false,
		},
		{
			name:    "print buffer",
			mode:    1,
			want:    []byte{common.GS, 'T', 1},
			wantErr: false,
		},
		{
			name:    "erase buffer ASCII",
			mode:    '0',
			want:    []byte{common.GS, 'T', '0'},
			wantErr: false,
		},
		{
			name:    "print buffer ASCII",
			mode:    '1',
			want:    []byte{common.GS, 'T', '1'},
			wantErr: false,
		},
		{
			name:    "invalid mode",
			mode:    2,
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cmd.SetPrintPositionBeginningLine(tt.mode)

			// Standardized error checking
			if (err != nil) != tt.wantErr {
				t.Errorf("SetPrintPositionBeginningLine(%v) error = %v, wantErr %v",
					tt.mode, err, tt.wantErr)
				return
			}

			// Check specific error type if expecting error
			if tt.wantErr && err != nil {
				if !errors.Is(err, printposition.ErrBeginLineMode) {
					t.Errorf("SetPrintPositionBeginningLine(%v) error = %v, want %v",
						tt.mode, err, printposition.ErrBeginLineMode)
				}
				return
			}

			// Check result if no error expected
			if !tt.wantErr && !bytes.Equal(got, tt.want) {
				t.Errorf("SetPrintPositionBeginningLine(%v) = %#v, want %#v",
					tt.mode, got, tt.want)
			}
		})
	}
}

// ============================================================================
// Page Mode Tests
// ============================================================================

func TestCommands_SelectPrintDirectionPageMode(t *testing.T) {
	cmd := printposition.NewCommands()

	tests := []struct {
		name      string
		direction printposition.PrintDirection
		want      []byte
		wantErr   bool
	}{
		{
			name:      "left to right",
			direction: 0,
			want:      []byte{common.ESC, 'T', 0},
			wantErr:   false,
		},
		{
			name:      "bottom to top",
			direction: 1,
			want:      []byte{common.ESC, 'T', 1},
			wantErr:   false,
		},
		{
			name:      "right to left",
			direction: 2,
			want:      []byte{common.ESC, 'T', 2},
			wantErr:   false,
		},
		{
			name:      "top to bottom",
			direction: 3,
			want:      []byte{common.ESC, 'T', 3},
			wantErr:   false,
		},
		{
			name:      "left to right ASCII",
			direction: '0',
			want:      []byte{common.ESC, 'T', '0'},
			wantErr:   false,
		},
		{
			name:      "bottom to top ASCII",
			direction: '1',
			want:      []byte{common.ESC, 'T', '1'},
			wantErr:   false,
		},
		{
			name:      "right to left ASCII",
			direction: '2',
			want:      []byte{common.ESC, 'T', '2'},
			wantErr:   false,
		},
		{
			name:      "top to bottom ASCII",
			direction: '3',
			want:      []byte{common.ESC, 'T', '3'},
			wantErr:   false,
		},
		{
			name:      "invalid direction",
			direction: 4,
			want:      nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cmd.SelectPrintDirectionPageMode(tt.direction)

			// Standardized error checking
			if (err != nil) != tt.wantErr {
				t.Errorf("SelectPrintDirectionPageMode(%v) error = %v, wantErr %v",
					tt.direction, err, tt.wantErr)
				return
			}

			// Check specific error type if expecting error
			if tt.wantErr && err != nil {
				if !errors.Is(err, printposition.ErrPrintDirection) {
					t.Errorf("SelectPrintDirectionPageMode(%v) error = %v, want %v",
						tt.direction, err, printposition.ErrPrintDirection)
				}
				return
			}

			// Check result if no error expected
			if !tt.wantErr && !bytes.Equal(got, tt.want) {
				t.Errorf("SelectPrintDirectionPageMode(%v) = %#v, want %#v",
					tt.direction, got, tt.want)
			}
		})
	}
}

func TestCommands_SetPrintAreaPageMode(t *testing.T) {
	cmd := printposition.NewCommands()

	tests := []struct {
		name    string
		x       uint16
		y       uint16
		width   uint16
		height  uint16
		want    []byte
		wantErr bool
	}{
		{
			name:    "zero origin, minimum size",
			x:       0,
			y:       0,
			width:   1,
			height:  1,
			want:    []byte{common.ESC, 'W', 0, 0, 0, 0, 1, 0, 1, 0},
			wantErr: false,
		},
		{
			name:    "typical area",
			x:       10,
			y:       20,
			width:   100,
			height:  200,
			want:    []byte{common.ESC, 'W', 10, 0, 20, 0, 100, 0, 200, 0},
			wantErr: false,
		},
		{
			name:    "values requiring high bytes",
			x:       256,
			y:       512,
			width:   1024,
			height:  2048,
			want:    []byte{common.ESC, 'W', 0, 1, 0, 2, 0, 4, 0, 8},
			wantErr: false,
		},
		{
			name:    "maximum values",
			x:       65535,
			y:       65535,
			width:   65535,
			height:  65535,
			want:    []byte{common.ESC, 'W', 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
			wantErr: false,
		},
		{
			name:    "zero width",
			x:       0,
			y:       0,
			width:   0,
			height:  1,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "zero height",
			x:       0,
			y:       0,
			width:   1,
			height:  0,
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cmd.SetPrintAreaPageMode(tt.x, tt.y, tt.width, tt.height)

			if (err != nil) != tt.wantErr {
				t.Errorf("SetPrintAreaPageMode(%d, %d, %d, %d) error = %v, wantErr %v",
					tt.x, tt.y, tt.width, tt.height, err, tt.wantErr)
				return
			}

			var baseErr error
			switch tt.name {
			case "zero width":
				baseErr = printposition.ErrPrintAreaWidthSize
			case "zero height":
				baseErr = printposition.ErrPrintAreaHeightSize
			default:
				baseErr = nil
			}

			if tt.wantErr && err != nil {
				if !errors.Is(err, baseErr) {
					t.Errorf("SetPrintAreaPageMode(%d, %d, %d, %d) error = %v, want %v",
						tt.x, tt.y, tt.width, tt.height, err, baseErr)
				}
				if !errors.Is(err, baseErr) {
					t.Errorf("SetPrintAreaPageMode(%d, %d) error = %v, want %v",
						tt.width, tt.height, err, baseErr)
				}
				return
			}

			if !tt.wantErr && !bytes.Equal(got, tt.want) {
				t.Errorf("BuildCharacterSize(%d, %d) = %v, want %v",
					tt.width, tt.height, got, tt.want)
			}
		})
	}
}

func TestCommands_SetAbsoluteVerticalPrintPosition(t *testing.T) {
	cmd := printposition.NewCommands()

	tests := []struct {
		name     string
		position uint16
		want     []byte
	}{
		{
			name:     "zero position",
			position: 0,
			want:     []byte{common.GS, '$', 0x00, 0x00},
		},
		{
			name:     "position 100",
			position: 100,
			want:     []byte{common.GS, '$', 0x64, 0x00},
		},
		{
			name:     "position 256",
			position: 256,
			want:     []byte{common.GS, '$', 0x00, 0x01},
		},
		{
			name:     "maximum position",
			position: 65535,
			want:     []byte{common.GS, '$', 0xFF, 0xFF},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cmd.SetAbsoluteVerticalPrintPosition(tt.position)
			if !bytes.Equal(got, tt.want) {
				t.Errorf("SetAbsoluteVerticalPrintPosition(%d) = %#v, want %#v",
					tt.position, got, tt.want)
			}
		})
	}
}

func TestCommands_SetRelativeVerticalPrintPosition(t *testing.T) {
	cmd := printposition.NewCommands()

	tests := []struct {
		name     string
		distance int16
		want     []byte
	}{
		{
			name:     "zero distance",
			distance: 0,
			want:     []byte{common.GS, '\\', 0x00, 0x00},
		},
		{
			name:     "positive distance",
			distance: 100,
			want:     []byte{common.GS, '\\', 0x64, 0x00},
		},
		{
			name:     "negative distance",
			distance: -100,
			want:     []byte{common.GS, '\\', 0x9C, 0xFF}, // Two's complement of -100
		},
		{
			name:     "maximum positive",
			distance: 32767,
			want:     []byte{common.GS, '\\', 0xFF, 0x7F},
		},
		{
			name:     "maximum negative",
			distance: -32768,
			want:     []byte{common.GS, '\\', 0x00, 0x80},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cmd.SetRelativeVerticalPrintPosition(tt.distance)
			if !bytes.Equal(got, tt.want) {
				t.Errorf("SetRelativeVerticalPrintPosition(%d) = %#v, want %#v",
					tt.distance, got, tt.want)
			}
		})
	}
}
