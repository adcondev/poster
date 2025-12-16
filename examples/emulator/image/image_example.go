// Package main demonstrates image embedding in the ESC/POS emulator
package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"path/filepath"

	"github.com/adcondev/poster/internal/load"
	"github.com/adcondev/poster/pkg/composer"
	"github.com/adcondev/poster/pkg/connection"
	"github.com/adcondev/poster/pkg/constants"
	"github.com/adcondev/poster/pkg/emulator"
	"github.com/adcondev/poster/pkg/graphics"
	"github.com/adcondev/poster/pkg/profile"
	"github.com/adcondev/poster/pkg/service"
)

func main() {
	// ============================================================================
	// Configuration
	// ============================================================================
	baseDir := "./examples/emulator/image"
	outputFile := "receipt_with_images.png"
	printerName := "58mm PT-210" // Adjust to your printer

	// Ensure output directory exists
	if err := os.MkdirAll(baseDir, 0755); err != nil { //nolint:gosec
		log.Fatalf("Failed to create directory: %v", err)
	}

	// ============================================================================
	// Part 1: Generate Receipt Image with Emulator
	// ============================================================================
	fmt.Println("=== Generating Receipt with Embedded Images ===")

	engine, err := emulator.New58mmEngine()
	if err != nil {
		log.Fatalf("Error creating engine: %v", err)
	}

	// Load test image
	logoImg := loadTestImage()

	// Build the receipt
	buildReceipt(engine, logoImg)

	// Save to file
	fullPath := filepath.Join(baseDir, outputFile)
	if err := saveReceipt(engine, fullPath); err != nil {
		log.Fatalf("Failed to save receipt: %v", err)
	}

	fmt.Printf("Receipt saved to:  %s\n", fullPath)

	// Print dimensions
	result := engine.RenderWithInfo()
	fmt.Printf("Image dimensions: %dx%d pixels\n", result.Width, result.Height)

	// ============================================================================
	// Part 2: Optional Physical Printing
	// ============================================================================
	fmt.Print("\nDo you want to print the receipt?  (y/N): ")
	var answer string
	if _, err := fmt.Scanln(&answer); err != nil {
		fmt.Println("Skipping print.  Goodbye!")
		return
	}

	if answer != "y" && answer != "Y" {
		fmt.Println("Skipping print. Goodbye!")
		return
	}

	if err := printReceipt(baseDir, outputFile, printerName); err != nil {
		log.Fatalf("Failed to print receipt: %v", err)
	}

	fmt.Println("Receipt printed successfully!")
}

// loadTestImage loads a test image from Base64 or returns nil
func loadTestImage() image.Image {
	base64Data := Base64ImageExample()

	if base64Data == "BASE_64_CODE_HERE" || base64Data == "" {
		fmt.Println("Note: No Base64 image provided. Using placeholder text instead.")
		return nil
	}

	data, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		log.Printf("Warning: Failed to decode base64 image: %v", err)
		return nil
	}

	img, format, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		log.Printf("Warning: Failed to decode image:  %v", err)
		return nil
	}

	fmt.Printf("Loaded image format: %s (%dx%d)\n", format, img.Bounds().Dx(), img.Bounds().Dy())
	return img
}

