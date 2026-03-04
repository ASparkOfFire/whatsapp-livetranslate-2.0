package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/asparkoffire/whatsapp-livetranslate-go/config"
	"github.com/asparkoffire/whatsapp-livetranslate-go/internal/constants"
	"github.com/asparkoffire/whatsapp-livetranslate-go/internal/services"
	"github.com/asparkoffire/whatsapp-livetranslate-go/internal/services/messagehandler"
	"github.com/asparkoffire/whatsapp-livetranslate-go/internal/services/openrouter"
	_ "github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
)

func main() {
	ctx := context.Background()
	container, err := sqlstore.New(ctx, "sqlite3", "file:./data/auth.db?_foreign_keys=on", nil)
	if err != nil {
		log.Fatalf("error while opening a database connection: %v\n", err)
		return
	}

	deviceStore, err := container.GetFirstDevice(ctx)
	if err != nil {
		log.Fatalf("error while getting the device store : %v\n", err)
		return
	}

	client := whatsmeow.NewClient(deviceStore, nil)
	translator := openrouter.NewOpenrouterTranslator(config.AppConfig.OpenrouterModel, config.AppConfig.OpenrouterBaseUrl, config.AppConfig.OpenrouterApiKey)
	imageGenerator := openrouter.NewOpenrouterImageGenerator(config.AppConfig.OpenrouterImageModel, config.AppConfig.OpenrouterBaseUrl, config.AppConfig.OpenrouterApiKey)

	// Initialize the language detector with supported languages
	detector := services.NewLinguaLangDetectService(constants.SupportedLanguages)

	// connect to the client and event handler
	evtHandler, err := messagehandler.NewWhatsMeowEventHandler(client, detector, translator, imageGenerator)
	if err != nil {
		log.Fatalf("error while setting up the event handler: %v\n", err)
		return
	}

	client.AddEventHandler(evtHandler.HandleEvents)
	fmt.Println("Server started, Listening for messages...")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	client.Disconnect()
}
