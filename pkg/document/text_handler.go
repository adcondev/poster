package document

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/adcondev/pos-printer/pkg/service"
)

// TextCommand represents a text command
type TextCommand struct {
	// Contenido principal
	Content Content `json:"content"`
	// Label opcional
	Label *Label `json:"label,omitempty"`
	// Control de línea
	NewLine *bool `json:"new_line,omitempty"` // default: true
}

// Label representa una etiqueta de texto
type Label struct {
	Text      string     `json:"text,omitempty"`        // texto del label default: ""
	Style     *TextStyle `json:"label_style,omitempty"` // estilo del label
	Separator *string    `json:"separator,omitempty"`   // default: ": "
	Align     *string    `json:"align,omitempty"`       // left, center, right default: left
}

// Content representa el contenido de texto
type Content struct {
	Text  string     `json:"text"`                    // texto del contenido
	Style *TextStyle `json:"content_style,omitempty"` // estilo del contenido
	Align *string    `json:"align,omitempty"`         // left, center, right default: left
}

// TextStyle estilo de texto
type TextStyle struct {
	Bold      *bool   `json:"bold,omitempty"`      // true, false default: false
	Size      *string `json:"size,omitempty"`      // 1x1, 2x2, 3x3 ... 8x8 default: 1x1
	Underline *string `json:"underline,omitempty"` // "1pt", "2pt", default: "0pt"
	Inverse   *bool   `json:"inverse,omitempty"`   // true, false default: false
	Font      *string `json:"font,omitempty"`      // A, B default: A
}

// handleText manages text commands
func (e *Executor) handleText(printer *service.Printer, data json.RawMessage) error {
	var cmd TextCommand
	if err := json.Unmarshal(data, &cmd); err != nil {
		return fmt.Errorf("failed to parse text command: %w", err)
	}

	// Si hay label con diferente alineación que el content
	if cmd.Label != nil {
		// Aplicar estilo del label
		if err := e.applyTextStyle(printer, cmd.Label.Style); err != nil {
			return err
		}

		// Construir texto del label
		labelText := cmd.Label.Text
		if cmd.Label.Separator != nil && *cmd.Label.Separator != "" {
			labelText += *cmd.Label.Separator
		} else {
			labelText += ": "
		}

		if err := e.applyAlign(printer, cmd.Label.Align); err != nil {
			return err
		}

		// Si las alineaciones son diferentes, imprimir label y resetear
		if !sameOrNil(cmd.Label.Align, cmd.Content.Align) {
			// Imprimir label con salto de línea
			if err := printer.PrintLine(labelText); err != nil {
				return err
			}
			// Reset completo antes del content
			if err := e.resetTextStyle(printer, cmd.Label.Style); err != nil {
				return err
			}
		} else {
			if err := printer.Print(labelText); err != nil {
				return err
			}
			// Reset solo los estilos diferentes
			if err := e.resetDifferingStyles(printer, cmd.Label.Style, cmd.Content.Style); err != nil {
				return err
			}
		}
	}

	// Imprimir contenido
	if cmd.Content.Text != "" {
		if err := e.applyAlign(printer, cmd.Content.Align); err != nil {
			return err
		}
		// Aplicar estilo del contenido
		if err := e.applyTextStyle(printer, cmd.Content.Style); err != nil {
			return fmt.Errorf("failed to apply content style: %w", err)
		}

		if cmd.NewLine != nil {
			if *cmd.NewLine {
				if err := printer.PrintLine(cmd.Content.Text); err != nil { // Con salto
					return err
				}
			} else {
				if err := printer.Print(cmd.Content.Text); err != nil { // Sin salto
					return err
				}
			}
		} else {
			// Default: true
			if err := printer.PrintLine(cmd.Content.Text); err != nil {
				return err
			}
		}
	}

	// Resetear estilos del contenido
	if err := e.resetTextStyle(printer, cmd.Content.Style); err != nil {
		return fmt.Errorf("failed to reset content style: %w", err)
	}

	// Resetear alineación a izquierda
	if err := e.applyAlign(printer, strPtr(left)); err != nil {
		return err
	}

	return nil
}

func (e *Executor) resetDifferingStyles(printer *service.Printer, labelStyle, contentStyle *TextStyle) error {
	if labelStyle == nil || contentStyle == nil {
		log.Printf("No text styles to compare for resetting")
		return nil
	}

	if !sameOrNil(labelStyle.Bold, contentStyle.Bold) {
		if err := printer.DisableBold(); err != nil {
			return err
		}
	}

	if !sameOrNil(labelStyle.Size, contentStyle.Size) {
		if err := printer.SingleSize(); err != nil {
			return err
		}
	}

	if !sameOrNil(labelStyle.Underline, contentStyle.Underline) {
		if err := printer.NoDot(); err != nil {
			return err
		}
	}

	if !sameOrNil(labelStyle.Inverse, contentStyle.Inverse) {
		if err := printer.InverseOff(); err != nil {
			return err
		}
	}

	if !sameOrNil(labelStyle.Font, contentStyle.Font) {
		if err := printer.FontA(); err != nil {
			return err
		}
	}

	return nil
}

