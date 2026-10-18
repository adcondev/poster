package qrcode

import (
	"errors"
	"fmt"

	"github.com/adcondev/poster/pkg/commands/shared"
)

// PDF417 command constants and limits.
const (
	// MaxPDF417Columns defines the maximum number of columns in a PDF417 barcode.
	MaxPDF417Columns = 30
	// MinPDF417Rows defines the minimum number of rows when specified.
	MinPDF417Rows = 3
	// MaxPDF417Rows defines the maximum number of rows in a PDF417 barcode.
	MaxPDF417Rows = 90
	// MinPDF417ModuleWidth defines the minimum module width in dots.
	MinPDF417ModuleWidth = 2
	// MaxPDF417ModuleWidth defines the maximum module width in dots.
	MaxPDF417ModuleWidth = 8
	// MinPDF417RowHeight defines the minimum row height multiplier.
	MinPDF417RowHeight = 2
	// MaxPDF417RowHeight defines the maximum row height multiplier.
	MaxPDF417RowHeight = 8
	// MaxPDF417ErrorCorrection defines the maximum error correction level.
	MaxPDF417ErrorCorrection = 8
	// MaxPDF417DataLength defines the maximum byte length of PDF417 data.
	MaxPDF417DataLength = 65532
)

var (
	// ErrPDF417Columns indicates an invalid column count (0-30).
	ErrPDF417Columns = errors.New("invalid PDF417 columns (try 0-30)")
	// ErrPDF417Rows indicates an invalid row count (0 or 3-90).
	ErrPDF417Rows = errors.New("invalid PDF417 rows (try 0, or 3-90)")
	// ErrPDF417ModuleWidth indicates an invalid module width (2-8).
	ErrPDF417ModuleWidth = errors.New("invalid PDF417 module width (try 2-8)")
	// ErrPDF417RowHeight indicates an invalid row height (2-8).
	ErrPDF417RowHeight = errors.New("invalid PDF417 row height (try 2-8)")
	// ErrPDF417ErrorCorrection indicates an invalid error correction level (0-8).
	ErrPDF417ErrorCorrection = errors.New("invalid PDF417 error correction level (try 0-8)")
	// ErrPDF417DataEmpty indicates empty PDF417 data payload.
	ErrPDF417DataEmpty = errors.New("PDF417 data cannot be empty")
	// ErrPDF417DataTooLong indicates data payload exceeds 65532 bytes.
	ErrPDF417DataTooLong = errors.New("PDF417 data too long (maximum 65532 bytes)")
)

// SetPDF417Columns sets the number of columns in the PDF417 data area (fn = 65).
// Range: 0 <= columns <= 30 (0 sets automatic processing).
func SetPDF417Columns(columns byte) ([]byte, error) {
	if columns > MaxPDF417Columns {
		return nil, fmt.Errorf("%w: columns=%d", ErrPDF417Columns, columns)
	}
	return []byte{shared.GS, '(', 'k', 0x03, 0x00, 48, 65, columns}, nil
}

// SetPDF417Rows sets the number of rows in the PDF417 data area (fn = 66).
// Range: rows == 0 || (rows >= 3 && rows <= 90) (0 sets automatic processing).
func SetPDF417Rows(rows byte) ([]byte, error) {
	if rows != 0 && (rows < MinPDF417Rows || rows > MaxPDF417Rows) {
		return nil, fmt.Errorf("%w: rows=%d", ErrPDF417Rows, rows)
	}
	return []byte{shared.GS, '(', 'k', 0x03, 0x00, 48, 66, rows}, nil
}

// SetPDF417ModuleWidth sets the module width of the PDF417 symbol (fn = 67).
// Range: 2 <= width <= 8 dots.
func SetPDF417ModuleWidth(width byte) ([]byte, error) {
	if width < MinPDF417ModuleWidth || width > MaxPDF417ModuleWidth {
		return nil, fmt.Errorf("%w: width=%d", ErrPDF417ModuleWidth, width)
	}
	return []byte{shared.GS, '(', 'k', 0x03, 0x00, 48, 67, width}, nil
}

