package graphics

import (
	"bytes"
	"fmt"
	"image"
	"log"

	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
	"golang.org/x/image/draw"

	"github.com/adcondev/poster/internal/load"
	posqr "github.com/adcondev/poster/pkg/commands/qrcode"
	"github.com/adcondev/poster/pkg/commands/shared"
)

const (
	// minBorderWidth es el quiet zone mínimo recomendado por el estándar QR (4 módulos)
	minBorderWidth = 4

	// maxPixelWidth es el ancho máximo soportado para impresoras térmicas de 80mm a 203 DPI
	// (80mm = 3.15in * 203 DPI = 639px, menos márgenes)
	maxPixelWidth = shared.Dpl80mm203dpi
	// minPixelWidth es el mínimo para QR Version 1 (21x21) con módulos de 3px + borders
	minPixelWidth = 87

	minGridSize = 21  // QR Version 1 (21x21 modules)
	maxGridSize = 177 // QR Version 40 (177x177 modules)

	// LogoData size multiplier limits
	minSizeMulti = 1
	maxSizeMulti = 5
)

// TODO: Check if /internals fits better for custom WriteCloser

// WriteCloser wraps bytes.Buffer to implement io.WriteCloser
// FIXME: /internals fits better for custom WriteCloser?
// TODO: Consider naming this BufferWriteCloser or similar for clarity.
type WriteCloser struct {
	buffer *bytes.Buffer
}

// NewWriteCloser creates a new WriteCloser instance
func NewWriteCloser() *WriteCloser {
	return &WriteCloser{
		buffer: bytes.NewBuffer(nil),
	}
}

// Write implements io.Writer
func (wc *WriteCloser) Write(p []byte) (n int, err error) {
	return wc.buffer.Write(p)
}

// Close implements io.Closer - NO RESET HERE!
func (wc *WriteCloser) Close() error {
	// NO hacer Reset() aquí porque perdemos los datos antes de leerlos
	// El defer en ProcessQRImage llama a Close() DESPUÉS de usar Save()
	// pero ANTES de leer los bytes con Bytes()
	// FIXME: Empty Close method might indicate a design smell or need for better abstraction.
	return nil
}

// Bytes returns the buffer contents
func (wc *WriteCloser) Bytes() []byte {
	return wc.buffer.Bytes()
}

// Reset manually resets the buffer if needed for reuse
func (wc *WriteCloser) Reset() {
	wc.buffer.Reset()
}

// QROptions contiene opciones para generar QR (nativo o imagen)
type QROptions struct {
	// === Opciones comunes ===
	Model           posqr.Model
	ErrorCorrection posqr.ErrorCorrection
	PixelWidth      int

	// === Opciones para QR como imagen ===
	LogoData    string // LogoData en base64
	CircleShape bool

	// === Campos internos calculados ===
	Qr   QrInfo
	Logo LogoInfo
}

// QrInfo contiene las dimensiones calculadas del QR
type QrInfo struct {
	gridSize       int              // Tamaño de la cuadrícula (21-177)
	moduleSize     posqr.ModuleSize // Tamaño del módulo (3-16)
	dataWidth      int              // Ancho del área de datos (px)
	borderWidth    int              // Ancho del borde (px)
	totalWidth     int              // Ancho total (px)
	requestedWidth int              // Ancho solicitado (px)
	diff           int              // Diferencia (px)
	scaleMode      string           // "exact", "scaled_up", "scaled_down"
}

// LogoInfo contiene las dimensiones calculadas del Logo
type LogoInfo struct {
	sizeMulti int         // Multiplicador de tamaño del Logo
	width     int         // Ancho del Logo (px)
	height    int         // Alto del Logo (px)
	Image     image.Image // Imagen del Logo cargada
	format    string      // Formato de la imagen (png, jpg, etc.)
}

