package graphics

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"strings"
	"testing"

	posqr "github.com/adcondev/poster/pkg/commands/qrcode"
	"github.com/adcondev/poster/pkg/constants"
)

func createTestPNG(t *testing.T, w, h int) (image.Image, string) {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.Set(x, y, color.RGBA{R: 200, G: 100, B: 50, A: 255})
		}
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("png encode failed: %v", err)
	}
	return img, base64.StdEncoding.EncodeToString(buf.Bytes())
}

func createTestJPEG(t *testing.T, w, h int) (image.Image, string) {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.Set(x, y, color.RGBA{R: 50, G: 100, B: 200, A: 255})
		}
	}
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, nil); err != nil {
		t.Fatalf("jpeg encode failed: %v", err)
	}
	return img, base64.StdEncoding.EncodeToString(buf.Bytes())
}

func TestWriteCloser(t *testing.T) {
	wc := NewWriteCloser()
	if wc == nil {
		t.Fatal("NewWriteCloser returned nil")
	}

	testData := []byte("hello qr world")
	n, err := wc.Write(testData)
	if err != nil {
		t.Fatalf("Write returned error: %v", err)
	}
	if n != len(testData) {
		t.Errorf("Write n = %d, want %d", n, len(testData))
	}

	if !bytes.Equal(wc.Bytes(), testData) {
		t.Errorf("Bytes() = %s, want %s", wc.Bytes(), testData)
	}

	if err := wc.Close(); err != nil {
		t.Errorf("Close() returned error: %v", err)
	}

	// Close must not reset buffer
	if !bytes.Equal(wc.Bytes(), testData) {
		t.Errorf("Bytes() after Close() = %s, want %s", wc.Bytes(), testData)
	}

	wc.Reset()
	if len(wc.Bytes()) != 0 {
		t.Errorf("Bytes() after Reset() = %d bytes, want 0", len(wc.Bytes()))
	}
}

func TestDefaultQROptions(t *testing.T) {
	opts := DefaultQROptions()
	if opts == nil {
		t.Fatal("DefaultQROptions returned nil")
	}
	if opts.Model != posqr.Model2 {
		t.Errorf("Model = %v, want %v", opts.Model, posqr.Model2)
	}
	if opts.ErrorCorrection != posqr.LevelM {
		t.Errorf("ErrorCorrection = %v, want %v", opts.ErrorCorrection, posqr.LevelM)
	}
	if opts.PixelWidth != 128 {
		t.Errorf("PixelWidth = %d, want 128", opts.PixelWidth)
	}
	if opts.MaxPixelWidth != constants.PaperPxWidth80mm {
		t.Errorf("MaxPixelWidth = %d, want %d", opts.MaxPixelWidth, constants.PaperPxWidth80mm)
	}
	if opts.CircleShape {
		t.Error("CircleShape = true, want false")
	}
	if opts.LogoData != "" {
		t.Errorf("LogoData = %q, want empty", opts.LogoData)
	}
}

func TestQrOptions_GetModuleSize(t *testing.T) {
	opts := DefaultQROptions()
	opts.Qr.moduleSize = posqr.ModuleSize(6)
	if got := opts.GetModuleSize(); got != posqr.ModuleSize(6) {
		t.Errorf("GetModuleSize() = %v, want 6", got)
	}
}