// SetPDF417RowHeight sets the row height of the PDF417 symbol (fn = 68).
// Range: 2 <= height <= 8 (height = module width × row height).
func SetPDF417RowHeight(height byte) ([]byte, error) {
	if height < MinPDF417RowHeight || height > MaxPDF417RowHeight {
		return nil, fmt.Errorf("%w: height=%d", ErrPDF417RowHeight, height)
	}
	return []byte{shared.GS, '(', 'k', 0x03, 0x00, 48, 68, height}, nil
}

// SetPDF417ErrorCorrection sets the error correction level for PDF417 (fn = 69, m = 48).
// Range: 0 <= level <= 8.
func SetPDF417ErrorCorrection(level byte) ([]byte, error) {
	if level > MaxPDF417ErrorCorrection {
		return nil, fmt.Errorf("%w: level=%d", ErrPDF417ErrorCorrection, level)
	}
	return []byte{shared.GS, '(', 'k', 0x04, 0x00, 48, 69, 48, 48 + level}, nil
}

// StorePDF417Data stores symbol data into the PDF417 symbol storage area (fn = 80, m = 48).
// Range: 1 <= len(data) <= 65532 bytes.
func StorePDF417Data(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, ErrPDF417DataEmpty
	}
	if len(data) > MaxPDF417DataLength {
		return nil, fmt.Errorf("%w: %d bytes", ErrPDF417DataTooLong, len(data))
	}

	length := len(data) + 3
	pL := byte(length & 0xFF)
	pH := byte((length >> 8) & 0xFF)

	cmd := make([]byte, 0, 8+len(data))
	cmd = append(cmd, shared.GS, '(', 'k', pL, pH, 48, 80, 48)
	cmd = append(cmd, data...)
	return cmd, nil
}

// PrintPDF417 prints the PDF417 symbol data stored in the symbol storage area (fn = 81, m = 48).
func PrintPDF417() []byte {
	return []byte{shared.GS, '(', 'k', 0x03, 0x00, 48, 81, 48}
}

// ComposePDF417 combines all setup commands, data storage, and the print command into a single atomic byte sequence.
func ComposePDF417(columns, rows, width, rowHeight, ecLevel byte, data []byte) ([]byte, error) {
	colCmd, err := SetPDF417Columns(columns)
	if err != nil {
		return nil, fmt.Errorf("pdf417 set columns: %w", err)
	}
	rowCmd, err := SetPDF417Rows(rows)
	if err != nil {
		return nil, fmt.Errorf("pdf417 set rows: %w", err)
	}
	widthCmd, err := SetPDF417ModuleWidth(width)
	if err != nil {
		return nil, fmt.Errorf("pdf417 set module width: %w", err)
	}
	heightCmd, err := SetPDF417RowHeight(rowHeight)
	if err != nil {
		return nil, fmt.Errorf("pdf417 set row height: %w", err)
	}
	ecCmd, err := SetPDF417ErrorCorrection(ecLevel)
	if err != nil {
		return nil, fmt.Errorf("pdf417 set error correction: %w", err)
	}
	storeCmd, err := StorePDF417Data(data)
	if err != nil {
		return nil, fmt.Errorf("pdf417 store data: %w", err)
	}
	printCmd := PrintPDF417()

	totalLen := len(colCmd) + len(rowCmd) + len(widthCmd) + len(heightCmd) + len(ecCmd) + len(storeCmd) + len(printCmd)
	buf := make([]byte, 0, totalLen)
	buf = append(buf, colCmd...)
	buf = append(buf, rowCmd...)
	buf = append(buf, widthCmd...)
	buf = append(buf, heightCmd...)
	buf = append(buf, ecCmd...)
	buf = append(buf, storeCmd...)
	buf = append(buf, printCmd...)
	return buf, nil
}
