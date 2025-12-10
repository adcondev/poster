package qrcode_test

import (
	"testing"

	"github.com/adcondev/poster/internal/testutils"
	"github.com/adcondev/poster/pkg/commands/qrcode"
)

// ============================================================================
// Test Data
// ============================================================================

var (
	numericData      = []byte("123456789")
	alphanumericData = []byte("HELLO123")
	mixedData        = []byte("Hello World!")
	urlData          = []byte("https://example.com")
	maxData          = testutils.RepeatByte(7089, 'A')
	overMaxData      = testutils.RepeatByte(7090, 'B')
)

// ============================================================================
// SelectQRCodeModel Tests
// ============================================================================

func TestCommands_SelectQRCodeModel(t *testing.T) {
	cmd := qrcode.NewCommands()
	prefix := []byte{0x1D, '(', 'k', 0x04, 0x00, 0x31, 0x41}

	tests := []struct {
		name    string
		n1      qrcode.Model
		n2      byte
		want    []byte
		wantErr error
	}{
		{
			name:    "Model 1",
			n1:      qrcode.Model1,
			n2:      0,
			want:    append(prefix, 49, 0),
			wantErr: nil,
		},
		{
			name:    "Model 2 (default)",
			n1:      qrcode.Model2,
			n2:      0,
			want:    append(prefix, 50, 0),
			wantErr: nil,
		},
		{
			name:    "Micro QR",
			n1:      qrcode.MicroQR,
			n2:      0,
			want:    append(prefix, 51, 0),
			wantErr: nil,
		},
		{
			name:    "invalid model 48",
			n1:      48,
			n2:      0,
			want:    nil,
			wantErr: qrcode.ErrQRModel,
		},
		{
			name:    "invalid model 52",
			n1:      52,
			n2:      0,
			want:    nil,
			wantErr: qrcode.ErrQRModel,
		},
		{
			name:    "invalid n2 parameter",
			n1:      qrcode.Model2,
			n2:      1,
			want:    nil,
			wantErr: qrcode.ErrParameter,
		},
		{
			name:    "invalid n2 parameter 255",
			n1:      qrcode.Model2,
			n2:      255,
			want:    nil,
			wantErr: qrcode.ErrParameter,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cmd.SelectQRCodeModel(tt.n1, tt.n2)

			if !testutils.AssertErrorOccurred(t, err, tt.wantErr != nil, "SelectQRCodeModel") {
				return
			}
			if tt.wantErr != nil {
				testutils.AssertError(t, err, tt.wantErr)
				return
			}

			testutils.AssertBytes(t, got, tt.want, "SelectQRCodeModel(%v, %v)", tt.n1, tt.n2)
		})
	}
}

// ============================================================================
// SetQRCodeModuleSize Tests
// ============================================================================

func TestCommands_SetQRCodeModuleSize(t *testing.T) {
	cmd := qrcode.NewCommands()
	prefix := []byte{0x1D, '(', 'k', 0x03, 0x00, 0x31, 0x43}

	tests := []struct {
		name    string
		size    qrcode.ModuleSize
		want    []byte
		wantErr error
	}{
		{
			name:    "minimum size 1",
			size:    qrcode.MinModuleSize,
			want:    append(prefix, 1),
			wantErr: nil,
		},
		{
			name:    "default size 3",
			size:    qrcode.DefaultModuleSize,
			want:    append(prefix, 3),
			wantErr: nil,
		},
		{
			name:    "medium size 8",
			size:    8,
			want:    append(prefix, 8),
			wantErr: nil,
		},
		{
			name:    "maximum size 16",
			size:    qrcode.MaxModuleSize,
			want:    append(prefix, 16),
			wantErr: nil,
		},
		{
			name:    "invalid size 0",
			size:    0,
			want:    nil,
			wantErr: qrcode.ErrModuleSize,
		},
		{
			name:    "invalid size 17",
			size:    17,
			want:    nil,
			wantErr: qrcode.ErrModuleSize,
		},
		{
			name:    "invalid size 255",
			size:    255,
			want:    nil,
			wantErr: qrcode.ErrModuleSize,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cmd.SetQRCodeModuleSize(tt.size)

			if !testutils.AssertErrorOccurred(t, err, tt.wantErr != nil, "SetQRCodeModuleSize") {
				return
			}
			if tt.wantErr != nil {
				testutils.AssertError(t, err, tt.wantErr)
				return
			}

			testutils.AssertBytes(t, got, tt.want, "SetQRCodeModuleSize(%v)", tt.size)
		})
	}
}

