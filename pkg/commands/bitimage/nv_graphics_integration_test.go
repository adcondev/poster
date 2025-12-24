package bitimage

import (
	"bytes"
	"testing"

	"github.com/adcondev/poster/internal/testutils"
	"github.com/adcondev/poster/pkg/commands/shared"
)

func TestIntegration_NVGraphics_CompleteWorkflow(t *testing.T) {
	cmd := NewNVGraphicsCommands()

	t.Run("define and print NV graphics", func(t *testing.T) {
		var buffer []byte

		// Step 1: Check NV capacity
		capacityCmd, err := cmd.GetNVGraphicsCapacity(NVFuncGetCapacity)
		if err != nil {
			t.Fatalf("GetNVGraphicsCapacity failed: %v", err)
		}
		buffer = append(buffer, capacityCmd...)

		// Step 2: Check remaining capacity
		remainingCmd, err := cmd.GetNVGraphicsRemainingCapacity(NVFuncGetRemaining)
		if err != nil {
			t.Fatalf("GetNVGraphicsRemainingCapacity failed: %v", err)
		}
		buffer = append(buffer, remainingCmd...)

		// Step 3: Define monochrome NV graphics
		width := uint16(300)
		height := uint16(150)
		widthBytes := (int(width) + 7) / 8
		dataSize := widthBytes * int(height)

		colorData := []NVGraphicsColorData{
			{
				Color: Color1,
				Data:  testutils.RepeatByte(dataSize, 0xF0),
			},
		}

		defineCmd, err := cmd.DefineNVRasterGraphics(
			Monochrome,
			'N', 'V',
			width, height,
			colorData,
		)
		if err != nil {
			t.Fatalf("DefineNVRasterGraphics failed: %v", err)
		}
		buffer = append(buffer, defineCmd...)

		// Step 4: Print the NV graphics
		printCmd, err := cmd.PrintNVGraphics('N', 'V', NormalScale, NormalScale)
		if err != nil {
			t.Fatalf("PrintNVGraphics failed: %v", err)
		}
		buffer = append(buffer, printCmd...)

		// Step 5: Get key code list
		listCmd := cmd.GetNVGraphicsKeyCodeList()
		buffer = append(buffer, listCmd...)

		// Verify commands
		if !bytes.Contains(buffer, []byte{shared.GS, '(', 'L'}) {
			t.Error("Buffer should contain NV graphics commands")
		}

		expectedMinSize := 7 + 7 + dataSize + 11 + 9 // capacity + remaining + define + print + list
		if len(buffer) < expectedMinSize {
			t.Errorf("Buffer size = %d, expected at least %d", len(buffer), expectedMinSize)
		}
	})

	t.Run("multiple tone NV graphics with color groups", func(t *testing.T) {
		var buffer []byte

		width := uint16(200)
		height := uint16(100)
		widthBytes := (int(width) + 7) / 8
		dataSize := widthBytes * int(height)

		// Define multiple color groups
		colorData := []NVGraphicsColorData{
			{Color: Color1, Data: testutils.RepeatByte(dataSize, 0x11)},
			{Color: Color2, Data: testutils.RepeatByte(dataSize, 0x22)},
			{Color: Color3, Data: testutils.RepeatByte(dataSize, 0x33)},
			{Color: Color4, Data: testutils.RepeatByte(dataSize, 0x44)},
		}

		defineCmd, err := cmd.DefineNVRasterGraphics(
			MultipleTone,
			'M', 'C',
			width, height,
			colorData,
		)
		if err != nil {
			t.Fatalf("DefineNVRasterGraphics multi-tone failed: %v", err)
		}
		buffer = append(buffer, defineCmd...)

		// Print with double scale
		printCmd, err := cmd.PrintNVGraphics('M', 'C', DoubleScale, DoubleScale)
		if err != nil {
			t.Fatalf("PrintNVGraphics failed: %v", err)
		}
		buffer = append(buffer, printCmd...)

		expectedSize := 9 + 4*(1+dataSize) + 11 // define + print
		if len(buffer) < expectedSize {
			t.Errorf("Buffer size = %d, expected at least %d", len(buffer), expectedSize)
		}
	})

	t.Run("column format NV graphics", func(t *testing.T) {
		var buffer []byte

		width := uint16(512)
		height := uint16(128)
		heightBytes := (int(height) + 7) / 8
		dataSize := int(width) * heightBytes

		// Column format with colors 1 and 2
		colorData := []NVGraphicsColorData{
			{Color: Color1, Data: testutils.RepeatByte(dataSize, 0xAA)},
			{Color: Color2, Data: testutils.RepeatByte(dataSize, 0x55)},
		}

		defineCmd, err := cmd.DefineNVColumnGraphics(
			'C', 'F',
			width, height,
			colorData,
		)
		if err != nil {
			t.Fatalf("DefineNVColumnGraphics failed: %v", err)
		}
		buffer = append(buffer, defineCmd...)

		// Column format with color 3 only
		colorData3 := []NVGraphicsColorData{
			{Color: Color3, Data: testutils.RepeatByte(dataSize, 0xFF)},
		}

		defineCmd3, err := cmd.DefineNVColumnGraphics(
			'C', '3',
			width, height,
			colorData3,
		)
		if err != nil {
			t.Fatalf("DefineNVColumnGraphics color 3 failed: %v", err)
		}
		buffer = append(buffer, defineCmd3...)

		// Print both graphics
		printCmd1, _ := cmd.PrintNVGraphics('C', 'F', NormalScale, DoubleScale)
		printCmd2, _ := cmd.PrintNVGraphics('C', '3', DoubleScale, NormalScale)

		buffer = append(buffer, printCmd1...)
		buffer = append(buffer, printCmd2...)

		if len(buffer) < 3*dataSize {
			t.Error("Buffer should contain column format data")
		}
	})

	t.Run("Windows BMP NV graphics", func(t *testing.T) {
		var buffer []byte

		// Create valid BMP header + small image data
		bmpData := make([]byte, 200)
		bmpData[0] = 'B'
		bmpData[1] = 'M'
		// Fill rest with pattern
		for i := 2; i < len(bmpData); i++ {
			bmpData[i] = byte(i % 256)
		}

		// Define monochrome BMP
		defineCmd, err := cmd.DefineWindowsBMPNVGraphics(
			'B', '1',
			Monochrome,
			bmpData,
		)
		if err != nil {
			t.Fatalf("DefineWindowsBMPNVGraphics monochrome failed: %v", err)
		}
		buffer = append(buffer, defineCmd...)

		// Define multiple tone BMP
		defineCmd2, err := cmd.DefineWindowsBMPNVGraphics(
			'B', '2',
			MultipleTone,
			bmpData,
		)
		if err != nil {
			t.Fatalf("DefineWindowsBMPNVGraphics multiple tone failed: %v", err)
		}
		buffer = append(buffer, defineCmd2...)

		// Print both BMPs
		printCmd1, _ := cmd.PrintNVGraphics('B', '1', NormalScale, NormalScale)
		printCmd2, _ := cmd.PrintNVGraphics('B', '2', DoubleScale, DoubleScale)

		buffer = append(buffer, printCmd1...)
		buffer = append(buffer, printCmd2...)

		if !bytes.Contains(buffer, []byte{shared.GS, 'D'}) {
			t.Error("Buffer should contain BMP definition commands")
		}
	})

	t.Run("NV graphics management workflow", func(t *testing.T) {
		var buffer []byte

		// Get key code list
		listCmd := cmd.GetNVGraphicsKeyCodeList()
		buffer = append(buffer, listCmd...)

		// Delete specific NV graphics
		deleteCmd, err := cmd.DeleteNVGraphicsByKeyCode('X', 'Y')
		if err != nil {
			t.Fatalf("DeleteNVGraphicsByKeyCode failed: %v", err)
		}
		buffer = append(buffer, deleteCmd...)

		// Delete all NV graphics
		deleteAllCmd := cmd.DeleteAllNVGraphics()
		buffer = append(buffer, deleteAllCmd...)

		if !bytes.Contains(buffer, []byte{'C', 'L', 'R'}) {
			t.Error("Buffer should contain delete all command")
		}

		expectedLen := 9 + 9 + 10 // list (9 bytes) + delete (9 bytes) + delete all (10 bytes)
		if len(buffer) != expectedLen {
			t.Errorf("Buffer length = %d, expected %d", len(buffer), expectedLen)
		}
	})
}

