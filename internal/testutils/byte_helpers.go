package testutils

// Create a test helper for buffer management

// BufferBuilder helps in building byte buffers for tests
type BufferBuilder struct {
	buffer []byte
}

// Append adds bytes to the buffer
func (b *BufferBuilder) Append(cmd []byte) *BufferBuilder {
	b.buffer = append(b.buffer, cmd...)
	return b
}

// GetBuffer returns the constructed byte buffer
func (b *BufferBuilder) GetBuffer() []byte {
	return b.buffer
}

// ============================================================================
// Byte Slice Generators and Manipulators
// ============================================================================

// RepeatByte creates a byte slice of specified length filled with value
func RepeatByte(length int, value byte) []byte {
	result := make([]byte, length)
	for i := range result {
		result[i] = value
	}
	return result
}

// GenerateString creates a string of specified length filled with value
func GenerateString(length int, value byte) string {
	return string(RepeatByte(length, value))
}

// ============================================================================
// Command Builders
// ============================================================================

// BuildCommand constructs a command with variable parameters
func BuildCommand(cmd byte, subcmd byte, params ...byte) []byte {
	result := []byte{cmd, subcmd}
	return append(result, params...)
}
