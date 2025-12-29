package character_test

import (
	"testing"

	"github.com/adcondev/poster/internal/testutils"
	"github.com/adcondev/poster/pkg/commands/character"
	"github.com/adcondev/poster/pkg/commands/shared"
)

// ============================================================================
// Effects Commands Tests
// ============================================================================

func TestEffectsCommands_SelectCharacterColor(t *testing.T) {
	// Setup
	ec := character.NewEffectsCommands()
	prefix := []byte{shared.GS, '(', 'N', 0x02, 0x00, 0x30}

	tests := []struct {
		name    string
		color   byte
		want    []byte
		wantErr error
	}{
		{
			name:    "no color",
			color:   character.CharColorNone,
			want:    append(prefix, '0'),
			wantErr: nil,
		},
		{
			name:    "color 1",
			color:   character.CharColor1,
			want:    append(prefix, '1'),
			wantErr: nil,
		},
		{
			name:    "color 2",
			color:   character.CharColor2,
			want:    append(prefix, '2'),
			wantErr: nil,
		},
		{
			name:    "color 3",
			color:   character.CharColor3,
			want:    append(prefix, '3'),
			wantErr: nil,
		},
		{
			name:    "invalid color",
			color:   '4',
			want:    nil,
			wantErr: character.ErrInvalidCharacterColor,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Execute
			got, err := ec.SelectCharacterColor(tt.color)

			// Verify error
			if !testutils.AssertErrorOccurred(t, err, tt.wantErr != nil, "SelectCharacterColor") {
				return
			}
			if tt.wantErr != nil {
				testutils.AssertError(t, err, tt.wantErr)
				return
			}

			// Verify result
			testutils.AssertBytes(t, got, tt.want, "SelectCharacterColor(%v)", tt.color)
		})
	}
}

func TestEffectsCommands_SelectBackgroundColor(t *testing.T) {
	// Setup
	ec := character.NewEffectsCommands()
	prefix := []byte{shared.GS, '(', 'N', 0x02, 0x00, 0x31}

	tests := []struct {
		name    string
		color   byte
		want    []byte
		wantErr error
	}{
		{
			name:    "no background",
			color:   character.BackgroundColorNone,
			want:    append(prefix, '0'),
			wantErr: nil,
		},
		{
			name:    "background color 1",
			color:   character.BackgroundColor1,
			want:    append(prefix, '1'),
			wantErr: nil,
		},
		{
			name:    "background color 2",
			color:   character.BackgroundColor2,
			want:    append(prefix, '2'),
			wantErr: nil,
		},
		{
			name:    "background color 3",
			color:   character.BackgroundColor3,
			want:    append(prefix, '3'),
			wantErr: nil,
		},
		{
			name:    "invalid background",
			color:   '5',
			want:    nil,
			wantErr: character.ErrInvalidBackgroundColor,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Execute
			got, err := ec.SelectBackgroundColor(tt.color)

			// Verify error
			if !testutils.AssertErrorOccurred(t, err, tt.wantErr != nil, "SelectBackgroundColor") {
				return
			}
			if tt.wantErr != nil {
				testutils.AssertError(t, err, tt.wantErr)
				return
			}

			// Verify result
			testutils.AssertBytes(t, got, tt.want, "SelectBackgroundColor(%v)", tt.color)
		})
	}
}

func TestEffectsCommands_SetCharacterShadowMode(t *testing.T) {
	// Setup
	ec := &character.EffectsCommands{}
	prefix := []byte{shared.GS, '(', 'N', 0x03, 0x00, 0x32}

	tests := []struct {
		name        string
		shadowMode  byte
		shadowColor byte
		want        []byte
		wantErr     error
	}{
		{
			name:        "shadow off no color",
			shadowMode:  character.ShadowModeOffByte,
			shadowColor: character.ShadowColorNone,
			want:        append(prefix, 0x00, '0'),
			wantErr:     nil,
		},
		{
			name:        "shadow on color 1",
			shadowMode:  character.ShadowModeOnByte,
			shadowColor: character.ShadowColor1,
			want:        append(prefix, 0x01, '1'),
			wantErr:     nil,
		},
		{
			name:        "shadow off ASCII",
			shadowMode:  character.ShadowModeOffASCII,
			shadowColor: character.ShadowColor2,
			want:        append(prefix, '0', '2'),
			wantErr:     nil,
		},
		{
			name:        "shadow on ASCII",
			shadowMode:  character.ShadowModeOnASCII,
			shadowColor: character.ShadowColor3,
			want:        append(prefix, '1', '3'),
			wantErr:     nil,
		},
		{
			name:        "invalid shadow mode",
			shadowMode:  2,
			shadowColor: character.ShadowColor1,
			want:        nil,
			wantErr:     character.ErrInvalidShadowMode,
		},
		{
			name:        "invalid shadow color",
			shadowMode:  character.ShadowModeOffByte,
			shadowColor: '4',
			want:        nil,
			wantErr:     character.ErrInvalidShadowColor,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Execute
			got, err := ec.SetCharacterShadowMode(tt.shadowMode, tt.shadowColor)

			// Verify error
			if !testutils.AssertErrorOccurred(t, err, tt.wantErr != nil, "SetCharacterShadowMode") {
				return
			}
			if tt.wantErr != nil {
				testutils.AssertError(t, err, tt.wantErr)
				return
			}

			// Verify result
			testutils.AssertBytes(t, got, tt.want, "SetCharacterShadowMode(%v, %v)", tt.shadowMode, tt.shadowColor)
		})
	}
}
