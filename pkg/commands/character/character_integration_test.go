package character

import (
	"bytes"
	"testing"

	"github.com/adcondev/poster/pkg/commands/shared"
)

func TestIntegration_Character_BasicFormatting(t *testing.T) {
	cmd := NewCommands()

	t.Run("receipt header formatting", func(t *testing.T) {
		var buffer = make([]byte, 0, 32)

		// Title formatting
		titleSize, _ := NewSize(2, 3)
		buffer = append(buffer, cmd.SelectCharacterSize(titleSize)...)
		buffer = append(buffer, cmd.SetEmphasizedMode(OnEm)...)

		// Subtitle formatting
		normalSize, _ := NewSize(1, 1)
		buffer = append(buffer, cmd.SelectCharacterSize(normalSize)...)
		buffer = append(buffer, cmd.SetEmphasizedMode(OffEm)...)
		underlineCmd, _ := cmd.SetUnderlineMode(OneDot)
		buffer = append(buffer, underlineCmd...)

		if len(buffer) < 15 {
			t.Error("Buffer should contain multiple formatting commands")
		}
	})

	t.Run("print modes combination", func(t *testing.T) {
		// All effects combined
		modes := EmphasizedOnPm |
			DoubleHeightOnPm |
			DoubleWidthOnPm |
			UnderlineOnPm

		result := cmd.SelectPrintModes(modes)

		expected := []byte{shared.ESC, '!', 0xB8}
		if !bytes.Equal(result, expected) {
			t.Errorf("SelectPrintModes = %#v, want %#v", result, expected)
		}
	})

	t.Run("character transformations", func(t *testing.T) {
		var buffer = make([]byte, 0, 16)

		rotationCmd, _ := cmd.Set90DegreeClockwiseRotationMode(On90Dot1)
		buffer = append(buffer, rotationCmd...)
		buffer = append(buffer, cmd.SetUpsideDownMode(OnUdm)...)
		buffer = append(buffer, cmd.SetWhiteBlackReverseMode(OnRm)...)

		if len(buffer) != 9 { // 3 commands × 3 bytes each
			t.Errorf("Buffer length = %d, want 9", len(buffer))
		}
	})
}

func TestIntegration_Character_ErrorHandling(t *testing.T) {
	cmd := NewCommands()

	// Test invalid parameters cascade
	_, err := cmd.SetUnderlineMode(99)
	if err == nil {
		t.Error("Invalid underline mode should return error")
	}

	_, err = NewSize(0, 1)
	if err == nil {
		t.Error("Invalid character size should return error")
	}
}
