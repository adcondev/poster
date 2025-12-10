package executor

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/adcondev/poster/pkg/constants"
	"github.com/adcondev/poster/pkg/service"
)

// handleText manages text commands
func (e *Executor) handleText(printer *service.Printer, data json.RawMessage) error {
	var cmd TextCommand
	if err := json.Unmarshal(data, &cmd); err != nil {
		return fmt.Errorf("failed to parse text command: %w", err)
	}

	// Procesar label si existe
	if cmd.Label != nil && cmd.Label.Text != "" {
		// Aplicar estilo del label solo si existe
		if cmd.Label.Style != nil {
			if err := e.applyTextStyle(printer, cmd.Label.Style); err != nil {
				return fmt.Errorf("failed to apply label style: %w", err)
			}
		}

		// Construir texto del label
		labelText := cmd.Label.Text
		if cmd.Label.Separator != nil && *cmd.Label.Separator != "" {
			labelText += *cmd.Label.Separator
		} else {
			labelText += ": "
		}

		// Aplicar alineación del label
		if err := e.applyAlign(printer, cmd.Label.Align); err != nil {
			return err
		}

		// Si las alineaciones son diferentes, imprimir label y resetear
		if !sameOrNil(cmd.Label.Align, cmd.Content.Align) {
			// Imprimir label con salto de línea
			if err := printer.PrintLine(labelText); err != nil {
				return err
			}
			// Reset completo antes del content solo si había estilo
			if cmd.Label.Style != nil {
				if err := e.resetTextStyle(printer, cmd.Label.Style); err != nil {
					return err
				}
			}
		} else {
			// Imprimir label sin salto
			if err := printer.Print(labelText); err != nil {
				return err
			}
			// Reset solo los estilos diferentes si ambos existen
			if cmd.Label.Style != nil && cmd.Content.Style != nil {
				if err := e.resetDifferingStyles(printer, cmd.Label.Style, cmd.Content.Style); err != nil {
					return err
				}
			} else if cmd.Label.Style != nil {
				// Si solo hay estilo de label, resetearlo
				if err := e.resetTextStyle(printer, cmd.Label.Style); err != nil {
					return err
				}
			}
		}
	}

	// Imprimir contenido
	if cmd.Content.Text != "" {
		// Aplicar alineación del contenido
		if err := e.applyAlign(printer, cmd.Content.Align); err != nil {
			return err
		}

		// Aplicar estilo del contenido solo si existe
		if cmd.Content.Style != nil {
			if err := e.applyTextStyle(printer, cmd.Content.Style); err != nil {
				return fmt.Errorf("failed to apply content style: %w", err)
			}
		}

		// Determinar si hacer salto de línea
		newLine := true // default
		if cmd.NewLine != nil {
			newLine = *cmd.NewLine
		}

		if newLine {
			if err := printer.PrintLine(cmd.Content.Text); err != nil {
				return err
			}
		} else {
			if err := printer.Print(cmd.Content.Text); err != nil {
				return err
			}
		}

		// Resetear estilos del contenido solo si había estilo aplicado
		if cmd.Content.Style != nil {
			if err := e.resetTextStyle(printer, cmd.Content.Style); err != nil {
				return fmt.Errorf("failed to reset content style: %w", err)
			}
		}
	}

	// Resetear alineación a izquierda al final
	return e.applyAlign(printer, strPtr(constants.Left.String()))
}

func (e *Executor) resetDifferingStyles(printer *service.Printer, labelStyle, contentStyle *TextStyle) error {
	if labelStyle == nil {
		return nil
	}

	// Si contentStyle es nil, resetear labelStyle
	if contentStyle == nil {
		return e.resetTextStyle(printer, labelStyle)
	}

	// Resetear solo los campos que difieren
	if !sameOrNil(labelStyle.Bold, contentStyle.Bold) {
		if labelStyle.Bold != nil && *labelStyle.Bold {
			if err := printer.DisableBold(); err != nil {
				return err
			}
		}
	}

	if !sameOrNil(labelStyle.Size, contentStyle.Size) {
		if labelStyle.Size != nil && *labelStyle.Size != constants.DefaultTextSize.String() {
			if err := printer.SingleSize(); err != nil {
				return err
			}
		}
	}

	if !sameOrNil(labelStyle.Underline, contentStyle.Underline) {
		if labelStyle.Underline != nil && *labelStyle.Underline != constants.NoDot.String() {
			if err := printer.NoDot(); err != nil {
				return err
			}
		}
	}

	if !sameOrNil(labelStyle.Inverse, contentStyle.Inverse) {
		if labelStyle.Inverse != nil && *labelStyle.Inverse {
			if err := printer.InverseOff(); err != nil {
				return err
			}
		}
	}

	if !sameOrNil(labelStyle.Font, contentStyle.Font) {
		if labelStyle.Font != nil && *labelStyle.Font != "A" {
			if err := printer.FontA(); err != nil {
				return err
			}
		}
	}

	return nil
}

// sameOrNil compara dos punteros a valores comparables
func sameOrNil[T comparable](a, b *T) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

func strPtr(s string) *string {
	return &s
}

func (e *Executor) applyAlign(printer *service.Printer, align *string) error {
	alignValue := constants.DefaultTextAlignment.String() // default
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
		log.Printf("Unknown alignment: %s, using left", alignValue)
		return printer.AlignLeft()
	}
}

