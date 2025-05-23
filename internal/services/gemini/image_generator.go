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
	"strings"
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
	fmt.Printf("Starting image generation with model %s\n", g.modelID)

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
	fmt.Printf("Sending request to %s\n", url)

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

	fmt.Println("Received response from Gemini API, extracting image data...")

	// ---------- stream & extract image ----------
	imageData, err := extractImageData(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("extract image data: %w", err)
	}

	fmt.Println("Successfully extracted image data, decoding base64...")
	decoded, err := base64.StdEncoding.DecodeString(imageData)
	if err != nil {
		return nil, fmt.Errorf("base64-decode image: %w", err)
	}

	fmt.Printf("Successfully decoded image (%d bytes)\n", len(decoded))
	return decoded, nil
}

// -----------------------------------------------------------------------------
// helpers
// -----------------------------------------------------------------------------

// extractImageData reads the JSON stream and returns the first inlineData.data.
func extractImageData(r io.Reader) (string, error) {
	decoder := json.NewDecoder(r)
	var imageData strings.Builder

	// First decode the array of events
	var events []event
	if err := decoder.Decode(&events); err != nil {
		return "", fmt.Errorf("decode events array: %w", err)
	}

	fmt.Printf("Received %d events from API\n", len(events))

	// Process each event in the array
	for _, event := range events {
		for _, cand := range event.Candidates {
			for _, p := range cand.Content.Parts {
				if p.InlineData != nil && p.InlineData.Data != "" {
					// Append the data chunk
					imageData.WriteString(p.InlineData.Data)
					fmt.Printf("Received image data chunk (%d bytes)\n", len(p.InlineData.Data))
				}
			}
		}
	}

	finalData := imageData.String()
	if finalData == "" {
		return "", errors.New("no inlineData.data found in response")
	}

	fmt.Printf("Total image data size: %d bytes\n", len(finalData))
	return finalData, nil
}
