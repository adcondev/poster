package composer_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/adcondev/poster/pkg/commands/barcode"
	"github.com/adcondev/poster/pkg/commands/shared"
	"github.com/adcondev/poster/pkg/composer"
	"github.com/adcondev/poster/pkg/graphics"
)

func TestGenerateBarcode(t *testing.T) {
	c := composer.NewEscpos()

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
				// Expected sequence:
				// 1. GS w n (width)
				// 2. GS h n (height)
				// 3. GS H n (hri pos)
				// 4. GS f n (hri font)
				// 5. GS k m n { B data... (print command for CODE128 with Set B)

				expectedPrefix := []byte{
					shared.GS, 'w', byte(barcode.DefaultWidth),
					shared.GS, 'h', byte(barcode.DefaultHeight),
					shared.GS, 'H', byte(barcode.HRIBelow),
					shared.GS, 'f', byte(barcode.HRIFontA),
				}
				assert.Equal(t, expectedPrefix, output[:len(expectedPrefix)], "Configuration commands mismatch")

				// Print command part
				printPart := output[len(expectedPrefix):]
				// GS k m n ...
				assert.Equal(t, byte(shared.GS), printPart[0])
				assert.Equal(t, byte('k'), printPart[1])
				assert.Equal(t, byte(barcode.CODE128), printPart[2])
				// Length = len("{B") + len("TEST1234") = 2 + 8 = 10
				assert.Equal(t, byte(10), printPart[3])
				// Data part: { B T E S T 1 2 3 4
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
			data: []byte("123456789012"), // 12 digits for EAN13
			validate: func(t *testing.T, output []byte) {
				expectedPrefix := []byte{
					shared.GS, 'w', 2,
					shared.GS, 'h', 50,
					shared.GS, 'H', byte(barcode.HRINotPrinted),
					shared.GS, 'f', byte(barcode.HRIFontB),
				}
				assert.Equal(t, expectedPrefix, output[:len(expectedPrefix)], "Configuration commands mismatch")

				printPart := output[len(expectedPrefix):]
				// GS k m n data...
				assert.Equal(t, byte(shared.GS), printPart[0])
				assert.Equal(t, byte('k'), printPart[1])
				assert.Equal(t, byte(barcode.EAN13), printPart[2])
				assert.Equal(t, byte(12), printPart[3]) // Length
				assert.Equal(t, []byte("123456789012"), printPart[4:], "Barcode data payload mismatch")
			},
		},
		{
			name: "Invalid Width",
			cfg: graphics.BarcodeConfig{
				Symbology: barcode.CODE128,
				Width:     barcode.Width(1), // Invalid (min is 2)
			},
			data:          []byte("TEST"),
			expectedError: "config width",
		},
		{
			name: "Invalid Height",
			cfg: graphics.BarcodeConfig{
				Symbology: barcode.CODE128,
				Width:     barcode.DefaultWidth,
				Height:    barcode.Height(0), // Invalid (min is 1)
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
				HRIPosition: barcode.HRIPosition(99), // Invalid
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
				HRIFont:     barcode.HRIFont(99), // Invalid
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
				CodeSet:     barcode.Code128Set(99), // Invalid
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
