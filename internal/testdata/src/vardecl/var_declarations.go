package vardecl

// Test variable declaration edge cases

// Status enum
type Status int // want "quasi-enum type Status uses int but has only 3 constants; consider using uint8 for memory optimization" "quasi-enum type Status lacks a String\\(\\) method" "quasi-enum type Status lacks an UnmarshalText\\(\\[\\]byte\\) error method"

const (
	StatusActive Status = iota
	StatusInactive
	StatusPending
)

// Priority enum
type Priority int // want "quasi-enum type Priority uses int but has only 2 constants; consider using uint8 for memory optimization" "quasi-enum type Priority lacks a String\\(\\) method" "quasi-enum type Priority lacks an UnmarshalText\\(\\[\\]byte\\) error method"

const (
	PriorityLow  Priority = 1
	PriorityHigh Priority = 2
)

func testVarDeclEdgeCases() {
	// Variable with no initial value
	var s1 Status
	_ = s1

	// Multiple variables with literals
	var s2, s3 Status = StatusActive, 5 // want "literal value assigned to quasi-enum type Status"
	_ = s2
	_ = s3

	// Cross-enum conversion
	var p Priority = PriorityLow
	var s4 Status = Status(p) // want "variable converted to quasi-enum type Status"
	_ = s4

	// Underlying type conversion
	var x int = 3
	var s5 Status = Status(x) // want "variable converted to quasi-enum type Status"
	_ = s5

	// Valid: conversion from constant
	const validConst Status = StatusActive
	var s6 Status = validConst // want "untyped constant assigned to quasi-enum type Status"
	_ = s6

	// Invalid: conversion from untyped constant
	const untypedConst = 7
	var s7 Status = untypedConst // want "untyped constant assigned to quasi-enum type Status"
	_ = s7
}

func testVarDeclWithConversion() {
	var underlying int = 5
	s := Status(underlying) // want "variable converted to quasi-enum type Status"
	_ = s
}
