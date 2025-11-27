package document_test

import (
	"encoding/json"
	"testing"

	"github.com/adcondev/pos-printer/pkg/document"
	"github.com/adcondev/pos-printer/pkg/tables"
)

func TestTableBuilder_FluentAPI(t *testing.T) {
	doc := document.NewBuilder().
		SetProfile("test", 80, "WPC1252").
		NewTableBuilder().
		AddColumn("Product", 20, "").
		AddColumn("Price", 10, tables.Right).
		SetPaperWidth(48).
		AddRow("Coffee", "$4.50").
		AddRow("Tea", "$3.00").
		ShowHeaders(true).
		Done().
		Build()

	if len(doc.Commands) != 1 {
		t.Errorf("expected 1 command, got %d", len(doc.Commands))
	}

	if doc.Commands[0].Type != "table" {
		t.Errorf("expected 'table' command, got %q", doc.Commands[0].Type)
	}
}

func TestTableBuilder_AddColumn(t *testing.T) {
	builder := document.NewBuilder()
	tableBuilder := builder.NewTableBuilder().
		AddColumn("A", 10, "").
		AddColumn("B", 15, tables.Center).
		AddColumn("C", 20, tables.Right)

	cmd := tableBuilder.Build()

	if len(cmd.Definition.Columns) != 3 {
		t.Errorf("expected 3 columns, got %d", len(cmd.Definition.Columns))
	}

	// Verify column details
	if cmd.Definition.Columns[0].Header != "A" || cmd.Definition.Columns[0].Width != 10 {
		t.Error("Column A not configured correctly")
	}
	if cmd.Definition.Columns[1].Align != tables.Center {
		t.Error("Column B alignment not set correctly")
	}
	if cmd.Definition.Columns[2].Align != tables.Right {
		t.Error("Column C alignment not set correctly")
	}
}

func TestTableBuilder_AddRow(t *testing.T) {
	builder := document.NewBuilder()
	tableBuilder := builder.NewTableBuilder().
		AddColumn("A", 10, "").
		AddRow("row1-a").
		AddRow("row2-a").
		AddRow("row3-a")

	cmd := tableBuilder.Build()

	if len(cmd.Rows) != 3 {
		t.Errorf("expected 3 rows, got %d", len(cmd.Rows))
	}
}

func TestTableBuilder_AddRows(t *testing.T) {
	builder := document.NewBuilder()
	rows := [][]string{
		{"a", "b"},
		{"c", "d"},
		{"e", "f"},
	}

	tableBuilder := builder.NewTableBuilder().
		AddColumn("Col1", 10, "").
		AddColumn("Col2", 10, "").
		AddRows(rows)

	cmd := tableBuilder.Build()

	if len(cmd.Rows) != 3 {
		t.Errorf("expected 3 rows, got %d", len(cmd.Rows))
	}
}

func TestTableBuilder_Options(t *testing.T) {
	builder := document.NewBuilder()
	tableBuilder := builder.NewTableBuilder().
		AddColumn("Test", 10, "").
		ShowHeaders(false).
		HeaderBold(false).
		WordWrap(false).
		ColumnSpacing(3).
		Align("center")

	cmd := tableBuilder.Build()

	if cmd.ShowHeaders != false {
		t.Error("ShowHeaders should be false")
	}
	if cmd.Options.HeaderBold != false {
		t.Error("HeaderBold should be false")
	}
	if cmd.Options.WordWrap != false {
		t.Error("WordWrap should be false")
	}
	if cmd.Options.ColumnSpacing != 3 {
		t.Errorf("ColumnSpacing should be 3, got %d", cmd.Options.ColumnSpacing)
	}
	if cmd.Options.Align != "center" {
		t.Errorf("Align should be 'center', got %q", cmd.Options.Align)
	}
}

func TestTableBuilder_NegativeColumnSpacing(t *testing.T) {
	builder := document.NewBuilder()
	tableBuilder := builder.NewTableBuilder().
		AddColumn("Test", 10, "").
		ColumnSpacing(-5)

	cmd := tableBuilder.Build()

	if cmd.Options.ColumnSpacing != 0 {
		t.Errorf("Negative spacing should be corrected to 0, got %d", cmd.Options.ColumnSpacing)
	}
}

func TestTableBuilder_SetPaperWidth(t *testing.T) {
	builder := document.NewBuilder()
	tableBuilder := builder.NewTableBuilder().
		AddColumn("Test", 10, "").
		SetPaperWidth(32)

	cmd := tableBuilder.Build()

	if cmd.Definition.PaperWidth != 32 {
		t.Errorf("PaperWidth should be 32, got %d", cmd.Definition.PaperWidth)
	}
}

