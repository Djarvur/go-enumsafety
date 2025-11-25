package main

// Test DC-003: Same File Check
// Type and constants must be in the same file

// Valid: Type and constants in same file
type StatusValid int

const (
	StatusValidActive StatusValid = iota
	StatusValidInactive
)

// Invalid: Type here, constants in other file (constraints_other/other.go)
// This test requires the type to be detected but constants elsewhere
// For simplicity, we'll document this as a cross-file test scenario
// The actual test would need multiple files in the package

// Note: This is a documentation placeholder for the constraint
// Real testing would require analyzing multiple files in a package
