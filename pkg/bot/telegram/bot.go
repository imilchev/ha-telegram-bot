package telegram

import (
	"fmt"
	"time"

	"github.com/imilchev/hass-telegram-bot/pkg/bot/telegram/model"
	"github.com/imilchev/hass-telegram-bot/pkg/config"
	"go.uber.org/zap"
	tb "gopkg.in/tucnak/telebot.v2"
)

type Bot struct {
	b            *tb.Bot
	botChan      chan struct{}
	config       config.TelegramConfig
	temperatures map[string]model.DataEntry
	humidities   map[string]model.DataEntry
}

func NewBot(config config.TelegramConfig, botChan chan struct{}) (*Bot, error) {
	b, err := tb.NewBot(tb.Settings{
		URL: "https://api.telegram.org",

		Token:  config.Token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		return nil, err
	}

	bot := &Bot{
		b:            b,
		botChan:      botChan,
		config:       config,
		temperatures: make(map[string]model.DataEntry),
		humidities:   make(map[string]model.DataEntry),
	}
	bot.registerHandles()

	return bot, nil
}

func (b *Bot) Start() {
	b.b.Start()
	zap.S().Infof("Telegram bot shut down.")
	close(b.botChan)
}

func (b *Bot) Stop() {
	b.b.Stop()
}

func (b *Bot) AddTemp(d model.DataEntry) {
	b.temperatures[d.EntityId] = d
}

func (b *Bot) AddHumid(d model.DataEntry) {
	b.humidities[d.EntityId] = d
}

func (b *Bot) registerHandles() {
	b.b.Handle("/hello", func(m *tb.Message) {
		if _, err := b.b.Send(m.Sender, "Hello World!"); err != nil {
			zap.S().Error(err)
		}
	})
	b.b.Handle("/temp", func(m *tb.Message) {
		b.getTemp(m)
	})
}

func (b *Bot) getTemp(m *tb.Message) {
	if err := b.ensureSenderAllowed(m); err != nil {
		zap.S().Warn(err)
		return
	}

	msgRows := make(map[string]string)
	for _, v := range b.temperatures {
		msgRows[v.FriendlyName] = "<b>" + v.FriendlyName + "</b> - " + v.Value + v.UnitOfMeasurement + ""
	}

	for _, v := range b.humidities {
		if _, ok := msgRows[v.FriendlyName]; ok {
			msgRows[v.FriendlyName] += " / " + v.Value + v.UnitOfMeasurement
		} else {
			msgRows[v.FriendlyName] = "<b>" + v.FriendlyName + "</b> - " + v.Value + v.UnitOfMeasurement + ""
		}
	}

	msg := ""
	keys := keysAlphabetically(msgRows)
	for _, k := range keys {
		msg += msgRows[k] + "\n"
	}
	if _, err := b.b.Send(m.Sender, msg, tb.ModeHTML); err != nil {
		zap.S().Error(err)
	}
}

func (b *Bot) ensureSenderAllowed(m *tb.Message) error {
	allowed := false
	for _, u := range b.config.AllowedUsers {
		if m.Sender.Username == u {
			allowed = true
			break
		}
	}

	if !allowed {
		if _, err := b.b.Send(m.Sender, "Fuck off!", tb.ModeHTML); err != nil {
			zap.S().Error(err)
		}
		return fmt.Errorf("User %q is not allowed to use this bot.", m.Sender.Username)
	}
	return nil
}
