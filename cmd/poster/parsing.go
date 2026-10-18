package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/adcondev/poster/internal/load"
	"github.com/adcondev/poster/pkg/constants"
	"github.com/adcondev/poster/pkg/document/schema"
	"github.com/adcondev/poster/pkg/profile"
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

	// Read file securely
	file, err := os.Open(securePath) //nolint:gosec
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()

	// Check file size
	fi, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}
	if fi.Size() > int64(constants.MaxJSONSize) {
		return nil, fmt.Errorf("file too large: %d bytes (max %d bytes)", fi.Size(), constants.MaxJSONSize)
	}

	// Read with limit reader to prevent TOCTOU issues if file grows
	reader := io.LimitReader(file, int64(constants.MaxJSONSize)+1)
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	if len(data) > constants.MaxJSONSize {
		return nil, fmt.Errorf("file too large (exceeded %d bytes during read)", constants.MaxJSONSize)
	}

	// Validate JSON
	var temp interface{}
	if err := json.Unmarshal(data, &temp); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	return data, nil
}

func validateJSON(filePath string) error {
	data, err := loadJSON(filePath)
	if err != nil {
		return err
	}

	var doc schema.Document
	if err := json.Unmarshal(data, &doc); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	if err := doc.Validate(); err != nil {
		return fmt.Errorf("invalid document: %w", err)
	}

	return validateDocument(&doc)
}

func validateDocument(doc *schema.Document) error {
	if err := doc.Validate(); err != nil {
		return fmt.Errorf("invalid document: %w", err)
	}
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
		fmt.Printf(" %s: %d\n", cmdType, count)
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

	return "", fmt.Errorf("no printer specified and none found in JSON profile.model")
}

func createProfile(doc *schema.Document) *profile.Escpos {
	var prof *profile.Escpos

	// Select profile based on paper width or model
	if doc.Profile.Model != "" {
		switch strings.ToLower(doc.Profile.Model) {
		case "80mm ec-pm-80250", "ec-pm-80250":
			prof = profile.CreateECPM80250()
		case "58mm pt-210", "pt-210":
			prof = profile.CreatePt210()
		case "58mm gp-58n", "gp-58n":
			prof = profile.CreateGP58N()
		default:
			if doc.Profile.PaperWidth >= 80 {
				prof = profile.CreateProfile80mm()
			} else {
				prof = profile.CreateProfile58mm()
			}
		}
	} else {
		// Default based on paper width
		if doc.Profile.PaperWidth >= 80 {
			prof = profile.CreateProfile80mm()
		} else {
			prof = profile.CreateProfile58mm()
		}
	}

	// Apply JSON overrides
	if doc.Profile.Model != "" {
		prof.Model = doc.Profile.Model
	}
	if doc.Profile.DPI > 0 {
		prof.DPI = doc.Profile.DPI
	}
	prof.HasQR = doc.Profile.HasQR

	return prof
}

