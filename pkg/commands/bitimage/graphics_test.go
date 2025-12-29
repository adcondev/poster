package bitimage_test

import (
	"testing"

	"github.com/adcondev/poster/internal/testutils"
	"github.com/adcondev/poster/pkg/commands/bitimage"
	"github.com/adcondev/poster/pkg/commands/shared"
)

// ============================================================================
// Graphics Commands Tests
// ============================================================================

func TestGraphicsCommands_SetGraphicsDotDensity(t *testing.T) {
	cmd := bitimage.NewGraphicsCommands()

	tests := []struct {
		name    string
		fn      bitimage.FunctionCode
		x       bitimage.DotDensity
		y       bitimage.DotDensity
		want    []byte
		wantErr error
	}{
		{
			name:    "180x180 dpi with function code 1",
			fn:      bitimage.FunctionCodeDensity1,
			x:       bitimage.Density180x180,
			y:       bitimage.Density180x180,
			want:    []byte{shared.GS, '(', 'L', 0x04, 0x00, 0x30, 1, 50, 50},
			wantErr: nil,
		},
		{
			name:    "180x180 dpi with function code 49",
			fn:      bitimage.FunctionCodeDensity49,
			x:       bitimage.Density180x180,
			y:       bitimage.Density180x180,
			want:    []byte{shared.GS, '(', 'L', 0x04, 0x00, 0x30, 49, 50, 50},
			wantErr: nil,
		},
		{
			name:    "360x360 dpi with function code 1",
			fn:      bitimage.FunctionCodeDensity1,
			x:       bitimage.Density360x360,
			y:       bitimage.Density360x360,
			want:    []byte{shared.GS, '(', 'L', 0x04, 0x00, 0x30, 1, 51, 51},
			wantErr: nil,
		},
		{
			name:    "360x360 dpi with function code 49",
			fn:      bitimage.FunctionCodeDensity49,
			x:       bitimage.Density360x360,
			y:       bitimage.Density360x360,
			want:    []byte{shared.GS, '(', 'L', 0x04, 0x00, 0x30, 49, 51, 51},
			wantErr: nil,
		},
		{
			name:    "invalid function code 2",
			fn:      bitimage.FunctionCodePrint2,
			x:       bitimage.Density180x180,
			y:       bitimage.Density180x180,
			want:    nil,
			wantErr: bitimage.ErrInvalidFunctionCode,
		},
		{
			name:    "invalid function code 50",
			fn:      bitimage.FunctionCodePrint50,
			x:       bitimage.Density180x180,
			y:       bitimage.Density180x180,
			want:    nil,
			wantErr: bitimage.ErrInvalidFunctionCode,
		},
		{
			name:    "invalid density value",
			fn:      bitimage.FunctionCodeDensity1,
			x:       52,
			y:       52,
			want:    nil,
			wantErr: bitimage.ErrInvalidDensityValue,
		},
		{
			name:    "mismatched density values",
			fn:      bitimage.FunctionCodeDensity1,
			x:       bitimage.Density180x180,
			y:       bitimage.Density360x360,
			want:    nil,
			wantErr: bitimage.ErrInvalidDensityCombination,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cmd.SetGraphicsDotDensity(tt.fn, tt.x, tt.y)

			if !testutils.AssertErrorOccurred(t, err, tt.wantErr != nil, "SetGraphicsDotDensity") {
				return
			}
			if tt.wantErr != nil {
				testutils.AssertError(t, err, tt.wantErr)
				return
			}

			testutils.AssertBytes(t, got, tt.want, "SetGraphicsDotDensity(%v, %v, %v)", tt.fn, tt.x, tt.y)
		})
	}
}

