package call_expr

// Test function call expression edge cases

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

func SetStatus(s Status) {
	_ = s
}

func SetMultiple(s Status, p Priority, name string) {
	_ = s
	_ = p
	_ = name
}

func SetMany(statuses ...Status) {
	_ = statuses
}

type Handler struct{}

func (h Handler) Handle(s Status) {
	_ = s
}

func testFunctionCalls() {
	// Valid
	SetStatus(StatusActive)

	// Literal argument
	SetStatus(5) // want "literal value passed as quasi-enum type Status"

	// Multiple parameters
	SetMultiple(StatusActive, 99, "test") // want "literal value passed as quasi-enum type Priority"

	// Variadic function - Note: variadic parameters not checked
	SetMany(StatusActive, 2, StatusPending)

	// Method call
	h := Handler{}
	h.Handle(7) // want "literal value passed as quasi-enum type Status"
}

func testEdgeCases() {
	SetStatus(StatusActive)

	const untypedConst = 5
	SetStatus(untypedConst) // want "untyped constant assigned to quasi-enum type Status"
}

func testTypeConversionInCall() {
	var x int = 3
	// Note: Type conversions in function calls wrap the conversion
	// These conversions are valid enum constants after the conversion
	SetStatus(Status(x))
	SetStatus(Status(5)) // want "literal value passed as quasi-enum type Status"
}
