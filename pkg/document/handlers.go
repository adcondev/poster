package document

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/adcondev/pos-printer/internal/load"
	posqr "github.com/adcondev/pos-printer/pkg/commands/qrcode"
	"github.com/adcondev/pos-printer/pkg/graphics"
	"github.com/adcondev/pos-printer/pkg/service"
	"github.com/adcondev/pos-printer/pkg/tables"
)

const (
	center = "center"
	right  = "right"
	left   = "left"
)

// handleImage manages image commands
func (e *Executor) handleImage(printer *service.Printer, data json.RawMessage) error {
	var cmd ImageCommand
	if err := json.Unmarshal(data, &cmd); err != nil {
		return fmt.Errorf("failed to parse image command: %w", err)
	}

	// Decodificar imagen desde base64
	img, format, err := load.ImgFromBase64(cmd.Code)
	if err != nil {
		return fmt.Errorf("failed to load image: %w", err)
	}

	cmd.Format = format
	log.Printf("Loaded image with format: %s", format)

	// Configurar opciones de procesamiento
	opts := &graphics.ImgOptions{
		PixelWidth:     cmd.PixelWidth,
		Threshold:      cmd.Threshold,
		PreserveAspect: true,
		AutoRotate:     false,
	}

	// Si no se especifica ancho, usar el ancho del perfil
	if opts.PixelWidth == 0 {
		opts.PixelWidth = 256
	}

	// Si no se especifica threshold, usar valor por defecto
	if opts.Threshold == 0 {
		opts.Threshold = 128
	}

	// Configurar dithering
	switch strings.ToLower(cmd.Dithering) {
	case "atkinson":
		opts.Dithering = graphics.Atkinson
	default:
		opts.Dithering = graphics.Threshold
	}

	// Aplicar alineación
	switch strings.ToLower(cmd.Align) {
	case center:
		err := printer.AlignCenter()
		if err != nil {
			return err
		}
	case right:
		err := printer.AlignRight()
		if err != nil {
			return err
		}
	default:
		err := printer.AlignLeft()
		if err != nil {
			return err
		}
	}

	// Procesar imagen
	pipeline := graphics.NewPipeline(opts)
	bitmap, err := pipeline.Process(img)
	if err != nil {
		return fmt.Errorf("failed to process image: %w", err)
	}

	// Imprimir bitmap
	if err := printer.PrintBitmap(bitmap); err != nil {
		return fmt.Errorf("failed to print bitmap: %w", err)
	}

	// Resetear alineación
	err = printer.AlignLeft()
	if err != nil {
		return err
	}

	return nil
}

// handleSeparator manages separator commands
func (e *Executor) handleSeparator(printer *service.Printer, data json.RawMessage) error {
	var cmd SeparatorCommand
	if err := json.Unmarshal(data, &cmd); err != nil {
		return fmt.Errorf("failed to parse separator command: %w", err)
	}

	// Valores por defecto
	if cmd.Char == "" {
		cmd.Char = "- "
	}
	if cmd.Length == 0 {
		// Usar ancho del papel en caracteres (aproximado)
		// TODO: Verify the following line for different fonts
		cmd.Length = e.printer.Profile.DotsPerLine / 12 // Aproximación para Font A
	}

	// Construir línea separadora
	line := strings.Repeat(cmd.Char, cmd.Length)

	return printer.PrintLine(line)
}

// handleFeed manages feed commands
func (e *Executor) handleFeed(printer *service.Printer, data json.RawMessage) error {
	var cmd FeedCommand
	if err := json.Unmarshal(data, &cmd); err != nil {
		return fmt.Errorf("failed to parse feed command: %w", err)
	}

	if cmd.Lines <= 0 {
		cmd.Lines = 1
	}

	return printer.FeedLines(byte(cmd.Lines))
}

// handleCut manages cut commands
func (e *Executor) handleCut(printer *service.Printer, data json.RawMessage) error {
	var cmd CutCommand
	if err := json.Unmarshal(data, &cmd); err != nil {
		return fmt.Errorf("failed to parse cut command: %w", err)
	}

	// Avance antes del corte si se especifica
	if cmd.Feed > 0 {
		err := printer.FeedLines(byte(cmd.Feed))
		if err != nil {
			return err
		}
	}

	// Ejecutar corte
	switch strings.ToLower(cmd.Mode) {
	case "full":
		return printer.FullFeedAndCut(0)
	default: // partial
		return printer.PartialFeedAndCut(0)
	}
}

// TODO: Manage text_under and text_above options instead of human_text

