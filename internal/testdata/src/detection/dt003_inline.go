package detection

// Test DT-003: Inline Comment Detection
// Types with inline "// enum" comment should be detected

// Should be detected (inline comment)
type Status int // enum

const (
	StatusActive Status = iota
	StatusInactive
)

// Should be detected (inline comment with extra text)
type Priority uint8 // enum - priority levels

const (
	PriorityLow  Priority = 1
	PriorityHigh Priority = 2
)

// Should NOT be detected (no inline comment)
type Color uint8

const (
	ColorRed Color = iota
	ColorGreen
)

// Should NOT be detected (comment but not "enum")
type Level int // levels

const (
	LevelLow  Level = 1
	LevelHigh Level = 2
)
