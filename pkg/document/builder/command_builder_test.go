package builder

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/adcondev/poster/pkg/document/schema"

	"github.com/adcondev/poster/pkg/constants"
)

const Release = constants.DefaultVersion

// ============================================================================
// NewDocument Tests
// ============================================================================

func TestNewDocument(t *testing.T) {
	doc := NewDocument()

	if doc == nil {
		t.Fatal("NewDocument() returned nil")
	}

	if doc.version != Release {
		t.Errorf("Expected version '1.0', got '%s'", doc.version)
	}

	if doc.commands == nil {
		t.Error("Expected commands to be initialized, got nil")
	}

	if len(doc.commands) != 0 {
		t.Errorf("Expected empty commands slice, got %d commands", len(doc.commands))
	}
}

// ============================================================================
// SetProfile Tests
// ============================================================================

func TestSetProfile(t *testing.T) {
	tests := []struct {
		name       string
		model      string
		paperWidth int
		codeTable  string
	}{
		{"58mm printer", "PT-210", 58, "PC850"},
		{"80mm printer", "EPSON TM-T20II", 80, "WPC1252"},
		{"custom printer", "Custom Model", 72, "PC437"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := NewDocument().SetProfile(tt.model, tt.paperWidth, tt.codeTable)

			if doc.profile.Model != tt.model {
				t.Errorf("Expected model '%s', got '%s'", tt.model, doc.profile.Model)
			}
			if doc.profile.PaperWidth != tt.paperWidth {
				t.Errorf("Expected paper width %d, got %d", tt.paperWidth, doc.profile.PaperWidth)
			}
			if doc.profile.CodeTable != tt.codeTable {
				t.Errorf("Expected code table '%s', got '%s'", tt.codeTable, doc.profile.CodeTable)
			}
		})
	}
}

func TestSetDPI(t *testing.T) {
	doc := NewDocument().SetProfile("Test", 80, "WPC1252").SetDPI(300)

	if doc.profile.DPI != 300 {
		t.Errorf("Expected DPI 300, got %d", doc.profile.DPI)
	}
}

func TestSetHasQR(t *testing.T) {
	doc := NewDocument().SetProfile("Test", 80, "WPC1252").SetHasQR(true)

	if !doc.profile.HasQR {
		t.Error("Expected HasQR to be true")
	}
}

func TestEnableDebug(t *testing.T) {
	doc := NewDocument().EnableDebug()

	if !doc.debugLog {
		t.Error("Expected debugLog to be true")
	}
}

// ============================================================================
// Command Chaining Tests
// ============================================================================

func TestCommandChaining(t *testing.T) {
	doc := NewDocument().
		SetProfile("Test", 80, "WPC1252").
		Text("Hello").End().
		Feed(3).
		Separator("-").
		Cut()

	result := doc.Build()

	if len(result.Commands) != 4 {
		t.Errorf("Expected 4 commands, got %d", len(result.Commands))
	}

	expectedTypes := []string{"text", "feed", "separator", "cut"}
	for i, expected := range expectedTypes {
		if result.Commands[i].Type != expected {
			t.Errorf("Command %d: expected type '%s', got '%s'", i, expected, result.Commands[i].Type)
		}
	}
}

func TestFluentAPIReturnsBuilder(t *testing.T) {
	// Ensure all methods return *DocumentBuilder for chaining
	builder := NewDocument()

	// Profile methods
	builder = builder.SetProfile("Test", 80, "WPC1252")
	builder = builder.SetDPI(203)
	builder = builder.SetHasQR(true)
	builder = builder.EnableDebug()

	// Simple commands
	builder = builder.Feed(1)
	builder = builder.Cut()
	builder = builder.FullCut()
	builder = builder.Separator("-")
	builder = builder.Pulse()
	builder = builder.PulseWithOptions(1, 100, 200) // ADD THIS
	builder = builder.Beep(1, 1)                    // ADD THIS

	if builder == nil {
		t.Error("Fluent API chain returned nil")
	}
}

// ============================================================================
// Build Tests
// ============================================================================

func TestBuild(t *testing.T) {
	doc := NewDocument().
		SetProfile("EPSON", 80, "WPC1252").
		SetDPI(203).
		SetHasQR(true).
		EnableDebug().
		Text("Test").End().
		Build()

	if doc.Version != Release {
		t.Errorf("Expected version '1.0', got '%s'", doc.Version)
	}

	if doc.Profile.Model != "EPSON" {
		t.Errorf("Expected model 'EPSON', got '%s'", doc.Profile.Model)
	}

	if !doc.DebugLog {
		t.Error("Expected DebugLog to be true")
	}

	if len(doc.Commands) != 1 {
		t.Errorf("Expected 1 command, got %d", len(doc.Commands))
	}
}

