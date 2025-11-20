package connection

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockPrinterService is a mock implementation of PrinterService
type MockPrinterService struct {
	mock.Mock
}

func (m *MockPrinterService) Open(name string) (uintptr, error) {
	args := m.Called(name)
	return args.Get(0).(uintptr), args.Error(1)
}

func (m *MockPrinterService) Close(handle uintptr) error {
	args := m.Called(handle)
	return args.Error(0)
}

func (m *MockPrinterService) StartDoc(handle uintptr, docName, dataType string) (uint32, error) {
	args := m.Called(handle, docName, dataType)
	return args.Get(0).(uint32), args.Error(1)
}

func (m *MockPrinterService) EndDoc(handle uintptr) error {
	args := m.Called(handle)
	return args.Error(0)
}

func (m *MockPrinterService) AbortDoc(handle uintptr) error {
	args := m.Called(handle)
	return args.Error(0)
}

func (m *MockPrinterService) Write(handle uintptr, data []byte) (uint32, error) {
	args := m.Called(handle, data)
	return args.Get(0).(uint32), args.Error(1)
}

// NewWindowsPrintConnectorWithService allows injecting a mock service for testing.
// This function is internal to the package (or exposed for testing if necessary).
// Since we are in the same package, we can just create the struct directly or modify NewWindowsPrintConnector to accept options.
// But NewWindowsPrintConnector uses getPlatformPrinterService().
// For testing, we want to bypass that.
func newTestWindowsPrintConnector(name string, service PrinterService) (*WindowsPrintConnector, error) {
	if name == "" {
		return nil, errors.New("el nombre de la impresora no puede estar vacío")
	}
	handle, err := service.Open(name)
	if err != nil {
		return nil, err
	}
	return &WindowsPrintConnector{
		printerName: name,
		service:     service,
		handle:      handle,
		jobStarted:  false,
	}, nil
}

func TestNewWindowsPrintConnector_Success(t *testing.T) {
	mockService := new(MockPrinterService)
	mockService.On("Open", "TestPrinter").Return(uintptr(123), nil)

	connector, err := newTestWindowsPrintConnector("TestPrinter", mockService)

	assert.NoError(t, err)
	assert.NotNil(t, connector)
	assert.Equal(t, uintptr(123), connector.handle)
	mockService.AssertExpectations(t)
}

func TestNewWindowsPrintConnector_EmptyName(t *testing.T) {
	mockService := new(MockPrinterService)
	connector, err := newTestWindowsPrintConnector("", mockService)

	assert.Error(t, err)
	assert.Nil(t, connector)
	assert.Equal(t, "el nombre de la impresora no puede estar vacío", err.Error())
}

func TestNewWindowsPrintConnector_OpenError(t *testing.T) {
	mockService := new(MockPrinterService)
	mockService.On("Open", "TestPrinter").Return(uintptr(0), errors.New("open failed"))

	connector, err := newTestWindowsPrintConnector("TestPrinter", mockService)

	assert.Error(t, err)
	assert.Nil(t, connector)
	mockService.AssertExpectations(t)
}

func TestWindowsPrintConnector_Write_Success(t *testing.T) {
	mockService := new(MockPrinterService)
	mockService.On("Open", "TestPrinter").Return(uintptr(123), nil)
	mockService.On("StartDoc", uintptr(123), "ESC/POS PrintDataInPageMode Job", "RAW").Return(uint32(1), nil)
	data := []byte("test data")
	mockService.On("Write", uintptr(123), data).Return(uint32(len(data)), nil)

	connector, _ := newTestWindowsPrintConnector("TestPrinter", mockService)

	n, err := connector.Write(data)

	assert.NoError(t, err)
	assert.Equal(t, len(data), n)
	assert.True(t, connector.jobStarted)
	mockService.AssertExpectations(t)
}

func TestWindowsPrintConnector_Write_StartDocError(t *testing.T) {
	mockService := new(MockPrinterService)
	mockService.On("Open", "TestPrinter").Return(uintptr(123), nil)
	mockService.On("StartDoc", uintptr(123), mock.Anything, mock.Anything).Return(uint32(0), errors.New("startdoc failed"))

	connector, _ := newTestWindowsPrintConnector("TestPrinter", mockService)

	data := []byte("test data")
	n, err := connector.Write(data)

	assert.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, connector.jobStarted)
	mockService.AssertExpectations(t)
}

