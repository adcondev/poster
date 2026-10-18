package emulator

import (
	"bytes"
	"fmt"

	"github.com/adcondev/poster/pkg/constants"
	"github.com/adcondev/poster/pkg/graphics"
)

// ParseESC parses a raw ESC/POS byte stream and renders the resulting formatting
// and content onto the emulator Engine canvas.
func (e *Engine) ParseESC(raw []byte) error {
	var textBuf bytes.Buffer
	flushText := func() {
		if textBuf.Len() > 0 {
			e.Print(textBuf.String())
			textBuf.Reset()
		}
	}

	i := 0
	for i < len(raw) {
		switch raw[i] {
		case 0x0A: // LF
			flushText()
			e.NewLine()
			i++
		case 0x0D: // CR
			flushText()
			if i+1 < len(raw) && raw[i+1] == 0x0A {
				i++
			} else {
				e.NewLine()
				i++
			}
		case 0x1B: // ESC
			flushText()
			i = e.parseESCCommand(raw, i)
		case 0x1D: // GS
			flushText()
			i = e.parseGSCommand(raw, i)
		default:
			if raw[i] >= 0x20 {
				textBuf.WriteByte(raw[i])
			}
			i++
		}

		if i > len(raw) {
			i = len(raw)
		}
	}

	flushText()
	return nil
}

// parseESCCommand processes an ESC (0x1B) command sequence and returns the new byte index.
func (e *Engine) parseESCCommand(raw []byte, i int) int {
	if i+1 >= len(raw) {
		return len(raw)
	}

	cmd := raw[i+1]
	switch cmd {
	case '@': // Initialize printer
		e.Reset()
		return i + 2
	case '2': // Default line spacing
		e.state.LineSpacing = float64(constants.DefaultLineSpacing)
		return i + 2
	case 'E', '-', 'M', 'd', '3', 't':
		if i+2 >= len(raw) {
			return len(raw)
		}
		e.handleESCSingleParam(cmd, raw[i+2])
		return i + 3
	case 'a':
		return e.handleESCAlign(raw, i)
	case '!':
		return e.handleESCPrintMode(raw, i)
	case 'p':
		return skipESCPulse(raw, i)
	case 'B':
		return skipESCBeep(raw, i)
	default:
		return i + 2
	}
}

func (e *Engine) handleESCSingleParam(cmd, arg byte) {
	switch cmd {
	case 'E': // Bold
		e.SetBold((arg & 0x01) != 0)
	case '-': // Underline
		mode := int(arg)
		if mode >= '0' && mode <= '2' {
			mode -= '0'
		}
		e.SetUnderline(mode)
	case 'M': // Font A / B
		if arg == 1 || arg == '1' {
			e.SetFont("B")
		} else {
			e.SetFont("A")
		}
	case 'd': // Feed lines
		e.Feed(int(arg))
	case '3': // Set line spacing
		e.state.LineSpacing = float64(arg)
	case 't': // Code table selection (ignored by emulator)
	}
}

func (e *Engine) handleESCAlign(raw []byte, i int) int {
	if i+2 >= len(raw) {
		return len(raw)
	}
	switch raw[i+2] {
	case 1, '1':
		e.AlignCenter()
	case 2, '2':
		e.AlignRight()
	default:
		e.AlignLeft()
	}
	return i + 3
}

func (e *Engine) handleESCPrintMode(raw []byte, i int) int {
	if i+2 >= len(raw) {
		return len(raw)
	}
	n := raw[i+2]
	if (n & 0x01) != 0 {
		e.SetFont("B")
	} else {
		e.SetFont("A")
	}
	e.SetBold((n & 0x08) != 0)
	w, h := 1, 1
	if (n & 0x10) != 0 {
		h = 2
	}
	if (n & 0x20) != 0 {
		w = 2
	}
	e.SetSize(w, h)
	if (n & 0x80) != 0 {
		e.SetUnderline(1)
	} else {
		e.SetUnderline(0)
	}
	return i + 3
}

func skipESCPulse(raw []byte, i int) int {
	if i+4 < len(raw) {
		return i + 5
	}
	return len(raw)
}

func skipESCBeep(raw []byte, i int) int {
	if i+3 < len(raw) {
		return i + 4
	}
	return len(raw)
}

// parseGSCommand processes a GS (0x1D) command sequence and returns the new byte index.
func (e *Engine) parseGSCommand(raw []byte, i int) int {
	if i+1 >= len(raw) {
		return len(raw)
	}

	switch raw[i+1] {
	case '!': // Select character size
		return e.handleGSSize(raw, i)
	case 'B': // White/black reverse mode
		return e.handleGSInverse(raw, i)
	case 'V': // Cut paper
		return e.handleGSCut(raw, i)
	case 'v': // Print raster bit image
		return e.handleGSRaster(raw, i)
	case 'h', 'w', 'H', 'f': // Barcode height, width, HRI position, HRI font
		if i+2 < len(raw) {
			return i + 3
		}
		return len(raw)
	case 'k': // Print barcode
		return skipGSBarcode(raw, i)
	case '(': // 2D barcode / special function
		return skipGS2DFunction(raw, i)
	default:
		return i + 2
	}
}

