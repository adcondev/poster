//go:build !windows

package connection

import "errors"

// ListAvailablePrinters is not available on non-Windows platforms
func ListAvailablePrinters() ([]PrinterDetail, error) {
	return nil, errors.New("printer enumeration is only available on Windows")
}

// FilterThermalPrinters stub
func FilterThermalPrinters(_ []PrinterDetail) []PrinterDetail {
	return nil
}

// FilterPhysicalPrinters stub
func FilterPhysicalPrinters(_ []PrinterDetail) []PrinterDetail {
	return nil
}
