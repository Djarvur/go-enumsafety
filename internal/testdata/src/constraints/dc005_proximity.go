package constraints

// Test DC-005: Proximity Check
// Type declaration and const block should be adjacent (no code between them)

// Valid: Type and const block are adjacent
type StatusValid int

const (
	StatusValidActive StatusValid = iota
	StatusValidInactive
)

// Valid: Empty lines and comments are allowed
type PriorityValid uint8

// This is a comment

const (
	PriorityValidLow  PriorityValid = 1
	PriorityValidHigh PriorityValid = 2
)

// Invalid: Code between type and const block (should warn)
// enum
type ColorInvalid uint8

var someVariable = 42 // Code between type and const

const (
	ColorInvalidRed ColorInvalid = iota
	ColorInvalidGreen
)
