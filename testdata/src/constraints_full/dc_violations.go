package constraints_full

// Test comprehensive constraint violations
// Focus on covering validation code paths

// Single constant enum
// enum
type SingleConstEnum int // want "quasi-enum type SingleConstEnum violates DC-001 \\(minimum 2 constants\\): must have at least 2 constants" "quasi-enum type SingleConstEnum uses int but has only 1 constants; consider using uint8 for memory optimization" "quasi-enum type SingleConstEnum lacks a String\\(\\) method" "quasi-enum type SingleConstEnum lacks an UnmarshalText\\(\\[\\]byte\\) error method"

const SingleConstEnumValue SingleConstEnum = 1

// Constants in different blocks
// enum
type SplitBlockEnum int // want "quasi-enum type SplitBlockEnum uses int but has only 2 constants; consider using uint8 for memory optimization" "quasi-enum type SplitBlockEnum lacks a String\\(\\) method" "quasi-enum type SplitBlockEnum lacks an UnmarshalText\\(\\[\\]byte\\) error method"

const SplitBlockEnumFirst SplitBlockEnum = 1

const SplitBlockEnumSecond SplitBlockEnum = 2

// Mixed constant block
// enum
type MixedBlockEnum int // want "quasi-enum type MixedBlockEnum violates DC-004 \\(exclusive const block\\): const block must contain only constants of this type" "quasi-enum type MixedBlockEnum uses int but has only 2 constants; consider using uint8 for memory optimization" "quasi-enum type MixedBlockEnum lacks a String\\(\\) method" "quasi-enum type MixedBlockEnum lacks an UnmarshalText\\(\\[\\]byte\\) error method"

const (
	MixedBlockEnumFirst  MixedBlockEnum = 1
	MixedBlockEnumSecond MixedBlockEnum = 2
	UnrelatedConstant    int            = 99
)

// Type and constants far apart
// enum
type FarApartEnum int // want "quasi-enum type FarApartEnum violates DC-005 \\(proximity\\): type definition and const block must be adjacent" "quasi-enum type FarApartEnum uses int but has only 2 constants; consider using uint8 for memory optimization" "quasi-enum type FarApartEnum lacks a String\\(\\) method" "quasi-enum type FarApartEnum lacks an UnmarshalText\\(\\[\\]byte\\) error method"

var spacer1 int
var spacer2 int
var spacer3 int
var spacer4 int
var spacer5 int
var spacer6 int
var spacer7 int
var spacer8 int
var spacer9 int
var spacer10 int

const (
	FarApartEnumFirst  FarApartEnum = 1
	FarApartEnumSecond FarApartEnum = 2
)

// Valid enum with uint8
// enum
type ValidEnum uint8 // want "quasi-enum type ValidEnum lacks a String\\(\\) method" "quasi-enum type ValidEnum lacks an UnmarshalText\\(\\[\\]byte\\) error method"

const (
	ValidEnumFirst  ValidEnum = 1
	ValidEnumSecond ValidEnum = 2
)