// ============================================================================
// SetQRCodeErrorCorrectionLevel Tests
// ============================================================================

func TestCommands_SetQRCodeErrorCorrectionLevel(t *testing.T) {
	cmd := qrcode.NewCommands()
	prefix := []byte{0x1D, '(', 'k', 0x03, 0x00, 0x31, 0x45}

	tests := []struct {
		name    string
		level   qrcode.ErrorCorrection
		want    []byte
		wantErr error
	}{
		{
			name:    "Level L (7% recovery)",
			level:   qrcode.LevelL,
			want:    append(prefix, 48),
			wantErr: nil,
		},
		{
			name:    "Level M (15% recovery)",
			level:   qrcode.LevelM,
			want:    append(prefix, 49),
			wantErr: nil,
		},
		{
			name:    "Level Q (25% recovery)",
			level:   qrcode.LevelQ,
			want:    append(prefix, 50),
			wantErr: nil,
		},
		{
			name:    "Level H (30% recovery)",
			level:   qrcode.LevelH,
			want:    append(prefix, 51),
			wantErr: nil,
		},
		{
			name:    "invalid level 47",
			level:   47,
			want:    nil,
			wantErr: qrcode.ErrErrorCorrection,
		},
		{
			name:    "invalid level 52",
			level:   52,
			want:    nil,
			wantErr: qrcode.ErrErrorCorrection,
		},
		{
			name:    "invalid level 255",
			level:   255,
			want:    nil,
			wantErr: qrcode.ErrErrorCorrection,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cmd.SetQRCodeErrorCorrectionLevel(tt.level)

			if !testutils.AssertErrorOccurred(t, err, tt.wantErr != nil, "SetQRCodeErrorCorrectionLevel") {
				return
			}
			if tt.wantErr != nil {
				testutils.AssertError(t, err, tt.wantErr)
				return
			}

			testutils.AssertBytes(t, got, tt.want, "SetQRCodeErrorCorrectionLevel(%v)", tt.level)
		})
	}
}

// ============================================================================
// StoreQRCodeData Tests
// ============================================================================

func TestCommands_StoreQRCodeData(t *testing.T) {
	cmd := qrcode.NewCommands()

	tests := []struct {
		name    string
		data    []byte
		wantErr error
	}{
		{
			name:    "numeric data",
			data:    numericData,
			wantErr: nil,
		},
		{
			name:    "alphanumeric data",
			data:    alphanumericData,
			wantErr: nil,
		},
		{
			name:    "mixed data",
			data:    mixedData,
			wantErr: nil,
		},
		{
			name:    "URL data",
			data:    urlData,
			wantErr: nil,
		},
		{
			name:    "minimum data (1 byte)",
			data:    []byte("A"),
			wantErr: nil,
		},
		{
			name:    "maximum data (7089 bytes)",
			data:    maxData,
			wantErr: nil,
		},
		{
			name:    "empty data",
			data:    []byte{},
			wantErr: qrcode.ErrDataTooShort,
		},
		{
			name:    "data too long (7090 bytes)",
			data:    overMaxData,
			wantErr: qrcode.ErrDataTooLong,
		},
		{
			name:    "binary data",
			data:    []byte{0x00, 0xFF, 0x7F, 0x80},
			wantErr: nil,
		},
		{
			name:    "special characters",
			data:    []byte("$%*+-./:"),
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cmd.StoreQRCodeData(tt.data)

			if !testutils.AssertErrorOccurred(t, err, tt.wantErr != nil, "StoreQRCodeData") {
				return
			}
			if tt.wantErr != nil {
				testutils.AssertError(t, err, tt.wantErr)
				return
			}

			// Verify command structure
			testutils.AssertHasPrefix(t, got, []byte{0x1D, '(', 'k'}, "StoreQRCodeData prefix")

			// Verify length bytes
			if len(got) >= 5 {
				totalLen := int(got[3]) | (int(got[4]) << 8)
				expectedLen := 3 + len(tt.data)
				if totalLen != expectedLen {
					t.Errorf("incorrect length bytes: got %d, want %d", totalLen, expectedLen)
				}
			}

			// Verify function parameters
			if len(got) >= 8 {
				if got[5] != 0x31 || got[6] != byte(80) || got[7] != 0x30 {
					t.Errorf("incorrect function parameters")
				}
			}

			// Verify data is appended
			if len(tt.data) > 0 {
				testutils.AssertContains(t, got, tt.data, "data should be in command")
			}
		})
	}
}

// ============================================================================
// PrintQRCode Tests
// ============================================================================

