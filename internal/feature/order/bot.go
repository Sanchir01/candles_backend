package order

import (
	"context"
	"os"
	"strconv"

	"github.com/Sanchir01/candles_backend/pkg/lib/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
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
func (b *OrderBot) SendAllordersTg(ctx context.Context, chatId int64) error {
	orders, err := b.OrderService.AllOrders(ctx)
	if err != nil {
		return err
	}
	orderText := utils.FormatOrders(orders)
	msg := tgbotapi.NewMessage(chatId, orderText)
	if _, err := b.bot.Send(msg); err != nil {
		return err
	}
	return nil
}

func (b *OrderBot) SendStatusOrder(ctx context.Context, productId uuid.UUID, chatId int64)
