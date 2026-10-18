package profile

import (
	"strings"

	"github.com/adcondev/poster/pkg/commands/character"
)

// CreateEpsonTMT20III creates a profile for the Epson TM-T20III thermal receipt printer.
func CreateEpsonTMT20III() *Escpos {
	p := CreateProfile80mm()
	p.Model = "Epson TM-T20III"
	p.CodeTable = character.PC858
	p.AutoSwitchCodePage = true
	return p
}

// CreateEpsonTMT88VI creates a profile for the high-speed Epson TM-T88VI thermal receipt printer.
func CreateEpsonTMT88VI() *Escpos {
	p := CreateProfile80mm()
	p.Model = "Epson TM-T88VI"
	p.CodeTable = character.PC858
	p.AutoSwitchCodePage = true
	return p
}

// CreateStarTSP100 creates a profile for the Star Micronics TSP100 thermal receipt printer.
func CreateStarTSP100() *Escpos {
	p := CreateProfile80mm()
	p.Model = "Star TSP100"
	p.CodeTable = character.WPC1252
	p.AutoSwitchCodePage = true
	return p
}

// CreateBixolonSRP350III creates a profile for the Bixolon SRP-350III thermal receipt printer.
func CreateBixolonSRP350III() *Escpos {
	p := CreateProfile80mm()
	p.Model = "Bixolon SRP-350III"
	p.CodeTable = character.PC858
	p.AutoSwitchCodePage = true
	return p
}

// CreateXprinterXP58 creates a profile for the Xprinter XP-58 58mm thermal receipt printer.
func CreateXprinterXP58() *Escpos {
	p := CreateProfile58mm()
	p.Model = "Xprinter XP-58"
	p.CodeTable = character.PC850
	p.HasQR = false
	p.SupportsCutter = false
	p.SupportsDrawer = true
	return p
}

// CreateXprinterXP80 creates a profile for the Xprinter XP-80 80mm thermal receipt printer.
func CreateXprinterXP80() *Escpos {
	p := CreateProfile80mm()
	p.Model = "Xprinter XP-80"
	p.CodeTable = character.PC858
	p.HasQR = true
	p.SupportsCutter = true
	p.SupportsDrawer = true
	return p
}

// FromModel resolves a printer model string into an Escpos profile.
// It is an alias for ProfileFromModel.
func FromModel(model string) (*Escpos, bool) {
	return ProfileFromModel(model)
}

// ProfileFromModel resolves a printer model string into an Escpos profile.
// Matching is case-insensitive and supports model name substrings.
//
//nolint:revive // ProfileFromModel is the canonical name specified for catalog resolution
func ProfileFromModel(model string) (*Escpos, bool) {
	m := strings.ToLower(strings.TrimSpace(model))
	if m == "" {
		return nil, false
	}

	switch {
	case strings.Contains(m, "pt-210") || strings.Contains(m, "pt210"):
		return CreatePt210(), true
	case strings.Contains(m, "gp-58n") || strings.Contains(m, "gp58n"):
		return CreateGP58N(), true
	case strings.Contains(m, "ec-pm-80250") || strings.Contains(m, "ecpm80250"):
		return CreateECPM80250(), true
	case strings.Contains(m, "tm-t20") || strings.Contains(m, "tmt20"):
		return CreateEpsonTMT20III(), true
	case strings.Contains(m, "tm-t88") || strings.Contains(m, "tmt88"):
		return CreateEpsonTMT88VI(), true
	case strings.Contains(m, "tsp100") || strings.Contains(m, "tsp-100"):
		return CreateStarTSP100(), true
	case strings.Contains(m, "srp-350") || strings.Contains(m, "srp350"):
		return CreateBixolonSRP350III(), true
	case strings.Contains(m, "xp-58") || strings.Contains(m, "xp58"):
		return CreateXprinterXP58(), true
	case strings.Contains(m, "xp-80") || strings.Contains(m, "xp80"):
		return CreateXprinterXP80(), true
	case strings.Contains(m, "58mm") || strings.Contains(m, "58"):
		return CreateProfile58mm(), true
	case strings.Contains(m, "80mm") || strings.Contains(m, "80"):
		return CreateProfile80mm(), true
	default:
		return nil, false
	}
}
