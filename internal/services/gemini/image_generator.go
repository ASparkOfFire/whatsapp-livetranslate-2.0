package gemini

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/asparkoffire/whatsapp-livetranslate-go/internal/services"
)

// -----------------------------------------------------------------------------
// types for (partial) response decoding
// -----------------------------------------------------------------------------

type inlineData struct {
	MimeType string `json:"mimeType"`
	Data     string `json:"data"`
}

type part struct {
	Text       string      `json:"text,omitempty"`
	InlineData *inlineData `json:"inlineData,omitempty"`
}

type content struct {
	Parts []part `json:"parts"`
	Role  string `json:"role"`
}

type candidate struct {
	Content content `json:"content"`
	Index   int     `json:"index"`
}

type event struct {
	Candidates []candidate `json:"candidates"`
}

// -----------------------------------------------------------------------------
// request payload
// -----------------------------------------------------------------------------

type requestPayload struct {
	Contents         []map[string]any `json:"contents"`
	GenerationConfig map[string]any   `json:"generationConfig"`
}

// -----------------------------------------------------------------------------
// geminiImageGenerator
// -----------------------------------------------------------------------------

type geminiImageGenerator struct {
	client             *http.Client
	apiKey             string
	modelID            string
	generateContentAPI string
	geminiAPIBaseURL   string
}

// NewGeminiImageGenerator returns a services.ImageGenerator implementation.
func NewGeminiImageGenerator(model, apiKey string) services.ImageGenerator {
	return &geminiImageGenerator{
		modelID:            model,
		generateContentAPI: "streamGenerateContent",
		geminiAPIBaseURL:   "https://generativelanguage.googleapis.com/v1beta/models",
		apiKey:             apiKey,
		client:             &http.Client{Timeout: 30 * time.Second},
	}
}

// GenerateImage sends the prompt and returns decoded image bytes.
func (g *geminiImageGenerator) GenerateImage(ctx context.Context, prompt string) ([]byte, error) {
	// ---------- build request body ----------
	payload := requestPayload{
		Contents: []map[string]any{
			{
				"role": "user",
				"parts": []map[string]any{
					{"text": prompt},
				},
			},
		},
		GenerationConfig: map[string]any{
			"responseModalities": []string{"IMAGE", "TEXT"},
			"responseMimeType":   "text/plain",
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	// ---------- HTTP request ----------
	url := fmt.Sprintf("%s/%s:%s?key=%s",
		g.geminiAPIBaseURL, g.modelID, g.generateContentAPI, g.apiKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBuf, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, string(respBuf))
	}

	// ---------- stream & extract image ----------
	imageData, err := extractImageData(resp.Body)
	if err != nil {
		return nil, err
	}

	decoded, err := base64.StdEncoding.DecodeString(imageData)
	if err != nil {
		return nil, fmt.Errorf("base64-decode image: %w", err)
	}

	return decoded, nil
}

// -----------------------------------------------------------------------------
// helpers
// -----------------------------------------------------------------------------

// extractImageData reads the JSON stream and returns the first inlineData.data.
func extractImageData(r io.Reader) (string, error) {
	var events []event
	if err := json.NewDecoder(r).Decode(&events); err != nil {
		return "", fmt.Errorf("decode events array: %w", err)
	}

	for _, ev := range events {
		for _, cand := range ev.Candidates {
			for _, p := range cand.Content.Parts {
				if p.InlineData != nil && p.InlineData.Data != "" {
					return p.InlineData.Data, nil
				}
			}
		}
	}

	return "", errors.New("no inlineData.data found in response")
}
