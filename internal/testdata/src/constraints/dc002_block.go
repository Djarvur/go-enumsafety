package constraints

// Test DC-002: Same Const Block Check
// All enum constants must be in the same const block

// Valid: All constants in same block
type StatusValid int

const (
	StatusValidActive StatusValid = iota
	StatusValidInactive
	StatusValidPending
)

// Invalid: Constants in different blocks (should warn)
// enum
type PriorityInvalid uint8

const PriorityInvalidLow PriorityInvalid = 1

const (
	PriorityInvalidMedium PriorityInvalid = 2
	PriorityInvalidHigh   PriorityInvalid = 3
)