func TestGraphicsCommands_PrintBufferedGraphics(t *testing.T) {
	cmd := bitimage.NewGraphicsCommands()

	tests := []struct {
		name    string
		fn      bitimage.FunctionCode
		want    []byte
		wantErr error
	}{
		{
			name:    "function code 2",
			fn:      bitimage.FunctionCodePrint2,
			want:    []byte{shared.GS, '(', 'L', 0x02, 0x00, 0x30, 2},
			wantErr: nil,
		},
		{
			name:    "function code 50",
			fn:      bitimage.FunctionCodePrint50,
			want:    []byte{shared.GS, '(', 'L', 0x02, 0x00, 0x30, 50},
			wantErr: nil,
		},
		{
			name:    "invalid function code 1",
			fn:      bitimage.FunctionCodeDensity1,
			want:    nil,
			wantErr: bitimage.ErrInvalidFunctionCode,
		},
		{
			name:    "invalid function code 49",
			fn:      bitimage.FunctionCodeDensity49,
			want:    nil,
			wantErr: bitimage.ErrInvalidFunctionCode,
		},
		{
			name:    "invalid function code 99",
			fn:      99,
			want:    nil,
			wantErr: bitimage.ErrInvalidFunctionCode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cmd.PrintBufferedGraphics(tt.fn)

			if !testutils.AssertErrorOccurred(t, err, tt.wantErr != nil, "PrintBufferedGraphics") {
				return
			}
			if tt.wantErr != nil {
				testutils.AssertError(t, err, tt.wantErr)
				return
			}

			testutils.AssertBytes(t, got, tt.want, "PrintBufferedGraphics(%v)", tt.fn)
		})
	}
}

