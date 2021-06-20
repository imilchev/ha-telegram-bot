package utils

import "go.uber.org/zap"

func InitLogger(enableDebug bool) error {
	var cfg zap.Config
	if enableDebug {
		cfg = zap.NewDevelopmentConfig()
	} else {
		cfg = zap.NewProductionConfig()
		cfg.Encoding = "console"
		cfg.EncoderConfig = zap.NewDevelopmentEncoderConfig()
	}

	logger, err := cfg.Build()
	zap.ReplaceGlobals(logger)

	return err
}
