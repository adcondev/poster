package bitimage

import (
	"bytes"
	"testing"

	"github.com/adcondev/poster/internal/testutils"
	"github.com/adcondev/poster/pkg/commands/common"
)

func TestIntegration_Graphics_CompleteWorkflow(t *testing.T) {
	cmd := NewGraphicsCommands()

	t.Run("raster graphics printing workflow", func(t *testing.T) {
		var buffer []byte

		// Step 1: Set dot density to 360x360 for high quality
		densityCmd, err := cmd.SetGraphicsDotDensity(
			FunctionCodeDensity49,
			Density360x360,
			Density360x360,
		)
		if err != nil {
			t.Fatalf("SetGraphicsDotDensity failed: %v", err)
		}
		buffer = append(buffer, densityCmd...)

		// Step 2: Store monochrome raster graphics
		width := uint16(400)
		height := uint16(200)
		widthBytes := (int(width) + 7) / 8
		dataSize := widthBytes * int(height)
		rasterData := testutils.RepeatByte(dataSize, 0xF0)

		storeCmd, err := cmd.StoreRasterGraphicsInBuffer(
			Monochrome,
			NormalScale,
			NormalScale,
			Color1,
			width,
			height,
			rasterData,
		)
		if err != nil {
			t.Fatalf("StoreRasterGraphicsInBuffer failed: %v", err)
		}
		buffer = append(buffer, storeCmd...)

		// Step 3: Print the graphics
		printCmd, err := cmd.PrintBufferedGraphics(FunctionCodePrint50)
		if err != nil {
			t.Fatalf("PrintBufferedGraphics failed: %v", err)
		}
		buffer = append(buffer, printCmd...)

		// Verify command structure
		if !bytes.Contains(buffer, []byte{common.GS, '(', 'L'}) {
			t.Error("Buffer should contain graphics commands")
		}

		expectedMinSize := 9 + dataSize + 15 + 7 // density + store + print
		if len(buffer) < expectedMinSize {
			t.Errorf("Buffer size = %d, expected at least %d", len(buffer), expectedMinSize)
		}
	})

	t.Run("column graphics workflow", func(t *testing.T) {
		var buffer []byte

		// Step 1: Set density to 180x180
		densityCmd, err := cmd.SetGraphicsDotDensity(
			FunctionCodeDensity1,
			Density180x180,
			Density180x180,
		)
		if err != nil {
			t.Fatalf("SetGraphicsDotDensity failed: %v", err)
		}
		buffer = append(buffer, densityCmd...)

		// Step 2: Store column graphics with scaling
		width := uint16(512)
		height := uint16(128)
		heightBytes := (int(height) + 7) / 8
		dataSize := int(width) * heightBytes
		columnData := testutils.RepeatByte(dataSize, 0xAA)

		storeCmd, err := cmd.StoreColumnGraphicsInBuffer(
			DoubleScale,
			NormalScale,
			Color2,
			width,
			height,
			columnData,
		)
		if err != nil {
			t.Fatalf("StoreColumnGraphicsInBuffer failed: %v", err)
		}
		buffer = append(buffer, storeCmd...)

		// Step 3: Print
		printCmd, err := cmd.PrintBufferedGraphics(FunctionCodePrint2)
		if err != nil {
			t.Fatalf("PrintBufferedGraphics failed: %v", err)
		}
		buffer = append(buffer, printCmd...)

		if len(buffer) != 9+dataSize+15+7 {
			t.Errorf("Buffer length = %d, expected %d", len(buffer), 9+dataSize+15+7)
		}
	})

	t.Run("multi-tone graphics with color overlay", func(t *testing.T) {
		var buffer []byte

		width := uint16(200)
		height := uint16(150)
		widthBytes := (int(width) + 7) / 8
		dataSize := widthBytes * int(height)

		// Store multiple colors for overlay effect
		colors := []GraphicsColor{
			Color1,
			Color2,
			Color3,
			Color4,
		}

		for i, color := range colors {
			// Create different patterns for each color
			pattern := byte(0x11 << uint(i)) //nolint:gosec
			data := testutils.RepeatByte(dataSize, pattern)

			storeCmd, err := cmd.StoreRasterGraphicsInBuffer(
				MultipleTone,
				NormalScale,
				NormalScale,
				color,
				width,
				height,
				data,
			)
			if err != nil {
				t.Fatalf("StoreRasterGraphicsInBuffer for color %v failed: %v", color, err)
			}
			buffer = append(buffer, storeCmd...)
		}

		// Print combined graphics
		printCmd, err := cmd.PrintBufferedGraphics(FunctionCodePrint50)
		if err != nil {
			t.Fatalf("PrintBufferedGraphics failed: %v", err)
		}
		buffer = append(buffer, printCmd...)

		expectedSize := len(colors)*(dataSize+15) + 7
		if len(buffer) != expectedSize {
			t.Errorf("Buffer size = %d, expected %d", len(buffer), expectedSize)
		}
	})

	t.Run("scaled graphics combinations", func(t *testing.T) {
		scaleCombinations := []struct {
			name       string
			horizontal GraphicsScale
			vertical   GraphicsScale
			tone       GraphicsTone
			maxHeight  uint16
		}{
			{"monochrome normal", NormalScale, NormalScale, Monochrome, 2400},
			{"monochrome double width", DoubleScale, NormalScale, Monochrome, 2400},
			{"monochrome double height", NormalScale, DoubleScale, Monochrome, 1200},
			{"monochrome quadruple", DoubleScale, DoubleScale, Monochrome, 1200},
			{"multi-tone normal", NormalScale, NormalScale, MultipleTone, 600},
			{"multi-tone double", DoubleScale, DoubleScale, MultipleTone, 300},
		}

		for _, sc := range scaleCombinations {
			t.Run(sc.name, func(t *testing.T) {
				width := uint16(100)
				height := uint16(100)
				if height > sc.maxHeight {
					height = sc.maxHeight
				}

				widthBytes := (int(width) + 7) / 8
				dataSize := widthBytes * int(height)
				data := testutils.RepeatByte(dataSize, 0xFF)

				storeCmd, err := cmd.StoreRasterGraphicsInBuffer(
					sc.tone,
					sc.horizontal,
					sc.vertical,
					Color1,
					width,
					height,
					data,
				)
				if err != nil {
					t.Fatalf("StoreRasterGraphicsInBuffer failed for %s: %v", sc.name, err)
				}

				// Verify scaling parameters are in command
				if storeCmd[8] != byte(sc.horizontal) {
					t.Errorf("Horizontal scale = %d, want %d", storeCmd[8], sc.horizontal)
				}
				if storeCmd[9] != byte(sc.vertical) {
					t.Errorf("Vertical scale = %d, want %d", storeCmd[9], sc.vertical)
				}
			})
		}
	})
}

