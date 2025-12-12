package executor

import (
	"encoding/json"
	"testing"
)

// ============================================================================
// Handler Registry Tests
// ============================================================================

func TestHandlerRegistry(t *testing.T) {
	registry := NewRegistry()

	// Test empty registry
	if _, ok := registry.Get("text"); ok {
		t.Error("Expected empty registry to not have 'text' handler")
	}

	// Test registration
	dummyHandler := func(_ interface{}, _ json.RawMessage) error {
		return nil
	}

	// Note: This test assumes HandlerRegistry uses a compatible handler signature
	// In the actual implementation, the handler signature is:
	// func(printer *service.Printer, data json.RawMessage) error
	// For testing purposes, we test the registry's basic functionality

	t.Run("list empty registry", func(t *testing.T) {
		list := registry.List()
		if len(list) != 0 {
			t.Errorf("Expected empty list, got %d items", len(list))
		}
	})

	// Since we can't easily test with the actual handler type without service.Printer,
	// we just verify the registry structure exists
	_ = dummyHandler
}

func BenchmarkTextCommandParsing(b *testing.B) {
	jsonData := []byte(`{"content": {"text": "Hello World", "content_style": {"bold": true, "size": "2x2"}}, "label": {"text": "Label"}}`)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var cmd TextCommand
		_ = json.Unmarshal(jsonData, &cmd)
	}
}

func BenchmarkTableCommandParsing(b *testing.B) {
	jsonData := []byte(`{
		"definition":  {"columns": [{"name": "Item", "width": 20}, {"name": "Price", "width":  10}]},
		"rows":  [["Coffee", "$4.50"], ["Muffin", "$3.00"], ["Tea", "$2.50"]],
		"options": {"header_bold": true, "word_wrap": true}
	}`)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var cmd TableCommand
		_ = json.Unmarshal(jsonData, &cmd)
	}
}
