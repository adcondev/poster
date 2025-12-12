package executor

import (
	"fmt"

	"github.com/adcondev/poster/pkg/constants"
	"github.com/adcondev/poster/pkg/service"
)

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
