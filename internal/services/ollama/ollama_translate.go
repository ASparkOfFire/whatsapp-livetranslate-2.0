package ollama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/asparkoffire/whatsapp-livetranslate-go/internal/services"
	lingua "github.com/pemistahl/lingua-go"
)

type OllamaTranslator struct {
	Model   string
	BaseUrl string
	client  http.Client
}

func NewOllamaTranslator(model string, baseUrl string) services.TranslateService {
	client := http.Client{Timeout: 10 * time.Second}
	ollamaTranslator := &OllamaTranslator{
		BaseUrl: baseUrl,
		client:  client,
		Model:   model,
	}

	return ollamaTranslator
}

// GetModel implements [TranslateService].
func (o *OllamaTranslator) GetModel() string {
	return o.Model
}

// GetTemperature implements [TranslateService].
func (o *OllamaTranslator) GetTemperature() float64 {
	return 0.2
}

// SetModel implements [TranslateService].
func (o *OllamaTranslator) SetModel(modelID string) error {
	return nil
}

// SetTemperature implements [TranslateService].
func (o *OllamaTranslator) SetTemperature(temp float64) error {
	return nil
}

// TranslateText implements [TranslateService].
func (o *OllamaTranslator) TranslateText(text string, sourceLang lingua.Language, targetLang lingua.Language) (string, error) {
	log.Printf("Recieved translation request for: %s to %s with text: %s", sourceLang.IsoCode639_1(), targetLang.IsoCode639_1(), text)
	req := OllamaTranslateRequestSchema{
		Model: o.Model,
		Messages: []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		}{
			{
				Role:    "user",
				Content: fmt.Sprintf("Translate the below text from %s to %s:\n\n%s", sourceLang.String(), targetLang.String(), text),
			},
		},
		Stream: false,
		Format: struct {
			Type       string `json:"type"`
			Properties struct {
				Output struct {
					Type    string `json:"type"`
					Example string `json:"example"`
				} `json:"output"`
			} `json:"properties"`
		}{
			Type: "object",
			Properties: struct {
				Output struct {
					Type    string `json:"type"`
					Example string `json:"example"`
				} `json:"output"`
			}{
				Output: struct {
					Type    string `json:"type"`
					Example string `json:"example"`
				}{
					Type:    "string",
					Example: "this is translated text",
				},
			},
		},
		Required: []string{"output"},
	}
	b, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	url := o.BaseUrl + "/api/chat"
	apiResp, err := o.client.Post(url, "application/json", bytes.NewReader(b))
	if err != nil {
		return "", err
	}
	b, err = io.ReadAll(apiResp.Body)
	if err != nil {
		return "", err
	}

	var resp OllamaTranslateResponseSchema
	if err := json.Unmarshal(b, &resp); err != nil {
		return "", err
	}
	var output ModelOutputSchema
	if err := json.Unmarshal([]byte(resp.Message.Content), &output); err != nil {
		return "", err
	}
	return output.Output, nil
}