// DefaultQROptions retorna opciones por defecto optimizadas para impresoras térmicas.
//
// PixelWidth por defecto: 288px total
//   - Para grid típico 25x25: total modules = 33, module size = 8px
//   - Área datos: 25 × 8 = 200px
//   - Quiet zone: 8 × 8 = 64px (32px por lado)
//   - Total: 264px (ajustado a múltiplo de module size)
func DefaultQROptions() *QROptions {
	return &QROptions{
		Model:           posqr.Model2,
		ErrorCorrection: posqr.LevelQ,
		PixelWidth:      128,
		LogoData:        "",
		CircleShape:     false,
		// Post-calculated fields
		Qr:   QrInfo{},
		Logo: LogoInfo{},
	}
}

// GetModuleSize retorna el tamaño del módulo calculado
func (qro *QROptions) GetModuleSize() posqr.ModuleSize {
	return qro.Qr.moduleSize
}

// GenerateQR calcula y establece el tamaño del módulo basado en PixelWidth y el tamaño de la cuadrícula del QR
// TODO: GenerateQR modifies QROptions struct (side effects). Consider returning a result struct or modified copy.
// TODO: Break down complex methods further if possible.
func (qro *QROptions) GenerateQR(data string) (*qrcode.QRCode, error) {
	if data == "" {
		return nil, fmt.Errorf("QR data cannot be empty")
	}
	if len(data) > posqr.MaxDataLength {
		return nil, fmt.Errorf("QR data too long: %d bytes (maximum %d)",
			len(data), posqr.MaxDataLength)
	}

	// FIXME: === Check if encoding affects QR Code whether is printed as image or natively ===

	// TODO: Validate/clamp PixelWidth explicitly returning errors instead of silent adjustment?
	qro.validateAndAdjustPixelWidth()

	if qro.ErrorCorrection < posqr.LevelL || qro.ErrorCorrection > posqr.LevelH {
		qro.ErrorCorrection = posqr.LevelQ
		log.Printf("QR: using default error correction level Q")
	}

	// Crear QR code
	qrc, err := qrcode.NewWith(data, mapEclOption(qro.ErrorCorrection))
	if err != nil {
		return nil, fmt.Errorf("failed to create QR code: %w", err)
	}

	qro.Qr.gridSize = qrc.Dimension()
	if err := qro.validateGridSize(); err != nil {
		return nil, err
	}

	// Calculate module size
	qro.calculateModuleSize(data)

	// Calcular y reportar tamaño final real
	qro.calculateFinalDimensions()

	// Auto-calcular LogoSizeMulti si hay Logo
	qro.loadLogo()

	return qrc, nil
}

func (qro *QROptions) validateAndAdjustPixelWidth() {
	switch {
	case qro.PixelWidth < minPixelWidth:
		log.Printf("warning: pixel_width %d < minimum %d, adjusting to minimum",
			qro.PixelWidth, minPixelWidth)
		qro.PixelWidth = minPixelWidth
	case qro.PixelWidth > maxPixelWidth:
		log.Printf("warning: pixel_width %d exceeds maximum %d, clamping",
			qro.PixelWidth, maxPixelWidth)
		qro.PixelWidth = maxPixelWidth
	default:
		qro.Qr.requestedWidth = qro.PixelWidth
		log.Printf("QR: using requested pixel_width %d", qro.PixelWidth)
	}
}

func (qro *QROptions) validateGridSize() error {
	if qro.Qr.gridSize < minGridSize {
		return fmt.Errorf("QR grid size %d is too small (minimum %d)",
			qro.Qr.gridSize, minGridSize)
	} else if qro.Qr.gridSize >= maxGridSize {
		log.Printf("warning: QR grid size %d is very large, maybe data is missing", qro.Qr.gridSize)
	}
	return nil
}

