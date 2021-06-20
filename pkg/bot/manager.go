package bot

import (
	"github.com/imilchev/hass-telegram-bot/pkg/bot/ha"
	"github.com/imilchev/hass-telegram-bot/pkg/config"
)

type BotManager struct {
	config   config.Config
	haClient *ha.WsClient
}

func NewBotManager(config config.Config) (*BotManager, error) {
	wsClient, err := ha.NewWsClient(config.HomeAssistant)
	if err != nil {
		return nil, err
	}
	return &BotManager{config: config, haClient: wsClient}, nil
}

func (m *BotManager) Start() error {
	return m.haClient.Start()
}
