package main

import (
	"github.com/imilchev/hass-telegram-bot/cmd"
)

func main() {
	rootCmd := cmd.NewRootCmd()
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
