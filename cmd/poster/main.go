// Package main es el ejecutor principal de la librería pos-printer
// Uso: pos-printer [options] <json_file> [printer_name]
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/adcondev/pos-printer/internal/load"
	"github.com/adcondev/pos-printer/pkg/composer"
	"github.com/adcondev/pos-printer/pkg/connection"
	"github.com/adcondev/pos-printer/pkg/document"
	"github.com/adcondev/pos-printer/pkg/profile"
	"github.com/adcondev/pos-printer/pkg/service"
)

const (
	AppName    = "pos-printer"
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
  - JSON files should follow the pos-printer document format
  - Use --dry-run to validate JSON without printing`)
}

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

	// Parse document
	var doc document.Document
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

	// Create executor
	executor := document.NewExecutor(printerService)

	// Execute document
	if config.Debug {
		log.Println("Executing document...")
	}

	return executor.Execute(&doc)
}

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

func createProfile(doc *document.Document) *profile.Escpos {
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

func validateDocument(doc *document.Document) error {
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
