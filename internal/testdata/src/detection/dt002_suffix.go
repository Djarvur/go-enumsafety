package detection

// Test DT-002: Name Suffix Detection
// Types ending in "enum" (case-insensitive) should be detected

// Should be detected (lowercase "enum")
type Statusenum int

const (
	StatusenumActive Statusenum = iota
	StatusenumInactive
)

// Should be detected (uppercase "ENUM")
type PriorityENUM uint8

const (
	PriorityENUMLow  PriorityENUM = 1
	PriorityENUMHigh PriorityENUM = 2
)

// Should be detected (mixed case "Enum")
type ColorEnum uint8

const (
	ColorEnumRed ColorEnum = iota
	ColorEnumGreen
)

// Should NOT be detected (no "enum" suffix)
type Level int

const (
	LevelLow  Level = 1
	LevelHigh Level = 2
)
