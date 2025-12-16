package messagehandler

import (
	"context"
	"fmt"
	"strings"

	framework "github.com/asparkoffire/whatsapp-livetranslate-go/internal/cmdframework"
	"github.com/asparkoffire/whatsapp-livetranslate-go/internal/constants"
	"github.com/asparkoffire/whatsapp-livetranslate-go/internal/handlers"
	"github.com/asparkoffire/whatsapp-livetranslate-go/internal/handlers/admin"
	"github.com/asparkoffire/whatsapp-livetranslate-go/internal/handlers/fun"
	"github.com/asparkoffire/whatsapp-livetranslate-go/internal/handlers/translation"
	"github.com/asparkoffire/whatsapp-livetranslate-go/internal/handlers/utility"
	waProto "go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
)

// handleMessage uses the new command system
func (h *WhatsMeowEventHandler) handleMessage(msg *waProto.Message, msgInfo types.MessageInfo) {
	text := extractText(msg)
	var cmdName string
	var args []string
	var rawArgs string

	if strings.HasPrefix(text, "s/") {
		cmdName = "s"
		rawArgs = text
		args = []string{text}
	} else if strings.HasPrefix(text, "/") {
		parts := strings.Fields(text)
		if len(parts) == 0 {
			return
		}
		cmdName = strings.TrimPrefix(parts[0], "/")
		args = parts[1:]
		if len(parts) > 1 {
			rawArgs = strings.Join(parts[1:], " ")
		}
	} else {
		// Handle non-command messages - check for AFK mode
		if h.IsAfkMode() {
			adapter := NewHandlerAdapter(h)
			response := `The person you are trying to reach is not available at the moment, in case of an urgency - Reach out via call.`
			_ = adapter.SendResponse(msgInfo, response)
		}
		return
	}

	if cmdName == "" {
		return
	}

	// Look up command in registry
	cmd, exists := h.commandRegistry.Get(cmdName)
	if !exists {
		// Check if it's a language code for translation
		if len(cmdName) == 2 {
			if _, isLang := constants.SupportedLanguages[cmdName]; isLang {
				cmd, _ = h.commandRegistry.Get(cmdName)
			}
		}

		if cmd == nil {
			return // Silent fail for unknown commands
		}
	}

	// Create command context
	adapter := NewHandlerAdapter(h)
	ctx := &framework.Context{
		Context:     context.Background(),
		Message:     msg,
		MessageInfo: msgInfo,
		Command:     cmdName,
		Args:        args,
		RawArgs:     rawArgs,
		Handler:     adapter,
	}

	// Execute command
	if err := cmd.Execute(ctx); err != nil {
		fmt.Printf("Command execution error: %v\n", err)
	}
}

// InitializeCommands sets up all commands in the registry
func (h *WhatsMeowEventHandler) InitializeCommands() error {
	registry := h.commandRegistry

	// Register help command
	helpCmd := handlers.NewHelpCommand(registry)
	if err := registry.Register(helpCmd); err != nil {
		return fmt.Errorf("failed to register help command: %w", err)
	}

	// Register utility commands
	if err := registry.Register(utility.NewPingCommand()); err != nil {
		return fmt.Errorf("failed to register ping command: %w", err)
	}

	if err := registry.Register(utility.NewSupportedLangsCommand()); err != nil {
		return fmt.Errorf("failed to register supportedlangs command: %w", err)
	}

	if err := registry.Register(utility.NewDownloadCommand()); err != nil {
		return fmt.Errorf("failed to register download command: %w", err)
	}

	if err := registry.Register(utility.NewSedCommand()); err != nil {
		return fmt.Errorf("failed to register sed command: %w", err)
	}

	if err := registry.Register(utility.NewHIBPCommand()); err != nil {
		return fmt.Errorf("failed to register hibp command: %w", err)
	}

	afkCmd := utility.NewAfkCommand(h)
	if err := registry.Register(afkCmd); err != nil {
		return fmt.Errorf("failed to register afk command: %w", err)
	}

	noAfkCmd := utility.NewNoAfkCommand(h)
	if err := registry.Register(noAfkCmd); err != nil {
		return fmt.Errorf("failed to register noafk command: %w", err)
	}

	// Register admin commands
	if err := registry.Register(admin.NewSetModelCommand()); err != nil {
		return fmt.Errorf("failed to register setmodel command: %w", err)
	}

	if err := registry.Register(admin.NewGetModelCommand()); err != nil {
		return fmt.Errorf("failed to register getmodel command: %w", err)
	}

	if err := registry.Register(admin.NewSetTempCommand()); err != nil {
		return fmt.Errorf("failed to register settemp command: %w", err)
	}

	if err := registry.Register(admin.NewGetTempCommand()); err != nil {
		return fmt.Errorf("failed to register gettemp command: %w", err)
	}

	// Register fun commands
	if err := registry.Register(fun.NewImageCommand()); err != nil {
		return fmt.Errorf("failed to register image command: %w", err)
	}

	if err := registry.Register(fun.NewMemeCommand()); err != nil {
		return fmt.Errorf("failed to register meme command: %w", err)
	}

	if err := registry.Register(fun.NewRandmojiCommand()); err != nil {
		return fmt.Errorf("failed to register randmoji command: %w", err)
	}

	if err := registry.Register(fun.NewHahaCommand()); err != nil {
		return fmt.Errorf("failed to register haha command: %w", err)
	}

	// Register translation commands for all supported languages
	if err := translation.RegisterTranslationCommands(registry); err != nil {
		return fmt.Errorf("failed to register translation commands: %w", err)
	}

	// Apply middleware to commands that need owner permissions based on metadata
	allCommands := registry.GetAll()
	for name, cmd := range allCommands {
		meta := cmd.Metadata()
		if meta.RequireOwner {
			wrappedCmd := framework.WithMiddleware(cmd, framework.RequireOwner())
			// Update the command with middleware
			registry.UpdateCommand(name, wrappedCmd)
		}
	}

	return nil
}
