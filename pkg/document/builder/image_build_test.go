package builder

import (
	"encoding/json"
	"testing"

	"github.com/adcondev/poster/pkg/constants"
)

func TestImageBuilder(t *testing.T) {
	doc := NewDocument().
		SetProfile("Test", 80, "WPC1252").
		Image("base64data").
		Width(256).
		Threshold(100).
		Dithering(constants.Atkinson).
		Scaling(constants.Bilinear).
		Center().
		End().
		Build()

	if doc.Commands[0].Type != "image" {
		t.Errorf("Expected type 'image', got '%s'", doc.Commands[0].Type)
	}

	var cmd imageCommand
	_ = json.Unmarshal(doc.Commands[0].Data, &cmd)

	if cmd.Code != "base64data" {
		t.Errorf("Expected code 'base64data', got '%s'", cmd.Code)
	}

	if cmd.PixelWidth != 256 {
		t.Errorf("Expected pixel width 256, got %d", cmd.PixelWidth)
	}

	if cmd.Threshold != 100 {
		t.Errorf("Expected threshold 100, got %d", cmd.Threshold)
	}

	if cmd.Dithering != constants.Atkinson.String() {
		t.Errorf("Expected dithering 'atkinson', got '%s'", cmd.Dithering)
	}

	if cmd.Scaling != constants.Bilinear.String() {
		t.Errorf("Expected scaling 'bilinear', got '%s'", cmd.Scaling)
	}
}

func TestImageBuilderDefaults(t *testing.T) {
	doc := NewDocument().
		SetProfile("Test", 80, "WPC1252").
		Image("base64data").End().
		Build()

	var cmd imageCommand
	_ = json.Unmarshal(doc.Commands[0].Data, &cmd)

	if cmd.PixelWidth != constants.DefaultImagePixelWidth {
		t.Errorf("Expected default pixel width %d, got %d", constants.DefaultImagePixelWidth, cmd.PixelWidth)
	}

	if cmd.Threshold != constants.DefaultImageThreshold {
		t.Errorf("Expected default threshold %d, got %d", constants.DefaultImageThreshold, cmd.Threshold)
	}

	if cmd.Align != constants.DefaultImageAlignment.String() {
		t.Errorf("Expected default align '%s', got '%s'", constants.DefaultImageAlignment.String(), cmd.Align)
	}

	if cmd.Dithering != constants.Atkinson.String() {
		t.Errorf("Expected default dithering '%s', got '%s'", constants.Atkinson.String(), cmd.Dithering)
	}

	if cmd.Scaling != constants.Bilinear.String() {
		t.Errorf("Expected default scaling '%s', got '%s'", constants.Bilinear.String(), cmd.Scaling)
	}
}

func TestImageBuilderAlignment(t *testing.T) {
	tests := []struct {
		name     string
		align    func(*ImageBuilder) *ImageBuilder
		expected string
	}{
		{"Left", func(ib *ImageBuilder) *ImageBuilder { return ib.Left() }, constants.Left.String()},
		{"Center", func(ib *ImageBuilder) *ImageBuilder { return ib.Center() }, constants.Center.String()},
		{"Right", func(ib *ImageBuilder) *ImageBuilder { return ib.Right() }, constants.Right.String()},
		{"Align", func(ib *ImageBuilder) *ImageBuilder { return ib.Align(constants.Left) }, constants.Left.String()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ib := newImageBuilder(NewDocument(), "data")
			tt.align(ib)

			if ib.align != tt.expected {
				t.Errorf("Expected align '%s', got '%s'", tt.expected, ib.align)
			}
		})
	}
}
