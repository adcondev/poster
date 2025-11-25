package profile_test

import (
	"testing"

	"github.com/adcondev/pos-printer/pkg/commands/character"
	"github.com/adcondev/pos-printer/pkg/profile"
)

func TestCreatePt210(t *testing.T) {
	p := profile.CreatePt210()

	if p.Model != "58mm PT-210" {
		t.Errorf("expected model '58mm PT-210', got '%s'", p.Model)
	}
	if p.PaperWidth != 58 {
		t.Errorf("expected paper width 58, got %f", p.PaperWidth)
	}
	if p.CodeTable != character.PC850 {
		t.Errorf("expected code table PC850, got %v", p.CodeTable)
	}
	if p.HasQR != true {
		t.Error("expected HasQR to be true")
	}
	if p.QRMaxSize != 19 {
		t.Errorf("expected QRMaxSize 19, got %d", p.QRMaxSize)
	}
}

func TestCreateGP58N(t *testing.T) {
	p := profile.CreateGP58N()

	if p.Model != "58mm GP-58N" {
		t.Errorf("expected model '58mm GP-58N', got '%s'", p.Model)
	}
	if p.CodeTable != character.PC850 {
		t.Errorf("expected code table PC850, got %v", p.CodeTable)
	}
}

func TestCreateProfile58mm(t *testing.T) {
	p := profile.CreateProfile58mm()

	if p.Model != "Generic 58mm" {
		t.Errorf("expected model 'Generic 58mm', got '%s'", p.Model)
	}
	if p.PaperWidth != 58 {
		t.Errorf("expected paper width 58, got %f", p.PaperWidth)
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
	if !p.SupportsGraphics {
		t.Error("expected SupportsGraphics to be true")
	}
	if !p.SupportsBarcode {
		t.Error("expected SupportsBarcode to be true")
	}
	if p.HasQR {
		t.Error("expected HasQR to be false")
	}
	if p.SupportsCutter {
		t.Error("expected SupportsCutter to be false")
	}
	if p.SupportsDrawer {
		t.Error("expected SupportsDrawer to be false")
	}
	if p.CodeTable != character.PC850 {
		t.Errorf("expected code table PC850, got %v", p.CodeTable)
	}
}

func TestCreateECPM80250(t *testing.T) {
	p := profile.CreateECPM80250()

	if p.Model != "80mm EC-PM-80250" {
		t.Errorf("expected model '80mm EC-PM-80250', got '%s'", p.Model)
	}
}

func TestCreateProfile80mm(t *testing.T) {
	p := profile.CreateProfile80mm()

	if p.Model != "Generic 80mm" {
		t.Errorf("expected model 'Generic 80mm', got '%s'", p.Model)
	}
	if p.PaperWidth != 80 {
		t.Errorf("expected paper width 80, got %f", p.PaperWidth)
	}
	if p.DPI != 203 {
		t.Errorf("expected DPI 203, got %d", p.DPI)
	}
	if p.DotsPerLine != 576 {
		t.Errorf("expected DotsPerLine 576, got %d", p.DotsPerLine)
	}
	if !p.SupportsGraphics {
		t.Error("expected SupportsGraphics to be true")
	}
	if !p.SupportsBarcode {
		t.Error("expected SupportsBarcode to be true")
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
	if p.CodeTable != character.PC850 {
		t.Errorf("expected code table PC850, got %v", p.CodeTable)
	}
	if p.ImageThreshold != 128 {
		t.Errorf("expected ImageThreshold 128, got %d", p.ImageThreshold)
	}
}
