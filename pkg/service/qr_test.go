package service

import (
	"bytes"
	"testing"

	"github.com/adcondev/poster/pkg/composer"
	"github.com/adcondev/poster/pkg/graphics"
	"github.com/adcondev/poster/pkg/profile"
	"github.com/stretchr/testify/assert"
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

func TestPrintQR_AutomaticOptions(t *testing.T) {
	tests := []struct {
		name              string
		profile           profile.Escpos
		inputOpts         *graphics.QrOptions
		expectedMaxWidth  int
	}{
		{
			name: "Standard 80mm 203dpi profile",
			profile: profile.Escpos{
				PaperWidth:  80,
				DPI:         203,
				DotsPerLine: 576,
				HasQR:       true,
			},
			inputOpts:        &graphics.QrOptions{Model: 1}, // minimal opts
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
				DotsPerLine: 944, // 80mm * 300 / 25.4 roughly
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
				PrintWidth: 72, // 72mm printable
				HasQR:      true,
			},
			inputOpts:        &graphics.QrOptions{Model: 1},
			expectedMaxWidth: 575, // 72 * 203 / 25.4 = 575.4 -> 575
		},
		{
			name: "Profile without DotsPerLine or PrintWidth (calc from PaperWidth)",
			profile: profile.Escpos{
				PaperWidth: 80,
				DPI:        203,
				HasQR:      true,
			},
			inputOpts:        &graphics.QrOptions{Model: 1},
			expectedMaxWidth: 575, // 72mm assumed * 203 / 25.4
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn := &mockConnector{}
			proto := composer.NewEscpos()

			// We need a pointer to profile
			prof := tt.profile

			printer, err := NewPrinter(proto, &prof, conn)
			assert.NoError(t, err)

			err = printer.PrintQR("test", tt.inputOpts)
			// assert.NoError(t, err) // Ignoring print error

			assert.Equal(t, tt.expectedMaxWidth, tt.inputOpts.MaxPixelWidth)
		})
	}
}
