package builder

import (
	"encoding/json"
	"testing"

	"github.com/adcondev/poster/pkg/constants"
)

// ============================================================================
// Table Builder Tests
// ============================================================================

func TestTableBuilder(t *testing.T) {
	doc := NewDocument().
		SetProfile("Test", 80, "WPC1252").
		Table().
		Column("Item", 20, constants.Left).
		Column("Price", 10, constants.Right).
		Row("Coffee", "$3.50").
		Row("Muffin", "$4.25").
		PaperWidth(48).
		End().
		Build()

	if doc.Commands[0].Type != "table" {
		t.Errorf("Expected type 'table', got '%s'", doc.Commands[0].Type)
	}

	var cmd tableCommand
	_ = json.Unmarshal(doc.Commands[0].Data, &cmd)

	if len(cmd.Definition.Columns) != 2 {
		t.Errorf("Expected 2 columns, got %d", len(cmd.Definition.Columns))
	}

	if len(cmd.Rows) != 2 {
		t.Errorf("Expected 2 rows, got %d", len(cmd.Rows))
	}
}

func TestTableBuilderOptions(t *testing.T) {
	doc := NewDocument().
		SetProfile("Test", 80, "WPC1252").
		Table().
		Column("Test", 10).
		Row("Data").
		HideHeaders().
		NoHeaderBold().
		NoWordWrap().
		ColumnSpacing(2).
		Align(constants.Center).
		End().
		Build()

	var cmd tableCommand
	_ = json.Unmarshal(doc.Commands[0].Data, &cmd)

	if cmd.ShowHeaders != false {
		t.Error("Expected ShowHeaders to be false")
	}

	if cmd.Options.HeaderBold != false {
		t.Error("Expected HeaderBold to be false")
	}

	if cmd.Options.WordWrap != false {
		t.Error("Expected WordWrap to be false")
	}

	if cmd.Options.ColumnSpacing != 2 {
		t.Errorf("Expected ColumnSpacing 2, got %d", cmd.Options.ColumnSpacing)
	}
}

func TestTableBuilderRows(t *testing.T) {
	rows := [][]string{
		{"A", "1"},
		{"B", "2"},
		{"C", "3"},
	}

	doc := NewDocument().
		SetProfile("Test", 80, "WPC1252").
		Table().
		Column("Letter", 10).
		Column("Number", 10).
		Rows(rows).
		End().
		Build()

	var cmd tableCommand
	_ = json.Unmarshal(doc.Commands[0].Data, &cmd)

	if len(cmd.Rows) != 3 {
		t.Errorf("Expected 3 rows, got %d", len(cmd.Rows))
	}
}
