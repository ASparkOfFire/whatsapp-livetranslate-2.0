package ollama

import "time"

type OllamaTranslateRequestSchema struct {
	Model    string `json:"model"`
	Messages []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
	Stream bool `json:"stream"`
	Format struct {
		Type       string `json:"type"`
		Properties struct {
			Output struct {
				Type    string `json:"type"`
				Example string `json:"example"`
			} `json:"output"`
		} `json:"properties"`
	} `json:"format"`
	Required []string `json:"required"`
}

type OllamaTranslateResponseSchema struct {
	Model     string    `json:"model"`
	CreatedAt time.Time `json:"created_at"`
	Message   struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"message"`
	Done               bool   `json:"done"`
	DoneReason         string `json:"done_reason"`
	TotalDuration      int64  `json:"total_duration"`
	LoadDuration       int    `json:"load_duration"`
	PromptEvalCount    int    `json:"prompt_eval_count"`
	PromptEvalDuration int    `json:"prompt_eval_duration"`
	EvalCount          int    `json:"eval_count"`
	EvalDuration       int64  `json:"eval_duration"`
}

type ModelOutputSchema struct {
	Output string `json:"output"`
}
