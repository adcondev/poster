# Poster - Thermal Printer Driver & Utility

## Project Overview

**Poster** is a robust, production-ready Go library and command-line utility designed to bridge the gap between modern web/backend applications and legacy ESC/POS thermal printers. Unlike simple raw-byte senders, this project implements a full abstraction layer that decouples business logic from hardware specifics, allowing developers to define print jobs using a clean, versioned JSON schema.

The system features a modular command architecture, a custom graphics engine for high-quality image dithering, a dynamic table layout engine, a visual emulator for testing, and direct integration with the Windows Print Spooler API, making it suitable for retail and POS environments where reliability and print quality are paramount.

## Technical Architecture

The project follows a strict **Layered Architecture** to ensure maintainability and testability:

### 1. API Layer (`api/v1`)

Defines the contract. Uses **JSON Schema** to validate print jobs before they reach the core logic, ensuring type safety and structural integrity.

### 2. Command Layer (`pkg/commands`)

Organizes ESC/POS commands into logical modules:
- **`barcode`** - Barcode generation commands
- **`bitimage`** - Raster image commands
- **`character`** - Character styling and selection
- **`linespacing`** - Line spacing control
- **`mechanismcontrol`** - Paper cutting management
- **`print`** - Print mode commands
- **`printposition`** - Alignment and positioning
- **`qrcode`** - QR code generation

### 3. Parsing & Builder Layer (`pkg/document/builder`)

- **Builder Pattern**: Provides a fluent API (`builder.Text("...").Bold().Build()`) for programmatic document creation.
- Separate builders for each command type: text, image, barcode, QR, table, raw, and basic commands.

### 4. Schema Layer (`pkg/document/schema`)

Defines Go structures that map directly to the JSON document format with validation.

### 5. Executor Layer (`pkg/document/executor`)

- **Command Pattern**: The `Executor` iterates through the document's command list. Each command type (`text`, `image`, `table`, `barcode`, `qr`, `raw`) has a registered handler function.
- **Registry Pattern**: Extensible command registry allows adding new commands without modifying the core execution loop.
- Includes text styling appliers and a text state manager for complex formatting.

### 6. Service Layer (`pkg/service`)

- **Facade Pattern**: The `Printer` struct acts as a facade, coordinating the `Protocol` (command generation), `Profile` (device capabilities), and `Connection` (I/O).
- Clean interface (`PrinterService`) for easy mocking in tests.

### 7. Protocol Layer (`pkg/composer`)

- `EscposComposer`: Responsible for generating the raw ESC/POS byte sequences.
- Abstracts the complexity of control codes (e.g., `ESC @`, `GS V`, `ESC *`).

### 8. Connection Layer (`pkg/connection`)

- **Strategy Pattern**: Defines a `Connector` interface with multiple implementations:
  - `WindowsPrintConnector` - Windows Print Spooler integration
  - Network connector - Raw TCP/9100 connections
  - Serial connector - COM port communication
  - File connector - Output to file for debugging

### 9. Emulator Layer (`pkg/emulator`)

- **Visual Emulator**: Renders print jobs as PNG images for preview and testing.
- Canvas-based rendering with font support.
- Handles text, images, barcodes, QR codes, and basic formatting commands.

## Key Technologies & Skills Demonstrated

### 1. Systems Programming with Go (Golang)

- **Windows API Integration**: Utilized `syscall` and `unsafe` packages to interact directly with `winspool.drv`. Implemented complex C-style struct mapping (`DOC_INFO_1`) to manage print jobs natively in Windows.
- **Interface-Driven Design**: Heavy use of interfaces (`Connector`, `PrinterService`) to decouple OS-specific implementations from business logic, facilitating mocking and testing.
- **Error Handling**: Robust error wrapping and propagation using `fmt.Errorf` and `%w` to provide meaningful context in stack traces.

### 2. Advanced Image Processing

- **Custom Graphics Pipeline**: Built a `graphics` package from scratch to handle image printing on thermal paper (which only supports black and white).
- **Dithering Algorithms**: Implemented **Atkinson Dithering** (an improvement over Floyd-Steinberg) to preserve detail in photos and logos when converting to 1-bit monochrome.
- **Scaling Algorithms**: Bilinear interpolation and nearest-neighbor scaling for image resizing.
- **Raster Graphics**: Manually constructed ESC/POS raster bit-image commands, handling byte alignment and bit-packing logic.

### 3. Protocol Implementation (ESC/POS)

- **Binary Protocol Engineering**: Deep understanding of the Epson ESC/POS standard. Implemented commands for:
    - Hardware initialization and reset
    - Text formatting (bold, underline, inverse, sizing 1x1 to 8x8, justification)
    - Font selection (Font A/B)
    - Barcode generation (Code128, EAN13, EAN8, UPC-A, UPC-E, CODE39, ITF, CODABAR)
    - QR Code generation (native firmware commands vs. software rendering fallback)
    - Image printing with multiple raster modes
    - Paper mechanism control (full/partial cuts)
    - Raw byte passthrough for custom commands

