package executor

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/adcondev/poster/pkg/commands/barcode"
	"github.com/adcondev/poster/pkg/constants"
	"github.com/adcondev/poster/pkg/graphics"
	"github.com/adcondev/poster/pkg/service"
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

	// Crear configuración con defaults
	cfg := graphics.DefaultBarcodeConfig()

	// Mapear la simbología del schema a la constante del paquete barcode
	symbology, err := graphics.MapSymbology(cmd.Symbology)
	if err != nil {
		symbology = barcode.CODE128 // Fallback seguro
		log.Printf("Unsupported barcode symbology '%s', defaulting to CODE128", cmd.Symbology)
	}
	cfg.Symbology = symbology

	// Aplicar configuración de ancho si se especifica
	if cmd.Width != nil {
		switch {
		case *cmd.Width < constants.MinBarcodeModuleWidth:
			cfg.Width = constants.MinBarcodeModuleWidth
			log.Printf("Barcode width too small, set to minimum %d", *cmd.Width)
		case *cmd.Width > constants.MaxBarcodeModuleWidth:
			cfg.Width = constants.MaxBarcodeModuleWidth
			log.Printf("Barcode width too large, set to maximum %d", *cmd.Width)
		default:
			cfg.Width = barcode.Width(*cmd.Width)
			log.Printf("Barcode width set to %d", *cmd.Width)
		}
	} else {
		// Si no se especifica, usar valor por defecto
		cfg.Width = constants.DefaultBarcodeModuleWidth
		log.Printf("Barcode width using default %d", cfg.Width)
	}

	// Aplicar configuración de altura si se especifica
	if cmd.Height != nil {
		switch {
		case *cmd.Height < constants.MinBarcodeHeight:
			cfg.Height = constants.MinBarcodeHeight
			log.Printf("Barcode height too small, set to minimum %d dots", *cmd.Height)
		case *cmd.Height > constants.MaxBarcodeHeight:
			cfg.Height = constants.MaxBarcodeHeight
			log.Printf("Barcode height too large, set to maximum %d dots", *cmd.Height)
		default:
			cfg.Height = barcode.Height(*cmd.Height)
			log.Printf("Barcode height set to %d dots", *cmd.Height)
		}
	} else {
		// Si no se especifica, usar valor por defecto
		cfg.Height = constants.DefaultBarcodeHeight
		log.Printf("Barcode height using default %d dots", cfg.Height)
	}

	// Configurar posición del HRI (Human Readable Interpretation)
	hriPos := constants.DefaultBarcodeHriPosition // default según schema
	if cmd.HRIPosition != nil {
		hriPos = constants.HriPosition(strings.ToLower(*cmd.HRIPosition))
	}
	cfg.HRIPosition = mapHRIPosition(hriPos)

	// Configurar fuente del HRI
	hriFont := constants.DefaultBarcodeHriFont // default según schema
	if cmd.HRIFont != nil {
		hriFont = constants.Font(strings.ToLower(*cmd.HRIFont))
	}
	cfg.HRIFont = mapHRIFont(hriFont)

	// Para CODE128, establecer el CodeSet por defecto si es necesario
	if symbology == barcode.CODE128 || symbology == barcode.GS1128 {
		// Validar que los datos sean compatibles con el code set B por defecto
		cfg.CodeSet = barcode.Code128SetB // Default para texto general

		// Si los datos contienen solo números y tienen longitud par, usar Code Set C
		if barcode.ValidateNumericData([]byte(cmd.Data)) && len(cmd.Data)%2 == 0 {
			cfg.CodeSet = barcode.Code128SetC
			log.Printf("CODE128 data is numeric with even length, using Code Set C")
		}
	}

	// Aplicar alineación
	align := constants.DefaultBarcodeAlignment // default según schema
	if cmd.Align != nil {
		align = constants.Alignment(strings.ToLower(*cmd.Align))
	}

	switch align {
	case constants.Center:
		if err := printer.AlignCenter(); err != nil {
			return err
		}
	case constants.Right:
		if err := printer.AlignRight(); err != nil {
			return err
		}
	case constants.Left:
		if err := printer.AlignLeft(); err != nil {
			return err
		}
	default:
		if err := printer.SetAlignment(constants.DefaultTextAlignment.String()); err != nil {
			return err
		}
	}

	// Imprimir el código de barras usando el méetodo stateless
	if err := printer.PrintBarcode(cfg, []byte(cmd.Data)); err != nil {
		return fmt.Errorf("failed to print barcode: %w", err)
	}

	// Restaurar alineación a la izquierda
	if err := printer.SetAlignment(constants.DefaultTextAlignment.String()); err != nil {
		return err
	}

	return nil
}

// mapHRIPosition mapea las posiciones del schema a las constantes del paquete barcode
func mapHRIPosition(pos constants.HriPosition) barcode.HRIPosition {
	posMap := map[constants.HriPosition]barcode.HRIPosition{
		constants.None:  barcode.HRINotPrinted,
		constants.Above: barcode.HRIAbove,
		constants.Below: barcode.HRIBelow,
		constants.Both:  barcode.HRIBoth,
	}

	hriPos, ok := posMap[pos]
	if !ok {
		return posMap[constants.DefaultBarcodeHriPosition] // Default
	}

	return hriPos
}

// mapHRIFont mapea las fuentes del schema a las constantes del paquete barcode
func mapHRIFont(font constants.Font) barcode.HRIFont {
	fontMap := map[constants.Font]barcode.HRIFont{
		constants.A: barcode.HRIFontA,
		constants.B: barcode.HRIFontB,
	}

	hriFont, ok := fontMap[font]
	if !ok {
		return fontMap[constants.DefaultBarcodeHriFont] // Default
	}

	return hriFont
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
