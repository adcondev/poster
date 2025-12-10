package test_test

import (
	"bytes"
	"testing"

	"github.com/adcondev/poster/internal/testutils"
	"github.com/adcondev/poster/pkg/commands/bitimage"
	"github.com/adcondev/poster/pkg/commands/common"
)

func TestIntegration_DownloadGraphics_CompleteWorkflow(t *testing.T) {
	cmd := bitimage.NewDownloadGraphicsCommands()

	t.Run("define and print download graphics", func(t *testing.T) {
		var buffer []byte

		// Step 1: Check remaining capacity
		capacityCmd, err := cmd.GetDownloadGraphicsRemainingCapacity(bitimage.DLFuncGetRemaining)
		if err != nil {
			t.Fatalf("GetDownloadGraphicsRemainingCapacity failed: %v", err)
		}
		buffer = append(buffer, capacityCmd...)

		// Step 2: Define monochrome graphics
		width := uint16(200)
		height := uint16(100)
		widthBytes := (int(width) + 7) / 8
		dataSize := widthBytes * int(height)

		colorData := []bitimage.DLGraphicsColorData{
			{
				Color: bitimage.Color1,
				Data:  testutils.RepeatByte(dataSize, 0xF0),
			},
		}

		defineCmd, err := cmd.DefineDownloadGraphics(
			bitimage.Monochrome,
			'L', 'G',
			width, height,
			colorData,
		)
		if err != nil {
			t.Fatalf("DefineDownloadGraphics failed: %v", err)
		}
		buffer = append(buffer, defineCmd...)

		// Step 3: Print the graphics
		printCmd, err := cmd.PrintDownloadGraphics('L', 'G', bitimage.NormalScale, bitimage.NormalScale)
		if err != nil {
			t.Fatalf("PrintDownloadGraphics failed: %v", err)
		}
		buffer = append(buffer, printCmd...)

		// Step 4: Get key code list
		listCmd := cmd.GetDownloadGraphicsKeyCodeList()
		buffer = append(buffer, listCmd...)

		// Verify commands
		if !bytes.Contains(buffer, []byte{common.GS, '(', 'L'}) {
			t.Error("Buffer should contain download graphics commands")
		}

		if len(buffer) < dataSize {
			t.Error("Buffer should contain graphics data")
		}
	})

	t.Run("multiple tone graphics with color groups", func(t *testing.T) {
		var buffer []byte

		width := uint16(100)
		height := uint16(50)
		widthBytes := (int(width) + 7) / 8
		dataSize := widthBytes * int(height)

		// Define multiple color groups
		colorData := []bitimage.DLGraphicsColorData{
			{Color: bitimage.Color1, Data: testutils.RepeatByte(dataSize, 0x11)},
			{Color: bitimage.Color2, Data: testutils.RepeatByte(dataSize, 0x22)},
			{Color: bitimage.Color3, Data: testutils.RepeatByte(dataSize, 0x33)},
			{Color: bitimage.Color4, Data: testutils.RepeatByte(dataSize, 0x44)},
		}

		defineCmd, err := cmd.DefineDownloadGraphics(
			bitimage.MultipleTone,
			'M', 'T',
			width, height,
			colorData,
		)
		if err != nil {
			t.Fatalf("DefineDownloadGraphics multi-tone failed: %v", err)
		}
		buffer = append(buffer, defineCmd...)

		// Print with double scale
		printCmd, err := cmd.PrintDownloadGraphics('M', 'T', bitimage.DoubleScale, bitimage.DoubleScale)
		if err != nil {
			t.Fatalf("PrintDownloadGraphics failed: %v", err)
		}
		buffer = append(buffer, printCmd...)

		expectedSize := 7 + 9 + 4*(1+dataSize) + 10 // capacity + define + print commands
		if len(buffer) < expectedSize {
			t.Errorf("Buffer size = %d, expected at least %d", len(buffer), expectedSize)
		}
	})

	t.Run("column format graphics", func(t *testing.T) {
		var buffer []byte

		width := uint16(256)
		height := uint16(64)
		heightBytes := (int(height) + 7) / 8
		dataSize := int(width) * heightBytes

		// Column format with colors 1 and 2
		colorData := []bitimage.DLGraphicsColorData{
			{Color: bitimage.Color1, Data: testutils.RepeatByte(dataSize, 0xAA)},
			{Color: bitimage.Color2, Data: testutils.RepeatByte(dataSize, 0x55)},
		}

		defineCmd, err := cmd.DefineDownloadGraphicsColumn(
			'C', 'L',
			width, height,
			colorData,
		)
		if err != nil {
			t.Fatalf("DefineDownloadGraphicsColumn failed: %v", err)
		}
		buffer = append(buffer, defineCmd...)

		// Column format with color 3 only
		colorData3 := []bitimage.DLGraphicsColorData{
			{Color: bitimage.Color3, Data: testutils.RepeatByte(dataSize, 0xFF)},
		}

		defineCmd3, err := cmd.DefineDownloadGraphicsColumn(
			'C', '3',
			width, height,
			colorData3,
		)
		if err != nil {
			t.Fatalf("DefineDownloadGraphicsColumn color 3 failed: %v", err)
		}
		buffer = append(buffer, defineCmd3...)

		if len(buffer) < 2*dataSize {
			t.Error("Buffer should contain column format data")
		}
	})

	t.Run("BMP graphics definition", func(t *testing.T) {
		var buffer []byte

		// Create minimal BMP header (54 bytes) + small image data
		bmpData := make([]byte, 100)
		bmpData[0] = 'B'
		bmpData[1] = 'M'
		// Fill rest with dummy data
		for i := 2; i < len(bmpData); i++ {
			bmpData[i] = byte(i)
		}

		defineCmd, err := cmd.DefineBMPDownloadGraphics(
			'B', 'M',
			bitimage.Monochrome,
			bmpData,
		)
		if err != nil {
			t.Fatalf("DefineBMPDownloadGraphics failed: %v", err)
		}
		buffer = append(buffer, defineCmd...)

		// Print the BMP graphics
		printCmd, err := cmd.PrintDownloadGraphics('B', 'M', bitimage.NormalScale, bitimage.NormalScale)
		if err != nil {
			t.Fatalf("PrintDownloadGraphics BMP failed: %v", err)
		}
		buffer = append(buffer, printCmd...)

		if !bytes.Contains(buffer, []byte{common.GS, 'D'}) {
			t.Error("Buffer should contain BMP definition command")
		}
	})

	t.Run("graphics deletion workflow", func(t *testing.T) {
		var buffer []byte

		// Delete specific graphics by key code
		deleteCmd, err := cmd.DeleteDownloadGraphicsByKeyCode('X', 'Y')
		if err != nil {
			t.Fatalf("DeleteDownloadGraphicsByKeyCode failed: %v", err)
		}
		buffer = append(buffer, deleteCmd...)

		// Delete all graphics
		deleteAllCmd := cmd.DeleteAllDownloadGraphics()
		buffer = append(buffer, deleteAllCmd...)

		if !bytes.Contains(buffer, []byte{'C', 'L', 'R'}) {
			t.Error("Buffer should contain delete all command")
		}

		expectedLen := 9 + 10 // delete by key (9 bytes) + delete all (10 bytes)
		if len(buffer) != expectedLen {
			t.Errorf("Buffer length = %d, expected %d", len(buffer), expectedLen)
		}
	})
}

