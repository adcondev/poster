package calculate

import "testing"

func TestDotsPerLine(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		paperWidthMM float64
		dpi          int
		want         int
	}{
		{
			name:         "80mm paper at 203 dpi",
			paperWidthMM: 80,
			dpi:          203,
			want:         639, // (80 * 203) / 25.4 = 639.37... -> 639
		},
		{
			name:         "58mm paper at 203 dpi",
			paperWidthMM: 58,
			dpi:          203,
			want:         463, // (58 * 203) / 25.4 = 463.54... -> 463
		},
		{
			name:         "80mm paper at 180 dpi",
			paperWidthMM: 80,
			dpi:          180,
			want:         566, // (80 * 180) / 25.4 = 566.92... -> 566
		},
		{
			name:         "Exact inch (25.4mm) at 100 dpi",
			paperWidthMM: 25.4,
			dpi:          100,
			want:         100,
		},
		{
			name:         "Zero width",
			paperWidthMM: 0,
			dpi:          203,
			want:         0,
		},
		{
			name:         "Zero DPI",
			paperWidthMM: 80,
			dpi:          0,
			want:         0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := DotsPerLine(tt.paperWidthMM, tt.dpi); got != tt.want {
				t.Errorf("DotsPerLine() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMmToDots(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		mm   float64
		dpi  int
		want int
	}{
		{
			name: "10mm at 203 dpi",
			mm:   10,
			dpi:  203,
			want: 79, // (10 * 203) / 25.4 = 79.92... -> 79
		},
		{
			name: "1 inch (25.4mm) at 203 dpi",
			mm:   25.4,
			dpi:  203,
			want: 203,
		},
		{
			name: "Zero mm",
			mm:   0,
			dpi:  203,
			want: 0,
		},
		{
			name: "Negative mm",
			mm:   -10,
			dpi:  203,
			want: -79,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := MmToDots(tt.mm, tt.dpi); got != tt.want {
				t.Errorf("MmToDots() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDotsToMm(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		dots int
		dpi  int
		want float64
	}{
		{
			name: "203 dots at 203 dpi (1 inch)",
			dots: 203,
			dpi:  203,
			want: 25.4,
		},
		{
			name: "100 dots at 203 dpi",
			dots: 100,
			dpi:  203,
			want: 12.51231527093596, // (100 * 25.4) / 203
		},
		{
			name: "Zero dots",
			dots: 0,
			dpi:  203,
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := DotsToMm(tt.dots, tt.dpi)
			// Using a small epsilon for float comparison
			if diff := got - tt.want; diff < -0.00001 || diff > 0.00001 {
				t.Errorf("DotsToMm() = %v, want %v", got, tt.want)
			}
		})
	}
}
