package mechanismcontrol_test

import (
	"testing"

	"github.com/adcondev/poster/internal/testutils"
	"github.com/adcondev/poster/pkg/commands/mechanismcontrol"
	"github.com/adcondev/poster/pkg/commands/shared"
)

// ============================================================================
// Print Head Control Tests
// ============================================================================

func TestCommands_ReturnHome(t *testing.T) {
	cmd := mechanismcontrol.NewCommands()

	got := cmd.ReturnHome()
	want := []byte{shared.ESC, '<'}

	testutils.AssertBytes(t, got, want, "ReturnHome()")
}

func TestCommands_SetUnidirectionalPrintMode(t *testing.T) {
	cmd := mechanismcontrol.NewCommands()
	prefix := testutils.BuildCommand(shared.ESC, 'U')

	tests := []struct {
		name string
		mode mechanismcontrol.UnidirectionalMode
		want []byte
	}{
		{
			name: "unidirectional off",
			mode: mechanismcontrol.UnidirOff,
			want: append(prefix, 0x00),
		},
		{
			name: "unidirectional on",
			mode: mechanismcontrol.UnidirOn,
			want: append(prefix, 0x01),
		},
		{
			name: "any even number (LSB=0)",
			mode: 0xFE,
			want: append(prefix, 0xFE),
		},
		{
			name: "any odd number (LSB=1)",
			mode: 0xFF,
			want: append(prefix, 0xFF),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cmd.SetUnidirectionalPrintMode(tt.mode)
			testutils.AssertBytes(t, got, tt.want, "SetUnidirectionalPrintMode(%v)", tt.mode)
		})
	}
}

// ============================================================================
// Modern Cut Commands Tests
// ============================================================================

func TestCommands_PaperCut(t *testing.T) {
	cmd := mechanismcontrol.NewCommands()
	prefix := []byte{shared.GS, '(', 'V', 0x02, 0x00, 0x30}

	tests := []struct {
		name    string
		cutType mechanismcontrol.CutType
		want    []byte
		wantErr error
	}{
		{
			name:    "full cut",
			cutType: mechanismcontrol.CutTypeFull,
			want:    append(prefix, 0x00),
			wantErr: nil,
		},
		{
			name:    "partial cut",
			cutType: mechanismcontrol.CutTypePartial,
			want:    append(prefix, 0x01),
			wantErr: nil,
		},
		{
			name:    "invalid cut type 2",
			cutType: 2,
			want:    nil,
			wantErr: mechanismcontrol.ErrCutType,
		},
		{
			name:    "invalid cut type 99",
			cutType: 99,
			want:    nil,
			wantErr: mechanismcontrol.ErrCutType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cmd.PaperCut(tt.cutType)

			if !testutils.AssertErrorOccurred(t, err, tt.wantErr != nil, "PartialPaperCut") {
				return
			}
			if tt.wantErr != nil {
				testutils.AssertError(t, err, tt.wantErr)
				return
			}

			testutils.AssertBytes(t, got, tt.want, "PartialPaperCut(%v)", tt.cutType)
		})
	}
}

func TestCommands_PaperFeedAndCut(t *testing.T) {
	cmd := mechanismcontrol.NewCommands()
	prefix := []byte{shared.GS, '(', 'V', 0x03, 0x00, 0x31}

	tests := []struct {
		name       string
		cutType    mechanismcontrol.CutType
		feedAmount byte
		want       []byte
		wantErr    error
	}{
		{
			name:       "full cut no feed",
			cutType:    mechanismcontrol.CutTypeFull,
			feedAmount: 0,
			want:       append(prefix, 0x00, 0),
			wantErr:    nil,
		},
		{
			name:       "partial cut no feed",
			cutType:    mechanismcontrol.CutTypePartial,
			feedAmount: 0,
			want:       append(prefix, 0x01, 0),
			wantErr:    nil,
		},
		{
			name:       "full cut with feed",
			cutType:    mechanismcontrol.CutTypeFull,
			feedAmount: 50,
			want:       append(prefix, 0x00, 50),
			wantErr:    nil,
		},
		{
			name:       "partial cut with feed",
			cutType:    mechanismcontrol.CutTypePartial,
			feedAmount: 100,
			want:       append(prefix, 0x01, 100),
			wantErr:    nil,
		},
		{
			name:       "full cut max feed",
			cutType:    mechanismcontrol.CutTypeFull,
			feedAmount: mechanismcontrol.MaxFeedAmount,
			want:       append(prefix, 0x00, 255),
			wantErr:    nil,
		},
		{
			name:       "invalid cut type",
			cutType:    5,
			feedAmount: 50,
			want:       nil,
			wantErr:    mechanismcontrol.ErrCutType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cmd.PaperFeedAndCut(tt.cutType, tt.feedAmount)

			if !testutils.AssertErrorOccurred(t, err, tt.wantErr != nil, "PaperFeedAndCut") {
				return
			}
			if tt.wantErr != nil {
				testutils.AssertError(t, err, tt.wantErr)
				return
			}

			testutils.AssertBytes(t, got, tt.want, "PaperFeedAndCut(%v, %d)", tt.cutType, tt.feedAmount)
		})
	}
}