func TestGraphicsCommands_StoreRasterGraphicsInBuffer(t *testing.T) {
	cmd := bitimage.NewGraphicsCommands()

	// Helper to create testutils data for raster format
	createRasterData := func(width, height uint16) []byte {
		widthBytes := (int(width) + 7) / 8
		return testutils.RepeatByte(widthBytes*int(height), 0xFF)
	}

	tests := []struct {
		name            string
		tone            bitimage.GraphicsTone
		horizontalScale bitimage.GraphicsScale
		verticalScale   bitimage.GraphicsScale
		color           bitimage.GraphicsColor
		width           uint16
		height          uint16
		data            []byte
		wantErr         error
	}{
		{
			name:            "monochrome normal scale",
			tone:            bitimage.Monochrome,
			horizontalScale: bitimage.NormalScale,
			verticalScale:   bitimage.NormalScale,
			color:           bitimage.Color1,
			width:           100,
			height:          50,
			data:            createRasterData(100, 50),
			wantErr:         nil,
		},
		{
			name:            "monochrome double width",
			tone:            bitimage.Monochrome,
			horizontalScale: bitimage.DoubleScale,
			verticalScale:   bitimage.NormalScale,
			color:           bitimage.Color2,
			width:           200,
			height:          100,
			data:            createRasterData(200, 100),
			wantErr:         nil,
		},
		{
			name:            "monochrome double height",
			tone:            bitimage.Monochrome,
			horizontalScale: bitimage.NormalScale,
			verticalScale:   bitimage.DoubleScale,
			color:           bitimage.Color3,
			width:           150,
			height:          1200,
			data:            createRasterData(150, 1200),
			wantErr:         nil,
		},
		{
			name:            "multiple tone normal scale",
			tone:            bitimage.MultipleTone,
			horizontalScale: bitimage.NormalScale,
			verticalScale:   bitimage.NormalScale,
			color:           bitimage.Color4,
			width:           100,
			height:          600,
			data:            createRasterData(100, 600),
			wantErr:         nil,
		},
		{
			name:            "multiple tone double scale",
			tone:            bitimage.MultipleTone,
			horizontalScale: bitimage.DoubleScale,
			verticalScale:   bitimage.DoubleScale,
			color:           bitimage.Color1,
			width:           100,
			height:          300,
			data:            createRasterData(100, 300),
			wantErr:         nil,
		},
		{
			name:            "maximum width",
			tone:            bitimage.Monochrome,
			horizontalScale: bitimage.NormalScale,
			verticalScale:   bitimage.NormalScale,
			color:           bitimage.Color1,
			width:           2400,
			height:          10,
			data:            createRasterData(2400, 10),
			wantErr:         nil,
		},
		{
			name:            "invalid tone",
			tone:            49,
			horizontalScale: bitimage.NormalScale,
			verticalScale:   bitimage.NormalScale,
			color:           bitimage.Color1,
			width:           100,
			height:          100,
			data:            createRasterData(100, 100),
			wantErr:         bitimage.ErrInvalidTone,
		},
		{
			name:            "invalid horizontal scale",
			tone:            bitimage.Monochrome,
			horizontalScale: 0,
			verticalScale:   bitimage.NormalScale,
			color:           bitimage.Color1,
			width:           100,
			height:          100,
			data:            createRasterData(100, 100),
			wantErr:         bitimage.ErrInvalidScale,
		},
		{
			name:            "invalid vertical scale",
			tone:            bitimage.Monochrome,
			horizontalScale: bitimage.NormalScale,
			verticalScale:   3,
			color:           bitimage.Color1,
			width:           100,
			height:          100,
			data:            createRasterData(100, 100),
			wantErr:         bitimage.ErrInvalidScale,
		},
		{
			name:            "invalid color",
			tone:            bitimage.Monochrome,
			horizontalScale: bitimage.NormalScale,
			verticalScale:   bitimage.NormalScale,
			color:           48,
			width:           100,
			height:          100,
			data:            createRasterData(100, 100),
			wantErr:         bitimage.ErrInvalidColor,
		},
		{
			name:            "width exceeds limit",
			tone:            bitimage.Monochrome,
			horizontalScale: bitimage.NormalScale,
			verticalScale:   bitimage.NormalScale,
			color:           bitimage.Color1,
			width:           2401,
			height:          100,
			data:            createRasterData(2401, 100),
			wantErr:         bitimage.ErrInvalidWidth,
		},
		{
			name:            "height exceeds limit for monochrome normal",
			tone:            bitimage.Monochrome,
			horizontalScale: bitimage.NormalScale,
			verticalScale:   bitimage.NormalScale,
			color:           bitimage.Color1,
			width:           100,
			height:          2401,
			data:            createRasterData(100, 2401),
			wantErr:         bitimage.ErrInvalidHeight,
		},
		{
			name:            "height exceeds limit for monochrome double",
			tone:            bitimage.Monochrome,
			horizontalScale: bitimage.NormalScale,
			verticalScale:   bitimage.DoubleScale,
			color:           bitimage.Color1,
			width:           100,
			height:          1201,
			data:            createRasterData(100, 1201),
			wantErr:         bitimage.ErrInvalidHeight,
		},
		{
			name:            "height exceeds limit for multiple tone normal",
			tone:            bitimage.MultipleTone,
			horizontalScale: bitimage.NormalScale,
			verticalScale:   bitimage.NormalScale,
			color:           bitimage.Color1,
			width:           100,
			height:          601,
			data:            createRasterData(100, 601),
			wantErr:         bitimage.ErrInvalidHeight,
		},
		{
			name:            "height exceeds limit for multiple tone double",
			tone:            bitimage.MultipleTone,
			horizontalScale: bitimage.NormalScale,
			verticalScale:   bitimage.DoubleScale,
			color:           bitimage.Color1,
			width:           100,
			height:          301,
			data:            createRasterData(100, 301),
			wantErr:         bitimage.ErrInvalidHeight,
		},
		{
			name:            "invalid data length",
			tone:            bitimage.Monochrome,
			horizontalScale: bitimage.NormalScale,
			verticalScale:   bitimage.NormalScale,
			color:           bitimage.Color1,
			width:           100,
			height:          100,
			data:            []byte{0xFF}, // Should be more bytes
			wantErr:         bitimage.ErrInvalidDataLength,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cmd.StoreRasterGraphicsInBuffer(tt.tone, tt.horizontalScale, tt.verticalScale,
				tt.color, tt.width, tt.height, tt.data)

			if !testutils.AssertErrorOccurred(t, err, tt.wantErr != nil, "StoreRasterGraphicsInBuffer") {
				return
			}
			if tt.wantErr != nil {
				testutils.AssertError(t, err, tt.wantErr)
				return
			}

			// Verify the command structure
			if got[0] != shared.GS || got[1] != '(' || got[2] != 'L' {
				t.Errorf("StoreRasterGraphicsInBuffer: invalid command prefix")
			}
		})
	}
}

