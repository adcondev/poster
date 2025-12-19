package character

import (
	"bytes"
	"testing"
)

func TestIntegration_UserDefined_CustomLogoWorkflow(t *testing.T) {
	cmd := NewCommands()

	t.Run("define and use custom characters", func(t *testing.T) {
		var buffer []byte

		// Create 4-part logo
		logoChars := make([]UserDefinedChar, 4)
		for i := range logoChars {
			logoChars[i] = UserDefinedChar{
				Width: 12,
				Data:  bytes.Repeat([]byte{byte(0x01 << i)}, 36), // Pattern for each part
			}
		}

		// Define characters 64-67
		defineCmd, err := cmd.UserDefined.DefineUserDefinedCharacters(3, 64, 67, logoChars)
		if err != nil {
			t.Fatalf("DefineUserDefinedCharacters: %v", err)
		}
		buffer = append(buffer, defineCmd...)

		// Enable user-defined character set
		buffer = append(buffer, cmd.UserDefined.SelectUserDefinedCharacterSet(UserDefinedOn)...)

		// Later disable user-defined set
		buffer = append(buffer, cmd.UserDefined.SelectUserDefinedCharacterSet(UserDefinedOff)...)

		if len(buffer) < 150 {
			t.Error("Buffer should contain logo definition commands")
		}

		// Verify command structure
		if !bytes.Contains(buffer[:10], []byte{0x1B, 0x26}) {
			t.Error("Buffer should start with define characters command")
		}
	})

	t.Run("character replacement workflow", func(t *testing.T) {
		// Replace a single character
		customChar := []UserDefinedChar{{
			Width: 8,
			Data:  bytes.Repeat([]byte{0xAA}, 24), // 8 width × 3 height
		}}

		defineCmd, err := cmd.UserDefined.DefineUserDefinedCharacters(3, 65, 65, customChar)
		if err != nil {
			t.Fatalf("DefineUserDefinedCharacters: %v", err)
		}

		// Cancel the character
		cancelCmd, err := cmd.UserDefined.CancelUserDefinedCharacter(65)
		if err != nil {
			t.Fatalf("CancelUserDefinedCharacter: %v", err)
		}

		if len(defineCmd) != 30 || len(cancelCmd) != 3 {
			t.Error("Command lengths incorrect")
		}
	})
}