func TestIntegration_DownloadGraphics_LargeDataHandling(t *testing.T) {
	cmd := bitimage.NewDownloadGraphicsCommands()

	t.Run("large raster format exceeding standard size", func(t *testing.T) {
		// Create data larger than 65535 bytes
		width := uint16(2000)
		height := uint16(1500)
		widthBytes := (int(width) + 7) / 8
		dataSize := widthBytes * int(height)

		colorData := []bitimage.DLGraphicsColorData{
			{
				Color: bitimage.Color1,
				Data:  testutils.RepeatByte(dataSize, 0xDD),
			},
		}

		defineCmd, err := cmd.DefineDownloadGraphicsLarge(
			bitimage.Monochrome,
			'L', 'D',
			width, height,
			colorData,
		)
		if err != nil {
			t.Fatalf("DefineDownloadGraphicsLarge failed: %v", err)
		}

		// Verify extended format (GS 8 L)
		if defineCmd[0] != common.GS || defineCmd[1] != '8' || defineCmd[2] != 'L' {
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
		width := uint16(2048)
		height := uint16(128)
		heightBytes := (int(height) + 7) / 8
		dataSize := int(width) * heightBytes

		// Large column data with two colors
		colorData := []bitimage.DLGraphicsColorData{
			{Color: bitimage.Color1, Data: testutils.RepeatByte(dataSize, 0x11)},
			{Color: bitimage.Color2, Data: testutils.RepeatByte(dataSize, 0x22)},
		}

		defineCmd, err := cmd.DefineDownloadGraphicsColumnLarge(
			'C', 'X',
			width, height,
			colorData,
		)
		if err != nil {
			t.Fatalf("DefineDownloadGraphicsColumnLarge failed: %v", err)
		}

		// Verify extended format
		if defineCmd[0] != common.GS || defineCmd[1] != '8' || defineCmd[2] != 'L' {
			t.Error("Large column data should use extended command format")
		}

		if len(defineCmd) < 10+2*(1+dataSize) {
			t.Error("Command should contain all color data")
		}
	})
}

func TestIntegration_DownloadGraphics_ErrorHandling(t *testing.T) {
	cmd := bitimage.NewDownloadGraphicsCommands()

	t.Run("invalid color combinations", func(t *testing.T) {
		width := uint16(100)
		height := uint16(50)
		widthBytes := (int(width) + 7) / 8
		dataSize := widthBytes * int(height)

		// Monochrome with multiple colors - should fail
		colorData := []bitimage.DLGraphicsColorData{
			{Color: bitimage.Color1, Data: testutils.RepeatByte(dataSize, 0x11)},
			{Color: bitimage.Color2, Data: testutils.RepeatByte(dataSize, 0x22)},
		}

		_, err := cmd.DefineDownloadGraphics(
			bitimage.Monochrome,
			'E', 'R',
			width, height,
			colorData,
		)
		if err == nil {
			t.Error("Monochrome with multiple colors should return error")
		}

		// Column format with invalid color combination
		heightBytes := (int(height) + 7) / 8
		columnDataSize := int(width) * heightBytes

		invalidColorData := []bitimage.DLGraphicsColorData{
			{Color: bitimage.Color3, Data: testutils.RepeatByte(columnDataSize, 0x33)},
			{Color: bitimage.Color1, Data: testutils.RepeatByte(columnDataSize, 0x11)},
		}

		_, err = cmd.DefineDownloadGraphicsColumn(
			'I', 'C',
			width, height,
			invalidColorData,
		)
		if err == nil {
			t.Error("Color 3 with other colors in column format should return error")
		}
	})

	t.Run("invalid key codes", func(t *testing.T) {
		// Key code out of range
		_, err := cmd.DeleteDownloadGraphicsByKeyCode(31, 'A')
		if err == nil {
			t.Error("Key code below 32 should return error")
		}

		_, err = cmd.DeleteDownloadGraphicsByKeyCode('A', 127)
		if err == nil {
			t.Error("Key code above 126 should return error")
		}

		// Print with invalid key codes
		_, err = cmd.PrintDownloadGraphics(200, 'X', bitimage.NormalScale, bitimage.NormalScale)
		if err == nil {
			t.Error("Invalid key code should return error")
		}
	})

	t.Run("invalid dimensions", func(t *testing.T) {
		colorData := []bitimage.DLGraphicsColorData{
			{Color: bitimage.Color1, Data: []byte{0xFF}},
		}

		// Width exceeds limit
		_, err := cmd.DefineDownloadGraphics(
			bitimage.Monochrome,
			'W', 'E',
			8193, 100,
			colorData,
		)
		if err == nil {
			t.Error("Width exceeding 8192 should return error")
		}

		// Height exceeds limit
		_, err = cmd.DefineDownloadGraphics(
			bitimage.Monochrome,
			'H', 'E',
			100, 2305,
			colorData,
		)
		if err == nil {
			t.Error("Height exceeding 2304 should return error")
		}

		// Zero dimensions
		_, err = cmd.DefineDownloadGraphics(
			bitimage.Monochrome,
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

		colorData := []bitimage.DLGraphicsColorData{
			{Color: bitimage.Color1, Data: wrongData},
		}

		_, err := cmd.DefineDownloadGraphics(
			bitimage.Monochrome,
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

		_, err := cmd.DefineBMPDownloadGraphics('B', 'X', bitimage.Monochrome, invalidBMP)
		if err == nil {
			t.Error("Invalid BMP data should return error")
		}

		// BMP with wrong signature
		wrongBMP := make([]byte, 54)
		wrongBMP[0] = 'X'
		wrongBMP[1] = 'Y'

		_, err = cmd.DefineBMPDownloadGraphics('B', 'Y', bitimage.Monochrome, wrongBMP)
		if err == nil {
			t.Error("BMP with wrong signature should return error")
		}
	})
}

func TestIntegration_DownloadGraphics_ScalingModes(t *testing.T) {
	cmd := bitimage.NewDownloadGraphicsCommands()

	// Define testutils graphics once
	width := uint16(100)
	height := uint16(50)
	widthBytes := (int(width) + 7) / 8
	dataSize := widthBytes * int(height)

	colorData := []bitimage.DLGraphicsColorData{
		{Color: bitimage.Color1, Data: testutils.RepeatByte(dataSize, 0xCC)},
	}

	defineCmd, err := cmd.DefineDownloadGraphics(
		bitimage.Monochrome,
		'S', 'C',
		width, height,
		colorData,
	)
	if err != nil {
		t.Fatalf("DefineDownloadGraphics failed: %v", err)
	}

	scales := []struct {
		name       string
		horizontal bitimage.GraphicsScale
		vertical   bitimage.GraphicsScale
		expectedX  byte
		expectedY  byte
	}{
		{"normal", bitimage.NormalScale, bitimage.NormalScale, 1, 1},
		{"double width", bitimage.DoubleScale, bitimage.NormalScale, 2, 1},
		{"double height", bitimage.NormalScale, bitimage.DoubleScale, 1, 2},
		{"quadruple", bitimage.DoubleScale, bitimage.DoubleScale, 2, 2},
	}

	for _, scale := range scales {
		t.Run(scale.name, func(t *testing.T) {
			printCmd, err := cmd.PrintDownloadGraphics(
				'S', 'C',
				scale.horizontal,
				scale.vertical,
			)
			if err != nil {
				t.Fatalf("PrintDownloadGraphics with %s failed: %v", scale.name, err)
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
		printCmd, _ := cmd.PrintDownloadGraphics('S', 'C', scale.horizontal, scale.vertical)
		fullBuffer = append(fullBuffer, printCmd...)
	}

	expectedLen := len(defineCmd) + 4*11 // define + 4 print commands
	if len(fullBuffer) != expectedLen {
		t.Errorf("Full buffer length = %d, expected %d", len(fullBuffer), expectedLen)
	}
}

func TestIntegration_DownloadGraphics_DuplicateColors(t *testing.T) {
	cmd := bitimage.NewDownloadGraphicsCommands()

	width := uint16(50)
	height := uint16(50)
	widthBytes := (int(width) + 7) / 8
	dataSize := widthBytes * int(height)

	// Attempt to define duplicate colors
	colorData := []bitimage.DLGraphicsColorData{
		{Color: bitimage.Color1, Data: testutils.RepeatByte(dataSize, 0x11)},
		{Color: bitimage.Color1, Data: testutils.RepeatByte(dataSize, 0x22)}, // Duplicate
	}

	_, err := cmd.DefineDownloadGraphics(
		bitimage.MultipleTone,
		'D', 'P',
		width, height,
		colorData,
	)

	if err == nil {
		t.Error("Duplicate colors should return error")
	}
}
