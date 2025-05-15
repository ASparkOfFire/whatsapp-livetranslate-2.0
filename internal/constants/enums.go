package constants

//go:generate stringer -type Message -trimprefix Message
type Message int

const (
	MessageText Message = iota
	MessageExtendedText
	MessageImage
	MessageVideo
	MessageDocument
	MessageAudio
	MessageLocation
	MessageContact
	MessagePoll
	MessageUnknown
)
