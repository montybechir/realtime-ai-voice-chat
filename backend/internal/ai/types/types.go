package types

type MessageType int

const (
	TextMessage MessageType = iota
	AudioMessage
	SystemMessage
)

type Message struct {
	SenderID   string
	ReceiverID string
	Payload    []byte
	Type       MessageType
}

// AIClient represents the AI client connection.
