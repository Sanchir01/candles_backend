package botkit

import (
	"context"
	"log/slog"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ViewFunc func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error
type BotKit struct {
	bot     *tgbotapi.BotAPI
	cmdView map[string]ViewFunc
}

func NewBotKit(bot *tgbotapi.BotAPI) *BotKit {
	return &BotKit{bot: bot}
}
func (b *BotKit) handleMessage(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, message.Text)
	b.bot.Send(msg)
}
func (b *BotKit) RegisterCmdView(cmd string, view ViewFunc) {
	if b.cmdView == nil {
		b.cmdView = make(map[string]ViewFunc)
	}
	b.cmdView[cmd] = view
}
func (b *BotKit) handleCommands() {
	b.RegisterCmdView("start", ViewCmdStart())
	b.RegisterCmdView("word", ViewCmdHolloWord())
	b.RegisterCmdView("allOrders", ViewCmdHolloWord())
	b.botCommands()
}

func (b *BotKit) SendOrder(path string, chatId int64) error {
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

func (b *BotKit) botCommands() {
	commands := []tgbotapi.BotCommand{
		{Command: "start", Description: "Запустить бота"},
		{Command: "help", Description: "Список доступных команд"},
	}

	if _, err := b.bot.Request(tgbotapi.NewSetMyCommands(commands...)); err != nil {
		slog.Warn(
			"error", err.Error(),
		)
	}
}

func (b *BotKit) ViewCmdAllOrders() ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		if _, err := bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Привет")); err != nil {
			return err
		}
		return nil
	}
}

// todo:delete for production
func (b *BotKit) SendMessage(chatId int64, msg string) error {
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
