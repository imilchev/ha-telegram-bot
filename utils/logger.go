package utils

import "go.uber.org/zap"

func InitLogger() {
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeCaller = nil
	logger, _ := cfg.Build()
	zap.ReplaceGlobals(logger)
}
