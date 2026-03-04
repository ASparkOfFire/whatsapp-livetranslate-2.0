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

	OpenrouterModel      string
	OpenrouterBaseUrl    string
	OpenrouterApiKey     string
	OpenrouterImageModel string
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

	AppConfig.OpenrouterBaseUrl = os.Getenv("OPENROUTER_BASEURL")
	AppConfig.OpenrouterModel = os.Getenv("OPENROUTER_MODEL")
	AppConfig.OpenrouterApiKey = os.Getenv("OPENROUTER_APIKEY")
	AppConfig.OpenrouterImageModel = os.Getenv("OPENROUTER_IMAGE_MODEL")
}