func TestQrOptions_GenerateQR(t *testing.T) {
	t.Run("empty data error", func(t *testing.T) {
		opts := DefaultQROptions()
		_, err := opts.GenerateQR("")
		if err == nil {
			t.Error("expected error on empty data, got nil")
		}
	})

	t.Run("data too long error", func(t *testing.T) {
		opts := DefaultQROptions()
		tooLong := strings.Repeat("A", posqr.MaxDataLength+1)
		_, err := opts.GenerateQR(tooLong)
		if err == nil {
			t.Error("expected error on data too long, got nil")
		}
	})

	t.Run("valid QR generation with out-of-range error correction reset", func(t *testing.T) {
		opts := DefaultQROptions()
		opts.ErrorCorrection = posqr.ErrorCorrection(99)
		qrc, err := opts.GenerateQR("https://example.com/poster")
		if err != nil {
			t.Fatalf("GenerateQR failed: %v", err)
		}
		if qrc == nil {
			t.Fatal("expected non-nil QRCode")
		}
		if opts.ErrorCorrection != posqr.LevelM {
			t.Errorf("ErrorCorrection = %v, want fallback LevelM", opts.ErrorCorrection)
		}
		if opts.Qr.gridSize < minGridSize {
			t.Errorf("gridSize = %d, want >= %d", opts.Qr.gridSize, minGridSize)
		}
	})

	t.Run("valid QR generation with logo data", func(t *testing.T) {
		_, b64 := createTestPNG(t, 16, 16)
		opts := DefaultQROptions()
		opts.LogoData = b64
		qrc, err := opts.GenerateQR("https://example.com/poster/logo")
		if err != nil {
			t.Fatalf("GenerateQR with logo failed: %v", err)
		}
		if qrc == nil {
			t.Fatal("expected non-nil QRCode")
		}
		if opts.Logo.Image == nil {
			t.Fatal("expected non-nil Logo.Image")
		}
		if opts.Logo.format != "png" {
			t.Errorf("Logo format = %q, want png", opts.Logo.format)
		}
	})
}

func TestQrOptions_ValidateAndAdjustPixelWidth(t *testing.T) {
	t.Run("below minPixelWidth", func(t *testing.T) {
		opts := DefaultQROptions()
		opts.PixelWidth = minPixelWidth - 20
		opts.validateAndAdjustPixelWidth()
		if opts.PixelWidth != minPixelWidth {
			t.Errorf("PixelWidth = %d, want %d", opts.PixelWidth, minPixelWidth)
		}
		if opts.Qr.requestedWidth != minPixelWidth {
			t.Errorf("requestedWidth = %d, want %d", opts.Qr.requestedWidth, minPixelWidth)
		}
	})

	t.Run("above custom MaxPixelWidth", func(t *testing.T) {
		opts := DefaultQROptions()
		opts.MaxPixelWidth = 250
		opts.PixelWidth = 350
		opts.validateAndAdjustPixelWidth()
		if opts.PixelWidth != 250 {
			t.Errorf("PixelWidth = %d, want 250", opts.PixelWidth)
		}
		if opts.Qr.requestedWidth != 250 {
			t.Errorf("requestedWidth = %d, want 250", opts.Qr.requestedWidth)
		}
	})

	t.Run("above default MaxPixelWidth when MaxPixelWidth is 0", func(t *testing.T) {
		opts := DefaultQROptions()
		opts.MaxPixelWidth = 0
		opts.PixelWidth = constants.PaperPxWidth80mm + 100
		opts.validateAndAdjustPixelWidth()
		if opts.PixelWidth != constants.PaperPxWidth80mm {
			t.Errorf("PixelWidth = %d, want %d", opts.PixelWidth, constants.PaperPxWidth80mm)
		}
	})

	t.Run("within valid bounds", func(t *testing.T) {
		opts := DefaultQROptions()
		opts.PixelWidth = 200
		opts.validateAndAdjustPixelWidth()
		if opts.PixelWidth != 200 {
			t.Errorf("PixelWidth = %d, want 200", opts.PixelWidth)
		}
	})
}

func TestQrOptions_ValidateGridSize(t *testing.T) {
	opts := DefaultQROptions()
	opts.Qr.gridSize = minGridSize - 1
	if err := opts.validateGridSize(); err == nil {
		t.Error("expected error for gridSize < minGridSize, got nil")
	}

	opts.Qr.gridSize = minGridSize
	if err := opts.validateGridSize(); err != nil {
		t.Errorf("unexpected error for minGridSize: %v", err)
	}
}

