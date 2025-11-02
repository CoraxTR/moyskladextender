package usecases

import (
	"mstorefgo/internal/moyskladapi"
	"mstorefgo/internal/order_processor"
)

func ChangeStatusToShiped(ms *moyskladapi.MoySkladProcessor, orders map[string]order_processor.ProcessedOrder) {
	for _, order := range orders {
		ms.SetOrderShipped(order.HREF)
	}
}
