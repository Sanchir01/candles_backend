package app

import (
	"context"
	"log"
	"log/slog"
	"os"
	"runtime/debug"
	"time"

	"github.com/Sanchir01/candles_backend/internal/feature/order"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	bot      *tgbotapi.BotAPI
	cmdView  map[string]ViewFunc
	OrderBot *order.OrderBot
}

func NewBot(Services *Services) *Bot {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}
	return &Bot{
		bot:      bot,
		OrderBot: order.NewBotService(Services.OrderService, bot),
	}
}

func (b *Bot) Start(ctx context.Context) error {
	b.bot.Debug = true
	updates, err := b.initUpdatesChannel()
	if err != nil {
		return err
	}
	b.HandleCommands()
	b.botCommands()
	for {
		select {
		case update := <-updates:
			if !update.Message.IsCommand() {
				continue
			}
			updateCtx, updateChannel := context.WithTimeout(ctx, 5*time.Second)
			b.handleUpdateCommand(updateCtx, update)
			updateChannel()
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (b *Bot) handleUpdateCommand(ctx context.Context, update tgbotapi.Update) {
	defer func() {
		if p := recover(); p != nil {
			log.Printf("panic recoverL: %v\n%s", p, string(debug.Stack()))
		}
	}()
	var view ViewFunc

	cmd := update.Message.Command()
	cmdView, ok := b.cmdView[cmd]
	if !ok {
		if _, err := b.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "ошибка выполнения команды")); err != nil {
			log.Println(err)
		}
		return
	}

	view = cmdView
	if err := view(ctx, b.bot, update); err != nil {
		slog.Error(err.Error())
		if _, err := b.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "ошибка выполнения команды")); err != nil {
			log.Println(err)
		}
	}
}

func (b *Bot) initUpdatesChannel() (tgbotapi.UpdatesChannel, error) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	return b.bot.GetUpdatesChan(u), nil
}

type ViewFunc func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error

func (b *Bot) handleMessage(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, message.Text)
	b.bot.Send(msg)
}

func (b *Bot) RegisterCmdView(cmd string, view ViewFunc) {
	if b.cmdView == nil {
		b.cmdView = make(map[string]ViewFunc)
	}
	b.cmdView[cmd] = view
}

func (b *Bot) HandleCommands() {
	b.RegisterCmdView("start", b.ViewCmdStart())
	b.RegisterCmdView("orders", b.ViewCmdAllOrders())
	b.RegisterCmdView("help", b.ViewCmdHelp())
	b.RegisterCmdView("order_status", b.ViewCmdOrderStatus())
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

func (b *Bot) botCommands() {
	commands := []tgbotapi.BotCommand{
		{Command: "start", Description: "Запустить бота"},
		{Command: "help", Description: "Список доступных команд"},
		{Command: "orders", Description: "Список всех заказов"},
		{Command: "order_status", Description: "Получить статус заказа"},
	}

	if _, err := b.bot.Request(tgbotapi.NewSetMyCommands(commands...)); err != nil {
		slog.Warn(
			"error", err.Error(),
		)
	}
}

// todo:delete for production
func (b *Bot) SendMessage(chatId int64, msg string) error {
	message := tgbotapi.NewMessage(chatId, msg)
	if _, err := b.bot.Send(message); err != nil {
		return err
	}
	return nil
}
