package utils

import (
	"fmt"
	"strings"

	"github.com/Sanchir01/candles_backend/internal/gql/model"
)

func FormatOrders(orders []*model.Orders) string {
	var builder strings.Builder
	builder.WriteString("📦 *Список заказов:*\n\n")

	for _, order := range orders {
		builder.WriteString(fmt.Sprintf(
			"🆔 *ID:* `%s`\n📌 *Товар:* %s\n🔢 *Количество:* %d\n💰 *Цена:* %.2f\n📄 *Статус:* %s\n\n",
			order.ID,
			escapeMarkdownV2(order.Status),
			order.TotalAmount,
			order.Version,
			order.TotalAmount,
		))
	}

	return builder.String()
}
func escapeMarkdownV2(text string) string {
	specialChars := "_*[]()~`>#+-=|{}.!"
	escapedText := strings.Builder{}

	for _, char := range text {
		if strings.ContainsRune(specialChars, char) {
			escapedText.WriteRune('\\')
		}
		escapedText.WriteRune(char)
	}

	return escapedText.String()
}
