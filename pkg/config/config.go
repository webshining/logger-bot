package config

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

var (
	k        = koanf.New(".")
	validate = validator.New()
)

func Load(cfg interface{}, defaultConfig map[string]interface{}) error {

	k.Load(confmap.Provider(defaultConfig, "_"), nil)
	k.Load(file.Provider(".env"), dotenv.ParserEnv("", "_", func(k string) string {
		return k
	}))
	if err := errors.Join(
		k.Unmarshal("", cfg),
		validate.Struct(cfg),
	); err != nil {
		return err
	}

	return nil
}

func MustLoad(cfg interface{}, defaultConfig map[string]interface{}, logger *zap.Logger) {
	if err := Load(cfg, defaultConfig); err != nil {
		logger.Named("Config").Fatal("failed to load config", zap.Error(err))
	}
}
