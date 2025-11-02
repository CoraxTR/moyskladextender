package usecases

import (
	"fmt"
	"mstorefgo/internal/moyskladapi"
	"mstorefgo/internal/order_processor"
	"strconv"
)

func ShipOrders(m *moyskladapi.MoySkladProcessor, orders map[string]order_processor.ProcessedOrder) {
	for _, order := range orders {

		refnumber, err := strconv.Atoi(order.RefGoNumber)
		fmt.Printf("refnumber is %v\n", refnumber)
		if err != nil {
			fmt.Println(err)
		}
		m.SetOrderSellTypetoOther(order.HREF)
		m.SetOrderRefGoNumber(order.HREF, refnumber)
		/*	if order.PaymentMethod == "Наличные" || order.PaymentMethod == "Терминал" {
				m.ShipOrder(order.HREF, order.Counterpartyhref)
			}
		*/
	}
}
