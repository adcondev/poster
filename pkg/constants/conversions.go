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

// MaxCharsForPaper calculates the maximum characters per line
// based on dots per line and font width in dots.
//
// Example:  384 dots / 12 dots per char = 32 characters for 58mm Font A
func MaxCharsForPaper(dotsPerLine, fontWidthDots int) int {
	if fontWidthDots <= 0 {
		return 0
	}
	return dotsPerLine / fontWidthDots
}

// MaxCharsForPaperFontA is a convenience function for Font A calculations.
// Font A is the standard table font (12 dots wide).
func MaxCharsForPaperFontA(dotsPerLine int) int {
	return MaxCharsForPaper(dotsPerLine, FontAWidth)
}

// MaxCharsForPaperFontB is a convenience function for Font B calculations.
// Font B is smaller (9 dots wide), allowing more characters per line.
func MaxCharsForPaperFontB(dotsPerLine int) int {
	return MaxCharsForPaper(dotsPerLine, FontBWidth)
}
