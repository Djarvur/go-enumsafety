package detection

// Test DT-004: Preceding Comment Detection
// Types with "// enum" in doc comment should be detected

// enum
// Status represents the status of an entity
type Status int

const (
	StatusActive Status = iota
	StatusInactive
)

// Priority levels
// enum
type Priority uint8

const (
	PriorityLow  Priority = 1
	PriorityHigh Priority = 2
)

// Should NOT be detected (no "enum" in comment)
// Color represents a color value
type Color uint8

const (
	ColorRed Color = iota
	ColorGreen
)

// Should be detected (enum at start of multi-line comment)
// enum
// Level represents difficulty levels
type Level int

const (
	LevelEasy Level = 1
	LevelHard Level = 2
)
