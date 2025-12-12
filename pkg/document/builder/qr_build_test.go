package builder

import (
	"encoding/json"
	"testing"

	"github.com/adcondev/poster/pkg/constants"
)

func TestQRBuilder(t *testing.T) {
	doc := NewDocument().
		SetProfile("Test", 80, "WPC1252").
		QR("https://example.com").
		Size(200).
		Correction("H").
		WithText("Scan me").
		WithLogo("logobase64").
		Center().
		End().
		Build()

	if doc.Commands[0].Type != "qr" {
		t.Errorf("Expected type 'qr', got '%s'", doc.Commands[0].Type)
	}

	var cmd qrCommand
	_ = json.Unmarshal(doc.Commands[0].Data, &cmd)

	if cmd.Data != "https://example.com" {
		t.Errorf("Expected data 'https://example.com', got '%s'", cmd.Data)
	}

	if cmd.PixelWidth != 200 {
		t.Errorf("Expected pixel width 200, got %d", cmd.PixelWidth)
	}

	if cmd.Correction != "H" {
		t.Errorf("Expected correction 'H', got '%s'", cmd.Correction)
	}

	if cmd.HumanText != "Scan me" {
		t.Errorf("Expected human text 'Scan me', got '%s'", cmd.HumanText)
	}

	if cmd.Logo != "logobase64" {
		t.Errorf("Expected logo 'logobase64', got '%s'", cmd.Logo)
	}

	if cmd.Align != constants.Center.String() {
		t.Errorf("Expected align 'center', got '%s'", cmd.Align)
	}
}

func TestQRBuilderDefaults(t *testing.T) {
	doc := NewDocument().
		SetProfile("Test", 80, "WPC1252").
		QR("data").End().
		Build()

	var cmd qrCommand
	_ = json.Unmarshal(doc.Commands[0].Data, &cmd)

	if cmd.PixelWidth != 128 {
		t.Errorf("Expected default pixel width 128, got %d", cmd.PixelWidth)
	}

	if cmd.Correction != "Q" {
		t.Errorf("Expected default correction 'Q', got '%s'", cmd.Correction)
	}

	if cmd.Align != constants.Center.String() {
		t.Errorf("Expected default align 'center', got '%s'", cmd.Align)
	}

	if cmd.HumanText != "" {
		t.Errorf("Expected empty human text, got '%s'", cmd.HumanText)
	}

	if cmd.Logo != "" {
		t.Errorf("Expected empty logo, got '%s'", cmd.Logo)
	}

	if cmd.CircleShape != false {
		t.Error("Expected CircleShape to be false by default")
	}
}

func TestQRBuilderAlignment(t *testing.T) {
	tests := []struct {
		name     string
		align    func(*QRBuilder) *QRBuilder
		expected string
	}{
		{"Left", func(qb *QRBuilder) *QRBuilder { return qb.Left() }, constants.Left.String()},
		{"Center", func(qb *QRBuilder) *QRBuilder { return qb.Center() }, constants.Center.String()},
		{"Right", func(qb *QRBuilder) *QRBuilder { return qb.Right() }, constants.Right.String()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qb := newQRBuilder(NewDocument(), "data")
			tt.align(qb)

			if qb.align != tt.expected {
				t.Errorf("Expected align '%s', got '%s'", tt.expected, qb.align)
			}
		})
	}
}
