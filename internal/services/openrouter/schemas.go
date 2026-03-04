package openrouter

type OpenrouterTranslateRequestSchema struct {
	Model       string  `json:"model"`
	Stream      bool    `json:"stream"`
	Temperature float64 `json:"temperature"`
	Messages    []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
	ResponseFormat struct {
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
	} `json:"response_format"`
}

type OpenrouterTranslateResponseSchema struct {
	ID       string `json:"id"`
	Object   string `json:"object"`
	Created  int    `json:"created"`
	Model    string `json:"model"`
	Provider string `json:"provider"`
	Choices  []struct {
		Index              int    `json:"index"`
		Logprobs           any    `json:"logprobs"`
		FinishReason       string `json:"finish_reason"`
		NativeFinishReason string `json:"native_finish_reason"`
		Message            struct {
			Role      string `json:"role"`
			Content   string `json:"content"`
			Refusal   any    `json:"refusal"`
			Reasoning any    `json:"reasoning"`
		} `json:"message"`
	} `json:"choices"`
	Usage struct {
		PromptTokens        int  `json:"prompt_tokens"`
		CompletionTokens    int  `json:"completion_tokens"`
		TotalTokens         int  `json:"total_tokens"`
		Cost                int  `json:"cost"`
		IsByok              bool `json:"is_byok"`
		PromptTokensDetails struct {
			CachedTokens     int `json:"cached_tokens"`
			CacheWriteTokens int `json:"cache_write_tokens"`
			AudioTokens      int `json:"audio_tokens"`
			VideoTokens      int `json:"video_tokens"`
		} `json:"prompt_tokens_details"`
		CostDetails struct {
			UpstreamInferenceCost            any `json:"upstream_inference_cost"`
			UpstreamInferencePromptCost      int `json:"upstream_inference_prompt_cost"`
			UpstreamInferenceCompletionsCost int `json:"upstream_inference_completions_cost"`
		} `json:"cost_details"`
		CompletionTokensDetails struct {
			ReasoningTokens int `json:"reasoning_tokens"`
			ImageTokens     int `json:"image_tokens"`
			AudioTokens     int `json:"audio_tokens"`
		} `json:"completion_tokens_details"`
	} `json:"usage"`
}

type ModelOutputSchema struct {
	Output string `json:"output"`
}

type OpenrouterImageGenerationRequestSchema struct {
	Model    string `json:"model"`
	Messages []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
}

type OpenrouterImageGenerationResponseSchema struct {
	ID       string `json:"id"`
	Object   string `json:"object"`
	Created  int    `json:"created"`
	Model    string `json:"model"`
	Provider string `json:"provider"`
	Choices  []struct {
		Index              int    `json:"index"`
		Logprobs           any    `json:"logprobs"`
		FinishReason       string `json:"finish_reason"`
		NativeFinishReason any    `json:"native_finish_reason"`
		Message            struct {
			Role      string `json:"role"`
			Content   any    `json:"content"`
			Refusal   any    `json:"refusal"`
			Reasoning any    `json:"reasoning"`
			Images    []struct {
				Type     string `json:"type"`
				ImageURL struct {
					URL string `json:"url"`
				} `json:"image_url"`
			} `json:"images"`
		} `json:"message"`
	} `json:"choices"`
	Usage struct {
		PromptTokens        int     `json:"prompt_tokens"`
		CompletionTokens    int     `json:"completion_tokens"`
		TotalTokens         int     `json:"total_tokens"`
		Cost                float64 `json:"cost"`
		IsByok              bool    `json:"is_byok"`
		PromptTokensDetails struct {
			CachedTokens     int `json:"cached_tokens"`
			CacheWriteTokens int `json:"cache_write_tokens"`
			AudioTokens      int `json:"audio_tokens"`
			VideoTokens      int `json:"video_tokens"`
		} `json:"prompt_tokens_details"`
		CostDetails struct {
			UpstreamInferenceCost            float64 `json:"upstream_inference_cost"`
			UpstreamInferencePromptCost      int     `json:"upstream_inference_prompt_cost"`
			UpstreamInferenceCompletionsCost float64 `json:"upstream_inference_completions_cost"`
		} `json:"cost_details"`
		CompletionTokensDetails struct {
			ReasoningTokens int `json:"reasoning_tokens"`
			ImageTokens     int `json:"image_tokens"`
			AudioTokens     int `json:"audio_tokens"`
		} `json:"completion_tokens_details"`
	} `json:"usage"`
}
