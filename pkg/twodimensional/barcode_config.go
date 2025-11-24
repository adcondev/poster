// Package twodimensional contains configurations and utilities for 2D printing commands.
package twodimensional

import (
	"github.com/adcondev/pos-printer/pkg/commands/barcode"
)

// BarcodeConfig agrupa todos los parámetros visuales de un código de barras
type BarcodeConfig struct {
	Symbology   barcode.Symbology
	Width       barcode.Width       // Ancho del módulo
	Height      barcode.Height      // Altura en dots
	HRIPosition barcode.HRIPosition // Posición del texto
	HRIFont     barcode.HRIFont     // Fuente del texto
	CodeSet     barcode.Code128Set  // Opcional: Específico para CODE128 manual
}

// DefaultConfig devuelve una configuración segura
func DefaultConfig() BarcodeConfig {
	return BarcodeConfig{
		Symbology:   barcode.UPCA,
		Width:       barcode.DefaultWidth,
		Height:      barcode.DefaultHeight,
		HRIPosition: barcode.HRIBelow,
		HRIFont:     barcode.HRIFontA,
		CodeSet:     barcode.Code128SetB,
	}
}
