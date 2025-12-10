// Package main implements an example of printing a document using the poster library in JSON print job.
package main

import (
	"log"
	"os"

	"github.com/adcondev/poster/pkg/composer"
	"github.com/adcondev/poster/pkg/connection"
	"github.com/adcondev/poster/pkg/document/executor"
	"github.com/adcondev/poster/pkg/profile"
	"github.com/adcondev/poster/pkg/service"
)

func main() {
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
	exec := executor.NewExecutor(printer)

	// Opción A: Cargar documento JSON desde archivo
	jsonData, err := os.ReadFile("./examples/document/basic_ticket.json")
	if err != nil {
		log.Panicf("Failed to read JSON file: %v", err)
	}

	if err := exec.ExecuteJSON(jsonData); err != nil {
		log.Panicf("Failed to execute document: %v", err)
	}

	log.Println("Document printed successfully!")
}
