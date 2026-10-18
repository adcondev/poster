package shared_test

import (
	"testing"

	"github.com/adcondev/poster/pkg/commands/shared"
)

func TestToLittleEndian32(t *testing.T) {
	tests := []struct {
		number   uint32
		wantL    byte
		wantH    byte
		wantHH   byte
		wantHHH  byte
	}{
		{0, 0x00, 0x00, 0x00, 0x00},
		{1, 0x01, 0x00, 0x00, 0x00},
		{0x12345678, 0x78, 0x56, 0x34, 0x12},
		{0xFFFFFFFF, 0xFF, 0xFF, 0xFF, 0xFF},
	}

	for _, tt := range tests {
		nL, nH, nHH, nHHH := shared.ToLittleEndian32(tt.number)
		if nL != tt.wantL || nH != tt.wantH || nHH != tt.wantHH || nHHH != tt.wantHHH {
			t.Errorf("ToLittleEndian32(0x%x) = (%#x, %#x, %#x, %#x); want (%#x, %#x, %#x, %#x)",
				tt.number, nL, nH, nHH, nHHH, tt.wantL, tt.wantH, tt.wantHH, tt.wantHHH)
		}
	}
}

func TestSharedConstants(t *testing.T) {
	if shared.SP != 0x20 {
		t.Errorf("shared.SP = %#x, want 0x20", shared.SP)
	}
	if shared.FS != 0x1C {
		t.Errorf("shared.FS = %#x, want 0x1C", shared.FS)
	}
	if shared.NUL != 0x00 {
		t.Errorf("shared.NUL = %#x, want 0x00", shared.NUL)
	}
	if shared.ESC != 0x1B {
		t.Errorf("shared.ESC = %#x, want 0x1B", shared.ESC)
	}
	if shared.GS != 0x1D {
		t.Errorf("shared.GS = %#x, want 0x1D", shared.GS)
	}
	if shared.HT != 0x09 {
		t.Errorf("shared.HT = %#x, want 0x09", shared.HT)
	}
	if shared.Dpl80mm203dpi != 576 {
		t.Errorf("shared.Dpl80mm203dpi = %d, want 576", shared.Dpl80mm203dpi)
	}
	if shared.Dpl58mm203dpi != 384 {
		t.Errorf("shared.Dpl58mm203dpi = %d, want 384", shared.Dpl58mm203dpi)
	}
}