// handleQR manges QR code commands
func (e *Executor) handleQR(printer *service.Printer, data json.RawMessage) error {
	var cmd QRCommand
	if err := json.Unmarshal(data, &cmd); err != nil {
		return fmt.Errorf("failed to parse QR command: %w", err)
	}

	// Validación de datos requeridos
	if cmd.Data == "" {
		return fmt.Errorf("QR data cannot be empty")
	}
	if len(cmd.Data) > posqr.MaxDataLength {
		return fmt.Errorf("QR data too long: %d bytes (maximum %d)", len(cmd.Data), posqr.MaxDataLength)
	}

	// Construir opciones con defaults
	opts := graphics.DefaultQROptions()

	// Configurar tamaño del QR
	if cmd.PixelWidth > 0 {
		if cmd.PixelWidth < 87 { // Mínimo según schema
			cmd.PixelWidth = 87
			log.Printf("Warning: pixel_width too small, setting to minimum 87")
		}
		opts.PixelWidth = cmd.PixelWidth
	}

	// Mapear corrección de errores (default: Q según schema)
	correction := "Q"
	if cmd.Correction != "" {
		correction = strings.ToUpper(cmd.Correction)
	}

	switch correction {
	case "L":
		opts.ErrorCorrection = posqr.LevelL
	case "M":
		opts.ErrorCorrection = posqr.LevelM
	case "Q":
		opts.ErrorCorrection = posqr.LevelQ
	case "H":
		opts.ErrorCorrection = posqr.LevelH
	default:
		return fmt.Errorf("invalid QR correction level: %s (valid: L, M, Q, H)", cmd.Correction)
	}

	// Configurar logo si se proporciona
	if cmd.Logo != "" {
		opts.LogoData = cmd.Logo
		log.Printf("QR with logo enabled")
	}

	// Configurar forma circular (solo para QR grandes)
	if cmd.CircleShape {
		if opts.PixelWidth > 256 {
			opts.CircleShape = true
			log.Printf("QR with circular blocks enabled")
		} else {
			log.Printf("Warning: circle_shape ignored (requires pixel_width > 256)")
		}
	}

	// Aplicar alineación (default: center según schema)
	align := "center"
	if cmd.Align != "" {
		align = strings.ToLower(cmd.Align)
	}

	switch align {
	case center:
		if err := printer.AlignCenter(); err != nil {
			return err
		}
	case right:
		if err := printer.AlignRight(); err != nil {
			return err
		}
	case left:
		if err := printer.AlignLeft(); err != nil {
			return err
		}
	default:
		if err := printer.AlignCenter(); err != nil {
			return err
		}
	}

	// Imprimir QR (automáticamente intenta nativo y fallback a imagen)
	if err := printer.PrintQR(cmd.Data, opts); err != nil {
		return fmt.Errorf("failed to print QR: %w", err)
	}

	// Imprimir texto humano si existe
	if cmd.HumanText != "" {
		// Mantener la alineación del QR para el texto
		if err := printer.PrintLine(cmd.HumanText); err != nil {
			return fmt.Errorf("failed to print QR human text: %w", err)
		}
	}

	// Restaurar alineación a la izquierda
	if err := printer.AlignLeft(); err != nil {
		return err
	}

	return nil
}

// TODO: Consider a title fields for tables as field in TableCommand

// handleTable manages table commands
func (e *Executor) handleTable(printer *service.Printer, data json.RawMessage) error {
	var cmd TableCommand
	if err := json.Unmarshal(data, &cmd); err != nil {
		return fmt.Errorf("failed to parse table command: %w", err)
	}

	if cmd.Options.ColumnSpacing < 0 {
		// Validar que ColumnSpacing no sea negativo
		cmd.Options.ColumnSpacing = 0
		log.Printf("ColumnSpacing cannot be negative, using 0")
	}
	// Validate table command
	if len(cmd.Definition.Columns) == 0 {
		return fmt.Errorf("table must have at least one column defined")
	}

	// Create table options with defaults
	opts := &tables.Options{
		ShowHeaders:   cmd.ShowHeaders,
		WordWrap:      true,
		ColumnSpacing: 1,
		HeaderStyle:   tables.Style{Bold: true},
	}

	// Apply custom options if provided
	if cmd.Options != nil {
		opts.WordWrap = cmd.Options.WordWrap
		if cmd.Options.HeaderBold {
			opts.HeaderStyle.Bold = true
		}
		if cmd.Options.ColumnSpacing > 0 {
			opts.ColumnSpacing = cmd.Options.ColumnSpacing
		}
	}

	// Set paper width
	switch {
	case cmd.Definition.PaperWidth > 0:
		opts.PaperWidth = cmd.Definition.PaperWidth
	case printer.Profile.PrintWidth > 0:
		opts.PaperWidth = printer.Profile.PrintWidth
	default:
		if printer.Profile.PaperWidth >= 80 {
			opts.PaperWidth = tables.Width80mm203dpi
		} else {
			opts.PaperWidth = tables.Width58mm203dpi
		}
	}

	// Create table engine
	engine := tables.NewEngine(&cmd.Definition, opts)

	// Prepare table data
	tableData := &tables.Data{
		Definition:  cmd.Definition,
		ShowHeaders: cmd.ShowHeaders,
		Rows:        make([]tables.Row, len(cmd.Rows)),
	}

	// Convert rows
	for i, row := range cmd.Rows {
		tableData.Rows[i] = row
	}

	// Render table to string
	var buf strings.Builder
	if err := engine.Render(&buf, tableData); err != nil {
		return fmt.Errorf("failed to render table: %w", err)
	}

	// Aplicar alineación
	switch strings.ToLower(cmd.Options.Align) {
	case center:
		err := printer.AlignCenter()
		if err != nil {
			return err
		}
	case right:
		err := printer.AlignRight()
		if err != nil {
			return err
		}
	default:
		err := printer.AlignLeft()
		if err != nil {
			return err
		}
	}

	err := printer.Print(buf.String())
	if err != nil {
		return err
	}

	// Restaurar alineación
	err = printer.AlignLeft()
	if err != nil {
		return err
	}

	// Send the raw output (includes ESC/POS commands for bold)
	return nil
}
