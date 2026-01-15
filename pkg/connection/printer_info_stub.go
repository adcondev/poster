//go:build !windows

package connection

// PrinterDetail stub for non-Windows platforms
type PrinterDetail struct {
	Name        string `json:"name"`
	Port        string `json:"port"`
	Driver      string `json:"driver"`
	Status      string `json:"status"`
	StatusRaw   uint32 `json:"status_raw,omitempty"`
	IsDefault   bool   `json:"is_default"`
	IsVirtual   bool   `json:"is_virtual"`
	PrinterType string `json:"printer_type"`
}

// PrinterState stub
type PrinterState = string

const (
	StateReady   PrinterState = "ready"
	StateOffline PrinterState = "offline"
	StatePaused  PrinterState = "paused"
	StateError   PrinterState = "error"
	StateUnknown PrinterState = "unknown"
)
