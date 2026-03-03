package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type config struct {
	GeminiAPIKey  string
	OllamaModel   string
	OllamaBaseUrl string
}

var (
	AppConfig = new(config)
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Printf("error loading env variables, proceeding without it %v\n", err)
	}

	AppConfig.GeminiAPIKey = os.Getenv("GEMINI_API_KEY")

	AppConfig.OllamaModel = os.Getenv("OLLAMA_MODEL")
	AppConfig.OllamaBaseUrl = os.Getenv("OLLAMA_BASEURL")
}
