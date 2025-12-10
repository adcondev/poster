// Package main demuestra cómo enviar comandos RAW no estándar (específicos del fabricante).
//
// Caso de uso: Muchas impresoras genéricas/chinas utilizan 'ESC B' (1B 42) para el bíper,
// mientras que el estándar oficial ESC/POS no esta definido.
// Este ejemplo muestra cómo la librería permite enviar bytes crudos para cubrir estos casos de borde.
package main

import (
	"log"
	"time"

	"github.com/adcondev/poster/pkg/composer"
	"github.com/adcondev/poster/pkg/connection"
	"github.com/adcondev/poster/pkg/profile"
	"github.com/adcondev/poster/pkg/service"
)

func main() {
	// 1. Configuración del perfil
	// Usamos un perfil genérico de 80mm (EC-PM-80250)
	prof := profile.CreateECPM80250()

	// 2. Establecer conexión (Windows en este ejemplo)
	// Nota: Asegúrate de que el nombre de la impresora coincida con el de tu sistema
	conn, err := connection.NewWindowsPrintConnector(prof.Model)
	if err != nil {
		log.Fatalf("Error de conexión: %v", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("Error cerrando conexión: %v", err)
		}
	}()

	// 3. Inicializar protocolo y servicio
	proto := composer.NewEscpos()
	printer, err := service.NewPrinter(proto, prof, conn)
	if err != nil {
		log.Panicf("Error creando servicio de impresora: %v", err)
	}
	defer func() {
		if err := printer.Close(); err != nil {
			log.Printf("Error cerrando servicio: %v", err)
		}
	}()

	// Inicializar impresora (resetea buffer y estados)
	if err := printer.Initialize(); err != nil {
		log.Panic(err)
	}

	log.Println("Iniciando prueba de Beeper Genérico...")

	// 4. Construcción del Comando RAW
	// -------------------------------------------------------------------------
	// ADVERTENCIA DE COMPATIBILIDAD:
	// El estándar ESC/POS define 'ESC ( A' (Pulse) para el buzzer.
	// Sin embargo, este hardware específico usa 'ESC B' (no existe es ESC/POS).
	//
	// Formato: ESC B n t
	// Hex:     1B  42 09 02
	// -------------------------------------------------------------------------
	// 1B 42 -> Cabecera del comando (ESC B)
	// 09    -> n: Número de repeticiones (9 veces)
	// 02    -> t: Duración/Intervalo (factor de tiempo, aprox 100ms * t)
	// -------------------------------------------------------------------------
	beeperCmd := []byte{0x1B, 0x42, 0x09, 0x02}

	// Enviar bytes directamente (bypass del protocolo estándar)
	if err := printer.Write(beeperCmd); err != nil {
		log.Panicf("Fallo al enviar comando raw: %v", err)
	}

	// Pequeña pausa para permitir que el sonido termine antes de cortar
	time.Sleep(2 * time.Second)

	// 5. Imprimir confirmación visual usando métodos estándar
	if err := printer.PrintLine("Prueba de sonido completada."); err != nil {
		log.Printf("Error imprimiendo texto: %v", err)
	}

	if err := printer.PartialFeedAndCut(1); err != nil {
		log.Printf("Error cortando papel: %v", err)
	}

	log.Println("Finalizado con éxito.")
}