func TestCommands_ReservePaperCut(t *testing.T) {
	cmd := mechanismcontrol.NewCommands()
	prefix := []byte{shared.GS, '(', 'V', 0x03, 0x00, 0x33}

	tests := []struct {
		name       string
		cutType    mechanismcontrol.CutType
		feedAmount byte
		want       []byte
		wantErr    error
	}{
		{
			name:       "reserve full cut no feed",
			cutType:    mechanismcontrol.CutTypeFull,
			feedAmount: 0,
			want:       append(prefix, 0x00, 0),
			wantErr:    nil,
		},
		{
			name:       "reserve partial cut no feed",
			cutType:    mechanismcontrol.CutTypePartial,
			feedAmount: 0,
			want:       append(prefix, 0x01, 0),
			wantErr:    nil,
		},
		{
			name:       "reserve full cut with feed",
			cutType:    mechanismcontrol.CutTypeFull,
			feedAmount: 75,
			want:       append(prefix, 0x00, 75),
			wantErr:    nil,
		},
		{
			name:       "reserve partial cut with feed",
			cutType:    mechanismcontrol.CutTypePartial,
			feedAmount: 150,
			want:       append(prefix, 0x01, 150),
			wantErr:    nil,
		},
		{
			name:       "reserve full cut max feed",
			cutType:    mechanismcontrol.CutTypeFull,
			feedAmount: 255,
			want:       append(prefix, 0x00, 255),
			wantErr:    nil,
		},
		{
			name:       "invalid cut type",
			cutType:    10,
			feedAmount: 50,
			want:       nil,
			wantErr:    mechanismcontrol.ErrCutType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cmd.ReservePaperCut(tt.cutType, tt.feedAmount)

			if !testutils.AssertErrorOccurred(t, err, tt.wantErr != nil, "ReservePaperCut") {
				return
			}
			if tt.wantErr != nil {
				testutils.AssertError(t, err, tt.wantErr)
				return
			}

			testutils.AssertBytes(t, got, tt.want, "ReservePaperCut(%v, %d)", tt.cutType, tt.feedAmount)
		})
	}
}

// ============================================================================
// Simple Cut Commands Tests
// ============================================================================

func TestCommands_CutPaper(t *testing.T) {
	cmd := mechanismcontrol.NewCommands()
	prefix := testutils.BuildCommand(shared.GS, 'V')

	tests := []struct {
		name    string
		mode    mechanismcontrol.CutMode
		want    []byte
		wantErr error
	}{
		{
			name:    "full cut (0)",
			mode:    mechanismcontrol.CutModeFull,
			want:    append(prefix, 0x00),
			wantErr: nil,
		},
		{
			name:    "partial cut (1)",
			mode:    mechanismcontrol.CutModePartial,
			want:    append(prefix, 0x01),
			wantErr: nil,
		},
		{
			name:    "full cut ASCII (48)",
			mode:    mechanismcontrol.CutModeFullASCII,
			want:    append(prefix, '0'),
			wantErr: nil,
		},
		{
			name:    "partial cut ASCII (49)",
			mode:    mechanismcontrol.CutModePartialASCII,
			want:    append(prefix, '1'),
			wantErr: nil,
		},
		{
			name:    "invalid mode 2",
			mode:    2,
			want:    nil,
			wantErr: mechanismcontrol.ErrCutMode,
		},
		{
			name:    "invalid mode 47",
			mode:    47,
			want:    nil,
			wantErr: mechanismcontrol.ErrCutMode,
		},
		{
			name:    "invalid mode 50",
			mode:    50,
			want:    nil,
			wantErr: mechanismcontrol.ErrCutMode,
		},
		{
			name:    "invalid mode 99",
			mode:    99,
			want:    nil,
			wantErr: mechanismcontrol.ErrCutMode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cmd.CutPaper(tt.mode)

			if !testutils.AssertErrorOccurred(t, err, tt.wantErr != nil, "CutPaper") {
				return
			}
			if tt.wantErr != nil {
				testutils.AssertError(t, err, tt.wantErr)
				return
			}

			testutils.AssertBytes(t, got, tt.want, "CutPaper(%v)", tt.mode)
		})
	}
}

