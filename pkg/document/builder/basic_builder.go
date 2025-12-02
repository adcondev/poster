package builder

import (
	"fmt"
)

// Feed adds paper feed lines
func (b *DocumentBuilder) Feed(lines int) *DocumentBuilder {
	return b.addCommand("feed", feedCommand{Lines: lines})
}

// Cut adds a paper cut command
func (b *DocumentBuilder) Cut() *DocumentBuilder {
	return b.addCommand("cut", cutCommand{Mode: "partial", Feed: 2})
}

// FullCut adds a full paper cut
func (b *DocumentBuilder) FullCut() *DocumentBuilder {
	return b.addCommand("cut", cutCommand{Mode: "full", Feed: 2})
}

// CutWithFeed adds a cut with custom feed
func (b *DocumentBuilder) CutWithFeed(mode string, feed int) *DocumentBuilder {
	return b.addCommand("cut", cutCommand{Mode: mode, Feed: feed})
}

// Separator adds a separator line
func (b *DocumentBuilder) Separator(char string) *DocumentBuilder {
	return b.addCommand("separator", separatorCommand{Char: char, Length: 48})
}

// SeparatorWithLength adds a separator with custom length
func (b *DocumentBuilder) SeparatorWithLength(char string, length int) *DocumentBuilder {
	return b.addCommand("separator", separatorCommand{Char: char, Length: length})
}

// Pulse opens the cash drawer
func (b *DocumentBuilder) Pulse() *DocumentBuilder {
	return b.Raw("1B 70 00 32 64").Comment("Cash Drawer Pulse").End()
}

// Beep emits beep sounds
//
// ADVERTENCIA DE COMPATIBILIDAD:
//
// El estándar ESC/POS define 'ESC ( A' (Buzz) para el buzzer.
// Sin embargo, este hardware específico usa 'ESC B' (no existe es ESC/POS).
//
//	Formato: ESC B n t
//	Hex:     1B  42 n t
//
//	1B 42 -> Cabecera del comando (ESC B)
//	n     -> times: Número de repeticiones (9 veces)
//	t     -> lapse: Duración/Intervalo (factor de tiempo, aprox 100ms * t)
func (b *DocumentBuilder) Beep(times, lapse int) *DocumentBuilder {
	raw := fmt.Sprintf("1B 42 %02X %02X", times, lapse)
	return b.Raw(raw).Comment("Beep Sound").End()
}
