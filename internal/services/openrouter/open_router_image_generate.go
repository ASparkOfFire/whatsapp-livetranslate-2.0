package openrouter

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/asparkoffire/whatsapp-livetranslate-go/internal/services"
)

type OpenrouterImageGenerator struct {
	apiKey     string
	Model      string
	BaseUrl    string
	client     http.Client
	Temprature float64
}

func NewOpenrouterImageGenerator(model, baseurl, apikey string) services.ImageGenerator {
	return &OpenrouterImageGenerator{
		Model:   model,
		BaseUrl: baseurl,
		apiKey:  apikey,
		client:  http.Client{Timeout: time.Second * 10},
	}
}

// GenerateImage implements [services.ImageGenerator].
func (o *OpenrouterImageGenerator) GenerateImage(ctx context.Context, prompt string) ([]byte, error) {
	apiURL := o.BaseUrl + "/api/v1/chat/completions"

	payload := OpenrouterImageGenerationRequestSchema{
		Model: o.Model,
		Messages: []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		}{
			{Role: "user", Content: prompt},
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+o.apiKey)

	resp, err := o.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("api error (status %d): %s", resp.StatusCode, string(body))
	}

	var result OpenrouterImageGenerationResponseSchema
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(result.Choices) == 0 || len(result.Choices[0].Message.Images) == 0 {
		return nil, fmt.Errorf("no image returned in response")
	}

	rawURL := result.Choices[0].Message.Images[0].ImageURL.URL

	parts := strings.Split(rawURL, ",")
	base64Data := parts[len(parts)-1]

	imgBytes, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 image: %w", err)
	}

	return imgBytes, nil
}