// buildReceipt constructs the receipt content
func buildReceipt(engine *emulator.Engine, logoImg image.Image) {
	engine.SetFont("A")
	engine.AlignCenter()

	// ============================================================================
	// Header with Logo (Normal Preview - keeps colors/grayscale)
	// ============================================================================
	if logoImg != nil {
		fmt.Println("Embedding logo (normal preview)...")
		opts := emulator.DefaultImageOptions()
		opts.PixelWidth = 256
		opts.Align = constants.Center.String()
		if err := engine.PrintImageWithOptions(logoImg, opts); err != nil {
			log.Printf("Error printing logo: %v", err)
		}
		engine.Feed(1) // Extra space for larger text
	}

	// Store name with larger text
	engine.SetBold(true)
	engine.SetSize(2, 2)
	engine.PrintLine("RED2000")
	engine.SetSize(1, 1)
	engine.SetBold(false)
	engine.PrintLine("Stationery & Tech")
	engine.PrintLine("www.red2000.com")
	engine.Separator("=", 32)

	// ============================================================================
	// Order Details
	// ============================================================================
	engine.AlignLeft()
	engine.PrintLine("Date:  2024-12-16 14:30")
	engine.PrintLine("Order #: IMG-00456")
	engine.Separator("-", 32)

	engine.PrintLine("Bookmarks           $4.50")
	engine.PrintLine("Pencils             $3.00")
	engine.PrintLine("Nail polish         $8.99")
	engine.Separator("-", 32)

	engine.SetBold(true)
	engine.PrintLine("TOTAL              $16.49")
	engine.SetBold(false)
	engine.Separator("=", 32)

	// ============================================================================
	// QR Code Section (Thermal Preview - shows exact print output)
	// ============================================================================
	engine.AlignCenter()
	engine.PrintLine("Pay with QR Code:")
	engine.Feed(1)

	if logoImg != nil {
		// Use thermal preview to show exactly how image will print
		fmt.Println("Embedding QR placeholder (thermal preview)...")
		if err := engine.PrintImageThermalPreview(logoImg, 256); err != nil {
			log.Printf("Error printing QR preview: %v", err)
		}
	} else {
		engine.PrintLine("[QR Code Here]")
	}

	engine.Feed(1)
	engine.PrintLine("Scan to pay or leave review")
	engine.PrintLine("á é í ó ú ü ñ Á É Í Ó Ú Ü Ñ") // Test special characters
	engine.PrintLine("¿? ¡! @ # & * () [] {} < >")  // Test punctuation

	// ============================================================================
	// Footer
	// ============================================================================
	engine.Separator("-", 32)
	engine.PrintLine("Thank you for visiting!")
	engine.PrintLine("Follow us @red2000Stationery")
	engine.Feed(1)
	engine.Cut(true)
}

// saveReceipt writes the receipt image to a file
func saveReceipt(engine *emulator.Engine, path string) error {
	file, err := os.Create(path) //nolint:gosec
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Panicf("failed to close file: %v", err)
		}
	}(file)

	return engine.WritePNG(file)
}

// printReceipt sends the generated image to the physical printer
func printReceipt(baseDir, fileName, printerName string) error {
	fmt.Println("\n=== Printing to Physical Printer ===")

	// ============================================================================
	// Setup Printer
	// ============================================================================
	prof := profile.CreatePt210()
	fmt.Printf("Printer profile: %s (%d dots/line)\n", prof.Model, prof.DotsPerLine)

	proto := composer.NewEscpos()

	conn, err := connection.NewWindowsPrintConnector(printerName)
	if err != nil {
		return fmt.Errorf("connect to printer '%s': %w", printerName, err)
	}
	defer func(conn *connection.WindowsPrintConnector) {
		err := conn.Close()
		if err != nil {
			log.Panicf("failed to close printer connection: %v", err)
		}
	}(conn)

	printer, err := service.NewPrinter(proto, prof, conn)
	if err != nil {
		return fmt.Errorf("create printer service: %w", err)
	}
	defer func(printer *service.Printer) {
		err := printer.Close()
		if err != nil {
			log.Panicf("failed to close printer service: %v", err)
		}
	}(printer)

	if err := printer.Initialize(); err != nil {
		return fmt.Errorf("initialize printer: %w", err)
	}
	fmt.Println("Printer initialized")

	// ============================================================================
	// Load and Process Image
	// ============================================================================
	img, err := load.ImgFromFile(baseDir, fileName)
	if err != nil {
		return fmt.Errorf("load image:  %w", err)
	}
	fmt.Printf("Loaded receipt image: %dx%d pixels\n", img.Bounds().Dx(), img.Bounds().Dy())

	// Process for thermal printing
	imgOpts := &graphics.ImgOptions{
		PixelWidth:     prof.DotsPerLine,
		Threshold:      128,
		Dithering:      graphics.Atkinson, // Best for receipts with images
		Scaling:        graphics.BiLinear,
		PreserveAspect: true,
	}

	pipeline := graphics.NewPipeline(imgOpts)
	bitmap, err := pipeline.Process(img)
	if err != nil {
		return fmt.Errorf("process image: %w", err)
	}
	fmt.Printf("Processed bitmap: %dx%d pixels (%d bytes)\n",
		bitmap.Width, bitmap.Height, len(bitmap.GetRasterData()))

	// ============================================================================
	// Print
	// ============================================================================
	if err := printer.AlignCenter(); err != nil {
		return fmt.Errorf("set alignment: %w", err)
	}

	if err := printer.PrintBitmap(bitmap); err != nil {
		return fmt.Errorf("print bitmap: %w", err)
	}
	fmt.Println("Bitmap sent to printer")

	if err := printer.FeedLines(2); err != nil {
		return fmt.Errorf("feed lines: %w", err)
	}

	if err := printer.PartialFeedAndCut(2); err != nil {
		log.Printf("Warning: Cut failed (printer may not support it): %v", err)
	}

	return nil
}
