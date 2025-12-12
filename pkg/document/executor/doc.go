// Package executor parses and executes print documents on ESC/POS printers.
//
// This package is the runtime component of the poster library.  It takes JSON
// documents (created by the builder package or manually) and sends the
// corresponding commands to a physical or virtual printer.
//
// # Quick Start
//
//	conn, _ := connection.NewWindowsPrintConnector("EPSON TM-T20II")
//	printer, _ := service.NewPrinter(composer. NewEscpos(), profile.NewEscpos(), conn)
//	defer printer.Close()
//
//	exec := executor.NewExecutor(printer)
//	err := exec.ExecuteJSON(jsonData)
//
// # Architecture
//
// The executor uses a handler registry pattern.  Each command type has a
// dedicated handler that translates JSON data into printer commands:
//
//	Document JSON → ParseDocument() → Execute() → handlers → Printer
//
// # Built-in Commands
//
//	text        Styled text with labels
//	image       Base64 images
//	qr          QR codes (native/image fallback)
//	barcode     1D barcodes (CODE128, EAN13, etc.)
//	table       Formatted tables
//	separator   Line separators
//	feed        Paper advance
//	cut         Paper cutting (full/partial)
//	pulse       Cash drawer activation
//	beep        Buzzer sound
//	raw         Direct ESC/POS bytes
//
// # Error Handling
//
// Errors include context about which command failed:
//
//	command 3 (barcode) failed: barcode data is required
//
// # Testing
//
// Handler tests follow a standardized three-category structure:
//
//   - Parsing:  Verify JSON unmarshaling (2-4 cases)
//   - Defaults:  Verify nil/zero handling (1-3 cases)
//   - Validation: Verify invalid JSON rejection (2-3 cases)
//
// Tests verify JSON parsing only—handler business logic is tested separately.
//
// # Key Types
//
//	Document        Complete print document structure
//	ProfileConfig   Printer configuration (model, paper width, etc.)
//	Command         Single command with type and JSON data
//	Executor        Main executor with handler registry
//	CommandHandler  Function signature:  func(*service.Printer, json.RawMessage) error
//	HandlerRegistry Registry for custom command handlers
//
// # Thread Safety
//
// Executor instances are NOT safe for concurrent use. Create one per
// goroutine or synchronize access externally.
//
// # Related Packages
//
//	builder   Fluent API to create Document JSON
//	service   Low-level printer communication
//	profile   Printer capability profiles
package executor
