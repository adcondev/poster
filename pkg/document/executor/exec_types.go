package executor

// ============================================================================
// Command Data Structures (for handlers)
// ============================================================================

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
