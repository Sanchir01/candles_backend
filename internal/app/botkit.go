package app

import (
	"context"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
)

type UserState struct {
	orderId sync.Map
}

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

func (b *Bot) ViewCmdHelp() ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		chatId := update.FromChat().ID
		messageText := "Можете написать в тех поддержку: \n или позвонить нам \n" + "[@bayadigital(https://t.me/bayadigital)\n" + "41252"
		msg := tgbotapi.NewMessage(chatId, messageText)
		msg.ParseMode = "Markdown"
		if _, err := bot.Send(msg); err != nil {
			slog.Error("cmd help error", err.Error())
			return err
		}
		return nil
	}
}

func (b *Bot) ViewCmdOrderStatus() ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		arg := strings.TrimSpace(strings.TrimPrefix(update.Message.Text, "/order_status"))
		if arg == "" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Пожалуйста, укажите ID вашего заказа после команды /order_status.")
			if _, err := b.bot.Send(msg); err != nil {
				return err
			}
		}
		orderId, err := uuid.Parse(arg)
		if err != nil {
			slog.Error("error parse tg order id", err.Error())
			return err
		}
		slog.Warn("this order id rg msg", orderId)
		if err := b.OrderBot.SendStatusOrder(ctx, orderId, update.FromChat().ID); err != nil {
			return err
		}
		return nil
	}
}
