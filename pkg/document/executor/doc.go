// Package executor parses and executes print documents on ESC/POS printers.
//
// This package is the runtime component of the poster library.  It takes JSON
// documents (created by the builder package or manually) and sends the
// corresponding commands to a physical or virtual printer.
//
// # Quick Start
//
//	// Setup printer connection
//	conn, _ := connection.NewWindowsPrintConnector("EPSON TM-T20II")
//	printer, _ := service.NewPrinter(composer.NewEscpos(), profile.NewEscpos(), conn)
//	defer printer.Close()
//
//	// Execute a JSON document
//	exec := executor.NewExecutor(printer)
//	err := exec.ExecuteJSON(jsonData)
//
// # Architecture
//
// The executor uses a handler registry pattern.  Each command type ("text",
// "qr", "table", etc.) has a dedicated handler function that knows how to
// translate the JSON data into printer commands.
//
//	Document JSON → ParseDocument() → Execute() → handlers → Printer
//
// # Built-in Commands
//
//	Type        Description                 Handler
//	────────────────────────────────────────────────────
//	text        Styled text with labels     handleText
//	image       Base64 images               handleImage
//	qr          QR codes (native/image)     handleQR
//	barcode     1D barcodes                 handleBarcode
//	table       Formatted tables            handleTable
//	separator   Line separators             handleSeparator
//	feed        Paper advance               handleFeed
//	cut         Paper cutting               handleCut
//	pulse       Cash drawer activation      handlePulse
//	beep        Buzzer sound                handleBeep
//	raw         Direct ESC/POS bytes        handleRaw
//
// # Custom Handlers
//
// Extend the executor with custom command types:
//
//	exec := executor.NewExecutor(printer)
//	exec.registerHandler("logo", func(p *service.Printer, data json.RawMessage) error {
//	    var cmd LogoCommand
//	    json.Unmarshal(data, &cmd)
//	    // ...  process and print
//	    return nil
//	})
//
// # Error Handling
//
// All errors include context about which command failed:
//
//	command 3 (barcode) failed: barcode data is required
//
// # Key Types
//
//	Document       - Complete print document structure
//	ProfileConfig  - Printer configuration (model, paper width, etc.)
//	Command        - Single command with type and JSON data
//	Executor       - Main executor with handler registry
//	CommandHandler - Function signature for handlers
//
// # Thread Safety
//
// Executor instances are NOT safe for concurrent use. Create one executor
// per goroutine or synchronize access externally.
//
// # Related Packages
//
//	builder  - Fluent API to create Document JSON
//	service  - Low-level printer communication
//	profile  - Printer capability profiles
package executor
