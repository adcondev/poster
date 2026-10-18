package fonts_test

import (
	"testing"

	"github.com/adcondev/poster/pkg/emulator/fonts"
)

func TestFonts_Embedding(t *testing.T) {
	// HasEmbeddedFonts
	if !fonts.HasEmbeddedFonts() {
		t.Fatal("HasEmbeddedFonts() returned false, expected true")
	}

	// ListFonts
	fontList := fonts.ListFonts()
	if len(fontList) == 0 {
		t.Fatal("ListFonts() returned empty list")
	}

	foundJetBrains := false
	for _, f := range fontList {
		if f == "JetBrainsMono-Bold.ttf" {
			foundJetBrains = true
			break
		}
	}
	if !foundJetBrains {
		t.Errorf("ListFonts() = %v, expected to contain \"JetBrainsMono-Bold.ttf\"", fontList)
	}

	// LoadFontData valid
	data, err := fonts.LoadFontData("JetBrainsMono-Bold.ttf")
	if err != nil {
		t.Fatalf("LoadFontData() unexpected error = %v", err)
	}
	if len(data) == 0 {
		t.Fatal("LoadFontData() returned empty data")
	}

	// LoadFontData non-existent
	_, err = fonts.LoadFontData("NonExistentFont.ttf")
	if err == nil {
		t.Fatal("LoadFontData() expected error for non-existent font, got nil")
	}

	// DiagnoseEmbedding
	fonts.DiagnoseEmbedding()
}