func TestQrOptions_CalculateModuleSize(t *testing.T) {
	t.Run("totalModules <= 0 fallback", func(t *testing.T) {
		opts := DefaultQROptions()
		opts.Qr.gridSize = -10 // totalModules = -10 + 8 = -2
		opts.calculateModuleSize()
		if opts.Qr.moduleSize != posqr.DefaultModuleSize {
			t.Errorf("moduleSize = %v, want DefaultModuleSize", opts.Qr.moduleSize)
		}
	})

	t.Run("rawModSize < DefaultModuleSize", func(t *testing.T) {
		opts := DefaultQROptions()
		opts.Qr.gridSize = 177
		opts.Qr.requestedWidth = minPixelWidth // 87 / (177 + 8) = 0
		opts.calculateModuleSize()
		if opts.Qr.moduleSize != posqr.DefaultModuleSize {
			t.Errorf("moduleSize = %v, want DefaultModuleSize", opts.Qr.moduleSize)
		}
	})

	t.Run("rawModSize > MaxModuleSize", func(t *testing.T) {
		opts := DefaultQROptions()
		opts.Qr.gridSize = 21
		opts.Qr.requestedWidth = 1000 // 1000 / 29 = 34 > 16
		opts.calculateModuleSize()
		if opts.Qr.moduleSize != posqr.MaxModuleSize {
			t.Errorf("moduleSize = %v, want MaxModuleSize", opts.Qr.moduleSize)
		}
	})

	t.Run("normal rawModSize", func(t *testing.T) {
		opts := DefaultQROptions()
		opts.Qr.gridSize = 21
		opts.Qr.requestedWidth = 200 // 200 / 29 = 6
		opts.calculateModuleSize()
		if opts.Qr.moduleSize != posqr.ModuleSize(6) {
			t.Errorf("moduleSize = %v, want 6", opts.Qr.moduleSize)
		}
	})
}

func TestQrOptions_CalculateFinalDimensions(t *testing.T) {
	opts := DefaultQROptions()
	opts.Qr.gridSize = 21
	opts.Qr.moduleSize = posqr.ModuleSize(4)
	// totalModules = 21 + 8 = 29; totalWidth = 29 * 4 = 116

	t.Run("scaled_up mode", func(t *testing.T) {
		opts.Qr.requestedWidth = 100
		opts.calculateFinalDimensions()
		if opts.Qr.scaleMode != "scaled_up" {
			t.Errorf("scaleMode = %q, want scaled_up", opts.Qr.scaleMode)
		}
		if opts.Qr.diff != 16 {
			t.Errorf("diff = %d, want 16", opts.Qr.diff)
		}
	})

	t.Run("scaled_down mode", func(t *testing.T) {
		opts.Qr.requestedWidth = 120
		opts.calculateFinalDimensions()
		if opts.Qr.scaleMode != "scaled_down" {
			t.Errorf("scaleMode = %q, want scaled_down", opts.Qr.scaleMode)
		}
		if opts.Qr.diff != -4 {
			t.Errorf("diff = %d, want -4", opts.Qr.diff)
		}
	})

	t.Run("exact mode", func(t *testing.T) {
		opts.Qr.requestedWidth = 116
		opts.calculateFinalDimensions()
		// diff is 0, scaleMode remains unchanged or exact
		if opts.Qr.totalWidth != 116 {
			t.Errorf("totalWidth = %d, want 116", opts.Qr.totalWidth)
		}
	})
}

func TestQrOptions_LoadLogo(t *testing.T) {
	t.Run("empty LogoData", func(t *testing.T) {
		opts := DefaultQROptions()
		opts.LogoData = ""
		opts.loadLogo()
		if opts.Logo.Image != nil {
			t.Error("expected Logo.Image to be nil")
		}
	})

	t.Run("corrupt LogoData resets to empty", func(t *testing.T) {
		opts := DefaultQROptions()
		opts.LogoData = "not-a-valid-base64"
		opts.loadLogo()
		if opts.LogoData != "" {
			t.Errorf("LogoData = %q, want empty after error", opts.LogoData)
		}
		if opts.Logo.Image != nil {
			t.Error("expected Logo.Image to be nil")
		}
	})

	t.Run("valid base64 png loads successfully", func(t *testing.T) {
		_, b64 := createTestPNG(t, 24, 24)
		opts := DefaultQROptions()
		opts.LogoData = b64
		opts.loadLogo()
		if opts.Logo.Image == nil {
			t.Fatal("expected Logo.Image to be non-nil")
		}
		if opts.Logo.format != "png" {
			t.Errorf("format = %q, want png", opts.Logo.format)
		}
		if opts.Logo.width != 24 || opts.Logo.height != 24 {
			t.Errorf("dimensions = %dx%d, want 24x24", opts.Logo.width, opts.Logo.height)
		}
	})
}

