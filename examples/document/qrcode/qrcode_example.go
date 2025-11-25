// Package main implements an example of printing a document using the pos-printer library in JSON print job.
package main

import (
	"log"
	"os"

	"github.com/adcondev/pos-printer/pkg/composer"
	"github.com/adcondev/pos-printer/pkg/connection"
	"github.com/adcondev/pos-printer/pkg/document"
	"github.com/adcondev/pos-printer/pkg/profile"
	"github.com/adcondev/pos-printer/pkg/service"
)

func main() {
	// 1. Verificar archivos
	checkFile := func(path string) {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			log.Printf("⚠️  File not found: %s", path)
		} else {
			log.Printf("✅ File exists: %s", path)
		}
	}
	checkFile("./examples/document/qrcode/qr_test_advanced_1.json")
	checkFile("./examples/document/qrcode/qr_scenario_wifi.json")

	// ====== Iniciar impresión de documento JSON con QR avanzado =====

	fileName := "qr_test_advanced_1.json"
	jsonPath := "./examples/document/qrcode/" + fileName
	// Si el archivo no existe en esa ubicación, usar path alternativo
	if _, err := os.Stat(jsonPath); os.IsNotExist(err) {
		// Intentar en el directorio actual
		jsonPath = "./" + fileName
	}

	// 1. Crear perfil de impresora
	prof := profile.CreateECPM80250()

	// 2. Crear conexión
	conn, err := connection.NewWindowsPrintConnector(prof.Model)
	if err != nil {
		log.Fatalf("Failed to create connector: %v", err)
	}
	defer func(conn *connection.WindowsPrintConnector) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("Failed to close connector: %v", err)
		}
	}(conn)

	// 3. Crear protocolo
	proto := composer.NewEscpos()

	// 4. Crear servicio de impresora
	printer, err := service.NewPrinter(proto, prof, conn)
	if err != nil {
		log.Panicf("Failed to create printer: %v", err)
	}
	defer func(printer *service.Printer) {
		err := printer.Close()
		if err != nil {
			log.Fatalf("Failed to close printer: %v", err)
		}
	}(printer)

	// 5. Crear ejecutor de documentos
	executor := document.NewExecutor(printer)

	// Opción A: Cargar documento JSON desde archivo
	jsonData, err := os.ReadFile(jsonPath)
	if err != nil {
		log.Panicf("Failed to read JSON file: %v", err)
	}

	if err := executor.ExecuteJSON(jsonData); err != nil {
		log.Panicf("Failed to execute document: %v", err)
	}

	log.Println("Document printed successfully!")
}
