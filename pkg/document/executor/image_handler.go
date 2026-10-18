package executor

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/adcondev/poster/internal/load"
	"github.com/adcondev/poster/pkg/constants"
	"github.com/adcondev/poster/pkg/document/schema"
	"github.com/adcondev/poster/pkg/graphics"
	"github.com/adcondev/poster/pkg/service"
)

// ImageCommand for image handler
type ImageCommand = schema.ImageCommand

// handleImage manages image commands
func (e *Executor) handleImage(printer service.PrinterActions, data json.RawMessage) error {
	var cmd ImageCommand
	if err := json.Unmarshal(data, &cmd); err != nil {
		return fmt.Errorf("failed to parse image command: %w", err)
	}

	img, _, err := load.ImgFromBase64(cmd.Code)
	if err != nil {
		return fmt.Errorf("failed to load image: %w", err)
	}

	opts := &graphics.ImgOptions{
		PixelWidth:     cmd.PixelWidth,
		Threshold:      cmd.Threshold,
		PreserveAspect: true,
		AutoRotate:     false,
	}

	if opts.PixelWidth <= 0 {
		opts.PixelWidth = constants.DefaultImagePixelWidth
	}

	if prof := printer.GetProfile(); prof != nil && prof.DotsPerLine > 0 {
		opts.MaxWidth = prof.DotsPerLine
		if opts.PixelWidth > opts.MaxWidth {
			opts.PixelWidth = opts.MaxWidth
		}
	}

	if opts.Threshold == 0 {
		opts.Threshold = constants.DefaultImageThreshold
	}

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

	switch strings.ToLower(cmd.Align) {
	case constants.Center.String():
		if err := printer.AlignCenter(); err != nil {
			return err
		}
	case constants.Right.String():
		if err := printer.AlignRight(); err != nil {
			return err
		}
	default:
		if err := printer.AlignLeft(); err != nil {
			return err
		}
	}

	// Accept both "nns" (schema enum) and "nearest_neighbor"
	switch strings.ToLower(cmd.Scaling) {
	case "nns", constants.NearestNeighbor.String():
		opts.Scaling = graphics.NearestNeighbor
	case constants.Bilinear.String():
		opts.Scaling = graphics.BiLinear
	default:
		opts.Scaling = graphics.ScaleMap[constants.DefaultImageScaling]
	}

	pipeline := graphics.NewPipeline(opts)
	bitmap, err := pipeline.Process(img)
	if err != nil {
		return fmt.Errorf("failed to process image: %w", err)
	}

	if err := printer.PrintBitmap(bitmap); err != nil {
		return fmt.Errorf("failed to print bitmap: %w", err)
	}

	return printer.AlignLeft()
}
