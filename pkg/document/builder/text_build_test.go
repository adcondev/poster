package builder

import (
	"encoding/json"
	"testing"
)

// ============================================================================
// Text Builder Tests
// ============================================================================

func TestTextBuilder(t *testing.T) {
	doc := NewDocument().
		SetProfile("Test", 80, "WPC1252").
		Text("Hello").
		Bold().
		Size("2x2").
		Underline("1pt").
		Inverse().
		Font("B").
		Center().
		End().
		Build()

	if len(doc.Commands) != 1 {
		t.Fatalf("Expected 1 command, got %d", len(doc.Commands))
	}

	if doc.Commands[0].Type != "text" {
		t.Errorf("Expected type 'text', got '%s'", doc.Commands[0].Type)
	}

	// Parse the command data
	var cmd textCommand
	if err := json.Unmarshal(doc.Commands[0].Data, &cmd); err != nil {
		t.Fatalf("Failed to unmarshal text command: %v", err)
	}

	if cmd.Content.Text != "Hello" {
		t.Errorf("Expected text 'Hello', got '%s'", cmd.Content.Text)
	}
}

func TestTextWithLabel(t *testing.T) {
	doc := NewDocument().
		SetProfile("Test", 80, "WPC1252").
		Text("$100. 00").
		WithLabel("Total").
		LabelSeparator(": ").
		Bold().
		Right().
		End().
		Build()

	var cmd textCommand
	_ = json.Unmarshal(doc.Commands[0].Data, &cmd)

	if cmd.Label == nil {
		t.Fatal("Expected label to be set")
	}

	if cmd.Label.Text != "Total" {
		t.Errorf("Expected label 'Total', got '%s'", cmd.Label.Text)
	}
}

func TestTextNoNewLine(t *testing.T) {
	doc := NewDocument().
		SetProfile("Test", 80, "WPC1252").
		Text("No newline").NoNewLine().End().
		Build()

	var cmd textCommand
	_ = json.Unmarshal(doc.Commands[0].Data, &cmd)

	if cmd.NewLine == nil || *cmd.NewLine != false {
		t.Error("Expected NewLine to be false")
	}
}
