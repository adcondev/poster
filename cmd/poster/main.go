// Package main es el ejecutor principal de la librería poster
// Uso: poster [options] <json_file> [printer_name]
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/adcondev/poster/pkg/composer"
	"github.com/adcondev/poster/pkg/connection"
	"github.com/adcondev/poster/pkg/document/executor"
	"github.com/adcondev/poster/pkg/document/schema"
	"github.com/adcondev/poster/pkg/service"
)

const (
	AppName   = "poster"
	AppAuthor = "adcondev"
	win       = "windows"
)

// AppVersion is the CLI version string, overridden at build time via -ldflags "-X main.AppVersion=..."
var AppVersion = "dev"

type Config struct {
	JSONFile       string
	JSONOutput     bool
	PrinterName    string
	DryRun         bool
	Debug          bool
	Version        bool
	Help           bool
	ListPrinters   bool
	ListThermal    bool
	ListPhysical   bool
	ConnectionType string
	NetworkAddr    string
	SerialPort     string
	BaudRate       int
	OutputFile     string
}

func main() {
	os.Exit(run())
}

func run() int {
	config := parseArgs()

	if config.Version {
		showVersion()
		return 0
	}

	if config.Help {
		showHelp()
		return 0
	}

	if config.ListPrinters || config.ListThermal || config.ListPhysical {
		if runtime.GOOS != win {
			log.Printf("Printer listing is only available on Windows")
			return 1
		}
		if err := listPrinters(config); err != nil {
			log.Printf("Failed to enumerate printers: %v", err)
			return 1
		}
		return 0
	}

	if config.JSONFile == "" {
		log.Printf("Error: JSON file is required. Use -h for help")
		return 1
	}

	if err := executePrint(config); err != nil {
		log.Printf("Print failed: %v", err)
		return 1
	}

	log.Println("✅ Print completed successfully!")
	return 0
}

func parseArgs() *Config {
	config := &Config{}

	flag.StringVar(&config.JSONFile, "file", "", "JSON document file path")
	flag.StringVar(&config.JSONFile, "f", "", "JSON document file path (shorthand)")

	flag.BoolVar(&config.JSONOutput, "json", false, "Output in JSON format (for --list commands)")

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
	flag.BoolVar(&config.ListPrinters, "list", false, "List all available printers (Windows only)")
	flag.BoolVar(&config.ListThermal, "list-thermal", false, "List thermal/POS printers only (Windows only)")
	flag.BoolVar(&config.ListPhysical, "list-physical", false, "List physical printers only (Windows only)")
	flag.BoolVar(&config.Version, "version", false, "Show version")
	flag.BoolVar(&config.Version, "v", false, "Show version (short)")
	flag.BoolVar(&config.Help, "help", false, "Show help")
	flag.BoolVar(&config.Help, "h", false, "Show help (short)")

	flag.Parse()

	args := flag.Args()
	if len(args) > 0 && config.JSONFile == "" {
		config.JSONFile = args[0]
	}
	if len(args) > 1 && config.PrinterName == "" {
		config.PrinterName = args[1]
	}

	if !config.DryRun && config.PrinterName == "" && config.ConnectionType == win {
		config.PrinterName = detectPrinter()
	}

	return config
}

func showVersion() {
	fmt.Printf("%s %s\n", AppName, AppVersion)
	fmt.Printf("Author: %s\n", AppAuthor)
	fmt.Printf("Go: %s\n", runtime.Version())
	fmt.Printf("OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
}

func showHelp() {
	fmt.Printf(`%s %s - ESC/POS Printer Driver

USAGE:
  %s [options] <json_file> [printer_name]

EXAMPLES:
  %s ticket.json
  %s ticket.json "POS-80"
  %s -t file -output receipt.prn ticket.json
  %s --dry-run ticket.json
  %s --list
  %s --list-thermal
  %s --list-physical

OPTIONS:
`, AppName, AppVersion, AppName, AppName, AppName, AppName, AppName, AppName, AppName, AppName)

	flag.PrintDefaults()

	fmt.Println(`
CONNECTION TYPES:
  windows  - Windows Print Spooler (default)
  file     - Output raw ESC/POS bytes to file

PRINTER LISTING (Windows only):
  --list          List all installed printers
  --list-thermal  List only thermal/POS printers
  --list-physical List only physical (non-virtual) printers

NOTES:
  - If no printer is specified, attempts to auto-detect common models
  - JSON files should follow the poster v1.0 document format
  - Use --dry-run to validate JSON without printing`)
}

func listPrinters(config *Config) error {
	printers, err := connection.ListAvailablePrinters()
	if err != nil {
		if len(printers) == 0 {
			return err
		}
		log.Printf("Warning during printer enumeration: %v", err)
	}

	switch {
	case config.ListThermal:
		printers = connection.FilterThermalPrinters(printers)
		fmt.Println("Thermal/POS Printers:")
	case config.ListPhysical:
		printers = connection.FilterPhysicalPrinters(printers)
		fmt.Println("Physical Printers:")
	default:
		fmt.Println("All Available Printers:")
	}

	if len(printers) == 0 {
		fmt.Println("  No printers found.")
		return nil
	} else if config.JSONOutput {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(printers)
	}

	fmt.Println()
	for _, p := range printers {
		defaultMark := ""
		if p.IsDefault {
			defaultMark = " [DEFAULT]"
		}

		virtualMark := ""
		if p.IsVirtual {
			virtualMark = " (virtual)"
		}

		fmt.Printf("  Name:    %s%s%s\n", p.Name, defaultMark, virtualMark)
		fmt.Printf("  Port:    %s\n", p.Port)
		fmt.Printf("  Driver:  %s\n", p.Driver)
		fmt.Printf("  Status:  %s\n", p.Status)
		fmt.Printf("  Type:    %s\n", p.PrinterType)
		fmt.Println()
	}

	fmt.Printf("Total: %d printer(s)\n", len(printers))
	return nil
}

func executePrint(config *Config) error {
	jsonData, err := loadJSON(config.JSONFile)
	if err != nil {
		return fmt.Errorf("failed to load JSON: %w", err)
	}

	var doc schema.Document
	if err := json.Unmarshal(jsonData, &doc); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	if config.Debug {
		log.Printf("Document version: %s", doc.Version)
		log.Printf("Profile: %s (%dmm)", doc.Profile.Model, doc.Profile.PaperWidth)
		log.Printf("Commands: %d", len(doc.Commands))
		log.Printf("Connection type: %s", config.ConnectionType)
	}

	// Dry run validation executes before attempting hardware connection/detection
	if config.DryRun {
		return validateDocument(&doc)
	}

	printerName, err := getPrinterNameFromDocument(config, jsonData)
	if err != nil && config.ConnectionType == win {
		return err
	}
	config.PrinterName = printerName

	conn, err := createConnection(config)
	if err != nil {
		return fmt.Errorf("failed to create connection: %w", err)
	}

	prof := createProfile(&doc)
	proto := composer.NewEscpos()

	printerService, err := service.NewPrinter(proto, prof, conn)
	if err != nil {
		_ = conn.Close()
		return fmt.Errorf("failed to create printer: %w", err)
	}
	defer func() {
		if closeErr := printerService.Close(); closeErr != nil && config.Debug {
			log.Printf("warning: failed to close printer: %v", closeErr)
		}
	}()

	exec := executor.NewExecutor(printerService)

	if config.Debug {
		log.Println("Executing document...")
	}

	return exec.Execute(&doc)
}
