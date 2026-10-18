//go:build !windows

package connection

import "errors"

// ListAvailablePrinters is not available on non-Windows platforms
func ListAvailablePrinters() ([]PrinterDetail, error) {
	return nil, errors.New("printer enumeration is only available on Windows")
}

