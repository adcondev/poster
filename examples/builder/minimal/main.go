// Package main implements a minimal example of using the document builder
package main

import (
	"fmt"

	"github.com/adcondev/pos-printer/pkg/document/builder"
)

func main() {
	// Ejemplo mínimo - solo genera JSON
	doc := builder.NewDocument().
		SetProfile("POS-80", 80, "WPC1252").
		Text("Hello, Printer!").Bold().Center().End().
		Barcode("EAN13", "5901234123457").Height(80).End().
		QR("https://example.com").Size(128).End().
		Cut().
		Build()

	json, _ := doc.ToJSON()
	fmt.Println(string(json))
}
