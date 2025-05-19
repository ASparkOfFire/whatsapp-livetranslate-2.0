package messagehandler

import (
	"fmt"
	"log"
	"math/rand"
	"sync/atomic"
	"time"

	"github.com/asparkoffire/whatsapp-livetranslate-go/internal/constants"
	"go.mau.fi/whatsmeow/types"
)

var memeRunning int32 = 0

func getRandomEmoji() string {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	return constants.Emojis[rng.Intn(len(constants.Emojis))]
}

func randomEmoji(h *WhatsMeowEventHandler, msgInfo types.MessageInfo, times int) {
	// Try to set the flag from 0 to 1
	if !atomic.CompareAndSwapInt32(&memeRunning, 0, 1) {
		// Already running
		h.editMessageContent(msgInfo.Chat, msgInfo.ID, "Already Running", nil)
		return
	}
	defer atomic.StoreInt32(&memeRunning, 0)

	log.Printf("invoking randmoji routine with loop count: %d\n", times)

	for range times {
		for range 3 {
			time.Sleep(time.Millisecond * 500)
			h.editMessageContent(msgInfo.Chat, msgInfo.ID, getRandomEmoji(), nil)
		}
	}
}

func haha(h *WhatsMeowEventHandler, msgInfo types.MessageInfo, times int) {
	if !atomic.CompareAndSwapInt32(&memeRunning, 0, 1) {
		h.editMessageContent(msgInfo.Chat, msgInfo.ID, "Already Running", nil)
		return
	}
	defer atomic.StoreInt32(&memeRunning, 0)

	var hahaText string
	for range times {
		for range 3 {
			hahaText += "😂"
			time.Sleep(300 * time.Millisecond)
			h.editMessageContent(msgInfo.Chat, msgInfo.ID, fmt.Sprintf("```%s```", hahaText), nil)
		}
	}

	// Descend down to one emoji (length 1)
	for len(hahaText) > 1 {
		hahaText = hahaText[:len(hahaText)-1]
		time.Sleep(150 * time.Millisecond) // Faster descending
		h.editMessageContent(msgInfo.Chat, msgInfo.ID, fmt.Sprintf("```%s```", hahaText), nil)
	}
}
