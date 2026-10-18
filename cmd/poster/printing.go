//go:build windows

package main

import (
	"errors"
	"fmt"
	"log"
	"runtime"
	"strings"

	"github.com/adcondev/poster/pkg/connection"
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
		printerName := config.PrinterName
		if printerName == "" {
			printerName = detectPrinter()
		}
		if printerName == "" {
			return nil, fmt.Errorf("printer name required for Windows connection")
		}
		return connection.NewWindowsPrintConnector(printerName)

	case "network", "tcp":
		target := config.PrinterName
		if target == "" {
			target = config.NetworkAddr
		}
		if target == "" {
			return nil, fmt.Errorf("printer name or address required for network/tcp connection")
		}
		return connection.NewTCPConnector(target)

	case "serial":
		if config.SerialPort == "" {
			return nil, fmt.Errorf("serial port required")
		}
		// TODO: Implement when available in library
		return nil, fmt.Errorf("serial connection not yet implemented")

	case "file":
		targetPath := config.OutputFile
		if targetPath == "" {
			targetPath = config.PrinterName
		}
		if targetPath == "" {
			return nil, fmt.Errorf("output file path required for file connection (use -output or -p)")
		}
		return connection.NewFileConnector(targetPath)

	default:
		return nil, fmt.Errorf("unknown connection type: %s", config.ConnectionType)
	}
}
