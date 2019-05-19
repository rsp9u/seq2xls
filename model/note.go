package model

type Note struct {
	Assoc    *Message
	OnLeft   bool
	Text     string
	ColorHex string
}
