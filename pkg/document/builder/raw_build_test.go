package builder

import (
	"encoding/json"
	"testing"
)

// ============================================================================
// Raw Builder Tests
// ============================================================================

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
}

func TestRawBuilderBase64(t *testing.T) {
	doc := NewDocument().
		SetProfile("Test", 80, "WPC1252").
		Raw("G0A=").Base64().End().
		Build()

	var cmd rawCommand
	_ = json.Unmarshal(doc.Commands[0].Data, &cmd)

	if cmd.Format != "base64" {
		t.Errorf("Expected format 'base64', got '%s'", cmd.Format)
	}
}

func TestFromBytes(t *testing.T) {
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
}