func TestIntegration_NVGraphics_LargeDataHandling(t *testing.T) {
	cmd := NewNVGraphicsCommands()

	t.Run("large raster format exceeding standard size", func(t *testing.T) {
		// Create data larger than 65535 bytes
		width := uint16(2000)
		height := uint16(1500)
		widthBytes := (int(width) + 7) / 8
		dataSize := widthBytes * int(height)

		if dataSize <= 65535 {
			t.Skip("Data size not large enough for extended format testutils")
		}

		colorData := []NVGraphicsColorData{
			{
				Color: Color1,
				Data:  testutils.RepeatByte(dataSize, 0xEE),
			},
		}

		defineCmd, err := cmd.DefineNVRasterGraphicsLarge(
			Monochrome,
			'L', 'R',
			width, height,
			colorData,
		)
		if err != nil {
			t.Fatalf("DefineNVRasterGraphicsLarge failed: %v", err)
		}

		// Verify extended format (GS 8 L)
		if defineCmd[0] != shared.GS || defineCmd[1] != '8' || defineCmd[2] != 'L' {
			t.Error("Large data should use extended command format")
		}

		// Verify 32-bit size encoding
		totalSize := uint32(9 + 1 + dataSize) //nolint:gosec
		p1 := defineCmd[3]
		p2 := defineCmd[4]
		p3 := defineCmd[5]
		p4 := defineCmd[6]

		calculatedSize := uint32(p1) + uint32(p2)<<8 + uint32(p3)<<16 + uint32(p4)<<24
		if calculatedSize != totalSize {
			t.Errorf("Size encoding incorrect: got %d, want %d", calculatedSize, totalSize)
		}
	})

	t.Run("large column format", func(t *testing.T) {
		width := uint16(4096)
		height := uint16(256)
		heightBytes := (int(height) + 7) / 8
		dataSize := int(width) * heightBytes

		// Large column data with two colors
		colorData := []NVGraphicsColorData{
			{Color: Color1, Data: testutils.RepeatByte(dataSize, 0x11)},
			{Color: Color2, Data: testutils.RepeatByte(dataSize, 0x22)},
		}

		defineCmd, err := cmd.DefineNVColumnGraphicsLarge(
			'L', 'C',
			width, height,
			colorData,
		)
		if err != nil {
			t.Fatalf("DefineNVColumnGraphicsLarge failed: %v", err)
		}

		// Verify extended format
		if defineCmd[0] != shared.GS || defineCmd[1] != '8' || defineCmd[2] != 'L' {
			t.Error("Large column data should use extended command format")
		}

		if len(defineCmd) < 10+2*(1+dataSize) {
			t.Error("Command should contain all color data")
		}
	})
}