func sameOrNil[T comparable](a, b *T) bool {
	if a == nil || b == nil {
		return true
	}
	return *a == *b
}

func strPtr(s string) *string {
	return &s
}

func (e *Executor) applyAlign(printer *service.Printer, align *string) error {
	if align == nil {
		align = strPtr(left)
	}
	// Aplicar alineación
	switch strings.ToLower(*align) {
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
		log.Printf("Unknown alignment: %s, using left", *align)
		if err := printer.AlignLeft(); err != nil {
			return err
		}
	}
	return nil
}

func (e *Executor) applySize(printer *service.Printer, size string) error {
	if size != "" {
		switch ss := strings.ToLower(size); ss {
		case "1x1", "1":
			if err := printer.SingleSize(); err != nil {
				return err
			}
		case "2x2", "2":
			if err := printer.DoubleSize(); err != nil {
				return err
			}
		case "3x3", "3":
			if err := printer.TripleSize(); err != nil {
				return err
			}
		case "4x4", "4":
			if err := printer.QuadraSize(); err != nil {
				return err
			}
		case "5x5", "5":
			if err := printer.PentaSize(); err != nil {
				return err
			}
		case "6x6", "6":
			if err := printer.HexaSize(); err != nil {
				return err
			}
		case "7x7", "7":
			if err := printer.HeptaSize(); err != nil {
				return err
			}
		case "8x8", "8":
			if err := printer.OctaSize(); err != nil {
				return err
			}
		default:
			// Intentar parsear tamaño personalizado WxH
			if len(ss) == 3 && ss[1] == 'x' {
				parts := strings.Split(ss, "x")
				widthMultiplier := parts[0][0] - '0'
				heightMultiplier := parts[1][0] - '0'
				if err := printer.CustomSize(widthMultiplier, heightMultiplier); err != nil {
					return err
				}
				log.Printf("Applied custom text size: %s", size)
			} else {
				if err := printer.SingleSize(); err != nil {
					return err
				}
				log.Printf("Unknown text size: %s, using single size", size)
			}
		}
	} else {
		// Default size
		if err := printer.SingleSize(); err != nil {
			return err
		}
	}
	return nil
}

func (e *Executor) applyUnderline(printer *service.Printer, underline string) error {
	switch strings.ToLower(underline) {
	case "0", "0pt":
		// No underline
		err := printer.NoDot()
		if err != nil {
			return err
		}
	case "1", "1pt":
		err := printer.OneDot()
		if err != nil {
			return err
		}
	case "2", "2pt":
		err := printer.TwoDot()
		if err != nil {
			return err
		}
	default:
		err := printer.NoDot()
		if err != nil {
			return err
		}
		log.Printf("Unknown underline style: %s, using none", underline)
	}
	return nil
}

func (e *Executor) applyFont(printer *service.Printer, font string) error {
	switch strings.ToLower(font) {
	case "a":
		if err := printer.FontA(); err != nil {
			return err
		}
	case "b":
		if err := printer.FontB(); err != nil {
			return err
		}
	default:
		log.Printf("Unknown font: %s, using Font A", font)
		if err := printer.FontA(); err != nil {
			return err
		}
	}
	return nil
}

// applyTextStyle aplica los estilos de texto especificados
func (e *Executor) applyTextStyle(printer *service.Printer, style *TextStyle) error {
	if style == nil {
		log.Printf("No text style to apply")
		return nil
	}

	// Aplicar bold
	if style.Bold != nil && *style.Bold {
		if err := printer.EnableBold(); err != nil {
			return fmt.Errorf("failed to enable bold: %w", err)
		}
	} else {
		if err := printer.DisableBold(); err != nil {
			return fmt.Errorf("failed to disable bold: %w", err)
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
	if style.Inverse != nil && *style.Inverse {
		err := printer.InverseOn()
		if err != nil {
			return err
		}
	} else {
		err := printer.InverseOff()
		if err != nil {
			return err
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
	// Reset bold
	if style.Bold != nil && *style.Bold {
		if err := printer.DisableBold(); err != nil {
			return err
		}
	}

	// Reset size
	if style.Size != nil {
		if err := printer.SingleSize(); err != nil {
			return err
		}
	}

	// Reset underline
	if style.Underline != nil {
		if err := printer.NoDot(); err != nil {
			return err
		}
	}

	// Reset inverse
	if style.Inverse != nil && *style.Inverse {
		if err := printer.InverseOff(); err != nil {
			return err
		}
	}

	// Reset font
	if style.Font != nil {
		if err := printer.FontA(); err != nil {
			return err
		}
	}

	return nil
}
