package telegrambot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"runtime/debug"
)

func (b *Bot) handleUpdate(updates tgbotapi.UpdatesChannel) {
	defer func() {
		if p := recover(); p != nil {
			log.Printf("panic recoverL: %v\n%s", p, string(debug.Stack()))
		}
	}()
	for update := range updates {
		if update.Message != nil {
			b.handleMessage(update.Message)
		}
	}
}

func (b *Bot) handleMessage(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, message.Text)
	b.bot.Send(msg)
}
