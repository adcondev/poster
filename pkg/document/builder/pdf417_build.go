package builder

import (
	"github.com/adcondev/poster/pkg/constants"
	"github.com/adcondev/poster/pkg/document/schema"
)

// PDF417Builder constructs PDF417 2D barcode commands.
type PDF417Builder struct {
	parent          *DocumentBuilder
	data            string
	columns         int
	rows            int
	moduleWidth     int
	rowHeight       int
	errorCorrection int
	align           string
}

func newPDF417Builder(parent *DocumentBuilder, data string) *PDF417Builder {
	return &PDF417Builder{
		parent:          parent,
		data:            data,
		moduleWidth:     3,
		rowHeight:       3,
		errorCorrection: 2,
		align:           constants.Center.String(),
	}
}

// Columns sets the number of data columns (0-30, 0 = automatic).
func (pb *PDF417Builder) Columns(columns int) *PDF417Builder {
	pb.columns = columns
	return pb
}

// Rows sets the number of rows (0 or 3-90, 0 = automatic).
func (pb *PDF417Builder) Rows(rows int) *PDF417Builder {
	pb.rows = rows
	return pb
}

// ModuleWidth sets module width in dots (2-8).
func (pb *PDF417Builder) ModuleWidth(width int) *PDF417Builder {
	pb.moduleWidth = width
	return pb
}

// RowHeight sets row height multiplier (2-8).
func (pb *PDF417Builder) RowHeight(height int) *PDF417Builder {
	pb.rowHeight = height
	return pb
}

// ErrorCorrection sets error correction level (0-8).
func (pb *PDF417Builder) ErrorCorrection(level int) *PDF417Builder {
	pb.errorCorrection = level
	return pb
}

// Align sets the PDF417 alignment (left, center, right).
func (pb *PDF417Builder) Align(align string) *PDF417Builder {
	pb.align = align
	return pb
}

// Left aligns PDF417 barcode to the left.
func (pb *PDF417Builder) Left() *PDF417Builder {
	pb.align = constants.Left.String()
	return pb
}

// Center centers the PDF417 barcode (default).
func (pb *PDF417Builder) Center() *PDF417Builder {
	pb.align = constants.Center.String()
	return pb
}

// Right aligns PDF417 barcode to the right.
func (pb *PDF417Builder) Right() *PDF417Builder {
	pb.align = constants.Right.String()
	return pb
}

// End finishes the PDF417 command and appends it to the document.
func (pb *PDF417Builder) End() *DocumentBuilder {
	cmd := schema.PDF417Command{
		Data:            pb.data,
		Columns:         pb.columns,
		Rows:            pb.rows,
		ModuleWidth:     pb.moduleWidth,
		RowHeight:       pb.rowHeight,
		ErrorCorrection: pb.errorCorrection,
		Align:           pb.align,
	}
	return pb.parent.addCommand("pdf417", cmd)
}
