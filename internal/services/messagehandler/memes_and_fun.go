package messagehandler

import (
	"math/rand"
	"time"

	"go.mau.fi/whatsmeow/types"
)

// emojiRange represents a start and end of a Unicode range
type emojiRange struct {
	start, end rune
}

var emojiRanges = []emojiRange{
	{0x1F600, 0x1F64F}, // Emoticons
	{0x2600, 0x26FF},   // Miscellaneous Symbols
	{0x2700, 0x27BF},   // Dingbats
	{0x1F680, 0x1F6FF}, // Transport and Map
	{0x1F900, 0x1F9FF}, // Symbols and Pictographs Extended-A
	{0x1F1E6, 0x1F1FF}, // Regional Indicator Symbols
}

func getRandomEmoji() string {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Pick a random range
	r := emojiRanges[rng.Intn(len(emojiRanges))]

	// Pick a random rune within that range
	codePoint := rng.Intn(int(r.end-r.start)+1) + int(r.start)

	return string(rune(codePoint))
}

func randomEmoji(h *WhatsMeowEventHandler, msgInfo types.MessageInfo, duration int) {
	start := time.Now()
	for {
		if time.Since(start) > time.Duration(duration) {
			return
		}
		time.Sleep(time.Millisecond * 500)
		h.editMessageContent(msgInfo.Chat, msgInfo.ID, getRandomEmoji(), nil)
	}
}
