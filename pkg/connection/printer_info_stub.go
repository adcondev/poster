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
	// StateReady indicates the printer is ready to accept jobs
	StateReady PrinterState = "ready"
	// StateOffline indicates the printer is offline
	StateOffline PrinterState = "offline"
	// StatePaused indicates the printer is paused
	StatePaused PrinterState = "paused"
	// StateError indicates the printer is in an error state
	StateError PrinterState = "error"
	// StateUnknown indicates the printer state is unknown
	StateUnknown PrinterState = "unknown"
)
