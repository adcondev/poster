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
	"github.com/adcondev/poster/pkg/constants"
)

const (
	// minBorderWidth es el quiet zone mínimo recomendado por el estándar QR (4 módulos)
	minBorderWidth = 4
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

// QrOptions contiene opciones para generar QR (nativo o imagen)
type QrOptions struct {
	// === Opciones comunes ===
	Model           posqr.Model
	ErrorCorrection posqr.ErrorCorrection
	PixelWidth      int
	MaxPixelWidth   int // Ancho máximo permitido (override global)

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
func DefaultQROptions() *QrOptions {
	return &QrOptions{
		Model:           posqr.Model2,
		ErrorCorrection: posqr.LevelM,
		PixelWidth:      128,
		MaxPixelWidth:   constants.PaperPxWidth80mm,
		// Imagen
		LogoData:    "",
		CircleShape: false,
		// Post-calculated fields
		Qr:   QrInfo{},
		Logo: LogoInfo{},
	}
}

// GetModuleSize retorna el tamaño del módulo calculado
func (qo *QrOptions) GetModuleSize() posqr.ModuleSize {
	return qo.Qr.moduleSize
}

// GenerateQR calcula y establece el tamaño del módulo basado en PixelWidth y el tamaño de la cuadrícula del QR
// TODO: GenerateQR modifies QrOptions struct (side effects). Consider returning a result struct or modified copy.
// TODO: Break down complex methods further if possible.
func (qo *QrOptions) GenerateQR(data string) (*qrcode.QRCode, error) {
	if data == "" {
		return nil, fmt.Errorf("QR data cannot be empty")
	}
	if len(data) > posqr.MaxDataLength {
		return nil, fmt.Errorf("QR data too long: %d bytes (maximum %d)",
			len(data), posqr.MaxDataLength)
	}

	// FIXME: === Check if encoding affects QR Code whether is printed as image or natively ===

	// TODO: Validate/clamp PixelWidth explicitly returning errors instead of silent adjustment?
	qo.validateAndAdjustPixelWidth()

	if qo.ErrorCorrection < posqr.LevelL || qo.ErrorCorrection > posqr.LevelH {
		qo.ErrorCorrection = posqr.LevelM
		log.Printf("QR: using default error correction level M")
	}

	// Crear QR code
	qrc, err := qrcode.NewWith(data, mapEclOption(qo.ErrorCorrection))
	if err != nil {
		return nil, fmt.Errorf("failed to create QR code: %w", err)
	}

	qo.Qr.gridSize = qrc.Dimension()
	if err := qo.validateGridSize(); err != nil {
		return nil, err
	}

	// Calculate module size
	qo.calculateModuleSize(data)

	// Calcular y reportar tamaño final real
	qo.calculateFinalDimensions()

	// Auto-calcular LogoSizeMulti si hay Logo
	qo.loadLogo()

	return qrc, nil
}

// validateAndAdjustPixelWidth ajusta PixelWidth si está fuera de los límites permitidos
func (qo *QrOptions) validateAndAdjustPixelWidth() {
	// Use dynamic max width if set, otherwise fallback to global constant
	limit := constants.PaperPxWidth80mm
	if qo.MaxPixelWidth > 0 {
		limit = qo.MaxPixelWidth
	}

	switch {
	case qo.PixelWidth < minPixelWidth:
		log.Printf("warning: pixel_width %d < minimum %d, adjusting to minimum",
			qo.PixelWidth, minPixelWidth)
		qo.PixelWidth = minPixelWidth
		qo.Qr.requestedWidth = minPixelWidth
	case qo.PixelWidth > limit:
		log.Printf("warning: pixel_width %d exceeds limit %d, clamping",
			qo.PixelWidth, limit)
		qo.PixelWidth = limit
		qo.Qr.requestedWidth = limit
	default:
		log.Printf("QR: using requested pixel_width %d", qo.PixelWidth)
		qo.Qr.requestedWidth = qo.PixelWidth
	}
}

func (qo *QrOptions) validateGridSize() error {
	if qo.Qr.gridSize < minGridSize {
		return fmt.Errorf("QR grid size %d is too small (minimum %d)",
			qo.Qr.gridSize, minGridSize)
	} else if qo.Qr.gridSize >= maxGridSize {
		log.Printf("warning: QR grid size %d is very large, maybe data is missing", qo.Qr.gridSize)
	}
	return nil
}

func (qo *QrOptions) calculateModuleSize(data string) {
	totalModules := qo.Qr.gridSize + (2 * minBorderWidth)
	if totalModules <= 0 {
		// This case should be unreachable given minGridSize validation.
		panic(fmt.Sprintf("bug: invalid total modules after validation: %d", totalModules))
	}
	qo.Qr.moduleSize = posqr.ModuleSize(qo.Qr.requestedWidth / totalModules)

	log.Printf("QR: data='%s' (len=%d), grid=%dx%d, border=%d modules, total=%d modules",
		truncateString(data, 50), len(data), qo.Qr.gridSize, qo.Qr.gridSize, minBorderWidth, totalModules)
	log.Printf("QR: requested=%dpx, calculated_module_size=%d", qo.Qr.requestedWidth, qo.Qr.moduleSize)

	// Aplicar límites al module size
	switch {
	case qo.Qr.moduleSize < posqr.DefaultModuleSize:
		qo.Qr.moduleSize = posqr.DefaultModuleSize
		log.Printf("QR: calculated module size %d too small, using default %d",
			qo.Qr.moduleSize, posqr.DefaultModuleSize)
	case qo.Qr.moduleSize > posqr.MaxModuleSize:
		qo.Qr.moduleSize = posqr.MaxModuleSize
		log.Printf("QR: calculated module size %d too large, using maximum %d",
			qo.Qr.moduleSize, posqr.MaxModuleSize)
	default:
		log.Printf("QR: using calculated module size %d", qo.Qr.moduleSize)
	}
}

func (qo *QrOptions) calculateFinalDimensions() {
	totalModules := qo.Qr.gridSize + (2 * minBorderWidth)
	qo.Qr.totalWidth = totalModules * int(qo.Qr.moduleSize)
	qo.Qr.dataWidth = qo.Qr.gridSize * int(qo.Qr.moduleSize)
	qo.Qr.borderWidth = (2 * minBorderWidth) * int(qo.Qr.moduleSize)

	// Log detallado de dimensiones
	log.Printf("QR dimensions:")
	log.Printf("  - Grid:        %dx%d modules", qo.Qr.gridSize, qo.Qr.gridSize)
	log.Printf("  - Module size: %dpx", qo.Qr.moduleSize)
	log.Printf("  - Data area:   %dx%dpx", qo.Qr.dataWidth, qo.Qr.dataWidth)
	log.Printf("  - Border:      %dpx per side (%dpx total)",
		qo.Qr.borderWidth/2, qo.Qr.borderWidth)
	log.Printf("  - Total:       %dx%dpx", qo.Qr.totalWidth, qo.Qr.totalWidth)
	log.Printf("  - Requested:   %dpx", qo.Qr.requestedWidth)

	// Advertencia si el tamaño difiere del solicitado
	if qo.Qr.totalWidth != qo.Qr.requestedWidth {
		qo.Qr.diff = qo.Qr.totalWidth - qo.Qr.requestedWidth
		switch {
		case qo.Qr.diff > 0:
			qo.Qr.scaleMode = "scaled_up"
			log.Printf("info: actual QR size %dpx exceeds requested %dpx by %dpx",
				qo.Qr.totalWidth, qo.Qr.requestedWidth, qo.Qr.diff)
		case qo.Qr.diff < 0:
			qo.Qr.scaleMode = "scaled_down"
			log.Printf("info: actual QR size %dpx is smaller than requested %dpx by %dpx",
				qo.Qr.totalWidth, qo.Qr.requestedWidth, -qo.Qr.diff)
		default:
			qo.Qr.scaleMode = "exact"
			log.Printf("info: actual QR size matches requested size exactly")
		}
	}
}

func (qo *QrOptions) loadLogo() {
	if qo.LogoData != "" {
		// Load images if base64 strings are provided
		// FIXME: Define a single place where to load images from base64
		// TODO: Optimize Base64 handling? Decoding images here might be slow if repeated.
		logoImg, format, err := load.ImgFromBase64(qo.LogoData)
		if err != nil {
			log.Printf("warning: failed to load Logo: %v", err)
			qo.LogoData = ""
		} else {
			qo.Logo.Image = logoImg
			qo.Logo.format = format
			qo.Logo.sizeMulti = mapLogoSize(qo.ErrorCorrection)
			qo.Logo.width = logoImg.Bounds().Dx()
			qo.Logo.height = logoImg.Bounds().Dy()

			aspectRatio := float64(qo.Logo.width) / float64(qo.Logo.height)
			if aspectRatio < 0.8 || aspectRatio > 1.2 {
				log.Printf("warning: Logo aspect ratio %.2f not square, may cause distortion", aspectRatio)
			}

			qo.Logo.sizeMulti = mapLogoSize(qo.ErrorCorrection)
			log.Printf("QR: loaded Logo format=%s, size=%dx%d, aspect_ratio=%.2f, size_multi=%d",
				qo.Logo.format, qo.Logo.width, qo.Logo.height, aspectRatio, qo.Logo.sizeMulti)
		}
	}
}

// ProcessQRImage genera un QR code como imagen optimizada para impresora térmica
func ProcessQRImage(data string, opts *QrOptions) (image.Image, error) {
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
func buildImageOptions(opts *QrOptions) []standard.ImageOption {
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
