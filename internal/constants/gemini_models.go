package constants

type GeminiModel string

const (
	GeminiModel15Flash GeminiModel = "gemini-1.5-flash"
	GeminiModel20Flash GeminiModel = "gemini-2.0-flash"
	GeminiModel25Flash GeminiModel = "gemini-2.5-flash-preview-05-20"
)

var ValidGeminiModels = map[GeminiModel]bool{
	GeminiModel15Flash: true,
	GeminiModel20Flash: true,
	GeminiModel25Flash: true,
}
