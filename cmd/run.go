package cmd

import (
	"fmt"

	"github.com/imilchev/hass-telegram-bot/pkg/bot"
	"github.com/imilchev/hass-telegram-bot/pkg/config"
	"github.com/imilchev/hass-telegram-bot/pkg/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var debug bool

func newRunCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "run [configPath]",
		Short:        "Run the HA-Telegram bot.",
		SilenceUsage: true,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("It is mandatory to provide the config path.")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := utils.InitLogger(debug); err != nil {
				panic(err)
			}
			defer func() {
				if err := zap.S().Sync(); err != nil {
					zap.S().Info()
				}
			}()

			cfg, err := config.ReadConfig(args[0])
			if err != nil {
				return err
			}

			mgr, err := bot.NewBotManager(*cfg)
			if err != nil {
				return err
			}
			return mgr.Start()
		},
	}
	cmd.Flags().BoolVar(
		&debug,
		"debug",
		false,
		"Enable debug logging.")
	return cmd
}
