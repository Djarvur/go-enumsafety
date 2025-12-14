package constraints

// Test DC-004: Exclusive Const Block Check
// Const block should only contain constants of the enum type

// Valid: Block contains only Status constants
type StatusValid int

const (
	StatusValidActive StatusValid = iota
	StatusValidInactive
	StatusValidPending
)

// Invalid: Block contains mixed types (should warn)
// enum
type PriorityInvalid uint8

const (
	PriorityInvalidLow  PriorityInvalid = 1
	PriorityInvalidHigh PriorityInvalid = 2
	SomeOtherConst                      = "mixed" // Different type in same block
)
