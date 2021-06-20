package main

import (
	"github.com/imilchev/hass-telegram-bot/cmd"
	"github.com/imilchev/hass-telegram-bot/utils"
	"go.uber.org/zap"
)

func main() {
	utils.InitLogger()
	defer func() {
		if err := zap.S().Sync(); err != nil {
			panic(err)
		}
	}()

	rootCmd := cmd.NewRootCmd()
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
