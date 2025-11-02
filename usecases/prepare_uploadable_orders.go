package usecases

import (
	"mstorefgo/internal/moyskladapi"
	"mstorefgo/internal/order_processor"
)

func PrepareUploadableOrders(mp *moyskladapi.MoySkladProcessor) (count int, storage *map[string]order_processor.ProcessedOrder, err error) {
	count, storage, err = order_processor.ProcessOrders(mp)
	if err != nil {
		return count, storage, err
	}
	return count, storage, nil
}
