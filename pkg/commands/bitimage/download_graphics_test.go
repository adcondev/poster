package bitimage_test

import (
	"testing"

	"github.com/adcondev/poster/internal/testutils"
	"github.com/adcondev/poster/pkg/commands/bitimage"
	"github.com/adcondev/poster/pkg/commands/common"
)

// ============================================================================
// Download Graphics Commands Tests
// ============================================================================

func TestDownloadGraphicsCommands_GetDownloadGraphicsRemainingCapacity(t *testing.T) {
	cmd := bitimage.NewDownloadGraphicsCommands()

	tests := []struct {
		name    string
		fn      bitimage.DLFunctionCode
		want    []byte
		wantErr error
	}{
		{
			name:    "function code 4",
			fn:      bitimage.DLFuncGetRemaining,
			want:    []byte{common.GS, '(', 'L', 0x02, 0x00, 0x30, 4},
			wantErr: nil,
		},
		{
			name:    "function code 52 (ASCII)",
			fn:      bitimage.DLFuncGetRemainingASCII,
			want:    []byte{common.GS, '(', 'L', 0x02, 0x00, 0x30, 52},
			wantErr: nil,
		},
		{
			name:    "invalid function code 0",
			fn:      0,
			want:    nil,
			wantErr: bitimage.ErrInvalidDLFunctionCode,
		},
		{
			name:    "invalid function code 99",
			fn:      99,
			want:    nil,
			wantErr: bitimage.ErrInvalidDLFunctionCode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cmd.GetDownloadGraphicsRemainingCapacity(tt.fn)

			if !testutils.AssertErrorOccurred(t, err, tt.wantErr != nil, "GetDownloadGraphicsRemainingCapacity") {
				return
			}
			if tt.wantErr != nil {
				testutils.AssertError(t, err, tt.wantErr)
				return
			}

			testutils.AssertBytes(t, got, tt.want, "GetDownloadGraphicsRemainingCapacity(%v)", tt.fn)
		})
	}
}

func TestDownloadGraphicsCommands_GetDownloadGraphicsKeyCodeList(t *testing.T) {
	cmd := bitimage.NewDownloadGraphicsCommands()
	want := []byte{common.GS, '(', 'L', 0x04, 0x00, 0x30, 0x50, 'K', 'C'}

	got := cmd.GetDownloadGraphicsKeyCodeList()
	testutils.AssertBytes(t, got, want, "GetDownloadGraphicsKeyCodeList()")
}

func TestDownloadGraphicsCommands_DeleteAllDownloadGraphics(t *testing.T) {
	cmd := bitimage.NewDownloadGraphicsCommands()
	want := []byte{common.GS, '(', 'L', 0x05, 0x00, 0x30, 0x51, 'C', 'L', 'R'}

	got := cmd.DeleteAllDownloadGraphics()
	testutils.AssertBytes(t, got, want, "DeleteAllDownloadGraphics()")
}

func TestDownloadGraphicsCommands_DeleteDownloadGraphicsByKeyCode(t *testing.T) {
	cmd := bitimage.NewDownloadGraphicsCommands()

	tests := []struct {
		name    string
		kc1     byte
		kc2     byte
		want    []byte
		wantErr error
	}{
		{
			name:    "valid key codes minimum",
			kc1:     32,
			kc2:     32,
			want:    []byte{common.GS, '(', 'L', 0x04, 0x00, 0x30, 0x52, 32, 32},
			wantErr: nil,
		},
		{
			name:    "valid key codes typical",
			kc1:     'L',
			kc2:     '5',
			want:    []byte{common.GS, '(', 'L', 0x04, 0x00, 0x30, 0x52, 'L', '5'},
			wantErr: nil,
		},
		{
			name:    "valid key codes maximum",
			kc1:     126,
			kc2:     126,
			want:    []byte{common.GS, '(', 'L', 0x04, 0x00, 0x30, 0x52, 126, 126},
			wantErr: nil,
		},
		{
			name:    "invalid kc1",
			kc1:     31,
			kc2:     32,
			want:    nil,
			wantErr: bitimage.ErrInvalidKeyCode,
		},
		{
			name:    "invalid kc2",
			kc1:     32,
			kc2:     127,
			want:    nil,
			wantErr: bitimage.ErrInvalidKeyCode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cmd.DeleteDownloadGraphicsByKeyCode(tt.kc1, tt.kc2)

			if !testutils.AssertErrorOccurred(t, err, tt.wantErr != nil, "DeleteDownloadGraphicsByKeyCode") {
				return
			}
			if tt.wantErr != nil {
				testutils.AssertError(t, err, tt.wantErr)
				return
			}

			testutils.AssertBytes(t, got, tt.want, "DeleteDownloadGraphicsByKeyCode(%v, %v)", tt.kc1, tt.kc2)
		})
	}
}

