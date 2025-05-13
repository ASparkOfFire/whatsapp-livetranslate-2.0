package services

import (
	"github.com/asparkoffire/whatsapp-livetranslate-go/internal/constants"
	"github.com/pemistahl/lingua-go"
)

type LangDetectService interface {
	DetectLanguage(text string) (lingua.Language, bool)
}

type linguaLangDetectService struct {
	detector lingua.LanguageDetector
}

func NewLinguaLangDetectService(supportedLanguages map[string]lingua.Language) LangDetectService {
	var langs []lingua.Language
	for _, language := range constants.SupportedLanguages {
		langs = append(langs, language)
	}

	detector := lingua.NewLanguageDetectorBuilder().
		FromLanguages(langs...).Build()

	return &linguaLangDetectService{
		detector: detector,
	}
}

func (s *linguaLangDetectService) DetectLanguage(text string) (lingua.Language, bool) {
	return s.detector.DetectLanguageOf(text)
}