func TestCommands_PrintQRCode(t *testing.T) {
	cmd := qrcode.NewCommands()

	got := cmd.PrintQRCode()
	want := []byte{0x1D, '(', 'k', 0x03, 0x00, 0x31, 0x51, 0x30}

	testutils.AssertBytes(t, got, want, "PrintQRCode()")
}

// ============================================================================
// GetQRCodeSize Tests
// ============================================================================

func TestCommands_GetQRCodeSize(t *testing.T) {
	cmd := qrcode.NewCommands()

	got := cmd.GetQRCodeSize()
	want := []byte{0x1D, '(', 'k', 0x03, 0x00, 0x31, 0x52, 0x30}

	testutils.AssertBytes(t, got, want, "GetQRCodeSize()")
}

// ============================================================================
// Validation Helper Tests
// ============================================================================

func TestValidateQRModel(t *testing.T) {
	tests := []struct {
		name    string
		model   qrcode.Model
		wantErr error
	}{
		{
			name:    "valid Model1",
			model:   qrcode.Model1,
			wantErr: nil,
		},
		{
			name:    "valid Model2",
			model:   qrcode.Model2,
			wantErr: nil,
		},
		{
			name:    "valid MicroQR",
			model:   qrcode.MicroQR,
			wantErr: nil,
		},
		{
			name:    "invalid model 48",
			model:   48,
			wantErr: qrcode.ErrQRModel,
		},
		{
			name:    "invalid model 52",
			model:   52,
			wantErr: qrcode.ErrQRModel,
		},
		{
			name:    "invalid model 0",
			model:   0,
			wantErr: qrcode.ErrQRModel,
		},
		{
			name:    "invalid model 255",
			model:   255,
			wantErr: qrcode.ErrQRModel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := qrcode.ValidateQRModel(tt.model)
			testutils.AssertError(t, err, tt.wantErr)
		})
	}
}

func TestValidateModuleSize(t *testing.T) {
	tests := []struct {
		name    string
		size    qrcode.ModuleSize
		wantErr error
	}{
		{
			name:    "valid minimum 1",
			size:    1,
			wantErr: nil,
		},
		{
			name:    "valid default 3",
			size:    3,
			wantErr: nil,
		},
		{
			name:    "valid middle 8",
			size:    8,
			wantErr: nil,
		},
		{
			name:    "valid maximum 16",
			size:    16,
			wantErr: nil,
		},
		{
			name:    "invalid 0",
			size:    0,
			wantErr: qrcode.ErrModuleSize,
		},
		{
			name:    "invalid 17",
			size:    17,
			wantErr: qrcode.ErrModuleSize,
		},
		{
			name:    "invalid 100",
			size:    100,
			wantErr: qrcode.ErrModuleSize,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := qrcode.ValidateModuleSize(tt.size)
			testutils.AssertError(t, err, tt.wantErr)
		})
	}
}

func TestValidateErrorCorrection(t *testing.T) {
	tests := []struct {
		name    string
		level   qrcode.ErrorCorrection
		wantErr error
	}{
		{
			name:    "valid L",
			level:   qrcode.LevelL,
			wantErr: nil,
		},
		{
			name:    "valid M",
			level:   qrcode.LevelM,
			wantErr: nil,
		},
		{
			name:    "valid Q",
			level:   qrcode.LevelQ,
			wantErr: nil,
		},
		{
			name:    "valid H",
			level:   qrcode.LevelH,
			wantErr: nil,
		},
		{
			name:    "invalid 47",
			level:   47,
			wantErr: qrcode.ErrErrorCorrection,
		},
		{
			name:    "invalid 52",
			level:   52,
			wantErr: qrcode.ErrErrorCorrection,
		},
		{
			name:    "invalid 0",
			level:   0,
			wantErr: qrcode.ErrErrorCorrection,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := qrcode.ValidateErrorCorrection(tt.level)
			testutils.AssertError(t, err, tt.wantErr)
		})
	}
}

func TestValidateDataLength(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		wantErr error
	}{
		{
			name:    "valid minimum 1 byte",
			data:    []byte("A"),
			wantErr: nil,
		},
		{
			name:    "valid small data",
			data:    []byte("Hello"),
			wantErr: nil,
		},
		{
			name:    "valid large data 1000 bytes",
			data:    testutils.RepeatByte(1000, 'X'),
			wantErr: nil,
		},
		{
			name:    "valid maximum 7089 bytes",
			data:    testutils.RepeatByte(7089, 'Y'),
			wantErr: nil,
		},
		{
			name:    "invalid empty",
			data:    []byte{},
			wantErr: qrcode.ErrDataTooShort,
		},
		{
			name:    "invalid 7090 bytes",
			data:    testutils.RepeatByte(7090, 'Z'),
			wantErr: qrcode.ErrDataTooLong,
		},
		{
			name:    "invalid 10000 bytes",
			data:    testutils.RepeatByte(10000, 'W'),
			wantErr: qrcode.ErrDataTooLong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := qrcode.ValidateDataLength(tt.data)
			testutils.AssertError(t, err, tt.wantErr)
		})
	}
}

