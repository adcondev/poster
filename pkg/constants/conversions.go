package constants

// ============================================================================
// Helper Functions
// ============================================================================

// PaperWidthFromMm calculates pixel width from millimeters (printable area)
func PaperWidthFromMm(mm int) int {
	return mm * DefaultDotsPerMm
}

// CharsPerLineForFont calculates characters per line for a given paper width and font
func CharsPerLineForFont(paperWidthPx, fontWidth int) int {
	if fontWidth <= 0 {
		return 0
	}
	return paperWidthPx / fontWidth
}
