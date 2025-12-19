// Package tables provides table generation and rendering for ESC/POS thermal printers.
//
// This package implements a complete table layout engine designed for thermal receipt
// printers with the following key features:
//
// # Core Features
//
//   - Dynamic column width calculation with configurable spacing
//   - Automatic word wrapping for long text content
//   - Multiple column alignment options (left, center, right)
//   - Bold header styling with ESC/POS command injection
//   - Hardware-aware validation against paper width limits
//
// # Overflow Protection
//
// The package includes automatic overflow protection that prevents tables from
// exceeding the physical paper width:
//
//   - Validates total column widths against calculated character limits
//   - Auto-reduction algorithm shrinks widest columns first to preserve small columns
//   - Configurable minimum column width (default: 3 characters)
//   - Clear error messages when reduction is impossible
//
// # Paper Width Calculation
//
// Character limits are calculated from printer profile:
//
//	MaxChars = DotsPerLine / FontWidth
//
// Standard values at 203 DPI with Font A (12 dots):
//
//	58mm paper: 384 dots / 12 = 32 characters
//	80mm paper: 576 dots / 12 = 48 characters
//
// # Usage Example
//
//	definition := &tables.Definition{
//	    Columns: []tables.Column{
//	        {Name: "Item", Width: 20, Align: constants.Left},
//	        {Name: "Price", Width: 10, Align: constants.Right},
//	    },
//	}
//
//	opts := tables.DefaultOptions()
//	engine := tables.NewEngine(definition, opts)
//
//	data := &tables.Data{
//	    Definition:   *definition,
//	    ShowHeaders: true,
//	    Rows: []tables. Row{
//	        {"Coffee", "$3.50"},
//	        {"Sandwich", "$8.00"},
//	    },
//	}
//
//	var buf strings.Builder
//	engine.Render(&buf, data)
//
// # Auto-Reduction Example
//
// When columns exceed paper width, the ReduceToFit function automatically
// shrinks the widest columns:
//
//	columns := []tables.Column{
//	    {Name: "#", Width: 3},           // Preserved (small)
//	    {Name: "Description", Width: 30}, // Reduced first (largest)
//	    {Name: "Qty", Width: 3},          // Preserved (small)
//	}
//
//	result, err := tables. ReduceToFit(columns, 32, 1, 3)
//	// Result: [3, 24, 3] - Description absorbed all reductions
//
// # File Organization
//
//   - table_types.go: Core types (Column, Definition, Row, Data) and validation
//   - table_engine.go: TabEngine for rendering tables to io.Writer
//   - table_text.go: Text utilities (WrapText, PadString)
//   - table_reduce.go: Auto-reduction algorithm for overflow protection
package tables
