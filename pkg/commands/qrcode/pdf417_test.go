package qrcode

import (
	"bytes"
	"errors"
	"testing"

	"github.com/adcondev/poster/pkg/commands/shared"
)

func TestSetPDF417Columns(t *testing.T) {
	tests := []struct {
		name      string
		columns   byte
		wantBytes []byte
		wantErr   bool
	}{
		{"auto (0)", 0, []byte{shared.GS, '(', 'k', 0x03, 0x00, 48, 65, 0}, false},
		{"mid (15)", 15, []byte{shared.GS, '(', 'k', 0x03, 0x00, 48, 65, 15}, false},
		{"max (30)", 30, []byte{shared.GS, '(', 'k', 0x03, 0x00, 48, 65, 30}, false},
		{"over limit (31)", 31, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SetPDF417Columns(tt.columns)
			if (err != nil) != tt.wantErr {
				t.Fatalf("SetPDF417Columns(%d) error = %v, wantErr %v", tt.columns, err, tt.wantErr)
			}
			if !tt.wantErr && !bytes.Equal(got, tt.wantBytes) {
				t.Errorf("SetPDF417Columns(%d) = %v, want %v", tt.columns, got, tt.wantBytes)
			}
			if tt.wantErr && !errors.Is(err, ErrPDF417Columns) {
				t.Errorf("expected ErrPDF417Columns, got %v", err)
			}
		})
	}
}

func TestSetPDF417Rows(t *testing.T) {
	tests := []struct {
		name      string
		rows      byte
		wantBytes []byte
		wantErr   bool
	}{
		{"auto (0)", 0, []byte{shared.GS, '(', 'k', 0x03, 0x00, 48, 66, 0}, false},
		{"min (3)", 3, []byte{shared.GS, '(', 'k', 0x03, 0x00, 48, 66, 3}, false},
		{"mid (45)", 45, []byte{shared.GS, '(', 'k', 0x03, 0x00, 48, 66, 45}, false},
		{"max (90)", 90, []byte{shared.GS, '(', 'k', 0x03, 0x00, 48, 66, 90}, false},
		{"invalid below min (1)", 1, nil, true},
		{"invalid below min (2)", 2, nil, true},
		{"invalid above max (91)", 91, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SetPDF417Rows(tt.rows)
			if (err != nil) != tt.wantErr {
				t.Fatalf("SetPDF417Rows(%d) error = %v, wantErr %v", tt.rows, err, tt.wantErr)
			}
			if !tt.wantErr && !bytes.Equal(got, tt.wantBytes) {
				t.Errorf("SetPDF417Rows(%d) = %v, want %v", tt.rows, got, tt.wantBytes)
			}
			if tt.wantErr && !errors.Is(err, ErrPDF417Rows) {
				t.Errorf("expected ErrPDF417Rows, got %v", err)
			}
		})
	}
}

func TestSetPDF417ModuleWidth(t *testing.T) {
	tests := []struct {
		name      string
		width     byte
		wantBytes []byte
		wantErr   bool
	}{
		{"min (2)", 2, []byte{shared.GS, '(', 'k', 0x03, 0x00, 48, 67, 2}, false},
		{"mid (5)", 5, []byte{shared.GS, '(', 'k', 0x03, 0x00, 48, 67, 5}, false},
		{"max (8)", 8, []byte{shared.GS, '(', 'k', 0x03, 0x00, 48, 67, 8}, false},
		{"invalid below min (1)", 1, nil, true},
		{"invalid above max (9)", 9, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SetPDF417ModuleWidth(tt.width)
			if (err != nil) != tt.wantErr {
				t.Fatalf("SetPDF417ModuleWidth(%d) error = %v, wantErr %v", tt.width, err, tt.wantErr)
			}
			if !tt.wantErr && !bytes.Equal(got, tt.wantBytes) {
				t.Errorf("SetPDF417ModuleWidth(%d) = %v, want %v", tt.width, got, tt.wantBytes)
			}
			if tt.wantErr && !errors.Is(err, ErrPDF417ModuleWidth) {
				t.Errorf("expected ErrPDF417ModuleWidth, got %v", err)
			}
		})
	}
}

func TestSetPDF417RowHeight(t *testing.T) {
	tests := []struct {
		name      string
		height    byte
		wantBytes []byte
		wantErr   bool
	}{
		{"min (2)", 2, []byte{shared.GS, '(', 'k', 0x03, 0x00, 48, 68, 2}, false},
		{"mid (5)", 5, []byte{shared.GS, '(', 'k', 0x03, 0x00, 48, 68, 5}, false},
		{"max (8)", 8, []byte{shared.GS, '(', 'k', 0x03, 0x00, 48, 68, 8}, false},
		{"invalid below min (1)", 1, nil, true},
		{"invalid above max (9)", 9, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SetPDF417RowHeight(tt.height)
			if (err != nil) != tt.wantErr {
				t.Fatalf("SetPDF417RowHeight(%d) error = %v, wantErr %v", tt.height, err, tt.wantErr)
			}
			if !tt.wantErr && !bytes.Equal(got, tt.wantBytes) {
				t.Errorf("SetPDF417RowHeight(%d) = %v, want %v", tt.height, got, tt.wantBytes)
			}
			if tt.wantErr && !errors.Is(err, ErrPDF417RowHeight) {
				t.Errorf("expected ErrPDF417RowHeight, got %v", err)
			}
		})
	}
}

