//go:build !windows

package main

import (
	"fmt"
	"strings"

	"github.com/adcondev/poster/pkg/connection"
)

func detectPrinter() string {
	return ""
}

func createConnection(config *Config) (connection.Connector, error) {
	switch strings.ToLower(config.ConnectionType) {
	case win:
		return nil, fmt.Errorf("windows spooler connection is only supported on Windows (use -t tcp, -t file, --preview, or --dry-run on this OS)")

	case "network", "tcp":
		target := config.PrinterName
		if target == "" {
			target = config.NetworkAddr
		}
		if target == "" {
			return nil, fmt.Errorf("printer name or address required for network/tcp connection")
		}
		return connection.NewTCPConnector(target)

	case "serial":
		if config.SerialPort == "" {
			return nil, fmt.Errorf("serial port required")
		}
		return nil, fmt.Errorf("serial connection not yet implemented")

	case "file":
		targetPath := config.OutputFile
		if targetPath == "" {
			targetPath = config.PrinterName
		}
		if targetPath == "" {
			return nil, fmt.Errorf("output file path required for file connection (use -output or -p)")
		}
		return connection.NewFileConnector(targetPath)

	default:
		return nil, fmt.Errorf("unknown connection type: %s", config.ConnectionType)
	}
}

