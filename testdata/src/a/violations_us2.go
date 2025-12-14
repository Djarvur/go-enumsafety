package a

// Priority enum
type Priority uint8 // want "quasi-enum type Priority lacks a String\\(\\) method; consider using golang.org/x/tools/cmd/stringer or github.com/Djarvur/go-silly-enum to generate it" "quasi-enum type Priority lacks an UnmarshalText\\(\\[\\]byte\\) error method; consider using github.com/Djarvur/go-silly-enum to generate it"

const (
	PriorityLow  Priority = 1
	PriorityHigh Priority = 2
)

// Test US2: Untyped Constant Assignment Detection

// Untyped constants (not part of the enum)
const (
	myValue    = 3
	anotherVal = 5
)

func testUntypedConstantAssignment() {
	// Valid: defined enum constant
	var p1 Priority = PriorityLow

	// Invalid: untyped constant (not part of enum)
	var p2 Priority = myValue // want "untyped constant assigned to quasi-enum type Priority"

	// Invalid: assignment with untyped constant
	p3 := Priority(0) // want "literal value converted to quasi-enum type Priority"
	p3 = anotherVal   // want "untyped constant assigned to quasi-enum type Priority"

	_ = p1
	_ = p2
	_ = p3
}

func SetPriority(p Priority) {
	_ = p
}

func testUntypedConstantArgument() {
	// Valid: defined enum constant
	SetPriority(PriorityHigh)

	// Invalid: untyped constant as argument
	SetPriority(myValue) // want "untyped constant assigned to quasi-enum type Priority"
}
