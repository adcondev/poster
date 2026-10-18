package connection

import (
	"errors"
	"fmt"
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
		return nil, errors.New("printer name cannot be empty")
	}

	service := getPlatformPrinterService()

	handle, err := service.Open(printerName)
	if err != nil {
		return nil, fmt.Errorf("failed to open printer '%s': %w", printerName, err)
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
	if c.handle == 0 {
		return 0, errors.New("handle de impresora no válido")
	}

	if !c.jobStarted {
		docName := "ESC/POS PrintDataInPageMode Job"
		dataType := "RAW"

		_, err := c.service.StartDoc(c.handle, docName, dataType)
		if err != nil {
			return 0, fmt.Errorf("no se pudo iniciar el trabajo de impresión: %w", err)
		}
		c.jobStarted = true
	}

	bytesWritten, err := c.service.Write(c.handle, data)
	if err != nil {
		return int(bytesWritten), fmt.Errorf("falló al escribir en la impresora: %w", err)
	}

	if int(bytesWritten) != len(data) {
		return int(bytesWritten), fmt.Errorf("solo se escribieron %d de %d bytes", bytesWritten, len(data))
	}

	return int(bytesWritten), nil
}

// Close ends the print job and closes the printer handle.
func (c *WindowsPrintConnector) Close() error {
	var finalErr error

	if c.jobStarted {
		func() {
			defer func() {
				if r := recover(); r != nil {
					finalErr = fmt.Errorf("pánico durante EndDocPrinter: %v", r)
				}
			}()
			err := c.service.EndDoc(c.handle)
			if err != nil {
				if abortErr := c.service.AbortDoc(c.handle); abortErr != nil {
					finalErr = fmt.Errorf("falló EndDoc y AbortDoc: %w", abortErr)
				}
			}
		}()
	}

	if c.handle != 0 {
		if err := c.service.Close(c.handle); err != nil {
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
	return l, fmt.Errorf("lectura directa desde impresora Windows no soportada por este conector")
}