// ============================================================================
// Data Analysis Helper Tests
// ============================================================================

func TestIsNumericData(t *testing.T) {
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := qrcode.IsNumericData(tt.data)
			if got != tt.want {
				t.Errorf("IsNumericData(%q) = %v, want %v", tt.data, got, tt.want)
			}
		})
	}
}

func TestIsAlphanumericData(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		want bool
	}{
		{
			name: "numeric only",
			data: []byte("123456"),
			want: true,
		},
		{
			name: "uppercase letters",
			data: []byte("ABCDEF"),
			want: true,
		},
		{
			name: "mixed numeric and uppercase",
			data: []byte("ABC123"),
			want: true,
		},
		{
			name: "with space",
			data: []byte("ABC 123"),
			want: true,
		},
		{
			name: "with allowed special chars",
			data: []byte("$%*+-./:"),
			want: true,
		},
		{
			name: "lowercase letters",
			data: []byte("abc"),
			want: false,
		},
		{
			name: "mixed case",
			data: []byte("AbC123"),
			want: false,
		},
		{
			name: "invalid special char",
			data: []byte("ABC@123"),
			want: false,
		},
		{
			name: "invalid char #",
			data: []byte("ABC#123"),
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
			got := qrcode.IsAlphanumericData(tt.data)
			if got != tt.want {
				t.Errorf("IsAlphanumericData(%q) = %v, want %v", tt.data, got, tt.want)
			}
		})
	}
}

// ============================================================================
// Boundary Tests
// ============================================================================

func TestCommands_BoundaryValues(t *testing.T) {
	cmd := qrcode.NewCommands()

	t.Run("module size boundaries", func(t *testing.T) {
		// Test all valid module sizes
		for size := qrcode.ModuleSize(1); size <= 16; size++ {
			_, err := cmd.SetQRCodeModuleSize(size)
			if err != nil {
				t.Errorf("SetQRCodeModuleSize(%d) should be valid: %v", size, err)
			}
		}

		// Test invalid boundaries
		invalidSizes := []qrcode.ModuleSize{0, 17, 255}
		for _, size := range invalidSizes {
			_, err := cmd.SetQRCodeModuleSize(size)
			if err == nil {
				t.Errorf("SetQRCodeModuleSize(%d) should be invalid", size)
			}
		}
	})

	t.Run("data length boundaries", func(t *testing.T) {
		// Minimum valid data
		minData := []byte("A")
		_, err := cmd.StoreQRCodeData(minData)
		testutils.AssertError(t, err, nil)

		// Maximum valid data
		maxValidData := testutils.RepeatByte(7089, 'X')
		_, err = cmd.StoreQRCodeData(maxValidData)
		testutils.AssertError(t, err, nil)

		// Just over maximum
		overMaxData := testutils.RepeatByte(7090, 'Y')
		_, err = cmd.StoreQRCodeData(overMaxData)
		testutils.AssertError(t, err, qrcode.ErrDataTooLong)
	})

	t.Run("all models", func(t *testing.T) {
		models := []qrcode.Model{
			qrcode.Model1,
			qrcode.Model2,
			qrcode.MicroQR,
		}

		for _, model := range models {
			_, err := cmd.SelectQRCodeModel(model, 0)
			if err != nil {
				t.Errorf("SelectQRCodeModel(%d, 0) should be valid: %v", model, err)
			}
		}
	})

	t.Run("all error correction levels", func(t *testing.T) {
		levels := []qrcode.ErrorCorrection{
			qrcode.LevelL,
			qrcode.LevelM,
			qrcode.LevelQ,
			qrcode.LevelH,
		}

		for _, level := range levels {
			_, err := cmd.SetQRCodeErrorCorrectionLevel(level)
			if err != nil {
				t.Errorf("SetQRCodeErrorCorrectionLevel(%d) should be valid: %v", level, err)
			}
		}
	})
}

// ============================================================================
// Command Structure Tests
// ============================================================================