func TestGraphicsCommands_StoreRasterGraphicsInBufferLarge(t *testing.T) {
	cmd := bitimage.NewGraphicsCommands()

	// Helper to create testutils data for raster format
	createRasterData := func(width, height uint16) []byte {
		widthBytes := (int(width) + 7) / 8
		return testutils.RepeatByte(widthBytes*int(height), 0xFF)
	}

	tests := []struct {
		name            string
		tone            bitimage.GraphicsTone
		horizontalScale bitimage.GraphicsScale
		verticalScale   bitimage.GraphicsScale
		color           bitimage.GraphicsColor
		width           uint16
		height          uint16
		data            []byte
		wantErr         error
	}{
		{
			name:            "large monochrome data",
			tone:            bitimage.Monochrome,
			horizontalScale: bitimage.NormalScale,
			verticalScale:   bitimage.NormalScale,
			color:           bitimage.Color1,
			width:           2000,
			height:          2000,
			data:            createRasterData(2000, 2000),
			wantErr:         nil,
		},
		{
			name:            "invalid parameters same as standard",
			tone:            49,
			horizontalScale: bitimage.NormalScale,
			verticalScale:   bitimage.NormalScale,
			color:           bitimage.Color1,
			width:           100,
			height:          100,
			data:            createRasterData(100, 100),
			wantErr:         bitimage.ErrInvalidTone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cmd.StoreRasterGraphicsInBufferLarge(tt.tone, tt.horizontalScale, tt.verticalScale,
				tt.color, tt.width, tt.height, tt.data)

			if !testutils.AssertErrorOccurred(t, err, tt.wantErr != nil, "StoreRasterGraphicsInBufferLarge") {
				return
			}
			if tt.wantErr != nil {
				testutils.AssertError(t, err, tt.wantErr)
				return
			}

			// Verify the command structure for large format
			if got[0] != shared.GS || got[1] != '8' || got[2] != 'L' {
				t.Errorf("StoreRasterGraphicsInBufferLarge: invalid command prefix")
			}
		})
	}
}

