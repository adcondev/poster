package connection

// PrinterService defines the interface for OS-specific printer operations.
// This allows mocking the low-level printer API for unit testing.
type PrinterService interface {
	Open(name string) (uintptr, error)
	Close(handle uintptr) error
	StartDoc(handle uintptr, docName, dataType string) (uint32, error)
	EndDoc(handle uintptr) error
	AbortDoc(handle uintptr) error
	Write(handle uintptr, data []byte) (uint32, error)
}
