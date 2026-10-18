package connection

import (
	"fmt"
	"strings"
)

// Printer type string constants.
const (
	// PrinterTypeThermal indicates a thermal ESC/POS or receipt printer.
	PrinterTypeThermal = "thermal"
	// PrinterTypeVirtual indicates a virtual or software printer (PDF, XPS, etc.).
	PrinterTypeVirtual = "virtual"
	// PrinterTypeNetwork indicates a network-connected printer.
	PrinterTypeNetwork = "network"
	// PrinterTypeUSB indicates a USB-connected printer.
	PrinterTypeUSB = "usb"
	// PrinterTypeSerial indicates a serial/parallel connected printer.
	PrinterTypeSerial = "serial"
	// PrinterTypeUnknown indicates an unknown printer type.
	PrinterTypeUnknown = "unknown"
)

// Printer status flags (from winspool.h)
const (
	// PrinterStatusPaused indicates the printer is paused.
	PrinterStatusPaused uint32 = 0x00000001
	// PrinterStatusError indicates the printer is in an error state.
	PrinterStatusError uint32 = 0x00000002
	// PrinterStatusPendingDeletion indicates the printer is pending deletion.
	PrinterStatusPendingDeletion uint32 = 0x00000004
	// PrinterStatusPaperJam indicates paper is jammed in the printer.
	PrinterStatusPaperJam uint32 = 0x00000008
	// PrinterStatusPaperOut indicates the printer is out of paper.
	PrinterStatusPaperOut uint32 = 0x00000010
	// PrinterStatusManualFeed indicates the printer is in manual feed mode.
	PrinterStatusManualFeed uint32 = 0x00000020
	// PrinterStatusPaperProblem indicates a paper problem in the printer.
	PrinterStatusPaperProblem uint32 = 0x00000040
	// PrinterStatusOffline indicates the printer is offline.
	PrinterStatusOffline uint32 = 0x00000080
	// PrinterStatusIOActive indicates the printer I/O is active.
	PrinterStatusIOActive uint32 = 0x00000100
	// PrinterStatusBusy indicates the printer is busy.
	PrinterStatusBusy uint32 = 0x00000200
	// PrinterStatusPrinting indicates the printer is currently printing.
	PrinterStatusPrinting uint32 = 0x00000400
	// PrinterStatusOutputBinFull indicates the output bin is full.
	PrinterStatusOutputBinFull uint32 = 0x00000800
	// PrinterStatusNotAvailable indicates the printer is not available.
	PrinterStatusNotAvailable uint32 = 0x00001000
	// PrinterStatusWaiting indicates the printer is waiting.
	PrinterStatusWaiting uint32 = 0x00002000
	// PrinterStatusProcessing indicates the printer is processing a job.
	PrinterStatusProcessing uint32 = 0x00004000
	// PrinterStatusInitializing indicates the printer is initializing.
	PrinterStatusInitializing uint32 = 0x00008000
	// PrinterStatusWarmingUp indicates the printer is warming up.
	PrinterStatusWarmingUp uint32 = 0x00010000
	// PrinterStatusTonerLow indicates the toner is low.
	PrinterStatusTonerLow uint32 = 0x00020000
	// PrinterStatusNoToner indicates the printer has no toner.
	PrinterStatusNoToner uint32 = 0x00040000
	// PrinterStatusPagePunt indicates the printer cannot print the current page.
	PrinterStatusPagePunt uint32 = 0x00080000
	// PrinterStatusUserIntervention indicates user intervention is required.
	PrinterStatusUserIntervention uint32 = 0x00100000
	// PrinterStatusOutOfMemory indicates the printer is out of memory.
	PrinterStatusOutOfMemory uint32 = 0x00200000
	// PrinterStatusDoorOpen indicates the printer door is open.
	PrinterStatusDoorOpen uint32 = 0x00400000
	// PrinterStatusServerUnknown indicates the printer server is unknown.
	PrinterStatusServerUnknown uint32 = 0x00800000
	// PrinterStatusPowerSave indicates the printer is in power save mode.
	PrinterStatusPowerSave uint32 = 0x01000000
)

