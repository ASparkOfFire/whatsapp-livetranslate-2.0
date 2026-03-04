package openrouter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/asparkoffire/whatsapp-livetranslate-go/internal/constants"
	"github.com/asparkoffire/whatsapp-livetranslate-go/internal/services"
	lingua "github.com/pemistahl/lingua-go"
)

type OpenrouterTranslator struct {
	apiKey     string
	Model      string
	BaseUrl    string
	client     http.Client
	Temprature float64
}

func NewOpenrouterTranslator(model string, baseurl string, apiKey string) services.TranslateService {
	return &OpenrouterTranslator{
		Model:   model,
		BaseUrl: baseurl,
		client:  http.Client{Timeout: time.Second * 15},
		apiKey:  apiKey,
	}
}

// GetModel implements [services.TranslateService].
func (o *OpenrouterTranslator) GetModel() string {
	return o.Model
}

// GetTemperature implements [services.TranslateService].
func (o *OpenrouterTranslator) GetTemperature() float64 {
	return o.Temprature
}

// SetModel implements [services.TranslateService].
func (o *OpenrouterTranslator) SetModel(modelID string) error {
	o.Model = modelID
	return nil
}

// SetTemperature implements [services.TranslateService].
func (o *OpenrouterTranslator) SetTemperature(temp float64) error {
	o.Temprature = temp
	return nil
}

// TranslateText implements [services.TranslateService].
func (o *OpenrouterTranslator) TranslateText(text string, sourceLang lingua.Language, targetLang lingua.Language) (string, error) {
	url := o.BaseUrl + "/api/v1/chat/completions"
	auth := "Bearer " + o.apiKey

	body := OpenrouterTranslateRequestSchema{
		Model:       o.Model,
		Stream:      false,
		Temperature: o.Temprature,
		Messages: []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		}{
			{
				Role:    "system",
				Content: constants.SystemPromptMessage,
			},
			{
				Role:    "user",
				Content: fmt.Sprintf("Translate the following text from %s to %s:\n\n%s", sourceLang.String(), targetLang.String(), text),
			},
		},
		ResponseFormat: struct {
			Type       string `json:"type"`
			JSONSchema struct {
				Name   string `json:"name"`
				Strict bool   `json:"strict"`
				Schema struct {
					Type       string `json:"type"`
					Properties struct {
						Output struct {
							Type        string `json:"type"`
							Description string `json:"description"`
						} `json:"output"`
					} `json:"properties"`
					Required             []string `json:"required"`
					AdditionalProperties bool     `json:"additionalProperties"`
				} `json:"schema"`
			} `json:"json_schema"`
		}{
			Type: "json_schema",
			JSONSchema: struct {
				Name   string `json:"name"`
				Strict bool   `json:"strict"`
				Schema struct {
					Type       string `json:"type"`
					Properties struct {
						Output struct {
							Type        string `json:"type"`
							Description string `json:"description"`
						} `json:"output"`
					} `json:"properties"`
					Required             []string `json:"required"`
					AdditionalProperties bool     `json:"additionalProperties"`
				} `json:"schema"`
			}{
				Name:   "translation_response",
				Strict: true,
				Schema: struct {
					Type       string `json:"type"`
					Properties struct {
						Output struct {
							Type        string `json:"type"`
							Description string `json:"description"`
						} `json:"output"`
					} `json:"properties"`
					Required             []string `json:"required"`
					AdditionalProperties bool     `json:"additionalProperties"`
				}{
					Type: "object",
					Properties: struct {
						Output struct {
							Type        string `json:"type"`
							Description string `json:"description"`
						} `json:"output"`
					}{
						Output: struct {
							Type        string `json:"type"`
							Description string `json:"description"`
						}{
							Type:        "string",
							Description: "The translated text in Punjabi",
						},
					},
					Required:             []string{"output"},
					AdditionalProperties: false,
				},
			},
		},
	}

	b, err := json.Marshal(body)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", auth)

	resp, err := o.client.Do(req)
	if err != nil {
		return "", err
	}

	b, err = io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result OpenrouterTranslateResponseSchema
	if err := json.Unmarshal(b, &result); err != nil {
		return "", err
	}
	modelResult := result.Choices[0].Message.Content
	output := ModelOutputSchema{}
	if err := json.Unmarshal([]byte(modelResult), &output); err != nil {
		return "", err
	}
	return output.Output, nil
}