func TestCommands_CommandStructure(t *testing.T) {
	cmd := qrcode.NewCommands()

	t.Run("verify length encoding", func(t *testing.T) {
		// Test various data lengths to verify pL/pH encoding
		testCases := []struct {
			dataLen     int
			expectedPL  byte
			expectedPH  byte
			description string
		}{
			{1, 0x04, 0x00, "minimum data"},
			{100, 0x67, 0x00, "small data"},
			{255, 0x02, 0x01, "boundary at 255"},
			{256, 0x03, 0x01, "boundary at 256"},
			{1000, 0xEB, 0x03, "medium data"},
			{7089, 0xB4, 0x1B, "maximum data"},
		}

		for _, tc := range testCases {
			t.Run(tc.description, func(t *testing.T) {
				data := testutils.RepeatByte(tc.dataLen, 'T')
				result, err := cmd.StoreQRCodeData(data)
				testutils.AssertError(t, err, nil)

				if len(result) >= 5 {
					if result[3] != tc.expectedPL || result[4] != tc.expectedPH {
						t.Errorf("incorrect length bytes for %d bytes: got pL=%#x pH=%#x, want pL=%#x pH=%#x",
							tc.dataLen, result[3], result[4], tc.expectedPL, tc.expectedPH)
					}
				}
			})
		}
	})

	t.Run("fixed length commands", func(t *testing.T) {
		// Commands with fixed 3-byte parameter length
		fixedLengthCommands := []struct {
			name string
			cmd  []byte
			fn   byte
		}{
			{"SetModuleSize", func() []byte { c, _ := cmd.SetQRCodeModuleSize(3); return c }(), 0x43},
			{"SetErrorCorrection", func() []byte { c, _ := cmd.SetQRCodeErrorCorrectionLevel(qrcode.LevelM); return c }(), 0x45},
			{"Print", cmd.PrintQRCode(), 0x51},
			{"GetSize", cmd.GetQRCodeSize(), 0x52},
		}

		for _, tc := range fixedLengthCommands {
			t.Run(tc.name, func(t *testing.T) {
				if len(tc.cmd) < 7 {
					t.Errorf("%s command too short: %d bytes", tc.name, len(tc.cmd))
					return
				}
				// Check pL=3, pH=0
				if tc.cmd[3] != 0x03 || tc.cmd[4] != 0x00 {
					t.Errorf("%s incorrect length bytes: pL=%#x pH=%#x", tc.name, tc.cmd[3], tc.cmd[4])
				}
				// Check cn=49
				if tc.cmd[5] != 0x31 {
					t.Errorf("%s incorrect cn: %#x", tc.name, tc.cmd[5])
				}
				// Check function number
				if tc.cmd[6] != tc.fn {
					t.Errorf("%s incorrect fn: got %#x, want %#x", tc.name, tc.cmd[6], tc.fn)
				}
			})
		}
	})

	t.Run("model command has 4-byte length", func(t *testing.T) {
		modelCmd, err := cmd.SelectQRCodeModel(qrcode.Model2, 0)
		testutils.AssertError(t, err, nil)

		if len(modelCmd) < 9 {
			t.Errorf("model command too short: %d bytes", len(modelCmd))
			return
		}

		// Check pL=4, pH=0
		if modelCmd[3] != 0x04 || modelCmd[4] != 0x00 {
			t.Errorf("model command incorrect length bytes: pL=%#x pH=%#x", modelCmd[3], modelCmd[4])
		}
	})
}

// ============================================================================
// Error Handling Tests
// ============================================================================

func TestCommands_ErrorMessages(t *testing.T) {
	cmd := qrcode.NewCommands()

	t.Run("error messages contain values", func(t *testing.T) {
		// Invalid model
		_, err := cmd.SelectQRCodeModel(99, 0)
		if err == nil || !testutils.ContainsAny([]byte(err.Error()), []byte("99")) {
			t.Error("error message should contain invalid value")
		}

		// Invalid n2 parameter
		_, err = cmd.SelectQRCodeModel(qrcode.Model2, 5)
		if err == nil || !testutils.ContainsAny([]byte(err.Error()), []byte("5")) {
			t.Error("error message should contain invalid n2 value")
		}

		// Invalid module size
		_, err = cmd.SetQRCodeModuleSize(20)
		if err == nil || !testutils.ContainsAny([]byte(err.Error()), []byte("20")) {
			t.Error("error message should contain invalid size")
		}

		// Invalid error correction
		_, err = cmd.SetQRCodeErrorCorrectionLevel(55)
		if err == nil || !testutils.ContainsAny([]byte(err.Error()), []byte("55")) {
			t.Error("error message should contain invalid level")
		}
	})
}
