package shared

// ESCPOS and Dots Per Line constants

const (
	// SP represent the byte de "Space" en ESC/POS.
	SP byte = 0x20 // Espacio (carácter de espacio en blanco)
	// FS represent the byte de "File Separator" en ESC/POS.
	FS byte = 0x1C
	// NUL represents the byte de "Null" en ESC/POS.
	NUL byte = 0x00
	// ESC represents the byte de "Escape" en ESC/POS.
	ESC byte = 0x1B
	// GS represents the byte de "Group Separator" en ESC/POS.
	GS byte = 0x1D
	// HT represents the byte de "Horizontal Tab" en ESC/POS.
	HT byte = 0x09

	// Dpl80mm203dpi represents the dots per line for 80mm paper at 203 dpi.
	Dpl80mm203dpi = 576
	// Dpl58mm203dpi represents the dots per line for 58mm paper at 203 dpi.
	Dpl58mm203dpi = 384
)
