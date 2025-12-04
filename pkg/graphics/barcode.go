package graphics

import (
	"fmt"
	"strings"

	"github.com/adcondev/pos-printer/pkg/commands/barcode"
	"github.com/adcondev/pos-printer/pkg/constants"
)

var (
	// TODO: Agregar soporte para más simbologías si el paquete barcode lo permite

	// symMap mapea nombres de simbologías a constantes del paquete barcode
	symMap = map[string]barcode.Symbology{
		constants.UPCA.String():    barcode.UPCA,
		constants.UPCE.String():    barcode.UPCE,
		constants.EAN13.String():   barcode.EAN13,
		constants.EAN8.String():    barcode.EAN8,
		constants.CODE39.String():  barcode.CODE39,
		constants.ITF.String():     barcode.ITF,
		constants.CODE128.String(): barcode.CODE128,
		constants.CODABAR.String(): barcode.CODABAR,
		// Agregar soporte para sinónimos comunes
		constants.JAN13.String(): barcode.EAN13, // JAN13 es equivalente a EAN13
		constants.JAN8.String():  barcode.EAN8,  // JAN8 es equivalente a EAN8
	}

	// hriPosMap mapea posiciones HRI a constantes del paquete barcode
	hriPosMap = map[string]barcode.HRIPosition{
		constants.None.String():  barcode.HRINotPrinted,
		constants.Above.String(): barcode.HRIAbove,
		constants.Below.String(): barcode.HRIBelow,
		constants.Both.String():  barcode.HRIBoth,
	}

	// hriFontMap mapea fuentes HRI a constantes del paquete barcode
	hriFontMap = map[string]barcode.HRIFont{
		constants.A.String(): barcode.HRIFontA,
		constants.B.String(): barcode.HRIFontB,
	}
)

// MapSymbology mapea los nombres del schema a las constantes del paquete barcode
func MapSymbology(sym string) (barcode.Symbology, error) {
	// Normalizar el input
	normalized := strings.ToLower(strings.ReplaceAll(sym, "-", ""))

	symbology, ok := symMap[normalized]
	if !ok {
		// Listar simbologías soportadas en el error
		return 0, fmt.Errorf("unsupported barcode symbology: %s (supported: UPC-A, UPC-E, EAN13, EAN8, CODE39, CODE128, ITF, CODABAR)", sym)
	}

	return symbology, nil
}

// MapHriPosition mapea las posiciones HRI del schema a las constantes del paquete barcode
func MapHriPosition(pos string) (barcode.HRIPosition, error) {
	normalized := strings.ToLower(pos)

	hriPos, ok := hriPosMap[normalized]
	if !ok {
		return 0, fmt.Errorf("unsupported HRI position: %s (supported: none, above, below, both)", pos)
	}

	return hriPos, nil
}

// MapHriFont mapea las fuentes HRI del schema a las constantes del paquete barcode
func MapHriFont(font string) (barcode.HRIFont, error) {
	normalized := strings.ToLower(font)

	hriFont, ok := hriFontMap[normalized]
	if !ok {
		return 0, fmt.Errorf("unsupported HRI font: %s (supported: A, B)", font)
	}

	return hriFont, nil
}

// IntToBarcodeModuleWidth convierte un entero a barcode.Width
func IntToBarcodeModuleWidth(w int) barcode.Width {
	switch {
	case w <= constants.MinBarcodeModuleWidth:
		return barcode.MinWidth
	case w >= constants.MaxBarcodeModuleWidth:
		return barcode.MaxWidth
	}
	return barcode.Width(w)
}

// IntToBarcodeHeight convierte un entero a barcode.Height
func IntToBarcodeHeight(h int) barcode.Height {
	switch {
	case h <= constants.MinBarcodeHeight:
		return barcode.MinHeight
	case h >= constants.MaxBarcodeHeight:
		return barcode.MaxHeight
	}
	return barcode.Height(h)
}

// BarcodeConfig agrupa todos los parámetros visuales de un código de barras
type BarcodeConfig struct {
	Symbology   barcode.Symbology
	Width       barcode.Width       // Ancho del módulo
	Height      barcode.Height      // Altura en dots
	HRIPosition barcode.HRIPosition // Posición del texto
	HRIFont     barcode.HRIFont     // Fuente del texto
	// TODO: Check CodeSet only applies to CODE128 and do it automatically
	CodeSet barcode.Code128Set // Opcional: Específico para CODE128 manual
}

// DefaultBarcodeConfig devuelve una configuración segura
func DefaultBarcodeConfig() BarcodeConfig {
	return BarcodeConfig{
		Symbology:   symMap[constants.CODE128.String()],
		Width:       IntToBarcodeModuleWidth(constants.DefaultBarcodeModuleWidth),
		Height:      IntToBarcodeHeight(constants.DefaultBarcodeHeight),
		HRIPosition: hriPosMap[constants.DefaultBarcodeHriPosition.String()],
		HRIFont:     hriFontMap[constants.DefaultBarcodeHriFont.String()],
		CodeSet:     barcode.Code128SetB,
	}
}
