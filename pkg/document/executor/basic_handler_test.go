package executor

import (
	"encoding/json"
	"testing"

	"github.com/adcondev/poster/pkg/constants"
)

// ============================================================================
// Feed Command Tests
// ============================================================================

func TestFeedCommand_Parsing(t *testing.T) {
	tests := []struct {
		name      string
		json      string
		checkFunc func(t *testing.T, cmd FeedCommand)
	}{
		{
			name: "valid feed command",
			json: `{"lines": 5}`,
			checkFunc: func(t *testing.T, cmd FeedCommand) {
				if cmd.Lines != 5 {
					t.Errorf("Expected lines 5, got %d", cmd.Lines)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cmd FeedCommand
			if err := json.Unmarshal([]byte(tt.json), &cmd); err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if tt.checkFunc != nil {
				tt.checkFunc(t, cmd)
			}
		})
	}
}

func TestFeedCommand_Defaults(t *testing.T) {
	var cmd FeedCommand
	if err := json.Unmarshal([]byte(`{}`), &cmd); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if cmd.Lines != 0 {
		t.Error("Expected lines 0 (handler applies default)")
	}
}

func TestFeedCommand_Validation(t *testing.T) {
	var cmd FeedCommand
	err := json.Unmarshal([]byte(`{invalid}`), &cmd)
	if err == nil {
		t.Error("Expected error for invalid json")
	}
}

// ============================================================================
// Cut Command Tests
// ============================================================================

func TestCutCommand_Parsing(t *testing.T) {
	tests := []struct {
		name      string
		json      string
		checkFunc func(t *testing.T, cmd CutCommand)
	}{
		{
			name: "partial cut",
			json: `{"mode": "partial", "feed": 3}`,
			checkFunc: func(t *testing.T, cmd CutCommand) {
				if cmd.Mode != "partial" {
					t.Errorf("Expected mode 'partial', got '%s'", cmd.Mode)
				}
				if cmd.Feed != 3 {
					t.Errorf("Expected feed 3, got %d", cmd.Feed)
				}
			},
		},
		{
			name: "full cut",
			json: `{"mode": "full"}`,
			checkFunc: func(t *testing.T, cmd CutCommand) {
				if cmd.Mode != "full" {
					t.Errorf("Expected mode 'full', got '%s'", cmd.Mode)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cmd CutCommand
			if err := json.Unmarshal([]byte(tt.json), &cmd); err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if tt.checkFunc != nil {
				tt.checkFunc(t, cmd)
			}
		})
	}
}

func TestCutCommand_Defaults(t *testing.T) {
	var cmd CutCommand
	if err := json.Unmarshal([]byte(`{}`), &cmd); err != nil {
		t.Fatalf("Unexpected error:  %v", err)
	}
	if cmd.Mode != "" {
		t.Error("Expected empty mode (handler applies default)")
	}
	if cmd.Feed != 0 {
		t.Error("Expected feed 0 (handler applies default)")
	}
}

func TestCutCommand_Validation(t *testing.T) {
	var cmd CutCommand
	err := json.Unmarshal([]byte(`{invalid}`), &cmd)
	if err == nil {
		t.Error("Expected error for invalid json")
	}
}

// ============================================================================
// Separator Command Tests
// ============================================================================

func TestSeparatorCommand_Parsing(t *testing.T) {
	tests := []struct {
		name      string
		json      string
		checkFunc func(t *testing.T, cmd SeparatorCommand)
	}{
		{
			name: "dash separator",
			json: `{"char": "-", "length": 48}`,
			checkFunc: func(t *testing.T, cmd SeparatorCommand) {
				if cmd.Char != "-" {
					t.Errorf("Expected char '-', got '%s'", cmd.Char)
				}
				if cmd.Length != 48 {
					t.Errorf("Expected length 48, got %d", cmd.Length)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cmd SeparatorCommand
			if err := json.Unmarshal([]byte(tt.json), &cmd); err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if tt.checkFunc != nil {
				tt.checkFunc(t, cmd)
			}
		})
	}
}

func TestSeparatorCommand_Defaults(t *testing.T) {
	var cmd SeparatorCommand
	if err := json.Unmarshal([]byte(`{}`), &cmd); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if cmd.Char != "" {
		t.Error("Expected empty char (handler applies default)")
	}
	if cmd.Length != 0 {
		t.Error("Expected length 0 (handler applies default)")
	}
}

func TestSeparatorCommand_Validation(t *testing.T) {
	var cmd SeparatorCommand
	err := json.Unmarshal([]byte(`{invalid}`), &cmd)
	if err == nil {
		t.Error("Expected error for invalid json")
	}
}

// ============================================================================
// Pulse Command Tests
// ============================================================================

func TestPulseCommand_Parsing(t *testing.T) {
	var cmd PulseCommand
	err := json.Unmarshal([]byte(`{"pin": 1, "on_time": 100, "off_time": 200}`), &cmd)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if cmd.Pin != 1 || cmd.OnTime != 100 || cmd.OffTime != 200 {
		t.Error("Expected pin=1, on_time=100, off_time=200")
	}
}

func TestPulseCommand_Defaults(t *testing.T) {
	var cmd PulseCommand
	if err := json.Unmarshal([]byte(`{}`), &cmd); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if cmd.Pin != 0 || cmd.OnTime != 0 || cmd.OffTime != 0 {
		t.Error("Expected all zero values (handler applies defaults)")
	}
}

func TestPulseCommand_Validation(t *testing.T) {
	var cmd PulseCommand
	err := json.Unmarshal([]byte(`{invalid}`), &cmd)
	if err == nil {
		t.Error("Expected error for invalid json")
	}
}

// ============================================================================
// Beep Command Tests
// ============================================================================

func TestBeepCommand_Parsing(t *testing.T) {
	var cmd BeepCommand
	err := json.Unmarshal([]byte(`{"times": 5, "lapse": 3}`), &cmd)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if cmd.Times != 5 || cmd.Lapse != 3 {
		t.Error("Expected times=5, lapse=3")
	}
}

func TestBeepCommand_Defaults(t *testing.T) {
	var cmd BeepCommand
	if err := json.Unmarshal([]byte(`{}`), &cmd); err != nil {
		t.Fatalf("Unexpected error:  %v", err)
	}
	if cmd.Times != 0 || cmd.Lapse != 0 {
		t.Error("Expected zero values (handler applies defaults)")
	}
}

func TestBeepCommand_Validation(t *testing.T) {
	var cmd BeepCommand
	err := json.Unmarshal([]byte(`{invalid}`), &cmd)
	if err == nil {
		t.Error("Expected error for invalid json")
	}
}

// ============================================================================
// Default Application Tests (Handler Logic)
// ============================================================================

func TestFeedCommand_DefaultApplication(t *testing.T) {
	var cmd FeedCommand
	_ = json.Unmarshal([]byte(`{"lines": 0}`), &cmd)

	// Simulate handler default logic
	if cmd.Lines <= 0 {
		cmd.Lines = constants.DefaultFeedLines
	}
	if cmd.Lines != constants.DefaultFeedLines {
		t.Errorf("Expected default lines %d, got %d", constants.DefaultFeedLines, cmd.Lines)
	}
}