func TestIntegration_NVGraphics_ColorRestrictions(t *testing.T) {
	cmd := NewNVGraphicsCommands()

	t.Run("monochrome tone color restrictions", func(t *testing.T) {
		width := uint16(100)
		height := uint16(50)
		widthBytes := (int(width) + 7) / 8
		dataSize := widthBytes * int(height)

		// Monochrome should accept only Color1 or Color2
		validColors := []GraphicsColor{
			Color1,
			Color2,
		}

		for _, color := range validColors {
			colorData := []NVGraphicsColorData{
				{Color: color, Data: testutils.RepeatByte(dataSize, 0xFF)},
			}

			_, err := cmd.DefineNVRasterGraphics(
				Monochrome,
				'T', byte('0'+color),
				width, height,
				colorData,
			)
			if err != nil {
				t.Errorf("Monochrome should accept color %v: %v", color, err)
			}
		}

		// Monochrome should reject multiple colors
		multiColorData := []NVGraphicsColorData{
			{Color: Color1, Data: testutils.RepeatByte(dataSize, 0x11)},
			{Color: Color2, Data: testutils.RepeatByte(dataSize, 0x22)},
		}

		_, err := cmd.DefineNVRasterGraphics(
			Monochrome,
			'E', 'M',
			width, height,
			multiColorData,
		)
		if err == nil {
			t.Error("Monochrome should reject multiple colors")
		}
	})

	t.Run("column format color restrictions", func(t *testing.T) {
		width := uint16(100)
		height := uint16(64)
		heightBytes := (int(height) + 7) / 8
		dataSize := int(width) * heightBytes

		// Test valid combinations for column format
		validCombinations := [][]GraphicsColor{
			{Color1},
			{Color2},
			{Color3},
			{Color1, Color2},
		}

		for i, colors := range validCombinations {
			var colorData []NVGraphicsColorData
			for _, color := range colors {
				colorData = append(colorData, NVGraphicsColorData{
					Color: color,
					Data:  testutils.RepeatByte(dataSize, byte(color)),
				})
			}

			_, err := cmd.DefineNVColumnGraphics(
				'V', byte('0'+i),
				width, height,
				colorData,
			)
			if err != nil {
				t.Errorf("Column format should accept combination %v: %v", colors, err)
			}
		}

		// Test invalid combination (Color3 with others)
		invalidColorData := []NVGraphicsColorData{
			{Color: Color3, Data: testutils.RepeatByte(dataSize, 0x33)},
			{Color: Color1, Data: testutils.RepeatByte(dataSize, 0x11)},
		}

		_, err := cmd.DefineNVColumnGraphics(
			'I', 'V',
			width, height,
			invalidColorData,
		)
		if err == nil {
			t.Error("Column format should reject Color3 with other colors")
		}
	})

	t.Run("duplicate colors detection", func(t *testing.T) {
		width := uint16(50)
		height := uint16(50)
		widthBytes := (int(width) + 7) / 8
		dataSize := widthBytes * int(height)

		// Try to define duplicate colors
		duplicateColorData := []NVGraphicsColorData{
			{Color: Color1, Data: testutils.RepeatByte(dataSize, 0x11)},
			{Color: Color1, Data: testutils.RepeatByte(dataSize, 0x22)}, // Duplicate
		}

		_, err := cmd.DefineNVRasterGraphics(
			MultipleTone,
			'D', 'P',
			width, height,
			duplicateColorData,
		)
		if err == nil {
			t.Error("Should reject duplicate colors")
		}
	})
}

