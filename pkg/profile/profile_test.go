package profile_test

import (
	"testing"

	"github.com/adcondev/poster/pkg/commands/character"
	"github.com/adcondev/poster/pkg/profile"
)

func TestCreateEpsonTMT20III(t *testing.T) {
	p := profile.CreateEpsonTMT20III()

	if p.Model != "Epson TM-T20III" {
		t.Errorf("expected model 'Epson TM-T20III', got '%s'", p.Model)
	}
	if p.PaperWidth != 80 {
		t.Errorf("expected PaperWidth 80, got %f", p.PaperWidth)
	}
	if p.DPI != 203 {
		t.Errorf("expected DPI 203, got %d", p.DPI)
	}
	if p.DotsPerLine != 576 {
		t.Errorf("expected DotsPerLine 576, got %d", p.DotsPerLine)
	}
	if p.PrintWidth != 72 {
		t.Errorf("expected PrintWidth 72, got %d", p.PrintWidth)
	}
	if p.CodeTable != character.PC858 {
		t.Errorf("expected CodeTable PC858, got %v", p.CodeTable)
	}
	if !p.HasQR {
		t.Error("expected HasQR to be true")
	}
	if !p.SupportsCutter {
		t.Error("expected SupportsCutter to be true")
	}
	if !p.SupportsDrawer {
		t.Error("expected SupportsDrawer to be true")
	}
	if !p.AutoSwitchCodePage {
		t.Error("expected AutoSwitchCodePage to be true")
	}
}

func TestCreateEpsonTMT88VI(t *testing.T) {
	p := profile.CreateEpsonTMT88VI()

	if p.Model != "Epson TM-T88VI" {
		t.Errorf("expected model 'Epson TM-T88VI', got '%s'", p.Model)
	}
	if p.PaperWidth != 80 {
		t.Errorf("expected PaperWidth 80, got %f", p.PaperWidth)
	}
	if p.DPI != 203 {
		t.Errorf("expected DPI 203, got %d", p.DPI)
	}
	if p.DotsPerLine != 576 {
		t.Errorf("expected DotsPerLine 576, got %d", p.DotsPerLine)
	}
	if p.PrintWidth != 72 {
		t.Errorf("expected PrintWidth 72, got %d", p.PrintWidth)
	}
	if p.CodeTable != character.PC858 {
		t.Errorf("expected CodeTable PC858, got %v", p.CodeTable)
	}
	if !p.HasQR {
		t.Error("expected HasQR to be true")
	}
	if !p.SupportsCutter {
		t.Error("expected SupportsCutter to be true")
	}
	if !p.SupportsDrawer {
		t.Error("expected SupportsDrawer to be true")
	}
	if !p.AutoSwitchCodePage {
		t.Error("expected AutoSwitchCodePage to be true")
	}
}

func TestCreateStarTSP100(t *testing.T) {
	p := profile.CreateStarTSP100()

	if p.Model != "Star TSP100" {
		t.Errorf("expected model 'Star TSP100', got '%s'", p.Model)
	}
	if p.PaperWidth != 80 {
		t.Errorf("expected PaperWidth 80, got %f", p.PaperWidth)
	}
	if p.DPI != 203 {
		t.Errorf("expected DPI 203, got %d", p.DPI)
	}
	if p.DotsPerLine != 576 {
		t.Errorf("expected DotsPerLine 576, got %d", p.DotsPerLine)
	}
	if p.PrintWidth != 72 {
		t.Errorf("expected PrintWidth 72, got %d", p.PrintWidth)
	}
	if p.CodeTable != character.WPC1252 {
		t.Errorf("expected CodeTable WPC1252, got %v", p.CodeTable)
	}
	if !p.HasQR {
		t.Error("expected HasQR to be true")
	}
	if !p.SupportsCutter {
		t.Error("expected SupportsCutter to be true")
	}
	if !p.SupportsDrawer {
		t.Error("expected SupportsDrawer to be true")
	}
	if !p.AutoSwitchCodePage {
		t.Error("expected AutoSwitchCodePage to be true")
	}
}

func TestCreateBixolonSRP350III(t *testing.T) {
	p := profile.CreateBixolonSRP350III()

	if p.Model != "Bixolon SRP-350III" {
		t.Errorf("expected model 'Bixolon SRP-350III', got '%s'", p.Model)
	}
	if p.PaperWidth != 80 {
		t.Errorf("expected PaperWidth 80, got %f", p.PaperWidth)
	}
	if p.DPI != 203 {
		t.Errorf("expected DPI 203, got %d", p.DPI)
	}
	if p.DotsPerLine != 576 {
		t.Errorf("expected DotsPerLine 576, got %d", p.DotsPerLine)
	}
	if p.PrintWidth != 72 {
		t.Errorf("expected PrintWidth 72, got %d", p.PrintWidth)
	}
	if p.CodeTable != character.PC858 {
		t.Errorf("expected CodeTable PC858, got %v", p.CodeTable)
	}
	if !p.HasQR {
		t.Error("expected HasQR to be true")
	}
	if !p.SupportsCutter {
		t.Error("expected SupportsCutter to be true")
	}
	if !p.SupportsDrawer {
		t.Error("expected SupportsDrawer to be true")
	}
	if !p.AutoSwitchCodePage {
		t.Error("expected AutoSwitchCodePage to be true")
	}
}

