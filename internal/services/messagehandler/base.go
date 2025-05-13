package messagehandler

import (
	"github.com/asparkoffire/whatsapp-livetranslate-go/internal/services"
	"go.mau.fi/whatsmeow"
)

type WhatsMeowEventHandler struct {
	client     *whatsmeow.Client
	detector   services.LangDetectService
	translator services.TranslateService
}

func NewWhatsMeowEventHandler(client *whatsmeow.Client, detector services.LangDetectService, translator services.TranslateService) (*WhatsMeowEventHandler, error) {
	handler := &WhatsMeowEventHandler{
		client:     client,
		detector:   detector,
		translator: translator,
	}
	if handler.client.Store.ID == nil {
		if err := handler.setupQRLogin(); err != nil {
			return nil, err
		}
	} else {
		if err := client.Connect(); err != nil {
			return nil, err
		}
	}
	return handler, nil
}
