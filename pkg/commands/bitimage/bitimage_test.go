package bitimage_test

import (
	"testing"

	"github.com/adcondev/poster/internal/testutils"
	"github.com/adcondev/poster/pkg/commands/bitimage"
	"github.com/adcondev/poster/pkg/commands/common"
)

// ============================================================================
// Legacy Bit Image Commands Tests (SelectBitImageMode)
// ============================================================================

func TestCommands_SelectBitImageMode(t *testing.T) {
	cmd := bitimage.NewCommands()
	prefix := []byte{common.ESC, '*'}

	// Helper to create test data
	createTestData := func(mode bitimage.Mode, width uint16) []byte {
		length := bitimage.CalculateDataLength(mode, width)
		return testutils.RepeatByte(length, 0xFF)
	}

	tests := []struct {
		name    string
		mode    bitimage.Mode
		width   uint16
		data    []byte
		want    []byte
		wantErr error
	}{
		{
			name:    "8-dot single-density mode",
			mode:    bitimage.SingleDensity8,
			width:   10,
			data:    createTestData(bitimage.SingleDensity8, 10),
			want:    append(append(prefix, 0, 10, 0), createTestData(bitimage.SingleDensity8, 10)...),
			wantErr: nil,
		},
		{
			name:    "8-dot double-density mode",
			mode:    bitimage.DoubleDensity8,
			width:   20,
			data:    createTestData(bitimage.DoubleDensity8, 20),
			want:    append(append(prefix, 1, 20, 0), createTestData(bitimage.DoubleDensity8, 20)...),
			wantErr: nil,
		},
		{
			name:    "24-dot single-density mode",
			mode:    bitimage.SingleDensity24,
			width:   15,
			data:    createTestData(bitimage.SingleDensity24, 15),
			want:    append(append(prefix, 32, 15, 0), createTestData(bitimage.SingleDensity24, 15)...),
			wantErr: nil,
		},
		{
			name:    "24-dot double-density mode",
			mode:    bitimage.DoubleDensity24,
			width:   25,
			data:    createTestData(bitimage.DoubleDensity24, 25),
			want:    append(append(prefix, 33, 25, 0), createTestData(bitimage.DoubleDensity24, 25)...),
			wantErr: nil,
		},
		{
			name:    "maximum width",
			mode:    bitimage.SingleDensity8,
			width:   2400,
			data:    createTestData(bitimage.SingleDensity8, 2400),
			want:    append(append(prefix, 0, 0x60, 0x09), createTestData(bitimage.SingleDensity8, 2400)...),
			wantErr: nil,
		},
		{
			name:    "invalid mode",
			mode:    99,
			width:   10,
			data:    []byte{0xFF},
			want:    nil,
			wantErr: bitimage.ErrBitImageMode,
		},
		{
			name:    "width too large",
			mode:    bitimage.SingleDensity8,
			width:   2401,
			data:    []byte{0xFF},
			want:    nil,
			wantErr: bitimage.ErrHorizontalDotsRange,
		},
		{
			name:    "width zero",
			mode:    bitimage.SingleDensity8,
			width:   0,
			data:    []byte{},
			want:    nil,
			wantErr: bitimage.ErrHorizontalDotsRange,
		},
		{
			name:    "data length mismatch for 8-dot",
			mode:    bitimage.SingleDensity8,
			width:   10,
			data:    []byte{0xFF}, // Should be 10 bytes
			want:    nil,
			wantErr: bitimage.ErrDataLength,
		},
		{
			name:    "data length mismatch for 24-dot",
			mode:    bitimage.SingleDensity24,
			width:   10,
			data:    []byte{0xFF}, // Should be 30 bytes
			want:    nil,
			wantErr: bitimage.ErrDataLength,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cmd.SelectBitImageMode(tt.mode, tt.width, tt.data)

			if !testutils.AssertErrorOccurred(t, err, tt.wantErr != nil, "SelectBitImageMode") {
				return
			}
			if tt.wantErr != nil {
				testutils.AssertError(t, err, tt.wantErr)
				return
			}

			testutils.AssertBytes(t, got, tt.want, "SelectBitImageMode(%v, %d, data[%d])",
				tt.mode, tt.width, len(tt.data))
		})
	}
}

