package builder

import (
	"encoding/json"
	"testing"

	"github.com/adcondev/poster/pkg/constants"
)

// ============================================================================
// Image Builder Tests
// ============================================================================

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

	if cmd.PixelWidth != 256 {
		t.Errorf("Expected pixel width 256, got %d", cmd.PixelWidth)
	}

	if cmd.Threshold != 100 {
		t.Errorf("Expected threshold 100, got %d", cmd.Threshold)
	}

	if cmd.Dithering != constants.Atkinson.String() {
		t.Errorf("Expected dithering 'atkinson', got '%s'", cmd.Dithering)
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
}
