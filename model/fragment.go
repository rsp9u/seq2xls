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
)

type Fragment struct {
	Begin, End *Message
	Contain    []*Lifeline
	Type       FragmentType
}
