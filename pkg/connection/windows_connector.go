package connection

import (
	"errors"
	"fmt"
	"log"
)

// WindowsPrintConnector implements a connector for Windows printers using the Windows API.
// It uses a PrinterService to abstract the underlying OS calls.
type WindowsPrintConnector struct {
	printerName string
	service     PrinterService
	handle      uintptr
	jobStarted  bool
}

// NewWindowsPrintConnector creates a new connector for the specified printer name.
func NewWindowsPrintConnector(printerName string) (*WindowsPrintConnector, error) {
	if printerName == "" {
		return nil, errors.New("el nombre de la impresora no puede estar vacío")
	}

	// Get the platform-specific service implementation
	// TODO: Check CONN_NITPICKS.md for details on platform-specific architecture issues
	service := getPlatformPrinterService()

	// If we are on a non-windows platform, the stub service might return error on use,
	// but here we want to check if we can even open the printer.
	// However, to maintain backward compatibility with the stub's behavior (which returned error immediately),
	// we can check if the service is functional.
	// Ideally, Open() will fail on the stub.

	handle, err := service.Open(printerName)
	if err != nil {
		// On non-windows, this will likely fail with "not available"
		return nil, fmt.Errorf("no se pudo abrir la impresora '%s': %w", printerName, err)
	}

	return &WindowsPrintConnector{
		printerName: printerName,
		service:     service,
		handle:      handle,
		jobStarted:  false,
	}, nil
}

// Write writes data to the printer.
func (c *WindowsPrintConnector) Write(data []byte) (int, error) {
	// TODO: Refactor handling of invalid handle
	if c.handle == 0 {
		return 0, errors.New("handle de impresora no válido")
	}

	if !c.jobStarted {
		// Default values as per original implementation
		docName := "ESC/POS PrintDataInPageMode Job"
		dataType := "RAW"

		jobID, err := c.service.StartDoc(c.handle, docName, dataType)
		if err != nil {
			return 0, fmt.Errorf("no se pudo iniciar el trabajo de impresión: %w", err)
		}
		log.Printf("Trabajo de impresión iniciado (ID: %d)", jobID)
		c.jobStarted = true
	}

	bytesWritten, err := c.service.Write(c.handle, data)
	if err != nil {
		return int(bytesWritten), fmt.Errorf("falló al escribir en la impresora: %w", err)
	}

	// TODO: Verify if strict equality check is appropriate for all drivers
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
		err := c.service.EndDoc(c.handle)
		if err != nil {
			log.Printf("Falló EndDocPrinter: %v, intentando AbortDocPrinter...", err)
			if abortErr := c.service.AbortDoc(c.handle); abortErr != nil {
				log.Printf("Falló AbortDocPrinter: %v", abortErr)
				finalErr = fmt.Errorf("falló EndDoc y AbortDoc: %v", abortErr)
			}
		} else {
			log.Println("Trabajo de impresión finalizado correctamente.")
		}
	}

	if c.handle != 0 {
		if err := c.service.Close(c.handle); err != nil {
			log.Printf("Falló ClosePrinter: %v", err)
			if finalErr == nil {
				finalErr = fmt.Errorf("falló ClosePrinter: %w", err)
			}
		}
		c.handle = 0
		c.jobStarted = false
	}

	return finalErr
}

// Read implements the io.Reader interface but is not supported by the Windows Spooler API directly.
func (c *WindowsPrintConnector) Read(buf []byte) (int, error) {
	l := len(buf)
	// TODO: Check if bi-directional communication is possible or necessary
	// FIXME: Implement a way to read status if possible, or clarify limitation in docs
	return l, errors.New("spooler no soporta lectura de estado de impresora directamente")
}
