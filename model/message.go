package model

type MessageType int

const (
	Synchronous MessageType = iota
	Asynchronous
	Reply
	Found
	Lost
	SelfReference
)

type Message struct {
	Index    int
	From     *Lifeline
	To       *Lifeline
	Type     MessageType
	ColorHex string
}
