package executor

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/adcondev/poster/pkg/commands/shared"
	"github.com/adcondev/poster/pkg/document/schema"
	"github.com/adcondev/poster/pkg/service"
)

func TestHandlePDF417_Success(t *testing.T) {
	mock := service.NewMockPrinter()
	exec := NewExecutor(mock)

	data := json.RawMessage(`{"data": "TEST-PDF417", "align": "center"}`)
	err := exec.handlePDF417(mock, data)
	if err != nil {
		t.Fatalf("handlePDF417 failed: %v", err)
	}

	// Verify alignment calls: AlignCenter followed by Write followed by AlignLeft
	calls := mock.Calls
	if len(calls) < 3 {
		t.Fatalf("expected at least 3 calls on mock printer, got %d", len(calls))
	}

	hasCenter := false
	hasWrite := false
	for _, call := range calls {
		if call.Method == "AlignCenter" {
			hasCenter = true
		}
		if call.Method == "Write" {
			hasWrite = true
			writtenBytes := call.Args[0].([]byte)
			if !bytes.Contains(writtenBytes, []byte("TEST-PDF417")) {
				t.Errorf("written bytes missing payload data")
			}
		}
	}

	if !hasCenter {
		t.Error("expected AlignCenter call")
	}
	if !hasWrite {
		t.Error("expected Write call")
	}

	lastCall := mock.LastCall()
	if lastCall == nil || lastCall.Method != "AlignLeft" {
		t.Errorf("expected last call to be AlignLeft, got %v", lastCall)
	}
}

func TestHandlePDF417_EmptyData(t *testing.T) {
	mock := service.NewMockPrinter()
	exec := NewExecutor(mock)

	data := json.RawMessage(`{"data": ""}`)
	err := exec.handlePDF417(mock, data)
	if err == nil {
		t.Fatal("expected error for empty PDF417 data, got nil")
	}
}

func TestHandlePDF417_InvalidJSON(t *testing.T) {
	mock := service.NewMockPrinter()
	exec := NewExecutor(mock)

	data := json.RawMessage(`{not-valid-json}`)
	err := exec.handlePDF417(mock, data)
	if err == nil {
		t.Fatal("expected error for invalid json, got nil")
	}
}

func TestHandlePDF417_DefaultsApplied(t *testing.T) {
	mock := service.NewMockPrinter()
	exec := NewExecutor(mock)

	data := json.RawMessage(`{"data": "DEFAULT-DATA"}`)
	err := exec.handlePDF417(mock, data)
	if err != nil {
		t.Fatalf("handlePDF417 failed: %v", err)
	}

	var writtenBytes []byte
	for _, call := range mock.Calls {
		if call.Method == "Write" {
			writtenBytes = call.Args[0].([]byte)
			break
		}
	}
	if writtenBytes == nil {
		t.Fatal("no Write call found on mock printer")
	}

	// Verify defaults:
	// columns: 0 -> []byte{shared.GS, '(', 'k', 0x03, 0x00, 48, 65, 0}
	wantCol := []byte{shared.GS, '(', 'k', 0x03, 0x00, 48, 65, 0}
	if !bytes.Contains(writtenBytes, wantCol) {
		t.Errorf("expected default columns cmd %v", wantCol)
	}

	// rows: 0 -> []byte{shared.GS, '(', 'k', 0x03, 0x00, 48, 66, 0}
	wantRows := []byte{shared.GS, '(', 'k', 0x03, 0x00, 48, 66, 0}
	if !bytes.Contains(writtenBytes, wantRows) {
		t.Errorf("expected default rows cmd %v", wantRows)
	}

	// moduleWidth: 3 -> []byte{shared.GS, '(', 'k', 0x03, 0x00, 48, 67, 3}
	wantWidth := []byte{shared.GS, '(', 'k', 0x03, 0x00, 48, 67, 3}
	if !bytes.Contains(writtenBytes, wantWidth) {
		t.Errorf("expected default moduleWidth cmd %v", wantWidth)
	}

	// rowHeight: 3 -> []byte{shared.GS, '(', 'k', 0x03, 0x00, 48, 68, 3}
	wantHeight := []byte{shared.GS, '(', 'k', 0x03, 0x00, 48, 68, 3}
	if !bytes.Contains(writtenBytes, wantHeight) {
		t.Errorf("expected default rowHeight cmd %v", wantHeight)
	}

	// ecLevel: 2 -> []byte{shared.GS, '(', 'k', 0x04, 0x00, 48, 69, 48, 50}
	wantEC := []byte{shared.GS, '(', 'k', 0x04, 0x00, 48, 69, 48, 50}
	if !bytes.Contains(writtenBytes, wantEC) {
		t.Errorf("expected default ecLevel cmd %v", wantEC)
	}
}