func TestSetPDF417ErrorCorrection(t *testing.T) {
	tests := []struct {
		name      string
		level     byte
		wantBytes []byte
		wantErr   bool
	}{
		{"level 0", 0, []byte{shared.GS, '(', 'k', 0x04, 0x00, 48, 69, 48, 48}, false},
		{"level 2", 2, []byte{shared.GS, '(', 'k', 0x04, 0x00, 48, 69, 48, 50}, false},
		{"level 8", 8, []byte{shared.GS, '(', 'k', 0x04, 0x00, 48, 69, 48, 56}, false},
		{"invalid level (9)", 9, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SetPDF417ErrorCorrection(tt.level)
			if (err != nil) != tt.wantErr {
				t.Fatalf("SetPDF417ErrorCorrection(%d) error = %v, wantErr %v", tt.level, err, tt.wantErr)
			}
			if !tt.wantErr && !bytes.Equal(got, tt.wantBytes) {
				t.Errorf("SetPDF417ErrorCorrection(%d) = %v, want %v", tt.level, got, tt.wantBytes)
			}
			if tt.wantErr && !errors.Is(err, ErrPDF417ErrorCorrection) {
				t.Errorf("expected ErrPDF417ErrorCorrection, got %v", err)
			}
		})
	}
}

func TestStorePDF417Data(t *testing.T) {
	t.Run("valid single byte", func(t *testing.T) {
		data := []byte("A")
		got, err := StorePDF417Data(data)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// length = 1 + 3 = 4 -> pL = 0x04, pH = 0x00
		want := []byte{shared.GS, '(', 'k', 0x04, 0x00, 48, 80, 48, 'A'}
		if !bytes.Equal(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("empty data", func(t *testing.T) {
		_, err := StorePDF417Data(nil)
		if !errors.Is(err, ErrPDF417DataEmpty) {
			t.Errorf("expected ErrPDF417DataEmpty, got %v", err)
		}
	})

	t.Run("data too long", func(t *testing.T) {
		data := make([]byte, 65533)
		_, err := StorePDF417Data(data)
		if !errors.Is(err, ErrPDF417DataTooLong) {
			t.Errorf("expected ErrPDF417DataTooLong, got %v", err)
		}
	})

	t.Run("max valid data length", func(t *testing.T) {
		data := make([]byte, MaxPDF417DataLength)
		got, err := StorePDF417Data(data)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// length = 65532 + 3 = 65535 (0xFFFF) -> pL = 0xFF, pH = 0xFF
		if got[3] != 0xFF || got[4] != 0xFF {
			t.Errorf("expected length bytes 0xFF 0xFF, got %02X %02X", got[3], got[4])
		}
		if len(got) != 8+MaxPDF417DataLength {
			t.Errorf("expected total length %d, got %d", 8+MaxPDF417DataLength, len(got))
		}
	})
}

func TestPrintPDF417(t *testing.T) {
	got := PrintPDF417()
	want := []byte{shared.GS, '(', 'k', 0x03, 0x00, 48, 81, 48}
	if !bytes.Equal(got, want) {
		t.Errorf("PrintPDF417() = %v, want %v", got, want)
	}
}

func TestComposePDF417(t *testing.T) {
	t.Run("valid compose", func(t *testing.T) {
		data := []byte("TEST-PDF417-DATA")
		got, err := ComposePDF417(5, 10, 3, 4, 2, data)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Verify first command is columns
		wantCol := []byte{shared.GS, '(', 'k', 0x03, 0x00, 48, 65, 5}
		if !bytes.HasPrefix(got, wantCol) {
			t.Errorf("expected prefix %v, got %v", wantCol, got[:len(wantCol)])
		}

		// Verify last command is print
		wantPrint := []byte{shared.GS, '(', 'k', 0x03, 0x00, 48, 81, 48}
		if !bytes.HasSuffix(got, wantPrint) {
			t.Errorf("expected suffix %v, got %v", wantPrint, got[len(got)-len(wantPrint):])
		}

		// Verify data is contained
		if !bytes.Contains(got, data) {
			t.Errorf("expected composed bytes to contain data %s", string(data))
		}
	})

	t.Run("invalid column fails", func(t *testing.T) {
		_, err := ComposePDF417(35, 0, 3, 3, 2, []byte("ok"))
		if err == nil {
			t.Fatal("expected error on invalid columns")
		}
	})

	t.Run("invalid rows fails", func(t *testing.T) {
		_, err := ComposePDF417(0, 1, 3, 3, 2, []byte("ok"))
		if err == nil {
			t.Fatal("expected error on invalid rows")
		}
	})

	t.Run("invalid module width fails", func(t *testing.T) {
		_, err := ComposePDF417(0, 0, 1, 3, 2, []byte("ok"))
		if err == nil {
			t.Fatal("expected error on invalid module width")
		}
	})

	t.Run("invalid row height fails", func(t *testing.T) {
		_, err := ComposePDF417(0, 0, 3, 10, 2, []byte("ok"))
		if err == nil {
			t.Fatal("expected error on invalid row height")
		}
	})

	t.Run("invalid error correction fails", func(t *testing.T) {
		_, err := ComposePDF417(0, 0, 3, 3, 9, []byte("ok"))
		if err == nil {
			t.Fatal("expected error on invalid error correction")
		}
	})

	t.Run("empty data fails", func(t *testing.T) {
		_, err := ComposePDF417(0, 0, 3, 3, 2, nil)
		if err == nil {
			t.Fatal("expected error on empty data")
		}
	})
}
