package builder

import (
	"encoding/json"
	"testing"

	"github.com/adcondev/poster/pkg/constants"
)

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

	if cmd.Definition.Columns[0].Name != "Item" {
		t.Errorf("Expected column name 'Item', got '%s'", cmd.Definition.Columns[0].Name)
	}

	if cmd.Definition.Columns[1].Align != string(constants.Right) {
		t.Errorf("Expected column align 'right', got '%s'", cmd.Definition.Columns[1].Align)
	}

	if len(cmd.Rows) != 2 {
		t.Errorf("Expected 2 rows, got %d", len(cmd.Rows))
	}

	if cmd.Definition.PaperWidth != 48 {
		t.Errorf("Expected paper width 48, got %d", cmd.Definition.PaperWidth)
	}
}

func TestTableBuilder_Defaults(t *testing.T) {
	doc := NewDocument().
		SetProfile("Test", 80, "WPC1252").
		Table().
		Column("Test", 10).
		Row("Data").
		End().
		Build()

	var cmd tableCommand
	_ = json.Unmarshal(doc.Commands[0].Data, &cmd)

	if cmd.ShowHeaders != true {
		t.Error("Expected ShowHeaders to be true by default")
	}

	if cmd.Options == nil {
		t.Fatal("Expected Options to be set")
	}

	if cmd.Options.HeaderBold != true {
		t.Error("Expected HeaderBold to be true by default")
	}

	if cmd.Options.WordWrap != true {
		t.Error("Expected WordWrap to be true by default")
	}

	if cmd.Options.ColumnSpacing != 0 {
		t.Errorf("Expected ColumnSpacing 0, got %d", cmd.Options.ColumnSpacing)
	}
}

func TestTableBuilder_Options(t *testing.T) {
	rows := [][]string{
		{"A", "1"},
		{"B", "2"},
	}

	doc := NewDocument().
		SetProfile("Test", 80, "WPC1252").
		Table().
		Column("Letter", 10).
		Column("Number", 10).
		Rows(rows).
		HideHeaders().
		NoHeaderBold().
		NoWordWrap().
		ColumnSpacing(2).
		Align(constants.Center).
		End().
		Build()

	var cmd tableCommand
	_ = json.Unmarshal(doc.Commands[0].Data, &cmd)

	if len(cmd.Rows) != 2 {
		t.Errorf("Expected 2 rows, got %d", len(cmd.Rows))
	}

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

	if cmd.Options.Align != string(constants.Center) {
		t.Errorf("Expected Align 'center', got '%s'", cmd.Options.Align)
	}
}

func TestTableBuilder_AutoReduce_ExplicitTrue(t *testing.T) {
	doc := NewDocument().
		SetProfile("Test", 58, "PC850").
		Table().
		Column("Item", 20, constants.Left).
		Column("Price", 15, constants.Right).
		Row("Coffee", "$3.50").
		AutoReduce().
		End().
		Build()

	var cmd tableCommand
	_ = json.Unmarshal(doc.Commands[0].Data, &cmd)

	if cmd.Options == nil {
		t.Fatal("expected options to be set")
	}
	if cmd.Options.AutoReduce == nil {
		t.Fatal("expected auto_reduce to be set")
	}
	if *cmd.Options.AutoReduce != true {
		t.Fatalf("expected auto_reduce=true, got %v", *cmd.Options.AutoReduce)
	}
}

func TestTableBuilder_AutoReduce_ExplicitFalse(t *testing.T) {
	doc := NewDocument().
		SetProfile("Test", 58, "PC850").
		Table().
		Column("Item", 20, constants.Left).
		Column("Price", 15, constants.Right).
		Row("Coffee", "$3.50").
		NoAutoReduce().
		End().
		Build()

	var cmd tableCommand
	_ = json.Unmarshal(doc.Commands[0].Data, &cmd)

	if cmd.Options == nil {
		t.Fatal("expected options to be set")
	}
	if cmd.Options.AutoReduce == nil {
		t.Fatal("expected auto_reduce to be set")
	}
	if *cmd.Options.AutoReduce != false {
		t.Fatalf("expected auto_reduce=false, got %v", *cmd.Options.AutoReduce)
	}
}

func TestTableBuilder_AutoReduce_DefaultNil(t *testing.T) {
	// If user doesn't call AutoReduce/NoAutoReduce, the builder should NOT force a value,
	// so the field should be omitted => AutoReduce stays nil.
	doc := NewDocument().
		SetProfile("Test", 58, "PC850").
		Table().
		Column("Item", 10).
		Row("Coffee").
		End().
		Build()

	var cmd tableCommand
	_ = json.Unmarshal(doc.Commands[0].Data, &cmd)

	if cmd.Options == nil {
		t.Fatal("expected options to be set")
	}
	if cmd.Options.AutoReduce != nil {
		t.Fatalf("expected auto_reduce to be nil (omitted), got %v", *cmd.Options.AutoReduce)
	}
}
