package executor

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/adcondev/pos-printer/pkg/commands/barcode"
	"github.com/adcondev/pos-printer/pkg/constants"
	"github.com/adcondev/pos-printer/pkg/graphics"
	"github.com/adcondev/pos-printer/pkg/service"
)

// handleBarcode manages barcode commands
func (e *Executor) handleBarcode(printer *service.Printer, data json.RawMessage) error {
	var cmd BarcodeCommand
	if err := json.Unmarshal(data, &cmd); err != nil {
		return fmt.Errorf("failed to parse barcode command: %w", err)
	}

	// Validación de campos requeridos
	if cmd.Symbology == "" {
		return fmt.Errorf("barcode symbology is required")
	}
	if cmd.Data == "" {
		return fmt.Errorf("barcode data is required")
	}

	e.analyzeBarcodeRisk(cmd.Data)

	// Mapear la simbología del schema a la constante del paquete barcode
	symbology, err := mapSymbology(cmd.Symbology)
	if err != nil {
		return err
	}

	// Crear configuración con defaults
	cfg := graphics.DefaultConfig()
	cfg.Symbology = symbology

	// Aplicar configuración de ancho si se especifica
	if cmd.Width != nil {
		if *cmd.Width < 2 || *cmd.Width > 6 {
			return fmt.Errorf("barcode width must be between 2 and 6, got %d", *cmd.Width)
		}
		cfg.Width = barcode.Width(*cmd.Width)
		log.Printf("Barcode width set to %d", *cmd.Width)
	}

	// Aplicar configuración de altura si se especifica
	if cmd.Height != nil {
		if *cmd.Height < 1 || *cmd.Height > 255 {
			return fmt.Errorf("barcode height must be between 1 and 255, got %d", *cmd.Height)
		}
		cfg.Height = barcode.Height(*cmd.Height)
		log.Printf("Barcode height set to %d dots", *cmd.Height)
	}

	// Configurar posición del HRI (Human Readable Interpretation)
	hriPos := "below" // default según schema
	if cmd.HRIPosition != nil {
		hriPos = strings.ToLower(*cmd.HRIPosition)
	}

	hriPosition, err := mapHRIPosition(hriPos)
	if err != nil {
		return err
	}
	cfg.HRIPosition = hriPosition

	// Configurar fuente del HRI
	hriFont := "A" // default según schema
	if cmd.HRIFont != nil {
		hriFont = strings.ToUpper(*cmd.HRIFont)
	}

	hriF, err := mapHRIFont(hriFont)
	if err != nil {
		return err
	}
	cfg.HRIFont = hriF

	// Para CODE128, establecer el CodeSet por defecto si es necesario
	if symbology == barcode.CODE128 || symbology == barcode.GS1128 {
		// Validar que los datos sean compatibles con el code set B por defecto
		cfg.CodeSet = barcode.Code128SetB // Default para texto general

		// Si los datos contienen solo números y tienen longitud par, sugerir Code Set C
		if barcode.ValidateNumericData([]byte(cmd.Data)) && len(cmd.Data)%2 == 0 {
			log.Printf("warning: consider using Code Set C for numeric data '%s' for better density", cmd.Data)
		}
	}

	// Aplicar alineación
	align := "center" // default según schema
	if cmd.Align != nil {
		align = strings.ToLower(*cmd.Align)
	}

	switch align {
	case constants.AlignCenter.String():
		if err := printer.AlignCenter(); err != nil {
			return err
		}
	case constants.AlignRight.String():
		if err := printer.AlignRight(); err != nil {
			return err
		}
	case constants.AlignLeft.String():
		if err := printer.AlignLeft(); err != nil {
			return err
		}
	default:
		if err := printer.AlignCenter(); err != nil {
			return err
		}
	}

	// Imprimir el código de barras usando el método stateless
	if err := printer.PrintBarcode(cfg, []byte(cmd.Data)); err != nil {
		return fmt.Errorf("failed to print barcode: %w", err)
	}

	// Restaurar alineación a la izquierda
	if err := printer.AlignLeft(); err != nil {
		return err
	}

	return nil
}

