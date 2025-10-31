package middleware

import (
	"encoding/json"
	"go_blog/pkg/config"
	"go_blog/pkg/utils"
	"path"

	ginI18n "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

func I18nMiddleware(conf *config.Config) gin.HandlerFunc {
	// default language
	i18nDefaultLanguage := language.Make(conf.I18n.DefaultLanguage)
	// accept languages
	i18nAcceptLanguages := make([]language.Tag, len(conf.I18n.AcceptLanguages))
	for i, lang := range conf.I18n.AcceptLanguages {
		i18nAcceptLanguages[i] = language.Make(lang)
	}
	// load language files
	rootPath := path.Join(utils.GetRootDir(), conf.I18n.LanguageFilePath)
	var unmarshalFunc i18n.UnmarshalFunc
	if conf.I18n.LanguageFileFormat == "json" {
		unmarshalFunc = json.Unmarshal
	} else if conf.I18n.LanguageFileFormat == "yaml" {
		unmarshalFunc = yaml.Unmarshal
	}
	return ginI18n.Localize(ginI18n.WithBundle(&ginI18n.BundleCfg{
		RootPath:         rootPath,
		AcceptLanguage:   i18nAcceptLanguages,
		DefaultLanguage:  i18nDefaultLanguage,
		UnmarshalFunc:    unmarshalFunc,
		FormatBundleFile: conf.I18n.LanguageFileFormat,
	}))
}
