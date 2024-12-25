package order

import (
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type OrderBot struct {
	bot          *tgbotapi.BotAPI
	OrderService *Service
}

func NewBotService(OrderService *Service, bot *tgbotapi.BotAPI) *OrderBot {
	return &OrderBot{
		OrderService: OrderService,
		bot:          bot,
	}
}

func (b *OrderBot) SendOrder(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	defer file.Close()
	chatId, err := strconv.ParseInt(os.Getenv("TELEGRAM_CHATID"), 10, 64)

	if err != nil {
		return err
	}
	doc := tgbotapi.NewDocument(chatId, tgbotapi.FileReader{
		Name:   "Order.xlsx",
		Reader: file,
	})

	if _, err := b.bot.Send(doc); err != nil {
		return err
	}
	return nil
}
