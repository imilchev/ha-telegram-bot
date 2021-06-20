package cmd

import "github.com/spf13/cobra"

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "hass-telegram-bot",
		Short:        "HASS Telegram bot.",
		SilenceUsage: true,
		// This breaks completion for 'helm help <TAB>'
		// The Cobra release following 1.0 will fix this
		//ValidArgsFunction: noCompletions, // Disable file completion
	}
	cmd.AddCommand(newRunCmd())
	return cmd
}