func TestIntegration_NVGraphics_ErrorHandling(t *testing.T) {
	cmd := NewNVGraphicsCommands()

	t.Run("invalid key codes", func(t *testing.T) {
		// Key code out of range
		_, err := cmd.DeleteNVGraphicsByKeyCode(31, 'A')
		if err == nil {
			t.Error("Key code below 32 should return error")
		}

		_, err = cmd.DeleteNVGraphicsByKeyCode('A', 127)
		if err == nil {
			t.Error("Key code above 126 should return error")
		}

		// Print with invalid key codes
		_, err = cmd.PrintNVGraphics(200, 'X', NormalScale, NormalScale)
		if err == nil {
			t.Error("Invalid key code should return error")
		}
	})

	t.Run("invalid dimensions", func(t *testing.T) {
		colorData := []NVGraphicsColorData{
			{Color: Color1, Data: []byte{0xFF}},
		}

		// Width exceeds limit
		_, err := cmd.DefineNVRasterGraphics(
			Monochrome,
			'W', 'E',
			8193, 100,
			colorData,
		)
		if err == nil {
			t.Error("Width exceeding 8192 should return error")
		}

		// Height exceeds limit
		_, err = cmd.DefineNVRasterGraphics(
			Monochrome,
			'H', 'E',
			100, 2305,
			colorData,
		)
		if err == nil {
			t.Error("Height exceeding 2304 should return error")
		}

		// Zero dimensions
		_, err = cmd.DefineNVRasterGraphics(
			Monochrome,
			'Z', 'D',
			0, 0,
			colorData,
		)
		if err == nil {
			t.Error("Zero dimensions should return error")
		}
	})

	t.Run("data size mismatch", func(t *testing.T) {
		width := uint16(100)
		height := uint16(50)
		wrongData := []byte{0xFF, 0xFF} // Too small

		colorData := []NVGraphicsColorData{
			{Color: Color1, Data: wrongData},
		}

		_, err := cmd.DefineNVRasterGraphics(
			Monochrome,
			'D', 'M',
			width, height,
			colorData,
		)
		if err == nil {
			t.Error("Incorrect data size should return error")
		}
	})

	t.Run("invalid BMP data", func(t *testing.T) {
		// BMP without proper header
		invalidBMP := []byte{0xFF, 0xFF, 0xFF}

		_, err := cmd.DefineWindowsBMPNVGraphics('B', 'X', Monochrome, invalidBMP)
		if err == nil {
			t.Error("Invalid BMP data should return error")
		}

		// BMP with wrong signature
		wrongBMP := make([]byte, 54)
		wrongBMP[0] = 'X'
		wrongBMP[1] = 'Y'

		_, err = cmd.DefineWindowsBMPNVGraphics('B', 'Y', Monochrome, wrongBMP)
		if err == nil {
			t.Error("BMP with wrong signature should return error")
		}

		// BMP too small
		smallBMP := []byte{'B', 'M'}

		_, err = cmd.DefineWindowsBMPNVGraphics('B', 'S', Monochrome, smallBMP)
		if err == nil {
			t.Error("BMP too small should return error")
		}
	})
}

