package executor

import (
	"github.com/adcondev/pos-printer/pkg/tables"
)

// ============================================================================
// Command Data Structures (for handlers)
// ============================================================================

// TextCommand for text handler
type TextCommand struct {
	Content Content `json:"content"`
	Label   *Label  `json:"label,omitempty"`
	NewLine *bool   `json:"new_line,omitempty"`
}

// Content for text
type Content struct {
	Text  string     `json:"text"`
	Style *TextStyle `json:"content_style,omitempty"`
	Align *string    `json:"align,omitempty"`
}

// Label for text
type Label struct {
	Text      string     `json:"text,omitempty"`
	Style     *TextStyle `json:"label_style,omitempty"`
	Separator *string    `json:"separator,omitempty"`
	Align     *string    `json:"align,omitempty"`
}

// TextStyle for text formatting
type TextStyle struct {
	Bold      *bool   `json:"bold,omitempty"`
	Size      *string `json:"size,omitempty"`
	Underline *string `json:"underline,omitempty"`
	Inverse   *bool   `json:"inverse,omitempty"`
	Font      *string `json:"font,omitempty"`
}

// ImageCommand for image handler
type ImageCommand struct {
	Code string `json:"code"`
	// Format     string `json:"format,omitempty"`
	PixelWidth int    `json:"pixel_width,omitempty"`
	Align      string `json:"align,omitempty"`
	Threshold  byte   `json:"threshold,omitempty"`
	Dithering  string `json:"dithering,omitempty"`
	Scaling    string `json:"scaling,omitempty"`
}

// SeparatorCommand for separator handler
type SeparatorCommand struct {
	Char   string `json:"char,omitempty"`
	Length int    `json:"length,omitempty"`
	// TODO: Add TextStyle if needed
	// TODO: Add Align if needed
	// TODO: Add Pre and Post feed line if needed
}

// FeedCommand for feed handler
type FeedCommand struct {
	Lines int `json:"lines"`
}

// CutCommand for cut handler
type CutCommand struct {
	Mode string `json:"mode,omitempty"`
	Feed int    `json:"feed,omitempty"`
}

// QRCommand for QR handler
type QRCommand struct {
	Data        string `json:"data"`
	HumanText   string `json:"human_text,omitempty"`
	PixelWidth  int    `json:"pixel_width,omitempty"`
	Correction  string `json:"correction,omitempty"`
	Align       string `json:"align,omitempty"`
	Logo        string `json:"logo,omitempty"`
	CircleShape bool   `json:"circle_shape,omitempty"`
}

// BarcodeCommand for barcode handler
type BarcodeCommand struct {
	Symbology   string  `json:"symbology"`
	Data        string  `json:"data"`
	Width       *int    `json:"width,omitempty"`
	Height      *int    `json:"height,omitempty"`
	HRIPosition *string `json:"hri_position,omitempty"`
	HRIFont     *string `json:"hri_font,omitempty"`
	Align       *string `json:"align,omitempty"`
	// TODO: Check if CodeSet is needed
}

// TableCommand for table handler
type TableCommand struct {
	Definition  tables.Definition `json:"definition"`
	ShowHeaders bool              `json:"show_headers,omitempty"`
	Rows        [][]string        `json:"rows"`
	Options     *TableOptions     `json:"options,omitempty"`
}

// TableOptions for table configuration
type TableOptions struct {
	HeaderBold    bool   `json:"header_bold,omitempty"`
	WordWrap      bool   `json:"word_wrap,omitempty"`
	ColumnSpacing int    `json:"column_spacing,omitempty"`
	Align         string `json:"align,omitempty"`
}

// RawCommand for raw handler
type RawCommand struct {
	Hex      string `json:"hex"`
	Format   string `json:"format,omitempty"`
	Comment  string `json:"comment,omitempty"`
	SafeMode bool   `json:"safe_mode,omitempty"`
}
