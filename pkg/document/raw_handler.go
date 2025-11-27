package document

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/adcondev/pos-printer/pkg/service"
)

// Known dangerous command patterns (for safe mode)
var dangerousPatterns = []struct {
	pattern []byte
	name    string
	risk    string
}{
	{[]byte{0x1B, 0x40}, "ESC @", "Full reset - clears all settings"},
	{[]byte{0x1B, 0x3D, 0x00}, "ESC = 0", "Disable printer - requires manual intervention"},
	{[]byte{0x10, 0x05, 0x04}, "DLE ENQ 4", "Real-time status - may cause unexpected responses"},
	{[]byte{0x1B, 0x70}, "ESC p", "Cash drawer - physical hardware activation"},
}

// Known commands that expect responses
var bidirectionalCommands = []struct {
	pattern []byte
	name    string
}{
	{[]byte{0x1D, 0x49}, "GS I - Transmit printer ID"},
	{[]byte{0x1D, 0x72}, "GS r - Transmit status"},
	{[]byte{0x10, 0x04}, "DLE EOT - Real-time status"},
	{[]byte{0x1D, 0x61}, "GS a - Enable/disable automatic status"},
	{[]byte{0x1B, 0x75}, "ESC u - Transmit peripheral device status"},
	{[]byte{0x1B, 0x76}, "ESC v - Transmit paper sensor status"},
}

// RawCommand represents a raw ESC/POS command
type RawCommand struct {
	Hex      string `json:"hex"`                 // Hex string
	Format   string `json:"format,omitempty"`    // "hex" or "base64"
	Comment  string `json:"comment,omitempty"`   // Documentation
	SafeMode bool   `json:"safe_mode,omitempty"` // Enable safety checks
}

// HandleRaw manages raw command execution
func (e *Executor) HandleRaw(printer *service.Printer, data json.RawMessage) error {
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
		format = "hex" // default
	}

	// Parse bytes based on format
	var bytes []byte
	var err error

	switch format {
	case "hex":
		bytes, err = ParseHexString(cmd.Hex)
		if err != nil {
			return fmt.Errorf("failed to parse hex string: %w", err)
		}
	case "base64":
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
	if len(bytes) > 4096 { // 4KB limit
		return fmt.Errorf("raw command too large: %d bytes (max 4096)", len(bytes))
	}

	// Safety checks if enabled (AFTER parsing bytes)
	if cmd.SafeMode {
		if warning := CheckDangerousCommands(bytes); warning != "" {
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

// ParseHexString converts various hex string formats to bytes
func ParseHexString(hexStr string) ([]byte, error) {
	// Clean the string: remove spaces, commas, 0x prefixes
	cleaned := CleanHexString(hexStr)

	// Validate not empty
	if len(cleaned) == 0 {
		return nil, fmt.Errorf("hex string is empty or contains no valid hex characters")
	}

	// Ensure even number of characters
	if len(cleaned)%2 != 0 {
		return nil, fmt.Errorf("hex string must have even number of characters")
	}

	// Decode hex
	bytes, err := hex.DecodeString(cleaned)
	if err != nil {
		return nil, fmt.Errorf("invalid hex string: %w", err)
	}

	return bytes, nil
}

// CleanHexString optimizado para manejar múltiples formatos de entrada
func CleanHexString(s string) string {
	if len(s) == 0 {
		return ""
	}

	// Pre-allocate con tamaño estimado
	var sb strings.Builder
	sb.Grow(len(s))

	s = strings.ToUpper(s)

	// Single-pass filtering con manejo de prefijos 0X
	i := 0
	for i < len(s) {
		c := s[i]

		// Detectar y saltar prefijo "0X"
		if c == '0' && i+1 < len(s) && s[i+1] == 'X' {
			i += 2 // Saltar "0X"
			continue
		}

		// Solo aceptar hex válido (0-9, A-F)
		if (c >= '0' && c <= '9') || (c >= 'A' && c <= 'F') {
			sb.WriteByte(c)
		}
		// Ignorar silenciosamente espacios, comas, colones, guiones, etc.

		i++
	}

	return sb.String()
}

// CheckDangerousCommands checks for known risky command sequences
func CheckDangerousCommands(bytes []byte) string {
	for _, danger := range dangerousPatterns {
		if ContainsSequence(bytes, danger.pattern) {
			return fmt.Sprintf("%s - %s", danger.name, danger.risk)
		}
	}
	return ""
}

// ContainsSequence checks if bytes contains the pattern
func ContainsSequence(bytes, pattern []byte) bool {
	if len(pattern) == 0 || len(bytes) < len(pattern) {
		return false
	}

	for i := 0; i <= len(bytes)-len(pattern); i++ {
		match := true
		for j := 0; j < len(pattern); j++ {
			if bytes[i+j] != pattern[j] {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}

// ContainsBidirectionalCommand checks for commands that expect responses
func ContainsBidirectionalCommand(bytes []byte) bool {
	for _, cmd := range bidirectionalCommands {
		if ContainsSequence(bytes, cmd.pattern) {
			return true
		}
	}
	return false
}