func TestCommands_FeedAndCutPaper(t *testing.T) {
	cmd := mechanismcontrol.NewCommands()
	prefix := testutils.BuildCommand(shared.GS, 'V')

	tests := []struct {
		name       string
		mode       mechanismcontrol.FeedCut
		feedAmount byte
		want       []byte
		wantErr    error
	}{
		{
			name:       "full cut no feed",
			mode:       mechanismcontrol.FeedCutFull,
			feedAmount: 0,
			want:       append(prefix, 65, 0),
			wantErr:    nil,
		},
		{
			name:       "partial cut no feed",
			mode:       mechanismcontrol.FeedCutPartial,
			feedAmount: 0,
			want:       append(prefix, 66, 0),
			wantErr:    nil,
		},
		{
			name:       "full cut with feed",
			mode:       mechanismcontrol.FeedCutFull,
			feedAmount: 30,
			want:       append(prefix, 65, 30),
			wantErr:    nil,
		},
		{
			name:       "partial cut with feed",
			mode:       mechanismcontrol.FeedCutPartial,
			feedAmount: 60,
			want:       append(prefix, 66, 60),
			wantErr:    nil,
		},
		{
			name:       "full cut max feed",
			mode:       mechanismcontrol.FeedCutFull,
			feedAmount: 255,
			want:       append(prefix, 65, 255),
			wantErr:    nil,
		},
		{
			name:       "invalid mode 64",
			mode:       64,
			feedAmount: 50,
			want:       nil,
			wantErr:    mechanismcontrol.ErrFeedCutMode,
		},
		{
			name:       "invalid mode 67",
			mode:       67,
			feedAmount: 50,
			want:       nil,
			wantErr:    mechanismcontrol.ErrFeedCutMode,
		},
		{
			name:       "invalid mode 0",
			mode:       0,
			feedAmount: 50,
			want:       nil,
			wantErr:    mechanismcontrol.ErrFeedCutMode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cmd.FeedAndCutPaper(tt.mode, tt.feedAmount)

			if !testutils.AssertErrorOccurred(t, err, tt.wantErr != nil, "FeedAndCutPaper") {
				return
			}
			if tt.wantErr != nil {
				testutils.AssertError(t, err, tt.wantErr)
				return
			}

			testutils.AssertBytes(t, got, tt.want, "FeedAndCutPaper(%v, %d)", tt.mode, tt.feedAmount)
		})
	}
}

func TestCommands_SetCutPosition(t *testing.T) {
	cmd := mechanismcontrol.NewCommands()
	prefix := testutils.BuildCommand(shared.GS, 'V')

	tests := []struct {
		name     string
		mode     mechanismcontrol.PositionCut
		position byte
		want     []byte
		wantErr  error
	}{
		{
			name:     "full cut at position 0",
			mode:     mechanismcontrol.PositionCutFull,
			position: 0,
			want:     append(prefix, 97, 0),
			wantErr:  nil,
		},
		{
			name:     "partial cut at position 0",
			mode:     mechanismcontrol.PositionCutPartial,
			position: 0,
			want:     append(prefix, 98, 0),
			wantErr:  nil,
		},
		{
			name:     "full cut at position 50",
			mode:     mechanismcontrol.PositionCutFull,
			position: 50,
			want:     append(prefix, 97, 50),
			wantErr:  nil,
		},
		{
			name:     "partial cut at position 100",
			mode:     mechanismcontrol.PositionCutPartial,
			position: 100,
			want:     append(prefix, 98, 100),
			wantErr:  nil,
		},
		{
			name:     "full cut at max position",
			mode:     mechanismcontrol.PositionCutFull,
			position: 255,
			want:     append(prefix, 97, 255),
			wantErr:  nil,
		},
		{
			name:     "invalid mode 96",
			mode:     96,
			position: 50,
			want:     nil,
			wantErr:  mechanismcontrol.ErrPositionCutMode,
		},
		{
			name:     "invalid mode 99",
			mode:     99,
			position: 50,
			want:     nil,
			wantErr:  mechanismcontrol.ErrPositionCutMode,
		},
		{
			name:     "invalid mode 0",
			mode:     0,
			position: 50,
			want:     nil,
			wantErr:  mechanismcontrol.ErrPositionCutMode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cmd.SetCutPosition(tt.mode, tt.position)

			if !testutils.AssertErrorOccurred(t, err, tt.wantErr != nil, "SetCutPosition") {
				return
			}
			if tt.wantErr != nil {
				testutils.AssertError(t, err, tt.wantErr)
				return
			}

			testutils.AssertBytes(t, got, tt.want, "SetCutPosition(%v, %d)", tt.mode, tt.position)
		})
	}
}

