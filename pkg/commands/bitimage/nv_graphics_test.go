package bitimage_test

import (
	"testing"

	"github.com/adcondev/poster/internal/testutils"
	"github.com/adcondev/poster/pkg/commands/bitimage"
	"github.com/adcondev/poster/pkg/commands/shared"
)

// ============================================================================
// NV Graphics Commands Tests
// ============================================================================

func TestNVGraphicsCommands_GetNVGraphicsCapacity(t *testing.T) {
	cmd := bitimage.NewNVGraphicsCommands()

	tests := []struct {
		name    string
		fn      bitimage.NVFunctionCode
		want    []byte
		wantErr error
	}{
		{
			name:    "function code 0",
			fn:      bitimage.NVFuncGetCapacity,
			want:    []byte{shared.GS, '(', 'L', 0x02, 0x00, 0x30, 0},
			wantErr: nil,
		},
		{
			name:    "function code 48",
			fn:      bitimage.NVFuncGetCapacityASCII,
			want:    []byte{shared.GS, '(', 'L', 0x02, 0x00, 0x30, 48},
			wantErr: nil,
		},
		{
			name:    "invalid function code 1",
			fn:      1,
			want:    nil,
			wantErr: bitimage.ErrInvalidNVFunctionCode,
		},
		{
			name:    "invalid function code 99",
			fn:      99,
			want:    nil,
			wantErr: bitimage.ErrInvalidNVFunctionCode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cmd.GetNVGraphicsCapacity(tt.fn)

			if !testutils.AssertErrorOccurred(t, err, tt.wantErr != nil, "GetNVGraphicsCapacity") {
				return
			}
			if tt.wantErr != nil {
				testutils.AssertError(t, err, tt.wantErr)
				return
			}

			testutils.AssertBytes(t, got, tt.want, "GetNVGraphicsCapacity(%v)", tt.fn)
		})
	}
}

func TestNVGraphicsCommands_GetNVGraphicsRemainingCapacity(t *testing.T) {
	cmd := bitimage.NewNVGraphicsCommands()

	tests := []struct {
		name    string
		fn      bitimage.NVFunctionCode
		want    []byte
		wantErr error
	}{
		{
			name:    "function code 3",
			fn:      bitimage.NVFuncGetRemaining,
			want:    []byte{shared.GS, '(', 'L', 0x02, 0x00, 0x30, 3},
			wantErr: nil,
		},
		{
			name:    "function code 51",
			fn:      bitimage.NVFuncGetRemainingASCII,
			want:    []byte{shared.GS, '(', 'L', 0x02, 0x00, 0x30, 51},
			wantErr: nil,
		},
		{
			name:    "invalid function code 0",
			fn:      0,
			want:    nil,
			wantErr: bitimage.ErrInvalidNVFunctionCode,
		},
		{
			name:    "invalid function code 48",
			fn:      48,
			want:    nil,
			wantErr: bitimage.ErrInvalidNVFunctionCode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cmd.GetNVGraphicsRemainingCapacity(tt.fn)

			if !testutils.AssertErrorOccurred(t, err, tt.wantErr != nil, "GetNVGraphicsRemainingCapacity") {
				return
			}
			if tt.wantErr != nil {
				testutils.AssertError(t, err, tt.wantErr)
				return
			}

			testutils.AssertBytes(t, got, tt.want, "GetNVGraphicsRemainingCapacity(%v)", tt.fn)
		})
	}
}

func TestNVGraphicsCommands_GetNVGraphicsKeyCodeList(t *testing.T) {
	cmd := bitimage.NewNVGraphicsCommands()
	want := []byte{shared.GS, '(', 'L', 0x04, 0x00, 0x30, 0x40, 'K', 'C'}

	got := cmd.GetNVGraphicsKeyCodeList()
	testutils.AssertBytes(t, got, want, "GetNVGraphicsKeyCodeList()")
}

func TestNVGraphicsCommands_DeleteAllNVGraphics(t *testing.T) {
	cmd := bitimage.NewNVGraphicsCommands()
	want := []byte{shared.GS, '(', 'L', 0x05, 0x00, 0x30, 0x41, 'C', 'L', 'R'}

	got := cmd.DeleteAllNVGraphics()
	testutils.AssertBytes(t, got, want, "DeleteAllNVGraphics()")
}