// mapSymbology mapea los nombres del schema a las constantes del paquete barcode
func mapSymbology(sym string) (barcode.Symbology, error) {
	// Normalizar el input
	normalized := strings.ToUpper(strings.ReplaceAll(sym, "-", ""))

	symbologyMap := map[string]barcode.Symbology{
		"UPCA":    barcode.UPCAB,    // Function B version para mayor compatibilidad
		"UPCE":    barcode.UPCEB,    // Function B version
		"EAN13":   barcode.EAN13,    // Function B version (ya está en 67)
		"EAN8":    barcode.EAN8,     // Function B version (ya está en 68)
		"CODE39":  barcode.CODE39B,  // Function B version
		"CODE128": barcode.CODE128,  // Function B con code set
		"ITF":     barcode.ITFB,     // Function B version
		"CODABAR": barcode.CODABARB, // Function B version
		// Agregar soporte para sinónimos comunes
		"JAN13": barcode.EAN13, // JAN13 es equivalente a EAN13
		"JAN8":  barcode.EAN8,  // JAN8 es equivalente a EAN8
	}

	symbology, ok := symbologyMap[normalized]
	if !ok {
		// Listar simbologías soportadas en el error
		return 0, fmt.Errorf("unsupported barcode symbology: %s (supported: UPC-A, UPC-E, EAN13, EAN8, CODE39, CODE128, ITF, CODABAR)", sym)
	}

	return symbology, nil
}

// mapHRIPosition mapea las posiciones del schema a las constantes del paquete barcode
func mapHRIPosition(pos string) (barcode.HRIPosition, error) {
	posMap := map[string]barcode.HRIPosition{
		"none":  barcode.HRINotPrinted,
		"above": barcode.HRIAbove,
		"below": barcode.HRIBelow,
		"both":  barcode.HRIBoth,
	}

	hriPos, ok := posMap[pos]
	if !ok {
		return 0, fmt.Errorf("invalid HRI position: %s (valid: none, above, below, both)", pos)
	}

	return hriPos, nil
}

// mapHRIFont mapea las fuentes del schema a las constantes del paquete barcode
func mapHRIFont(font string) (barcode.HRIFont, error) {
	fontMap := map[string]barcode.HRIFont{
		"A": barcode.HRIFontA,
		"B": barcode.HRIFontB,
	}

	hriFont, ok := fontMap[font]
	if !ok {
		return 0, fmt.Errorf("invalid HRI font: %s (valid: A, B)", font)
	}

	return hriFont, nil
}

// analyzeBarcodeRisk evalúa la longitud de los datos contra el ancho del papel
// y emite logs de advertencia o error según la física de impresión térmica.
func (e *Executor) analyzeBarcodeRisk(data string) {
	length := len(data)
	// Obtenemos el ancho del papel del perfil (o default a 80 si es 0)
	paperWidth := e.printer.Profile.PaperWidth
	if paperWidth == 0 {
		paperWidth = 80
	}

	const (
		LimitPhysicalMax = 30 // Límite absoluto físico
		LimitSafe58mm    = 12 // Límite seguro para 58mm
		LimitSafe80mm    = 22 // Límite seguro para 80mm
	)

	// Lógica de Logs (La Tabla de Decisión)
	switch {
	case length > LimitPhysicalMax:
		log.Printf("[CRITICAL] BARCODE OVERFLOW: Data length (%d) exceeds physical limit (%d). 1D Barcode will be unscannable. ACTION: FORCING QR STRONGLY RECOMMENDED.", length, LimitPhysicalMax)

	case paperWidth <= 58 && length > LimitSafe58mm:
		log.Printf("[WARNING] BARCODE RISK: Data length (%d) is unsafe for 58mm paper (Safe limit: %d). Code may be clipped. RECOMMENDATION: Use QR.", length, LimitSafe58mm)

	case length > LimitSafe80mm:
		log.Printf("[WARNING] BARCODE DENSITY: Data length (%d) is high for 1D barcode (Safe limit: %d). Scan reliability may drop.", length, LimitSafe80mm)

	default:
		// Opcional: Log positivo si toodo está bien (útil para debug mode)
		if e.printer.Profile.DebugLog { // Asumiendo que tengas un flag de debug
			log.Printf("[INFO] Barcode size safe: %d chars on %vmm paper.", length, paperWidth)
		}
	}
}
