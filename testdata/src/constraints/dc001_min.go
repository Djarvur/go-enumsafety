package constraints

// Test DC-001: Minimum Constants Check
// Quasi-enums must have at least 2 constants

// Valid: 2 constants
type StatusValid int

const (
	StatusValidActive StatusValid = iota
	StatusValidInactive
)

// Valid: 3 constants
type PriorityValid uint8

const (
	PriorityValidLow    PriorityValid = 1
	PriorityValidMedium PriorityValid = 2
	PriorityValidHigh   PriorityValid = 3
)

// Invalid: Only 1 constant (should warn about constraint violation)
// Note: This won't be detected as quasi-enum by constants-based detection
// but could be detected by other techniques
// enum
type ColorInvalid uint8

const ColorInvalidRed ColorInvalid = 1
