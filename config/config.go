package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type config struct {
	GeminiAPIKey string `validate:"required"`
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
}
