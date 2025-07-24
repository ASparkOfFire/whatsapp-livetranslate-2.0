package fun

import (
	"fmt"
	"sync/atomic"
	"time"

	framework "github.com/asparkoffire/whatsapp-livetranslate-go/internal/cmdframework"
)

var hahaRunning int32 = 0

type HahaCommand struct{}

func NewHahaCommand() *HahaCommand {
	return &HahaCommand{}
}

func (c *HahaCommand) Execute(ctx *framework.Context) error {
	// Only allow one instance to run at a time
	if !atomic.CompareAndSwapInt32(&hahaRunning, 0, 1) {
		return ctx.Handler.SendResponse(ctx.MessageInfo,
			framework.Warning("Haha is already running"))
	}
	defer atomic.StoreInt32(&hahaRunning, 0)

	// Run haha animation in background
	go func() {
		var hahaText string
		for range 3 {
			for range 3 {
				hahaText += "😂"
				time.Sleep(300 * time.Millisecond)
				ctx.Handler.EditMessage(ctx.MessageInfo, fmt.Sprintf("```%s```", hahaText))
			}
		}
	}()

	return nil
}

func (c *HahaCommand) Metadata() *framework.Metadata {
	return &framework.Metadata{
		Name:         "haha",
		Description:  "Laughing emoji animation",
		Category:     "Fun",
		Usage:        "/haha",
		RequireOwner: true,
		Hidden:       true, // Hide from help as it's a fun easter egg
	}
}
