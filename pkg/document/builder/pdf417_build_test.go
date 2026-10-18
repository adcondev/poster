package builder

import (
	"encoding/json"
	"testing"

	"github.com/adcondev/poster/pkg/constants"
	"github.com/adcondev/poster/pkg/document/schema"
)

func TestPDF417Builder(t *testing.T) {
	doc := NewDocument().
		SetProfile("Test", 80, "WPC1252").
		PDF417("https://example.com/ticket").
		Columns(5).
		Rows(10).
		ModuleWidth(4).
		RowHeight(5).
		ErrorCorrection(3).
		Center().
		End().
		Build()

	if len(doc.Commands) != 1 {
		t.Fatalf("expected 1 command, got %d", len(doc.Commands))
	}

	if doc.Commands[0].Type != "pdf417" {
		t.Errorf("expected type 'pdf417', got '%s'", doc.Commands[0].Type)
	}

	var cmd schema.PDF417Command
	if err := json.Unmarshal(doc.Commands[0].Data, &cmd); err != nil {
		t.Fatalf("failed to unmarshal command data: %v", err)
	}

	if cmd.Data != "https://example.com/ticket" {
		t.Errorf("expected data 'https://example.com/ticket', got '%s'", cmd.Data)
	}
	if cmd.Columns != 5 {
		t.Errorf("expected columns 5, got %d", cmd.Columns)
	}
	if cmd.Rows != 10 {
		t.Errorf("expected rows 10, got %d", cmd.Rows)
	}
	if cmd.ModuleWidth != 4 {
		t.Errorf("expected module width 4, got %d", cmd.ModuleWidth)
	}
	if cmd.RowHeight != 5 {
		t.Errorf("expected row height 5, got %d", cmd.RowHeight)
	}
	if cmd.ErrorCorrection != 3 {
		t.Errorf("expected error correction 3, got %d", cmd.ErrorCorrection)
	}
	if cmd.Align != constants.Center.String() {
		t.Errorf("expected align 'center', got '%s'", cmd.Align)
	}
}

func TestPDF417BuilderDefaults(t *testing.T) {
	doc := NewDocument().
		SetProfile("Test", 80, "WPC1252").
		PDF417("default-data").
		End().
		Build()

	if len(doc.Commands) != 1 {
		t.Fatalf("expected 1 command, got %d", len(doc.Commands))
	}

	var cmd schema.PDF417Command
	if err := json.Unmarshal(doc.Commands[0].Data, &cmd); err != nil {
		t.Fatalf("failed to unmarshal command data: %v", err)
	}

	if cmd.Data != "default-data" {
		t.Errorf("expected data 'default-data', got '%s'", cmd.Data)
	}
	if cmd.Columns != 0 {
		t.Errorf("expected default columns 0, got %d", cmd.Columns)
	}
	if cmd.Rows != 0 {
		t.Errorf("expected default rows 0, got %d", cmd.Rows)
	}
	if cmd.ModuleWidth != 3 {
		t.Errorf("expected default module width 3, got %d", cmd.ModuleWidth)
	}
	if cmd.RowHeight != 3 {
		t.Errorf("expected default row height 3, got %d", cmd.RowHeight)
	}
	if cmd.ErrorCorrection != 2 {
		t.Errorf("expected default error correction 2, got %d", cmd.ErrorCorrection)
	}
	if cmd.Align != constants.Center.String() {
		t.Errorf("expected default align 'center', got '%s'", cmd.Align)
	}
}

func TestPDF417BuilderAlignment(t *testing.T) {
	tests := []struct {
		name     string
		align    func(*PDF417Builder) *PDF417Builder
		expected string
	}{
		{"Left", func(pb *PDF417Builder) *PDF417Builder { return pb.Left() }, constants.Left.String()},
		{"Center", func(pb *PDF417Builder) *PDF417Builder { return pb.Center() }, constants.Center.String()},
		{"Right", func(pb *PDF417Builder) *PDF417Builder { return pb.Right() }, constants.Right.String()},
		{"Custom Align", func(pb *PDF417Builder) *PDF417Builder { return pb.Align("right") }, "right"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := NewDocument()
			tt.align(doc.PDF417("data")).End()
			bDoc := doc.Build()

			var cmd schema.PDF417Command
			_ = json.Unmarshal(bDoc.Commands[0].Data, &cmd)
			if cmd.Align != tt.expected {
				t.Errorf("expected align '%s', got '%s'", tt.expected, cmd.Align)
			}
		})
	}
}
