package schema

import (
	"encoding/json"
	"testing"
)

func TestDocument_Validate(t *testing.T) {
	tests := []struct {
		name    string
		doc     Document
		wantErr bool
		errMsg  string
	}{
		// Version format tests
		{
			name: "valid version format 1.0",
			doc: Document{
				Version:  "1.0",
				Profile:  ProfileConfig{Model: "TestPrinter"},
				Commands: []Command{{Type: "text", Data: json.RawMessage(`{}`)}},
			},
			wantErr: false,
		},
		{
			name: "valid version format 2.1",
			doc: Document{
				Version:  "2.1",
				Profile:  ProfileConfig{Model: "TestPrinter"},
				Commands: []Command{{Type: "text", Data: json.RawMessage(`{}`)}},
			},
			wantErr: false,
		},
		{
			name: "valid version format 10.25",
			doc: Document{
				Version:  "10.25",
				Profile:  ProfileConfig{Model: "TestPrinter"},
				Commands: []Command{{Type: "text", Data: json.RawMessage(`{}`)}},
			},
			wantErr: false,
		},
		{
			name: "invalid version format - missing minor",
			doc: Document{
				Version:  "1",
				Profile:  ProfileConfig{Model: "TestPrinter"},
				Commands: []Command{{Type: "text", Data: json.RawMessage(`{}`)}},
			},
			wantErr: true,
			errMsg:  "invalid version format",
		},
		{
			name: "invalid version format - three parts",
			doc: Document{
				Version:  "1.0.0",
				Profile:  ProfileConfig{Model: "TestPrinter"},
				Commands: []Command{{Type: "text", Data: json.RawMessage(`{}`)}},
			},
			wantErr: true,
			errMsg:  "invalid version format",
		},
		{
			name: "invalid version format - letters",
			doc: Document{
				Version:  "v1.0",
				Profile:  ProfileConfig{Model: "TestPrinter"},
				Commands: []Command{{Type: "text", Data: json.RawMessage(`{}`)}},
			},
			wantErr: true,
			errMsg:  "invalid version format",
		},
		{
			name: "invalid version format - empty",
			doc: Document{
				Version:  "",
				Profile:  ProfileConfig{Model: "TestPrinter"},
				Commands: []Command{{Type: "text", Data: json.RawMessage(`{}`)}},
			},
			wantErr: true,
			errMsg:  "version is required",
		},

		// Profile.model tests
		{
			name: "missing profile.model",
			doc: Document{
				Version:  "1.0",
				Profile:  ProfileConfig{},
				Commands: []Command{{Type: "text", Data: json.RawMessage(`{}`)}},
			},
			wantErr: true,
			errMsg:  "profile.model is required",
		},
		{
			name: "valid profile.model",
			doc: Document{
				Version:  "1.0",
				Profile:  ProfileConfig{Model: "EPSON TM-T20II"},
				Commands: []Command{{Type: "text", Data: json.RawMessage(`{}`)}},
			},
			wantErr: false,
		},

		// Commands array tests
		{
			name: "empty commands array",
			doc: Document{
				Version:  "1.0",
				Profile:  ProfileConfig{Model: "TestPrinter"},
				Commands: []Command{},
			},
			wantErr: true,
			errMsg:  "document must contain at least one command",
		},
		{
			name: "nil commands array",
			doc: Document{
				Version:  "1.0",
				Profile:  ProfileConfig{Model: "TestPrinter"},
				Commands: nil,
			},
			wantErr: true,
			errMsg:  "document must contain at least one command",
		},

		// Paper width tests
		{
			name: "valid paper_width 58",
			doc: Document{
				Version:  "1.0",
				Profile:  ProfileConfig{Model: "TestPrinter", PaperWidth: 58},
				Commands: []Command{{Type: "text", Data: json.RawMessage(`{}`)}},
			},
			wantErr: false,
		},
		{
			name: "valid paper_width 72",
			doc: Document{
				Version:  "1.0",
				Profile:  ProfileConfig{Model: "TestPrinter", PaperWidth: 72},
				Commands: []Command{{Type: "text", Data: json.RawMessage(`{}`)}},
			},
			wantErr: false,
		},
		{
			name: "valid paper_width 80",
			doc: Document{
				Version:  "1.0",
				Profile:  ProfileConfig{Model: "TestPrinter", PaperWidth: 80},
				Commands: []Command{{Type: "text", Data: json.RawMessage(`{}`)}},
			},
			wantErr: false,
		},
		{
			name: "valid paper_width 100",
			doc: Document{
				Version:  "1.0",
				Profile:  ProfileConfig{Model: "TestPrinter", PaperWidth: 100},
				Commands: []Command{{Type: "text", Data: json.RawMessage(`{}`)}},
			},
			wantErr: false,
		},
		{
			name: "valid paper_width 112",
			doc: Document{
				Version:  "1.0",
				Profile:  ProfileConfig{Model: "TestPrinter", PaperWidth: 112},
				Commands: []Command{{Type: "text", Data: json.RawMessage(`{}`)}},
			},
			wantErr: false,
		},
		{
			name: "valid paper_width 120",
			doc: Document{
				Version:  "1.0",
				Profile:  ProfileConfig{Model: "TestPrinter", PaperWidth: 120},
				Commands: []Command{{Type: "text", Data: json.RawMessage(`{}`)}},
			},
			wantErr: false,
		},
		{
			name: "valid paper_width 0 (default)",
			doc: Document{
				Version:  "1.0",
				Profile:  ProfileConfig{Model: "TestPrinter", PaperWidth: 0},
				Commands: []Command{{Type: "text", Data: json.RawMessage(`{}`)}},
			},
			wantErr: false,
		},
		{
			name: "invalid paper_width 50",
			doc: Document{
				Version:  "1.0",
				Profile:  ProfileConfig{Model: "TestPrinter", PaperWidth: 50},
				Commands: []Command{{Type: "text", Data: json.RawMessage(`{}`)}},
			},
			wantErr: true,
			errMsg:  "invalid paper_width",
		},
		{
			name: "invalid paper_width 60",
			doc: Document{
				Version:  "1.0",
				Profile:  ProfileConfig{Model: "TestPrinter", PaperWidth: 60},
				Commands: []Command{{Type: "text", Data: json.RawMessage(`{}`)}},
			},
			wantErr: true,
			errMsg:  "invalid paper_width",
		},
		{
			name: "invalid paper_width 90",
			doc: Document{
				Version:  "1.0",
				Profile:  ProfileConfig{Model: "TestPrinter", PaperWidth: 90},
				Commands: []Command{{Type: "text", Data: json.RawMessage(`{}`)}},
			},
			wantErr: true,
			errMsg:  "invalid paper_width",
		},

		// DPI tests
		{
			name: "valid DPI 203",
			doc: Document{
				Version:  "1.0",
				Profile:  ProfileConfig{Model: "TestPrinter", DPI: 203},
				Commands: []Command{{Type: "text", Data: json.RawMessage(`{}`)}},
			},
			wantErr: false,
		},
		{
			name: "valid DPI 300",
			doc: Document{
				Version:  "1.0",
				Profile:  ProfileConfig{Model: "TestPrinter", DPI: 300},
				Commands: []Command{{Type: "text", Data: json.RawMessage(`{}`)}},
			},
			wantErr: false,
		},
		{
			name: "valid DPI 600",
			doc: Document{
				Version:  "1.0",
				Profile:  ProfileConfig{Model: "TestPrinter", DPI: 600},
				Commands: []Command{{Type: "text", Data: json.RawMessage(`{}`)}},
			},
			wantErr: false,
		},
		{
			name: "valid DPI 0 (default)",
			doc: Document{
				Version:  "1.0",
				Profile:  ProfileConfig{Model: "TestPrinter", DPI: 0},
				Commands: []Command{{Type: "text", Data: json.RawMessage(`{}`)}},
			},
			wantErr: false,
		},
		{
			name: "invalid DPI 150",
			doc: Document{
				Version:  "1.0",
				Profile:  ProfileConfig{Model: "TestPrinter", DPI: 150},
				Commands: []Command{{Type: "text", Data: json.RawMessage(`{}`)}},
			},
			wantErr: true,
			errMsg:  "invalid dpi",
		},
		{
			name: "invalid DPI 400",
			doc: Document{
				Version:  "1.0",
				Profile:  ProfileConfig{Model: "TestPrinter", DPI: 400},
				Commands: []Command{{Type: "text", Data: json.RawMessage(`{}`)}},
			},
			wantErr: true,
			errMsg:  "invalid dpi",
		},
		{
			name: "invalid DPI 72",
			doc: Document{
				Version:  "1.0",
				Profile:  ProfileConfig{Model: "TestPrinter", DPI: 72},
				Commands: []Command{{Type: "text", Data: json.RawMessage(`{}`)}},
			},
			wantErr: true,
			errMsg:  "invalid dpi",
		},

		// Combined valid document
		{
			name: "fully valid document with all fields",
			doc: Document{
				Version: "1.0",
				Profile: ProfileConfig{
					Model:      "EPSON TM-T88VI",
					PaperWidth: 80,
					CodeTable:  "WPC1252",
					DPI:        203,
					HasQR:      true,
				},
				DebugLog: true,
				Commands: []Command{
					{Type: "text", Data: json.RawMessage(`{"content":{"text":"Hello"}}`)},
					{Type: "cut", Data: json.RawMessage(`{"mode":"partial"}`)},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.doc.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Document.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && tt.errMsg != "" {
				if !contains(err.Error(), tt.errMsg) {
					t.Errorf("Document.Validate() error = %v, expected to contain %q", err, tt.errMsg)
				}
			}
		})
	}
}

// contains checks if s contains substr
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestDocument_Validate_MultipleErrors(t *testing.T) {
	// Test that validation returns the first error encountered
	doc := Document{
		Version:  "invalid",
		Profile:  ProfileConfig{}, // Missing model
		Commands: nil,             // Empty commands
	}

	err := doc.Validate()
	if err == nil {
		t.Error("Expected error for invalid document, got nil")
	}

	// Should fail on version first
	if err != nil && !contains(err.Error(), "version") {
		t.Errorf("Expected version error first, got: %v", err)
	}
}

func BenchmarkDocument_Validate(b *testing.B) {
	doc := Document{
		Version: "1.0",
		Profile: ProfileConfig{
			Model:      "TestPrinter",
			PaperWidth: 80,
			DPI:        203,
		},
		Commands: []Command{
			{Type: "text", Data: json.RawMessage(`{}`)},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = doc.Validate()
	}
}
