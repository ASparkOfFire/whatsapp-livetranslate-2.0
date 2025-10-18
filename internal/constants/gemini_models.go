package constants

type GeminiModel string

const (
	Gemini20Flash     GeminiModel = "gemini-2.0-flash"
	Gemini20FlashLite GeminiModel = "gemini-2.0-flash-lite"
	Gemini25Flash     GeminiModel = "gemini-2.5-flash"
	Gemini25FlashLite GeminiModel = "gemini-2.5-flash-lite"

	// Temperature constants
	MinTemperature     float64 = 0.0
	MaxTemperature     float64 = 1.0
	DefaultTemperature float64 = 0.2
)

const (
	GeminiModelImageGenerator GeminiModel = "gemini-2.5-flash-image"
)

var ValidGeminiModels = map[GeminiModel]bool{
	Gemini20Flash:     true,
	Gemini20FlashLite: true,
	Gemini25Flash:     true,
	Gemini25FlashLite: true,
}