func TestMapLogoSize(t *testing.T) {
	tests := []struct {
		ecl  posqr.ErrorCorrection
		want int
	}{
		{posqr.LevelL, 7},
		{posqr.LevelM, 6},
		{posqr.LevelQ, 5},
		{posqr.LevelH, 4},
		{posqr.ErrorCorrection(99), 5},
	}

	for _, tt := range tests {
		if got := mapLogoSize(tt.ecl); got != tt.want {
			t.Errorf("mapLogoSize(%v) = %d, want %d", tt.ecl, got, tt.want)
		}
	}
}

func TestMapEclOption(t *testing.T) {
	levels := []posqr.ErrorCorrection{
		posqr.LevelL,
		posqr.LevelM,
		posqr.LevelQ,
		posqr.LevelH,
		posqr.ErrorCorrection(100), // default
	}

	for _, lvl := range levels {
		opt := mapEclOption(lvl)
		if opt == nil {
			t.Errorf("mapEclOption(%v) returned nil", lvl)
		}
	}
}

func TestBuildImageOptions(t *testing.T) {
	t.Run("without logo, without circle", func(t *testing.T) {
		opts := DefaultQROptions()
		opts.Qr.moduleSize = posqr.ModuleSize(4)
		imgOpts := buildImageOptions(opts)
		if len(imgOpts) != 2 {
			t.Errorf("len(imgOpts) = %d, want 2", len(imgOpts))
		}
	})

	t.Run("with circle shape", func(t *testing.T) {
		opts := DefaultQROptions()
		opts.Qr.moduleSize = posqr.ModuleSize(4)
		opts.CircleShape = true
		imgOpts := buildImageOptions(opts)
		if len(imgOpts) != 3 {
			t.Errorf("len(imgOpts) = %d, want 3", len(imgOpts))
		}
	})

	t.Run("with PNG logo and safe zone", func(t *testing.T) {
		img, b64 := createTestPNG(t, 20, 20)
		opts := DefaultQROptions()
		opts.Qr.moduleSize = posqr.ModuleSize(4)
		opts.Qr.dataWidth = 84
		opts.LogoData = b64
		opts.Logo.Image = img
		opts.Logo.format = "png"
		opts.Logo.sizeMulti = 2
		opts.ErrorCorrection = posqr.LevelM

		imgOpts := buildImageOptions(opts)
		if len(imgOpts) == 0 {
			t.Error("expected non-empty image options")
		}
	})

	t.Run("with sizeMulti 0 falls back to no logo", func(t *testing.T) {
		img, b64 := createTestPNG(t, 20, 20)
		opts := DefaultQROptions()
		opts.Qr.moduleSize = posqr.ModuleSize(4)
		opts.Qr.dataWidth = 84
		opts.LogoData = b64
		opts.Logo.Image = img
		opts.Logo.format = "png"
		opts.Logo.sizeMulti = 0

		imgOpts := buildImageOptions(opts)
		if opts.LogoData != "" || opts.Logo.Image != nil {
			t.Error("expected LogoData and Logo.Image to be cleared when sizeMulti is 0")
		}
		if len(imgOpts) != 2 {
			t.Errorf("len(imgOpts) = %d, want 2 (only qr width and border)", len(imgOpts))
		}
	})

	t.Run("with JPEG logo and max sizeMulti clamping", func(t *testing.T) {
		img, b64 := createTestJPEG(t, 30, 30)
		opts := DefaultQROptions()
		opts.Qr.moduleSize = posqr.ModuleSize(4)
		opts.Qr.dataWidth = 100
		opts.LogoData = b64
		opts.Logo.Image = img
		opts.Logo.format = "jpeg"
		opts.Logo.sizeMulti = 10 // above maxSizeMulti -> clamped to 5
		opts.ErrorCorrection = posqr.LevelH

		imgOpts := buildImageOptions(opts)
		if opts.Logo.sizeMulti != maxSizeMulti {
			t.Errorf("sizeMulti = %d, want %d", opts.Logo.sizeMulti, maxSizeMulti)
		}
		if len(imgOpts) == 0 {
			t.Error("expected non-empty image options")
		}
	})
}

