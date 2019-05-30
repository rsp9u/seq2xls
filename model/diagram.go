package model

// SequenceDiagram is a data model of the sequence diagram.
type SequenceDiagram struct {
	Lifelines  []*Lifeline
	ExecSpecs  []*ExecSpec
	Messages   []*Message
	Fragments  []*Fragment
	Notes      []*Note
	Separators []*Separator
}