func TestGraphicsCommands_StoreColumnGraphicsInBuffer(t *testing.T) {
	cmd := bitimage.NewGraphicsCommands()

	// Helper to create testutils data for column format
	createColumnData := func(width, height uint16) []byte {
		heightBytes := (int(height) + 7) / 8
		return testutils.RepeatByte(int(width)*heightBytes, 0xFF)
	}

	tests := []struct {
		name            string
		horizontalScale bitimage.GraphicsScale
		verticalScale   bitimage.GraphicsScale
		color           bitimage.GraphicsColor
		width           uint16
		height          uint16
		data            []byte
		wantErr         error
	}{
		{
			name:            "normal scale color 1",
			horizontalScale: bitimage.NormalScale,
			verticalScale:   bitimage.NormalScale,
			color:           bitimage.Color1,
			width:           100,
			height:          50,
			data:            createColumnData(100, 50),
			wantErr:         nil,
		},
		{
			name:            "double width color 2",
			horizontalScale: bitimage.DoubleScale,
			verticalScale:   bitimage.NormalScale,
			color:           bitimage.Color2,
			width:           200,
			height:          100,
			data:            createColumnData(200, 100),
			wantErr:         nil,
		},
		{
			name:            "double height color 3",
			horizontalScale: bitimage.NormalScale,
			verticalScale:   bitimage.DoubleScale,
			color:           bitimage.Color3,
			width:           150,
			height:          128,
			data:            createColumnData(150, 128),
			wantErr:         nil,
		},
		{
			name:            "maximum dimensions",
			horizontalScale: bitimage.NormalScale,
			verticalScale:   bitimage.NormalScale,
			color:           bitimage.Color1,
			width:           2048,
			height:          128,
			data:            createColumnData(2048, 128),
			wantErr:         nil,
		},
		{
			name:            "color 4 not supported",
			horizontalScale: bitimage.NormalScale,
			verticalScale:   bitimage.NormalScale,
			color:           bitimage.Color4,
			width:           100,
			height:          100,
			data:            createColumnData(100, 100),
			wantErr:         bitimage.ErrInvalidColor,
		},
		{
			name:            "invalid horizontal scale",
			horizontalScale: 0,
			verticalScale:   bitimage.NormalScale,
			color:           bitimage.Color1,
			width:           100,
			height:          100,
			data:            createColumnData(100, 100),
			wantErr:         bitimage.ErrInvalidScale,
		},
		{
			name:            "invalid vertical scale",
			horizontalScale: bitimage.NormalScale,
			verticalScale:   3,
			color:           bitimage.Color1,
			width:           100,
			height:          100,
			data:            createColumnData(100, 100),
			wantErr:         bitimage.ErrInvalidScale,
		},
		{
			name:            "width exceeds limit",
			horizontalScale: bitimage.NormalScale,
			verticalScale:   bitimage.NormalScale,
			color:           bitimage.Color1,
			width:           2049,
			height:          100,
			data:            createColumnData(2049, 100),
			wantErr:         bitimage.ErrInvalidWidth,
		},
		{
			name:            "height exceeds limit",
			horizontalScale: bitimage.NormalScale,
			verticalScale:   bitimage.NormalScale,
			color:           bitimage.Color1,
			width:           100,
			height:          129,
			data:            createColumnData(100, 129),
			wantErr:         bitimage.ErrInvalidHeight,
		},
		{
			name:            "invalid data length",
			horizontalScale: bitimage.NormalScale,
			verticalScale:   bitimage.NormalScale,
			color:           bitimage.Color1,
			width:           100,
			height:          100,
			data:            []byte{0xFF}, // Should be more bytes
			wantErr:         bitimage.ErrInvalidDataLength,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cmd.StoreColumnGraphicsInBuffer(tt.horizontalScale, tt.verticalScale,
				tt.color, tt.width, tt.height, tt.data)

			if !testutils.AssertErrorOccurred(t, err, tt.wantErr != nil, "StoreColumnGraphicsInBuffer") {
				return
			}
			if tt.wantErr != nil {
				testutils.AssertError(t, err, tt.wantErr)
				return
			}

			// Verify the command structure
			if got[0] != shared.GS || got[1] != '(' || got[2] != 'L' {
				t.Errorf("StoreColumnGraphicsInBuffer: invalid command prefix")
			}
		})
	}
}

func TestGraphicsCommands_StoreColumnGraphicsInBufferLarge(t *testing.T) {
	cmd := bitimage.NewGraphicsCommands()

	// Helper to create testutils data for column format
	createColumnData := func(width, height uint16) []byte {
		heightBytes := (int(height) + 7) / 8
		return testutils.RepeatByte(int(width)*heightBytes, 0xFF)
	}

	tests := []struct {
		name            string
		horizontalScale bitimage.GraphicsScale
		verticalScale   bitimage.GraphicsScale
		color           bitimage.GraphicsColor
		width           uint16
		height          uint16
		data            []byte
		wantErr         error
	}{
		{
			name:            "large column data",
			horizontalScale: bitimage.NormalScale,
			verticalScale:   bitimage.NormalScale,
			color:           bitimage.Color1,
			width:           2048,
			height:          128,
			data:            createColumnData(2048, 128),
			wantErr:         nil,
		},
		{
			name:            "color 4 not supported",
			horizontalScale: bitimage.NormalScale,
			verticalScale:   bitimage.NormalScale,
			color:           bitimage.Color4,
			width:           100,
			height:          100,
			data:            createColumnData(100, 100),
			wantErr:         bitimage.ErrInvalidColor,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cmd.StoreColumnGraphicsInBufferLarge(tt.horizontalScale, tt.verticalScale,
				tt.color, tt.width, tt.height, tt.data)

			if !testutils.AssertErrorOccurred(t, err, tt.wantErr != nil, "StoreColumnGraphicsInBufferLarge") {
				return
			}
			if tt.wantErr != nil {
				testutils.AssertError(t, err, tt.wantErr)
				return
			}

			// Verify the command structure for large format
			if got[0] != shared.GS || got[1] != '8' || got[2] != 'L' {
				t.Errorf("StoreColumnGraphicsInBufferLarge: invalid command prefix")
			}
		})
	}
}

