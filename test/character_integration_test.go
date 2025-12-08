package test_test

import (
	"bytes"
	"testing"

	"github.com/adcondev/poster/pkg/commands/character"
	"github.com/adcondev/poster/pkg/commands/common"
)

func TestIntegration_Character_BasicFormatting(t *testing.T) {
	cmd := character.NewCommands()

	t.Run("receipt header formatting", func(t *testing.T) {
		var buffer []byte

		// Title formatting
		titleSize, _ := character.NewSize(2, 3)
		buffer = append(buffer, cmd.SelectCharacterSize(titleSize)...)
		buffer = append(buffer, cmd.SetEmphasizedMode(character.OnEm)...)

		// Subtitle formatting
		normalSize, _ := character.NewSize(1, 1)
		buffer = append(buffer, cmd.SelectCharacterSize(normalSize)...)
		buffer = append(buffer, cmd.SetEmphasizedMode(character.OffEm)...)
		underlineCmd, _ := cmd.SetUnderlineMode(character.OneDot)
		buffer = append(buffer, underlineCmd...)

		if len(buffer) < 15 {
			t.Error("Buffer should contain multiple formatting commands")
		}
	})

	t.Run("print modes combination", func(t *testing.T) {
		// All effects combined
		modes := character.EmphasizedOnPm |
			character.DoubleHeightOnPm |
			character.DoubleWidthOnPm |
			character.UnderlineOnPm

		result := cmd.SelectPrintModes(modes)

		expected := []byte{common.ESC, '!', 0xB8}
		if !bytes.Equal(result, expected) {
			t.Errorf("SelectPrintModes = %#v, want %#v", result, expected)
		}
	})

	t.Run("character transformations", func(t *testing.T) {
		var buffer []byte

		rotationCmd, _ := cmd.Set90DegreeClockwiseRotationMode(character.On90Dot1)
		buffer = append(buffer, rotationCmd...)
		buffer = append(buffer, cmd.SetUpsideDownMode(character.OnUdm)...)
		buffer = append(buffer, cmd.SetWhiteBlackReverseMode(character.OnRm)...)

		if len(buffer) != 9 { // 3 commands × 3 bytes each
			t.Errorf("Buffer length = %d, want 9", len(buffer))
		}
	})
}

func TestIntegration_Character_ErrorHandling(t *testing.T) {
	cmd := character.NewCommands()

	// Test invalid parameters cascade
	_, err := cmd.SetUnderlineMode(99)
	if err == nil {
		t.Error("Invalid underline mode should return error")
	}

	_, err = character.NewSize(0, 1)
	if err == nil {
		t.Error("Invalid character size should return error")
	}
}
