// Package main in an example demonstrating the ESC/POS emulator functionality
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/adcondev/poster/internal/load"

	"github.com/adcondev/poster/pkg/composer"
	"github.com/adcondev/poster/pkg/connection"
	"github.com/adcondev/poster/pkg/emulator"
	"github.com/adcondev/poster/pkg/graphics"
	"github.com/adcondev/poster/pkg/profile"
	"github.com/adcondev/poster/pkg/service"
)

func main() {
	// ============================================================================
	// Part 1: Generate Receipt Image with Emulator
	// ============================================================================

	// Create a new emulator engine for 58mm paper
	engine, err := emulator.New58mmEngine()
	if err != nil {
		fmt.Printf("Error creating engine: %v\n", err)
		os.Exit(1)
	}

	// ============================================================================
	// Header
	// ============================================================================
	engine.SetFont("A")
	engine.AlignCenter()
	engine.SetBold(true)
	engine.Feed(1)
	engine.SetSize(2, 2) // Double size
	engine.Print("RED2000")
	engine.SetSize(1, 1) // Normal size
	engine.Feed(1)
	engine.PrintLine("ÁÉÍÓÚÜÑ áéíóúüñ") // Test special characters
	engine.SetBold(false)
	engine.PrintLine("RED 2000 Coffee Shop")
	engine.PrintLine("www.red2000.com")
	engine.PrintLine("Tel: (123) 456-7890")
	engine.PrintLine("Date: 2024-06-01 14:30")
	engine.PrintLine("Receipt #: 000123")
	engine.PrintLine("123 Main Street")
	engine.AlignLeft()

	// Separator
	engine.Separator("=", 32)

	// ============================================================================
	// Items
	// ============================================================================
	engine.PrintLine("Coffee              $3.50")
	engine.PrintLine("Sandwich            $8.00")
	engine.PrintLine("Cookie              $2.50")

	// Separator
	engine.Separator("-", 32)

	// ============================================================================
	// Total
	// ============================================================================
	engine.SetBold(true)
	engine.PrintLine("TOTAL              $14.00")
	engine.SetBold(false)

	engine.Separator("=", 32)

	// ============================================================================
	// Footer
	// ============================================================================
	engine.AlignCenter()
	engine.PrintLine("Gracias por su compra!")
	engine.SetSize(1, 1)
	engine.PrintLine("Visitenos en www.example.com")
	engine.AlignLeft()

	// cut
	engine.Cut(true) // Partial cut

	// ============================================================================
	// Save to file
	// ============================================================================
	baseDir := "./examples/emulator/basic"
	relPath := "receipt_emulated.png"

	file, err := os.Create(filepath.Join(baseDir, relPath))
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		os.Exit(1)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Panicf("failed to close file: %v\n", err)
		}
	}(file)

	if err := engine.WritePNG(file); err != nil {
		log.Panicf("failed to write PNG: %v\n", err)
	}

	fmt.Println("Receipt saved to receipt_emulated.png")

	// Print some info
	result := engine.RenderWithInfo()
	fmt.Printf("Image dimensions: %dx%d pixels\n", result.Width, result.Height)

	// ============================================================================
	// Part 2: Print the Generated Image to PT-210 Printer
	// ============================================================================

	// Ask user if they want to print
	fmt.Print("\nDo you want to print the receipt?  (y/N): ")
	var answer string
	_, err = fmt.Scanln(&answer)
	if err != nil {
		return
	}
	if answer != "y" && answer != "Y" {
		fmt.Println("Skipping print.  Goodbye!")
		return
	}

	// Print the generated image
	if err := printReceipt(baseDir, relPath); err != nil {
		log.Panicf("Failed to print receipt: %v\n", err)
	}

	fmt.Println("Receipt printed successfully!")
}

// printReceipt loads the generated image and prints it to the PT-210 printer
func printReceipt(baseDir, relPath string) error {
	// ============================================================================
	// Setup Printer Profile, Protocol, and Connection
	// ============================================================================

	// Create PT-210 profile (58mm thermal printer)
	prof := profile.CreatePt210()
	log.Printf("Using printer profile: %s", prof.Model)

	// Create ESC/POS protocol composer
	proto := composer.NewEscpos()

	// Create connection to printer
	// Note: On Windows, use the printer's shared name or USB name
	// On Linux/macOS, this would be a different connector (USB, Serial, Network)
	printerName := "58mm PT-210" // Adjust this to match your printer's name
	conn, err := connection.NewWindowsPrintConnector(printerName)
	if err != nil {
		return fmt.Errorf("failed to connect to printer '%s': %w", printerName, err)
	}
	defer func(conn *connection.WindowsPrintConnector) {
		err := conn.Close()
		if err != nil {
			log.Panicf("failed to close printer connection: %v", err)
		}
	}(conn)

	// Create printer service
	printer, err := service.NewPrinter(proto, prof, conn)
	if err != nil {
		return fmt.Errorf("failed to create printer service: %w", err)
	}

	// ============================================================================
	// Initialize Printer
	// ============================================================================

	if err := printer.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize printer: %w", err)
	}
	log.Println("Printer initialized")

	// ============================================================================
	// Load and Process Image
	// ============================================================================

	// Load the generated receipt image
	// Using load package for secure file loading

	img, err := load.ImgFromFile(baseDir, relPath)
	if err != nil {
		return fmt.Errorf("failed to load image: %w", err)
	}
	log.Printf("Loaded image: %dx%d pixels", img.Bounds().Dx(), img.Bounds().Dy())

	// Configure image processing options for 58mm printer
	imgOpts := &graphics.ImgOptions{
		PixelWidth:     prof.DotsPerLine, // 384 pixels for 58mm @ 203 DPI
		Threshold:      128,
		Dithering:      graphics.Threshold,
		Scaling:        graphics.NearestNeighbor,
		PreserveAspect: true,
	}

	// Process image through pipeline
	pipeline := graphics.NewPipeline(imgOpts)
	bitmap, err := pipeline.Process(img)
	if err != nil {
		return fmt.Errorf("failed to process image: %w", err)
	}
	log.Printf("Processed bitmap: %dx%d pixels (%d bytes)",
		bitmap.Width, bitmap.Height, len(bitmap.GetRasterData()))

	// ============================================================================
	// Print Image
	// ============================================================================

	// Center the image
	if err := printer.AlignCenter(); err != nil {
		return fmt.Errorf("failed to set alignment: %w", err)
	}

	// Print the bitmap
	if err := printer.PrintBitmap(bitmap); err != nil {
		return fmt.Errorf("failed to print bitmap: %w", err)
	}
	log.Println("Bitmap sent to printer")

	// Feed some lines after the image
	if err := printer.FeedLines(3); err != nil {
		return fmt.Errorf("failed to feed lines: %w", err)
	}

	// Partial cut
	if err := printer.PartialFeedAndCut(2); err != nil {
		log.Printf("Warning: cut command failed (printer may not support it): %v", err)
		// Not a fatal error - some printers don't have cutters
	}

	return nil
}
