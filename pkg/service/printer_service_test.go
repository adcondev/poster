package service

import (
	"errors"
	"testing"

	"github.com/adcondev/poster/pkg/commands/character"
	"github.com/adcondev/poster/pkg/composer"
	"github.com/adcondev/poster/pkg/connection"
	"github.com/adcondev/poster/pkg/constants"
	"github.com/adcondev/poster/pkg/graphics"
	"github.com/adcondev/poster/pkg/profile"
)

func TestNewPrinter_Validation(t *testing.T) {
	conn := connection.NewBufferConnector()
	proto := composer.NewEscpos()
	prof := profile.CreateProfile80mm()

	t.Run("nil proto", func(t *testing.T) {
		_, err := NewPrinter(nil, prof, conn)
		if err == nil {
			t.Fatal("expected error for nil proto, got nil")
		}
	})

	t.Run("nil prof", func(t *testing.T) {
		_, err := NewPrinter(proto, nil, conn)
		if err == nil {
			t.Fatal("expected error for nil prof, got nil")
		}
	})

	t.Run("nil conn", func(t *testing.T) {
		_, err := NewPrinter(proto, prof, nil)
		if err == nil {
			t.Fatal("expected error for nil conn, got nil")
		}
	})

	t.Run("valid parameters", func(t *testing.T) {
		p, err := NewPrinter(proto, prof, conn)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if p == nil {
			t.Fatal("expected printer instance, got nil")
		}
	})
}