func TestCommands_FeedCutAndReturnPaper(t *testing.T) {
	cmd := mechanismcontrol.NewCommands()
	prefix := testutils.BuildCommand(shared.GS, 'V')

	tests := []struct {
		name       string
		mode       mechanismcontrol.FeedCutReturn
		feedAmount byte
		want       []byte
		wantErr    error
	}{
		{
			name:       "full cut and return no feed",
			mode:       mechanismcontrol.FeedCutReturnFull,
			feedAmount: 0,
			want:       append(prefix, 103, 0),
			wantErr:    nil,
		},
		{
			name:       "partial cut and return no feed",
			mode:       mechanismcontrol.FeedCutReturnPartial,
			feedAmount: 0,
			want:       append(prefix, 104, 0),
			wantErr:    nil,
		},
		{
			name:       "full cut and return with feed",
			mode:       mechanismcontrol.FeedCutReturnFull,
			feedAmount: 25,
			want:       append(prefix, 103, 25),
			wantErr:    nil,
		},
		{
			name:       "partial cut and return with feed",
			mode:       mechanismcontrol.FeedCutReturnPartial,
			feedAmount: 75,
			want:       append(prefix, 104, 75),
			wantErr:    nil,
		},
		{
			name:       "full cut and return max feed",
			mode:       mechanismcontrol.FeedCutReturnFull,
			feedAmount: 255,
			want:       append(prefix, 103, 255),
			wantErr:    nil,
		},
		{
			name:       "invalid mode 102",
			mode:       102,
			feedAmount: 50,
			want:       nil,
			wantErr:    mechanismcontrol.ErrFeedCutReturnMode,
		},
		{
			name:       "invalid mode 105",
			mode:       105,
			feedAmount: 50,
			want:       nil,
			wantErr:    mechanismcontrol.ErrFeedCutReturnMode,
		},
		{
			name:       "invalid mode 0",
			mode:       0,
			feedAmount: 50,
			want:       nil,
			wantErr:    mechanismcontrol.ErrFeedCutReturnMode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cmd.FeedCutAndReturnPaper(tt.mode, tt.feedAmount)

			if !testutils.AssertErrorOccurred(t, err, tt.wantErr != nil, "FeedCutAndReturnPaper") {
				return
			}
			if tt.wantErr != nil {
				testutils.AssertError(t, err, tt.wantErr)
				return
			}

			testutils.AssertBytes(t, got, tt.want, "FeedCutAndReturnPaper(%v, %d)", tt.mode, tt.feedAmount)
		})
	}
}

// ============================================================================
// Validation Helper Functions Tests
// ============================================================================

func TestValidateCutMode(t *testing.T) {
	tests := []struct {
		name    string
		mode    mechanismcontrol.CutMode
		wantErr error
	}{
		{
			name:    "valid full cut (0)",
			mode:    mechanismcontrol.CutModeFull,
			wantErr: nil,
		},
		{
			name:    "valid partial cut (1)",
			mode:    mechanismcontrol.CutModePartial,
			wantErr: nil,
		},
		{
			name:    "valid full cut ASCII (48)",
			mode:    mechanismcontrol.CutModeFullASCII,
			wantErr: nil,
		},
		{
			name:    "valid partial cut ASCII (49)",
			mode:    mechanismcontrol.CutModePartialASCII,
			wantErr: nil,
		},
		{
			name:    "invalid mode 2",
			mode:    2,
			wantErr: mechanismcontrol.ErrCutMode,
		},
		{
			name:    "invalid mode 47",
			mode:    47,
			wantErr: mechanismcontrol.ErrCutMode,
		},
		{
			name:    "invalid mode 50",
			mode:    50,
			wantErr: mechanismcontrol.ErrCutMode,
		},
		{
			name:    "invalid mode 255",
			mode:    255,
			wantErr: mechanismcontrol.ErrCutMode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mechanismcontrol.ValidateCutMode(tt.mode)
			testutils.AssertError(t, err, tt.wantErr)
		})
	}
}

