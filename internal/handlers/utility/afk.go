package utility

import (
	framework "github.com/asparkoffire/whatsapp-livetranslate-go/internal/cmdframework"
)

type AfkModeController interface {
	SetAfkMode(bool)
	IsAfkMode() bool
}

type AfkCommand struct {
	eventHandler AfkModeController
}

type NoAfkCommand struct {
	eventHandler AfkModeController
}

func NewAfkCommand(eventHandler AfkModeController) *AfkCommand {
	return &AfkCommand{
		eventHandler: eventHandler,
	}
}

func NewNoAfkCommand(eventHandler AfkModeController) *NoAfkCommand {
	return &NoAfkCommand{
		eventHandler: eventHandler,
	}
}

func (c *AfkCommand) Execute(ctx *framework.Context) error {
	c.eventHandler.SetAfkMode(true)
	response := "✅ AFK mode enabled. Anyone messaging will receive your AFK message."
	return ctx.Handler.SendResponse(ctx.MessageInfo, response)
}

func (c *NoAfkCommand) Execute(ctx *framework.Context) error {
	c.eventHandler.SetAfkMode(false)
	response := "❌ AFK mode disabled."
	return ctx.Handler.SendResponse(ctx.MessageInfo, response)
}

func (c *AfkCommand) Metadata() *framework.Metadata {
	return &framework.Metadata{
		Name:         "afk",
		Aliases:      []string{},
		Description:  "Enable AFK mode to send automatic responses",
		Category:     "Utility",
		Usage:        "/afk",
		Examples:     []string{"/afk"},
		RequireOwner: true,
	}
}

func (c *NoAfkCommand) Metadata() *framework.Metadata {
	return &framework.Metadata{
		Name:         "noafk",
		Aliases:      []string{},
		Description:  "Disable AFK mode",
		Category:     "Utility",
		Usage:        "/noafk",
		Examples:     []string{"/noafk"},
		RequireOwner: true,
	}
}