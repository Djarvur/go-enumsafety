package helpers

// Test US5 & US6: String() and UnmarshalText() Method Checks

// Missing both methods - should warn twice
type Priority uint8 // want "quasi-enum type Priority lacks a String\\(\\) method; consider using golang.org/x/tools/cmd/stringer or github.com/Djarvur/go-silly-enum to generate it" "quasi-enum type Priority lacks an UnmarshalText\\(\\[\\]byte\\) error method; consider using github.com/Djarvur/go-silly-enum to generate it"

const (
	PriorityLow Priority = iota
	PriorityHigh
)

// Has String() - should warn for UnmarshalText and uint8 optimization
type Status int // want "quasi-enum type Status uses int but has only 3 constants; consider using uint8 for memory optimization" "quasi-enum type Status lacks an UnmarshalText\\(\\[\\]byte\\) error method; consider using github.com/Djarvur/go-silly-enum to generate it"

const (
	StatusActive Status = iota
	StatusInactive
	StatusPending
)

func (s Status) String() string {
	switch s {
	case StatusActive:
		return "Active"
	case StatusInactive:
		return "Inactive"
	case StatusPending:
		return "Pending"
	default:
		return "Unknown"
	}
}

// Has both - should NOT warn for methods
type Level uint8

const (
	LevelLow Level = iota
	LevelHigh
)

func (l Level) String() string {
	switch l {
	case LevelLow:
		return "Low"
	case LevelHigh:
		return "High"
	default:
		return "Unknown"
	}
}

func (l *Level) UnmarshalText(text []byte) error {
	switch string(text) {
	case "Low":
		*l = LevelLow
	case "High":
		*l = LevelHigh
	default:
		return nil
	}
	return nil
}

// Has UnmarshalText but not String - should warn for String
type Color uint8 // want "quasi-enum type Color lacks a String\\(\\) method; consider using golang.org/x/tools/cmd/stringer or github.com/Djarvur/go-silly-enum to generate it"

const (
	ColorRed Color = iota
	ColorGreen
	ColorBlue
)

func (c *Color) UnmarshalText(text []byte) error {
	switch string(text) {
	case "Red":
		*c = ColorRed
	case "Green":
		*c = ColorGreen
	case "Blue":
		*c = ColorBlue
	}
	return nil
}
