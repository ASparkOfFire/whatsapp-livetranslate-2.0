package memegenerator

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type MemeResponse struct {
	Count int `json:"count"`
	Memes []struct {
		PostLink  string   `json:"postLink"`
		Subreddit string   `json:"subreddit"`
		Title     string   `json:"title"`
		URL       string   `json:"url"`
		NSFW      bool     `json:"nsfw"`
		Spoiler   bool     `json:"spoiler"`
		Author    string   `json:"author"`
		Ups       int      `json:"ups"`
		Preview   []string `json:"preview"`
	} `json:"memes"`
}

type MemeGenerator struct {
	client *http.Client
}

func NewMemeGenerator() *MemeGenerator {
	return &MemeGenerator{
		client: &http.Client{},
	}
}

func (g *MemeGenerator) GetRandomMeme(ctx context.Context, subreddit string) (*MemeResponse, error) {
	url := "https://meme-api.com/gimme/1"
	if subreddit != "" {
		url = fmt.Sprintf("https://meme-api.com/gimme/%s/1", subreddit)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error response from meme API: %s - %s", resp.Status, string(body))
	}

	var memeResp MemeResponse
	if err := json.NewDecoder(resp.Body).Decode(&memeResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	if len(memeResp.Memes) == 0 {
		return nil, fmt.Errorf("no memes found")
	}

	return &memeResp, nil
}