func TestHandlePDF417_CustomOptions(t *testing.T) {
	mock := service.NewMockPrinter()
	exec := NewExecutor(mock)

	data := json.RawMessage(`{
		"data": "CUSTOM-DATA",
		"columns": 5,
		"rows": 10,
		"module_width": 4,
		"row_height": 5,
		"error_correction": 4,
		"align": "right"
	}`)
	err := exec.handlePDF417(mock, data)
	if err != nil {
		t.Fatalf("handlePDF417 failed: %v", err)
	}

	hasRight := false
	var writtenBytes []byte
	for _, call := range mock.Calls {
		if call.Method == "AlignRight" {
			hasRight = true
		}
		if call.Method == "Write" {
			writtenBytes = call.Args[0].([]byte)
		}
	}
	if !hasRight {
		t.Error("expected AlignRight call")
	}
	if writtenBytes == nil {
		t.Fatal("expected Write call")
	}

	wantCol := []byte{shared.GS, '(', 'k', 0x03, 0x00, 48, 65, 5}
	if !bytes.Contains(writtenBytes, wantCol) {
		t.Errorf("expected custom columns cmd %v", wantCol)
	}

	wantRows := []byte{shared.GS, '(', 'k', 0x03, 0x00, 48, 66, 10}
	if !bytes.Contains(writtenBytes, wantRows) {
		t.Errorf("expected custom rows cmd %v", wantRows)
	}

	wantWidth := []byte{shared.GS, '(', 'k', 0x03, 0x00, 48, 67, 4}
	if !bytes.Contains(writtenBytes, wantWidth) {
		t.Errorf("expected custom module width cmd %v", wantWidth)
	}

	wantHeight := []byte{shared.GS, '(', 'k', 0x03, 0x00, 48, 68, 5}
	if !bytes.Contains(writtenBytes, wantHeight) {
		t.Errorf("expected custom row height cmd %v", wantHeight)
	}

	wantEC := []byte{shared.GS, '(', 'k', 0x04, 0x00, 48, 69, 48, 52}
	if !bytes.Contains(writtenBytes, wantEC) {
		t.Errorf("expected custom ecLevel cmd %v", wantEC)
	}
}

func TestHandlePDF417_AlignLeft(t *testing.T) {
	mock := service.NewMockPrinter()
	exec := NewExecutor(mock)

	data := json.RawMessage(`{"data": "LEFT-DATA", "align": "left"}`)
	err := exec.handlePDF417(mock, data)
	if err != nil {
		t.Fatalf("handlePDF417 failed: %v", err)
	}

	alignLeftCount := 0
	for _, call := range mock.Calls {
		if call.Method == "AlignLeft" {
			alignLeftCount++
		}
	}
	// Once before write, once to reset after write
	if alignLeftCount < 2 {
		t.Errorf("expected at least 2 AlignLeft calls, got %d", alignLeftCount)
	}
}

func TestExecutor_Execute_PDF417Document(t *testing.T) {
	mock := service.NewMockPrinter()
	exec := NewExecutor(mock)

	doc := &schema.Document{
		Version: "1.0",
		Profile: schema.ProfileConfig{
			Model: "POS-80",
		},
		Commands: []schema.Command{
			{
				Type: "pdf417",
				Data: json.RawMessage(`{"data": "END-TO-END-PDF417"}`),
			},
		},
	}

	if err := exec.Execute(doc); err != nil {
		t.Fatalf("exec.Execute failed: %v", err)
	}

	hasWrite := false
	for _, call := range mock.Calls {
		if call.Method == "Write" {
			writtenBytes := call.Args[0].([]byte)
			if bytes.Contains(writtenBytes, []byte("END-TO-END-PDF417")) {
				hasWrite = true
				break
			}
		}
	}
	if !hasWrite {
		t.Error("expected PDF417 command to be written during document execution")
	}
}
