package utils

import (
	"strings"

	"github.com/asparkoffire/whatsapp-livetranslate-go/internal/constants"
	"github.com/pemistahl/lingua-go"
)

func GetLangByCode(code string) lingua.Language {
	code = strings.ToLower(code)
	if lang, exists := constants.SupportedLanguages[code]; !exists {
		return lingua.Unknown
	} else {
		return lang
	}
}
