package bot

import (
	"os"
	"os/signal"

	"github.com/imilchev/hass-telegram-bot/pkg/bot/ha"
	"github.com/imilchev/hass-telegram-bot/pkg/bot/ha/model"
	"github.com/imilchev/hass-telegram-bot/pkg/config"
	"go.uber.org/zap"
)

type BotManager struct {
	config    config.Config
	haClient  *ha.WsClient
	eventChan chan model.EventMessage
}

func NewBotManager(config config.Config) (*BotManager, error) {
	eventChan := make(chan model.EventMessage, 100)
	wsClient, err := ha.NewWsClient(config.HomeAssistant, eventChan)
	if err != nil {
		return nil, err
	}
	return &BotManager{config: config, haClient: wsClient, eventChan: eventChan}, nil
}

func (m *BotManager) Start() error {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	go func() {
		if err := m.haClient.Start(); err != nil {
			zap.S().Error(err)
		}
	}()

	for {
		select {
		case msg := <-m.eventChan:
			// TODO: filter only relevant messages and store the results
			zap.S().Debug(msg)
		case <-interrupt:
			zap.S().Info("Shutting down...")
			if err := m.haClient.Stop(); err != nil {
				return err
			}
			<-m.eventChan // Wait until evenChan is closed by haClient
			zap.S().Info("Exit")
			return nil

		}
	}
}