func TestIntegration_Graphics_LargeDataHandling(t *testing.T) {
	cmd := NewGraphicsCommands()

	t.Run("large raster data requiring extended format", func(t *testing.T) {
		// Create data exceeding standard command size limit
		width := uint16(2400)
		height := uint16(2000)
		widthBytes := (int(width) + 7) / 8
		dataSize := widthBytes * int(height)

		if dataSize <= 65535 {
			t.Skip("Data size not large enough for extended format testutils")
		}

		largeData := testutils.RepeatByte(dataSize, 0x55)

		storeCmd, err := cmd.StoreRasterGraphicsInBufferLarge(
			Monochrome,
			NormalScale,
			NormalScale,
			Color1,
			width,
			height,
			largeData,
		)
		if err != nil {
			t.Fatalf("StoreRasterGraphicsInBufferLarge failed: %v", err)
		}

		// Verify extended format
		if storeCmd[0] != common.GS || storeCmd[1] != '8' || storeCmd[2] != 'L' {
			t.Error("Large data should use GS 8 L extended format")
		}

		// Verify 32-bit size parameters
		totalSize := uint32(11 + dataSize) //nolint:gosec
		p1 := storeCmd[3]
		p2 := storeCmd[4]
		p3 := storeCmd[5]
		p4 := storeCmd[6]

		calculatedSize := uint32(p1) + uint32(p2)<<8 + uint32(p3)<<16 + uint32(p4)<<24
		if calculatedSize != totalSize {
			t.Errorf("Size encoding = %d, want %d", calculatedSize, totalSize)
		}
	})

	t.Run("large column data requiring extended format", func(t *testing.T) {
		// Column format max width is 2048, so use maximum dimensions to get large data
		width := uint16(2048)
		height := uint16(128)
		heightBytes := (int(height) + 7) / 8
		dataSize := int(width) * heightBytes

		// Verify we need extended format (dataSize + 11 header bytes > 65535)
		if dataSize+11 <= 65535 {
			t.Skip("Data size not large enough for extended format testutils")
		}

		largeData := testutils.RepeatByte(dataSize, 0xCC)

		storeCmd, err := cmd.StoreColumnGraphicsInBufferLarge(
			NormalScale,
			DoubleScale,
			Color3,
			width,
			height,
			largeData,
		)
		if err != nil {
			t.Fatalf("StoreColumnGraphicsInBufferLarge failed: %v", err)
		}

		// Verify extended format
		if storeCmd[0] != common.GS || storeCmd[1] != '8' || storeCmd[2] != 'L' {
			t.Error("Large column data should use GS 8 L extended format")
		}
	})
}

