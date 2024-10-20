package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleMessage(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, message.Text)
	b.bot.Send(msg)
}
func (b *Bot) handleCommands() {
	b.RegisterCmdView("start", ViewCmdStart())
}

func ViewCmdStart() ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		if _, err := bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Привет")); err != nil {
			return err
		}
		return nil
	}
}
