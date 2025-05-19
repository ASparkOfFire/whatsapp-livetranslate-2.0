package messagehandler

import (
	"log"
	"math/rand"
	"sync/atomic"
	"time"

	"github.com/asparkoffire/whatsapp-livetranslate-go/internal/constants"
	"go.mau.fi/whatsmeow/types"
)

var randmojiRunning int32 = 0

func getRandomEmoji() string {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	return constants.Emojis[rng.Intn(len(constants.Emojis))]
}

func randomEmoji(h *WhatsMeowEventHandler, msgInfo types.MessageInfo, duration int) {
	// Try to set the flag from 0 to 1
	if !atomic.CompareAndSwapInt32(&randmojiRunning, 0, 1) {
		// Already running
		h.editMessageContent(msgInfo.Chat, msgInfo.ID, "Already Running", nil)
		return
	}
	defer atomic.StoreInt32(&randmojiRunning, 0)

	log.Printf("invoking randmoji routine with duration: %d seconds\n", duration)
	start := time.Now()
	for {
		if time.Since(start) > time.Duration(duration)*time.Second {
			return
		}
		time.Sleep(time.Millisecond * 500)
		h.editMessageContent(msgInfo.Chat, msgInfo.ID, getRandomEmoji(), nil)
	}
}
