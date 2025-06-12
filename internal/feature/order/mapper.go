package order

import orderv1 "github.com/Sanchir01/auth-proto/gen/go/order"

func MapProductsToOrderRequest(products []ProductWithQuantity) *orderv1.OrderRequest {
	orderRequest := &orderv1.OrderRequest{
		Orders: make([]*orderv1.OrderData, 0, len(products)),
	}

	for _, product := range products {
		orderData := &orderv1.OrderData{
			Orders: &orderv1.OrderItem{
				Id:    product.Candle.ID.String(),
				Title: product.Candle.Title,
				Slug:  product.Candle.Slug,
				Price: int32(product.Candle.Price),
			},
			Quantity: int32(product.Quantity),
		}
		orderRequest.Orders = append(orderRequest.Orders, orderData)
	}

	return orderRequest
}
