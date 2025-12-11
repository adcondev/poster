package executor

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/adcondev/poster/pkg/constants"
	"github.com/adcondev/poster/pkg/service"
)

// RawCommand for raw handler
type RawCommand struct {
	Hex      string `json:"hex"`
	Format   string `json:"format,omitempty"`
	Comment  string `json:"comment,omitempty"`
	SafeMode bool   `json:"safe_mode,omitempty"`
}

// TODO: Identify advanced patterns as RegExp, for now Execute Macro: (1D 5E ** ** >01<)

var criticalCommands = []struct {
	pattern []byte
	name    string
	risk    string
}{
	{[]byte{0x1B, 0x40}, "ESC @", "Full reset - clears all settings"},
	{[]byte{0x1B, 0x3D, 0x00}, "ESC = 0", "Disable printer - requires manual intervention"},
	{[]byte{0x1B, 0x70}, "ESC p", "Cash drawer - physical hardware activation"},
	{[]byte{0x1D, 0x3A}, "GS :", "Start/End Macro - execution (1D 5E ** ** >01<) mode = 1 may cause motor overheating if idle"},
}

// Known commands that may cause unexpected responses
var bidirectionalCommands = []struct {
	pattern []byte
	name    string
}{
	{[]byte{0x10, 0x05, 0x04}, "DLE ENQ 4 - Real-time status"},
	{[]byte{0x1D, 0x49}, "GS I - Transmit printer ID"},
	{[]byte{0x1D, 0x72}, "GS r - Transmit status"},
	{[]byte{0x10, 0x04}, "DLE EOT - Real-time status"},
	{[]byte{0x1D, 0x61}, "GS a - Enable/disable automatic status"},
	{[]byte{0x1B, 0x75}, "ESC u - Transmit peripheral device status"},
	{[]byte{0x1B, 0x76}, "ESC v - Transmit paper sensor status"},
}

// handleRaw manages raw command execution
func (e *Executor) handleRaw(printer *service.Printer, data json.RawMessage) error {
	var cmd RawCommand
	if err := json.Unmarshal(data, &cmd); err != nil {
		return fmt.Errorf("failed to parse raw command: %w", err)
	}

	// Log the raw command attempt
	if cmd.Comment != "" {
		log.Printf("Raw command: %s", cmd.Comment)
	}

	// Determine format
	format := strings.ToLower(cmd.Format)
	if format == "" {
		format = constants.DefaultRawFormat.String() // default
	}

	// Parse bytes based on format
	var bytes []byte
	var err error

	switch format {
	case constants.Hex.String():
		bytes, err = ParseHexString(cmd.Hex)
		if err != nil {
			return fmt.Errorf("failed to parse hex string: %w", err)
		}
	// TODO: Verify if base64 handling works as intended
	case constants.Base64.String():
		bytes, err = base64.StdEncoding.DecodeString(cmd.Hex)
		if err != nil {
			return fmt.Errorf("failed to decode base64: %w", err)
		}
	default:
		return fmt.Errorf("unsupported format: %s (valid: hex, base64)", cmd.Format)
	}

	// Validate byte length
	if len(bytes) == 0 {
		return fmt.Errorf("raw command cannot be empty")
	}
	if len(bytes) > constants.RawMaxBytes { // 4KB limit
		return fmt.Errorf("raw command too large: %d bytes (max %d)", len(bytes), constants.RawMaxBytes)
	}

	// Safety checks if enabled (AFTER parsing bytes)
	if cmd.SafeMode {
		if warning := CheckCriticalCommands(bytes); warning != "" {
			// CRÍTICO: Bloquear ejecución en modo seguro
			err := fmt.Errorf("unsafe command blocked in safe mode: %s", warning)
			log.Printf("[SECURITY] %v", err)
			log.Printf("[SECURITY] To bypass, set safe_mode to false (AT YOUR OWN RISK)")
			return err
		}
		log.Printf("[SECURITY] Raw command passed safety checks")
	}

	// Advertencia sobre comandos bidireccionales (solo log, no bloquea)
	if ContainsBidirectionalCommand(bytes) {
		log.Printf("[WARNING] This raw command may expect a response.  Raw commands are write-only.")
		log.Printf("[WARNING] Any printer responses will remain unread and may cause issues.")
	}

	// Log hex representation for debugging
	if printer.Profile.DebugLog {
		log.Printf("Raw bytes to send: % 02X", bytes)
		log.Printf("Raw bytes length: %d", len(bytes))
	}

	// Send raw bytes directly to printer
	if err := printer.Write(bytes); err != nil {
		return fmt.Errorf("failed to write raw command: %w", err)
	}

	return nil
}
