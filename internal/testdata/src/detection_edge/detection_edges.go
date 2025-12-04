package detection_edge

// Test detection technique edge cases

// DT-004: Multiple comment lines with enum not at start
// This is a long comment
// that spans multiple lines
// enum
// and has the keyword in the middle
type MultiLineComment int // want "quasi-enum type MultiLineComment uses int but has only 2 constants; consider using uint8 for memory optimization" "quasi-enum type MultiLineComment lacks a String\\(\\) method" "quasi-enum type MultiLineComment lacks an UnmarshalText\\(\\[\\]byte\\) error method"

const (
	MultiLineCommentFirst  MultiLineComment = 1
	MultiLineCommentSecond MultiLineComment = 2
)

// DT-004: Comment with enum keyword but not at the beginning
// This type is an enum for testing
type NotAtStart int // want "quasi-enum type NotAtStart uses int but has only 2 constants; consider using uint8 for memory optimization" "quasi-enum type NotAtStart lacks a String\\(\\) method" "quasi-enum type NotAtStart lacks an UnmarshalText\\(\\[\\]byte\\) error method"

const (
	NotAtStartFirst  NotAtStart = 1
	NotAtStartSecond NotAtStart = 2
)

// Should be detected: enum at the very start
// enum - this is a valid enum
type ValidPreceding int // want "quasi-enum type ValidPreceding uses int but has only 2 constants; consider using uint8 for memory optimization" "quasi-enum type ValidPreceding lacks a String\\(\\) method" "quasi-enum type ValidPreceding lacks an UnmarshalText\\(\\[\\]byte\\) error method"

const (
	ValidPrecedingFirst  ValidPreceding = 1
	ValidPrecedingSecond ValidPreceding = 2
)

// Test case for type conversion with selector expression
type Package struct {
	Status int
}

// enum
type SelectorTest int // want "quasi-enum type SelectorTest uses int but has only 2 constants; consider using uint8 for memory optimization" "quasi-enum type SelectorTest lacks a String\\(\\) method" "quasi-enum type SelectorTest lacks an UnmarshalText\\(\\[\\]byte\\) error method"

const (
	SelectorTestFirst  SelectorTest = 1
	SelectorTestSecond SelectorTest = 2
)

func testSelectorConversion() {
	pkg := Package{Status: 1}
	// Note: Selector expression conversions are not currently detected
	s := SelectorTest(pkg.Status)
	_ = s
}