// ============================================================================
// Helper Functions Tests
// ============================================================================

func TestCalculateDataLength(t *testing.T) {
	tests := []struct {
		name  string
		mode  bitimage.Mode
		width uint16
		want  int
	}{
		{
			name:  "8-dot single-density",
			mode:  bitimage.SingleDensity8,
			width: 100,
			want:  100,
		},
		{
			name:  "8-dot double-density",
			mode:  bitimage.DoubleDensity8,
			width: 100,
			want:  100,
		},
		{
			name:  "24-dot single-density",
			mode:  bitimage.SingleDensity24,
			width: 100,
			want:  300,
		},
		{
			name:  "24-dot double-density",
			mode:  bitimage.DoubleDensity24,
			width: 100,
			want:  300,
		},
		{
			name:  "invalid mode",
			mode:  99,
			width: 100,
			want:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Using the unexported function through the command
			// We test this indirectly through SelectBitImageMode
			cmd := bitimage.NewCommands()

			// Create data with the expected length
			data := testutils.RepeatByte(tt.want, 0xFF)

			if tt.want > 0 {
				// Should succeed with correct data length
				_, err := cmd.SelectBitImageMode(tt.mode, tt.width, data)
				if tt.mode != 99 && err != nil {
					t.Errorf("CalculateDataLength validation failed for mode %v: %v", tt.mode, err)
				}

				// Should fail with incorrect data length
				if tt.want > 1 {
					wrongData := testutils.RepeatByte(tt.want-1, 0xFF)
					_, err = cmd.SelectBitImageMode(tt.mode, tt.width, wrongData)
					if err == nil {
						t.Errorf("CalculateDataLength should have failed for incorrect data length")
					}
				}
			}
		})
	}
}

// ============================================================================
// Validation Functions Tests
// ============================================================================

func TestValidateBitImageMode(t *testing.T) {
	tests := []struct {
		name    string
		mode    bitimage.Mode
		wantErr bool
	}{
		{"single density 8", bitimage.SingleDensity8, false},
		{"double density 8", bitimage.DoubleDensity8, false},
		{"single density 24", bitimage.SingleDensity24, false},
		{"double density 24", bitimage.DoubleDensity24, false},
		{"invalid mode 2", 2, true},
		{"invalid mode 31", 31, true},
		{"invalid mode 34", 34, true},
		{"invalid mode 99", 99, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := bitimage.ValidateBitImageMode(tt.mode)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateBitImageMode(%v) error = %v, wantErr %v", tt.mode, err, tt.wantErr)
			}
		})
	}
}

func TestValidatePrintMode(t *testing.T) {
	tests := []struct {
		name    string
		mode    bitimage.PrintMode
		wantErr bool
	}{
		{"normal", bitimage.Normal, false},
		{"double width", bitimage.DoubleWidth, false},
		{"double height", bitimage.DoubleHeight, false},
		{"quadruple", bitimage.Quadruple, false},
		{"normal ASCII", bitimage.NormalASCII, false},
		{"double width ASCII", bitimage.DoubleWidthASCII, false},
		{"double height ASCII", bitimage.DoubleHeightASCII, false},
		{"quadruple ASCII", bitimage.QuadrupleASCII, false},
		{"invalid 4", 4, true},
		{"invalid 47", 47, true},
		{"invalid 52", 52, true},
		{"invalid 99", 99, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := bitimage.ValidatePrintMode(tt.mode)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePrintMode(%v) error = %v, wantErr %v", tt.mode, err, tt.wantErr)
			}
		})
	}
}

func TestValidateHorizontalDots(t *testing.T) {
	tests := []struct {
		name    string
		dots    uint16
		wantErr bool
	}{
		{"minimum valid", 1, false},
		{"typical value", 100, false},
		{"high value", 2000, false},
		{"maximum valid", 2400, false},
		{"zero invalid", 0, true},
		{"exceeds maximum", 2401, true},
		{"way over maximum", 5000, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := bitimage.ValidateHorizontalDots(tt.dots)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateHorizontalDots(%v) error = %v, wantErr %v", tt.dots, err, tt.wantErr)
			}
		})
	}
}

