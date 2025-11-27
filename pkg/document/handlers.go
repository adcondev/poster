package document

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/adcondev/pos-printer/internal/load"
	"github.com/adcondev/pos-printer/pkg/graphics"
	"github.com/adcondev/pos-printer/pkg/service"
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
		cmd.Char = "-"
	}
	if cmd.Length == 0 {
		// Usar ancho del papel en caracteres (aproximado)
		// TODO: Verify the following line for different fonts, constrained to Font A
		cmd.Length = e.printer.Profile.DotsPerLine / 12 // Aproximación para Font A
	}

	err := printer.AlignCenter()
	if err != nil {
		return err
	}

	// Construir línea separadora
	line := strings.Repeat(cmd.Char, cmd.Length)
	err = printer.PrintLine(line[:cmd.Length-1]) // Trim to exact length
	if err != nil {
		return err
	}

	err = printer.AlignLeft()
	if err != nil {
		return err
	}

	return nil
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

	// Default feed antes del corte según schema
	if cmd.Feed == 0 {
		cmd.Feed = 2 // Default según schema
	}

	// Avance antes del corte si se especifica
	if cmd.Feed > 0 {
		if err := printer.FeedLines(byte(cmd.Feed)); err != nil {
			return err
		}
	}

	// TODO: Implement field Dots, as this command make the feed based lines only.
	// Instead of 0 for Full and Partial cut, it should be the number of dots to feed before cutting.
	// Maybe calculation from mm to dpl is needed.

	// Ejecutar corte
	switch strings.ToLower(cmd.Mode) {
	case "full":
		return printer.FullFeedAndCut(0)
	default: // partial
		return printer.PartialFeedAndCut(0)
	}
}
