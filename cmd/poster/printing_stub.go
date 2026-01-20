//go:build !windows

package main

import (
	"fmt"

	"github.com/adcondev/poster/pkg/connection"
	"github.com/adcondev/poster/pkg/document/schema"
	"github.com/adcondev/poster/pkg/profile"
)

func detectPrinter() string {
	return ""
}

func createConnection(_ *Config) (connection.Connector, error) {
	return nil, fmt.Errorf("printing is only supported on Windows")
}

func createProfile(_ *schema.Document) *profile.Escpos {
	// Retornar un perfil dummy o nil para que compile
	return &profile.Escpos{}
}
