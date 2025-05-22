package messagehandler

import (
	"errors"
	"fmt"

	"github.com/asparkoffire/whatsapp-livetranslate-go/internal/constants"
	waProto "go.mau.fi/whatsmeow/proto/waE2E"
)

// Function to check if a message is a media message
func isMediaMessage(msg *waProto.Message) bool {
	msgType := getMessageType(msg)
	if msgType == constants.MessageImage ||
		msgType == constants.MessageAudio ||
		msgType == constants.MessageVideo ||
		msgType == constants.MessageDocument {
		return true
	}

	return false
}

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
	case msg.GetAudioMessage() != nil:
		return constants.MessageAudio
	case msg.GetLocationMessage() != nil:
		return constants.MessageLocation
	case msg.GetContactMessage() != nil:
		return constants.MessageContact
	case msg.GetPollCreationMessage() != nil:
		return constants.MessagePoll
	default:
		return constants.MessageUnknown
	}
}

func extractText(msg *waProto.Message) string {
	switch getMessageType(msg) {
	case constants.MessageText:
		return msg.GetConversation()
	case constants.MessageExtendedText:
		return msg.GetExtendedTextMessage().GetText()
	case constants.MessageImage:
		return msg.GetImageMessage().GetCaption()
	case constants.MessageVideo:
		return msg.GetVideoMessage().GetCaption()
	case constants.MessageDocument:
		return msg.GetDocumentMessage().GetCaption()
	case constants.MessageAudio:
		return ""
	case constants.MessageLocation:
		loc := msg.GetLocationMessage()
		return loc.GetName() + " (" + loc.GetAddress() + ")"
	case constants.MessageContact:
		return msg.GetContactMessage().GetDisplayName()
	case constants.MessagePoll:
		return msg.GetPollCreationMessage().GetName()
	default:
		return ""
	}
}

func getQuotedMessageAndType(msg *waProto.Message) (*waProto.Message, constants.Message, error) {
	if msg == nil || msg.GetExtendedTextMessage() == nil || msg.GetExtendedTextMessage().GetContextInfo() == nil {
		return nil, constants.MessageText, errors.New("failed to get quoted message")
	}

	quoted := msg.GetExtendedTextMessage().GetContextInfo().GetQuotedMessage()
	if quoted == nil {
		return nil, constants.MessageText, errors.New("failed to get quoted message")
	}

	return quoted, getMessageType(quoted), nil
}

func extractQuotedText(msg *waProto.Message) string {
	if msg == nil || msg.GetExtendedTextMessage() == nil || msg.GetExtendedTextMessage().GetContextInfo() == nil {
		return ""
	}

	quoted := msg.GetExtendedTextMessage().GetContextInfo().GetQuotedMessage()
	if quoted == nil {
		return ""
	}

	switch {
	case quoted.GetConversation() != "":
		return quoted.GetConversation()
	case quoted.GetExtendedTextMessage() != nil:
		return quoted.GetExtendedTextMessage().GetText()
	case quoted.GetImageMessage() != nil:
		return quoted.GetImageMessage().GetCaption()
	case quoted.GetVideoMessage() != nil:
		return quoted.GetVideoMessage().GetCaption()
	case quoted.GetDocumentMessage() != nil:
		return quoted.GetDocumentMessage().GetCaption()
	case quoted.GetAudioMessage() != nil:
		return ""
	case quoted.GetLocationMessage() != nil:
		loc := quoted.GetLocationMessage()
		return "Location: " + loc.GetName() + " (" + loc.GetAddress() + ")"
	case quoted.GetContactMessage() != nil:
		return "Contact: " + quoted.GetContactMessage().GetDisplayName()
	case quoted.GetPollCreationMessage() != nil:
		return "Poll: " + quoted.GetPollCreationMessage().GetName()
	default:
		return ""
	}
}

func getSupportedLanguages() string {
	supportedLangsString := "Supported Languages:\n"
	for code, lang := range constants.SupportedLanguages {
		supportedLangsString += fmt.Sprintf("\n`%s` - %s", code, lang.String())
	}
	return supportedLangsString
}

const HelpMessage string = "WhatsApp Live Translation and Meme bot by Kabir Kalsi (https://github.com/ASparkOfFire)\n\nAvailable Commands:\n\n`/<lang code>` - Translate from one language to another, works inline and in reply to other message.\n`/supportedlangs` - List all the supported languages.\n`/help` - Display this message.\n`/ping` - For healthcheck.\n`/randmoji` - Spams a random emoji that keeps changing for specified duration.\n`/haha` - Spams laughing emoji and later removes it one by one.\n`/setmodel` - Set the translation model (gemini-1.5-flash, gemini-2.0-flash, gemini-2.5-flash).\n`/getmodel` - Show the currently active translation model."
