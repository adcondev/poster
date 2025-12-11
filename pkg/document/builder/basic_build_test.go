package builder

import (
	"encoding/json"
	"testing"
)

// ============================================================================
// Basic Command Tests
// ============================================================================

func TestFeed(t *testing.T) {
	doc := NewDocument().
		SetProfile("Test", 80, "WPC1252").
		Feed(5).
		Build()

	var cmd feedCommand
	_ = json.Unmarshal(doc.Commands[0].Data, &cmd)

	if cmd.Lines != 5 {
		t.Errorf("Expected lines 5, got %d", cmd.Lines)
	}
}

func TestCut(t *testing.T) {
	doc := NewDocument().
		SetProfile("Test", 80, "WPC1252").
		Cut().
		Build()

	var cmd cutCommand
	_ = json.Unmarshal(doc.Commands[0].Data, &cmd)

	if cmd.Mode != "partial" {
		t.Errorf("Expected mode 'partial', got '%s'", cmd.Mode)
	}

	if cmd.Feed != 2 {
		t.Errorf("Expected feed 2, got %d", cmd.Feed)
	}
}

func TestFullCut(t *testing.T) {
	doc := NewDocument().
		SetProfile("Test", 80, "WPC1252").
		FullCut().
		Build()

	var cmd cutCommand
	_ = json.Unmarshal(doc.Commands[0].Data, &cmd)

	if cmd.Mode != "full" {
		t.Errorf("Expected mode 'full', got '%s'", cmd.Mode)
	}
}

func TestCutWithFeed(t *testing.T) {
	doc := NewDocument().
		SetProfile("Test", 80, "WPC1252").
		CutWithFeed("full", 10).
		Build()

	var cmd cutCommand
	_ = json.Unmarshal(doc.Commands[0].Data, &cmd)

	if cmd.Mode != "full" {
		t.Errorf("Expected mode 'full', got '%s'", cmd.Mode)
	}

	if cmd.Feed != 10 {
		t.Errorf("Expected feed 10, got %d", cmd.Feed)
	}
}

func TestSeparator(t *testing.T) {
	doc := NewDocument().
		SetProfile("Test", 80, "WPC1252").
		Separator("-").
		Build()

	var cmd separatorCommand
	_ = json.Unmarshal(doc.Commands[0].Data, &cmd)

	if cmd.Char != "-" {
		t.Errorf("Expected char '-', got '%s'", cmd.Char)
	}

	if cmd.Length != 48 {
		t.Errorf("Expected default length 48, got %d", cmd.Length)
	}
}

func TestSeparatorWithLength(t *testing.T) {
	doc := NewDocument().
		SetProfile("Test", 80, "WPC1252").
		SeparatorWithLength("=", 32).
		Build()

	var cmd separatorCommand
	_ = json.Unmarshal(doc.Commands[0].Data, &cmd)

	if cmd.Char != "=" {
		t.Errorf("Expected char '=', got '%s'", cmd.Char)
	}

	if cmd.Length != 32 {
		t.Errorf("Expected length 32, got %d", cmd.Length)
	}
}

func TestPulse(t *testing.T) {
	doc := NewDocument().
		SetProfile("Test", 80, "WPC1252").
		Pulse().
		Build()

	if doc.Commands[0].Type != "pulse" {
		t.Errorf("Expected type 'pulse', got '%s'", doc.Commands[0].Type)
	}

	var cmd pulseCommand
	_ = json.Unmarshal(doc.Commands[0].Data, &cmd)

	if cmd.Pin != 0 {
		t.Errorf("Expected pin 0, got %d", cmd.Pin)
	}

	if cmd.OnTime != 50 {
		t.Errorf("Expected on_time 50, got %d", cmd.OnTime)
	}

	if cmd.OffTime != 100 {
		t.Errorf("Expected off_time 100, got %d", cmd.OffTime)
	}
}

func TestPulseWithOptions(t *testing.T) {
	doc := NewDocument().
		SetProfile("Test", 80, "WPC1252").
		PulseWithOptions(1, 100, 200).
		Build()

	var cmd pulseCommand
	_ = json.Unmarshal(doc.Commands[0].Data, &cmd)

	if cmd.Pin != 1 {
		t.Errorf("Expected pin 1, got %d", cmd.Pin)
	}

	if cmd.OnTime != 100 {
		t.Errorf("Expected on_time 100, got %d", cmd.OnTime)
	}

	if cmd.OffTime != 200 {
		t.Errorf("Expected off_time 200, got %d", cmd.OffTime)
	}
}

func TestBeep(t *testing.T) {
	doc := NewDocument().
		SetProfile("Test", 80, "WPC1252").
		Beep(3, 5).
		Build()

	if doc.Commands[0].Type != "beep" {
		t.Errorf("Expected type 'beep', got '%s'", doc.Commands[0].Type)
	}

	var cmd beepCommand
	_ = json.Unmarshal(doc.Commands[0].Data, &cmd)

	if cmd.Times != 3 {
		t.Errorf("Expected times 3, got %d", cmd.Times)
	}

	if cmd.Lapse != 5 {
		t.Errorf("Expected lapse 5, got %d", cmd.Lapse)
	}
}
