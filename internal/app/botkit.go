package app

import (
	"context"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) ViewCmdStart() ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		if _, err := bot.Send(tgbotapi.NewMessage(update.FromChat().ID,
			"Привет, это бот магазина махакала тут вы можете посмотреть свои заказы и связаться с нами")); err != nil {
			return err
		}
		return nil
	}
}
func (b *Bot) ViewCmdAllOrders() ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		chatId, err := strconv.ParseInt(os.Getenv("TELEGRAM_CHATID"), 10, 64)
		if err != nil {
			return err
		}
		if err := b.OrderBot.SendAllordersTg(ctx, chatId); err != nil {
			return err
		}
		return nil
	}
}
