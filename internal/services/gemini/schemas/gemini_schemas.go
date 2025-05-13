package gemini

import (
	"fmt"

	"github.com/asparkoffire/whatsapp-livetranslate-go/internal/constants"
	"github.com/pemistahl/lingua-go"
)

type geminiLLMInferenceRequest struct {
	GenerationConfig  geminiModelConfig       `json:"generationConfig"`
	Contents          []geminiMessage         `json:"contents"`
	SystemInstruction geminiSystemInstruction `json:"systemInstruction"`
}

type geminiModelConfig struct {
	Temperature      float64              `json:"temperature"`
	TopP             float64              `json:"topP"`
	TopK             int                  `json:"topK"`
	MaxOutputTokens  int                  `json:"maxOutputTokens"`
	ResponseMimeType string               `json:"responseMimeType"`
	ResponseSchema   geminiResponseSchema `json:"responseSchema"`
}

type geminiResponseSchema struct {
	Type       string                   `json:"type"`
	Properties geminiResponseProperties `json:"properties"`
}

type geminiResponseProperties struct {
	Output geminiResponseProperty `json:"output"`
}

type geminiResponseProperty struct {
	Type string `json:"type"`
}

type geminiMessage struct {
	Role  string       `json:"role"`
	Parts []geminiPart `json:"parts"`
}

type geminiPart struct {
	Text string `json:"text"`
}

type geminiSystemInstruction struct {
	Role  string       `json:"role"`
	Parts []geminiPart `json:"parts"`
}

func NewGeminiLLMInferenceRequest(inputText string, sourceLang lingua.Language, targetLang lingua.Language) geminiLLMInferenceRequest {
	return geminiLLMInferenceRequest{
		GenerationConfig: geminiModelConfig{
			Temperature:      0.5,
			TopP:             0.95,
			TopK:             40,
			MaxOutputTokens:  8192,
			ResponseMimeType: "application/json",
			ResponseSchema: geminiResponseSchema{
				Type: "object",
				Properties: geminiResponseProperties{
					Output: geminiResponseProperty{
						Type: "string",
					},
				},
			},
		},
		Contents: []geminiMessage{
			{
				Role: "user",
				Parts: []geminiPart{
					{Text: fmt.Sprintf("Translate the below text from %s to %s:\n\n%s", sourceLang.String(), targetLang.String(), inputText)},
				},
			},
		},
		SystemInstruction: geminiSystemInstruction{
			Role: "user",
			Parts: []geminiPart{
				{Text: constants.SystemPromptMessage},
			},
		},
	}
}

type GeminiResponses struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

type OutputText struct {
	Output string `json:"output"`
}
