package a

// Status enum
type Status int // want "quasi-enum type Status uses int but has only 3 constants; consider using uint8 for memory optimization" "quasi-enum type Status lacks a String\\(\\) method; consider using golang.org/x/tools/cmd/stringer or github.com/Djarvur/go-silly-enum to generate it" "quasi-enum type Status lacks an UnmarshalText\\(\\[\\]byte\\) error method; consider using github.com/Djarvur/go-silly-enum to generate it"

const (
	StatusActive Status = iota
	StatusInactive
	StatusPending
)

// Test US1: Literal Assignment Detection

func testLiteralAssignment() {
	// Valid: constant assignment
	var s1 Status = StatusActive

	// Invalid: literal assignment
	var s2 Status = 5 // want "literal value assigned to quasi-enum type Status"

	// Invalid: literal conversion
	s3 := Status(3) // want "literal value converted to quasi-enum type Status"

	_ = s1
	_ = s2
	_ = s3
}

func SetStatus(s Status) {
	_ = s
}

func testLiteralArgument() {
	// Valid: constant argument
	SetStatus(StatusActive)

	// Invalid: literal argument
	SetStatus(2) // want "literal value passed as quasi-enum type Status"
}

type Config struct {
	Status Status
}

func testCompositeLiteral() {
	// Valid: constant in composite literal
	c1 := Config{Status: StatusActive}

	// Invalid: literal in composite literal
	c2 := Config{Status: 1} // want "literal value in composite literal for quasi-enum type Status"

	_ = c1
	_ = c2
}