func TestIntegration_Graphics_ErrorRecovery(t *testing.T) {
	cmd := NewGraphicsCommands()

	t.Run("invalid density settings recovery", func(t *testing.T) {
		// Try invalid density combination
		_, err := cmd.SetGraphicsDotDensity(
			FunctionCodeDensity1,
			Density180x180,
			Density360x360,
		)
		if err == nil {
			t.Error("Mismatched densities should return error")
		}

		// Should be able to set valid density after error
		densityCmd, err := cmd.SetGraphicsDotDensity(
			FunctionCodeDensity1,
			Density180x180,
			Density180x180,
		)
		if err != nil {
			t.Fatalf("Valid density setting failed after error: %v", err)
		}

		if len(densityCmd) != 9 {
			t.Error("Density command should be 9 bytes")
		}
	})

	t.Run("data size mismatch recovery", func(t *testing.T) {
		width := uint16(100)
		height := uint16(50)

		// Try with wrong data size
		wrongData := []byte{0xFF}
		_, err := cmd.StoreRasterGraphicsInBuffer(
			Monochrome,
			NormalScale,
			NormalScale,
			Color1,
			width,
			height,
			wrongData,
		)
		if err == nil {
			t.Error("Wrong data size should return error")
		}

		// Should work with correct data size
		widthBytes := (int(width) + 7) / 8
		correctData := testutils.RepeatByte(widthBytes*int(height), 0xAA)

		storeCmd, err := cmd.StoreRasterGraphicsInBuffer(
			Monochrome,
			NormalScale,
			NormalScale,
			Color1,
			width,
			height,
			correctData,
		)
		if err != nil {
			t.Fatalf("Valid data storage failed after error: %v", err)
		}

		if len(storeCmd) != 15+len(correctData) {
			t.Error("Store command has incorrect length")
		}
	})
}

func TestIntegration_Graphics_ColorRestrictions(t *testing.T) {
	cmd := NewGraphicsCommands()

	t.Run("column format color restrictions", func(t *testing.T) {
		width := uint16(100)
		height := uint16(64)
		heightBytes := (int(height) + 7) / 8
		dataSize := int(width) * heightBytes
		data := testutils.RepeatByte(dataSize, 0xFF)

		// Colors 1-3 should work for column format
		validColors := []GraphicsColor{
			Color1,
			Color2,
			Color3,
		}

		for _, color := range validColors {
			_, err := cmd.StoreColumnGraphicsInBuffer(
				NormalScale,
				NormalScale,
				color,
				width,
				height,
				data,
			)
			if err != nil {
				t.Errorf("Column format should accept color %v: %v", color, err)
			}
		}

		// Color 4 should not work for column format
		_, err := cmd.StoreColumnGraphicsInBuffer(
			NormalScale,
			NormalScale,
			Color4,
			width,
			height,
			data,
		)
		if err == nil {
			t.Error("Column format should reject Color4")
		}
	})

	t.Run("raster format accepts all colors", func(t *testing.T) {
		width := uint16(100)
		height := uint16(50)
		widthBytes := (int(width) + 7) / 8
		dataSize := widthBytes * int(height)
		data := testutils.RepeatByte(dataSize, 0xFF)

		// All colors should work for raster format
		allColors := []GraphicsColor{
			Color1,
			Color2,
			Color3,
			Color4,
		}

		for _, color := range allColors {
			_, err := cmd.StoreRasterGraphicsInBuffer(
				MultipleTone,
				NormalScale,
				NormalScale,
				color,
				width,
				height,
				data,
			)
			if err != nil {
				t.Errorf("Raster format should accept color %v: %v", color, err)
			}
		}
	})
}

func TestIntegration_Graphics_DimensionLimits(t *testing.T) {
	cmd := NewGraphicsCommands()

	t.Run("raster format dimension limits", func(t *testing.T) {
		// Test maximum width
		maxWidth := uint16(2400)
		height := uint16(10)
		widthBytes := (int(maxWidth) + 7) / 8
		data := testutils.RepeatByte(widthBytes*int(height), 0xFF)

		_, err := cmd.StoreRasterGraphicsInBuffer(
			Monochrome,
			NormalScale,
			NormalScale,
			Color1,
			maxWidth,
			height,
			data,
		)
		if err != nil {
			t.Errorf("Should accept max width %d: %v", maxWidth, err)
		}

		// Test width exceeding limit
		exceedWidth := uint16(2401)
		_, err = cmd.StoreRasterGraphicsInBuffer(
			Monochrome,
			NormalScale,
			NormalScale,
			Color1,
			exceedWidth,
			height,
			data,
		)
		if err == nil {
			t.Error("Should reject width exceeding 2400")
		}
	})

	t.Run("column format dimension limits", func(t *testing.T) {
		// Test maximum dimensions
		maxWidth := uint16(2048)
		maxHeight := uint16(128)
		heightBytes := (int(maxHeight) + 7) / 8
		data := testutils.RepeatByte(int(maxWidth)*heightBytes, 0xFF)

		_, err := cmd.StoreColumnGraphicsInBuffer(
			NormalScale,
			NormalScale,
			Color1,
			maxWidth,
			maxHeight,
			data,
		)
		if err != nil {
			t.Errorf("Should accept max dimensions %dx%d: %v", maxWidth, maxHeight, err)
		}

		// Test exceeding width
		_, err = cmd.StoreColumnGraphicsInBuffer(
			NormalScale,
			NormalScale,
			Color1,
			2049,
			maxHeight,
			data,
		)
		if err == nil {
			t.Error("Should reject width exceeding 2048")
		}

		// Test exceeding height
		_, err = cmd.StoreColumnGraphicsInBuffer(
			NormalScale,
			NormalScale,
			Color1,
			maxWidth,
			129,
			data,
		)
		if err == nil {
			t.Error("Should reject height exceeding 128")
		}
	})
}
