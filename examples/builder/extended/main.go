// Package main demonstrates the fluent API for building print documents
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/adcondev/pos-printer/pkg/composer"
	"github.com/adcondev/pos-printer/pkg/connection"
	"github.com/adcondev/pos-printer/pkg/constants"
	"github.com/adcondev/pos-printer/pkg/document/builder"
	"github.com/adcondev/pos-printer/pkg/document/executor"
	"github.com/adcondev/pos-printer/pkg/document/schema"
	"github.com/adcondev/pos-printer/pkg/profile"
	"github.com/adcondev/pos-printer/pkg/service"
)

func main() {
	// =========================================================================
	// PARTE 1: Construir documento con API fluent
	// =========================================================================

	doc := builder.NewDocument().
		SetProfile("80mm EC-PM-80250", 80, "WPC1252").
		SetDPI(203).

		// Encabezado
		Text("TIENDA EJEMPLO").Bold().Size("2x2").Center().End().
		Text("Av. Principal #123").Center().End().
		Text("Tel: 555-1234").Center().End().
		Feed(1).
		Separator("=").

		// Información del ticket
		Text("TICKET DE VENTA").Bold().Center().End().
		Text("2025-01-15 14:30:00").Center().End().
		Separator("-").

		// Tabla de productos
		Table().
		Column("Producto", 20).
		Column("Cant", 5, constants.Center).
		Column("Precio", 10, constants.Right).
		PaperWidth(48).
		Row("Café Americano", "2", "$90.00").
		Row("Pan de chocolate", "3", "$45.00").
		Row("Jugo de naranja", "1", "$35.00").
		End().
		Separator("-").

		// Totales
		Text("$170.00").Bold().Right().WithLabel("SUBTOTAL").End().
		Text("$27.20").Right().WithLabel("IVA 16%").End().
		Separator("=").
		Text("$197.20").Bold().Size("2x2").Right().WithLabel("TOTAL").End().
		Feed(1).

		// QR de verificación
		QR("https://factura.ejemplo.com/ABC123").
		Size(150).
		Correction("M").
		WithText("Escanea para tu factura").
		Center().
		End().
		Feed(1).
		Text("¡Gracias por su compra!").Center().End().
		Text("Vuelva pronto").Center().End().
		Feed(3).
		Cut().
		Build()

	// =========================================================================
	// PARTE 2: Opciones de uso
	// =========================================================================

	// Opción A: Solo generar JSON (para debug o enviar a otro sistema)
	if len(os.Args) > 1 && os.Args[1] == "--json" {
		jsonBytes, err := doc.ToJSON()
		if err != nil {
			log.Fatalf("Error generating JSON: %v", err)
		}
		fmt.Println(string(jsonBytes))
		return
	}

	// Opción B: Imprimir directamente
	printerName := "80mm EC-PM-80250"
	if len(os.Args) > 1 {
		printerName = os.Args[1]
	}

	if err := printDocument(doc, printerName); err != nil {
		log.Fatalf("Print failed: %v", err)
	}

	log.Println("✅ Document printed successfully!")
}

func printDocument(doc *schema.Document, printerName string) error {
	// Crear conexión
	conn, err := connection.NewWindowsPrintConnector(printerName)
	if err != nil {
		return fmt.Errorf("connection failed: %w", err)
	}
	defer func(conn *connection.WindowsPrintConnector) {
		err := conn.Close()
		if err != nil {
			log.Printf("failed to close connection: %v", err)
		}
	}(conn)

	// Crear profile basado en documento
	var prof *profile.Escpos
	if doc.Profile.PaperWidth >= 80 {
		prof = profile.CreateProfile80mm()
	} else {
		prof = profile.CreateProfile58mm()
	}
	prof.Model = doc.Profile.Model

	// Crear printer y executor
	printer, err := service.NewPrinter(composer.NewEscpos(), prof, conn)
	if err != nil {
		return fmt.Errorf("printer creation failed: %w", err)
	}
	defer func(printer *service.Printer) {
		err := printer.Close()
		if err != nil {
			log.Printf("failed to close printer service: %v", err)
		}
	}(printer)

	exec := executor.NewExecutor(printer)
	return exec.Execute(doc)
}
