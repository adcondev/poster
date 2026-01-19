//go:build windows

package connection

import (
	"strings"
	"syscall"
	"unsafe"
)

// Windows API constants for EnumPrinters
const (
	PrinterEnumLocal       = 0x00000002
	PrinterEnumConnections = 0x00000004
)

// Printer status flags (from winspool.h)
const (
	PrinterStatusPaused     = 0x00000001
	PrinterStatusError      = 0x00000002
	PrinterStatusPaperJam   = 0x00000008
	PrinterStatusPaperOut   = 0x00000010
	PrinterStatusOffline    = 0x00000080
	PrinterStatusBusy       = 0x00000200
	PrinterStatusPrinting   = 0x00000400
	PrinterStatusProcessing = 0x00004000
	PrinterStatusDoorOpen   = 0x00400000
)

// PRINTER_INFO_2 structure (matches Windows API layout)
type printerInfo2 struct {
	pServerName         *uint16
	pPrinterName        *uint16
	pShareName          *uint16
	pPortName           *uint16
	pDriverName         *uint16
	pComment            *uint16
	pLocation           *uint16
	pDevMode            uintptr
	pSepFile            *uint16
	pPrintProcessor     *uint16
	pDatatype           *uint16
	pParameters         *uint16
	pSecurityDescriptor uintptr
	Attributes          uint32
	Priority            uint32
	DefaultPriority     uint32
	StartTime           uint32
	UntilTime           uint32
	Status              uint32
	cJobs               uint32
	AveragePPM          uint32
}

// Additional procs - winspool is already defined in windows.go
var (
	procEnumPrinters      = winspool.NewProc("EnumPrintersW")
	procGetDefaultPrinter = winspool.NewProc("GetDefaultPrinterW")
)

// ListAvailablePrinters enumerates all installed printers on the system.
func ListAvailablePrinters() ([]PrinterDetail, error) {
	// First call:  get required buffer size
	var needed, returned uint32
	_, _, err := procEnumPrinters.Call(
		uintptr(PrinterEnumLocal|PrinterEnumConnections),
		0, 2, 0, 0,
		//nolint:gosec // Required for Windows API
		uintptr(unsafe.Pointer(&needed)),
		//nolint:gosec // Required for Windows API
		uintptr(unsafe.Pointer(&returned)),
	)
	if err != nil {
		return nil, err
	}

	if needed == 0 {
		return []PrinterDetail{}, nil
	}

	// Second call: get actual data
	buf := make([]byte, needed)
	r1, _, err := procEnumPrinters.Call(
		uintptr(PrinterEnumLocal|PrinterEnumConnections),
		0, 2,
		//nolint:gosec // Required for Windows API
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(needed),
		//nolint:gosec // Required for Windows API
		uintptr(unsafe.Pointer(&needed)),
		//nolint:gosec // Required for Windows API
		uintptr(unsafe.Pointer(&returned)),
	)

	if r1 == 0 {
		return nil, err
	}

	defaultPrinter := getDefaultPrinterName()
	if defaultPrinter == "" {
		defaultPrinter = "<none>"
	}
	printers := make([]PrinterDetail, 0, returned)
	infoSize := unsafe.Sizeof(printerInfo2{})

	for i := uint32(0); i < returned; i++ {
		//nolint:gosec // Required for Windows API - casting buffer to struct pointer
		info := (*printerInfo2)(unsafe.Pointer(&buf[uintptr(i)*infoSize]))

		name := utf16PtrToString(info.pPrinterName)
		port := utf16PtrToString(info.pPortName)
		driver := utf16PtrToString(info.pDriverName)

		printers = append(printers, PrinterDetail{
			Name:        name,
			Port:        port,
			Driver:      driver,
			Status:      interpretStatus(info.Status),
			StatusRaw:   info.Status,
			IsDefault:   name == defaultPrinter,
			IsVirtual:   isVirtualPrinter(name, port),
			PrinterType: detectPrinterType(name, port, driver),
		})
	}

	return printers, nil
}

func getDefaultPrinterName() string {
	var size uint32
	//nolint:gosec // Required for Windows API
	_, _, err := procGetDefaultPrinter.Call(0, uintptr(unsafe.Pointer(&size)))
	if err != nil {
		return ""
	}

	if size == 0 {
		return ""
	}

	buf := make([]uint16, size)
	r1, _, _ := procGetDefaultPrinter.Call(
		//nolint:gosec // Required for Windows API
		uintptr(unsafe.Pointer(&buf[0])),
		//nolint:gosec // Required for Windows API
		uintptr(unsafe.Pointer(&size)),
	)

	if r1 == 0 {
		return ""
	}
	return syscall.UTF16ToString(buf)
}

func utf16PtrToString(p *uint16) string {
	if p == nil {
		return ""
	}
	// Find length
	//nolint:gosec // Required for Windows API - pointer arithmetic for UTF-16 string
	end := unsafe.Pointer(p)
	n := 0
	//nolint:gosec // Required for Windows API - reading UTF-16 characters
	for *(*uint16)(unsafe.Pointer(uintptr(end) + uintptr(n)*2)) != 0 {
		n++
	}
	s := make([]uint16, n)
	for i := 0; i < n; i++ {
		//nolint:gosec // Required for Windows API - copying UTF-16 characters
		s[i] = *(*uint16)(unsafe.Pointer(uintptr(end) + uintptr(i)*2))
	}
	return syscall.UTF16ToString(s)
}

func interpretStatus(status uint32) PrinterState {
	switch {
	case status == 0:
		return StateReady
	case status&PrinterStatusOffline != 0:
		return StateOffline
	case status&PrinterStatusPaused != 0:
		return StatePaused
	case status&(PrinterStatusError|PrinterStatusPaperJam|PrinterStatusPaperOut|PrinterStatusDoorOpen) != 0:
		return StateError
	case status&(PrinterStatusPrinting|PrinterStatusProcessing|PrinterStatusBusy) != 0:
		return StateReady
	default:
		return StateUnknown
	}
}

func isVirtualPrinter(name, port string) bool {
	nameLower := strings.ToLower(name)
	portLower := strings.ToLower(port)

	for _, pattern := range virtualPrinterPatterns {
		if strings.Contains(nameLower, pattern) {
			return true
		}
	}

	virtualPorts := []string{"file:", "portprompt:", "nul:", "xpsport:"}
	for _, vp := range virtualPorts {
		if strings.HasPrefix(portLower, vp) {
			return true
		}
	}
	return false
}

func detectPrinterType(name, port, driver string) string {
	if isVirtualPrinter(name, port) {
		return "virtual"
	}

	portLower := strings.ToLower(port)
	if strings.HasPrefix(portLower, "\\\\") || strings.Contains(portLower, "wsd") {
		return "network"
	}

	combined := strings.ToLower(name + " " + driver)
	for _, pattern := range thermalDriverPatterns {
		if strings.Contains(combined, pattern) {
			return "thermal"
		}
	}
	return "unknown"
}

// FilterThermalPrinters returns only thermal/POS printers
func FilterThermalPrinters(printers []PrinterDetail) []PrinterDetail {
	result := make([]PrinterDetail, 0)
	for _, p := range printers {
		if p.PrinterType == "thermal" {
			result = append(result, p)
		}
	}
	return result
}

// FilterPhysicalPrinters returns non-virtual printers
func FilterPhysicalPrinters(printers []PrinterDetail) []PrinterDetail {
	result := make([]PrinterDetail, 0)
	for _, p := range printers {
		if !p.IsVirtual {
			result = append(result, p)
		}
	}
	return result
}
