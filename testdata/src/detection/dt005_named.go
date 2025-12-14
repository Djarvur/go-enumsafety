package detection

// Test DT-005: Named Comment Detection
// Types with "// TypeName enum" pattern should be detected

// Status enum
type Status int

const (
	StatusActive Status = iota
	StatusInactive
)

// Priority enum - priority levels
type Priority uint8

const (
	PriorityLow  Priority = 1
	PriorityHigh Priority = 2
)

// Should NOT be detected (wrong type name in comment)
// WrongName enum
type Color uint8

const (
	ColorRed Color = iota
	ColorGreen
)

// Should be detected (exact match)
// Level enum
type Level int

const (
	LevelLow  Level = 1
	LevelHigh Level = 2
)
