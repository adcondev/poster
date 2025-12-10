package common_test

import (
	"errors"
	"testing"

	"github.com/adcondev/poster/pkg/commands/common"
)

func TestUtils_IsBufOk_ValidInput(t *testing.T) {
	tests := []struct {
		name    string
		buf     []byte
		wantErr error
	}{
		{"empty buffer", []byte{}, common.ErrEmptyBuffer},
		{"valid buffer", []byte{1, 2, 3}, nil},
		{"max buffer", make([]byte, common.MaxBuf), nil},
		{"overflow buffer", make([]byte, common.MaxBuf+1), common.ErrBufferOverflow},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := common.IsBufLenOk(tt.buf)
			if !errors.Is(tt.wantErr, err) {
				t.Errorf("IsBufLenOk len=%d error = %v; want %v", len(tt.buf), err, tt.wantErr)
			}
		})
	}
}

func TestUtils_LengthLowHigh_ValidInput(t *testing.T) {
	tests := []struct {
		length uint16
		wantDL byte
		wantDH byte
	}{
		{0, 0, 0},
		{1, 1, 0},
		{0x1234, 0x34, 0x12},
		{0xFFFF, 0xFF, 0xFF},
	}
	for _, tt := range tests {
		dL, dH := common.ToLittleEndian(tt.length)

		if dL != tt.wantDL || dH != tt.wantDH {
			t.Errorf("ToLittleEndian(%d) = (%#x,%#x); want (%#x,%#x)", tt.length, dL, dH, tt.wantDL, tt.wantDH)
		}

	}
}
