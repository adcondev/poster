# Agent Guidelines: Poster (`adcondev/poster`)

Operational context, architectural invariants, and strict guardrails for AI agents and contributors working in this repository.

## 1. Project Overview & Tech Stack

Poster is a Go library and CLI utility for controlling ESC/POS thermal receipt printers via a versioned JSON document protocol (`v1.0`), programmatic Go builders, pluggable `io.WriteCloser` transports (Windows Spooler, raw TCP, file/device, in-memory buffer), a multi-algorithm dithering graphics pipeline, an auto-balancing pure-text table engine, multi-codepage text encoding, and a 2D PNG emulator that can render both documents (via `PrinterActions`) and raw ESC/POS byte streams.

- **Language / Runtime**: Go `1.24.6+` (`module github.com/adcondev/poster`)
- **Target Platforms**:
  - **Windows (`//go:build windows`)**: Native Spooler printing and printer discovery via `winspool.drv`.
  - **Cross-Platform (`linux`, `darwin`, `windows`)**: Everything else — command builders, composer, schema/builder/executor, graphics, tables, profiles, emulator, and `TCPConnector` / `FileConnector` / `BufferConnector`.
- **Release & Commit Tooling**: Node.js (`standard-version`, `@commitlint/cli`, `prettier`), Go `Taskfile.yml`, `golangci-lint` v2, and GitHub Actions.

---

## 2. Package Map & Responsibilities

| Path | Responsibility |
|---|---|
| `api/v1/` | **Canonical Contract**: `ticket_document.schema.json` (JSON Schema Draft-07) and `TICKET_DOCUMENT_V1.md`. |
| `cmd/poster/` | CLI binary (`main.go`, `parsing.go`, `printing.go` [`windows`], `printing_stub.go` [`!windows`]). Flags: `-t windows\|tcp\|network\|file`, `-network`, `-o/-output`, `-preview`, `--dry-run`, `--list*`. Version injected via `-ldflags "-X main.AppVersion=..."`. |
| `internal/calculate/` | Pure math conversions between millimeters, DPI, and dots (`DotsPerLine`, `MmToDots`, `DotsToMm`). |
| `internal/load/` | Path-traversal-safe file loader (`SecureFilepath`) and image decoder (`PNG`, `JPEG`, `BMP` via `golang.org/x/image/bmp`). |
| `internal/testutils/` | Shared test assertions, byte slice validators, and fake/mock helpers. |
| `pkg/commands/*` | Stateless, pure ESC/POS byte generators (`barcode`, `bitimage`, `character`, `linespacing`, `mechanismcontrol`, `print`, `printposition`, `qrcode` [QR + PDF417], `shared`). Each subpackage exposes a `Capability` interface and `NewCommands()`. |
| `pkg/composer/` | High-level protocol facade (`EscposProtocol`) aggregating all `pkg/commands/*` capabilities into atomic ESC/POS sequences. |
| `pkg/connection/` | `Connector` (`io.WriteCloser`) implementations: `WindowsPrintConnector` (`winspool.drv`), `TCPConnector`, `FileConnector`, `BufferConnector`; cross-platform printer filtering + Windows enumeration. |
| `pkg/constants/` | Centralized typed string enums, paper/dot geometry constants (`576` dots for 80mm, `384` dots for 58mm), size limits (`MaxJSONSize`), and schema-aligned defaults. |
| `pkg/document/schema/` | Top-level JSON AST (`Document`, `ProfileConfig`, `Command`, per-command structs), `ParseDocument`, and structural/version validation (`Validate`). |
| `pkg/document/builder/` | Fluent Go DSL (`DocumentBuilder`) for programmatic construction of `schema.Document` payloads. |
| `pkg/document/executor/` | Document execution engine (`Executor`, `CommandHandler` registry, `RegisterHandler`) translating `schema.Document` commands into `service.PrinterActions` calls. |
| `pkg/emulator/` | 2D receipt renderer (`Engine`, `DynamicCanvas`, embedded TTF fonts), `emulator.Printer` (`service.PrinterActions` adapter), and raw byte parser (`Engine.ParseESC`, `RenderBytesToPNG`). |
| `pkg/graphics/` | Image loading/scaling, `Atkinson`/`FloydSteinberg`/`Ordered`/`Threshold` dithering into 1-bit `MonochromeBitmap` (+ `SubBand`), QR image generation (`yeqown/go-qrcode/v2`), and barcode symbology mapping. |
| `pkg/profile/` | Hardware profiles + model presets (`ProfileFromModel`), 25 code tables (`ParseCodeTable`), `NormalizeTypography`, and `EncodeBytesSafe` (multi-codepage `ESC t n` switching, `'?'` fallback, `StrictEncoding`). |
| `pkg/service/` | Stateful printer orchestrator (`Printer`, `NewPrinterWithOptions`), `PrinterActions` interface (including `Close() error`), and `MockPrinter`. |
| `pkg/tables/` | Pure-text multi-column layout engine (`TabEngine`, `Render`, `RenderStructured`) — no ESC/POS bytes; UTF-8 rune word-wrapping, alignment padding, widest-column-first `ReduceToFit`. |

