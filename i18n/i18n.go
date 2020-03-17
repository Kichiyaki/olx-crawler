package i18n

import (
	"os"
	"path/filepath"

	_i18n "github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var lang = language.Polish
var bundle = _i18n.NewBundle(lang)

func LoadMessageFiles(root string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if path != root {
			bundle.MustLoadMessageFile(path)
		}
		return nil
	})
}

func SetLanguage(l string) error {
	t, err := language.Parse(l)
	if err != nil {
		return err
	}
	lang = t
	return nil
}

func NewLocalizer() *_i18n.Localizer {
	return _i18n.NewLocalizer(bundle, lang.String())
}
