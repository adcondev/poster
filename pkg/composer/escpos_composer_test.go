package composer

import (
	"bytes"
	"errors"
	"testing"
)

// ============================================================================
// Mocks
// ============================================================================

type mockPrintCapability struct {
	textFn                     func(string) ([]byte, error)
	printAndLineFeedFn         func() []byte
	printAndCarriageReturnFn   func() []byte
	formFeedFn                 func() []byte
	printAndFeedPaperFn        func(byte) []byte
	printAndFeedLinesFn        func(byte) []byte
	printAndReverseFeedFn      func(byte) ([]byte, error)
	printAndReverseFeedLinesFn func(byte) ([]byte, error)
	printDataInPageModeFn      func() []byte
	cancelDataFn               func() []byte
}

func (m *mockPrintCapability) Text(text string) ([]byte, error) {
	if m.textFn != nil {
		return m.textFn(text)
	}
	return nil, nil
}

func (m *mockPrintCapability) PrintAndLineFeed() []byte {
	if m.printAndLineFeedFn != nil {
		return m.printAndLineFeedFn()
	}
	return nil
}

func (m *mockPrintCapability) PrintAndCarriageReturn() []byte {
	if m.printAndCarriageReturnFn != nil {
		return m.printAndCarriageReturnFn()
	}
	return nil
}

func (m *mockPrintCapability) FormFeed() []byte {
	if m.formFeedFn != nil {
		return m.formFeedFn()
	}
	return nil
}

func (m *mockPrintCapability) PrintAndFeedPaper(units byte) []byte {
	if m.printAndFeedPaperFn != nil {
		return m.printAndFeedPaperFn(units)
	}
	return nil
}

func (m *mockPrintCapability) PrintAndFeedLines(lines byte) []byte {
	if m.printAndFeedLinesFn != nil {
		return m.printAndFeedLinesFn(lines)
	}
	return nil
}

func (m *mockPrintCapability) PrintAndReverseFeed(units byte) ([]byte, error) {
	if m.printAndReverseFeedFn != nil {
		return m.printAndReverseFeedFn(units)
	}
	return nil, nil
}

func (m *mockPrintCapability) PrintAndReverseFeedLines(lines byte) ([]byte, error) {
	if m.printAndReverseFeedLinesFn != nil {
		return m.printAndReverseFeedLinesFn(lines)
	}
	return nil, nil
}

func (m *mockPrintCapability) PrintDataInPageMode() []byte {
	if m.printDataInPageModeFn != nil {
		return m.printDataInPageModeFn()
	}
	return nil
}

func (m *mockPrintCapability) CancelData() []byte {
	if m.cancelDataFn != nil {
		return m.cancelDataFn()
	}
	return nil
}

// ============================================================================
// Tests
// ============================================================================

func TestEscposProtocol_PrintLn(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		mockText      func(string) ([]byte, error)
		mockPrintFeed func() []byte
		want          []byte
		wantErr       bool
	}{
		{
			name:  "success",
			input: "Hello World",
			mockText: func(s string) ([]byte, error) {
				return []byte(s), nil
			},
			mockPrintFeed: func() []byte {
				return []byte{0x0A}
			},
			want:    []byte("Hello World\n"),
			wantErr: false,
		},
		{
			name:  "text error",
			input: "Error",
			mockText: func(s string) ([]byte, error) {
				return nil, errors.New("mock error")
			},
			mockPrintFeed: func() []byte {
				return []byte{0x0A}
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			mock := &mockPrintCapability{
				textFn:             tt.mockText,
				printAndLineFeedFn: tt.mockPrintFeed,
			}

			// Create protocol with mock
			c := &EscposProtocol{
				Print: mock,
			}

			got, err := c.PrintLn(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("PrintLn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !bytes.Equal(got, tt.want) {
				t.Errorf("PrintLn() = %v, want %v", got, tt.want)
			}
		})
	}
}
