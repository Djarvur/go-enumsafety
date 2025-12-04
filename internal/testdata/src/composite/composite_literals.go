package composite

// Test composite literal edge cases

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

// Simple config
type Config struct {
	Status   Status
	Priority Priority
	Name     string
}

// Nested config
type NestedConfig struct {
	Inner Config
	ID    int
}

func testCompositeLiteralEdgeCases() {
	// Multiple enum fields with mixed valid/invalid
	c1 := Config{
		Status:   StatusActive,
		Priority: 99, // want "literal value in composite literal for quasi-enum type Priority"
		Name:     "test",
	}
	_ = c1

	// Nested struct with enum field
	c2 := NestedConfig{
		Inner: Config{
			Status:   5, // want "literal value in composite literal for quasi-enum type Status"
			Priority: PriorityLow,
			Name:     "nested",
		},
		ID: 1,
	}
	_ = c2

	// Valid: all enum constants
	c3 := Config{
		Status:   StatusActive,
		Priority: PriorityHigh,
		Name:     "valid",
	}
	_ = c3
}

// Test struct with only enum fields
type EnumOnly struct {
	S Status
	P Priority
}

func testEnumOnlyStruct() {
	// Both fields with literals
	e1 := EnumOnly{
		S: 1, // want "literal value in composite literal for quasi-enum type Status"
		P: 2, // want "literal value in composite literal for quasi-enum type Priority"
	}
	_ = e1

	// Valid
	e2 := EnumOnly{
		S: StatusActive,
		P: PriorityLow,
	}
	_ = e2
}