func TestNVGraphicsCommands_DeleteNVGraphicsByKeyCode(t *testing.T) {
	cmd := bitimage.NewNVGraphicsCommands()

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
			want:    []byte{shared.GS, '(', 'L', 0x04, 0x00, 0x30, 0x42, 32, 32},
			wantErr: nil,
		},
		{
			name:    "valid key codes typical",
			kc1:     'A',
			kc2:     '1',
			want:    []byte{shared.GS, '(', 'L', 0x04, 0x00, 0x30, 0x42, 'A', '1'},
			wantErr: nil,
		},
		{
			name:    "valid key codes maximum",
			kc1:     126,
			kc2:     126,
			want:    []byte{shared.GS, '(', 'L', 0x04, 0x00, 0x30, 0x42, 126, 126},
			wantErr: nil,
		},
		{
			name:    "invalid kc1 too low",
			kc1:     31,
			kc2:     32,
			want:    nil,
			wantErr: bitimage.ErrInvalidKeyCode,
		},
		{
			name:    "invalid kc2 too high",
			kc1:     32,
			kc2:     127,
			want:    nil,
			wantErr: bitimage.ErrInvalidKeyCode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cmd.DeleteNVGraphicsByKeyCode(tt.kc1, tt.kc2)

			if !testutils.AssertErrorOccurred(t, err, tt.wantErr != nil, "DeleteNVGraphicsByKeyCode") {
				return
			}
			if tt.wantErr != nil {
				testutils.AssertError(t, err, tt.wantErr)
				return
			}

			testutils.AssertBytes(t, got, tt.want, "DeleteNVGraphicsByKeyCode(%v, %v)", tt.kc1, tt.kc2)
		})
	}
}

func TestNVGraphicsCommands_PrintNVGraphics(t *testing.T) {
	cmd := bitimage.NewNVGraphicsCommands()

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
			kc1:             'A',
			kc2:             '1',
			horizontalScale: bitimage.NormalScale,
			verticalScale:   bitimage.NormalScale,
			want:            []byte{shared.GS, '(', 'L', 0x06, 0x00, 0x30, 0x45, 'A', '1', 1, 1},
			wantErr:         nil,
		},
		{
			name:            "double horizontal scale",
			kc1:             'B',
			kc2:             '2',
			horizontalScale: bitimage.DoubleScale,
			verticalScale:   bitimage.NormalScale,
			want:            []byte{shared.GS, '(', 'L', 0x06, 0x00, 0x30, 0x45, 'B', '2', 2, 1},
			wantErr:         nil,
		},
		{
			name:            "double vertical scale",
			kc1:             'C',
			kc2:             '3',
			horizontalScale: bitimage.NormalScale,
			verticalScale:   bitimage.DoubleScale,
			want:            []byte{shared.GS, '(', 'L', 0x06, 0x00, 0x30, 0x45, 'C', '3', 1, 2},
			wantErr:         nil,
		},
		{
			name:            "double both scales",
			kc1:             'D',
			kc2:             '4',
			horizontalScale: bitimage.DoubleScale,
			verticalScale:   bitimage.DoubleScale,
			want:            []byte{shared.GS, '(', 'L', 0x06, 0x00, 0x30, 0x45, 'D', '4', 2, 2},
			wantErr:         nil,
		},
		{
			name:            "invalid key code",
			kc1:             31,
			kc2:             32,
			horizontalScale: bitimage.NormalScale,
			verticalScale:   bitimage.NormalScale,
			want:            nil,
			wantErr:         bitimage.ErrInvalidKeyCode,
		},
		{
			name:            "invalid horizontal scale",
			kc1:             'A',
			kc2:             '1',
			horizontalScale: 0,
			verticalScale:   bitimage.NormalScale,
			want:            nil,
			wantErr:         bitimage.ErrInvalidScale,
		},
		{
			name:            "invalid vertical scale",
			kc1:             'A',
			kc2:             '1',
			horizontalScale: bitimage.NormalScale,
			verticalScale:   3,
			want:            nil,
			wantErr:         bitimage.ErrInvalidScale,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cmd.PrintNVGraphics(tt.kc1, tt.kc2, tt.horizontalScale, tt.verticalScale)

			if !testutils.AssertErrorOccurred(t, err, tt.wantErr != nil, "PrintNVGraphics") {
				return
			}
			if tt.wantErr != nil {
				testutils.AssertError(t, err, tt.wantErr)
				return
			}

			testutils.AssertBytes(t, got, tt.want, "PrintNVGraphics(%v, %v, %v, %v)",
				tt.kc1, tt.kc2, tt.horizontalScale, tt.verticalScale)
		})
	}
}

