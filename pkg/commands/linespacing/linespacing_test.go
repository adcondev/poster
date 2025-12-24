package linespacing_test

import (
	"bytes"
	"testing"

	"github.com/adcondev/poster/pkg/commands/linespacing"
	"github.com/adcondev/poster/pkg/commands/shared"
)

// Naming Convention: Test{Struct}_{Method}_{Scenario}

func TestLineSpacingCommands_SetLineSpacing(t *testing.T) {
	lsc := linespacing.NewCommands()
	tests := []struct {
		name    string
		spacing linespacing.Spacing
		want    []byte
	}{
		{
			"minimum spacing (0 dots)",
			0,
			[]byte{shared.ESC, '3', 0},
		},
		{
			"typical spacing (30 dots)",
			30,
			[]byte{shared.ESC, '3', 30},
		},
		{
			"maximum spacing (255 dots)",
			255,
			[]byte{shared.ESC, '3', 255},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := lsc.SetLineSpacing(tt.spacing)
			if !bytes.Equal(got, tt.want) {
				t.Errorf("SetLineSpacing(%d) = %#v; want %#v", tt.spacing, got, tt.want)
			}
		})
	}
}

func TestLineSpacingCommands_SelectDefaultLineSpacing(t *testing.T) {
	lsc := linespacing.NewCommands()
	got := lsc.SelectDefaultLineSpacing()
	want := []byte{shared.ESC, '2'}
	if !bytes.Equal(got, want) {
		t.Errorf("SelectDefaultLineSpacing() = %#v; want %#v", got, want)
	}
}
