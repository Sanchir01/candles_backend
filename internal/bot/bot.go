package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"log/slog"
	"os"
	"runtime/debug"
	"time"
)

type Bot struct {
	bot     *tgbotapi.BotAPI
	cmdView map[string]ViewFunc
}

type ViewFunc func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error

func NewBot() *Bot {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}
	return &Bot{bot: bot}
}

func (b *Bot) Start(ctx context.Context) error {
	b.bot.Debug = true
	updates, err := b.initUpdatesChannel()
	if err != nil {
		return err
	}
	b.handleCommands()

	for {
		select {
		case update := <-updates:
			if !update.Message.IsCommand() {
				continue
			}
			updateCtx, updateCanndel := context.WithTimeout(ctx, 5*time.Second)
			b.handleUpdateCommand(updateCtx, update)
			updateCanndel()
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (b *Bot) RegisterCmdView(cmd string, view ViewFunc) {
	if b.cmdView == nil {
		b.cmdView = make(map[string]ViewFunc)
	}
	b.cmdView[cmd] = view
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