// ============================================================================
// Validation Functions Tests
// ============================================================================

func TestValidateDensityFunctionCode(t *testing.T) {
	tests := []struct {
		name    string
		fn      bitimage.FunctionCode
		wantErr bool
	}{
		{"valid code 1", bitimage.FunctionCodeDensity1, false},
		{"valid code 49", bitimage.FunctionCodeDensity49, false},
		{"invalid code 2", bitimage.FunctionCodePrint2, true},
		{"invalid code 50", bitimage.FunctionCodePrint50, true},
		{"invalid code 0", 0, true},
		{"invalid code 99", 99, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := bitimage.ValidateDensityFunctionCode(tt.fn)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDensityFunctionCode(%v) error = %v, wantErr %v", tt.fn, err, tt.wantErr)
			}
		})
	}
}

func TestValidatePrintFunctionCode(t *testing.T) {
	tests := []struct {
		name    string
		fn      bitimage.FunctionCode
		wantErr bool
	}{
		{"valid code 2", bitimage.FunctionCodePrint2, false},
		{"valid code 50", bitimage.FunctionCodePrint50, false},
		{"invalid code 1", bitimage.FunctionCodeDensity1, true},
		{"invalid code 49", bitimage.FunctionCodeDensity49, true},
		{"invalid code 0", 0, true},
		{"invalid code 99", 99, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := bitimage.ValidatePrintFunctionCode(tt.fn)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePrintFunctionCode(%v) error = %v, wantErr %v", tt.fn, err, tt.wantErr)
			}
		})
	}
}

func TestValidateDotDensity(t *testing.T) {
	tests := []struct {
		name    string
		x       bitimage.DotDensity
		y       bitimage.DotDensity
		wantErr bool
	}{
		{"valid 180x180", bitimage.Density180x180, bitimage.Density180x180, false},
		{"valid 360x360", bitimage.Density360x360, bitimage.Density360x360, false},
		{"invalid x value", 49, bitimage.Density180x180, true},
		{"invalid y value", bitimage.Density180x180, 52, true},
		{"mismatched values", bitimage.Density180x180, bitimage.Density360x360, true},
		{"both invalid", 49, 52, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := bitimage.ValidateDotDensity(tt.x, tt.y)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDotDensity(%v, %v) error = %v, wantErr %v", tt.x, tt.y, err, tt.wantErr)
			}
		})
	}
}

func TestValidateGraphicsTone(t *testing.T) {
	tests := []struct {
		name    string
		tone    bitimage.GraphicsTone
		wantErr bool
	}{
		{"valid monochrome", bitimage.Monochrome, false},
		{"valid multiple tone", bitimage.MultipleTone, false},
		{"invalid 49", 49, true},
		{"invalid 51", 51, true},
		{"invalid 0", 0, true},
		{"invalid 99", 99, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := bitimage.ValidateGraphicsTone(tt.tone)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGraphicsTone(%v) error = %v, wantErr %v", tt.tone, err, tt.wantErr)
			}
		})
	}
}

func TestValidateGraphicsScale(t *testing.T) {
	tests := []struct {
		name    string
		scale   bitimage.GraphicsScale
		wantErr bool
	}{
		{"valid normal", bitimage.NormalScale, false},
		{"valid double", bitimage.DoubleScale, false},
		{"invalid 0", 0, true},
		{"invalid 3", 3, true},
		{"invalid 99", 99, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := bitimage.ValidateGraphicsScale(tt.scale)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGraphicsScale(%v) error = %v, wantErr %v", tt.scale, err, tt.wantErr)
			}
		})
	}
}

