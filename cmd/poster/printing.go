package main

import (
	"fmt"
	"log"
	"runtime"
	"strings"

	"github.com/adcondev/poster/pkg/connection"
	"github.com/adcondev/poster/pkg/document/schema"
	"github.com/adcondev/poster/pkg/profile"
)

func detectPrinter() string {
	// Common printer names to search for
	commonNames := []string{
		"POS-80",
		"POS-58",
		"80mm EC-PM-80250",
		"PT-210",
		"GP-58N",
		"Generic",
		"Receipt",
		"Thermal",
		"EPSON",
	}

	if runtime.GOOS == win {
		// TODO: Implement Windows printer enumeration
		for _, name := range commonNames {
			// Check if printer exists
			log.Printf("Checking for printer: %s", name)
		}
	}

	return ""
}

func createConnection(config *Config) (connection.Connector, error) {
	switch strings.ToLower(config.ConnectionType) {
	case win:
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