// Windows winspool.h PRINTER_STATUS_* aliases
//
//nolint:revive // Windows winspool.h API compatibility constants
const (
	// PRINTER_STATUS_PAUSED indicates the printer is paused.
	PRINTER_STATUS_PAUSED = PrinterStatusPaused
	// PRINTER_STATUS_ERROR indicates the printer is in an error state.
	PRINTER_STATUS_ERROR = PrinterStatusError
	// PRINTER_STATUS_PENDING_DELETION indicates the printer is pending deletion.
	PRINTER_STATUS_PENDING_DELETION = PrinterStatusPendingDeletion
	// PRINTER_STATUS_PAPER_JAM indicates paper is jammed in the printer.
	PRINTER_STATUS_PAPER_JAM = PrinterStatusPaperJam
	// PRINTER_STATUS_PAPER_OUT indicates the printer is out of paper.
	PRINTER_STATUS_PAPER_OUT = PrinterStatusPaperOut
	// PRINTER_STATUS_MANUAL_FEED indicates the printer is in manual feed mode.
	PRINTER_STATUS_MANUAL_FEED = PrinterStatusManualFeed
	// PRINTER_STATUS_PAPER_PROBLEM indicates a paper problem in the printer.
	PRINTER_STATUS_PAPER_PROBLEM = PrinterStatusPaperProblem
	// PRINTER_STATUS_OFFLINE indicates the printer is offline.
	PRINTER_STATUS_OFFLINE = PrinterStatusOffline
	// PRINTER_STATUS_IO_ACTIVE indicates the printer I/O is active.
	PRINTER_STATUS_IO_ACTIVE = PrinterStatusIOActive
	// PRINTER_STATUS_BUSY indicates the printer is busy.
	PRINTER_STATUS_BUSY = PrinterStatusBusy
	// PRINTER_STATUS_PRINTING indicates the printer is currently printing.
	PRINTER_STATUS_PRINTING = PrinterStatusPrinting
	// PRINTER_STATUS_OUTPUT_BIN_FULL indicates the output bin is full.
	PRINTER_STATUS_OUTPUT_BIN_FULL = PrinterStatusOutputBinFull
	// PRINTER_STATUS_NOT_AVAILABLE indicates the printer is not available.
	PRINTER_STATUS_NOT_AVAILABLE = PrinterStatusNotAvailable
	// PRINTER_STATUS_WAITING indicates the printer is waiting.
	PRINTER_STATUS_WAITING = PrinterStatusWaiting
	// PRINTER_STATUS_PROCESSING indicates the printer is processing a job.
	PRINTER_STATUS_PROCESSING = PrinterStatusProcessing
	// PRINTER_STATUS_INITIALIZING indicates the printer is initializing.
	PRINTER_STATUS_INITIALIZING = PrinterStatusInitializing
	// PRINTER_STATUS_WARMING_UP indicates the printer is warming up.
	PRINTER_STATUS_WARMING_UP = PrinterStatusWarmingUp
	// PRINTER_STATUS_TONER_LOW indicates the toner is low.
	PRINTER_STATUS_TONER_LOW = PrinterStatusTonerLow
	// PRINTER_STATUS_NO_TONER indicates the printer has no toner.
	PRINTER_STATUS_NO_TONER = PrinterStatusNoToner
	// PRINTER_STATUS_PAGE_PUNT indicates the printer cannot print the current page.
	PRINTER_STATUS_PAGE_PUNT = PrinterStatusPagePunt
	// PRINTER_STATUS_USER_INTERVENTION indicates user intervention is required.
	PRINTER_STATUS_USER_INTERVENTION = PrinterStatusUserIntervention
	// PRINTER_STATUS_OUT_OF_MEMORY indicates the printer is out of memory.
	PRINTER_STATUS_OUT_OF_MEMORY = PrinterStatusOutOfMemory
	// PRINTER_STATUS_DOOR_OPEN indicates the printer door is open.
	PRINTER_STATUS_DOOR_OPEN = PrinterStatusDoorOpen
	// PRINTER_STATUS_SERVER_UNKNOWN indicates the printer server is unknown.
	PRINTER_STATUS_SERVER_UNKNOWN = PrinterStatusServerUnknown
	// PRINTER_STATUS_POWER_SAVE indicates the printer is in power save mode.
	PRINTER_STATUS_POWER_SAVE = PrinterStatusPowerSave
)

