package messagehandler

import (
	"fmt"
	"strings"

	"github.com/asparkoffire/whatsapp-livetranslate-go/internal/constants"
	"github.com/asparkoffire/whatsapp-livetranslate-go/internal/utils"
	"github.com/pemistahl/lingua-go"
	waProto "go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
)

// Update the handler function to handle all media types
func (h *WhatsMeowEventHandler) handleMediaCaption(msg *waProto.Message, msgInfo types.MessageInfo) bool {
	// Extract caption based on media type
	var caption string
	var mediaType string

	if msg.GetImageMessage() != nil && msg.GetImageMessage().GetCaption() != "" {
		caption = msg.GetImageMessage().GetCaption()
		mediaType = "image"
	} else if msg.GetVideoMessage() != nil && msg.GetVideoMessage().GetCaption() != "" {
		caption = msg.GetVideoMessage().GetCaption()
		mediaType = "video"
	} else if msg.GetDocumentMessage() != nil && msg.GetDocumentMessage().GetCaption() != "" {
		caption = msg.GetDocumentMessage().GetCaption()
		mediaType = "document"
	} else {
		return false
	}

	fmt.Printf("Processing %s caption: %s\n", mediaType, caption)

	// Check if caption starts with a language code
	if !strings.HasPrefix(caption, "/") {
		return false
	}

	parts := strings.SplitN(caption, " ", 2)
	if len(parts) < 2 {
		fmt.Println("Caption starts with / but has no space - not a translation command")
		return false
	}

	langCode := strings.TrimPrefix(parts[0], "/")
	targetLang := utils.GetLangByCode(langCode)
	if targetLang == lingua.Unknown {
		fmt.Println("Not a valid language code in media caption - ignoring")
		return false
	}

	mediaCaption := parts[1]
	detectedLang, ok := h.detector.DetectLanguage(mediaCaption)
	if !ok {
		fmt.Printf("%s language detection failed:\n", mediaType)
		return false
	}

	translated, err := h.translator.TranslateText(mediaCaption, detectedLang, targetLang)
	if err != nil {
		fmt.Printf("%s caption translation failed: %v\n", mediaType, err)
		return false
	}

	if msgInfo.IsFromMe {
		// If it's our own message, edit the media caption
		if err := h.editMessageContent(msgInfo.Chat, msgInfo.ID, translated, msg); err != nil {
			fmt.Printf("%s caption edit failed: %v\n", mediaType, err)
			// Fall back to replying with the translation
			h.sendResponse(msgInfo, translated)
		}
	} else {
		// If it's someone else's message, reply with the translated caption
		h.sendResponse(msgInfo, translated)
	}

	return true
}

