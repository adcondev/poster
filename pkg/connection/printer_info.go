//go:build windows

package connection

// PrinterDetail contains detailed information about an installed printer.
// Note: Status reflects the last known state from the Windows Spooler,
// which may not reflect real-time connectivity for USB/Serial printers.
type PrinterDetail struct {
	Name        string       `json:"name"`
	Port        string       `json:"port"`
	Driver      string       `json:"driver"`
	Status      PrinterState `json:"status"`
	StatusRaw   uint32       `json:"status_raw,omitempty"`
	IsDefault   bool         `json:"is_default"`
	IsVirtual   bool         `json:"is_virtual"`
	PrinterType string       `json:"printer_type"` // "thermal", "virtual", "network", "unknown"
}

// PrinterState represents the interpreted printer status
type PrinterState string

const (
	// StateReady indicates the printer is ready to print
	StateReady PrinterState = "ready"
	// StateOffline indicates the printer is offline
	StateOffline PrinterState = "offline"
	// StatePaused indicates the printer is paused
	StatePaused PrinterState = "paused"
	// StateError indicates the printer is in an error state
	StateError PrinterState = "error"
	// StateUnknown indicates the printer state is unknown
	StateUnknown PrinterState = "unknown"
)

// virtualPrinterPatterns - names that indicate virtual/software printers
var virtualPrinterPatterns = []string{
	"microsoft print to pdf",
	"microsoft xps document writer",
	"onenote",
	"fax",
	"send to onenote",
	"adobe pdf",
	"cutepdf",
	"pdfcreator",
	"foxit",
	"dopdf",
	"bullzip",
	"primopdf",
	"nitro",
}

// thermalDriverPatterns - patterns that suggest ESC/POS thermal printers
var thermalDriverPatterns = []string{
	"epson",
	"star",
	"citizen",
	"bixolon",
	"sewoo",
	"pos-",
	"thermal",
	"receipt",
	"esc/pos",
	"zj-",
	"xp-",
	"pt-",
	"ec-pm",
	"gp-",
	"tsc",
	"zebra",
	"honeywell",
	"datamax",
	"58mm",
	"80mm",
}
