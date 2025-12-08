// Package main es el ejecutor principal de la librería poster
// Uso: poster [options] <json_file> [printer_name]
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"runtime"

	"github.com/adcondev/poster/pkg/composer"
	"github.com/adcondev/poster/pkg/connection"
	"github.com/adcondev/poster/pkg/document/executor"
	"github.com/adcondev/poster/pkg/document/schema"
	"github.com/adcondev/poster/pkg/service"
)

const (
	AppName    = "poster"
	AppVersion = "1.0.0"
	AppAuthor  = "adcondev"
	win        = "windows"
)

type Config struct {
	JSONFile       string
	PrinterName    string
	DryRun         bool
	Debug          bool
	Version        bool
	Help           bool
	ListPrinters   bool
	ConnectionType string // windows, network, serial, file
	NetworkAddr    string
	SerialPort     string
	BaudRate       int
	OutputFile     string
}

func main() {
	config := parseArgs()

	if config.Version {
		showVersion()
		return
	}

	if config.Help {
		showHelp()
		return
	}

	if config.ListPrinters && runtime.GOOS == win {
		listWindowsPrinters()
		return
	}

	if config.JSONFile == "" {
		log.Fatal("Error: JSON file is required. Use -h for help")
	}

	if err := executePrint(config); err != nil {
		log.Fatalf("Print failed: %v", err)
	}

	log.Println("✅ Print completed successfully!")
}

func parseArgs() *Config {
	config := &Config{}

	// Define flags
	flag.StringVar(&config.JSONFile, "file", "", "JSON document file path")
	flag.StringVar(&config.JSONFile, "f", "", "JSON document file path (short)")

	flag.StringVar(&config.PrinterName, "printer", "", "Printer name")
	flag.StringVar(&config.PrinterName, "p", "", "Printer name (short)")

	flag.StringVar(&config.ConnectionType, "type", win, "Connection type: windows, network, serial, file")
	flag.StringVar(&config.ConnectionType, "t", win, "Connection type (short)")

	flag.StringVar(&config.NetworkAddr, "network", "", "Network address (e.g., 192.168.1.100:9100)")
	flag.StringVar(&config.SerialPort, "serial", "", "Serial port (e.g., COM1, /dev/ttyUSB0)")
	flag.IntVar(&config.BaudRate, "baud", 9600, "Serial baud rate")
	flag.StringVar(&config.OutputFile, "output", "output.prn", "Output file for file type")

	flag.BoolVar(&config.DryRun, "dry-run", false, "Validate without printing")
	flag.BoolVar(&config.Debug, "debug", false, "Enable debug logging")
	flag.BoolVar(&config.ListPrinters, "list", false, "List available printers (Windows only)")
	flag.BoolVar(&config.Version, "version", false, "Show version")
	flag.BoolVar(&config.Version, "v", false, "Show version (short)")
	flag.BoolVar(&config.Help, "help", false, "Show help")
	flag.BoolVar(&config.Help, "h", false, "Show help (short)")

	flag.Parse()

	// Handle positional arguments
	args := flag.Args()
	if len(args) > 0 && config.JSONFile == "" {
		config.JSONFile = args[0]
	}
	if len(args) > 1 && config.PrinterName == "" {
		config.PrinterName = args[1]
	}

	// Auto-detect printer name for common models
	if config.PrinterName == "" && config.ConnectionType == win {
		config.PrinterName = detectPrinter()
	}

	return config
}

func showVersion() {
	fmt.Printf("%s v%s\n", AppName, AppVersion)
	fmt.Printf("Author: %s\n", AppAuthor)
	fmt.Printf("Go: %s\n", runtime.Version())
	fmt.Printf("OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
}

func showHelp() {
	fmt.Printf(`%s v%s - ESC/POS Printer Driver

USAGE:
  %s [options] <json_file> [printer_name]

EXAMPLES:
  %s ticket.json
  %s ticket.json "POS-80"
  %s -t network -network 192.168.1.100:9100 ticket.json
  %s -t serial -serial COM1 -baud 115200 ticket.json
  %s -t file -output receipt.prn ticket.json
  %s --dry-run ticket.json

OPTIONS:
`, AppName, AppVersion, AppName, AppName, AppName, AppName, AppName, AppName, AppName)

	flag.PrintDefaults()

	fmt.Println(`
CONNECTION TYPES:
  windows  - Windows printer (default)
  network  - Network printer
  serial   - Serial/USB printer  
  file     - Output to file

NOTES:
  - If no printer is specified, attempts to auto-detect common models
  - JSON files should follow the poster document format
  - Use --dry-run to validate JSON without printing`)
}

func listWindowsPrinters() {
	fmt.Println("Available printers:")
	// TODO: Implement Windows printer listing using winspool.dll
	fmt.Println("  (Feature not yet implemented)")
}

func executePrint(config *Config) error {
	// Load and validate JSON
	jsonData, err := loadJSON(config.JSONFile)
	if err != nil {
		return fmt.Errorf("failed to load JSON: %w", err)
	}

	// Obtener nombre de impresora del JSON si no se especificó
	printerName, err := getPrinterNameFromDocument(config, jsonData)
	if err != nil && config.ConnectionType == win {
		return err
	}
	config.PrinterName = printerName

	// Parse document
	var doc schema.Document
	if err := json.Unmarshal(jsonData, &doc); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Debug info
	if config.Debug {
		log.Printf("Document version: %s", doc.Version)
		log.Printf("Profile: %s (%dmm)", doc.Profile.Model, doc.Profile.PaperWidth)
		log.Printf("Commands: %d", len(doc.Commands))
		log.Printf("Connection type: %s", config.ConnectionType)
	}

	// Dry run validation
	if config.DryRun {
		return validateDocument(&doc)
	}

	// Create connection
	conn, err := createConnection(config)
	if err != nil {
		return fmt.Errorf("failed to create connection: %w", err)
	}
	defer func(conn connection.Connector) {
		err := conn.Close()
		if err != nil {
			log.Panicf("failed to close connection: %v", err)
		}
	}(conn)

	// Create profile
	prof := createProfile(&doc)

	// Create protocol
	proto := composer.NewEscpos()

	// Create printer
	printerService, err := service.NewPrinter(proto, prof, conn)
	if err != nil {
		return fmt.Errorf("failed to create printer: %w", err)
	}
	defer func(printerService *service.Printer) {
		err := printerService.Close()
		if err != nil {
			log.Panicf("failed to close printer: %v", err)
		}
	}(printerService)

	// Create exec
	exec := executor.NewExecutor(printerService)

	// Execute document
	if config.Debug {
		log.Println("Executing document...")
	}

	return exec.Execute(&doc)
}