func (h *WhatsMeowEventHandler) handleQuotedMessageTranslation(msg *waProto.Message, text string, msgInfo types.MessageInfo) bool {
	// Check if message is of form "/<lang>" and is replying to another
	if !strings.HasPrefix(text, "/") || msg.GetExtendedTextMessage() == nil || msg.GetExtendedTextMessage().GetContextInfo() == nil {
		return false
	}

	// For quoted messages, the command is usually just "/en" with no space
	langCode := strings.TrimPrefix(text, "/")
	// Check if there's a space - if so, it's not a pure language code
	if strings.Contains(langCode, " ") {
		fmt.Println("Not a valid quoted translation command (has spaces):", text)
		return false
	}

	targetLang := utils.GetLangByCode(langCode)
	if targetLang == lingua.Unknown {
		fmt.Println("Not a valid language code in quoted message - ignoring")
		return false
	}

	// Extract quoted message
	quoted := msg.GetExtendedTextMessage().GetContextInfo().GetQuotedMessage()
	if quoted == nil {
		fmt.Println("No quoted message found.")
		return true
	}

	// Get information about the quoted message
	quotedMsgID := ""
	quotedSender := ""
	if msg.GetExtendedTextMessage().GetContextInfo().GetStanzaID() != "" {
		quotedMsgID = msg.GetExtendedTextMessage().GetContextInfo().GetStanzaID()
	}
	if msg.GetExtendedTextMessage().GetContextInfo().GetParticipant() != "" {
		quotedSender = msg.GetExtendedTextMessage().GetContextInfo().GetParticipant()
		fmt.Println("Quoted message sender:", quotedSender)
	}

	// Check if the quoted message is a media message
	isMedia := isMediaMessage(quoted)
	mediaType := ""

	if quoted.GetImageMessage() != nil {
		fmt.Println("Quoted message is an image with caption:", quoted.GetImageMessage().GetCaption())
		mediaType = "image"
	} else if quoted.GetVideoMessage() != nil {
		fmt.Println("Quoted message is a video with caption:", quoted.GetVideoMessage().GetCaption())
		mediaType = "video"
	} else if quoted.GetDocumentMessage() != nil {
		fmt.Println("Quoted message is a document with caption:", quoted.GetDocumentMessage().GetCaption())
		mediaType = "document"
	} else if quoted.GetAudioMessage() != nil {
		fmt.Println("Quoted message is an audio message")
		mediaType = "audio"
	}

	quotedText := extractText(quoted)
	if quotedText == "" {
		fmt.Println("Quoted message has no translatable text.")
		return true
	}

	// Detect the language first
	detectedLang, ok := h.detector.DetectLanguage(quotedText)
	if !ok {
		fmt.Println("Could not detect source language.")
		return false
	}

	// Then translate with the detected language
	translated, err := h.translator.TranslateText(quotedText, detectedLang, targetLang)
	if err != nil {
		fmt.Println("Translation error:", err)
		return false
	}

	if msgInfo.IsFromMe && isMedia && quotedMsgID != "" {
		// Try to edit the caption if it's our own media
		fmt.Printf("Attempting to edit %s caption\n", mediaType)
		if err := h.editMessageContent(msgInfo.Chat, quotedMsgID, translated, quoted); err != nil {
			fmt.Printf("Quoted %s caption edit failed: %v\n", mediaType, err)
			// Fall back to editing the original message
			if err := h.editMessageContent(msgInfo.Chat, msgInfo.ID, translated, nil); err != nil {
				fmt.Println("Edit failed:", err)
			}
		}
	} else if msgInfo.IsFromMe {
		// For regular text messages we own
		if err := h.editMessageContent(msgInfo.Chat, msgInfo.ID, translated, nil); err != nil {
			fmt.Println("Edit failed:", err)
		}
	} else {
		// For messages from others, we quote their message
		if err := h.sendReplyMessage(msgInfo.Chat, translated, msgInfo.ID); err != nil {
			fmt.Println("Reply failed:", err)
		}
	}

	return true
}

func (h *WhatsMeowEventHandler) handleDirectTranslation(text string, msgInfo types.MessageInfo) bool {
	// Check if the message starts with "/"
	if !strings.HasPrefix(text, "/") {
		return false
	}

	parts := strings.SplitN(text, " ", 2)
	if len(parts) < 2 {
		fmt.Println("Message starts with / but has no space - not a translation command")
		return false
	}

	langCode := strings.TrimPrefix(parts[0], "/")
	targetLang := utils.GetLangByCode(langCode)
	if targetLang == lingua.Unknown {
		fmt.Println("Not a valid language code - ignoring")
		return false
	}

	textToTranslate := parts[1]

	// Detect the language first
	detectedLang, ok := h.detector.DetectLanguage(textToTranslate)
	if !ok {
		fmt.Println("Could not detect source language.")
		return false
	}

	// Then translate with the detected language
	translated, err := h.translator.TranslateText(textToTranslate, detectedLang, targetLang)
	if err != nil {
		fmt.Println(err)
		return false
	}

	// For direct translation, we're translating the command message's text
	// So we use the regular sendResponse
	h.sendResponse(msgInfo, translated)
	return true
}

func (h *WhatsMeowEventHandler) handleTranslation(msg *waProto.Message, text string, msgInfo types.MessageInfo) bool {
	switch getMessageType(msg) {
	case constants.MessageImage, constants.MessageVideo, constants.MessageDocument:
		// First check if it's a media message with a caption that needs translation
		if h.handleMediaCaption(msg, msgInfo) {
			return true
		}
	default:
		// Try to handle as a quoted message translation
		if h.handleQuotedMessageTranslation(msg, text, msgInfo) {
			return true
		}

		// if all fails, try direct translation
		return h.handleDirectTranslation(text, msgInfo)
	}

	return false
}
