package messagehandler

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	waProto "go.mau.fi/whatsmeow/proto/waE2E"

	"github.com/asparkoffire/whatsapp-livetranslate-go/internal/constants"
	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow/types"
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
	case "randmoji":
		if msgInfo.IsFromMe {
			duration := 10 // default duration

			if len(parts) > 1 {
				if d, err := strconv.Atoi(parts[1]); err == nil && d > 0 && d <= 10 {
					duration = d
				}
			}

			go randomEmoji(h, msgInfo, duration)
		}
	case "ping":
		if msgInfo.IsFromMe {
			h.SendResponse(msgInfo, fmt.Sprintf("Pong: %s", time.Since(start).String()))
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