func TestWindowsPrintConnector_Write_WriteError(t *testing.T) {
	mockService := new(MockPrinterService)
	mockService.On("Open", "TestPrinter").Return(uintptr(123), nil)
	mockService.On("StartDoc", uintptr(123), mock.Anything, mock.Anything).Return(uint32(1), nil)
	data := []byte("test data")
	mockService.On("Write", uintptr(123), data).Return(uint32(0), errors.New("write failed"))

	connector, _ := newTestWindowsPrintConnector("TestPrinter", mockService)

	n, err := connector.Write(data)

	assert.Error(t, err)
	assert.Equal(t, 0, n)
	mockService.AssertExpectations(t)
}

func TestWindowsPrintConnector_Close_Success(t *testing.T) {
	mockService := new(MockPrinterService)
	mockService.On("Open", "TestPrinter").Return(uintptr(123), nil)
	mockService.On("StartDoc", uintptr(123), mock.Anything, mock.Anything).Return(uint32(1), nil)
	mockService.On("Write", uintptr(123), mock.Anything).Return(uint32(4), nil)
	mockService.On("EndDoc", uintptr(123)).Return(nil)
	mockService.On("Close", uintptr(123)).Return(nil)

	connector, _ := newTestWindowsPrintConnector("TestPrinter", mockService)
	connector.Write([]byte("test")) // Start job

	err := connector.Close()

	assert.NoError(t, err)
	assert.False(t, connector.jobStarted)
	assert.Equal(t, uintptr(0), connector.handle)
	mockService.AssertExpectations(t)
}

func TestWindowsPrintConnector_Close_WithoutJob(t *testing.T) {
	mockService := new(MockPrinterService)
	mockService.On("Open", "TestPrinter").Return(uintptr(123), nil)
	mockService.On("Close", uintptr(123)).Return(nil)

	connector, _ := newTestWindowsPrintConnector("TestPrinter", mockService)

	err := connector.Close()

	assert.NoError(t, err)
	assert.Equal(t, uintptr(0), connector.handle)
	mockService.AssertExpectations(t)
}

func TestWindowsPrintConnector_Close_EndDocFail_AbortSucceeds(t *testing.T) {
	mockService := new(MockPrinterService)
	mockService.On("Open", "TestPrinter").Return(uintptr(123), nil)
	mockService.On("StartDoc", uintptr(123), mock.Anything, mock.Anything).Return(uint32(1), nil)
	mockService.On("Write", uintptr(123), mock.Anything).Return(uint32(4), nil)

	mockService.On("EndDoc", uintptr(123)).Return(errors.New("EndDoc failed"))
	mockService.On("AbortDoc", uintptr(123)).Return(nil)
	mockService.On("Close", uintptr(123)).Return(nil)

	connector, _ := newTestWindowsPrintConnector("TestPrinter", mockService)
	connector.Write([]byte("test"))

	err := connector.Close()

	assert.NoError(t, err) // Should be nil as abort succeeded?
	// Code says: if abortErr != nil { finalErr = ... }
	// So if abort succeeds, finalErr is nil (unless close fails)
	mockService.AssertExpectations(t)
}

func TestWindowsPrintConnector_Close_EndDocFail_AbortFails(t *testing.T) {
	mockService := new(MockPrinterService)
	mockService.On("Open", "TestPrinter").Return(uintptr(123), nil)
	mockService.On("StartDoc", uintptr(123), mock.Anything, mock.Anything).Return(uint32(1), nil)
	mockService.On("Write", uintptr(123), mock.Anything).Return(uint32(4), nil)

	mockService.On("EndDoc", uintptr(123)).Return(errors.New("EndDoc failed"))
	mockService.On("AbortDoc", uintptr(123)).Return(errors.New("AbortDoc failed"))
	mockService.On("Close", uintptr(123)).Return(nil)

	connector, _ := newTestWindowsPrintConnector("TestPrinter", mockService)
	connector.Write([]byte("test"))

	err := connector.Close()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "AbortDoc")
	mockService.AssertExpectations(t)
}

func TestWindowsPrintConnector_Read(t *testing.T) {
	mockService := new(MockPrinterService)
	mockService.On("Open", "TestPrinter").Return(uintptr(123), nil)
	connector, _ := newTestWindowsPrintConnector("TestPrinter", mockService)

	n, err := connector.Read([]byte{})
	assert.Error(t, err)
	assert.Equal(t, 0, n)
}
