package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/adcondev/pos-printer/internal/load"
	"github.com/adcondev/pos-printer/pkg/document/schema"
)

func loadJSON(filePath string) ([]byte, error) {
	// Validate path
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, fmt.Errorf("invalid path: %w", err)
	}

	// Security check
	baseDir := filepath.Dir(absPath)
	fileName := filepath.Base(absPath)
	securePath, err := load.SecureFilepath(baseDir, fileName)
	if err != nil {
		return nil, fmt.Errorf("security check failed: %w", err)
	}

	// Read file
	data, err := os.ReadFile(securePath) //nolint:gosec
	if err != nil {
		return nil, err
	}

	// Validate JSON
	var temp interface{}
	if err := json.Unmarshal(data, &temp); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	return data, nil
}

func validateDocument(doc *schema.Document) error {
	fmt.Println("=== Document Validation ===")
	fmt.Printf("Version: %s\n", doc.Version)
	fmt.Printf("Profile: %s (%dmm)\n", doc.Profile.Model, doc.Profile.PaperWidth)
	fmt.Printf("Commands: %d\n", len(doc.Commands))

	// Validate commands
	commandCounts := make(map[string]int)
	for i, cmd := range doc.Commands {
		if cmd.Type == "" {
			return fmt.Errorf("command %d: missing type", i)
		}
		commandCounts[cmd.Type]++

		// Basic data validation
		if len(cmd.Data) == 0 {
			return fmt.Errorf("command %d (%s): missing data", i, cmd.Type)
		}
	}

	// Show command summary
	fmt.Println("\nCommand Summary:")
	for cmdType, count := range commandCounts {
		fmt.Printf("  %s: %d\n", cmdType, count)
	}

	fmt.Println("\n✅ Validation passed")
	return nil
}

// getPrinterNameFromDocument extrae el nombre de la impresora del documento JSON
// Prioridad: 1) --printer flag, 2) profile.model del JSON, 3) auto-detect
func getPrinterNameFromDocument(config *Config, jsonData []byte) (string, error) {
	// Si ya se especificó por flag, usar ese
	if config.PrinterName != "" {
		return config.PrinterName, nil
	}

	// Intentar extraer del JSON
	var doc struct {
		Profile struct {
			Model string `json:"model"`
		} `json:"profile"`
	}

	if err := json.Unmarshal(jsonData, &doc); err != nil {
		return "", fmt.Errorf("failed to parse JSON for printer name: %w", err)
	}

	if doc.Profile.Model != "" {
		log.Printf("Using printer name from JSON: %s", doc.Profile.Model)
		return doc.Profile.Model, nil
	}

	// Fallback a auto-detect
	detected := detectPrinter()
	if detected != "" {
		return detected, nil
	}

	return "", fmt.Errorf("no printer specified and none found in JSON profile. model")
}
