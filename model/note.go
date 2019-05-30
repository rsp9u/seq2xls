package model

// Note is a data model of the note.
type Note struct {
	Assoc    *Message
	OnLeft   bool
	Text     string
	ColorHex string
}