func TestValidateCutType(t *testing.T) {
	tests := []struct {
		name    string
		cutType mechanismcontrol.CutType
		wantErr error
	}{
		{
			name:    "valid full cut",
			cutType: mechanismcontrol.CutTypeFull,
			wantErr: nil,
		},
		{
			name:    "valid partial cut",
			cutType: mechanismcontrol.CutTypePartial,
			wantErr: nil,
		},
		{
			name:    "invalid cut type 2",
			cutType: 2,
			wantErr: mechanismcontrol.ErrCutType,
		},
		{
			name:    "invalid cut type 10",
			cutType: 10,
			wantErr: mechanismcontrol.ErrCutType,
		},
		{
			name:    "invalid cut type 255",
			cutType: 255,
			wantErr: mechanismcontrol.ErrCutType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mechanismcontrol.ValidateCutType(tt.cutType)
			testutils.AssertError(t, err, tt.wantErr)
		})
	}
}

func TestValidateFeedCutMode(t *testing.T) {
	tests := []struct {
		name    string
		mode    mechanismcontrol.FeedCut
		wantErr error
	}{
		{
			name:    "valid full cut mode",
			mode:    mechanismcontrol.FeedCutFull,
			wantErr: nil,
		},
		{
			name:    "valid partial cut mode",
			mode:    mechanismcontrol.FeedCutPartial,
			wantErr: nil,
		},
		{
			name:    "invalid mode 64",
			mode:    64,
			wantErr: mechanismcontrol.ErrFeedCutMode,
		},
		{
			name:    "invalid mode 67",
			mode:    67,
			wantErr: mechanismcontrol.ErrFeedCutMode,
		},
		{
			name:    "invalid mode 0",
			mode:    0,
			wantErr: mechanismcontrol.ErrFeedCutMode,
		},
		{
			name:    "invalid mode 255",
			mode:    255,
			wantErr: mechanismcontrol.ErrFeedCutMode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mechanismcontrol.ValidateFeedCutMode(tt.mode)
			testutils.AssertError(t, err, tt.wantErr)
		})
	}
}

func TestValidatePositionCutMode(t *testing.T) {
	tests := []struct {
		name    string
		mode    mechanismcontrol.PositionCut
		wantErr error
	}{
		{
			name:    "valid full cut mode",
			mode:    mechanismcontrol.PositionCutFull,
			wantErr: nil,
		},
		{
			name:    "valid partial cut mode",
			mode:    mechanismcontrol.PositionCutPartial,
			wantErr: nil,
		},
		{
			name:    "invalid mode 96",
			mode:    96,
			wantErr: mechanismcontrol.ErrPositionCutMode,
		},
		{
			name:    "invalid mode 99",
			mode:    99,
			wantErr: mechanismcontrol.ErrPositionCutMode,
		},
		{
			name:    "invalid mode 0",
			mode:    0,
			wantErr: mechanismcontrol.ErrPositionCutMode,
		},
		{
			name:    "invalid mode 255",
			mode:    255,
			wantErr: mechanismcontrol.ErrPositionCutMode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mechanismcontrol.ValidatePositionCutMode(tt.mode)
			testutils.AssertError(t, err, tt.wantErr)
		})
	}
}

