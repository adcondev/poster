package executor

import (
	"encoding/hex"
	"fmt"
	"strings"
)

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

// CheckCriticalCommands checks for known risky command sequences
func CheckCriticalCommands(bytes []byte) string {
	for _, danger := range criticalCommands {
		if ContainsSequence(bytes, danger.pattern) {
			return fmt.Sprintf("%s - %s", danger.name, danger.risk)
		}
	}
	return ""
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
