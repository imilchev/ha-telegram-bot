package bot

import (
	"os"
	"os/signal"

	"github.com/imilchev/hass-telegram-bot/pkg/bot/ha"
	"github.com/imilchev/hass-telegram-bot/pkg/bot/telegram"
	"github.com/imilchev/hass-telegram-bot/pkg/bot/telegram/model"
	"github.com/imilchev/hass-telegram-bot/pkg/config"
	"go.uber.org/zap"
)

type BotManager struct {
	config   config.Config
	haClient *ha.WsClient
	bot      *telegram.Bot
	botChan  chan struct{}
}

func NewBotManager(config config.Config) (*BotManager, error) {
	wsClient, err := ha.NewWsClient(config.HomeAssistant)
	if err != nil {
		return nil, err
	}

	botChan := make(chan struct{})
	bot, err := telegram.NewBot(config.Telegram, botChan)
	if err != nil {
		return nil, err
	}

	return &BotManager{
		config:   config,
		haClient: wsClient,
		bot:      bot,
		botChan:  botChan,
	}, nil
}

func (m *BotManager) Start() error {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	go func() {
		if err := m.haClient.Start(); err != nil {
			zap.S().Error(err)
		}
	}()

	go m.bot.Start()

	for {
		select {
		case msg := <-m.haClient.TempChan:
			// TODO: filter only relevant messages and store the results
			m.bot.AddTemp(model.DataEntry{
				EntityId:          msg.Event.Data.EntityId,
				FriendlyName:      msg.Event.Data.NewState.Attributes.FriendlyName,
				Value:             msg.Event.Data.NewState.State,
				UnitOfMeasurement: msg.Event.Data.NewState.Attributes.UnitOfMeasurement,
			})
			zap.S().Debug(msg)
		case msg := <-m.haClient.HumidChan:
			// TODO: filter only relevant messages and store the results
			m.bot.AddHumid(model.DataEntry{
				EntityId:          msg.Event.Data.EntityId,
				FriendlyName:      msg.Event.Data.NewState.Attributes.FriendlyName,
				Value:             msg.Event.Data.NewState.State,
				UnitOfMeasurement: msg.Event.Data.NewState.Attributes.UnitOfMeasurement,
			})
			zap.S().Debug(msg)
		case <-interrupt:
			zap.S().Info("Shutting down...")
			if err := m.haClient.Stop(); err != nil {
				return err
			}
			m.bot.Stop()
			<-m.haClient.TempChan  // Wait until eventChan is closed by haClient
			<-m.haClient.HumidChan // Wait until botChan is closed by bot
			zap.S().Info("Exit")
			return nil

		}
	}
}