func TestValidateFeedCutReturnMode(t *testing.T) {
	tests := []struct {
		name    string
		mode    mechanismcontrol.FeedCutReturn
		wantErr error
	}{
		{
			name:    "valid full cut mode",
			mode:    mechanismcontrol.FeedCutReturnFull,
			wantErr: nil,
		},
		{
			name:    "valid partial cut mode",
			mode:    mechanismcontrol.FeedCutReturnPartial,
			wantErr: nil,
		},
		{
			name:    "invalid mode 102",
			mode:    102,
			wantErr: mechanismcontrol.ErrFeedCutReturnMode,
		},
		{
			name:    "invalid mode 105",
			mode:    105,
			wantErr: mechanismcontrol.ErrFeedCutReturnMode,
		},
		{
			name:    "invalid mode 0",
			mode:    0,
			wantErr: mechanismcontrol.ErrFeedCutReturnMode,
		},
		{
			name:    "invalid mode 255",
			mode:    255,
			wantErr: mechanismcontrol.ErrFeedCutReturnMode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mechanismcontrol.ValidateFeedCutReturnMode(tt.mode)
			testutils.AssertError(t, err, tt.wantErr)
		})
	}
}

// ============================================================================
// Boundary and Edge Case Tests
// ============================================================================

func TestCommands_BoundaryValues(t *testing.T) {
	cmd := mechanismcontrol.NewCommands()

	t.Run("feed amount boundaries", func(t *testing.T) {
		// Test minimum feed amount
		got, err := cmd.PaperFeedAndCut(mechanismcontrol.CutTypeFull, 0)
		testutils.AssertError(t, err, nil)
		testutils.AssertLength(t, got, 8, "minimum feed amount command length")

		// Test maximum feed amount
		got, err = cmd.PaperFeedAndCut(mechanismcontrol.CutTypeFull, 255)
		testutils.AssertError(t, err, nil)
		testutils.AssertLength(t, got, 8, "maximum feed amount command length")
	})

	t.Run("position boundaries", func(t *testing.T) {
		// Test minimum position
		got, err := cmd.SetCutPosition(mechanismcontrol.PositionCutFull, 0)
		testutils.AssertError(t, err, nil)
		testutils.AssertLength(t, got, 4, "minimum position command length")

		// Test maximum position
		got, err = cmd.SetCutPosition(mechanismcontrol.PositionCutFull, 255)
		testutils.AssertError(t, err, nil)
		testutils.AssertLength(t, got, 4, "maximum position command length")
	})
}

func TestCommands_AllCutModeCombinations(t *testing.T) {
	cmd := mechanismcontrol.NewCommands()

	// Test all valid cut modes with CutPaper
	validModes := []mechanismcontrol.CutMode{
		mechanismcontrol.CutModeFull,
		mechanismcontrol.CutModePartial,
		mechanismcontrol.CutModeFullASCII,
		mechanismcontrol.CutModePartialASCII,
	}

	for _, mode := range validModes {
		t.Run(string(rune(mode)), func(t *testing.T) {
			got, err := cmd.CutPaper(mode)
			testutils.AssertError(t, err, nil)
			testutils.AssertNotEmpty(t, got, "CutPaper should return non-empty result for valid mode")
			testutils.AssertHasPrefix(t, got, []byte{shared.GS, 'V'}, "CutPaper command prefix")
		})
	}
}

func TestCommands_ConsistencyBetweenModes(t *testing.T) {
	cmd := mechanismcontrol.NewCommands()

	t.Run("full cut consistency", func(t *testing.T) {
		// CutPaper full cut modes should produce consistent results
		got1, _ := cmd.CutPaper(mechanismcontrol.CutModeFull)
		got2, _ := cmd.CutPaper(mechanismcontrol.CutModeFullASCII)

		// Both should be GS V commands but with different parameter values
		testutils.AssertHasPrefix(t, got1, []byte{shared.GS, 'V'}, "full cut (0)")
		testutils.AssertHasPrefix(t, got2, []byte{shared.GS, 'V'}, "full cut ASCII (48)")

		// Values should be different
		if got1[2] == got2[2] {
			t.Error("CutTypeFull cut modes should have different parameter values")
		}
	})

	t.Run("partial cut consistency", func(t *testing.T) {
		// CutPaper partial cut modes should produce consistent results
		got1, _ := cmd.CutPaper(mechanismcontrol.CutModePartial)
		got2, _ := cmd.CutPaper(mechanismcontrol.CutModePartialASCII)

		// Both should be GS V commands but with different parameter values
		testutils.AssertHasPrefix(t, got1, []byte{shared.GS, 'V'}, "partial cut (1)")
		testutils.AssertHasPrefix(t, got2, []byte{shared.GS, 'V'}, "partial cut ASCII (49)")

		// Values should be different
		if got1[2] == got2[2] {
			t.Error("CutTypePartial cut modes should have different parameter values")
		}
	})
}
