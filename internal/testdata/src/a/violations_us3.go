package a

// Color enum
type Color uint8 // want "quasi-enum type Color lacks a String\\(\\) method; consider using golang.org/x/tools/cmd/stringer or github.com/Djarvur/go-silly-enum to generate it" "quasi-enum type Color lacks an UnmarshalText\\(\\[\\]byte\\) error method; consider using github.com/Djarvur/go-silly-enum to generate it"

const (
	ColorRed Color = iota
	ColorGreen
	ColorBlue
)

// Test US3: Variable Conversion Detection

func testVariableConversion() {
	// Valid: enum constant
	var c1 Color = ColorRed

	// Invalid: variable of underlying type converted to enum
	var x uint8 = 5
	c2 := Color(x) // want "variable converted to quasi-enum type Color"

	// Invalid: variable conversion in assignment
	var y uint8 = 10
	var c3 Color
	c3 = Color(y) // want "variable converted to quasi-enum type Color"

	_ = c1
	_ = c2
	_ = c3
}

// Test cross-enum conversion
type Level uint8 // want "quasi-enum type Level lacks a String\\(\\) method; consider using golang.org/x/tools/cmd/stringer or github.com/Djarvur/go-silly-enum to generate it" "quasi-enum type Level lacks an UnmarshalText\\(\\[\\]byte\\) error method; consider using github.com/Djarvur/go-silly-enum to generate it"

const (
	LevelLow  Level = 1
	LevelHigh Level = 2
)

func testCrossEnumConversion() {
	// Invalid: converting from one enum type to another
	var c Color = ColorRed
	l := Level(c) // want "variable converted to quasi-enum type Level"

	_ = l
}

func SetColor(c Color) {
	_ = c
}

func testVariableInFunctionCall() {
	// This should NOT be flagged - function calls don't use type conversions
	// The variable itself is passed, not converted
	var c Color = ColorRed
	SetColor(c) // Valid - passing enum variable directly
}
