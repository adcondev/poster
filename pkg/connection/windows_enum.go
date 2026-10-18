//go:build windows

package connection

import (
	"syscall"
	"unsafe"
)

// Windows API constants for EnumPrinters
const (
	PrinterEnumLocal       = 0x00000002
	PrinterEnumConnections = 0x00000004
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
	_, _, _ = procEnumPrinters.Call(
		uintptr(PrinterEnumLocal|PrinterEnumConnections),
		0, 2, 0, 0,
		//nolint:gosec // Required for Windows API
		uintptr(unsafe.Pointer(&needed)),
		//nolint:gosec // Required for Windows API
		uintptr(unsafe.Pointer(&returned)),
	)

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
			Status:      interpretStatus(info.Status, info.Attributes),
			StatusRaw:   info.Status,
			IsDefault:   name == defaultPrinter,
			IsVirtual:   isVirtualPrinter(name, driver, port),
			PrinterType: detectPrinterType(port, driver, name),
		})
	}

	return printers, nil
}

func getDefaultPrinterName() string {
	var size uint32
	//nolint:gosec // Required for Windows API
	_, _, _ = procGetDefaultPrinter.Call(0, uintptr(unsafe.Pointer(&size)))

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

func init() {
	if len(virtualPrinterPatterns) > 0 {
		virtualPrinterPatternsList = virtualPrinterPatterns
	}
	if len(thermalDriverPatterns) > 0 {
		thermalDriverPatternsList = thermalDriverPatterns
	}
}

