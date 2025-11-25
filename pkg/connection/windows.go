//go:build windows

package connection

import (
	"errors"
	"fmt"
	"log"
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

// === Estructura del conector ===

// WindowsPrintConnector implements a connector for Windows printers using the Windows API.
type WindowsPrintConnector struct {
	printerName   string
	printerHandle syscall.Handle
	jobStarted    bool
	docInfo       *docInfo1
}

// === Constructor ===

// NewWindowsPrintConnector creates a new connector for the specified printer name.
func NewWindowsPrintConnector(printerName string) (*WindowsPrintConnector, error) {
	if printerName == "" {
		return nil, errors.New("el nombre de la impresora no puede estar vacío")
	}

	printerNameUTF16, err := syscall.UTF16PtrFromString(printerName)
	if err != nil {
		return nil, fmt.Errorf("error al convertir el nombre de la impresora: %w", err)
	}

	handle, err := openPrinter(printerNameUTF16)
	if err != nil {
		return nil, fmt.Errorf("no se pudo abrir la impresora '%s': %w", printerName, err)
	}

	docName, _ := syscall.UTF16PtrFromString("ESC/POS PrintDataInPageMode Job")
	dataType, _ := syscall.UTF16PtrFromString("RAW")

	doc := &docInfo1{
		DocName:    docName,
		OutputFile: nil,
		DataType:   dataType,
	}

	return &WindowsPrintConnector{
		printerName:   printerName,
		printerHandle: handle,
		jobStarted:    false,
		docInfo:       doc,
	}, nil
}

// === Métodos de la API ===

func (c *WindowsPrintConnector) Write(data []byte) (int, error) {
	if c.printerHandle == 0 {
		return 0, errors.New("handle de impresora no válido")
	}

	if !c.jobStarted {
		jobID, err := startDocPrinter(c.printerHandle, c.docInfo)
		if err != nil {
			return 0, fmt.Errorf("no se pudo iniciar el trabajo de impresión: %w", err)
		}
		log.Printf("Trabajo de impresión iniciado (ID: %d)", jobID)
		c.jobStarted = true
	}

	bytesWritten, err := writePrinter(c.printerHandle, data)
	if err != nil {
		return int(bytesWritten), fmt.Errorf("falló al escribir en la impresora: %w", err)
	}

	if int(bytesWritten) != len(data) {
		log.Printf("Advertencia: solo se escribieron %d de %d bytes", bytesWritten, len(data))
		return int(bytesWritten), fmt.Errorf("solo se escribieron %d de %d bytes", bytesWritten, len(data))
	}

	return int(bytesWritten), nil
}

// Close ends the print job and closes the printer handle.
func (c *WindowsPrintConnector) Close() error {
	var finalErr error

	if c.jobStarted {
		err := endDocPrinter(c.printerHandle)
		if err != nil {
			log.Printf("Falló EndDocPrinter: %v, intentando AbortDocPrinter...", err)
			if abortErr := abortDocPrinter(c.printerHandle); abortErr != nil {
				log.Printf("Falló AbortDocPrinter: %v", abortErr)
				finalErr = fmt.Errorf("falló EndDoc y AbortDoc: %v", abortErr)
			}
		} else if r := recover(); r != nil {
			// En caso de pánico, intentar abortar el trabajo
			log.Printf("Pánico detectado: %v, intentando AbortDocPrinter...", r)
			if abortErr := abortDocPrinter(c.printerHandle); abortErr != nil {
				log.Printf("Falló AbortDocPrinter: %v", abortErr)
				finalErr = fmt.Errorf("pánico: %v, falló AbortDoc: %w", r, abortErr)
			}
		} else {
			log.Printf("Trabajo de impresión finalizado correctamente")
		}
	}

	if c.printerHandle != 0 {
		if err := closePrinter(c.printerHandle); err != nil {
			log.Printf("Falló ClosePrinter: %v", err)
			if finalErr == nil {
				finalErr = fmt.Errorf("falló ClosePrinter: %w", err)
			}
		}
		c.printerHandle = 0
		c.jobStarted = false
	}

	return finalErr
}

// FUNCIONES AUXILIARES
func openPrinter(name *uint16) (handle syscall.Handle, err error) {
	var h syscall.Handle
	r1, _, err := procOpenPrinter.Call(
		uintptr(unsafe.Pointer(name)), //nolint:gosec
		uintptr(unsafe.Pointer(&h)),   //nolint:gosec
		0,
	)
	if r1 == 0 {
		return 0, err
	}
	return h, nil
}

func closePrinter(handle syscall.Handle) error {
	r1, _, err := procClosePrinter.Call(uintptr(handle))
	if r1 == 0 {
		return err
	}
	return nil
}

func startDocPrinter(handle syscall.Handle, docInfo *docInfo1) (uint32, error) {
	// Necesario para interactuar con la API de Windows
	r1, _, err := procStartDocPrinter.Call(uintptr(handle), 1, uintptr(unsafe.Pointer(docInfo))) //nolint:gosec
	if r1 == 0 {
		return 0, err
	}
	return uint32(r1), nil
}

func endDocPrinter(handle syscall.Handle) error {
	r1, _, err := procEndDocPrinter.Call(uintptr(handle))
	if r1 == 0 {
		return err
	}
	return nil
}

func abortDocPrinter(handle syscall.Handle) error {
	r1, _, err := procAbortDocPrinter.Call(uintptr(handle))
	if r1 == 0 {
		return err
	}
	return nil
}

func writePrinter(handle syscall.Handle, data []byte) (uint32, error) {
	//nolint:gosec // Necesario para interactuar con la API de Windows
	var bytesWritten uint32
	r1, _, err := procWritePrinter.Call(
		uintptr(handle),
		uintptr(unsafe.Pointer(&data[0])), //nolint:gosec
		uintptr(len(data)),
		uintptr(unsafe.Pointer(&bytesWritten)), //nolint:gosec
	)
	if r1 == 0 {
		return 0, err
	}
	return bytesWritten, nil
}

func (c *WindowsPrintConnector) Read(buf []byte) (int, error) {
	l := len(buf) // No implementado, ya que Spooler no soporta lectura de estado de impresora directamente
	return l, errors.New("spooler no soporta lectura de estado de impresora directamente")
}
