package telegrambot

import (
	"context"
	"github.com/Sanchir01/candles_backend/internal/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
)

type Bot struct {
	bot *tgbotapi.BotAPI
	lg  *slog.Logger
}

type ViewFunc func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error

func New(bot *tgbotapi.BotAPI, lg *slog.Logger) *Bot {
	return &Bot{bot: bot, lg: lg}
}

func (b *Bot) Start(cfg *config.Config) error {
	b.bot.Debug = true
	updates, err := b.initUpdatesChannel()
	if err != nil {
		return err
	}
	b.handleUpdate(updates)
	return nil
}

func (b *Bot) initUpdatesChannel() (tgbotapi.UpdatesChannel, error) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	return b.bot.GetUpdatesChan(u), nil
}
