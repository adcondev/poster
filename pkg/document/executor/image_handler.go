package executor

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/adcondev/pos-printer/internal/load"
	"github.com/adcondev/pos-printer/pkg/constants"
	"github.com/adcondev/pos-printer/pkg/graphics"
	"github.com/adcondev/pos-printer/pkg/service"
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
	log.Printf("Loaded image with format: %s", format)

	// Configurar opciones de procesamiento
	opts := &graphics.ImgOptions{
		PixelWidth:     cmd.PixelWidth,
		Threshold:      cmd.Threshold,
		PreserveAspect: true,
		AutoRotate:     false,
	}

	// Si no se especifica ancho, usar valor por defecto
	if opts.PixelWidth <= 0 {
		opts.PixelWidth = constants.DefaultImagePixelWidth
	}

	// Si no se especifica threshold, usar valor por defecto
	if opts.Threshold <= 0 {
		opts.Threshold = constants.DefaultImageThreshold
	}

	// Configurar dithering
	switch strings.ToLower(cmd.Dithering) {
	case constants.Atkinson.String():
		opts.Dithering = graphics.Atkinson
	case constants.Threshold.String():
		opts.Dithering = graphics.Threshold
	default:
		opts.Dithering = graphics.DitherMap[constants.DefaultImageDithering.String()]
	}

	if cmd.Align == "" {
		cmd.Align = constants.DefaultImageAlignment.String()
	}

	// Aplicar alineación
	switch strings.ToLower(cmd.Align) {
	case constants.Center.String():
		err := printer.AlignCenter()
		if err != nil {
			return err
		}
	case constants.Right.String():
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

	// Configurar escalado
	switch strings.ToLower(cmd.Scaling) {
	case constants.NearestNeighbor.String():
		opts.Scaling = graphics.NearestNeighbor
	case constants.Bilinear.String():
		opts.Scaling = graphics.BiLinear
	default:
		opts.Scaling = graphics.ScaleMap[constants.DefaultImageScaling]
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
