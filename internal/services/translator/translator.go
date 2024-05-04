package services

import (
	"github.com/WildEgor/e-shop-support-bot/internal/configs"
	"github.com/WildEgor/e-shop-support-bot/internal/models"
	"github.com/kataras/i18n"
	"os"
	"reflect"
)

type TranslatorService struct {
	i18n   *i18n.I18n
	Prefix string
}

func NewTranslatorService(cfg *configs.TranslatorConfig) (*TranslatorService, error) {
	dirDirs, err := os.ReadDir(cfg.LocalesDir())
	if err != nil {
		return nil, err
	}

	dirs := make([]string, 0)
	dirs = append(dirs, cfg.DefaultLocale)

	for _, dir := range dirDirs {
		if !dir.IsDir() {
			continue
		}
		if cfg.DefaultLocale == dir.Name() {
			continue
		}

		dirs = append(dirs, dir.Name())
	}

	i, err := i18n.New(
		i18n.Glob(
			cfg.LocalesFullPath(), i18n.LoaderConfig{},
		), dirs...,
	)
	if err != nil {
		return nil, err
	}

	return &TranslatorService{
		i18n:   i,
		Prefix: cfg.Prefix,
	}, nil
}

func (t *TranslatorService) GetLocalizedMessage(lang string, dictKey models.MessageKey, args ...interface{}) string {
	if len(args) > 0 {
		in := t.parsePayload(args[0])
		return t.Prefix + " " + t.i18n.Tr(lang, string(dictKey), in)
	}

	return t.Prefix + " " + t.i18n.Tr(lang, string(dictKey))
}

func (t *TranslatorService) GetMessageWithoutPrefix(lang string, dictKey models.MessageKey, args ...interface{}) string {
	if len(args) > 0 {
		in := t.parsePayload(args[0])
		return t.i18n.Tr(lang, string(dictKey), in)
	}

	return t.i18n.Tr(lang, string(dictKey))
}

func (t *TranslatorService) parsePayload(args interface{}) map[string]interface{} {
	out := make(map[string]interface{})

	v := reflect.ValueOf(args)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return out
	}

	typ := v.Type()
	for i := 0; i < typ.NumField(); i++ {
		// gets us a StructField
		fi := typ.Field(i)
		if tagv := fi.Tag.Get("template"); tagv != "" {
			out[tagv] = v.Field(i).Interface()
		}
	}

	return out
}
