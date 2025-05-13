package messagehandler

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	waProto "go.mau.fi/whatsmeow/proto/waE2E"

	"github.com/asparkoffire/whatsapp-livetranslate-go/internal/constants"
	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

func (h *WhatsMeowEventHandler) handleMessage(msg *waProto.Message, msgInfo types.MessageInfo) {
	start := time.Now()

	text := extractText(msg)

	if text == "" || !strings.HasPrefix(text, "/") {
		return
	}
	parts := strings.Split(text, " ")
	if parts[0] == "" {
		return
	}

	cmd := strings.TrimPrefix(parts[0], "/")
	switch cmd {
	case "ping":
		if msgInfo.IsFromMe {
			h.sendResponse(msgInfo, fmt.Sprintf("Pong: %s", time.Since(start).String()))
		}
	default:
		if len(cmd) == 2 { // it is a two digits language code.
			if _, ok := constants.SupportedLanguages[cmd]; !ok {
				return // dont handle if it is an invalid code
			}
			h.handleTranslation(msg, text, msgInfo)
		}
	}

}

func (h *WhatsMeowEventHandler) HandleEvents(evt any) {
	switch v := evt.(type) {
	case *events.Message:
		h.handleMessage(v.Message, v.Info)
	}
}

func (h *WhatsMeowEventHandler) setupQRLogin() error {
	qrChan, _ := h.client.GetQRChannel(context.Background())
	err := h.client.Connect()
	if err != nil {
		return err
	}

	for evt := range qrChan {
		if evt.Event == "code" {
			qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
		} else {
			fmt.Println("Login event:", evt.Event)
		}
	}
	return nil
}

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
