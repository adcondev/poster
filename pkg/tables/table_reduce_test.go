package tables_test

import (
	"testing"

	"github.com/adcondev/poster/pkg/tables"
)

func TestReduceToFit(t *testing.T) {
	tests := []struct {
		name           string
		columns        []tables.Column
		maxChars       int
		spacing        int
		minWidth       int
		wantSuccess    bool
		wantReductions int
	}{
		{
			name: "already fits - no reduction needed",
			columns: []tables.Column{
				{Name: "A", Width: 15},
				{Name: "B", Width: 15},
			},
			maxChars:       32,
			spacing:        1,
			minWidth:       3,
			wantSuccess:    true,
			wantReductions: 0,
		},
		{
			name: "simple reduction - reduces longest",
			columns: []tables.Column{
				{Name: "Long", Width: 20},
				{Name: "Short", Width: 5},
			},
			maxChars:       20,
			spacing:        1,
			minWidth:       3,
			wantSuccess:    true,
			wantReductions: 6,
		},
		{
			name: "alternating reduction - two similar columns",
			columns: []tables.Column{
				{Name: "A", Width: 15},
				{Name: "B", Width: 16},
			},
			maxChars:       25,
			spacing:        1,
			minWidth:       3,
			wantSuccess:    true,
			wantReductions: 7,
		},
		{
			name: "fails when hitting min width",
			columns: []tables.Column{
				{Name: "A", Width: 5},
				{Name: "B", Width: 5},
			},
			maxChars:    5,
			spacing:     1,
			minWidth:    4,
			wantSuccess: false,
		},
		{
			name: "preserves tiny columns",
			columns: []tables.Column{
				{Name: "ID", Width: 3},
				{Name: "Desc", Width: 30},
			},
			maxChars:       20,
			spacing:        1,
			minWidth:       3,
			wantSuccess:    true,
			wantReductions: 14,
		},
		{
			name: "four columns with mixed widths",
			columns: []tables.Column{
				{Name: "#", Width: 2},
				{Name: "Product", Width: 20},
				{Name: "Qty", Width: 4},
				{Name: "Price", Width: 10},
			},
			maxChars:       32,
			spacing:        1,
			minWidth:       2,
			wantSuccess:    true,
			wantReductions: 7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tables.ReduceToFit(tt.columns, tt.maxChars, tt.spacing, tt.minWidth)

			if tt.wantSuccess {
				if err != nil {
					t.Errorf("ReduceToFit() unexpected error: %v", err)
					return
				}
				if !result.Success {
					t.Errorf("ReduceToFit() Success = false, want true")
				}
				if tt.wantReductions > 0 && result.Reductions != tt.wantReductions {
					t.Errorf("ReduceToFit() Reductions = %d, want %d",
						result.Reductions, tt.wantReductions)
				}
				if result.ReducedWidth > tt.maxChars {
					t.Errorf("ReduceToFit() ReducedWidth %d > maxChars %d",
						result.ReducedWidth, tt.maxChars)
				}
			} else if err == nil {
				t.Errorf("ReduceToFit() expected error, got nil")
			}
		})
	}
}

func TestReduceToFit_PreservesSmallColumns(t *testing.T) {
	columns := []tables.Column{
		{Name: "ID", Width: 3},
		{Name: "Description", Width: 25},
		{Name: "Qty", Width: 3},
	}

	result, err := tables.ReduceToFit(columns, 20, 1, 3)

	if err != nil {
		t.Fatalf("ReduceToFit() error: %v", err)
	}

	// ID and Qty should remain at 3 (minimum)
	if result.Columns[0].Width != 3 {
		t.Errorf("ID column width = %d, want 3 (preserved)", result.Columns[0].Width)
	}
	if result.Columns[2].Width != 3 {
		t.Errorf("Qty column width = %d, want 3 (preserved)", result.Columns[2].Width)
	}

	// Description should have absorbed all reductions
	// Original: 3 + 25 + 3 + 2 gaps = 33
	// Target: 20
	// Reductions needed: 13
	// Description:  25 - 13 = 12
	if result.Columns[1].Width != 12 {
		t.Errorf("Description column width = %d, want 12 (reduced)", result.Columns[1].Width)
	}
}
