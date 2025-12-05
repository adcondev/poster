# LEARNING.md: POS Printer Driver & Utility

## Project Overview

**pos-printer** is a robust, production-ready Go library and command-line utility designed to bridge the gap between
modern web/backend applications and legacy ESC/POS thermal printers. Unlike simple raw-byte senders, this project
implements a full abstraction layer that decouples business logic from hardware specifics, allowing developers to define
print jobs using a clean, versioned JSON schema.

The system features a custom graphics engine for high-quality image dithering, a dynamic table layout engine, and direct
integration with the Windows Print Spooler API, making it suitable for retail and POS environments where reliability and
print quality are paramount.

## Technical Architecture

The project follows a strict **Layered Architecture** to ensure maintainability and testability:

1. **API Layer (`api/v1`)**: Defines the contract. Uses **JSON Schema** to validate print jobs before they reach the
   core logic, ensuring type safety and structural integrity.
2. **Parsing & Builder Layer (`cmd/poster`, `pkg/document/builder`)**:
    - **Builder Pattern**: Provides a fluent API (`builder.Text("...").Bold().Build()`) for programmatic document
      creation.
    - **Parser**: Marshals/unmarshals JSON documents into internal Go structs.
3. **Executor Layer (`pkg/document/executor`)**:
    - **Command Pattern**: The `Executor` iterates through the document's command list. Each command type (`text`,
      `image`, `table`) has a registered handler function, allowing for easy extension of new commands without modifying
      the core loop.
4. **Service Layer (`pkg/service`)**:
    - **Facade Pattern**: The `Printer` struct acts as a facade, coordinating the `Protocol` (command generation),
      `Profile` (device capabilities), and `Connection` (I/O).
5. **Protocol Layer (`pkg/composer`)**:
    - Responsible for generating the raw ESC/POS byte sequences. It abstracts the complexity of control codes (e.g.,
      `ESC @`, `GS V`, `ESC *`).
6. **Connection Layer (`pkg/connection`)**:
    - **Strategy Pattern**: Defines a `Connector` interface. The `WindowsPrintConnector` implements this interface using
      low-level system calls, but it can be easily swapped for `NetworkConnector` or `SerialConnector`.

## Key Technologies & Skills Demonstrated

### 1. Systems Programming with Go (Golang)

- **Windows API Integration**: Utilized `syscall` and `unsafe` packages to interact directly with `winspool.drv`.
  Implemented complex C-style struct mapping (`DOC_INFO_1`) to manage print jobs natively in Windows.
- **Interface-Driven Design**: Heavy use of interfaces (`Connector`, `PrinterService`) to decouple the OS-specific
  implementation from the business logic, facilitating mocking and testing.
- **Error Handling**: Robust error wrapping and propagation using `fmt.Errorf` and `%w` to provide meaningful context in
  stack traces.

### 2. Advanced Image Processing

- **Custom Graphics Pipeline**: Built a `graphics` package from scratch to handle image printing on thermal paper (which
  only supports black and white).
- **Dithering Algorithms**: Implemented **Atkinson Dithering** (an improvement over Floyd-Steinberg) to preserve detail
  in photos and logos when converting to 1-bit monochrome.
- **Raster Graphics**: Manually constructed ESC/POS raster bit-image commands, handling byte alignment and bit-packing
  logic.

### 3. Protocol Implementation (ESC/POS)

- **Binary Protocol Engineering**: Deep understanding of the Epson ESC/POS standard. Implemented commands for:
    - Hardware initialization and reset.
    - Text formatting (bold, invert, sizing, justification).
    - Barcode (Code128, EAN13) and QR Code generation (native firmware commands vs. software rendering fallback).
    - Paper mechanism control (full/partial cuts).

### 4. Data Structures & Algorithms

- **Dynamic Table Engine**: Created a layout engine (`pkg/tables`) that calculates column widths, handles word wrapping,
  and manages alignment for ASCII tables. This solves the common pain point of aligning receipt items on variable-width
  paper.
- **JSON Schema Validation**: Integrated strict schema validation to prevent runtime errors on the physical device.

## Key Features Implemented

### 🖨️ Native Windows Spooler Integration

Instead of treating the printer as a simple serial file, the application integrates with the Windows Print Spooler. This
allows:

- Proper job tracking in the Windows print queue.
- Support for shared network printers mapped in Windows.
- Driver-level compatibility.

### 📄 Abstract Document Format (JSON)

Defined a versioned JSON document format (`v1.0`) that describes *what* to print, not *how*.

- **Example**: `{ "type": "text", "data": { "content": "Hello", "align": "center" } }`
- This abstraction allows the backend to be agnostic of the specific printer model (e.g., switching from an Epson TM-T88
  to a generic Chinese printer only requires changing the `profile` config).

### 🖼️ Hybrid QR & Barcode Support

- **Smart Fallback**: The system checks the printer profile. If the hardware supports native QR commands (faster,
  sharper), it uses them. If not (common in cheap generic models), it automatically renders the QR code as a bitmap
  image in software and prints it as a graphic.

### 🧩 Extensible Command Registry

The `Executor` uses a map-based registry for command handlers. This makes the system "Open for Extension, Closed for
Modification"—adding a new command like `draw_line` only requires registering a new handler function, not changing the
execution engine.
