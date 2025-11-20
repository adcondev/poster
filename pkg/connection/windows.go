//go:build windows

package connection

import (
	"syscall"
	"unsafe"
)

// === DLL y funciones ===

var (
	winspool            = syscall.NewLazyDLL("winspool.drv")
	procOpenPrinter     = winspool.NewProc("OpenPrinterW")
	procClosePrinter    = winspool.NewProc("ClosePrinter")
	procStartDocPrinter = winspool.NewProc("StartDocPrinterW")
	procEndDocPrinter   = winspool.NewProc("EndDocPrinter")
	procAbortDocPrinter = winspool.NewProc("AbortDocPrinter")
	procWritePrinter    = winspool.NewProc("WritePrinter")
)

// === Estructura DOC_INFO_1 (corresponde a la API de Windows) ===

type docInfo1 struct {
	DocName    *uint16
	OutputFile *uint16
	DataType   *uint16
}

// RealPrinterService implements PrinterService using Windows API syscalls.
type RealPrinterService struct{}

func getPlatformPrinterService() PrinterService {
	return &RealPrinterService{}
}

func (s *RealPrinterService) Open(name string) (uintptr, error) {
	printerNameUTF16, err := syscall.UTF16PtrFromString(name)
	if err != nil {
		return 0, err
	}

	var h syscall.Handle
	r1, _, err := procOpenPrinter.Call(
		uintptr(unsafe.Pointer(printerNameUTF16)), //nolint:gosec
		uintptr(unsafe.Pointer(&h)),               //nolint:gosec
		0,
	)
	if r1 == 0 {
		return 0, err
	}
	return uintptr(h), nil
}

func (s *RealPrinterService) Close(handle uintptr) error {
	r1, _, err := procClosePrinter.Call(handle)
	if r1 == 0 {
		return err
	}
	return nil
}

func (s *RealPrinterService) StartDoc(handle uintptr, docName, dataType string) (uint32, error) {
	docNamePtr, _ := syscall.UTF16PtrFromString(docName)
	dataTypePtr, _ := syscall.UTF16PtrFromString(dataType)

	doc := &docInfo1{
		DocName:    docNamePtr,
		OutputFile: nil,
		DataType:   dataTypePtr,
	}

	r1, _, err := procStartDocPrinter.Call(handle, 1, uintptr(unsafe.Pointer(doc))) //nolint:gosec
	if r1 == 0 {
		return 0, err
	}
	return uint32(r1), nil
}

func (s *RealPrinterService) EndDoc(handle uintptr) error {
	r1, _, err := procEndDocPrinter.Call(handle)
	if r1 == 0 {
		return err
	}
	return nil
}

func (s *RealPrinterService) AbortDoc(handle uintptr) error {
	r1, _, err := procAbortDocPrinter.Call(handle)
	if r1 == 0 {
		return err
	}
	return nil
}

func (s *RealPrinterService) Write(handle uintptr, data []byte) (uint32, error) {
	//nolint:gosec // Necesario para interactuar con la API de Windows
	var bytesWritten uint32
	r1, _, err := procWritePrinter.Call(
		handle,
		uintptr(unsafe.Pointer(&data[0])), //nolint:gosec
		uintptr(len(data)),
		uintptr(unsafe.Pointer(&bytesWritten)), //nolint:gosec
	)
	if r1 == 0 {
		return 0, err
	}
	return bytesWritten, nil
}
