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
	case "help":
		h.SendResponse(msgInfo, HelpMessage)
	case "supportedlangs":
		h.SendResponse(msgInfo, getSupportedLanguages())
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
	case "haha":
		if msgInfo.IsFromMe {
			duration := 3 // default duration
			go haha(h, msgInfo, duration)
		}
	case "ping":
		if msgInfo.IsFromMe {
			h.SendResponse(msgInfo, fmt.Sprintf("Pong: %s", time.Since(start).String()))
		}
	case "setmodel":
		if msgInfo.IsFromMe {
			if len(parts) < 2 {
				h.SendResponse(msgInfo, "Please specify a model ID. Supported models: gemini-1.5-flash, gemini-2.0-flash, gemini-2.5-flash")
				return
			}
			modelID := parts[1]
			if err := h.translator.SetModel(modelID); err != nil {
				h.SendResponse(msgInfo, fmt.Sprintf("Error setting model: %v", err))
				return
			}
			h.SendResponse(msgInfo, fmt.Sprintf("Successfully set translation model to: %s", modelID))
		}
	case "getmodel":
		h.SendResponse(msgInfo, fmt.Sprintf("Current translation model: %s", h.translator.GetModel()))
	case "settemp":
		if msgInfo.IsFromMe {
			if len(parts) < 2 {
				h.SendResponse(msgInfo, "Please specify a temperature value between 0.0 and 1.0")
				return
			}
			temp, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				h.SendResponse(msgInfo, "Invalid temperature value. Please provide a number between 0.0 and 1.0")
				return
			}
			if err := h.translator.SetTemperature(temp); err != nil {
				h.SendResponse(msgInfo, fmt.Sprintf("Error setting temperature: %v", err))
				return
			}
			h.SendResponse(msgInfo, fmt.Sprintf("Successfully set temperature to: %.1f", temp))
		}
	case "gettemp":
		h.SendResponse(msgInfo, fmt.Sprintf("Current temperature: %.1f", h.translator.GetTemperature()))
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
