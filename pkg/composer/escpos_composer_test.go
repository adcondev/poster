package composer

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/adcondev/poster/pkg/commands/barcode"
	"github.com/adcondev/poster/pkg/commands/character"
	"github.com/adcondev/poster/pkg/commands/printposition"
	"github.com/adcondev/poster/pkg/commands/shared"
	"github.com/adcondev/poster/pkg/graphics"
)

// ============================================================================
// Mocks
// ============================================================================

type mockPrintCapability struct {
	textFn                     func(string) ([]byte, error)
	printAndLineFeedFn         func() []byte
	printAndCarriageReturnFn   func() []byte
	formFeedFn                 func() []byte
	printAndFeedPaperFn        func(byte) []byte
	printAndFeedLinesFn        func(byte) []byte
	printAndReverseFeedFn      func(byte) ([]byte, error)
	printAndReverseFeedLinesFn func(byte) ([]byte, error)
	printDataInPageModeFn      func() []byte
	cancelDataFn               func() []byte
}

func (m *mockPrintCapability) Text(text string) ([]byte, error) {
	if m.textFn != nil {
		return m.textFn(text)
	}
	return nil, nil
}

func (m *mockPrintCapability) PrintAndLineFeed() []byte {
	if m.printAndLineFeedFn != nil {
		return m.printAndLineFeedFn()
	}
	return nil
}

func (m *mockPrintCapability) PrintAndCarriageReturn() []byte {
	if m.printAndCarriageReturnFn != nil {
		return m.printAndCarriageReturnFn()
	}
	return nil
}

func (m *mockPrintCapability) FormFeed() []byte {
	if m.formFeedFn != nil {
		return m.formFeedFn()
	}
	return nil
}

func (m *mockPrintCapability) PrintAndFeedPaper(units byte) []byte {
	if m.printAndFeedPaperFn != nil {
		return m.printAndFeedPaperFn(units)
	}
	return nil
}

func (m *mockPrintCapability) PrintAndFeedLines(lines byte) []byte {
	if m.printAndFeedLinesFn != nil {
		return m.printAndFeedLinesFn(lines)
	}
	return nil
}

func (m *mockPrintCapability) PrintAndReverseFeed(units byte) ([]byte, error) {
	if m.printAndReverseFeedFn != nil {
		return m.printAndReverseFeedFn(units)
	}
	return nil, nil
}

func (m *mockPrintCapability) PrintAndReverseFeedLines(lines byte) ([]byte, error) {
	if m.printAndReverseFeedLinesFn != nil {
		return m.printAndReverseFeedLinesFn(lines)
	}
	return nil, nil
}

func (m *mockPrintCapability) PrintDataInPageMode() []byte {
	if m.printDataInPageModeFn != nil {
		return m.printDataInPageModeFn()
	}
	return nil
}

func (m *mockPrintCapability) CancelData() []byte {
	if m.cancelDataFn != nil {
		return m.cancelDataFn()
	}
	return nil
}

// ============================================================================
// Tests
// ============================================================================

