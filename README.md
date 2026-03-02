# Poster: Thermal Printer Driver & ESC/POS Utility

![Go Version](https://img.shields.io/badge/go-1.24+-00ADD8?style=flat&logo=go)
![Platform](https://img.shields.io/badge/platform-windows-blue?style=flat&logo=windows)
![License](https://img.shields.io/badge/license-MIT-green)
[![CI](https://github.com/adcondev/poster/actions/workflows/ci.yml/badge.svg)](https://github.com/adcondev/poster/actions/workflows/ci.yml)

[![CodeQL](https://github.com/adcondev/poster/actions/workflows/codeql.yml/badge.svg)](https://github.com/adcondev/poster/security/code-scanning)
[![Release](https://github.com/adcondev/poster/actions/workflows/release.yml/badge.svg)](https://github.com/adcondev/poster/actions/workflows/release.yml)

[![codecov](https://codecov.io/github/adcondev/poster/graph/badge.svg?token=X9KVG99DDZ)](https://codecov.io/github/adcondev/poster)
[![Go Report Card](https://goreportcard.com/badge/github.com/adcondev/poster)](https://goreportcard.com/report/github.com/adcondev/poster)
[![Go Reference](https://pkg.go.dev/badge/github.com/adcondev/poster.svg)](https://pkg.go.dev/github.com/adcondev/poster)

[![Dependabot Status](https://flat.badgen.net/github/dependabot/adcondev/poster)](https://github.com/adcondev/poster/network/updates)
[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg?style=flat)](https://conventionalcommits.org)

<p align="center">
  <img src="assets/images/poster.png" alt="Poster Logo" width="576">
</p>

A professional, production-ready Go library and command-line utility for controlling ESC/POS thermal printers. Designed for high-reliability retail and POS environments, it features a robust JSON-based document protocol, native Windows Spooler integration, an advanced graphics engine, and a **visual emulator** for testing without physical hardware.

## 🚀 Key Features

- **JSON Document Protocol**: Define print jobs using a clean, versioned JSON schema (`v1.0`). Decouples business logic from hardware commands with full schema validation.
- **Native Windows Integration**: Prints directly via the Windows Print Spooler API (`winspool.drv`), supporting USB, Serial, and Network printers installed in Windows.
- **Advanced Graphics Engine**:
    - High-quality image printing with **Atkinson Dithering**.
    - Automatic scaling with bilinear interpolation.
    - Supports PNG, JPG, BMP formats.
- **Smart QR & Barcodes**: Automatically chooses between native printer firmware commands (fastest) or software rendering (maximum compatibility) based on the printer profile.
- **Dynamic Table Layout**: Built-in engine for generating perfectly aligned receipts with word wrapping, multi-column
  support, configurable spacing, **automatic overflow detection**, and **smart column auto-reduction** that preserves
  small columns while shrinking larger ones to fit paper width.
- **Visual Emulator**: Render print jobs as PNG images for preview and testing without physical hardware.
- **Hardware Agnostic**: Includes profiles for standard 80mm (Epson-compatible), 58mm (generic), PT-210, GP-58N, and
  EC-PM-80250 printers.
- **Smart Printer Discovery**: Auto-detects thermal/POS printers from installed Windows devices, filtering virtual
  printers and matching common hardware patterns.
- **11 Command Types**: `text`, `image`, `barcode`, `qr`, `table`, `separator`, `feed`, `cut`, `raw`, `pulse`, `beep`.
- **Raw Command Support**: Send raw ESC/POS bytes when full control is needed.

## 🏗️ Architecture

The project follows a layered architecture to ensure modularity and testability:

```mermaid
graph TD
    JSON["JSON Document<br/>(v1.0 Schema)"] --> Parser["Parser & Validator<br/>(pkg/document/schema)"]
    Parser --> Builder["Command Builder<br/>(pkg/document/builder)"]
    Builder --> Executor["Command Executor<br/>(pkg/document/executor)"]
    Executor --> Service["Printer Service<br/>(pkg/service)"]

    subgraph "Core Logic"
        Service --> Protocol["ESC/POS Composer<br/>(pkg/composer)"]
        Service --> Profile["Device Profile<br/>(pkg/profile)"]
        Service --> Graphics["Graphics Engine<br/>(pkg/graphics)"]
        Service --> Tables["Table Engine<br/>(pkg/tables)"]
    end

    subgraph "Command Modules (pkg/commands)"
        Protocol --> Barcode["Barcode"]
        Protocol --> BitImage["Bit Image"]
        Protocol --> Character["Character"]
        Protocol --> LineSpacing["Line Spacing"]
        Protocol --> MechControl["Mechanism Control"]
        Protocol --> PrintCmd["Print"]
        Protocol --> PrintPos["Print Position"]
        Protocol --> QRCode["QR Code"]
    end

    Service --> Connector["Connection Interface<br/>(pkg/connection)"]

    subgraph "Output Layer"
        Connector --> WinAPI["Windows Spooler API<br/>(winspool.drv)"]
        Connector --> Network["Network<br/>(TCP/9100)"]
        Connector --> Serial["Serial<br/>(COM Port)"]
        Connector --> FileOut["File<br/>(Debug Output)"]
        Connector --> Emulator["Visual Emulator<br/>(pkg/emulator)"]
        WinAPI --> Device["Physical Printer"]
        Emulator --> PNG["PNG Image"]
    end
```

### Package Structure

| Package          | Description                                                                                                                 |
|------------------|-----------------------------------------------------------------------------------------------------------------------------|
| `pkg/commands`   | ESC/POS command implementations (barcode, bitimage, character, linespacing, mechanismcontrol, print, printposition, qrcode) |
| `pkg/composer`   | ESC/POS byte sequence generation and protocol composition                                                                   |
| `pkg/connection` | Connection interfaces — Windows Spooler, Network, Serial, File                                                              |
| `pkg/constants`  | Shared constants, unit conversions, and default values                                                                      |
| `pkg/document`   | Document parsing, building, and execution (schema, builder, executor)                                                       |
| `pkg/emulator`   | Visual emulator for rendering print jobs as PNG images                                                                      |
| `pkg/graphics`   | Image processing pipeline — dithering, scaling, bitmap handling                                                             |
| `pkg/profile`    | Printer profiles and character encoding tables                                                                              |
| `pkg/service`    | High-level printer service facade (`PrinterActions` interface)                                                              |
| `pkg/tables`     | Table engine — column alignment, word wrapping, header styling, auto-reduction                                              |

## 📦 Installation

### Prerequisites

- **Go** 1.24 or higher
- **Windows OS** (for native spooler support; cross-platform for library/emulator use)

### As a Go Module

```bash
go get github.com/adcondev/poster@latest
```

### Build from Source

```bash
# Clone the repository
git clone https://github.com/adcondev/poster.git
cd poster

# Build the binary
go build -o poster.exe ./cmd/poster

# Or use Taskfile
task build
```

## 📖 Usage

### Command Line Interface (CLI)

The `poster` utility takes a JSON document and sends it to a specified printer.

```bash
# Print a document to a specific printer
poster.exe -file receipt.json -printer "EPSON TM-T88V"

# Dry-run (validate JSON without printing)
poster.exe -file receipt.json --dry-run

# List all installed printers (Windows only)
poster.exe --list

# List thermal/POS printers only
poster.exe --list-thermal

# List physical printers only (exclude virtual)
poster.exe --list-physical

# Output printer list as JSON
poster.exe --list --json

# Show version
poster.exe -v
```

### JSON Document Example

Create a file named `ticket.json`:

```json
{
  "version": "1.0",
  "profile": {
    "model": "Generic 80mm",
    "paper_width": 80,
    "dpi": 203,
    "has_qr": true
  },
  "commands": [
    {
      "type": "text",
      "data": {
        "content": {
          "text": "STORE NAME",
          "content_style": {
            "bold": true,
            "size": "2x2",
            "align": "center"
          }
        }
      }
    },
    {
      "type": "separator",
      "data": {
        "char": "-"
      }
    },
    {
      "type": "table",
      "data": {
        "definition": {
          "columns": [
            { "header": "Item", "width": 20, "align":  "left" },
            { "header":  "Price", "width": 10, "align": "right" }
          ]
        },
        "show_headers": true,
        "rows": [
          ["Coffee", "$3.50"],
          ["Sandwich", "$8.00"]
        ],
        "options": {
          "header_bold": true,
          "word_wrap": true,
          "column_spacing": 1,
          "auto_reduce": true
        }
      }
    },
    {
      "type": "qr",
      "data": {
        "data": "https://example.com",
        "align": "center",
        "pixel_width": 150
      }
    },
    {
      "type": "cut",
      "data": {
        "feed": 3
      }
    }
  ]
}
```

### Library Usage (Go)

You can also use the packages directly in your Go application:

```go
package main

import (
	"github.com/adcondev/poster/pkg/composer"
	"github.com/adcondev/poster/pkg/connection"
	"github.com/adcondev/poster/pkg/profile"
	"github.com/adcondev/poster/pkg/service"
)

func main() {
	// 1. Configure Profile
	prof := profile.CreateProfile80mm()

	// 2. Connect to Printer
	conn, _ := connection.NewWindowsPrintConnector("POS-80")
	defer conn.Close()

	// 3. Initialize Service
	proto := composer.NewEscpos()
	printer, _ := service.NewPrinter(proto, prof, conn)
	defer printer.Close()

	// 4. Print
	printer.Initialize()
	printer.PrintLine("Hello World!")
	printer.FullFeedAndCut(2)
}
```

## 📋 Supported Commands

| Command     | Description                                                         |
|-------------|---------------------------------------------------------------------|
| `text`      | Print formatted text with styles (bold, underline, inverse, sizing) |
| `image`     | Print images with dithering and scaling options                     |
| `barcode`   | Generate barcodes (CODE128, EAN13, UPC-A, CODE39, etc.)             |
| `qr`        | Generate QR codes with optional logos and human-readable text       |
| `table`     | Create formatted tables with column alignment and word wrapping     |
| `separator` | Print separator lines                                               |
| `feed`      | Advance paper by specified lines                                    |
| `cut`       | Perform full or partial paper cut                                   |
| `raw`       | Send raw ESC/POS bytes directly                                     |
| `pulse`     | Send electrical pulse to open cash drawer                           |
| `beep`      | Emit buzzer alert sound                                             |

For complete documentation, see [api/v1/DOCUMENT_V1.md](api/v1/DOCUMENT_V1.md).

## ⚙️ Configuration

### Connection Types

| Type      | Status        | Description                                                                 |
|-----------|---------------|-----------------------------------------------------------------------------|
| `windows` | ✅ Implemented | Windows Print Spooler (default). USB/Network printers installed in Windows. |
| `network` | 🔜 Planned    | Direct network connection via Raw TCP/9100                                  |
| `serial`  | 🔜 Planned    | Serial/USB direct connection (COM ports)                                    |
| `file`    | 🔜 Planned    | Output to file for debugging or emulator testing                            |

### Printer Profiles

The library includes built-in profiles for common hardware:

| Profile               | Description                                |
|-----------------------|--------------------------------------------|
| `CreateProfile80mm()` | Standard ESC/POS 80mm (Epson TM-T88, etc.) |
| `CreateProfile58mm()` | Generic 58mm thermal printers              |
| `CreatePt210()`       | PT-210 portable printer with QR support    |
| `CreateGP58N()`       | GP-58N 58mm thermal printer                |
| `CreateECPM80250()`   | EC-PM-80250 80mm thermal printer           |

## 🎨 Visual Emulator

The `pkg/emulator` package provides a visual emulator that renders print jobs as PNG images:

```go
import (
	"github.com/adcondev/poster/pkg/emulator"
)

// Create emulator with configuration
emu, _ := emulator.NewEngine(emulator.DefaultConfig())

// Render text
emu.SetBold(true)
emu.SetSize(2, 2)
emu.AlignCenter()
emu.PrintLine("STORE NAME")

// Render and save as PNG
img := emu.Render()
```

## 📊 Code Coverage

<details>
<summary>📈 Coverage Trends</summary>

[![Coverage Graph](https://codecov.io/github/adcondev/poster/graphs/sunburst.svg?token=X9KVG99DDZ
)](https://codecov.io/gh/adcondev/poster)
</details>

## 🤝 Contributing

Contributions are welcome! Please read [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on how to submit pull requests.
This project uses [Conventional Commits](https://conventionalcommits.org) and enforces semantic PR titles.

## 📄 License

This project is licensed under the MIT License — see the [LICENSE](LICENSE) file for details.