func TestDownloadGraphicsCommands_PrintDownloadGraphics(t *testing.T) {
	cmd := bitimage.NewDownloadGraphicsCommands()

	tests := []struct {
		name            string
		kc1             byte
		kc2             byte
		horizontalScale bitimage.GraphicsScale
		verticalScale   bitimage.GraphicsScale
		want            []byte
		wantErr         error
	}{
		{
			name:            "normal scale",
			kc1:             'X',
			kc2:             'Y',
			horizontalScale: bitimage.NormalScale,
			verticalScale:   bitimage.NormalScale,
			want:            []byte{common.GS, '(', 'L', 0x06, 0x00, 0x30, 0x55, 'X', 'Y', 1, 1},
			wantErr:         nil,
		},
		{
			name:            "double width",
			kc1:             'A',
			kc2:             'B',
			horizontalScale: bitimage.DoubleScale,
			verticalScale:   bitimage.NormalScale,
			want:            []byte{common.GS, '(', 'L', 0x06, 0x00, 0x30, 0x55, 'A', 'B', 2, 1},
			wantErr:         nil,
		},
		{
			name:            "double height",
			kc1:             'C',
			kc2:             'D',
			horizontalScale: bitimage.NormalScale,
			verticalScale:   bitimage.DoubleScale,
			want:            []byte{common.GS, '(', 'L', 0x06, 0x00, 0x30, 0x55, 'C', 'D', 1, 2},
			wantErr:         nil,
		},
		{
			name:            "quadruple",
			kc1:             'E',
			kc2:             'F',
			horizontalScale: bitimage.DoubleScale,
			verticalScale:   bitimage.DoubleScale,
			want:            []byte{common.GS, '(', 'L', 0x06, 0x00, 0x30, 0x55, 'E', 'F', 2, 2},
			wantErr:         nil,
		},
		{
			name:            "invalid key code 1",
			kc1:             31,
			kc2:             32,
			horizontalScale: bitimage.NormalScale,
			verticalScale:   bitimage.NormalScale,
			want:            nil,
			wantErr:         bitimage.ErrInvalidKeyCode,
		},
		{
			name:            "invalid key code 2",
			kc1:             32,
			kc2:             127,
			horizontalScale: bitimage.NormalScale,
			verticalScale:   bitimage.NormalScale,
			want:            nil,
			wantErr:         bitimage.ErrInvalidKeyCode,
		},
		{
			name:            "invalid horizontal scale",
			kc1:             'A',
			kc2:             'B',
			horizontalScale: 0,
			verticalScale:   bitimage.NormalScale,
			want:            nil,
			wantErr:         bitimage.ErrInvalidScale,
		},
		{
			name:            "invalid vertical scale",
			kc1:             'A',
			kc2:             'B',
			horizontalScale: bitimage.NormalScale,
			verticalScale:   3,
			want:            nil,
			wantErr:         bitimage.ErrInvalidScale,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cmd.PrintDownloadGraphics(tt.kc1, tt.kc2, tt.horizontalScale, tt.verticalScale)

			if !testutils.AssertErrorOccurred(t, err, tt.wantErr != nil, "PrintDownloadGraphics") {
				return
			}
			if tt.wantErr != nil {
				testutils.AssertError(t, err, tt.wantErr)
				return
			}

			testutils.AssertBytes(t, got, tt.want, "PrintDownloadGraphics(%v, %v, %v, %v)",
				tt.kc1, tt.kc2, tt.horizontalScale, tt.verticalScale)
		})
	}
}

// ============================================================================
// Download Graphics Validation Functions Tests
// ============================================================================

func TestValidateDLRemainingFunctionCode(t *testing.T) {
	tests := []struct {
		name    string
		fn      bitimage.DLFunctionCode
		wantErr bool
	}{
		{"valid code 4", bitimage.DLFuncGetRemaining, false},
		{"valid code 52", bitimage.DLFuncGetRemainingASCII, false},
		{"invalid code 0", 0, true},
		{"invalid code 3", 3, true},
		{"invalid code 51", 51, true},
		{"invalid code 99", 99, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := bitimage.ValidateDLRemainingFunctionCode(tt.fn)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDLRemainingFunctionCode(%v) error = %v, wantErr %v", tt.fn, err, tt.wantErr)
			}
		})
	}
}

func TestValidateDLGraphicsDimensions(t *testing.T) {
	tests := []struct {
		name    string
		width   uint16
		height  uint16
		wantErr bool
	}{
		{"minimum valid", 1, 1, false},
		{"typical dimensions", 1000, 1000, false},
		{"maximum width", 8192, 1000, false},
		{"maximum height", 1000, 2304, false},
		{"maximum both", 8192, 2304, false},
		{"width zero", 0, 100, true},
		{"height zero", 100, 0, true},
		{"width exceeds", 8193, 100, true},
		{"height exceeds", 100, 2305, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := bitimage.ValidateDLGraphicsDimensions(tt.width, tt.height)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDLGraphicsDimensions(%v, %v) error = %v, wantErr %v",
					tt.width, tt.height, err, tt.wantErr)
			}
		})
	}
}