func (e *Executor) applySize(printer *service.Printer, size string) error {
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
		// Intentar parsear tamaño personalizado WxH
		if len(ss) == 3 && ss[1] == 'x' {
			parts := strings.Split(ss, "x")
			widthMultiplier := parts[0][0] - '0'
			heightMultiplier := parts[1][0] - '0'
			if widthMultiplier >= 1 && widthMultiplier <= 8 &&
				heightMultiplier >= 1 && heightMultiplier <= 8 {
				return printer.CustomSize(widthMultiplier, heightMultiplier)
			}
		}
		log.Printf("Unknown text size: %s, using single size", size)
		return printer.SingleSize()
	}
}

func (e *Executor) applyUnderline(printer *service.Printer, underline string) error {
	switch strings.ToLower(underline) {
	case constants.NoDot.String():
		return printer.NoDot()
	case constants.OneDot.String():
		return printer.OneDot()
	case constants.TwoDot.String():
		return printer.TwoDot()
	default:
		log.Printf("Unknown underline style: %s, using none", underline)
		return printer.NoDot()
	}
}

func (e *Executor) applyFont(printer *service.Printer, font string) error {
	switch strings.ToLower(font) {
	case "", "a":
		return printer.FontA()
	case "b":
		return printer.FontB()
	default:
		log.Printf("Unknown font: %s, using Font A", font)
		return printer.FontA()
	}
}

// TODO: Normalize string comparisons
// .ToLower, .TrimSpace, etc.

// applyTextStyle aplica los estilos de texto especificados
func (e *Executor) applyTextStyle(printer *service.Printer, style *TextStyle) error {
	if style == nil {
		return nil
	}

	if printer == nil {
		return fmt.Errorf("printer is nil")
	}

	// Aplicar bold
	if style.Bold != nil {
		if *style.Bold {
			if err := printer.EnableBold(); err != nil {
				return fmt.Errorf("failed to enable bold: %w", err)
			}
		} else {
			if err := printer.DisableBold(); err != nil {
				return fmt.Errorf("failed to disable bold: %w", err)
			}
		}
	}

	// Aplicar tamaño
	if style.Size != nil {
		if err := e.applySize(printer, *style.Size); err != nil {
			return fmt.Errorf("failed to apply size: %w", err)
		}
	}

	// Aplicar underline
	if style.Underline != nil {
		if err := e.applyUnderline(printer, *style.Underline); err != nil {
			return fmt.Errorf("failed to apply underline: %w", err)
		}
	}

	// Aplicar inverse
	if style.Inverse != nil {
		if *style.Inverse {
			if err := printer.InverseOn(); err != nil {
				return fmt.Errorf("failed to enable inverse: %w", err)
			}
		} else {
			if err := printer.InverseOff(); err != nil {
				return fmt.Errorf("failed to disable inverse: %w", err)
			}
		}
	}

	// Apply font
	if style.Font != nil {
		if err := e.applyFont(printer, *style.Font); err != nil {
			return fmt.Errorf("failed to apply font: %w", err)
		}
	}

	return nil
}

// resetTextStyle resetea los estilos aplicados
func (e *Executor) resetTextStyle(printer *service.Printer, style *TextStyle) error {
	// Validaciones de seguridad
	if printer == nil {
		return fmt.Errorf("printer is nil")
	}

	if style == nil {
		return nil
	}

	// Reset bold solo si estaba activado
	if style.Bold != nil && *style.Bold {
		if err := printer.DisableBold(); err != nil {
			return fmt.Errorf("failed to reset bold: %w", err)
		}
	}

	// Reset size solo si no es el default
	if style.Size != nil && *style.Size != "" && *style.Size != constants.Normal.String() && *style.Size != "1" {
		if err := printer.SingleSize(); err != nil {
			return fmt.Errorf("failed to reset size: %w", err)
		}
	}

	// Reset underline solo si no es el default
	if style.Underline != nil && *style.Underline != "" && *style.Underline != "0" && *style.Underline != constants.NoDot.String() {
		if err := printer.NoDot(); err != nil {
			return fmt.Errorf("failed to reset underline: %w", err)
		}
	}

	// Reset inverse solo si estaba activado
	if style.Inverse != nil && *style.Inverse {
		if err := printer.InverseOff(); err != nil {
			return fmt.Errorf("failed to reset inverse: %w", err)
		}
	}

	// Reset font solo si no es el default
	if style.Font != nil && *style.Font != "" && *style.Font != "A" && *style.Font != "a" {
		if err := printer.FontA(); err != nil {
			return fmt.Errorf("failed to reset font: %w", err)
		}
	}

	return nil
}

// TODO: Check if still used, otherwise think in another DRY way to set alignments (aling into left or defer)

/* applyAlignment aplica alineación a la impresora y retorna una función para restaurar
func (e *Executor) applyAlignInto(printer *service.Printer, align string) (restore func() error, err error) {
	switch align {
	case constants.Center.String():
		err = printer.Center()
	case constants.Right.String():
		err = printer.Right()
	default:
		err = printer.Left()
	}

	if err != nil {
		return nil, err
	}

	// Retornar función para restaurar a izquierda
	return func() error {
		return printer.Left()
	}, nil
}
*/