func TestValidateNVBitImageNumber(t *testing.T) {
	tests := []struct {
		name    string
		n       byte
		wantErr bool
	}{
		{"minimum valid", 1, false},
		{"typical value", 10, false},
		{"high value", 200, false},
		{"maximum valid", 255, false},
		{"zero invalid", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := bitimage.ValidateNVBitImageNumber(tt.n)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateNVBitImageNumber(%v) error = %v, wantErr %v", tt.n, err, tt.wantErr)
			}
		})
	}
}

func TestValidateNVImageDimensions(t *testing.T) {
	tests := []struct {
		name    string
		width   uint16
		height  uint16
		wantErr bool
	}{
		{"minimum valid", 1, 1, false},
		{"typical dimensions", 100, 100, false},
		{"maximum width", 1023, 100, false},
		{"maximum height", 100, 288, false},
		{"maximum both", 1023, 288, false},
		{"width zero", 0, 100, true},
		{"height zero", 100, 0, true},
		{"width exceeds", 1024, 100, true},
		{"height exceeds", 100, 289, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := bitimage.ValidateNVImageDimensions(tt.width, tt.height)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateNVImageDimensions(%v, %v) error = %v, wantErr %v",
					tt.width, tt.height, err, tt.wantErr)
			}
		})
	}
}

func TestValidateDownloadedImageDimensions(t *testing.T) {
	tests := []struct {
		name    string
		x       byte
		y       byte
		wantErr bool
	}{
		{"minimum valid", 1, 1, false},
		{"typical dimensions", 10, 10, false},
		{"maximum y", 10, 48, false},
		{"maximum product", 32, 48, false}, // 1536
		{"x zero", 0, 10, true},
		{"y zero", 10, 0, true},
		{"y exceeds", 10, 49, true},
		{"product exceeds", 40, 40, true}, // 1600 > 1536
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := bitimage.ValidateDownloadedImageDimensions(tt.x, tt.y)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDownloadedImageDimensions(%v, %v) error = %v, wantErr %v",
					tt.x, tt.y, err, tt.wantErr)
			}
		})
	}
}

func TestValidateVariableImageDimensions(t *testing.T) {
	tests := []struct {
		name    string
		width   uint16
		height  uint16
		mode    bitimage.PrintMode
		wantErr bool
	}{
		{"minimum valid", 1, 1, bitimage.Normal, false},
		{"typical dimensions", 100, 10, bitimage.Normal, false},
		{"maximum width", 4256, 10, bitimage.Normal, false},
		{"maximum height", 100, 16, bitimage.Normal, false},
		{"maximum both", 4256, 16, bitimage.Normal, false},
		{"width zero", 0, 10, bitimage.Normal, true},
		{"height zero", 100, 0, bitimage.Normal, true},
		{"width exceeds", 4257, 10, bitimage.Normal, true},
		{"height exceeds", 100, 17, bitimage.Normal, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := bitimage.ValidateVariableImageDimensions(tt.width, tt.height, tt.mode)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateVariableImageDimensions(%v, %v, %v) error = %v, wantErr %v",
					tt.width, tt.height, tt.mode, err, tt.wantErr)
			}
		})
	}
}

func TestValidateRasterImageDimensions(t *testing.T) {
	tests := []struct {
		name    string
		width   uint16
		height  uint16
		mode    bitimage.PrintMode
		wantErr bool
	}{
		{"minimum valid", 1, 1, bitimage.Normal, false},
		{"typical dimensions", 100, 100, bitimage.Normal, false},
		{"maximum width", 65535, 100, bitimage.Normal, false},
		{"maximum height", 100, 2303, bitimage.Normal, false},
		{"width zero", 0, 100, bitimage.Normal, true},
		{"height zero", 100, 0, bitimage.Normal, true},
		{"height exceeds", 100, 2304, bitimage.Normal, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := bitimage.ValidateRasterImageDimensions(tt.width, tt.height, tt.mode)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateRasterImageDimensions(%v, %v, %v) error = %v, wantErr %v",
					tt.width, tt.height, tt.mode, err, tt.wantErr)
			}
		})
	}
}