func TestValidateGraphicsColor(t *testing.T) {
	tests := []struct {
		name    string
		color   bitimage.GraphicsColor
		wantErr bool
	}{
		{"valid color 1", bitimage.Color1, false},
		{"valid color 2", bitimage.Color2, false},
		{"valid color 3", bitimage.Color3, false},
		{"valid color 4", bitimage.Color4, false},
		{"invalid 48", 48, true},
		{"invalid 53", 53, true},
		{"invalid 0", 0, true},
		{"invalid 99", 99, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := bitimage.ValidateGraphicsColor(tt.color)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGraphicsColor(%v) error = %v, wantErr %v", tt.color, err, tt.wantErr)
			}
		})
	}
}

func TestValidateRasterDimensions(t *testing.T) {
	tests := []struct {
		name          string
		width         uint16
		height        uint16
		tone          bitimage.GraphicsTone
		verticalScale bitimage.GraphicsScale
		wantErr       bool
	}{
		{"valid monochrome normal", 100, 100, bitimage.Monochrome, bitimage.NormalScale, false},
		{"valid monochrome double", 100, 1200, bitimage.Monochrome, bitimage.DoubleScale, false},
		{"valid multiple tone normal", 100, 600, bitimage.MultipleTone, bitimage.NormalScale, false},
		{"valid multiple tone double", 100, 300, bitimage.MultipleTone, bitimage.DoubleScale, false},
		{"max width", 2400, 100, bitimage.Monochrome, bitimage.NormalScale, false},
		{"max monochrome height normal", 100, 2400, bitimage.Monochrome, bitimage.NormalScale, false},
		{"max monochrome height double", 100, 1200, bitimage.Monochrome, bitimage.DoubleScale, false},
		{"max multiple tone height normal", 100, 600, bitimage.MultipleTone, bitimage.NormalScale, false},
		{"max multiple tone height double", 100, 300, bitimage.MultipleTone, bitimage.DoubleScale, false},
		{"width zero", 0, 100, bitimage.Monochrome, bitimage.NormalScale, true},
		{"height zero", 100, 0, bitimage.Monochrome, bitimage.NormalScale, true},
		{"width exceeds", 2401, 100, bitimage.Monochrome, bitimage.NormalScale, true},
		{"monochrome height exceeds normal", 100, 2401, bitimage.Monochrome, bitimage.NormalScale, true},
		{"monochrome height exceeds double", 100, 1201, bitimage.Monochrome, bitimage.DoubleScale, true},
		{"multiple tone height exceeds normal", 100, 601, bitimage.MultipleTone, bitimage.NormalScale, true},
		{"multiple tone height exceeds double", 100, 301, bitimage.MultipleTone, bitimage.DoubleScale, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := bitimage.ValidateRasterDimensions(tt.width, tt.height, tt.tone, tt.verticalScale)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateRasterDimensions(%v, %v, %v, %v) error = %v, wantErr %v",
					tt.width, tt.height, tt.tone, tt.verticalScale, err, tt.wantErr)
			}
		})
	}
}

func TestValidateColumnDimensions(t *testing.T) {
	tests := []struct {
		name    string
		width   uint16
		height  uint16
		wantErr bool
	}{
		{"minimum valid", 1, 1, false},
		{"typical dimensions", 100, 100, false},
		{"maximum width", 2048, 100, false},
		{"maximum height", 100, 128, false},
		{"maximum both", 2048, 128, false},
		{"width zero", 0, 100, true},
		{"height zero", 100, 0, true},
		{"width exceeds", 2049, 100, true},
		{"height exceeds", 100, 129, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := bitimage.ValidateColumnDimensions(tt.width, tt.height)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateColumnDimensions(%v, %v) error = %v, wantErr %v",
					tt.width, tt.height, err, tt.wantErr)
			}
		})
	}
}

// ============================================================================
// Helper Functions Tests
// ============================================================================

