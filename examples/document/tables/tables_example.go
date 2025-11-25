// Package main is an example of printing a document with tables using the pos-printer library in JSON format.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/adcondev/pos-printer/internal/load"
	"github.com/adcondev/pos-printer/pkg/composer"
	"github.com/adcondev/pos-printer/pkg/connection"
	"github.com/adcondev/pos-printer/pkg/document"
	"github.com/adcondev/pos-printer/pkg/profile"
	"github.com/adcondev/pos-printer/pkg/service"
)

func main() {
	// Check command line arguments
	if len(os.Args) < 4 {
		fmt.Println("Usage: table_example <base_path> <rel_path> <printer_name>")
		fmt.Println("Example: table_example ./path_to_file ./file.json 'POS-80'")
		os.Exit(1)
	}

	baseDir := os.Args[1]
	relPath := os.Args[2]
	printerName := os.Args[3]

	// Read JSON file
	// gosec: Path traversal check
	jsonFile, err := load.SecureFilepath(baseDir, relPath)
	if err != nil {
		log.Fatalf("Invalid JSON file path: %v", err)
	}
	jsonData, err := os.ReadFile(jsonFile) //nolint:gosec
	if err != nil {
		log.Fatalf("Failed to read JSON file: %v", err)
	}

	// Parse JSON document
	var doc document.Document
	if err := json.Unmarshal(jsonData, &doc); err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	// Validate document has at least one table command
	hasTable := false
	for _, cmd := range doc.Commands {
		if cmd.Type == "table" {
			hasTable = true
			break
		}
	}

	if !hasTable {
		log.Println("Warning: No table commands found in document")
	}

	// Create printer connection
	conn, err := connection.NewWindowsPrintConnector(printerName)
	if err != nil {
		log.Fatalf("Failed to connect to printer: %v", err)
	}
	defer func(conn *connection.WindowsPrintConnector) {
		err := conn.Close()
		if err != nil {
			log.Panicf("Failed to close printer connection: %v", err)
		}
	}(conn)

	// Create ESC/POS protocol
	protocol := composer.NewEscpos()

	// Create printer profile based on document configuration
	var printerProfile *profile.Escpos
	switch doc.Profile.PaperWidth {
	case 58:
		printerProfile = profile.CreateProfile58mm()
	case 80:
		printerProfile = profile.CreateProfile80mm()
	default:
		// Default to 80mm if not specified
		printerProfile = profile.CreateProfile80mm()
	}

	// Update profile from document if specified
	if doc.Profile.Model != "" {
		printerProfile.Model = doc.Profile.Model
	}

	// Create printer service
	printer, err := service.NewPrinter(protocol, printerProfile, conn)
	if err != nil {
		log.Panicf("Failed to create printer service: %v", err)
	}

	// Create document executor
	executor := document.NewExecutor(printer)

	// Execute the document
	fmt.Printf("Printing table example to %s...\n", printerName)
	if err := executor.Execute(&doc); err != nil {
		log.Panicf("Failed to print document: %v", err)
	}

	fmt.Println("Table printing completed successfully!")
}