func TestEscposProtocol_PrintLn(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		mockText      func(string) ([]byte, error)
		mockPrintFeed func() []byte
		want          []byte
		wantErr       bool
	}{
		{
			name:  "success",
			input: "Hello World",
			mockText: func(s string) ([]byte, error) {
				return []byte(s), nil
			},
			mockPrintFeed: func() []byte {
				return []byte{0x0A}
			},
			want:    []byte("Hello World\n"),
			wantErr: false,
		},
		{
			name:  "text error",
			input: "Error",
			mockText: func(_ string) ([]byte, error) {
				return nil, errors.New("mock error")
			},
			mockPrintFeed: func() []byte {
				return []byte{0x0A}
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			mock := &mockPrintCapability{
				textFn:             tt.mockText,
				printAndLineFeedFn: tt.mockPrintFeed,
			}

			// Create protocol with mock
			c := &EscposProtocol{
				Print: mock,
			}

			got, err := c.PrintLn(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("PrintLn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !bytes.Equal(got, tt.want) {
				t.Errorf("PrintLn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateBarcode(t *testing.T) {
	c := NewEscpos()

	tests := []struct {
		name          string
		cfg           graphics.BarcodeConfig
		data          []byte
		expectedError string
		validate      func(t *testing.T, output []byte)
	}{
		{
			name: "Valid CODE128 Barcode",
			cfg: graphics.BarcodeConfig{
				Symbology:   barcode.CODE128,
				Width:       barcode.DefaultWidth,
				Height:      barcode.DefaultHeight,
				HRIPosition: barcode.HRIBelow,
				HRIFont:     barcode.HRIFontA,
				CodeSet:     barcode.Code128SetB,
			},
			data: []byte("TEST1234"),
			validate: func(t *testing.T, output []byte) {
				expectedPrefix := []byte{
					shared.GS, 'w', byte(barcode.DefaultWidth),
					shared.GS, 'h', byte(barcode.DefaultHeight),
					shared.GS, 'H', byte(barcode.HRIBelow),
					shared.GS, 'f', byte(barcode.HRIFontA),
				}
				assert.Equal(t, expectedPrefix, output[:len(expectedPrefix)], "Configuration commands mismatch")

				printPart := output[len(expectedPrefix):]
				assert.Equal(t, shared.GS, printPart[0])
				assert.Equal(t, byte('k'), printPart[1])
				assert.Equal(t, byte(barcode.CODE128), printPart[2])
				assert.Equal(t, byte(10), printPart[3])
				expectedData := append([]byte{'{', byte(barcode.Code128SetB)}, []byte("TEST1234")...)
				assert.Equal(t, expectedData, printPart[4:], "Barcode data payload mismatch")
			},
		},
		{
			name: "Valid EAN13 Barcode",
			cfg: graphics.BarcodeConfig{
				Symbology:   barcode.EAN13,
				Width:       barcode.Width(2),
				Height:      barcode.Height(50),
				HRIPosition: barcode.HRINotPrinted,
				HRIFont:     barcode.HRIFontB,
			},
			data: []byte("123456789012"),
			validate: func(t *testing.T, output []byte) {
				expectedPrefix := []byte{
					shared.GS, 'w', 2,
					shared.GS, 'h', 50,
					shared.GS, 'H', byte(barcode.HRINotPrinted),
					shared.GS, 'f', byte(barcode.HRIFontB),
				}
				assert.Equal(t, expectedPrefix, output[:len(expectedPrefix)], "Configuration commands mismatch")

				printPart := output[len(expectedPrefix):]
				assert.Equal(t, shared.GS, printPart[0])
				assert.Equal(t, byte('k'), printPart[1])
				assert.Equal(t, byte(barcode.EAN13), printPart[2])
				assert.Equal(t, byte(12), printPart[3])
				assert.Equal(t, []byte("123456789012"), printPart[4:], "Barcode data payload mismatch")
			},
		},
		{
			name: "Invalid Width",
			cfg: graphics.BarcodeConfig{
				Symbology: barcode.CODE128,
				Width:     barcode.Width(1),
			},
			data:          []byte("TEST"),
			expectedError: "config width",
		},
		{
			name: "Invalid Height",
			cfg: graphics.BarcodeConfig{
				Symbology: barcode.CODE128,
				Width:     barcode.DefaultWidth,
				Height:    barcode.Height(0),
			},
			data:          []byte("TEST"),
			expectedError: "config height",
		},
		{
			name: "Invalid HRI Position",
			cfg: graphics.BarcodeConfig{
				Symbology:   barcode.CODE128,
				Width:       barcode.DefaultWidth,
				Height:      barcode.DefaultHeight,
				HRIPosition: barcode.HRIPosition(99),
			},
			data:          []byte("TEST"),
			expectedError: "config HRI pos",
		},
		{
			name: "Invalid HRI Font",
			cfg: graphics.BarcodeConfig{
				Symbology:   barcode.CODE128,
				Width:       barcode.DefaultWidth,
				Height:      barcode.DefaultHeight,
				HRIPosition: barcode.HRIBelow,
				HRIFont:     barcode.HRIFont(99),
			},
			data:          []byte("TEST"),
			expectedError: "config HRI font",
		},
		{
			name: "Empty Data",
			cfg: graphics.BarcodeConfig{
				Symbology:   barcode.EAN13,
				Width:       barcode.DefaultWidth,
				Height:      barcode.DefaultHeight,
				HRIPosition: barcode.HRIBelow,
				HRIFont:     barcode.HRIFontA,
			},
			data:          []byte{},
			expectedError: "generate payload",
		},
		{
			name: "Invalid CodeSet for CODE128",
			cfg: graphics.BarcodeConfig{
				Symbology:   barcode.CODE128,
				Width:       barcode.DefaultWidth,
				Height:      barcode.DefaultHeight,
				HRIPosition: barcode.HRIBelow,
				HRIFont:     barcode.HRIFontA,
				CodeSet:     barcode.Code128Set(99),
			},
			data:          []byte("TEST"),
			expectedError: "generate payload",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := c.GenerateBarcode(tt.cfg, tt.data)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, output)
			} else {
				require.NoError(t, err)
				require.NotNil(t, output)
				if tt.validate != nil {
					tt.validate(t, output)
				}
			}
		})
	}
}

func TestEscposProtocol_InitializePrinter(t *testing.T) {
	c := NewEscpos()
	require.NotNil(t, c)
	assert.Equal(t, []byte{shared.ESC, '@'}, c.InitializePrinter())
}

func TestEscposProtocol_MarginsAndWidth(t *testing.T) {
	c := NewEscpos()

	marginCmd := c.LeftMargin(16)
	expectedMargin := []byte{shared.GS, 'L', 16, 0}
	assert.Equal(t, expectedMargin, marginCmd)

	widthCmd := c.PrintWidth(384)
	expectedWidth := []byte{shared.GS, 'W', byte(384 & 0xFF), byte(384 >> 8)}
	assert.Equal(t, expectedWidth, widthCmd)
}

func TestEscposProtocol_Alignment(t *testing.T) {
	c := NewEscpos()

	t.Run("SetAlign valid", func(t *testing.T) {
		cmd, err := c.SetAlign(printposition.Center)
		require.NoError(t, err)
		assert.Equal(t, []byte{shared.ESC, 'a', byte(printposition.Center)}, cmd)
	})

	t.Run("SetAlign invalid", func(t *testing.T) {
		cmd, err := c.SetAlign(printposition.Justification(99))
		require.Error(t, err)
		assert.Nil(t, cmd)
	})

	t.Run("LeftAlign", func(t *testing.T) {
		assert.Equal(t, []byte{shared.ESC, 'a', byte(printposition.Left)}, c.LeftAlign())
	})

	t.Run("CenterAlign", func(t *testing.T) {
		assert.Equal(t, []byte{shared.ESC, 'a', byte(printposition.Center)}, c.CenterAlign())
	})

	t.Run("RightAlign", func(t *testing.T) {
		assert.Equal(t, []byte{shared.ESC, 'a', byte(printposition.Right)}, c.RightAlign())
	})
}

func TestEscposProtocol_TextStyling(t *testing.T) {
	c := NewEscpos()

	assert.Equal(t, []byte{shared.ESC, 'E', 1}, c.EnableBold())
	assert.Equal(t, []byte{shared.ESC, 'E', 0}, c.DisableBold())

	assert.Equal(t, []byte{shared.ESC, '-', 1}, c.OneDotUnderline())
	assert.Equal(t, []byte{shared.ESC, '-', 2}, c.TwoDotUnderline())
	assert.Equal(t, []byte{shared.ESC, '-', 0}, c.DisableUnderline())

	assert.Equal(t, []byte{shared.GS, 'B', 1}, c.EnableReverseMode())
	assert.Equal(t, []byte{shared.GS, 'B', 0}, c.DisableReverseMode())

	assert.Equal(t, []byte{shared.ESC, 'M', byte(character.FontA)}, c.SetFontA())
	assert.Equal(t, []byte{shared.ESC, 'M', byte(character.FontB)}, c.SetFontB())
}

func TestEscposProtocol_TextSizing(t *testing.T) {
	c := NewEscpos()

	assert.Equal(t, []byte{shared.GS, '!', 0x00}, c.SingleSizeText())
	assert.Equal(t, []byte{shared.GS, '!', 0x11}, c.DoubleSizeText())
	assert.Equal(t, []byte{shared.GS, '!', 0x22}, c.TripleSizeText())
	assert.Equal(t, []byte{shared.GS, '!', 0x33}, c.QuadraSizeText())
	assert.Equal(t, []byte{shared.GS, '!', 0x44}, c.PentaSizeText())
	assert.Equal(t, []byte{shared.GS, '!', 0x55}, c.HexaSizeText())
	assert.Equal(t, []byte{shared.GS, '!', 0x66}, c.HeptaSizeText())
	assert.Equal(t, []byte{shared.GS, '!', 0x77}, c.OctaSizeText())
}

func TestEscposProtocol_CustomSizeText(t *testing.T) {
	c := NewEscpos()

	t.Run("clamped below 1", func(t *testing.T) {
		assert.Equal(t, c.SingleSizeText(), c.CustomSizeText(0, 0))
	})

	t.Run("clamped above 8", func(t *testing.T) {
		assert.Equal(t, c.OctaSizeText(), c.CustomSizeText(9, 10))
	})

	t.Run("custom valid size", func(t *testing.T) {
		expected := []byte{shared.GS, '!', byte((2-1)<<4 | (4 - 1))}
		assert.Equal(t, expected, c.CustomSizeText(2, 4))
	})
}

func TestEscposProtocol_LineSpacing(t *testing.T) {
	c := NewEscpos()

	assert.Equal(t, []byte{shared.ESC, '3', 30}, c.SetLineSpacing(30))
	assert.Equal(t, []byte{shared.ESC, '2'}, c.SelectDefaultLineSpacing())
}

func TestEscposProtocol_MechanismControl(t *testing.T) {
	c := NewEscpos()

	pulseCmd, err := c.Pulse(0, 10, 10)
	require.NoError(t, err)
	assert.NotEmpty(t, pulseCmd)

	beepCmd, err := c.Beep(2, 5)
	require.NoError(t, err)
	assert.NotEmpty(t, beepCmd)
}

