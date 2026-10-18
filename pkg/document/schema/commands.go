package schema

import "github.com/adcondev/poster/pkg/tables"

// TextCommand represents a text command payload.
type TextCommand struct {
	Content Content `json:"content"`
	Label   *Label  `json:"label,omitempty"`
	NewLine *bool   `json:"new_line,omitempty"`
}

// Content represents text content and optional style/alignment.
type Content struct {
	Text  string     `json:"text"`
	Style *TextStyle `json:"content_style,omitempty"`
	Align *string    `json:"align,omitempty"`
}

// Label represents an optional prefix label for a text command.
type Label struct {
	Text      string     `json:"text,omitempty"`
	Style     *TextStyle `json:"label_style,omitempty"`
	Separator *string    `json:"separator,omitempty"`
	Align     *string    `json:"align,omitempty"`
}

// TextStyle represents formatting attributes for text or labels.
type TextStyle struct {
	Bold      *bool   `json:"bold,omitempty"`
	Size      *string `json:"size,omitempty"`
	Underline *string `json:"underline,omitempty"`
	Inverse   *bool   `json:"inverse,omitempty"`
	Font      *string `json:"font,omitempty"`
}

// SeparatorCommand represents a horizontal separator line command payload.
type SeparatorCommand struct {
	Char   string `json:"char,omitempty"`
	Length int    `json:"length,omitempty"`
}

// FeedCommand represents a paper feed command payload.
type FeedCommand struct {
	Lines int `json:"lines"`
}

// CutCommand represents a paper cut command payload.
type CutCommand struct {
	Mode string `json:"mode,omitempty"`
	Feed *int   `json:"feed,omitempty"`
}

// PulseCommand represents a cash drawer pulse command payload.
type PulseCommand struct {
	Pin     int `json:"pin,omitempty"`
	OnTime  int `json:"on_time,omitempty"`
	OffTime int `json:"off_time,omitempty"`
}

// BeepCommand represents a printer buzzer beep command payload.
type BeepCommand struct {
	Times int `json:"times,omitempty"`
	Lapse int `json:"lapse,omitempty"`
}

// ImageCommand represents a raster bitmap image command payload.
type ImageCommand struct {
	Code       string `json:"code"`
	PixelWidth int    `json:"pixel_width,omitempty"`
	Align      string `json:"align,omitempty"`
	Threshold  byte   `json:"threshold,omitempty"`
	Dithering  string `json:"dithering,omitempty"`
	Scaling    string `json:"scaling,omitempty"`
}

// QRCommand represents a 2D QR code command payload.
type QRCommand struct {
	Data        string `json:"data"`
	HumanText   string `json:"human_text,omitempty"`
	PixelWidth  int    `json:"pixel_width,omitempty"`
	Correction  string `json:"correction,omitempty"`
	Align       string `json:"align,omitempty"`
	Logo        string `json:"logo,omitempty"`
	CircleShape bool   `json:"circle_shape,omitempty"`
}

// PDF417Command represents a 2D PDF417 barcode command payload.
type PDF417Command struct {
	Data            string `json:"data"`
	Columns         int    `json:"columns,omitempty"`
	Rows            int    `json:"rows,omitempty"`
	ModuleWidth     int    `json:"module_width,omitempty"`
	RowHeight       int    `json:"row_height,omitempty"`
	ErrorCorrection int    `json:"error_correction,omitempty"`
	Align           string `json:"align,omitempty"`
}

// BarcodeCommand represents a 1D barcode command payload.
type BarcodeCommand struct {
	Symbology   string  `json:"symbology"`
	Data        string  `json:"data"`
	Width       *int    `json:"width,omitempty"`
	Height      *int    `json:"height,omitempty"`
	HRIPosition *string `json:"hri_position,omitempty"`
	HRIFont     *string `json:"hri_font,omitempty"`
	Align       *string `json:"align,omitempty"`
}

// TableCommand represents a multi-column table command payload.
type TableCommand struct {
	Definition  tables.Definition `json:"definition"`
	ShowHeaders *bool             `json:"show_headers,omitempty"`
	Rows        [][]string        `json:"rows"`
	Options     *TableOptions     `json:"options,omitempty"`
}

// TableOptions represents layout options for a table command.
type TableOptions struct {
	HeaderBold    *bool  `json:"header_bold,omitempty"`
	WordWrap      *bool  `json:"word_wrap,omitempty"`
	ColumnSpacing *int   `json:"column_spacing,omitempty"`
	Align         string `json:"align,omitempty"`
	Font          string `json:"font,omitempty"`
	AutoReduce    *bool  `json:"auto_reduce,omitempty"`
}

// RawCommand represents a raw ESC/POS byte sequence command payload.
type RawCommand struct {
	Hex      string `json:"hex"`
	Format   string `json:"format,omitempty"`
	Comment  string `json:"comment,omitempty"`
	SafeMode bool   `json:"safe_mode,omitempty"`
}
