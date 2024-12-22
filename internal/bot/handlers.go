package bot

import (
	"context"
	"log/slog"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleMessage(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, message.Text)
	b.bot.Send(msg)
}
func (b *Bot) handleCommands() {
	b.RegisterCmdView("start", ViewCmdStart())
	b.RegisterCmdView("word", ViewCmdHolloWord())
}

func (b *Bot) SendOrder(path string, chatId int64) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	document := tgbotapi.NewDocument(chatId, tgbotapi.FileReader{Name: "products.xlsx", Reader: file})
	if _, err := b.bot.Send(document); err != nil {
		return err
	}
	slog.Warn("Документ успешно отправлен!")
	return nil
}

// todo:delete for production
func (b *Bot) SendMessage(chatId int64, msg string) error {
	message := tgbotapi.NewMessage(chatId, msg)
	if _, err := b.bot.Send(message); err != nil {
		return err
	}
	return nil
}

func ViewCmdHolloWord() ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		if _, err := bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Hollo Word")); err != nil {
			return err
		}
		return nil
	}
}
func ViewCmdStart() ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		if _, err := bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Привет")); err != nil {
			return err
		}
		return nil
	}
}
