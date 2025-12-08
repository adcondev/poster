package character_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/adcondev/poster/internal/testutils"
	"github.com/adcondev/poster/pkg/commands/character"
)

// ============================================================================
// User Defined Commands Tests
// ============================================================================

func TestUserDefined_SelectUserDefinedCharacterSet(t *testing.T) {
	// Setup
	udc := &character.UserDefinedCommands{}
	prefix := []byte{0x1B, 0x25}

	tests := []struct {
		name    string
		charSet byte
		want    []byte
	}{
		{
			name:    "user-defined off",
			charSet: character.UserDefinedOff,
			want:    append(prefix, 0x00),
		},
		{
			name:    "user-defined on",
			charSet: character.UserDefinedOn,
			want:    append(prefix, 0x01),
		},
		{
			name:    "user-defined off ASCII",
			charSet: character.UserDefinedOffASCII,
			want:    append(prefix, '0'),
		},
		{
			name:    "user-defined on ASCII",
			charSet: character.UserDefinedOnASCII,
			want:    append(prefix, '1'),
		},
		{
			name:    "any even number (LSB=0)",
			charSet: 0xFE,
			want:    append(prefix, 0xFE),
		},
		{
			name:    "any odd number (LSB=1)",
			charSet: 0xFF,
			want:    append(prefix, 0xFF),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Execute
			got := udc.SelectUserDefinedCharacterSet(tt.charSet)

			// Verify
			testutils.AssertBytes(t, got, tt.want, "SelectUserDefinedCharacterSet(%d)", tt.charSet)
		})
	}
}

func TestUserDefinedCommands_DefineUserDefinedCharacters(t *testing.T) {
	// Setup
	udc := &character.UserDefinedCommands{}

	tests := []struct {
		name        string
		height      byte
		startCode   byte
		endCode     byte
		definitions []character.UserDefinedChar
		wantPrefix  []byte
		wantErr     error
	}{
		{
			name:      "single character definition",
			height:    3,
			startCode: 65,
			endCode:   65,
			definitions: []character.UserDefinedChar{
				{Width: 5, Data: testutils.RepeatByte(15, 0xFF)}, // 3 height × 5 width
			},
			wantPrefix: []byte{0x1B, 0x26, 3, 65, 65},
			wantErr:    nil,
		},
		{
			name:      "multiple character definitions",
			height:    2,
			startCode: 65,
			endCode:   66,
			definitions: []character.UserDefinedChar{
				{Width: 2, Data: []byte{0xFF, 0x00, 0xFF, 0x00}},
				{Width: 3, Data: []byte{0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF}},
			},
			wantPrefix: []byte{0x1B, 0x26, 2, 65, 66},
			wantErr:    nil,
		},
		{
			name:      "zero width character",
			height:    3,
			startCode: character.UserDefinedMinCode,
			endCode:   character.UserDefinedMinCode,
			definitions: []character.UserDefinedChar{
				{Width: 0, Data: nil},
			},
			wantPrefix: []byte{0x1B, 0x26, 3, character.UserDefinedMinCode, character.UserDefinedMinCode},
			wantErr:    nil,
		},
		{
			name:      "maximum character code",
			height:    1,
			startCode: character.UserDefinedMaxCode,
			endCode:   character.UserDefinedMaxCode,
			definitions: []character.UserDefinedChar{
				{Width: 1, Data: []byte{0xFF}},
			},
			wantPrefix: []byte{0x1B, 0x26, 1, character.UserDefinedMaxCode, character.UserDefinedMaxCode},
			wantErr:    nil,
		},
		{
			name:      "invalid y value",
			height:    0,
			startCode: 65,
			endCode:   65,
			definitions: []character.UserDefinedChar{
				{Width: 5, Data: []byte{0xFF}},
			},
			wantPrefix: nil,
			wantErr:    character.ErrYValue,
		},
		{
			name:      "invalid character code",
			height:    3,
			startCode: 31, // Below minimum
			endCode:   31,
			definitions: []character.UserDefinedChar{
				{Width: 5, Data: []byte{0xFF}},
			},
			wantPrefix: nil,
			wantErr:    character.ErrCharacterCode,
		},
		{
			name:      "invalid code range",
			height:    3,
			startCode: 66,
			endCode:   65, // End before start
			definitions: []character.UserDefinedChar{
				{Width: 5, Data: []byte{0xFF}},
			},
			wantPrefix: nil,
			wantErr:    character.ErrCodeRange,
		},
		{
			name:      "definition count mismatch",
			height:    3,
			startCode: 65,
			endCode:   67, // Range of 3
			definitions: []character.UserDefinedChar{
				{Width: 5, Data: testutils.RepeatByte(15, 0xFF)}, // Only 1 definition
			},
			wantPrefix: nil,
			wantErr:    character.ErrDefinition,
		},
		{
			name:      "invalid data length",
			height:    3,
			startCode: 65,
			endCode:   65,
			definitions: []character.UserDefinedChar{
				{Width: 2, Data: []byte{0xFF}}, // Should be 6 bytes (3*2)
			},
			wantPrefix: nil,
			wantErr:    character.ErrDataLength,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Execute
			got, err := udc.DefineUserDefinedCharacters(tt.height, tt.startCode, tt.endCode, tt.definitions)

			// Verify error
			if !testutils.AssertErrorOccurred(t, err, tt.wantErr != nil, "DefineUserDefinedCharacters") {
				return
			}
			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("DefineUserDefinedCharacters() error = %v, want error containing %v", err, tt.wantErr)
				return
			}

			// Verify result prefix if no error
			if tt.wantErr == nil && len(got) >= len(tt.wantPrefix) {
				if !bytes.Equal(got[:len(tt.wantPrefix)], tt.wantPrefix) {
					t.Errorf("DefineUserDefinedCharacters() prefix = %#v, want %#v",
						got[:len(tt.wantPrefix)], tt.wantPrefix)
				}
			}
		})
	}
}

func TestUserDefined_CancelUserDefinedCharacter(t *testing.T) {
	// Setup
	udc := &character.UserDefinedCommands{}
	prefix := []byte{0x1B, 0x3F}

	tests := []struct {
		name     string
		charCode byte
		want     []byte
		wantErr  error
	}{
		{
			name:     "cancel minimum code",
			charCode: character.UserDefinedMinCode,
			want:     append(prefix, 0x20),
			wantErr:  nil,
		},
		{
			name:     "cancel typical code",
			charCode: 65,
			want:     append(prefix, 0x41),
			wantErr:  nil,
		},
		{
			name:     "cancel maximum code",
			charCode: character.UserDefinedMaxCode,
			want:     append(prefix, 0x7E),
			wantErr:  nil,
		},
		{
			name:     "invalid code too low",
			charCode: 31,
			want:     nil,
			wantErr:  character.ErrCharacterCode,
		},
		{
			name:     "invalid code too high",
			charCode: 127,
			want:     nil,
			wantErr:  character.ErrCharacterCode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Execute
			got, err := udc.CancelUserDefinedCharacter(tt.charCode)

			// Verify error
			if !testutils.AssertErrorOccurred(t, err, tt.wantErr != nil, "CancelUserDefinedCharacter") {
				return
			}
			if tt.wantErr != nil {
				testutils.AssertError(t, err, tt.wantErr)
				return
			}

			// Verify result
			testutils.AssertBytes(t, got, tt.want, "CancelUserDefinedCharacter(%v)", tt.charCode)
		})
	}
}
