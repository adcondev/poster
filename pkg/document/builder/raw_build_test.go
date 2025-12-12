package builder

import (
	"encoding/json"
	"testing"

	"github.com/adcondev/poster/pkg/constants"
)

func TestRawBuilder(t *testing.T) {
	doc := NewDocument().
		SetProfile("Test", 80, "WPC1252").
		Raw("1B 40").
		Comment("Initialize printer").
		Safe().
		End().
		Build()

	if doc.Commands[0].Type != "raw" {
		t.Errorf("Expected type 'raw', got '%s'", doc.Commands[0].Type)
	}

	var cmd rawCommand
	_ = json.Unmarshal(doc.Commands[0].Data, &cmd)

	if cmd.Hex != "1B 40" {
		t.Errorf("Expected hex '1B 40', got '%s'", cmd.Hex)
	}

	if cmd.Comment != "Initialize printer" {
		t.Errorf("Expected comment 'Initialize printer', got '%s'", cmd.Comment)
	}

	if !cmd.SafeMode {
		t.Error("Expected SafeMode to be true")
	}

	if cmd.Format != constants.Hex.String() {
		t.Errorf("Expected format 'hex', got '%s'", cmd.Format)
	}
}

func TestRawBuilderDefaults(t *testing.T) {
	doc := NewDocument().
		SetProfile("Test", 80, "WPC1252").
		Raw("1B40").End().
		Build()

	var cmd rawCommand
	_ = json.Unmarshal(doc.Commands[0].Data, &cmd)

	if cmd.Format != "hex" {
		t.Errorf("Expected default format 'hex', got '%s'", cmd.Format)
	}

	if cmd.Comment != "" {
		t.Errorf("Expected empty comment, got '%s'", cmd.Comment)
	}

	if cmd.SafeMode != false {
		t.Error("Expected SafeMode to be false by default")
	}
}

func TestRawBuilderFormats(t *testing.T) {
	t.Run("Base64Format", func(t *testing.T) {
		doc := NewDocument().
			SetProfile("Test", 80, "WPC1252").
			Raw("G0A=").Base64().End().
			Build()

		var cmd rawCommand
		_ = json.Unmarshal(doc.Commands[0].Data, &cmd)

		if cmd.Format != "base64" {
			t.Errorf("Expected format 'base64', got '%s'", cmd.Format)
		}
	})

	t.Run("FromBytes", func(t *testing.T) {
		data := []byte{0x1B, 0x40}
		doc := NewDocument().
			SetProfile("Test", 80, "WPC1252").
			FromBytes(data).End().
			Build()

		var cmd rawCommand
		_ = json.Unmarshal(doc.Commands[0].Data, &cmd)

		if cmd.Hex != "1b40" {
			t.Errorf("Expected hex '1b40', got '%s'", cmd.Hex)
		}

		if cmd.Format != "hex" {
			t.Errorf("Expected format 'hex', got '%s'", cmd.Format)
		}
	})
}
