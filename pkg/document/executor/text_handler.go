package executor

import (
	"encoding/json"
	"fmt"

	"github.com/adcondev/poster/pkg/constants"
	"github.com/adcondev/poster/pkg/service"
)

// TextCommand for text handler
type TextCommand struct {
	Content Content `json:"content"`
	Label   *Label  `json:"label,omitempty"`
	NewLine *bool   `json:"new_line,omitempty"`
}

// Content for text
type Content struct {
	Text  string     `json:"text"`
	Style *TextStyle `json:"content_style,omitempty"`
	Align *string    `json:"align,omitempty"`
}

// Label for text
type Label struct {
	Text      string     `json:"text,omitempty"`
	Style     *TextStyle `json:"label_style,omitempty"`
	Separator *string    `json:"separator,omitempty"`
	Align     *string    `json:"align,omitempty"`
}

// TextStyle for text formatting
type TextStyle struct {
	Bold      *bool   `json:"bold,omitempty"`
	Size      *string `json:"size,omitempty"`
	Underline *string `json:"underline,omitempty"`
	Inverse   *bool   `json:"inverse,omitempty"`
	Font      *string `json:"font,omitempty"`
}

func strPtr(s string) *string {
	return &s
}

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

// TODO: Normalize string comparisons
// .ToLower, .TrimSpace, etc.

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
