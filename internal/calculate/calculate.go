// Package calculate proporciona funciones para cálculos relacionados con la impresión.
package calculate

// DotsPerLine calcula los puntos por línea basándose en el ancho del papel en mm y la resolución en dpi.
func DotsPerLine(paperWidthMM float64, dpi int) int {
	return int((paperWidthMM * float64(dpi)) / 25.4)
}

// MmToDots convierte una medida en milímetros a puntos, según la resolución en dpi.
func MmToDots(mm float64, dpi int) int {
	return int((mm * float64(dpi)) / 25.4)
}

// DotsToMm convierte una medida en puntos a milímetros, según la resolución en dpi.
func DotsToMm(dots int, dpi int) float64 {
	return (float64(dots) * 25.4) / float64(dpi)
}
