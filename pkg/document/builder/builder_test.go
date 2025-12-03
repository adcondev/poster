package builder

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/adcondev/pos-printer/pkg/document/schema"

	"github.com/adcondev/pos-printer/pkg/constants"
)

const Release = "1.0"

// ============================================================================
// NewDocument Tests
// ============================================================================

func TestNewDocument(t *testing.T) {
	doc := NewDocument()

	if doc == nil {
		t.Fatal("NewDocument() returned nil")
	}

	if doc.version != Release {
		t.Errorf("Expected version '1. 0', got '%s'", doc.version)
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
		t.Errorf("Expected version '1. 0', got '%s'", doc.Version)
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
// Text Builder Tests
// ============================================================================

func TestTextBuilder(t *testing.T) {
	doc := NewDocument().
		SetProfile("Test", 80, "WPC1252").
		Text("Hello").
		Bold().
		Size("2x2").
		Underline("1pt").
		Inverse().
		Font("B").
		Center().
		End().
		Build()

	if len(doc.Commands) != 1 {
		t.Fatalf("Expected 1 command, got %d", len(doc.Commands))
	}

	if doc.Commands[0].Type != "text" {
		t.Errorf("Expected type 'text', got '%s'", doc.Commands[0].Type)
	}

	// Parse the command data
	var cmd textCommand
	if err := json.Unmarshal(doc.Commands[0].Data, &cmd); err != nil {
		t.Fatalf("Failed to unmarshal text command: %v", err)
	}

	if cmd.Content.Text != "Hello" {
		t.Errorf("Expected text 'Hello', got '%s'", cmd.Content.Text)
	}
}

func TestTextWithLabel(t *testing.T) {
	doc := NewDocument().
		SetProfile("Test", 80, "WPC1252").
		Text("$100. 00").
		WithLabel("Total").
		LabelSeparator(": ").
		Bold().
		Right().
		End().
		Build()

	var cmd textCommand
	_ = json.Unmarshal(doc.Commands[0].Data, &cmd)

	if cmd.Label == nil {
		t.Fatal("Expected label to be set")
	}

	if cmd.Label.Text != "Total" {
		t.Errorf("Expected label 'Total', got '%s'", cmd.Label.Text)
	}
}

func TestTextNoNewLine(t *testing.T) {
	doc := NewDocument().
		SetProfile("Test", 80, "WPC1252").
		Text("No newline").NoNewLine().End().
		Build()

	var cmd textCommand
	_ = json.Unmarshal(doc.Commands[0].Data, &cmd)

	if cmd.NewLine == nil || *cmd.NewLine != false {
		t.Error("Expected NewLine to be false")
	}
}

// ============================================================================
// Image Builder Tests
// ============================================================================

func TestImageBuilder(t *testing.T) {
	doc := NewDocument().
		SetProfile("Test", 80, "WPC1252").
		Image("base64data").
		Width(256).
		Threshold(100).
		Dithering(constants.Atkinson).
		Scaling(constants.Bilinear).
		Center().
		End().
		Build()

	if doc.Commands[0].Type != "image" {
		t.Errorf("Expected type 'image', got '%s'", doc.Commands[0].Type)
	}

	var cmd imageCommand
	_ = json.Unmarshal(doc.Commands[0].Data, &cmd)

	if cmd.PixelWidth != 256 {
		t.Errorf("Expected pixel width 256, got %d", cmd.PixelWidth)
	}

	if cmd.Threshold != 100 {
		t.Errorf("Expected threshold 100, got %d", cmd.Threshold)
	}

	if cmd.Dithering != constants.Atkinson.String() {
		t.Errorf("Expected dithering 'atkinson', got '%s'", cmd.Dithering)
	}
}

func TestImageBuilderDefaults(t *testing.T) {
	doc := NewDocument().
		SetProfile("Test", 80, "WPC1252").
		Image("base64data").End().
		Build()

	var cmd imageCommand
	_ = json.Unmarshal(doc.Commands[0].Data, &cmd)

	if cmd.PixelWidth != constants.DefaultImagePixelWidth {
		t.Errorf("Expected default pixel width %d, got %d", constants.DefaultImagePixelWidth, cmd.PixelWidth)
	}

	if cmd.Threshold != constants.DefaultImageThreshold {
		t.Errorf("Expected default threshold %d, got %d", constants.DefaultImageThreshold, cmd.Threshold)
	}

	if cmd.Align != constants.DefaultImageAlignment.String() {
		t.Errorf("Expected default align '%s', got '%s'", constants.DefaultImageAlignment.String(), cmd.Align)
	}
}

// ============================================================================
// QR Builder Tests
// ============================================================================

func TestQRBuilder(t *testing.T) {
	doc := NewDocument().
		SetProfile("Test", 80, "WPC1252").
		QR("https://example.com").
		Size(200).
		Correction("H").
		WithText("Scan me").
		WithLogo("logobase64").
		Center().
		End().
		Build()

	if doc.Commands[0].Type != "qr" {
		t.Errorf("Expected type 'qr', got '%s'", doc.Commands[0].Type)
	}

	var cmd qrCommand
	_ = json.Unmarshal(doc.Commands[0].Data, &cmd)

	if cmd.Data != "https://example.com" {
		t.Errorf("Expected data 'https://example. com', got '%s'", cmd.Data)
	}

	if cmd.PixelWidth != 200 {
		t.Errorf("Expected pixel width 200, got %d", cmd.PixelWidth)
	}

	if cmd.Correction != "H" {
		t.Errorf("Expected correction 'H', got '%s'", cmd.Correction)
	}

	if cmd.HumanText != "Scan me" {
		t.Errorf("Expected human text 'Scan me', got '%s'", cmd.HumanText)
	}
}

func TestQRBuilderDefaults(t *testing.T) {
	doc := NewDocument().
		SetProfile("Test", 80, "WPC1252").
		QR("data").End().
		Build()

	var cmd qrCommand
	_ = json.Unmarshal(doc.Commands[0].Data, &cmd)

	if cmd.PixelWidth != 128 {
		t.Errorf("Expected default pixel width 128, got %d", cmd.PixelWidth)
	}

	if cmd.Correction != "Q" {
		t.Errorf("Expected default correction 'Q', got '%s'", cmd.Correction)
	}

	if cmd.Align != constants.Center.String() {
		t.Errorf("Expected default align 'center', got '%s'", cmd.Align)
	}
}

// ============================================================================
// Barcode Builder Tests
// ============================================================================

func TestBarcodeBuilder(t *testing.T) {
	width := 3
	height := 100

	doc := NewDocument().
		SetProfile("Test", 80, "WPC1252").
		Barcode("CODE128", "12345").
		Width(width).
		Height(height).
		HRIPosition("below").
		HRIFont("A").
		Center().
		End().
		Build()

	if doc.Commands[0].Type != "barcode" {
		t.Errorf("Expected type 'barcode', got '%s'", doc.Commands[0].Type)
	}

	var cmd barcodeCommand
	_ = json.Unmarshal(doc.Commands[0].Data, &cmd)

	if cmd.Symbology != "CODE128" {
		t.Errorf("Expected symbology 'CODE128', got '%s'", cmd.Symbology)
	}

	if cmd.Data != "12345" {
		t.Errorf("Expected data '12345', got '%s'", cmd.Data)
	}

	if *cmd.Width != width {
		t.Errorf("Expected width %d, got %d", width, *cmd.Width)
	}
}

func TestBarcodeBuilderAlign(t *testing.T) {
	doc := NewDocument().
		SetProfile("Test", 80, "WPC1252").
		Barcode("EAN13", "1234567890123").
		Align(constants.Right).
		End().
		Build()

	var cmd barcodeCommand
	_ = json.Unmarshal(doc.Commands[0].Data, &cmd)

	if *cmd.Align != constants.Right.String() {
		t.Errorf("Expected align 'right', got '%s'", *cmd.Align)
	}
}

// ============================================================================
// Table Builder Tests
// ============================================================================

func TestTableBuilder(t *testing.T) {
	doc := NewDocument().
		SetProfile("Test", 80, "WPC1252").
		Table().
		Column("Item", 20, constants.Left).
		Column("Price", 10, constants.Right).
		Row("Coffee", "$3.50").
		Row("Muffin", "$4.25").
		PaperWidth(48).
		End().
		Build()

	if doc.Commands[0].Type != "table" {
		t.Errorf("Expected type 'table', got '%s'", doc.Commands[0].Type)
	}

	var cmd tableCommand
	_ = json.Unmarshal(doc.Commands[0].Data, &cmd)

	if len(cmd.Definition.Columns) != 2 {
		t.Errorf("Expected 2 columns, got %d", len(cmd.Definition.Columns))
	}

	if len(cmd.Rows) != 2 {
		t.Errorf("Expected 2 rows, got %d", len(cmd.Rows))
	}
}

func TestTableBuilderOptions(t *testing.T) {
	doc := NewDocument().
		SetProfile("Test", 80, "WPC1252").
		Table().
		Column("Test", 10).
		Row("Data").
		HideHeaders().
		NoHeaderBold().
		NoWordWrap().
		ColumnSpacing(2).
		Align(constants.Center).
		End().
		Build()

	var cmd tableCommand
	_ = json.Unmarshal(doc.Commands[0].Data, &cmd)

	if cmd.ShowHeaders != false {
		t.Error("Expected ShowHeaders to be false")
	}

	if cmd.Options.HeaderBold != false {
		t.Error("Expected HeaderBold to be false")
	}

	if cmd.Options.WordWrap != false {
		t.Error("Expected WordWrap to be false")
	}

	if cmd.Options.ColumnSpacing != 2 {
		t.Errorf("Expected ColumnSpacing 2, got %d", cmd.Options.ColumnSpacing)
	}
}

func TestTableBuilderRows(t *testing.T) {
	rows := [][]string{
		{"A", "1"},
		{"B", "2"},
		{"C", "3"},
	}

	doc := NewDocument().
		SetProfile("Test", 80, "WPC1252").
		Table().
		Column("Letter", 10).
		Column("Number", 10).
		Rows(rows).
		End().
		Build()

	var cmd tableCommand
	_ = json.Unmarshal(doc.Commands[0].Data, &cmd)

	if len(cmd.Rows) != 3 {
		t.Errorf("Expected 3 rows, got %d", len(cmd.Rows))
	}
}

// ============================================================================
// Raw Builder Tests
// ============================================================================

func TestRawBuilder(t *testing.T) {
	doc := NewDocument().
		SetProfile("Test", 80, "WPC1252").
		Raw("1B 40").
		Comment("Initialize printer").
		Safe().
		End().
		Build()

	if doc.Commands[0].Type != "raw" {
		t.Errorf("Expected type 'raw', got '%s'", doc.Commands[0].Type)
	}

	var cmd rawCommand
	_ = json.Unmarshal(doc.Commands[0].Data, &cmd)

	if cmd.Hex != "1B 40" {
		t.Errorf("Expected hex '1B 40', got '%s'", cmd.Hex)
	}

	if cmd.Comment != "Initialize printer" {
		t.Errorf("Expected comment 'Initialize printer', got '%s'", cmd.Comment)
	}

	if !cmd.SafeMode {
		t.Error("Expected SafeMode to be true")
	}
}

func TestRawBuilderBase64(t *testing.T) {
	doc := NewDocument().
		SetProfile("Test", 80, "WPC1252").
		Raw("G0A=").Base64().End().
		Build()

	var cmd rawCommand
	_ = json.Unmarshal(doc.Commands[0].Data, &cmd)

	if cmd.Format != "base64" {
		t.Errorf("Expected format 'base64', got '%s'", cmd.Format)
	}
}

func TestFromBytes(t *testing.T) {
	data := []byte{0x1B, 0x40}
	doc := NewDocument().
		SetProfile("Test", 80, "WPC1252").
		FromBytes(data).End().
		Build()

	var cmd rawCommand
	_ = json.Unmarshal(doc.Commands[0].Data, &cmd)

	if cmd.Hex != "1b40" {
		t.Errorf("Expected hex '1b40', got '%s'", cmd.Hex)
	}
}

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
		Row("Muffin", "$3. 00").
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
		t.Errorf("Expected version '1. 0', got '%s'", doc.Version)
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
			Text("Header").Bold().Center().End().
			Separator("=").
			Table().
			Column("Item", 20).
			Column("Price", 10).
			Row("Item 1", "$10. 00").
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