---

## 3. Core Architectural Invariants

1. **Canonical Schema Truth (`api/v1/ticket_document.schema.json`)**:
   - `api/v1/ticket_document.schema.json` is the single source of truth for command fields, defaults, and enums.
   - `pkg/document/builder` and `pkg/document/executor` MUST share canonical command structs (no duplicate private structs in `builder`).
   - Optional JSON fields whose schema default is `true` or non-zero (e.g., `table` booleans `show_headers`, `header_bold`, `word_wrap`, `auto_reduce`, `column_spacing`, and `cut.feed`) MUST use pointer types (`*bool`, `*int`) so explicit `false` or `0` values in JSON are never overwritten by Go zero-value defaults.
   - Both `"nns"` (schema enum) and `"nearest_neighbor"` must be accepted for nearest-neighbor image scaling.
2. **Strict Error Handling & Silent Libraries**:
   - Never call `panic`, `log.Panicf`, `log.Fatalf`, or `os.Exit` anywhere in `pkg/`, `internal/`, or `cmd/` (except a single exit code return at the top of `cmd/poster/main.go`).
   - Always check and return errors wrapped with context: `fmt.Errorf("context: %w", err)`. Never discard errors with `_` unless provably infallible and documented.
   - Library packages (`pkg/`, `internal/`) MUST NOT emit noisy debug logs via global `log.Printf`.
   - `executor.Execute` MUST fail fast (return an error) on unsupported `doc.Version`, unknown command types, invalid command payloads, or profile configuration errors.
3. **Hardware, Service & Emulator Boundaries**:
   - `executor.Executor` and all `CommandHandler` implementations MUST depend on the `service.PrinterActions` interface (including `Close() error`), never the concrete `*service.Printer` struct.
   - `pkg/emulator` is a preview engine, not a `connection.Connector`. It integrates via `emulator.Printer` (implements `service.PrinterActions`) and `ParseESC` / `RenderBytesToPNG` (consumes raw bytes). `ParseESC` MUST never panic or index out of bounds on truncated/malformed input.
   - `pkg/connection` transports are `io.WriteCloser` only (Windows Spooler, TCP, file/device, buffer). OS-specific files (`//go:build windows` vs `//go:build !windows`) MUST export identical struct field types and signatures so cross-platform builds and tests never diverge.
   - `pkg/tables` MUST remain pure text (no `pkg/composer` import, no raw ESC/POS bytes). Styling (e.g., bold headers) is applied by callers via `RenderStructured` + `PrinterActions`.
4. **Standard Paper & Dot Geometry (203 DPI / 8 dpmm)**:
   - **80mm paper** (`72mm` printable width): `576` dots per line (`48` Font A chars, `64` Font B chars).
   - **58mm paper** (`48mm` printable width): `384` dots per line (`32` Font A chars, `42` Font B chars).
   - All packages (`pkg/constants`, `pkg/profile`, `pkg/commands/shared`, `pkg/document/executor`, `internal/calculate`) MUST resolve to these exact standard dot widths. Image/QR max widths MUST derive from `Profile.DotsPerLine`, never a hardcoded 80mm constant.
