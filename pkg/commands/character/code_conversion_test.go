package character_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/adcondev/poster/internal/testutils"
	"github.com/adcondev/poster/pkg/commands/character"
	"github.com/adcondev/poster/pkg/commands/common"
)

// TODO: Swap out hardcoded bytes and values with constants from the character package
// TODO: Swap boilerplate assertion code with calls to utils/testutils/assertions.go
// TODO: Define prefix bytes for commands to reduce repetition
// TODO: Swap for specific type, the convert to byte. This simulates real usage better.

// ============================================================================
// Code Conversion Commands Tests
// ============================================================================

func TestCodeConversionCommands_SelectCharacterEncodeSystem(t *testing.T) {
	cc := character.NewCodeConversionCommands()
	prefix := []byte{common.FS, '(', 'C', 0x02, 0x00, 0x30}

	tests := []struct {
		name     string
		encoding character.EncodeSystem
		want     []byte
		wantErr  error
	}{
		{
			name:     "1-byte encoding",
			encoding: character.OneByte,
			want:     append(prefix, 1),
			wantErr:  nil,
		},
		{
			name:     "UTF-8 encoding",
			encoding: character.UTF8,
			want:     append(prefix, 2),
			wantErr:  nil,
		},
		{
			name:     "1-byte encoding ASCII",
			encoding: character.OneByteASCII,
			want:     append(prefix, '1'),
			wantErr:  nil,
		},
		{
			name:     "UTF-8 encoding ASCII",
			encoding: character.UTF8Ascii,
			want:     append(prefix, '2'),
			wantErr:  nil,
		},
		{
			name:     "invalid encoding",
			encoding: 3,
			want:     nil,
			wantErr:  character.ErrEncoding,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cc.SelectCharacterEncodeSystem(tt.encoding)

			// Standardized error checking
			if !testutils.AssertErrorOccurred(t, err, tt.wantErr != nil, "SelectCharacterEncodeSystem") {
				return
			}

			// Check specific error type if expecting error
			if tt.wantErr != nil {
				testutils.AssertError(t, err, tt.wantErr)
				return
			}

			// Check result if no error expected
			testutils.AssertBytes(t, got, tt.want, "SelectCharacterEncodeSystem(%v)", tt.encoding)
		})
	}
}

func TestCodeConversionCommands_SetFontPriority(t *testing.T) {
	cc := character.NewCodeConversionCommands()

	tests := []struct {
		name     string
		priority character.FontPriority
		function character.FontFunction
		want     []byte
		wantErr  bool
	}{
		{
			name:     "first priority AnkSansSerif font",
			priority: 0,
			function: 0,
			want:     []byte{common.FS, '(', 'C', 0x03, 0x00, 0x3C, 0, 0},
			wantErr:  false,
		},
		{
			name:     "second priority Japanese",
			priority: 1,
			function: 11,
			want:     []byte{common.FS, '(', 'C', 0x03, 0x00, 0x3C, 1, 11},
			wantErr:  false,
		},
		{
			name:     "first priority Simplified Chinese",
			priority: 0,
			function: 20,
			want:     []byte{common.FS, '(', 'C', 0x03, 0x00, 0x3C, 0, 20},
			wantErr:  false,
		},
		{
			name:     "invalid priority",
			priority: 2,
			function: 0,
			want:     nil,
			wantErr:  true,
		},
		{
			name:     "invalid font type",
			priority: 0,
			function: 99,
			want:     nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cc.SetFontPriority(tt.priority, tt.function)

			// Standardized error checking
			if (err != nil) != tt.wantErr {
				t.Errorf("SetFontPriority(%v, %v) error = %v, wantErr %v",
					tt.priority, tt.function, err, tt.wantErr)
				return
			}

			var baseErr error
			switch tt.name {
			case "invalid priority":
				baseErr = character.ErrFontPriority
			case "invalid font type":
				baseErr = character.ErrFontType
			default:
				baseErr = nil
			}

			// Check specific error type if expecting error
			if tt.wantErr && err != nil {
				if !errors.Is(err, baseErr) {
					t.Errorf("SetFontPriority(%v, %v) error = %v, want %v",
						tt.priority, tt.function, err, baseErr)
				}
				if !errors.Is(err, baseErr) {
					t.Errorf("SetFontPriority(%v, %v) error = %v, want %v",
						tt.priority, tt.function, err, baseErr)
				}
				return
			}

			// Check result if no error expected
			if !tt.wantErr && !bytes.Equal(got, tt.want) {
				t.Errorf("SetFontPriority(%v, %v) = %#v, want %#v",
					tt.priority, tt.function, got, tt.want)
			}
		})
	}
}
