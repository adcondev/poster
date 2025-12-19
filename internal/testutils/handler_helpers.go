package testutils

// MockConnector para testing
type MockConnector struct {
	WriteFunc func(data []byte) (int, error)
	CloseFunc func() error
}

// Write escribe datos al conector simulado
func (m *MockConnector) Write(data []byte) (int, error) {
	if m.WriteFunc != nil {
		return m.WriteFunc(data)
	}
	return len(data), nil
}

// Close cierra el conector simulado
func (m *MockConnector) Close() error {
	if m.CloseFunc != nil {
		return m.CloseFunc()
	}
	return nil
}

// Read lee datos del conector simulado
func (m *MockConnector) Read(_ []byte) (int, error) {
	return 0, nil
}
