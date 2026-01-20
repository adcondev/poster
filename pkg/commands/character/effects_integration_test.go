package character

import (
	"testing"
)

func TestIntegration_Effects_ColorCombinations(t *testing.T) {
	cmd := NewCommands()

	t.Run("promotional text with all effects", func(t *testing.T) {
		var buffer = make([]byte, 0, 32)

		// Apply character color
		charColor, err := cmd.Effects.SelectCharacterColor(CharColor2)
		if err != nil {
			t.Fatalf("SelectCharacterColor: %v", err)
		}
		buffer = append(buffer, charColor...)

		// Apply background color
		bgColor, err := cmd.Effects.SelectBackgroundColor(BackgroundColor1)
		if err != nil {
			t.Fatalf("SelectBackgroundColor: %v", err)
		}
		buffer = append(buffer, bgColor...)

		// Enable shadow
		shadow, err := cmd.Effects.SetCharacterShadowMode(
			ShadowModeOnByte,
			ShadowColor3,
		)
		if err != nil {
			t.Fatalf("SetCharacterShadowMode: %v", err)
		}
		buffer = append(buffer, shadow...)

		// Combine with reverse mode
		buffer = append(buffer, cmd.SetWhiteBlackReverseMode(OnRm)...)

		if len(buffer) != 25 { // 7 + 7 + 8 + 3 bytes
			t.Errorf("Buffer length = %d, want 25", len(buffer))
		}
	})

	t.Run("effect reset workflow", func(t *testing.T) {
		// Turn off all effects
		charCmd, _ := cmd.Effects.SelectCharacterColor(CharColorNone)
		bgCmd, _ := cmd.Effects.SelectBackgroundColor(BackgroundColorNone)
		shadowCmd, _ := cmd.Effects.SetCharacterShadowMode(
			ShadowModeOffByte,
			ShadowColorNone,
		)

		totalLen := len(charCmd) + len(bgCmd) + len(shadowCmd)
		if totalLen != 22 { // 7 + 7 + 8 bytes
			t.Errorf("Reset commands length = %d, want 22", totalLen)
		}
	})
}
