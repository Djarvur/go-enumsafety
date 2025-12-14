package optimization

// Test US4: uint8 Optimization Suggestion

// Should suggest uint8 (int with 3 constants)
type StatusInt int // want "quasi-enum type StatusInt uses int but has only 3 constants; consider using uint8 for memory optimization" "quasi-enum type StatusInt lacks a String\\(\\) method; consider using golang.org/x/tools/cmd/stringer or github.com/Djarvur/go-silly-enum to generate it" "quasi-enum type StatusInt lacks an UnmarshalText\\(\\[\\]byte\\) error method; consider using github.com/Djarvur/go-silly-enum to generate it"

const (
	StatusActive StatusInt = iota
	StatusInactive
	StatusPending
)

// Should NOT suggest (already uint8)
type ColorUint8 uint8 // want "quasi-enum type ColorUint8 lacks a String\\(\\) method; consider using golang.org/x/tools/cmd/stringer or github.com/Djarvur/go-silly-enum to generate it" "quasi-enum type ColorUint8 lacks an UnmarshalText\\(\\[\\]byte\\) error method; consider using github.com/Djarvur/go-silly-enum to generate it"

const (
	ColorRed ColorUint8 = iota
	ColorGreen
	ColorBlue
)

// Should suggest uint8 (uint16 with only 3 constants)
type LargeEnum uint16 // want "quasi-enum type LargeEnum uses uint16 but has only 3 constants; consider using uint8 for memory optimization" "quasi-enum type LargeEnum lacks a String\\(\\) method; consider using golang.org/x/tools/cmd/stringer or github.com/Djarvur/go-silly-enum to generate it" "quasi-enum type LargeEnum lacks an UnmarshalText\\(\\[\\]byte\\) error method; consider using github.com/Djarvur/go-silly-enum to generate it"

const (
	Large0 LargeEnum = iota
	Large1
	Large2
)

// Should suggest uint8 (uint with 2 constants)
type PriorityUint uint // want "quasi-enum type PriorityUint uses uint but has only 2 constants; consider using uint8 for memory optimization" "quasi-enum type PriorityUint lacks a String\\(\\) method; consider using golang.org/x/tools/cmd/stringer or github.com/Djarvur/go-silly-enum to generate it" "quasi-enum type PriorityUint lacks an UnmarshalText\\(\\[\\]byte\\) error method; consider using github.com/Djarvur/go-silly-enum to generate it"

const (
	PriorityLow  PriorityUint = 1
	PriorityHigh PriorityUint = 2
)

// Should suggest uint8 (int32 with 4 constants)
type Level int32 // want "quasi-enum type Level uses int32 but has only 4 constants; consider using uint8 for memory optimization" "quasi-enum type Level lacks a String\\(\\) method; consider using golang.org/x/tools/cmd/stringer or github.com/Djarvur/go-silly-enum to generate it" "quasi-enum type Level lacks an UnmarshalText\\(\\[\\]byte\\) error method; consider using github.com/Djarvur/go-silly-enum to generate it"

const (
	Level1 Level = iota
	Level2
	Level3
	Level4
)
