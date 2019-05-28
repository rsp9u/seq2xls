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

func (t FragmentType) String() string {
	switch t {
	case Ref:
		return "ref"
	case Alt:
		return "alt"
	case Opt:
		return "opt"
	case Par:
		return "par"
	case Loop:
		return "loop"
	case Break:
		return "break"
	case Critical:
		return "critical"
	case Assert:
		return "assert"
	case Neg:
		return "neg"
	case Ignore:
		return "ignore"
	case Consider:
		return "consider"
	}
	return "unknown"
}
