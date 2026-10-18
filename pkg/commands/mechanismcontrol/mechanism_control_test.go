package mechanismcontrol_test

import (
	"errors"
	"testing"

	"github.com/adcondev/poster/internal/testutils"
	"github.com/adcondev/poster/pkg/commands/mechanismcontrol"
	"github.com/adcondev/poster/pkg/commands/shared"
)

func TestCommands_Pulse(t *testing.T) {
	cmd := mechanismcontrol.NewCommands()

	tests := []struct {
		name        string
		pin         byte
		onTime      byte
		offTime     byte
		want        []byte
		expectedErr error
	}{
		{
			name:        "valid pin 0",
			pin:         0,
			onTime:      50,
			offTime:     50,
			want:        []byte{shared.ESC, 0x70, 0, 50, 50},
			expectedErr: nil,
		},
		{
			name:        "valid pin 1",
			pin:         1,
			onTime:      100,
			offTime:     200,
			want:        []byte{shared.ESC, 0x70, 1, 100, 200},
			expectedErr: nil,
		},
		{
			name:        "invalid pin 2",
			pin:         2,
			onTime:      50,
			offTime:     50,
			want:        nil,
			expectedErr: mechanismcontrol.ErrInvalidDrawerPin,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cmd.Pulse(tt.pin, tt.onTime, tt.offTime)
			if tt.expectedErr != nil {
				if !errors.Is(err, tt.expectedErr) {
					t.Fatalf("Pulse() error = %v, want %v", err, tt.expectedErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Pulse() unexpected error: %v", err)
			}
			testutils.AssertBytes(t, got, tt.want, "Pulse()")
		})
	}
}

func TestCommands_Beep(t *testing.T) {
	cmd := mechanismcontrol.NewCommands()

	tests := []struct {
		name        string
		times       byte
		lapse       byte
		want        []byte
		expectedErr error
	}{
		{
			name:        "valid beep",
			times:       2,
			lapse:       3,
			want:        []byte{shared.ESC, 0x42, 2, 3},
			expectedErr: nil,
		},
		{
			name:        "invalid times zero",
			times:       0,
			lapse:       3,
			want:        nil,
			expectedErr: mechanismcontrol.ErrInvalidBeepParams,
		},
		{
			name:        "invalid lapse zero",
			times:       2,
			lapse:       0,
			want:        nil,
			expectedErr: mechanismcontrol.ErrInvalidBeepParams,
		},
		{
			name:        "invalid times and lapse zero",
			times:       0,
			lapse:       0,
			want:        nil,
			expectedErr: mechanismcontrol.ErrInvalidBeepParams,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cmd.Beep(tt.times, tt.lapse)
			if tt.expectedErr != nil {
				if !errors.Is(err, tt.expectedErr) {
					t.Fatalf("Beep() error = %v, want %v", err, tt.expectedErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Beep() unexpected error: %v", err)
			}
			testutils.AssertBytes(t, got, tt.want, "Beep()")
		})
	}
}