func (qro *QROptions) calculateModuleSize(data string) {
	totalModules := qro.Qr.gridSize + (2 * minBorderWidth)
	if totalModules <= 0 {
		// This case should be unreachable given minGridSize validation.
		panic(fmt.Sprintf("bug: invalid total modules after validation: %d", totalModules))
	}
	qro.Qr.moduleSize = posqr.ModuleSize(qro.Qr.requestedWidth / totalModules)

	log.Printf("QR: data='%s' (len=%d), grid=%dx%d, border=%d modules, total=%d modules",
		truncateString(data, 50), len(data), qro.Qr.gridSize, qro.Qr.gridSize, minBorderWidth, totalModules)
	log.Printf("QR: requested=%dpx, calculated_module_size=%d", qro.Qr.requestedWidth, qro.Qr.moduleSize)

	// Aplicar límites al module size
	switch {
	case qro.Qr.moduleSize < posqr.DefaultModuleSize:
		qro.Qr.moduleSize = posqr.DefaultModuleSize
		log.Printf("QR: calculated module size %d too small, using default %d",
			qro.Qr.moduleSize, posqr.DefaultModuleSize)
	case qro.Qr.moduleSize > posqr.MaxModuleSize:
		qro.Qr.moduleSize = posqr.MaxModuleSize
		log.Printf("QR: calculated module size %d too large, using maximum %d",
			qro.Qr.moduleSize, posqr.MaxModuleSize)
	default:
		log.Printf("QR: using calculated module size %d", qro.Qr.moduleSize)
	}
}

func (qro *QROptions) calculateFinalDimensions() {
	totalModules := qro.Qr.gridSize + (2 * minBorderWidth)
	qro.Qr.totalWidth = totalModules * int(qro.Qr.moduleSize)
	qro.Qr.dataWidth = qro.Qr.gridSize * int(qro.Qr.moduleSize)
	qro.Qr.borderWidth = (2 * minBorderWidth) * int(qro.Qr.moduleSize)

	// Log detallado de dimensiones
	log.Printf("QR dimensions:")
	log.Printf("  - Grid:        %dx%d modules", qro.Qr.gridSize, qro.Qr.gridSize)
	log.Printf("  - Module size: %dpx", qro.Qr.moduleSize)
	log.Printf("  - Data area:   %dx%dpx", qro.Qr.dataWidth, qro.Qr.dataWidth)
	log.Printf("  - Border:      %dpx per side (%dpx total)",
		qro.Qr.borderWidth/2, qro.Qr.borderWidth)
	log.Printf("  - Total:       %dx%dpx", qro.Qr.totalWidth, qro.Qr.totalWidth)
	log.Printf("  - Requested:   %dpx", qro.Qr.requestedWidth)

	// Advertencia si el tamaño difiere del solicitado
	if qro.Qr.totalWidth != qro.Qr.requestedWidth {
		qro.Qr.diff = qro.Qr.totalWidth - qro.Qr.requestedWidth
		switch {
		case qro.Qr.diff > 0:
			qro.Qr.scaleMode = "scaled_up"
			log.Printf("info: actual QR size %dpx exceeds requested %dpx by %dpx",
				qro.Qr.totalWidth, qro.Qr.requestedWidth, qro.Qr.diff)
		case qro.Qr.diff < 0:
			qro.Qr.scaleMode = "scaled_down"
			log.Printf("info: actual QR size %dpx is smaller than requested %dpx by %dpx",
				qro.Qr.totalWidth, qro.Qr.requestedWidth, -qro.Qr.diff)
		default:
			qro.Qr.scaleMode = "exact"
			log.Printf("info: actual QR size matches requested size exactly")
		}
	}
}

func (qro *QROptions) loadLogo() {
	if qro.LogoData != "" {
		// Load images if base64 strings are provided
		// FIXME: Define a single place where to load images from base64
		// TODO: Optimize Base64 handling? Decoding images here might be slow if repeated.
		logoImg, format, err := load.ImgFromBase64(qro.LogoData)
		if err != nil {
			log.Printf("warning: failed to load Logo: %v", err)
			qro.LogoData = ""
		} else {
			qro.Logo.Image = logoImg
			qro.Logo.format = format
			qro.Logo.sizeMulti = mapLogoSize(qro.ErrorCorrection)
			qro.Logo.width = logoImg.Bounds().Dx()
			qro.Logo.height = logoImg.Bounds().Dy()

			aspectRatio := float64(qro.Logo.width) / float64(qro.Logo.height)
			if aspectRatio < 0.8 || aspectRatio > 1.2 {
				log.Printf("warning: Logo aspect ratio %.2f not square, may cause distortion", aspectRatio)
			}

			qro.Logo.sizeMulti = mapLogoSize(qro.ErrorCorrection)
			log.Printf("QR: loaded Logo format=%s, size=%dx%d, aspect_ratio=%.2f, size_multi=%d",
				qro.Logo.format, qro.Logo.width, qro.Logo.height, aspectRatio, qro.Logo.sizeMulti)
		}
	}
}

