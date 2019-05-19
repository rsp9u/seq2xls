package model

type SequenceDiagram struct {
	Lifelines []*Lifeline
	ExecSpecs []*ExecSpec
	Messages  []*Message
	Fragments []*Fragment
}
