package messagehandler

import (
	"fmt"

	"github.com/asparkoffire/whatsapp-livetranslate-go/internal/constants"
	waProto "go.mau.fi/whatsmeow/proto/waE2E"
)

// getMessageType returns the message type based on the waProto.Message structure
func getMessageType(msg *waProto.Message) constants.Message {
	switch {
	case msg.GetConversation() != "":
		return constants.MessageText
	case msg.GetExtendedTextMessage() != nil:
		return constants.MessageExtendedText
	case msg.GetImageMessage() != nil:
		return constants.MessageImage
	case msg.GetVideoMessage() != nil:
		return constants.MessageVideo
	case msg.GetDocumentMessage() != nil:
		return constants.MessageDocument
	default:
		return constants.MessageText // Default to text message
	}
}

func extractText(msg *waProto.Message) string {
	switch {
	case msg.GetConversation() != "":
		return msg.GetConversation()
	case msg.GetExtendedTextMessage() != nil:
		return msg.GetExtendedTextMessage().GetText()
	case msg.GetImageMessage() != nil:
		caption := msg.GetImageMessage().GetCaption()
		fmt.Println("Found image with caption:", caption)
		return caption
	case msg.GetVideoMessage() != nil:
		caption := msg.GetVideoMessage().GetCaption()
		fmt.Println("Found video with caption:", caption)
		return caption
	case msg.GetDocumentMessage() != nil:
		caption := msg.GetDocumentMessage().GetCaption()
		fmt.Println("Found document with caption:", caption)
		return caption
	default:
		return ""
	}
}