// ProcessQRImage genera un QR code como imagen optimizada para impresora térmica
func ProcessQRImage(data string, opts *QROptions) (image.Image, error) {
	if data == "" {
		return nil, fmt.Errorf("QR data cannot be empty")
	}
	if opts == nil {
		opts = DefaultQROptions()
	}

	// El objetivo es hacer el QR tan grande y legible como sea posible
	// mientras sea más pequeño o igual al PixelWidth solicitado. Si el generado es más
	// grande que el solicitado, preferir generado y no escalar hacia abajo.

	// Generar QR y calcular tamaño del módulo
	// TODO: Error Handling: GenerateQR modifies opts. If it fails partially, opts might be inconsistent.
	qrc, err := opts.GenerateQR(data)
	if err != nil {
		log.Printf("error setting module size: %v", err)
		return nil, fmt.Errorf("set module size: %w", err)
	}
	if qrc == nil {
		return nil, fmt.Errorf("QR code generation returned nil")
	}

	// Construir opciones de imagen
	imgOpts := buildImageOptions(opts)

	// Generar imagen en memoria
	buf := NewWriteCloser()

	w := standard.NewWithWriter(buf, imgOpts...)
	defer func(w *standard.Writer) {
		err := w.Close()
		if err != nil {
			log.Printf("error closing QR writer: %v", err)
		}
	}(w)

	if err := qrc.Save(w); err != nil {
		return nil, fmt.Errorf("generate QR image: %w", err)
	}

	// Decodificar imagen
	img, format, err := image.Decode(bytes.NewReader(buf.Bytes()))
	if err != nil {
		return nil, fmt.Errorf("decode QR image: %w", err)
	}

	log.Printf("QR: decoded image format=%s", format)

	// Validar dimensiones finales
	finalBounds := img.Bounds()
	finalWidth := finalBounds.Dx()

	if finalWidth != opts.Qr.totalWidth {
		log.Printf("⚠️  Final image size %dpx differs from calculated %dpx (library adjustment)",
			finalWidth, opts.Qr.totalWidth)
	}

	return img, nil
}

// mapLogoSize calcula el tamaño óptimo del Logo basado en corrección
// TODO: Avoid magic numbers (7, 6, 5, 4). Define constants for these ratios.
func mapLogoSize(errorCorrection posqr.ErrorCorrection) int {
	switch errorCorrection {
	case posqr.LevelL:
		return 7 // 1/7 px del QR
	case posqr.LevelM:
		return 6 // 1/6 px del QR
	case posqr.LevelQ:
		return 5 // 1/5 px del QR
	case posqr.LevelH:
		return 4 // 1/4 px del QR
	default:
		return 5 // 1/5 px del QR
	}
}