// Printer attribute flags (from winspool.h)
const (
	// PrinterAttributeQueued indicates print jobs are queued.
	PrinterAttributeQueued uint32 = 0x00000001
	// PrinterAttributeDirect indicates jobs are sent directly to the printer.
	PrinterAttributeDirect uint32 = 0x00000002
	// PrinterAttributeDefault indicates this is the default printer.
	PrinterAttributeDefault uint32 = 0x00000004
	// PrinterAttributeShared indicates the printer is shared.
	PrinterAttributeShared uint32 = 0x00000008
	// PrinterAttributeNetwork indicates the printer is a network printer.
	PrinterAttributeNetwork uint32 = 0x00000010
	// PrinterAttributeHidden indicates the printer is hidden.
	PrinterAttributeHidden uint32 = 0x00000020
	// PrinterAttributeLocal indicates the printer is directly connected.
	PrinterAttributeLocal uint32 = 0x00000040
	// PrinterAttributeEnableDevq indicates devquery is enabled.
	PrinterAttributeEnableDevq uint32 = 0x00000080
	// PrinterAttributeKeepprintedjobs indicates printed jobs should be kept.
	PrinterAttributeKeepprintedjobs uint32 = 0x00000100
	// PrinterAttributeDoCompleteFirst indicates jobs should be completed before others.
	PrinterAttributeDoCompleteFirst uint32 = 0x00000200
	// PrinterAttributeWorkOffline indicates the printer is configured to work offline.
	PrinterAttributeWorkOffline uint32 = 0x00000400
	// PrinterAttributeEnableBidi indicates bidirectional communication is enabled.
	PrinterAttributeEnableBidi uint32 = 0x00000800
	// PrinterAttributeRawOnly indicates only raw data type is accepted.
	PrinterAttributeRawOnly uint32 = 0x00001000
	// PrinterAttributePublished indicates the printer is published in directory service.
	PrinterAttributePublished uint32 = 0x00002000
)

// Windows winspool.h PRINTER_ATTRIBUTE_* aliases
//
//nolint:revive // Windows winspool.h API compatibility constants
const (
	// PRINTER_ATTRIBUTE_QUEUED indicates print jobs are queued.
	PRINTER_ATTRIBUTE_QUEUED = PrinterAttributeQueued
	// PRINTER_ATTRIBUTE_DIRECT indicates jobs are sent directly to the printer.
	PRINTER_ATTRIBUTE_DIRECT = PrinterAttributeDirect
	// PRINTER_ATTRIBUTE_DEFAULT indicates this is the default printer.
	PRINTER_ATTRIBUTE_DEFAULT = PrinterAttributeDefault
	// PRINTER_ATTRIBUTE_SHARED indicates the printer is shared.
	PRINTER_ATTRIBUTE_SHARED = PrinterAttributeShared
	// PRINTER_ATTRIBUTE_NETWORK indicates the printer is a network printer.
	PRINTER_ATTRIBUTE_NETWORK = PrinterAttributeNetwork
	// PRINTER_ATTRIBUTE_HIDDEN indicates the printer is hidden.
	PRINTER_ATTRIBUTE_HIDDEN = PrinterAttributeHidden
	// PRINTER_ATTRIBUTE_LOCAL indicates the printer is directly connected.
	PRINTER_ATTRIBUTE_LOCAL = PrinterAttributeLocal
	// PRINTER_ATTRIBUTE_ENABLE_DEVQ indicates devquery is enabled.
	PRINTER_ATTRIBUTE_ENABLE_DEVQ = PrinterAttributeEnableDevq
	// PRINTER_ATTRIBUTE_KEEPPRINTEDJOBS indicates printed jobs should be kept.
	PRINTER_ATTRIBUTE_KEEPPRINTEDJOBS = PrinterAttributeKeepprintedjobs
	// PRINTER_ATTRIBUTE_DO_COMPLETE_FIRST indicates jobs should be completed before others.
	PRINTER_ATTRIBUTE_DO_COMPLETE_FIRST = PrinterAttributeDoCompleteFirst
	// PRINTER_ATTRIBUTE_WORK_OFFLINE indicates the printer is configured to work offline.
	PRINTER_ATTRIBUTE_WORK_OFFLINE = PrinterAttributeWorkOffline
	// PRINTER_ATTRIBUTE_ENABLE_BIDI indicates bidirectional communication is enabled.
	PRINTER_ATTRIBUTE_ENABLE_BIDI = PrinterAttributeEnableBidi
	// PRINTER_ATTRIBUTE_RAW_ONLY indicates only raw data type is accepted.
	PRINTER_ATTRIBUTE_RAW_ONLY = PrinterAttributeRawOnly
	// PRINTER_ATTRIBUTE_PUBLISHED indicates the printer is published in directory service.
	PRINTER_ATTRIBUTE_PUBLISHED = PrinterAttributePublished
)

