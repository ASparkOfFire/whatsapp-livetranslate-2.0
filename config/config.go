package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type config struct {
	GeminiAPIKey string `validate:"required"`
	HIBPToken    string // HIBP API token for dark web searches
	HIBPURL      string
}

var (
	AppConfig = new(config)
)

func init() {
	if strings.ToLower(os.Getenv("IS_DOCKER")) != "true" {
		if err := godotenv.Load(); err != nil {
			log.Fatalf("error loading env variables: %v\n", err)
			return
		}
	}

	AppConfig.GeminiAPIKey = os.Getenv("GEMINI_API_KEY")
	AppConfig.HIBPToken = os.Getenv("HIBP_TOKEN")
	AppConfig.HIBPURL = os.Getenv("HIBP_URL")
}
