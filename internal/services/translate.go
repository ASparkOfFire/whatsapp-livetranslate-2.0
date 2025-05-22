package services

import (
	"github.com/pemistahl/lingua-go"
)

type TranslateService interface {
	TranslateText(text string, sourceLang lingua.Language, targetLang lingua.Language) (string, error)
	SetModel(modelID string) error
	GetModel() string
}
