// Package builder provides a fluent API for creating print documents.
//
// This package is the authoring component of the poster library. It allows
// constructing print documents programmatically using method chaining,
// which are then serialized to JSON for the executor package.
//
// # Quick Start
//
//	doc := builder.NewDocument().
//	    SetProfile("POS-80", 80, "WPC1252").
//	    Text("Hello World").Bold().Center().End().
//	    Cut().
//	    Build()
//
//	jsonBytes, _ := doc.ToJSON()
//
// # Architecture
//
// The builder uses a parent-child pattern. Complex commands return
// specialized builders that chain back to the document builder:
//
//	DocumentBuilder
//	    ├── Text()    → TextBuilder    → End() → DocumentBuilder
//	    ├── Table()   → TableBuilder   → End() → DocumentBuilder
//	    ├── QR()      → QRBuilder      → End() → DocumentBuilder
//	    ├── Barcode() → BarcodeBuilder → End() → DocumentBuilder
//	    ├── Image()   → ImageBuilder   → End() → DocumentBuilder
//	    └── Raw()     → RawBuilder     → End() → DocumentBuilder
//
// Simple commands return directly to DocumentBuilder:
//
//	. Feed(3). Separator("-").Cut()
//
// # Available Commands
//
//	Method          Returns           Description
//	─────────────────────────────────────────────────────
//	Text(s)         *TextBuilder      Styled text with labels
//	Table()         *TableBuilder     Multi-column tables
//	QR(data)        *QRBuilder        QR codes with logos
//	Barcode(s,d)    *BarcodeBuilder   1D barcodes
//	Image(b64)      *ImageBuilder     Images with dithering
//	Raw(hex)        *RawBuilder       Direct ESC/POS bytes
//	Feed(n)         *DocumentBuilder  Paper advance
//	Cut()           *DocumentBuilder  Partial cut
//	FullCut()       *DocumentBuilder  Full cut
//	Separator(c)    *DocumentBuilder  Line separator
//	Pulse()         *DocumentBuilder  Cash drawer
//	Beep(n,t)       *DocumentBuilder  Buzzer sound
//	Pulse()         *DocumentBuilder  Cash drawer
//	Beep(n,t)       *DocumentBuilder  Buzzer sound
//
// # Example: Complete Receipt
//
//	doc := builder.NewDocument().
//	    SetProfile("EPSON", 80, "WPC1252").
//	    Text("STORE NAME").Bold().Size("2x2").Center().End().
//	    Separator("=").
//	    Table().
//	        Column("Item", 20).
//	        Column("Price", 10, builder.Right).
//	        Row("Coffee", "$4.50").
//	        Row("Muffin", "$3.00").
//	    End().
//	    Separator("-").
//	    Text("$7.50").Bold().Right().WithLabel("TOTAL").End().
//	    Feed(2).
//	    QR("https://receipt.example.com/123").Size(200).End().
//	    Feed(3).
//	    Cut().
//	    Build()
//
// # Alignment Constants
//
//	builder.Left    "left"
//	builder.Center  "center"
//	builder.Right   "right"
//
// # Output Formats
//
//	doc.Build()   - Returns *Document struct
//	doc.ToJSON()  - Returns []byte (JSON)
//
// # Related Packages
//
//	executor - Parses and executes the generated JSON
//	service  - Low-level printer communication
package builder