// virtualPrinterPatternsList contains substrings indicating virtual/software printers.
var virtualPrinterPatternsList = []string{
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

// thermalDriverPatternsList contains patterns suggesting ESC/POS or receipt thermal printers.
var thermalDriverPatternsList = []string{
	"epson",
	"star",
	"citizen",
	"bixolon",
	"sewoo",
	"pos-",
	"pos ",
	PrinterTypeThermal,
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

// isVirtualPrinter checks whether the printer name, driver, or port corresponds to a virtual printer.
func isVirtualPrinter(name, driver, port string) bool {
	combined := strings.ToLower(name + " " + driver)
	for _, pattern := range virtualPrinterPatternsList {
		if strings.Contains(combined, pattern) {
			return true
		}
	}

	portLower := strings.ToLower(port)
	virtualPorts := []string{"file:", "portprompt:", "nul:", "xpsport:"}
	for _, vp := range virtualPorts {
		if strings.HasPrefix(portLower, vp) {
			return true
		}
	}
	return false
}

// detectPrinterType categorizes a printer as "virtual", "thermal", "network", "usb", "serial", or "unknown".
func detectPrinterType(portName, driverName, printerName string) string {
	if isVirtualPrinter(printerName, driverName, portName) {
		return PrinterTypeVirtual
	}

	combined := strings.ToLower(printerName + " " + driverName)
	for _, pattern := range thermalDriverPatternsList {
		if strings.Contains(combined, pattern) {
			return PrinterTypeThermal
		}
	}

	portLower := strings.ToLower(portName)
	if strings.HasPrefix(portLower, `\\`) || strings.Contains(portLower, "wsd") ||
		strings.HasPrefix(portLower, "ip_") || strings.HasPrefix(portLower, "tcp") {
		return PrinterTypeNetwork
	}

	if strings.HasPrefix(portLower, "usb") {
		return PrinterTypeUSB
	}

	if strings.HasPrefix(portLower, "com") || strings.HasPrefix(portLower, "lpt") ||
		strings.Contains(portLower, "serial") {
		return PrinterTypeSerial
	}

	return PrinterTypeUnknown
}

// interpretStatus converts Windows printer status and attribute bitmasks into a PrinterState.
func interpretStatus(status, attributes uint32) PrinterState {
	switch {
	case status&PrinterStatusOffline != 0 || attributes&PrinterAttributeWorkOffline != 0:
		return StateOffline
	case status&PrinterStatusPaused != 0:
		return StatePaused
	case status&(PrinterStatusError|PrinterStatusPaperJam|PrinterStatusPaperOut|PrinterStatusDoorOpen|
		PrinterStatusPaperProblem|PrinterStatusNoToner|PrinterStatusOutOfMemory|PrinterStatusUserIntervention) != 0:
		return StateError
	case status&(PrinterStatusPrinting|PrinterStatusProcessing|PrinterStatusBusy) != 0:
		return StateReady
	case status == 0:
		return StateReady
	default:
		return StateUnknown
	}
}

// FilterPhysicalPrinters returns non-virtual printers.
func FilterPhysicalPrinters(printers []PrinterDetail) []PrinterDetail {
	if len(printers) == 0 {
		return []PrinterDetail{}
	}
	result := make([]PrinterDetail, 0, len(printers))
	for _, p := range printers {
		if !p.IsVirtual {
			result = append(result, p)
		}
	}
	return result
}

// FilterThermalPrinters returns only thermal/POS printers.
func FilterThermalPrinters(printers []PrinterDetail) []PrinterDetail {
	if len(printers) == 0 {
		return []PrinterDetail{}
	}
	result := make([]PrinterDetail, 0, len(printers))
	for _, p := range printers {
		if p.PrinterType == PrinterTypeThermal {
			result = append(result, p)
		}
	}
	return result
}

// FindThermalPrinters enumerates system printers and returns only thermal/POS printers.
func FindThermalPrinters() ([]PrinterDetail, error) {
	printers, err := ListAvailablePrinters()
	if err != nil {
		return nil, fmt.Errorf("list available printers: %w", err)
	}
	return FilterThermalPrinters(printers), nil
}
