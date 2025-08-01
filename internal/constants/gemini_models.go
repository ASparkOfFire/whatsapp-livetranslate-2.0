package constants

type GeminiModel string

const (
	GeminiModel15Flash GeminiModel = "gemini-1.5-flash"
	GeminiModel20Flash GeminiModel = "gemini-2.0-flash"
	GeminiModel25Flash GeminiModel = "gemini-2.5-flash-preview-05-20"

	// Temperature constants
	MinTemperature     float64 = 0.0
	MaxTemperature     float64 = 1.0
	DefaultTemperature float64 = 0.2
)

const (
	GeminiModelImageGenerator GeminiModel = "gemini-2.0-flash-preview-image-generation"
)

var ValidGeminiModels = map[GeminiModel]bool{
	GeminiModel15Flash: true,
	GeminiModel20Flash: true,
	GeminiModel25Flash: true,
}
