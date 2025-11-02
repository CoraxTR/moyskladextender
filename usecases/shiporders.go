package usecases

import (
	"mstorefgo/internal/moyskladapi"
	"mstorefgo/internal/order_processor"
)

func ShipOrders(m *moyskladapi.MoySkladProcessor, orders map[string]order_processor.ProcessedOrder) {
	for _, order := range orders {
		if order.PaymentMethod == "Наличные" || order.PaymentMethod == "Терминал" {
			m.ShipOrder(order.HREF, order.Counterpartyhref)
		}
	}
}
