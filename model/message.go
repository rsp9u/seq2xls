package model

// MessageType is a type of the kind of message.
type MessageType int

const (
	// Synchronous is the synchronous message type.
	Synchronous MessageType = iota
	// Asynchronous is the asynchronous message type.
	Asynchronous
	// Reply is the reply message type.
	Reply
	// Found is the found message type.
	Found
	// Lost is the lost message type.
	Lost
	// SelfReference is the self-reference message type.
	SelfReference
)

// Message is a data model of the message.
type Message struct {
	Index    int
	From     *Lifeline
	To       *Lifeline
	Type     MessageType
	ColorHex string
	Text     string
}