func TestTableBuilder_Done_AddsToDocument(t *testing.T) {
	doc := document.NewBuilder().
		NewTableBuilder().
		AddColumn("Test", 10, "").
		AddRow("Value").
		Done().
		AddFeed(2). // Continue building after table
		Build()

	if len(doc.Commands) != 2 {
		t.Errorf("expected 2 commands, got %d", len(doc.Commands))
	}
	if doc.Commands[0].Type != "table" {
		t.Errorf("first command should be 'table', got %q", doc.Commands[0].Type)
	}
	if doc.Commands[1].Type != "feed" {
		t.Errorf("second command should be 'feed', got %q", doc.Commands[1].Type)
	}
}

func TestTableBuilder_Done_EmptyTable_Skipped(t *testing.T) {
	doc := document.NewBuilder().
		NewTableBuilder().
		// No columns added
		Done().
		Build()

	if len(doc.Commands) != 0 {
		t.Errorf("expected 0 commands (empty table skipped), got %d", len(doc.Commands))
	}
}

func TestTableBuilder_Build_Standalone(t *testing.T) {
	builder := document.NewBuilder()
	tableBuilder := builder.NewTableBuilder().
		AddColumn("A", 10, "").
		AddRow("1")

	// Use Build() to get the command without adding to parent
	cmd := tableBuilder.Build()

	if cmd == nil {
		t.Fatal("Build() returned nil")
	}
	if len(cmd.Definition.Columns) != 1 {
		t.Error("Command should have 1 column")
	}

	// Parent document should still be empty
	doc := builder.Build()
	if len(doc.Commands) != 0 {
		t.Error("Parent document should be empty when using Build() instead of Done()")
	}
}

func TestTableBuilder_MultipleTables(t *testing.T) {
	doc := document.NewBuilder().
		NewTableBuilder().
		AddColumn("Table1", 10, "").
		AddRow("Data1").
		Done().
		AddSeparator("-", 48).
		NewTableBuilder().
		AddColumn("Table2", 15, "").
		AddRow("Data2").
		Done().
		Build()

	if len(doc.Commands) != 3 {
		t.Errorf("expected 3 commands, got %d", len(doc.Commands))
	}

	// Verify first table
	var cmd1 document.TableCommand
	if err := json.Unmarshal(doc.Commands[0].Data, &cmd1); err != nil {
		t.Fatalf("failed to unmarshal first table: %v", err)
	}
	if cmd1.Definition.Columns[0].Header != "Table1" {
		t.Error("First table header incorrect")
	}

	// Verify separator
	if doc.Commands[1].Type != "separator" {
		t.Error("Second command should be separator")
	}

	// Verify second table
	var cmd2 document.TableCommand
	if err := json.Unmarshal(doc.Commands[2].Data, &cmd2); err != nil {
		t.Fatalf("failed to unmarshal second table: %v", err)
	}
	if cmd2.Definition.Columns[0].Header != "Table2" {
		t.Error("Second table header incorrect")
	}
}

func TestAddTableCommand(t *testing.T) {
	// Test the AddTableCommand method
	tableCmd := &document.TableCommand{
		Definition: tables.Definition{
			Columns: []tables.Column{
				{Header: "Pre-built", Width: 15},
			},
			PaperWidth: 48,
		},
		ShowHeaders: true,
		Rows:        [][]string{{"Value"}},
	}

	doc := document.NewBuilder().
		AddTableCommand(tableCmd).
		Build()

	if len(doc.Commands) != 1 {
		t.Errorf("expected 1 command, got %d", len(doc.Commands))
	}
}

func TestAddTableCommand_Nil(t *testing.T) {
	doc := document.NewBuilder().
		AddTableCommand(nil).
		Build()

	if len(doc.Commands) != 0 {
		t.Errorf("expected 0 commands for nil table, got %d", len(doc.Commands))
	}
}

// Benchmark comparing old vs new API
func BenchmarkTableBuilder_NewAPI(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = document.NewBuilder().
			NewTableBuilder().
			AddColumn("Product", 20, "").
			AddColumn("Qty", 5, tables.Center).
			AddColumn("Price", 10, tables.Right).
			SetPaperWidth(48).
			AddRow("Coffee", "2", "$4.50").
			AddRow("Tea", "1", "$3.00").
			ShowHeaders(true).
			Done().
			Build()
	}
}