// ============================================================================
// ToJSON Tests
// ============================================================================

func TestToJSON(t *testing.T) {
	jsonBytes, err := NewDocument().
		SetProfile("Test", 80, "WPC1252").
		Text("Hello World").Center().Bold().End().
		Cut().
		ToJSON()

	if err != nil {
		t.Fatalf("ToJSON() error: %v", err)
	}

	if len(jsonBytes) == 0 {
		t.Error("ToJSON() returned empty bytes")
	}

	// Verify it's valid JSON
	var parsed schema.Document
	if err := json.Unmarshal(jsonBytes, &parsed); err != nil {
		t.Fatalf("ToJSON() produced invalid JSON: %v", err)
	}

	if parsed.Version != Release {
		t.Errorf("Expected version '1.0' in JSON, got '%s'", parsed.Version)
	}
}

func TestToJSONContainsExpectedFields(t *testing.T) {
	jsonBytes, _ := NewDocument().
		SetProfile("TestPrinter", 58, "PC850").
		SetDPI(203).
		Text("Test").End().
		ToJSON()

	jsonStr := string(jsonBytes)

	expectedFields := []string{
		`"version"`,
		`"profile"`,
		`"model"`,
		`"paper_width"`,
		`"commands"`,
	}

	for _, field := range expectedFields {
		if !strings.Contains(jsonStr, field) {
			t.Errorf("Expected JSON to contain %s", field)
		}
	}
}

// ============================================================================
// Complex Document Tests
// ============================================================================

func TestCompleteReceipt(t *testing.T) {
	doc := NewDocument().
		SetProfile("EPSON TM-T20II", 80, "WPC1252").
		SetDPI(203).
		Text("STORE NAME").Bold().Size("2x2").Center().End().
		Separator("=").
		Table().
		Column("Item", 20, constants.Left).
		Column("Price", 10, constants.Right).
		Row("Coffee", "$4.50").
		Row("Muffin", "$3.00").
		End().
		Separator("-").
		Text("$7.50").Bold().Right().WithLabel("TOTAL").End().
		Feed(2).
		QR("https://receipt.example.com/123").Size(200).Center().End().
		Feed(3).
		Cut().
		Build()

	// Verify document structure
	if doc.Version != Release {
		t.Errorf("Expected version '1.0', got '%s'", doc.Version)
	}

	if doc.Profile.Model != "EPSON TM-T20II" {
		t.Errorf("Expected model 'EPSON TM-T20II', got '%s'", doc.Profile.Model)
	}

	// Expected command sequence
	expectedTypes := []string{
		"text",      // STORE NAME
		"separator", // =
		"table",     // items
		"separator", // -
		"text",      // TOTAL
		"feed",      // 2 lines
		"qr",        // receipt QR
		"feed",      // 3 lines
		"cut",       // partial cut
	}

	if len(doc.Commands) != len(expectedTypes) {
		t.Fatalf("Expected %d commands, got %d", len(expectedTypes), len(doc.Commands))
	}

	for i, expected := range expectedTypes {
		if doc.Commands[i].Type != expected {
			t.Errorf("Command %d: expected type '%s', got '%s'", i, expected, doc.Commands[i].Type)
		}
	}
}

// ============================================================================
// Benchmark Tests
// ============================================================================

func BenchmarkNewDocument(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewDocument()
	}
}

func BenchmarkBuildSimple(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewDocument().
			SetProfile("Test", 80, "WPC1252").
			Text("Hello").End().
			Build()
	}
}

func BenchmarkBuildComplex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewDocument().
			SetProfile("Test", 80, "WPC1252").
			Text("Name").Bold().Center().End().
			Separator("=").
			Table().
			Column("Item", 20).
			Column("Price", 10).
			Row("Item 1", "$10.00").
			Row("Item 2", "$20.00").
			End().
			QR("https://example.com").End().
			Cut().
			Build()
	}
}

func BenchmarkToJSON(b *testing.B) {
	builder := NewDocument().
		SetProfile("Test", 80, "WPC1252").
		Text("Test").End()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = builder.ToJSON()
	}
}
