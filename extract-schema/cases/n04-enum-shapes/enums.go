package app

type Color string

const (
	ColorRed   Color = "red"
	ColorGreen Color = "green"
	ColorBlue  Color = "blue"
)

type Perm int

const (
	PermRead Perm = 1 << iota
	PermWrite
	PermExec
	_ // blank member
)

// Solo has exactly one member.
type Solo string

const Only Solo = "only"

// NoMembers is a named type with no const block.
type NoMembers int

// Mixed picks up only the members typed Mixed — Free stays out.
type Mixed int

const (
	M1 Mixed = 1
	M2 Mixed = 2
	Free     = 3
)

// Split declares members across two const blocks.
type Split string

const S1 Split = "s1"

const S2 Split = "s2"
