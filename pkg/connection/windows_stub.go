//go:build !windows

package connection

import (
	"errors"
)

// StubPrinterService implements PrinterService for non-Windows platforms (always returns errors).
type StubPrinterService struct{}

func getPlatformPrinterService() PrinterService {
	return &StubPrinterService{}
}

func (s *StubPrinterService) Open(name string) (uintptr, error) {
	return 0, errors.New("WindowsPrintConnector no está disponible en este sistema operativo")
}

func (s *StubPrinterService) Close(handle uintptr) error {
	return errors.New("WindowsPrintConnector no está disponible en este sistema operativo")
}

func (s *StubPrinterService) StartDoc(handle uintptr, docName, dataType string) (uint32, error) {
	return 0, errors.New("WindowsPrintConnector no está disponible en este sistema operativo")
}

func (s *StubPrinterService) EndDoc(handle uintptr) error {
	return errors.New("WindowsPrintConnector no está disponible en este sistema operativo")
}

func (s *StubPrinterService) AbortDoc(handle uintptr) error {
	return errors.New("WindowsPrintConnector no está disponible en este sistema operativo")
}

func (s *StubPrinterService) Write(handle uintptr, data []byte) (uint32, error) {
	return 0, errors.New("WindowsPrintConnector no está disponible en este sistema operativo")
}