func TestCreateXprinterXP58(t *testing.T) {
	p := profile.CreateXprinterXP58()

	if p.Model != "Xprinter XP-58" {
		t.Errorf("expected model 'Xprinter XP-58', got '%s'", p.Model)
	}
	if p.PaperWidth != 58 {
		t.Errorf("expected PaperWidth 58, got %f", p.PaperWidth)
	}
	if p.DPI != 203 {
		t.Errorf("expected DPI 203, got %d", p.DPI)
	}
	if p.DotsPerLine != 384 {
		t.Errorf("expected DotsPerLine 384, got %d", p.DotsPerLine)
	}
	if p.PrintWidth != 48 {
		t.Errorf("expected PrintWidth 48, got %d", p.PrintWidth)
	}
	if p.CodeTable != character.PC850 {
		t.Errorf("expected CodeTable PC850, got %v", p.CodeTable)
	}
	if p.HasQR {
		t.Error("expected HasQR to be false")
	}
	if p.SupportsCutter {
		t.Error("expected SupportsCutter to be false")
	}
	if !p.SupportsDrawer {
		t.Error("expected SupportsDrawer to be true")
	}
}

func TestCreateXprinterXP80(t *testing.T) {
	p := profile.CreateXprinterXP80()

	if p.Model != "Xprinter XP-80" {
		t.Errorf("expected model 'Xprinter XP-80', got '%s'", p.Model)
	}
	if p.PaperWidth != 80 {
		t.Errorf("expected PaperWidth 80, got %f", p.PaperWidth)
	}
	if p.DPI != 203 {
		t.Errorf("expected DPI 203, got %d", p.DPI)
	}
	if p.DotsPerLine != 576 {
		t.Errorf("expected DotsPerLine 576, got %d", p.DotsPerLine)
	}
	if p.PrintWidth != 72 {
		t.Errorf("expected PrintWidth 72, got %d", p.PrintWidth)
	}
	if p.CodeTable != character.PC858 {
		t.Errorf("expected CodeTable PC858, got %v", p.CodeTable)
	}
	if !p.HasQR {
		t.Error("expected HasQR to be true")
	}
	if !p.SupportsCutter {
		t.Error("expected SupportsCutter to be true")
	}
	if !p.SupportsDrawer {
		t.Error("expected SupportsDrawer to be true")
	}
}

func TestProfileFromModel(t *testing.T) {
	tests := []struct {
		input     string
		wantModel string
		wantFound bool
	}{
		{"pt-210", "58mm PT-210", true},
		{"PT210", "58mm PT-210", true},
		{"58mm pt-210 portable", "58mm PT-210", true},
		{"gp-58n", "58mm GP-58N", true},
		{"GP58N", "58mm GP-58N", true},
		{"ec-pm-80250", "80mm EC-PM-80250", true},
		{"ECPM80250", "80mm EC-PM-80250", true},
		{"tm-t20", "Epson TM-T20III", true},
		{"tmt20", "Epson TM-T20III", true},
		{"Epson TM-T20III", "Epson TM-T20III", true},
		{"tm-t88", "Epson TM-T88VI", true},
		{"tmt88", "Epson TM-T88VI", true},
		{"Epson TM-T88VI Receipt", "Epson TM-T88VI", true},
		{"tsp100", "Star TSP100", true},
		{"tsp-100", "Star TSP100", true},
		{"Star TSP100 Cutter", "Star TSP100", true},
		{"srp-350", "Bixolon SRP-350III", true},
		{"srp350", "Bixolon SRP-350III", true},
		{"xp-58", "Xprinter XP-58", true},
		{"xp58", "Xprinter XP-58", true},
		{"xp-80", "Xprinter XP-80", true},
		{"xp80", "Xprinter XP-80", true},
		{"58mm", "Generic 58mm", true},
		{"58", "Generic 58mm", true},
		{"80mm", "Generic 80mm", true},
		{"80", "Generic 80mm", true},
		{"Unknown Model", "", false},
		{"", "", false},
		{"   ", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, found := profile.ProfileFromModel(tt.input)
			if found != tt.wantFound {
				t.Fatalf("ProfileFromModel(%q) found = %v, want %v", tt.input, found, tt.wantFound)
			}
			if found && got.Model != tt.wantModel {
				t.Errorf("ProfileFromModel(%q) model = %q, want %q", tt.input, got.Model, tt.wantModel)
			}
			if !found && got != nil {
				t.Errorf("ProfileFromModel(%q) got non-nil profile when not found: %+v", tt.input, got)
			}
		})
	}
}

func TestEscposProfileOptions(t *testing.T) {
	p58 := profile.CreateProfile58mm()
	if p58.PrintWidth != 48 {
		t.Errorf("expected Generic 58mm PrintWidth 48, got %d", p58.PrintWidth)
	}
	if p58.ImageThreshold != 128 {
		t.Errorf("expected Generic 58mm ImageThreshold 128, got %d", p58.ImageThreshold)
	}
	if p58.StrictEncoding {
		t.Error("expected default StrictEncoding to be false")
	}
	if p58.AutoSwitchCodePage {
		t.Error("expected default AutoSwitchCodePage to be false")
	}

	p80 := profile.CreateProfile80mm()
	if p80.PrintWidth != 72 {
		t.Errorf("expected Generic 80mm PrintWidth 72, got %d", p80.PrintWidth)
	}
	if p80.ImageThreshold != 128 {
		t.Errorf("expected Generic 80mm ImageThreshold 128, got %d", p80.ImageThreshold)
	}
}

func TestFromModel(t *testing.T) {
	p, ok := profile.FromModel("tm-t20")
	if !ok || p == nil || p.Model != "Epson TM-T20III" {
		t.Errorf("FromModel(\"tm-t20\") failed: got (%+v, %v)", p, ok)
	}
}