5. **QR, PDF417 & Barcode Rendering Policy**:
   - **QR Codes (`service.PrintQR`)**: Use native ESC/POS (`GS ( k`) only when `Profile.HasQR == true` AND neither `logo` nor `circle_shape` is enabled. Whenever `logo` or `circle_shape` is requested (or `Profile.HasQR == false`), automatically render via the software bitmap pipeline (`pkg/graphics`).
   - **PDF417**: Native `GS ( k` (`cn=48`) only; no software fallback yet.
   - **1D Barcodes (`service.PrintBarcode`)**: Use native ESC/POS (`GS k`) hardware commands. Do not claim software rasterization for 1D barcodes on hardware unless implemented in `pkg/graphics` (the emulator renders them in software for previews only).
6. **Text Encoding Policy**:
   - `service.Printer.Print` / `PrintLine` MUST encode via `Profile.EncodeBytesSafe` (typography normalization → active code table → optional `ESC t n` auto-switch → `'?'` fallback, or error when `StrictEncoding`).
   - New code tables MUST be added to `profile.ParseCodeTable` and `constants.ValidCodeTables` together.
7. **Resource Limits**: JSON input ≤ `constants.MaxJSONSize` (10MB); image height ≤ `graphics.MaxImagePixelHeight` (32000); tall bitmaps are band-split when `Profile.MaxBandHeight > 0`.

---

## 4. Development, Testing & Linting Commands

Always run verification across **all** packages (`./...`, including `cmd/`, `internal/`, and `pkg/`):

```bash
# Run full test suite with race detector and atomic coverage
go test -race -coverprofile=coverage.txt -covermode=atomic ./...

# Run golangci-lint (v2 config in .golangci.yml)
golangci-lint run --timeout=5m

# Build CLI with version injection
go build -ldflags "-X main.AppVersion=$(git describe --tags --always --dirty)" -o ./bin/poster ./cmd/poster

# Taskfile shortcuts (must use portable PATH binaries, never hardcoded user paths)
task test
task lint
task build
```

---

## 5. Commit & PR Governance (`.commitlintrc.json`)

All commits and PR titles MUST follow Conventional Commits:
`<type>(<scope>): <lowercase subject without trailing period>` (header max 100 chars).

- **Allowed `type`**: `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `build`, `ci`, `chore`, `revert`, `deps`.
- **Required `scope`**: `barcode`, `bitimage`, `builder`, `calculate`, `character`, `ci`, `commands`, `common`, `composer`, `connection`, `constants`, `docs`, `document`, `emulator`, `executor`, `gh-actions`, `github`, `gomod`, `graphics`, `linespacing`, `load`, `mechanismcontrol`, `miscellaneous`, `models`, `npm`, `poster`, `print`, `printposition`, `profile`, `qrcode`, `release`, `schema`, `security`, `service`, `shared`, `tables`.

---

## 6. Hard "Must Not Do" Boundaries (Context & Safety)

1. **Token Economy**: NEVER `cat` or read binary files, images (`assets/images/poster.png`), large base64/text dumps (`assets/images/imagen.txt`, `assets/images/logo.txt`), or full `CHANGELOG.md` (use `head -n 40 CHANGELOG.md`). Prefer `go doc`, `grep -E "^(type|func) "`, and targeted line ranges.
2. **No Unsafe Slice Indexing**: Never index `&data[0]` or slice `s[:n-1]` without explicit length/bounds checks (`len(data) == 0`).
3. **No Hardcoded Local Paths or Versions**: Never commit machine-specific paths in `Taskfile.yml` or hardcoded version strings in `cmd/poster/main.go`.
4. **No Unimplemented Transport Claims**: Do not document ZPL, Scale protocols, Serial COM, Bluetooth, USB-native (libusb), or bidirectional status (`DLE EOT` / ASB) as active features until implemented and tested. (Raw TCP/9100 and file/device output ARE implemented.)
5. **Coverage Floor**: Do not merge changes that drop any library package below ~85% statement coverage without an explicit, documented reason.
