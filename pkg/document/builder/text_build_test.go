package builder

import (
	"encoding/json"
	"testing"

	"github.com/adcondev/poster/pkg/constants"
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

	var cmd textCommand
	if err := json.Unmarshal(doc.Commands[0].Data, &cmd); err != nil {
		t.Fatalf("Failed to unmarshal text command: %v", err)
	}

	if cmd.Content.Text != "Hello" {
		t.Errorf("Expected text 'Hello', got '%s'", cmd.Content.Text)
	}

	if cmd.Content.Style == nil {
		t.Fatal("Expected style to be set")
	}

	if cmd.Content.Style.Bold == nil || !*cmd.Content.Style.Bold {
		t.Error("Expected bold to be true")
	}

	if cmd.Content.Style.Size == nil || *cmd.Content.Style.Size != "2x2" {
		t.Errorf("Expected size '2x2', got '%v'", cmd.Content.Style.Size)
	}

	if cmd.Content.Align == nil || *cmd.Content.Align != "center" {
		t.Errorf("Expected align 'center', got '%v'", cmd.Content.Align)
	}
}

func TestTextBuilderDefaults(t *testing.T) {
	doc := NewDocument().
		SetProfile("Test", 80, "WPC1252").
		Text("Plain text").End().
		Build()

	var cmd textCommand
	_ = json.Unmarshal(doc.Commands[0].Data, &cmd)

	if cmd.Content.Text != "Plain text" {
		t.Errorf("Expected text 'Plain text', got '%s'", cmd.Content.Text)
	}

	if cmd.Content.Style != nil {
		t.Error("Expected style to be nil for default text")
	}

	if cmd.Content.Align != nil {
		t.Error("Expected align to be nil for default text")
	}

	if cmd.Label != nil {
		t.Error("Expected label to be nil for default text")
	}

	if cmd.NewLine != nil {
		t.Error("Expected newLine to be nil for default text")
	}
}

func TestTextBuilderWithLabel(t *testing.T) {
	doc := NewDocument().
		SetProfile("Test", 80, "WPC1252").
		Text("$100.00").
		WithLabel("Total").
		LabelSeparator(": ").
		Bold().
		Right().
		NoNewLine().
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

	if cmd.Label.Separator == nil || *cmd.Label.Separator != ": " {
		t.Errorf("Expected separator ': ', got '%v'", cmd.Label.Separator)
	}

	if cmd.Content.Align == nil || *cmd.Content.Align != constants.Right.String() {
		t.Errorf("Expected align 'right', got '%v'", cmd.Content.Align)
	}

	if cmd.NewLine == nil || *cmd.NewLine != false {
		t.Error("Expected NewLine to be false")
	}
}
