package cmdframework

import (
	"context"
	
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
	waProto "go.mau.fi/whatsmeow/proto/waE2E"
)

type Command interface {
	Execute(ctx *Context) error
	Metadata() *Metadata
}

type Metadata struct {
	Name         string
	Aliases      []string
	Description  string
	Category     string
	Usage        string
	Examples     []string
	RequireOwner bool
	Hidden       bool
	Parameters   []Parameter
}

type Parameter struct {
	Name        string
	Type        ParameterType
	Description string
	Required    bool
	Default     interface{}
	Validator   func(value string) error
}

type ParameterType int

const (
	StringParam ParameterType = iota
	IntParam
	FloatParam
	BoolParam
	DurationParam
)

type Context struct {
	context.Context
	
	// Message info
	Message     *waProto.Message
	MessageInfo types.MessageInfo
	Command     string
	Args        []string
	RawArgs     string
	
	// Services
	Handler HandlerInterface
}

type HandlerInterface interface {
	SendResponse(msgInfo types.MessageInfo, text string) error
	SendMedia(msgInfo types.MessageInfo, mediaType MediaType, data []byte, caption string) error
	SendImage(msgInfo types.MessageInfo, upload UploadResponse, caption string) error
	SendVideo(msgInfo types.MessageInfo, upload UploadResponse, caption string) error
	SendDocument(msgInfo types.MessageInfo, upload UploadResponse, caption string) error
	EditMessage(msgInfo types.MessageInfo, newText string) error
	EditMessageWithOriginal(msgInfo types.MessageInfo, newText string, originalMsg *waProto.Message) error
	GetClient() ClientInterface
	GetTranslator() TranslatorInterface
	GetImageGenerator() ImageGeneratorInterface
	GetMemeGenerator() MemeGeneratorInterface
	GetLangDetector() LangDetectorInterface
}

type MediaType int

const (
	MediaImage MediaType = iota
	MediaVideo
	MediaDocument
	MediaAudio
)

type ClientInterface interface {
	SendMessage(ctx context.Context, to types.JID, message *waProto.Message) (resp whatsmeow.SendResponse, err error)
	Upload(ctx context.Context, data []byte, appInfo MediaType) (uploadResponse UploadResponse, err error)
}

type UploadResponse struct {
	URL           string
	DirectPath    string
	MediaKey      []byte
	FileEncSHA256 []byte
	FileSHA256    []byte
	FileLength    uint64
}

type TranslatorInterface interface {
	TranslateText(ctx context.Context, text, sourceLang, targetLang string) (string, error)
	SetModel(modelID string) error
	GetModel() string
	SetTemperature(temp float64) error
	GetTemperature() float64
}

type ImageGeneratorInterface interface {
	GenerateImage(ctx context.Context, prompt string) ([]byte, error)
}

type MemeGeneratorInterface interface {
	GetRandomMeme(ctx context.Context, subreddit string) (*MemeResponse, error)
}

type MemeResponse struct {
	Memes []Meme `json:"memes"`
}

type Meme struct {
	Title     string `json:"title"`
	URL       string `json:"url"`
	Subreddit string `json:"subreddit"`
}

type LangDetectorInterface interface {
	DetectLanguage(text string) (string, error)
}