func TestPrinterOptions(t *testing.T) {
	t.Run("nil connection returns error", func(t *testing.T) {
		_, err := NewPrinterWithOptions(nil)
		if err == nil {
			t.Fatal("expected error when connection is nil, got nil")
		}
	})

	t.Run("WithProfile nil returns error", func(t *testing.T) {
		conn := connection.NewBufferConnector()
		_, err := NewPrinterWithOptions(conn, WithProfile(nil))
		if err == nil {
			t.Fatal("expected error with nil profile, got nil")
		}
	})

	t.Run("WithProtocol nil returns error", func(t *testing.T) {
		conn := connection.NewBufferConnector()
		_, err := NewPrinterWithOptions(conn, WithProtocol(nil))
		if err == nil {
			t.Fatal("expected error with nil protocol, got nil")
		}
	})

	t.Run("with custom options success", func(t *testing.T) {
		conn := connection.NewBufferConnector()
		customProfile := profile.CreateProfile58mm()
		customProtocol := composer.NewEscpos()

		p, err := NewPrinterWithOptions(
			conn,
			WithProfile(customProfile),
			WithProtocol(customProtocol),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if p.Profile.PaperWidth != customProfile.PaperWidth {
			t.Errorf("expected custom profile paper width %v, got %v", customProfile.PaperWidth, p.Profile.PaperWidth)
		}
	})
}

func TestPrinter_BasicControl(t *testing.T) {
	conn := connection.NewBufferConnector()
	p, err := NewPrinterWithOptions(conn)
	if err != nil {
		t.Fatalf("NewPrinterWithOptions failed: %v", err)
	}

	t.Run("GetProfile", func(t *testing.T) {
		prof := p.GetProfile()
		if prof == nil {
			t.Fatal("expected non-nil profile")
		}
	})

	t.Run("Initialize with supported code table", func(t *testing.T) {
		p.Profile.CodeTable = character.WPC1252
		if err := p.Initialize(); err != nil {
			t.Fatalf("Initialize failed: %v", err)
		}
		if conn.Len() == 0 {
			t.Error("expected bytes written on initialize")
		}
	})

	t.Run("Initialize with unsupported code table fallback", func(t *testing.T) {
		conn.Reset()
		p.Profile.CodeTable = character.CodeTable(250)
		if err := p.Initialize(); err != nil {
			t.Fatalf("Initialize failed: %v", err)
		}
		if conn.Len() == 0 {
			t.Error("expected bytes written on initialize fallback")
		}
	})

	t.Run("Write empty and non-empty", func(t *testing.T) {
		conn.Reset()
		if err := p.Write([]byte{}); err != nil {
			t.Errorf("empty Write failed: %v", err)
		}
		if conn.Len() != 0 {
			t.Errorf("expected 0 bytes written, got %d", conn.Len())
		}

		if err := p.Write([]byte{0x1B, 0x40}); err != nil {
			t.Errorf("Write failed: %v", err)
		}
		if conn.Len() != 2 {
			t.Errorf("expected 2 bytes written, got %d", conn.Len())
		}
	})

	t.Run("Close", func(t *testing.T) {
		if err := p.Close(); err != nil {
			t.Errorf("Close failed: %v", err)
		}
		// Subsequent write on closed connection should error
		if err := p.Write([]byte{0x01}); err == nil {
			t.Error("expected error writing to closed connection")
		}
	})
}

func TestPrinter_TextPrintingAndFormatting(t *testing.T) {
	conn := connection.NewBufferConnector()
	p, err := NewPrinterWithOptions(conn)
	if err != nil {
		t.Fatalf("NewPrinterWithOptions failed: %v", err)
	}

	t.Run("Print and PrintLine and FeedLines", func(t *testing.T) {
		if err := p.Print("Hello"); err != nil {
			t.Errorf("Print failed: %v", err)
		}
		if err := p.PrintLine(" World"); err != nil {
			t.Errorf("PrintLine failed: %v", err)
		}
		if err := p.FeedLines(2); err != nil {
			t.Errorf("FeedLines failed: %v", err)
		}
	})

	t.Run("Fonts and Bold", func(t *testing.T) {
		if err := p.FontA(); err != nil {
			t.Errorf("FontA failed: %v", err)
		}
		if err := p.FontB(); err != nil {
			t.Errorf("FontB failed: %v", err)
		}
		if err := p.EnableBold(); err != nil {
			t.Errorf("EnableBold failed: %v", err)
		}
		if err := p.DisableBold(); err != nil {
			t.Errorf("DisableBold failed: %v", err)
		}
	})

	t.Run("Alignments", func(t *testing.T) {
		if err := p.AlignLeft(); err != nil {
			t.Errorf("AlignLeft failed: %v", err)
		}
		if err := p.AlignCenter(); err != nil {
			t.Errorf("AlignCenter failed: %v", err)
		}
		if err := p.AlignRight(); err != nil {
			t.Errorf("AlignRight failed: %v", err)
		}
		alignments := []string{"left", "center", "right", "unknown_default"}
		for _, align := range alignments {
			if err := p.SetAlignment(align); err != nil {
				t.Errorf("SetAlignment(%s) failed: %v", align, err)
			}
		}
	})

	t.Run("Sizes", func(t *testing.T) {
		sizes := []func() error{
			p.SingleSize,
			p.DoubleSize,
			p.TripleSize,
			p.QuadraSize,
			p.PentaSize,
			p.HexaSize,
			p.HeptaSize,
			p.OctaSize,
		}
		for _, sizeFunc := range sizes {
			if err := sizeFunc(); err != nil {
				t.Errorf("size function failed: %v", err)
			}
		}
		if err := p.CustomSize(2, 3); err != nil {
			t.Errorf("CustomSize failed: %v", err)
		}
	})

	t.Run("Underline and Inverse", func(t *testing.T) {
		if err := p.NoDot(); err != nil {
			t.Errorf("NoDot failed: %v", err)
		}
		if err := p.OneDot(); err != nil {
			t.Errorf("OneDot failed: %v", err)
		}
		if err := p.TwoDot(); err != nil {
			t.Errorf("TwoDot failed: %v", err)
		}
		if err := p.InverseOn(); err != nil {
			t.Errorf("InverseOn failed: %v", err)
		}
		if err := p.InverseOff(); err != nil {
			t.Errorf("InverseOff failed: %v", err)
		}
	})

	t.Run("Feed and Cut, Pulse, Beep, Spacing", func(t *testing.T) {
		if err := p.FullFeedAndCut(3); err != nil {
			t.Errorf("FullFeedAndCut failed: %v", err)
		}
		if err := p.PartialFeedAndCut(3); err != nil {
			t.Errorf("PartialFeedAndCut failed: %v", err)
		}
		if err := p.Pulse(0, 50, 50); err != nil {
			t.Errorf("Pulse failed: %v", err)
		}
		if err := p.Beep(2, 10); err != nil {
			t.Errorf("Beep failed: %v", err)
		}
		if err := p.SetLineSpacing(30); err != nil {
			t.Errorf("SetLineSpacing failed: %v", err)
		}
		if err := p.SelectDefaultLineSpacing(); err != nil {
			t.Errorf("SelectDefaultLineSpacing failed: %v", err)
		}
	})
}

func TestPrinter_CodeTable(t *testing.T) {
	conn := connection.NewBufferConnector()
	p, err := NewPrinterWithOptions(conn)
	if err != nil {
		t.Fatalf("NewPrinterWithOptions failed: %v", err)
	}

	t.Run("supported code table", func(t *testing.T) {
		if err := p.SetCodeTable(character.PC437); err != nil {
			t.Errorf("SetCodeTable(PC437) failed: %v", err)
		}
	})

	t.Run("unsupported code table falls back to WPC1252", func(t *testing.T) {
		if err := p.SetCodeTable(character.CodeTable(250)); err != nil {
			t.Errorf("SetCodeTable(250) failed: %v", err)
		}
	})
}

func TestPrinter_PrintBitmap(t *testing.T) {
	conn := connection.NewBufferConnector()
	p, err := NewPrinterWithOptions(conn)
	if err != nil {
		t.Fatalf("NewPrinterWithOptions failed: %v", err)
	}

	t.Run("nil bitmap returns error", func(t *testing.T) {
		if err := p.PrintBitmap(nil); err == nil {
			t.Fatal("expected error for nil bitmap, got nil")
		}
	})

	t.Run("standard non-banded bitmap", func(t *testing.T) {
		bm := graphics.NewMonochromeBitmap(8, 4)
		if err := p.PrintBitmap(bm); err != nil {
			t.Fatalf("PrintBitmap failed: %v", err)
		}
	})

	t.Run("banded bitmap", func(t *testing.T) {
		p.Profile.MaxBandHeight = 4
		bm := graphics.NewMonochromeBitmap(8, 10)
		if err := p.PrintBitmap(bm); err != nil {
			t.Fatalf("PrintBitmap banded failed: %v", err)
		}
	})
}

func TestPrinter_PrintBarcode(t *testing.T) {
	conn := connection.NewBufferConnector()
	p, err := NewPrinterWithOptions(conn)
	if err != nil {
		t.Fatalf("NewPrinterWithOptions failed: %v", err)
	}

	t.Run("empty data returns error", func(t *testing.T) {
		cfg := graphics.DefaultBarcodeConfig()
		if err := p.PrintBarcode(cfg, nil); err == nil {
			t.Fatal("expected error for empty barcode data, got nil")
		}
	})

	t.Run("valid barcode", func(t *testing.T) {
		cfg := graphics.DefaultBarcodeConfig()
		if err := p.PrintBarcode(cfg, []byte("12345678")); err != nil {
			t.Fatalf("PrintBarcode failed: %v", err)
		}
	})
}

func TestMockPrinter_BasicAndText(t *testing.T) {
	mock := NewMockPrinter()

	if prof := mock.GetProfile(); prof == nil {
		t.Error("expected non-nil profile")
	}
	if err := mock.Initialize(); err != nil {
		t.Errorf("Initialize failed: %v", err)
	}
	if err := mock.Close(); err != nil {
		t.Errorf("Close failed: %v", err)
	}
	if err := mock.Write([]byte{0x01}); err != nil {
		t.Errorf("Write failed: %v", err)
	}
	if err := mock.Print("test"); err != nil {
		t.Errorf("Print failed: %v", err)
	}
	if err := mock.PrintLine("line"); err != nil {
		t.Errorf("PrintLine failed: %v", err)
	}
	if err := mock.FeedLines(1); err != nil {
		t.Errorf("FeedLines failed: %v", err)
	}
	if mock.CallCount("Print") != 1 {
		t.Errorf("expected 1 Print call, got %d", mock.CallCount("Print"))
	}
}

func TestMockPrinter_FormattingAndSizes(t *testing.T) {
	mock := NewMockPrinter()

	methods := []func() error{
		mock.FontA, mock.FontB, mock.EnableBold, mock.DisableBold,
		mock.AlignLeft, mock.AlignCenter, mock.AlignRight,
		mock.SingleSize, mock.DoubleSize, mock.TripleSize, mock.QuadraSize,
		mock.PentaSize, mock.HexaSize, mock.HeptaSize, mock.OctaSize,
		mock.NoDot, mock.OneDot, mock.TwoDot,
		mock.InverseOn, mock.InverseOff,
	}
	for _, m := range methods {
		if err := m(); err != nil {
			t.Errorf("method failed: %v", err)
		}
	}
	if err := mock.SetAlignment("center"); err != nil {
		t.Errorf("SetAlignment failed: %v", err)
	}
	if err := mock.CustomSize(2, 2); err != nil {
		t.Errorf("CustomSize failed: %v", err)
	}
}

func TestMockPrinter_AdvancedAndControl(t *testing.T) {
	mock := NewMockPrinter()

	if err := mock.FullFeedAndCut(3); err != nil {
		t.Errorf("FullFeedAndCut failed: %v", err)
	}
	if err := mock.PartialFeedAndCut(3); err != nil {
		t.Errorf("PartialFeedAndCut failed: %v", err)
	}
	if err := mock.Pulse(0, 10, 10); err != nil {
		t.Errorf("Pulse failed: %v", err)
	}
	if err := mock.Beep(1, 1); err != nil {
		t.Errorf("Beep failed: %v", err)
	}
	if err := mock.SetLineSpacing(30); err != nil {
		t.Errorf("SetLineSpacing failed: %v", err)
	}
	if err := mock.SelectDefaultLineSpacing(); err != nil {
		t.Errorf("SelectDefaultLineSpacing failed: %v", err)
	}
	if err := mock.PrintBitmap(graphics.NewMonochromeBitmap(8, 8)); err != nil {
		t.Errorf("PrintBitmap failed: %v", err)
	}
	if err := mock.SetCodeTable(character.WPC1252); err != nil {
		t.Errorf("SetCodeTable failed: %v", err)
	}
	if err := mock.PrintQR("data", graphics.DefaultQROptions()); err != nil {
		t.Errorf("PrintQR failed: %v", err)
	}
	if err := mock.PrintBarcode(graphics.DefaultBarcodeConfig(), []byte("123")); err != nil {
		t.Errorf("PrintBarcode failed: %v", err)
	}

	last := mock.LastCall()
	if last == nil || last.Method != "PrintBarcode" {
		t.Errorf("expected last call PrintBarcode, got %v", last)
	}

	// Test FailOn error injection
	mock.FailOn["Write"] = errors.New("injected error")
	if err := mock.Write([]byte{0x00}); err == nil {
		t.Error("expected injected error from Write, got nil")
	}

	// Test Reset
	mock.Reset()
	if len(mock.Calls) != 0 {
		t.Errorf("expected 0 calls after Reset, got %d", len(mock.Calls))
	}
	if mock.LastCall() != nil {
		t.Errorf("expected nil LastCall after Reset, got %v", mock.LastCall())
	}
	if mock.CurrentAlignment != constants.Left.String() {
		t.Errorf("expected reset alignment 'left', got %s", mock.CurrentAlignment)
	}
}
