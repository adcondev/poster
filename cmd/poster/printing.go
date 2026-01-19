//go:build windows

package main

import (
	"errors"
	"fmt"
	"log"
	"runtime"
	"strings"

	"github.com/adcondev/poster/pkg/connection"
	"github.com/adcondev/poster/pkg/document/schema"
	"github.com/adcondev/poster/pkg/profile"
)

func detectPrinter() string {
	if runtime.GOOS != win {
		return ""
	}

	// Use the new ListAvailablePrinters function
	printers, err := connection.ListAvailablePrinters()
	if err != nil {
		log.Printf("Warning: failed to enumerate printers: %v", err)
		return ""
	}

	// First, try to find a thermal printer
	thermalPrinters := connection.FilterThermalPrinters(printers)
	if len(thermalPrinters) > 0 {
		// Prefer default thermal printer if available
		for _, p := range thermalPrinters {
			if p.IsDefault {
				log.Printf("Auto-detected default thermal printer: %s", p.Name)
				return p.Name
			}
		}
		// Otherwise use first thermal printer
		log.Printf("Auto-detected thermal printer: %s", thermalPrinters[0].Name)
		return thermalPrinters[0].Name
	}

	// Fallback: search for common printer name patterns
	commonPatterns := []string{
		"pos-80", "pos-58", "80mm", "58mm",
		"pt-210", "gp-58", "ec-pm",
		"receipt", "thermal", "epson",
	}

	physicalPrinters := connection.FilterPhysicalPrinters(printers)
	for _, p := range physicalPrinters {
		nameLower := strings.ToLower(p.Name)
		for _, pattern := range commonPatterns {
			if strings.Contains(nameLower, pattern) {
				log.Printf("Auto-detected printer by name pattern: %s", p.Name)
				return p.Name
			}
		}
	}

	// Last resort: use default physical printer
	for _, p := range physicalPrinters {
		if p.IsDefault {
			log.Printf("Using default physical printer: %s", p.Name)
			return p.Name
		}
	}

	return ""
}

func createConnection(config *Config) (connection.Connector, error) {
	switch strings.ToLower(config.ConnectionType) {
	case win:
		if runtime.GOOS != win {
			return nil, errors.New("windows connection type requires Windows OS")
		}
		if config.PrinterName == "" {
			return nil, fmt.Errorf("printer name required for Windows connection")
		}
		return connection.NewWindowsPrintConnector(config.PrinterName)

	case "network":
		if config.NetworkAddr == "" {
			return nil, fmt.Errorf("network address required")
		}
		// TODO: Implement when available in library
		return nil, fmt.Errorf("network connection not yet implemented")

	case "serial":
		if config.SerialPort == "" {
			return nil, fmt.Errorf("serial port required")
		}
		// TODO: Implement when available in library
		return nil, fmt.Errorf("serial connection not yet implemented")

	case "file":
		// TODO: Implement file output connector
		return nil, fmt.Errorf("file output not yet implemented")

	default:
		return nil, fmt.Errorf("unknown connection type: %s", config.ConnectionType)
	}
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