### 4. Data Structures & Algorithms

- **Dynamic Table Engine**: Created a layout engine (`pkg/tables`) that:
  - **Hardware-aware validation**: Calculates max characters from `DotsPerLine / FontWidth`
  - **Auto-reduction algorithm**: "Reduce longest first" strategy that shrinks oversized columns while preserving small
    ones (ID, Qty columns stay intact)
  - **Overflow protection**: Rejects tables exceeding paper width with clear error messages

- **JSON Schema Validation**: Integrated strict schema validation to prevent runtime errors on the physical device.
- **Registry Pattern**: Map-based command registry for extensible command handling.

### 5. Visual Rendering System

- **Canvas-Based Rendering**: The emulator uses a canvas system to compose print output.
- **Font Management**: Handles multiple fonts with fallback support.
- **State Machine**: Tracks printer state (text styles, position, alignment) for accurate emulation.

## Key Features Implemented

### 🖨️ Native Windows Spooler Integration

Instead of treating the printer as a simple serial file, the application integrates with the Windows Print Spooler. This allows:
- Proper job tracking in the Windows print queue
- Support for shared network printers mapped in Windows
- Driver-level compatibility

### 📄 Abstract Document Format (JSON v1.0)

Defined a versioned JSON document format that describes *what* to print, not *how*:
- **9 Command Types**: text, image, barcode, qr, table, separator, feed, cut, raw
- Profile configuration for printer capabilities (paper width, DPI, QR support, code tables)
- Debug logging support
- This abstraction allows the backend to be agnostic of the specific printer model

### 🖼️ Hybrid QR & Barcode Support

- **Smart Fallback**: The system checks the printer profile. If the hardware supports native QR commands (faster, sharper), it uses them. If not (common in cheap generic models), it automatically renders the QR code as a bitmap image in software.
- **Logo Support**: QR codes can include embedded logos.
- **Multiple Symbologies**: Full support for CODE128, EAN13, EAN8, UPC-A, UPC-E, CODE39, ITF, CODABAR.

### 🧩 Extensible Command Registry

The `Executor` uses a map-based registry for command handlers. This makes the system "Open for Extension, Closed for Modification"—adding a new command type only requires registering a new handler function, not changing the execution engine.

### 🎨 Visual Emulator

- Preview print jobs as PNG images without physical hardware
- Accurate rendering of text, images, barcodes, and QR codes
- Useful for development, testing, and documentation

### 🔧 Multiple Connection Types

- **Windows Spooler**: Primary production connection
- **Network**: Raw TCP/9100 for direct network printers
- **Serial**: COM port support for USB/Serial printers
- **File**: Output capture for debugging and emulator input

### 📊 Smart Table Engine

The table system goes beyond simple formatting to provide hardware-aware rendering:

- **Automatic Width Calculation**: Derives character limits from printer profile (`DotsPerLine / FontAWidth`)
- **Overflow Protection**: Validates `Σ(columns) + (n-1) × spacing ≤ maxChars` before rendering
- **Auto-Reduction**: When tables exceed paper width, automatically shrinks the widest columns first
  - Preserves small columns (like "#" or "Qty") at minimum width (default:  3 chars)
  - Configurable via `auto_reduce` option in JSON
- **Font A Enforcement**: Consistent 12-dot character width for predictable layouts
- **Paper Support**: 58mm (32 chars) and 80mm (48 chars) at 203 DPI

### 🧪 Quality Assurance & Testing

- **Table-Driven Testing**: Implemented comprehensive table-driven unit tests for the `TableEngine` to verify rendering logic, ensuring support for:
  - Header visibility toggling (explicit vs. data-driven).
  - Column alignment (Left, Center, Right) with precise padding calculations.
  - Word wrapping functionality for long text.
  - Edge cases like nil data, invalid definitions, and empty rows.
- **Mocking & Buffering**: Utilized `bytes.Buffer` as an `io.Writer` to capture and inspect ESC/POS command output without physical hardware, enabling deterministic verification of control codes (e.g., bold toggling `ESC E`).

```go
// Example:  Table that would overflow gets auto-reduced
// Original: [20, 5, 12] = 39 chars (exceeds 32 max for 58mm)
// After:     [13, 5, 12] = 32 chars (fits perfectly)
// Log:  "Table auto-reduced:  39 → 32 chars (7 reductions applied)"
```

## Module Dependencies

```
github.com/adcondev/poster
├── github.com/stretchr/testify (testing)
├── github.com/yeqown/go-qrcode/v2 (QR generation)
├── github.com/fogleman/gg (2D graphics - via go-qrcode)
├── golang.org/x/image (image processing)
└── golang.org/x/text (text encoding)
```
