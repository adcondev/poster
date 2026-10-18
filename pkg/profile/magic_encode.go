package profile

import (
	"strings"

	"github.com/adcondev/poster/pkg/commands/character"
)

// NormalizeTypography converts typographical/smart Unicode characters (smart quotes,
// dashes, ellipsis, non-breaking spaces) to their standard ASCII printer-friendly equivalents.
func NormalizeTypography(text string) string {
	if text == "" {
		return ""
	}

	var sb strings.Builder
	sb.Grow(len(text))

	for _, r := range text {
		switch r {
		case '\u2018', '\u2019', '\u201A', '\u201B': // ‘, ’, ‚, ‛
			sb.WriteRune('\'')
		case '\u201C', '\u201D', '\u201E', '\u201F': // “, ”, „, ‟
			sb.WriteRune('"')
		case '\u2013', '\u2014', '\u2015', '\u2212': // –, —, ―, −
			sb.WriteRune('-')
		case '\u2026': // …
			sb.WriteString("...")
		case '\u00A0', '\u202F': // NBSP, narrow NBSP
			sb.WriteRune(' ')
		default:
			sb.WriteRune(r)
		}
	}

	return sb.String()
}

// EncodeBytesSafe encodes text to printer-compatible bytes.
//
// If StrictEncoding is true, it encodes text strictly using the active code table
// and returns an error if any character cannot be represented.
//
// Otherwise, it performs safe rune-by-rune encoding:
//   - If a rune cannot be encoded in the current code page and AutoSwitchCodePage
//     is enabled, it attempts to switch to an appropriate fallback code page
//     (emitting the ESC t n command: 0x1B, 0x74, byte(cand)).
//   - If no candidate code page supports the rune, it substitutes '?' (0x3F).
func (p *Escpos) EncodeBytesSafe(text string) ([]byte, error) {
	text = NormalizeTypography(text)

	if p.StrictEncoding {
		s, err := p.EncodeString(text)
		if err != nil {
			return nil, err
		}
		return []byte(s), nil
	}

	activeTable := p.CodeTable
	if !p.IsSupported(activeTable) {
		activeTable = character.WPC1252
	}

	var out []byte

	candidates := p.FallbackCodePages
	if len(candidates) == 0 {
		candidates = []character.CodeTable{
			character.PC858,
			character.WPC1252,
			character.PC437,
			character.PC850,
			character.PC866,
			character.WPC1251,
			character.ISO88597,
		}
	}

	for _, r := range text {
		rStr := string(r)
		enc, ok := codeTableMap[activeTable]
		if ok {
			if encoded, err := enc.NewEncoder().String(rStr); err == nil {
				out = append(out, encoded...)
				continue
			}
		}

		if p.AutoSwitchCodePage {
			switched := false
			for _, cand := range candidates {
				if cand == activeTable {
					continue
				}
				candEnc, cOk := codeTableMap[cand]
				if !cOk {
					continue
				}
				if candEncoded, cErr := candEnc.NewEncoder().String(rStr); cErr == nil {
					out = append(out, 0x1B, 0x74, byte(cand))
					activeTable = cand
					p.CodeTable = cand
					out = append(out, candEncoded...)
					switched = true
					break
				}
			}
			if switched {
				continue
			}
		}

		out = append(out, '?')
	}

	return out, nil
}
