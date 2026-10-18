package executor

import (
	"strings"

	"github.com/adcondev/poster/pkg/constants"
	"github.com/adcondev/poster/pkg/service"
)

func (e *Executor) applyAlign(printer service.PrinterActions, align *string) error {
	alignValue := constants.DefaultTextAlignment.String()
	if align != nil {
		alignValue = strings.ToLower(*align)
	}

	switch alignValue {
	case constants.Center.String():
		return printer.AlignCenter()
	case constants.Right.String():
		return printer.AlignRight()
	case constants.Left.String():
		return printer.AlignLeft()
	default:
		return printer.AlignLeft()
	}
}

func (e *Executor) applySize(printer service.PrinterActions, size string) error {
	switch ss := strings.ToLower(size); ss {
	case constants.Normal.String():
		return printer.SingleSize()
	case constants.Double.String():
		return printer.DoubleSize()
	case constants.Triple.String():
		return printer.TripleSize()
	case constants.Quad.String():
		return printer.QuadraSize()
	case constants.Penta.String():
		return printer.PentaSize()
	case constants.Hexa.String():
		return printer.HexaSize()
	case constants.Hepta.String():
		return printer.HeptaSize()
	case constants.Octa.String():
		return printer.OctaSize()
	default:
		if len(ss) == 3 && ss[1] == 'x' {
			parts := strings.Split(ss, "x")
			widthMultiplier := parts[0][0] - '0'
			heightMultiplier := parts[1][0] - '0'
			if widthMultiplier >= 1 && widthMultiplier <= 8 &&
				heightMultiplier >= 1 && heightMultiplier <= 8 {
				return printer.CustomSize(widthMultiplier, heightMultiplier)
			}
		}
		return printer.SingleSize()
	}
}

func (e *Executor) applyUnderline(printer service.PrinterActions, underline string) error {
	switch strings.ToLower(underline) {
	case constants.NoDot.String():
		return printer.NoDot()
	case constants.OneDot.String():
		return printer.OneDot()
	case constants.TwoDot.String():
		return printer.TwoDot()
	default:
		return printer.NoDot()
	}
}

func (e *Executor) applyFont(printer service.PrinterActions, font string) error {
	switch strings.ToLower(font) {
	case "", "a":
		return printer.FontA()
	case "b":
		return printer.FontB()
	default:
		return printer.FontA()
	}
}
