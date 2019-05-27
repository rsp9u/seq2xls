package model

type FragmentType int

const (
	Ref FragmentType = iota
	Alt
	Opt
	Par
	Loop
	Break
	Critical
	Assert
	Neg
	Ignore
	Consider
	UnknownFragment
)

type Fragment struct {
	Index      int
	Begin, End *Message
	Type       FragmentType
}
