package constants

type Message int

const (
	MessageText Message = iota
	MessageExtendedText
	MessageImage
	MessageVideo
	MessageDocument
)
