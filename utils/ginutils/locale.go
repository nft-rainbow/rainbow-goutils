package ginutils

import (
	"strings"

	"golang.org/x/text/language"
)

const (
	LocaleEN = "en"
	LocaleZH = "zh"
)

var (
	zhBase, _ = language.Chinese.Base()
	enBase, _ = language.English.Base()
)

// LocaleFromAcceptLanguage resolves locale from Accept-Language.
// Supported locales: en, zh. Fallback: en.
func LocaleFromAcceptLanguage(header string) string {
	trimmed := strings.TrimSpace(header)
	if trimmed == "" {
		return LocaleEN
	}

	normalized := strings.ReplaceAll(trimmed, "_", "-")
	tags, _, _ := language.ParseAcceptLanguage(normalized)
	for _, tag := range tags {
		base, _ := tag.Base()
		switch base {
		case zhBase:
			return LocaleZH
		case enBase:
			return LocaleEN
		}
	}

	return LocaleEN
}
