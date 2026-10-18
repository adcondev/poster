package service

import (
	"bytes"
	"testing"

	"github.com/adcondev/poster/pkg/composer"
	"github.com/adcondev/poster/pkg/constants"
	"github.com/adcondev/poster/pkg/graphics"
	"github.com/adcondev/poster/pkg/profile"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockConnector struct {
	buffer bytes.Buffer
}

func (m *mockConnector) Write(p []byte) (n int, err error) {
	return m.buffer.Write(p)
}

func (m *mockConnector) Close() error {
	return nil
}

func TestCalculateMaxPixelWidth(t *testing.T) {
	tests := []struct {
		name     string
		profile  profile.Escpos
		expected int
	}{
		{
			name: "Profile with DotsPerLine set (80mm standard)",
			profile: profile.Escpos{
				PaperWidth:  80,
				DPI:         203,
				DotsPerLine: 576,
			},
			expected: 576,
		},
		{
			name: "Profile with DotsPerLine set (58mm standard)",
			profile: profile.Escpos{
				PaperWidth:  58,
				DPI:         203,
				DotsPerLine: 384,
			},
			expected: 384,
		},
		{
			name: "High DPI profile with DotsPerLine",
			profile: profile.Escpos{
				PaperWidth:  80,
				DPI:         300,
				DotsPerLine: 944,
			},
			expected: 944,
		},
		{
			name: "Profile without DotsPerLine (calculated from PrintWidth)",
			profile: profile.Escpos{
				PaperWidth: 80,
				DPI:        203,
				PrintWidth: 72,
			},
			expected: 575, // 72 * 203 / 25.4 = 575.4 -> 575
		},
		{
			name: "Profile without DotsPerLine or PrintWidth (calc from 80mm PaperWidth)",
			profile: profile.Escpos{
				PaperWidth: 80,
				DPI:        203,
			},
			expected: 575, // 72mm assumed * 203 / 25.4
		},
		{
			name: "Profile without DotsPerLine or PrintWidth (calc from 58mm PaperWidth)",
			profile: profile.Escpos{
				PaperWidth: 58,
				DPI:        203,
			},
			expected: 383, // 48mm assumed * 203 / 25.4 = 383.6 -> 383
		},
		{
			name:     "Profile with no dimension information (fallback to 80mm default)",
			profile:  profile.Escpos{},
			expected: constants.PaperPxWidth80mm,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn := &mockConnector{}
			proto := composer.NewEscpos()
			prof := tt.profile

			printer, err := NewPrinter(proto, &prof, conn)
			require.NoError(t, err)

			actual := printer.calculateMaxPixelWidth()
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestPrintQR_AutomaticOptions(t *testing.T) {
	tests := []struct {
		name             string
		profile          profile.Escpos
		inputOpts        *graphics.QrOptions
		expectedMaxWidth int
	}{
		{
			name: "Standard 80mm 203dpi profile",
			profile: profile.Escpos{
				PaperWidth:  80,
				DPI:         203,
				DotsPerLine: 576,
				HasQR:       true,
			},
			inputOpts:        &graphics.QrOptions{Model: 1},
			expectedMaxWidth: 576,
		},
		{
			name: "Standard 58mm 203dpi profile",
			profile: profile.Escpos{
				PaperWidth:  58,
				DPI:         203,
				DotsPerLine: 384,
				HasQR:       true,
			},
			inputOpts:        &graphics.QrOptions{Model: 1},
			expectedMaxWidth: 384,
		},
		{
			name: "High DPI 80mm profile (300dpi)",
			profile: profile.Escpos{
				PaperWidth:  80,
				DPI:         300,
				DotsPerLine: 944,
				HasQR:       true,
			},
			inputOpts:        &graphics.QrOptions{Model: 1},
			expectedMaxWidth: 944,
		},
		{
			name: "Profile without DotsPerLine (calc from PrintWidth)",
			profile: profile.Escpos{
				PaperWidth: 80,
				DPI:        203,
				PrintWidth: 72,
				HasQR:      true,
			},
			inputOpts:        &graphics.QrOptions{Model: 1},
			expectedMaxWidth: 575,
		},
		{
			name: "Profile without DotsPerLine or PrintWidth (calc from PaperWidth)",
			profile: profile.Escpos{
				PaperWidth: 80,
				DPI:        203,
				HasQR:      true,
			},
			inputOpts:        &graphics.QrOptions{Model: 1},
			expectedMaxWidth: 575,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn := &mockConnector{}
			proto := composer.NewEscpos()
			prof := tt.profile

			printer, err := NewPrinter(proto, &prof, conn)
			require.NoError(t, err)

			err = printer.PrintQR("https://example.com", tt.inputOpts)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedMaxWidth, tt.inputOpts.MaxPixelWidth)
		})
	}
}

func TestPrintQR_NativeVsFallback(t *testing.T) {
	nativeQRHeader := []byte{0x1D, 0x28, 0x6B}
	rasterBitImageHeader := []byte{0x1D, 0x76, 0x30}

	t.Run("Native QR used when HasQR is true and no logo/circle", func(t *testing.T) {
		conn := &mockConnector{}
		proto := composer.NewEscpos()
		prof := profile.Escpos{
			PaperWidth:  80,
			DPI:         203,
			DotsPerLine: 576,
			HasQR:       true,
		}

		printer, err := NewPrinter(proto, &prof, conn)
		require.NoError(t, err)

		opts := graphics.DefaultQROptions()
		err = printer.PrintQR("https://example.com", opts)
		assert.NoError(t, err)

		buf := conn.buffer.Bytes()
		assert.True(t, bytes.Contains(buf, nativeQRHeader), "expected native QR commands in buffer")
		assert.False(t, bytes.Contains(buf, rasterBitImageHeader), "raster bit image commands should not be emitted")
	})

	t.Run("Software fallback used when HasQR is false", func(t *testing.T) {
		conn := &mockConnector{}
		proto := composer.NewEscpos()
		prof := profile.Escpos{
			PaperWidth:  80,
			DPI:         203,
			DotsPerLine: 576,
			HasQR:       false,
		}

		printer, err := NewPrinter(proto, &prof, conn)
		require.NoError(t, err)

		opts := graphics.DefaultQROptions()
		err = printer.PrintQR("https://example.com", opts)
		assert.NoError(t, err)

		buf := conn.buffer.Bytes()
		assert.False(t, bytes.Contains(buf, nativeQRHeader), "native QR commands should not be emitted when HasQR is false")
		assert.True(t, bytes.Contains(buf, rasterBitImageHeader), "expected raster bit image commands in fallback")
	})

	t.Run("Software fallback used when CircleShape is true even with HasQR", func(t *testing.T) {
		conn := &mockConnector{}
		proto := composer.NewEscpos()
		prof := profile.Escpos{
			PaperWidth:  80,
			DPI:         203,
			DotsPerLine: 576,
			HasQR:       true,
		}

		printer, err := NewPrinter(proto, &prof, conn)
		require.NoError(t, err)

		opts := graphics.DefaultQROptions()
		opts.CircleShape = true
		err = printer.PrintQR("https://example.com", opts)
		assert.NoError(t, err)

		buf := conn.buffer.Bytes()
		assert.False(t, bytes.Contains(buf, nativeQRHeader), "native QR commands should not be emitted when CircleShape is true")
		assert.True(t, bytes.Contains(buf, rasterBitImageHeader), "expected raster bit image commands in circle shape fallback")
	})

	t.Run("Default options when nil qropts passed", func(t *testing.T) {
		conn := &mockConnector{}
		proto := composer.NewEscpos()
		prof := profile.Escpos{
			PaperWidth:  80,
			DPI:         203,
			DotsPerLine: 576,
			HasQR:       true,
		}

		printer, err := NewPrinter(proto, &prof, conn)
		require.NoError(t, err)

		err = printer.PrintQR("https://example.com", nil)
		assert.NoError(t, err)

		buf := conn.buffer.Bytes()
		assert.True(t, bytes.Contains(buf, nativeQRHeader), "expected native QR commands for nil options with HasQR")
	})
}
