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

func (h *WhatsMeowEventHandler) handleMediaCaptionTranslation(msg *waProto.Message, msgInfo types.MessageInfo) bool {
	// check if its a media message
	if !isMediaMessage(msg) {
		return false
	}

	caption := extractText(msg)
	if caption == "" || !strings.HasPrefix(caption, "/") {
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

	textToTranslate := parts[1]
	detectedLang, ok := h.detector.DetectLanguage(textToTranslate)
	if !ok {
		fmt.Printf("language detection failed:\n")
		return false
	}

	translated, err := h.translator.TranslateText(textToTranslate, detectedLang, targetLang)
	if err != nil {
		fmt.Printf("caption translation failed: %v\n", err)
		return false
	}

	if msgInfo.IsFromMe {
		if err := h.editMessageContent(msgInfo.Chat, msgInfo.ID, translated, msg); err != nil {
			fmt.Printf("caption edit failed: %v\n", err)
			h.SendResponse(msgInfo, translated)
		}
	} else {
		h.SendResponse(msgInfo, translated)
	}

	return true
}

func (h *WhatsMeowEventHandler) handleQuotedMessageTranslation(msg *waProto.Message, text string, msgInfo types.MessageInfo) bool {
	if !strings.HasPrefix(text, "/") {
		return false
	}

	parts := strings.SplitN(text, " ", 2)
	if len(parts) > 1 {
		fmt.Printf("Inline translation command \"%s\": %s\n", parts[0], text)
		return false
	}

	langCode := strings.TrimPrefix(text, "/")
	targetLang := utils.GetLangByCode(langCode)
	if targetLang == lingua.Unknown {
		fmt.Printf("Not a valid language code \"%s\" in quoted message - ignoring\n", langCode)
		return false
	}

	quotedMsg, msgType, err := getQuotedMessageAndType(msg)
	if err != nil {
		fmt.Println("Quoted message not found or invalid:", err)
		return true
	}

	quotedText := extractText(quotedMsg)
	if quotedText == "" {
		fmt.Println("Quoted message has no translatable text.")
		return true
	}

	detectedLang, ok := h.detector.DetectLanguage(quotedText)
	if !ok {
		fmt.Println("Could not detect source language.")
		return false
	}

	translated, err := h.translator.TranslateText(quotedText, detectedLang, targetLang)
	if err != nil {
		fmt.Println("Translation error:", err)
		return false
	}

	quotedMsgID := msg.GetExtendedTextMessage().GetContextInfo().GetStanzaID()

	isMedia := isMediaMessage(msg)

	if msgInfo.IsFromMe && isMedia && quotedMsgID != "" {
		fmt.Printf("Attempting to edit %s caption\n", msgType)
		if err := h.editMessageContent(msgInfo.Chat, quotedMsgID, translated, quotedMsg); err != nil {
			fmt.Printf("Quoted %s caption edit failed: %v\n", msgType, err)
			if err := h.editMessageContent(msgInfo.Chat, msgInfo.ID, translated, nil); err != nil {
				fmt.Println("Edit fallback failed:", err)
			}
		}
	} else if msgInfo.IsFromMe {
		if err := h.editMessageContent(msgInfo.Chat, msgInfo.ID, translated, nil); err != nil {
			fmt.Println("Edit failed:", err)
		}
	} else {
		if err := h.sendReplyMessage(msgInfo.Chat, translated, msgInfo.ID); err != nil {
			fmt.Println("Reply failed:", err)
		}
	}

	return true
}

func (h *WhatsMeowEventHandler) handleInlineTranslation(text string, msgInfo types.MessageInfo) bool {
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

	detectedLang, ok := h.detector.DetectLanguage(textToTranslate)
	if !ok {
		fmt.Println("Could not detect source language.")
		return false
	}

	translated, err := h.translator.TranslateText(textToTranslate, detectedLang, targetLang)
	if err != nil {
		fmt.Println("Translation failed:", err)
		return false
	}

	h.SendResponse(msgInfo, translated)
	return true
}

func (h *WhatsMeowEventHandler) handleTranslation(msg *waProto.Message, text string, msgInfo types.MessageInfo) bool {
	switch getMessageType(msg) {
	case constants.MessageImage, constants.MessageVideo, constants.MessageDocument:
		if h.handleMediaCaptionTranslation(msg, msgInfo) {
			return true
		}
	}

	if h.handleQuotedMessageTranslation(msg, text, msgInfo) {
		return true
	}

	return h.handleInlineTranslation(text, msgInfo)
}
