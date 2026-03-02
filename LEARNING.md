# Poster — ESC/POS Thermal Printer Driver & Utility

## Project Overview

**Poster** is a production-ready Go library and CLI utility that provides a full abstraction layer for ESC/POS thermal
printers. It decouples business logic from hardware specifics through a versioned JSON document protocol, enabling
developers to define complex print jobs (text, images, barcodes, QR codes, tables) without writing raw byte sequences.
The system includes native Windows Print Spooler integration, a custom graphics pipeline with Atkinson dithering, a
dynamic table layout engine, and a visual emulator for testing without physical hardware.

## Tech Stack & Infrastructure

| Category          | Technologies                                                                                        |
|-------------------|-----------------------------------------------------------------------------------------------------|
| **Language**      | Go 1.24 (primary), Node.js (release tooling)                                                        |
| **Platform**      | Windows (primary, `winspool.drv` via `syscall`/`unsafe`), Linux/macOS (cross-platform test targets) |
| **Protocol**      | ESC/POS binary protocol (Epson standard)                                                            |
| **API / Schema**  | JSON Schema v1.0, versioned document format                                                         |
| **CI/CD**         | GitHub Actions — 11 workflows (CI, Release, CodeQL, Dependabot Automerge, Stale, Tag Protection)    |
| **Security**      | Trivy vulnerability scanner, GitHub CodeQL static analysis                                          |
| **Testing**       | Multi-OS matrix (Ubuntu, Windows, macOS), race detection, Codecov coverage, performance benchmarks  |
| **Linting**       | golangci-lint (exhaustive config: `gofmt`, `govet`, `errcheck`, `staticcheck`, `gosec`, etc.)       |
| **Release Mgmt**  | Conventional Commits, `standard-version`, automated semantic versioning via GitHub App tokens       |
| **Quality Gates** | Pre-commit hooks, commit linting (`@commitlint`), semantic PR title validation                      |
| **Task Runner**   | Taskfile v3 with modular per-package task includes                                                  |
| **Dependencies**  | Dependabot (Go modules + GitHub Actions), auto-merge for patch/minor updates                        |

## Notable Libraries

| Library                       | Purpose                                                                                               |
|-------------------------------|-------------------------------------------------------------------------------------------------------|
| `golang.org/x/image`          | High-quality image scaling (bilinear interpolation) and advanced `draw.Scaler` interface for resizing |
| `golang.org/x/text`           | Character encoding conversion for international code tables (CP437, CP850, WPC1252)                   |
| `github.com/yeqown/go-qrcode` | QR code generation with error correction levels (L/M/Q/H), embedded logos, and circular block shapes  |
| `github.com/fogleman/gg`      | 2D graphics context (transitive dependency via go-qrcode) for rendering QR codes as images            |
| `github.com/stretchr/testify` | Test assertions, mocking (`testify/mock`), and table-driven test patterns across ~50 test files       |

## CV-Ready Achievements

* **Architected** a layered, modular Go library following SOLID principles (Facade, Strategy, Command, Registry
  patterns) with 12+ packages organized by domain responsibility, enabling clean separation between protocol generation,
  device profiles, connection strategies, and document execution.

* **Engineered** a custom image processing pipeline featuring Atkinson dithering, bilinear interpolation scaling, and
  monochrome bitmap generation optimized for 1-bit thermal printer output, delivering high-fidelity rendering of photos
  and logos on hardware limited to black-and-white.

* **Developed** a comprehensive ESC/POS binary protocol implementation supporting 11 command types (text, image,
  barcode, QR, table, separator, feed, cut, raw, pulse, beep) with full text formatting (bold, underline, inverse, 8
  size multipliers, Font A/B) and 8 barcode symbologies.

* **Implemented** native Windows Print Spooler integration using `syscall` and `unsafe` packages to interact directly
  with `winspool.drv`, mapping C-style structs (`DOC_INFO_1`, `PRINTER_INFO_2`) for print job management and printer
  enumeration via `EnumPrintersW` — including smart auto-detection that filters thermal/POS printers from virtual ones
  and supports proper job tracking in production POS environments.

* **Designed** a hardware-aware dynamic table engine with automatic column width reduction ("reduce longest first"
  algorithm), word wrapping, configurable spacing, and overflow protection — validating
  `Σ(columns) + (n-1) × spacing ≤ maxChars` before rendering to prevent print failures on 58mm/80mm thermal paper.

* **Built** a visual emulator (`pkg/emulator`) that renders ESC/POS documents as PNG images using a canvas-based system
  with TTF font management (JetBrains Mono), state machine tracking, and bitmap fallback — enabling development and
  automated testing without physical printer hardware.

* **Established** a production-grade CI/CD pipeline with 11 GitHub Actions workflows including multi-OS test matrix with
  race detection, automated semantic versioning with changelog generation, CodeQL and Trivy security scanning, PR
  automation with benchmark comparison, Dependabot auto-merge, and Codecov coverage reporting.

* **Implemented** a secure JSON document parsing layer with file path traversal protection, document schema validation (
  version format, paper width, DPI constraints), and a multi-strategy printer resolution system that resolves printer
  names from CLI flags → JSON profile → auto-detection with thermal-pattern matching.

## Skills Demonstrated

Go (Golang), Systems Programming, Windows API Integration (syscall/unsafe), ESC/POS Protocol Engineering, Binary
Protocol Implementation, Image Processing (Dithering, Scaling, Monochrome Conversion), JSON Schema Design & Validation,
Interface-Driven Architecture, Design Patterns (Facade, Strategy, Command, Builder, Registry), Table-Driven Testing,
CI/CD Pipeline Design (GitHub Actions), Static Analysis & Security Scanning (CodeQL, Trivy, golangci-lint), Automated
Semantic Versioning, Cross-Platform Development, Hardware Abstraction, Conventional Commits
