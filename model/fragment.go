package model

// FragmentType is a type of the kind of fragment.
type FragmentType int

const (
	// Ref is the reference fragment type.
	Ref FragmentType = iota
	// Alt is the alternative fragment type.
	Alt
	// Opt is the options fragment type.
	Opt
	// Par is the parallel fragment type.
	Par
	// Loop is the loop fragment type.
	Loop
	// Break is the break fragment type.
	Break
	// Critical is the critial fragment type.
	Critical
	// Assert is the assert fragment type.
	Assert
	// Neg is the negative fragment type.
	Neg
	// Ignore is the ignore fragment type.
	Ignore
	// Consider is the consider fragment type.
	Consider
	// UnknownFragment is the unknown type of fragment.
	UnknownFragment
)

// Fragment is a data model of the fragment.
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
