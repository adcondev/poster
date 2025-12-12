package builder

import (
	"encoding/json"
	"testing"
)

func TestBasicFeedCommands(t *testing.T) {
	doc := NewDocument().
		SetProfile("Test", 80, "WPC1252").
		Feed(5).
		Build()

	if doc.Commands[0].Type != "feed" {
		t.Errorf("Expected type 'feed', got '%s'", doc.Commands[0].Type)
	}

	var cmd feedCommand
	_ = json.Unmarshal(doc.Commands[0].Data, &cmd)

	if cmd.Lines != 5 {
		t.Errorf("Expected lines 5, got %d", cmd.Lines)
	}
}

func TestBasicCutCommands(t *testing.T) {
	tests := []struct {
		name         string
		buildFunc    func(*DocumentBuilder) *DocumentBuilder
		expectedMode string
		expectedFeed int
	}{
		{"PartialCut", func(b *DocumentBuilder) *DocumentBuilder { return b.Cut() }, "partial", 2},
		{"FullCut", func(b *DocumentBuilder) *DocumentBuilder { return b.FullCut() }, "full", 2},
		{"CutWithFeed", func(b *DocumentBuilder) *DocumentBuilder { return b.CutWithFeed("full", 10) }, "full", 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := NewDocument().SetProfile("Test", 80, "WPC1252")
			doc = tt.buildFunc(doc)
			result := doc.Build()

			if result.Commands[0].Type != "cut" {
				t.Errorf("Expected type 'cut', got '%s'", result.Commands[0].Type)
			}

			var cmd cutCommand
			_ = json.Unmarshal(result.Commands[0].Data, &cmd)

			if cmd.Mode != tt.expectedMode {
				t.Errorf("Expected mode '%s', got '%s'", tt.expectedMode, cmd.Mode)
			}

			if cmd.Feed != tt.expectedFeed {
				t.Errorf("Expected feed %d, got %d", tt.expectedFeed, cmd.Feed)
			}
		})
	}
}

func TestBasicSeparatorCommands(t *testing.T) {
	tests := []struct {
		name           string
		buildFunc      func(*DocumentBuilder) *DocumentBuilder
		expectedChar   string
		expectedLength int
	}{
		{"DefaultLength", func(b *DocumentBuilder) *DocumentBuilder { return b.Separator("-") }, "-", 48},
		{"CustomLength", func(b *DocumentBuilder) *DocumentBuilder { return b.SeparatorWithLength("=", 32) }, "=", 32},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := NewDocument().SetProfile("Test", 80, "WPC1252")
			doc = tt.buildFunc(doc)
			result := doc.Build()

			if result.Commands[0].Type != "separator" {
				t.Errorf("Expected type 'separator', got '%s'", result.Commands[0].Type)
			}

			var cmd separatorCommand
			_ = json.Unmarshal(result.Commands[0].Data, &cmd)

			if cmd.Char != tt.expectedChar {
				t.Errorf("Expected char '%s', got '%s'", tt.expectedChar, cmd.Char)
			}

			if cmd.Length != tt.expectedLength {
				t.Errorf("Expected length %d, got %d", tt.expectedLength, cmd.Length)
			}
		})
	}
}

func TestBasicPulseCommands(t *testing.T) {
	tests := []struct {
		name            string
		buildFunc       func(*DocumentBuilder) *DocumentBuilder
		expectedPin     int
		expectedOnTime  int
		expectedOffTime int
	}{
		{"DefaultPulse", func(b *DocumentBuilder) *DocumentBuilder { return b.Pulse() }, 0, 50, 100},
		{"CustomPulse", func(b *DocumentBuilder) *DocumentBuilder { return b.PulseWithOptions(1, 100, 200) }, 1, 100, 200},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := NewDocument().SetProfile("Test", 80, "WPC1252")
			doc = tt.buildFunc(doc)
			result := doc.Build()

			if result.Commands[0].Type != "pulse" {
				t.Errorf("Expected type 'pulse', got '%s'", result.Commands[0].Type)
			}

			var cmd pulseCommand
			_ = json.Unmarshal(result.Commands[0].Data, &cmd)

			if cmd.Pin != tt.expectedPin {
				t.Errorf("Expected pin %d, got %d", tt.expectedPin, cmd.Pin)
			}

			if cmd.OnTime != tt.expectedOnTime {
				t.Errorf("Expected on_time %d, got %d", tt.expectedOnTime, cmd.OnTime)
			}

			if cmd.OffTime != tt.expectedOffTime {
				t.Errorf("Expected off_time %d, got %d", tt.expectedOffTime, cmd.OffTime)
			}
		})
	}
}

func TestBasicBeepCommands(t *testing.T) {
	tests := []struct {
		name          string
		buildFunc     func(*DocumentBuilder) *DocumentBuilder
		expectedTimes int
		expectedLapse int
	}{
		{"DefaultBeep", func(b *DocumentBuilder) *DocumentBuilder { return b.BeepOnce() }, 1, 1},
		{"CustomBeep", func(b *DocumentBuilder) *DocumentBuilder { return b.Beep(3, 5) }, 3, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := NewDocument().SetProfile("Test", 80, "WPC1252")
			doc = tt.buildFunc(doc)
			result := doc.Build()

			if result.Commands[0].Type != "beep" {
				t.Errorf("Expected type 'beep', got '%s'", result.Commands[0].Type)
			}

			var cmd beepCommand
			_ = json.Unmarshal(result.Commands[0].Data, &cmd)

			if cmd.Times != tt.expectedTimes {
				t.Errorf("Expected times %d, got %d", tt.expectedTimes, cmd.Times)
			}

			if cmd.Lapse != tt.expectedLapse {
				t.Errorf("Expected lapse %d, got %d", tt.expectedLapse, cmd.Lapse)
			}
		})
	}
}