func (e *Engine) handleGSSize(raw []byte, i int) int {
	if i+2 >= len(raw) {
		return len(raw)
	}
	n := raw[i+2]
	w := int((n>>4)&0x07) + 1
	h := int(n&0x07) + 1
	e.SetSize(w, h)
	return i + 3
}

func (e *Engine) handleGSInverse(raw []byte, i int) int {
	if i+2 >= len(raw) {
		return len(raw)
	}
	e.SetInverse((raw[i+2] & 0x01) != 0)
	return i + 3
}

func (e *Engine) handleGSCut(raw []byte, i int) int {
	if i+2 >= len(raw) {
		return len(raw)
	}
	m := raw[i+2]
	step := 3
	if m == 65 || m == 66 || m == 97 || m == 98 || m == 103 || m == 104 {
		if i+3 < len(raw) {
			e.Feed(int(raw[i+3]))
			step = 4
		}
	}
	partial := m != 0 && m != '0' && m != 65
	e.Cut(partial)
	return i + step
}

func (e *Engine) handleGSRaster(raw []byte, i int) int {
	if i+2 >= len(raw) || raw[i+2] != '0' {
		return i + 2
	}
	if i+7 >= len(raw) {
		return len(raw)
	}

	xBytes := int(raw[i+4]) | (int(raw[i+5]) << 8)
	yDots := int(raw[i+6]) | (int(raw[i+7]) << 8)
	dataLen := xBytes * yDots
	start := i + 8
	end := start + dataLen
	if end > len(raw) {
		end = len(raw)
	}
	if xBytes > 0 && yDots > 0 {
		e.renderRasterBitmap(raw[start:end], xBytes, yDots)
	}
	return end
}

func (e *Engine) renderRasterBitmap(raster []byte, xBytes, yDots int) {
	bmp := graphics.NewMonochromeBitmap(xBytes*8, yDots)
	for y := 0; y < yDots; y++ {
		for bx := 0; bx < xBytes; bx++ {
			idx := y*xBytes + bx
			if idx >= len(raster) {
				break
			}
			b := raster[idx]
			for bit := 0; bit < 8; bit++ {
				if (b & (1 << uint(7-bit))) != 0 {
					bmp.SetPixel(bx*8+bit, y, true)
				}
			}
		}
	}
	_ = e.PrintImageAligned(bmp.ToImage(), xBytes*8, e.state.Align)
}

func skipGSBarcode(raw []byte, i int) int {
	if i+2 >= len(raw) {
		return len(raw)
	}
	m := raw[i+2]
	if m <= 6 {
		j := i + 3
		for j < len(raw) && raw[j] != 0x00 {
			j++
		}
		if j < len(raw) {
			return j + 1
		}
		return len(raw)
	}
	if i+3 < len(raw) {
		n := int(raw[i+3])
		end := i + 4 + n
		if end > len(raw) {
			return len(raw)
		}
		return end
	}
	return len(raw)
}

func skipGS2DFunction(raw []byte, i int) int {
	if i+4 < len(raw) {
		pL := int(raw[i+3])
		pH := int(raw[i+4])
		end := i + 5 + (pL | (pH << 8))
		if end > len(raw) {
			return len(raw)
		}
		return end
	}
	return len(raw)
}

// RenderBytesToPNG parses raw ESC/POS commands with default 80mm configuration and returns PNG bytes.
func RenderBytesToPNG(raw []byte) ([]byte, error) {
	eng, err := NewDefaultEngine()
	if err != nil {
		return nil, fmt.Errorf("creating default engine: %w", err)
	}
	if err := eng.ParseESC(raw); err != nil {
		return nil, fmt.Errorf("parsing esc/pos: %w", err)
	}
	var buf bytes.Buffer
	if err := eng.WritePNG(&buf); err != nil {
		return nil, fmt.Errorf("writing png: %w", err)
	}
	return buf.Bytes(), nil
}

// RenderBytesToPNGWithConfig parses raw ESC/POS commands with a custom configuration and returns PNG bytes.
func RenderBytesToPNGWithConfig(raw []byte, cfg Config) ([]byte, error) {
	eng, err := NewEngine(cfg)
	if err != nil {
		return nil, fmt.Errorf("creating engine with config: %w", err)
	}
	if err := eng.ParseESC(raw); err != nil {
		return nil, fmt.Errorf("parsing esc/pos: %w", err)
	}
	var buf bytes.Buffer
	if err := eng.WritePNG(&buf); err != nil {
		return nil, fmt.Errorf("writing png: %w", err)
	}
	return buf.Bytes(), nil
}