func TestIntegration_NVGraphics_ScalingModes(t *testing.T) {
	cmd := NewNVGraphicsCommands()

	// Define testutils graphics once
	width := uint16(100)
	height := uint16(50)
	widthBytes := (int(width) + 7) / 8
	dataSize := widthBytes * int(height)

	colorData := []NVGraphicsColorData{
		{Color: Color1, Data: testutils.RepeatByte(dataSize, 0xCC)},
	}

	defineCmd, err := cmd.DefineNVRasterGraphics(
		Monochrome,
		'S', 'T',
		width, height,
		colorData,
	)
	if err != nil {
		t.Fatalf("DefineNVRasterGraphics failed: %v", err)
	}

	scales := []struct {
		name       string
		horizontal GraphicsScale
		vertical   GraphicsScale
		expectedX  byte
		expectedY  byte
	}{
		{"normal", NormalScale, NormalScale, 1, 1},
		{"double width", DoubleScale, NormalScale, 2, 1},
		{"double height", NormalScale, DoubleScale, 1, 2},
		{"quadruple", DoubleScale, DoubleScale, 2, 2},
	}

	for _, scale := range scales {
		t.Run(scale.name, func(t *testing.T) {
			printCmd, err := cmd.PrintNVGraphics(
				'S', 'T',
				scale.horizontal,
				scale.vertical,
			)
			if err != nil {
				t.Fatalf("PrintNVGraphics with %s failed: %v", scale.name, err)
			}

			// Verify scale parameters in command
			if printCmd[9] != scale.expectedX {
				t.Errorf("Horizontal scale = %d, want %d", printCmd[9], scale.expectedX)
			}
			if printCmd[10] != scale.expectedY {
				t.Errorf("Vertical scale = %d, want %d", printCmd[10], scale.expectedY)
			}
		})
	}

	// Combine define and print commands
	var fullBuffer []byte
	fullBuffer = append(fullBuffer, defineCmd...)

	for _, scale := range scales {
		printCmd, _ := cmd.PrintNVGraphics('S', 'T', scale.horizontal, scale.vertical)
		fullBuffer = append(fullBuffer, printCmd...)
	}

	expectedLen := len(defineCmd) + 4*11 // define + 4 print commands
	if len(fullBuffer) != expectedLen {
		t.Errorf("Full buffer length = %d, expected %d", len(fullBuffer), expectedLen)
	}
}
