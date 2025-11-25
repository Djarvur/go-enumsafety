package detection

// Test DT-001: Constants-Based Detection
// Types with 2+ constants should be detected as quasi-enums

// Should be detected (2 constants)
type Status int

const (
	StatusActive Status = iota
	StatusInactive
)

// Should be detected (3 constants)
type Priority uint8

const (
	PriorityLow    Priority = 1
	PriorityMedium Priority = 2
	PriorityHigh   Priority = 3
)

// Should NOT be detected (only 1 constant)
type SingleValue int

const OnlyOne SingleValue = 1

// Should be detected (many constants)
type Color uint8

const (
	ColorRed Color = iota
	ColorGreen
	ColorBlue
	ColorYellow
	ColorOrange
)