// truncateString helper para logs
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// buildImageOptions construye las opciones mínimas útiles para impresora térmica
func buildImageOptions(opts *QROptions) []standard.ImageOption {
	var imgOpts []standard.ImageOption

	imgOpts = append(imgOpts, standard.WithQRWidth(uint8(opts.Qr.moduleSize))) // Module PixelWidth
	imgOpts = append(imgOpts, standard.WithBorderWidth(minBorderWidth))        // Silence Zone

	// LogoPath si está habilitado y existe
	if opts.LogoData != "" && opts.Logo.Image != nil && opts.Logo.sizeMulti > 0 && opts.Logo.format != "" {
		log.Printf("info: applying Logo options")
		// Tamaño del Logo
		switch {
		case opts.Logo.sizeMulti < minSizeMulti:
			log.Printf("warning: logo_size_multi %d out of range [%d-%d], using minimum %d",
				opts.Logo.sizeMulti, minSizeMulti, maxSizeMulti, minSizeMulti)
			opts.Logo.sizeMulti = minSizeMulti
		case opts.Logo.sizeMulti > maxSizeMulti:
			log.Printf("warning: logo_size_multi %d out of range [%d-%d], using maximum %d",
				opts.Logo.sizeMulti, minSizeMulti, maxSizeMulti, maxSizeMulti)
			opts.Logo.sizeMulti = maxSizeMulti
		default:
			log.Printf("info: using logo_size_multi %d", opts.Logo.sizeMulti)
		}

		// Zona segura para el Logo
		if opts.ErrorCorrection == posqr.LevelL || opts.ErrorCorrection == posqr.LevelM {
			imgOpts = append(imgOpts, standard.WithLogoSafeZone())
		}

		// Validar tamaño del Logo contra dataWidth
		maxLogoSize := opts.Qr.dataWidth / opts.Logo.sizeMulti
		log.Printf("QR Logo: max size for 1/%d of data area (%dpx): %dpx",
			opts.Logo.sizeMulti, opts.Qr.dataWidth, maxLogoSize)

		scaledLogo := ScaleImageToWidth(opts.Logo.Image, maxLogoSize)
		opts.Logo.Image = scaledLogo
		opts.Logo.width = scaledLogo.Bounds().Dx()
		opts.Logo.height = scaledLogo.Bounds().Dy()

		log.Printf("QR Logo: final size %dx%dpx", opts.Logo.width, opts.Logo.height)

		// Imagen del Logo
		imgOpts = append(imgOpts, standard.WithLogoSizeMultiplier(opts.Logo.sizeMulti))
		imgOpts = append(imgOpts, standard.WithLogoImage(opts.Logo.Image))
		switch opts.Logo.format {
		case "png":
			imgOpts = append(imgOpts, standard.WithBuiltinImageEncoder(standard.PNG_FORMAT))
		case "jpeg", "jpg":
			imgOpts = append(imgOpts, standard.WithBuiltinImageEncoder(standard.JPEG_FORMAT))
		}
	} else {
		opts.LogoData = ""
		opts.Logo.Image = nil
		log.Printf("info: no valid Logo provided, skipping logo options")
	}

	if opts.CircleShape {
		imgOpts = append(imgOpts, standard.WithCircleShape())
	}

	return imgOpts
}

// mapEclOption convierte el nivel de corrección de errores poster a go-qrcode
func mapEclOption(level posqr.ErrorCorrection) qrcode.EncodeOption {
	switch level {
	case posqr.LevelL:
		return qrcode.WithErrorCorrectionLevel(qrcode.ErrorCorrectionLow)
	case posqr.LevelM:
		return qrcode.WithErrorCorrectionLevel(qrcode.ErrorCorrectionMedium)
	case posqr.LevelQ:
		return qrcode.WithErrorCorrectionLevel(qrcode.ErrorCorrectionQuart)
	case posqr.LevelH:
		return qrcode.WithErrorCorrectionLevel(qrcode.ErrorCorrectionHighest)
	default:
		return qrcode.WithErrorCorrectionLevel(qrcode.ErrorCorrectionMedium)
	}
}

// ScaleImageToWidth escala una imagen al ancho objetivo manteniendo el aspect ratio.
// Usa interpolación bilinear para mejor calidad.
func ScaleImageToWidth(img image.Image, targetWidth int) image.Image {
	if img == nil {
		return nil
	}
	if targetWidth <= 0 {
		log.Printf("warning: invalid targetWidth %d, returning original", targetWidth)
		return img
	}

	bounds := img.Bounds()
	originalWidth := bounds.Dx()
	originalHeight := bounds.Dy()
	// Si ya tiene el ancho objetivo, retornar sin cambios
	if originalWidth == targetWidth {
		return img
	}

	// Calcular altura proporcional manteniendo aspect ratio
	aspectRatio := float64(originalHeight) / float64(originalWidth)
	targetHeight := int(float64(targetWidth) * aspectRatio)

	// Crear imagen de destino
	dst := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))

	// Escalar con interpolación bilinear
	draw.BiLinear.Scale(dst, dst.Bounds(), img, bounds, draw.Over, nil)

	return dst
}

/*
var ecMap = map[string]posqr.ErrorCorrection{
	constants.L.String(): posqr.LevelL,
	constants.M.String(): posqr.LevelM,
	constants.Q.String(): posqr.LevelQ,
	constants.H.String(): posqr.LevelH,
}
*/

// TODO: Add error correction mapping