func TestProcessQRImage(t *testing.T) {
	t.Run("empty data error", func(t *testing.T) {
		_, err := ProcessQRImage("", nil)
		if err == nil {
			t.Error("expected error on empty data, got nil")
		}
	})

	t.Run("nil options uses defaults", func(t *testing.T) {
		img, err := ProcessQRImage("https://example.com/test", nil)
		if err != nil {
			t.Fatalf("ProcessQRImage failed with nil options: %v", err)
		}
		if img == nil {
			t.Fatal("expected non-nil image")
		}
		if img.Bounds().Dx() == 0 || img.Bounds().Dy() == 0 {
			t.Errorf("empty image dimensions: %v", img.Bounds())
		}
	})

	t.Run("with CircleShape", func(t *testing.T) {
		opts := DefaultQROptions()
		opts.CircleShape = true
		img, err := ProcessQRImage("https://example.com/circle", opts)
		if err != nil {
			t.Fatalf("ProcessQRImage with CircleShape failed: %v", err)
		}
		if img == nil {
			t.Fatal("expected non-nil image")
		}
	})

	t.Run("with PNG logo", func(t *testing.T) {
		_, b64 := createTestPNG(t, 16, 16)
		opts := DefaultQROptions()
		opts.LogoData = b64
		opts.ErrorCorrection = posqr.LevelQ
		img, err := ProcessQRImage("https://example.com/with-logo", opts)
		if err != nil {
			t.Fatalf("ProcessQRImage with logo failed: %v", err)
		}
		if img == nil {
			t.Fatal("expected non-nil image")
		}
	})

	t.Run("with all error correction levels", func(t *testing.T) {
		levels := []posqr.ErrorCorrection{
			posqr.LevelL,
			posqr.LevelM,
			posqr.LevelQ,
			posqr.LevelH,
		}
		for _, lvl := range levels {
			opts := DefaultQROptions()
			opts.ErrorCorrection = lvl
			img, err := ProcessQRImage("https://example.com/ecl", opts)
			if err != nil {
				t.Fatalf("ProcessQRImage failed for ECL %v: %v", lvl, err)
			}
			if img == nil {
				t.Fatalf("expected non-nil image for ECL %v", lvl)
			}
		}
	})
}

func TestScaleImageToWidth(t *testing.T) {
	t.Run("nil image returns nil", func(t *testing.T) {
		if got := ScaleImageToWidth(nil, 100); got != nil {
			t.Errorf("ScaleImageToWidth(nil, 100) = %v, want nil", got)
		}
	})

	t.Run("targetWidth <= 0 returns original image", func(t *testing.T) {
		img, _ := createTestPNG(t, 20, 10)
		if got := ScaleImageToWidth(img, 0); got != img {
			t.Error("expected original image when targetWidth <= 0")
		}
		if got := ScaleImageToWidth(img, -10); got != img {
			t.Error("expected original image when targetWidth < 0")
		}
	})

	t.Run("targetWidth == originalWidth returns original image", func(t *testing.T) {
		img, _ := createTestPNG(t, 20, 10)
		if got := ScaleImageToWidth(img, 20); got != img {
			t.Error("expected original image when targetWidth == originalWidth")
		}
	})

	t.Run("scales down proportionally", func(t *testing.T) {
		img, _ := createTestPNG(t, 40, 20)
		scaled := ScaleImageToWidth(img, 20)
		if scaled == nil {
			t.Fatal("expected non-nil scaled image")
		}
		if scaled.Bounds().Dx() != 20 {
			t.Errorf("scaled width = %d, want 20", scaled.Bounds().Dx())
		}
		if scaled.Bounds().Dy() != 10 {
			t.Errorf("scaled height = %d, want 10", scaled.Bounds().Dy())
		}
	})

	t.Run("scales up proportionally", func(t *testing.T) {
		img, _ := createTestPNG(t, 20, 10)
		scaled := ScaleImageToWidth(img, 40)
		if scaled == nil {
			t.Fatal("expected non-nil scaled image")
		}
		if scaled.Bounds().Dx() != 40 {
			t.Errorf("scaled width = %d, want 40", scaled.Bounds().Dx())
		}
		if scaled.Bounds().Dy() != 20 {
			t.Errorf("scaled height = %d, want 20", scaled.Bounds().Dy())
		}
	})
}
