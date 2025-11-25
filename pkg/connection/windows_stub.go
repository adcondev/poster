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

// Open always returns an error indicating unavailability on non-Windows systems.
func (s *StubPrinterService) Open(_ string) (uintptr, error) {
	return 0, errors.New("WindowsPrintConnector no está disponible en este sistema operativo")
}

// Close always returns an error indicating unavailability on non-Windows systems.
func (s *StubPrinterService) Close(_ uintptr) error {
	return errors.New("WindowsPrintConnector no está disponible en este sistema operativo")
}

// StartDoc always returns an error indicating unavailability on non-Windows systems.
func (s *StubPrinterService) StartDoc(_ uintptr, _, _ string) (uint32, error) {
	return 0, errors.New("WindowsPrintConnector no está disponible en este sistema operativo")
}

// EndDoc always returns an error indicating unavailability on non-Windows systems.
func (s *StubPrinterService) EndDoc(_ uintptr) error {
	return errors.New("WindowsPrintConnector no está disponible en este sistema operativo")
}

// AbortDoc always returns an error indicating unavailability on non-Windows systems.
func (s *StubPrinterService) AbortDoc(_ uintptr) error {
	return errors.New("WindowsPrintConnector no está disponible en este sistema operativo")
}

// Write always returns an error indicating unavailability on non-Windows systems.
func (s *StubPrinterService) Write(_ uintptr, _ []byte) (uint32, error) {
	return 0, errors.New("WindowsPrintConnector no está disponible en este sistema operativo")
}