func TestCalculateRasterDataSize(t *testing.T) {
	tests := []struct {
		name   string
		width  uint16
		height uint16
		want   int
	}{
		{"8x8 pixels", 8, 8, 8},
		{"16x16 pixels", 16, 16, 32},
		{"100x50 pixels", 100, 50, 650},
		{"7x1 pixels", 7, 1, 1},
		{"9x1 pixels", 9, 1, 2},
		{"2400x1 pixels", 2400, 1, 300},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test through the command to ensure data validation works
			cmd := bitimage.NewGraphicsCommands()

			// Calculate expected size
			widthBytes := (int(tt.width) + 7) / 8
			expectedSize := widthBytes * int(tt.height)

			// Create data with expected size
			data := testutils.RepeatByte(expectedSize, 0xFF)

			// Should succeed with correct data size
			_, err := cmd.StoreRasterGraphicsInBuffer(bitimage.Monochrome, bitimage.NormalScale,
				bitimage.NormalScale, bitimage.Color1, tt.width, tt.height, data)
			if err != nil {
				t.Errorf("calculateRasterDataSize validation failed: %v", err)
			}

			// Should fail with incorrect data size
			if expectedSize > 0 {
				wrongData := testutils.RepeatByte(expectedSize-1, 0xFF)
				_, err = cmd.StoreRasterGraphicsInBuffer(bitimage.Monochrome, bitimage.NormalScale,
					bitimage.NormalScale, bitimage.Color1, tt.width, tt.height, wrongData)
				if err == nil {
					t.Errorf("calculateRasterDataSize should have failed for incorrect data length")
				}
			}
		})
	}
}

func TestCalculateColumnDataSize(t *testing.T) {
	tests := []struct {
		name   string
		width  uint16
		height uint16
		want   int
	}{
		{"8x8 pixels", 8, 8, 8},
		{"16x16 pixels", 16, 16, 32},
		{"100x50 pixels", 100, 50, 700},
		{"1x7 pixels", 1, 7, 1},
		{"1x9 pixels", 1, 9, 2},
		{"1x128 pixels", 1, 128, 16},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test through the command to ensure data validation works
			cmd := bitimage.NewGraphicsCommands()

			// Calculate expected size
			heightBytes := (int(tt.height) + 7) / 8
			expectedSize := int(tt.width) * heightBytes

			// Create data with expected size
			data := testutils.RepeatByte(expectedSize, 0xFF)

			// Should succeed with correct data size
			_, err := cmd.StoreColumnGraphicsInBuffer(bitimage.NormalScale, bitimage.NormalScale,
				bitimage.Color1, tt.width, tt.height, data)
			if err != nil {
				t.Errorf("calculateColumnDataSize validation failed: %v", err)
			}

			// Should fail with incorrect data size
			if expectedSize > 0 {
				wrongData := testutils.RepeatByte(expectedSize-1, 0xFF)
				_, err = cmd.StoreColumnGraphicsInBuffer(bitimage.NormalScale, bitimage.NormalScale,
					bitimage.Color1, tt.width, tt.height, wrongData)
				if err == nil {
					t.Errorf("calculateColumnDataSize should have failed for incorrect data length")
				}
			}
		})
	}
}

func TestGetMaxHeight(t *testing.T) {
	tests := []struct {
		name          string
		tone          bitimage.GraphicsTone
		verticalScale bitimage.GraphicsScale
		want          uint16
	}{
		{"monochrome normal", bitimage.Monochrome, bitimage.NormalScale, 2400},
		{"monochrome double", bitimage.Monochrome, bitimage.DoubleScale, 1200},
		{"multiple tone normal", bitimage.MultipleTone, bitimage.NormalScale, 600},
		{"multiple tone double", bitimage.MultipleTone, bitimage.DoubleScale, 300},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test through validation
			// Should succeed with max height
			err := bitimage.ValidateRasterDimensions(100, tt.want, tt.tone, tt.verticalScale)
			if err != nil {
				t.Errorf("getMaxHeight: should accept height %v for tone %v scale %v: %v",
					tt.want, tt.tone, tt.verticalScale, err)
			}

			// Should fail with max height + 1
			err = bitimage.ValidateRasterDimensions(100, tt.want+1, tt.tone, tt.verticalScale)
			if err == nil {
				t.Errorf("getMaxHeight: should reject height %v for tone %v scale %v",
					tt.want+1, tt.tone, tt.verticalScale)
			}
		})
	}
}