// ============================================================================
// NV Graphics Validation Functions Tests
// ============================================================================

func TestValidateNVCapacityFunctionCode(t *testing.T) {
	tests := []struct {
		name    string
		fn      bitimage.NVFunctionCode
		wantErr bool
	}{
		{"valid code 0", bitimage.NVFuncGetCapacity, false},
		{"valid code 48", bitimage.NVFuncGetCapacityASCII, false},
		{"invalid code 1", 1, true},
		{"invalid code 3", bitimage.NVFuncGetRemaining, true},
		{"invalid code 51", bitimage.NVFuncGetRemainingASCII, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := bitimage.ValidateNVCapacityFunctionCode(tt.fn)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateNVCapacityFunctionCode(%v) error = %v, wantErr %v", tt.fn, err, tt.wantErr)
			}
		})
	}
}

func TestValidateNVRemainingFunctionCode(t *testing.T) {
	tests := []struct {
		name    string
		fn      bitimage.NVFunctionCode
		wantErr bool
	}{
		{"valid code 3", bitimage.NVFuncGetRemaining, false},
		{"valid code 51", bitimage.NVFuncGetRemainingASCII, false},
		{"invalid code 0", bitimage.NVFuncGetCapacity, true},
		{"invalid code 48", bitimage.NVFuncGetCapacityASCII, true},
		{"invalid code 99", 99, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := bitimage.ValidateNVRemainingFunctionCode(tt.fn)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateNVRemainingFunctionCode(%v) error = %v, wantErr %v", tt.fn, err, tt.wantErr)
			}
		})
	}
}

func TestValidateKeyCode(t *testing.T) {
	tests := []struct {
		name    string
		kc      byte
		wantErr bool
	}{
		{"minimum valid", 32, false},
		{"typical letter", 'A', false},
		{"typical number", '1', false},
		{"maximum valid", 126, false},
		{"too low", 31, true},
		{"too high", 127, true},
		{"zero", 0, true},
		{"max byte", 255, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := bitimage.ValidateKeyCode(tt.kc)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateKeyCode(%v) error = %v, wantErr %v", tt.kc, err, tt.wantErr)
			}
		})
	}
}

func TestValidateNVGraphicsDimensions(t *testing.T) {
	tests := []struct {
		name    string
		width   uint16
		height  uint16
		wantErr bool
	}{
		{"minimum valid", 1, 1, false},
		{"typical dimensions", 500, 500, false},
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
			err := bitimage.ValidateNVGraphicsDimensions(tt.width, tt.height)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateNVGraphicsDimensions(%v, %v) error = %v, wantErr %v",
					tt.width, tt.height, err, tt.wantErr)
			}
		})
	}
}

func TestValidateBMPData(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		wantErr bool
	}{
		{
			name:    "valid BMP header",
			data:    append([]byte{'B', 'M'}, testutils.RepeatByte(52, 0)...),
			wantErr: false,
		},
		{
			name:    "larger valid BMP",
			data:    append([]byte{'B', 'M'}, testutils.RepeatByte(100, 0)...),
			wantErr: false,
		},
		{
			name:    "too small",
			data:    []byte{'B', 'M'},
			wantErr: true,
		},
		{
			name:    "wrong signature",
			data:    append([]byte{'X', 'Y'}, testutils.RepeatByte(52, 0)...),
			wantErr: true,
		},
		{
			name:    "empty data",
			data:    []byte{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := bitimage.ValidateBMPData(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateBMPData(data[%d]) error = %v, wantErr %v",
					len(tt.data), err, tt.wantErr)
			}
		})
	}
}
