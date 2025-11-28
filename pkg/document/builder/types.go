package builder

// Alignment constants for text and table positioning
type Alignment string

const (
	// Left as alignment
	Left Alignment = "left"
	// Center as alignment
	Center Alignment = "center"
	// Right as alignment
	Right Alignment = "right"
)

// String returns the string representation of the alignment
func (a Alignment) String() string {
	return string(a)
}
