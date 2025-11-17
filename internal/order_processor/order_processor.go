package order_processor

import (
	"fmt"
	"math"
	"mstorefgo/internal/config"
	"mstorefgo/internal/moyskladapi"
	"mstorefgo/internal/unmarshaller"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type ProcessedOrder struct {
	HREF                  string
	name                  string
	Counterpartyhref      string
	RecieverName          string
	Description           string
	DeliveryPlannedDate   string
	ShipmentAddress       string
	DeliveryIntervalFrom  string
	DeliveryIntervalUntil string
	DeliveryRegion        string
	PaymentMethod         string
	RecieverPhoneNumber   int
	RefGoNumber           string
	chilledPositions      []Position
	frozenPositions       []Position
	servicePositions      []Position
	Sum                   float64
	ChilledWeight         float64
	FrozenWeight          float64
	FrozenBoxes           uint8
	ChilledBoxes          uint8
	errors                []error
}

func (u *ProcessedOrder) setOrderRefgoNumber(i int) {
	u.RefGoNumber = strconv.Itoa(i)
}

// Set Order HREF, no need to check for nil as it is automatically assigned by MoySklad
func (u *ProcessedOrder) setOrderHREF(o unmarshaller.Order) {
	u.HREF = o.Meta.HREF
}

// Set Order Name aka Number, no need to check for nil, as it is set automatically by MoySklad
func (u *ProcessedOrder) setOrderName(o unmarshaller.Order) {
	u.name = o.Name
}

// Set Order Description, no default, ok if nil
func (u *ProcessedOrder) setOrderDescription(o unmarshaller.Order) {
	u.Description = o.Description
}

// Set Order Delivery Planned Date, we trim the date to yyyy-mm-dd, and we don't need to check if the
// date is in the past, as we fetch only for 2 upcoming days
func (u *ProcessedOrder) setOrderDeliveryPlannedDate(o unmarshaller.Order) {
	dateLayout := "02.01.2006"
	tempDate, err := time.Parse("2006-01-02 15:04:05", o.DeliveryPlannedMoment)
	if err != nil {
		return
	}
	u.DeliveryPlannedDate = tempDate.Format(dateLayout)
}

func (u *ProcessedOrder) setOrderCounterpartyhref(o unmarshaller.Order) {
	u.Counterpartyhref = o.Agent.Meta.HREF
}

func processPhoneNumber(phone string) (int, error) {
	re := regexp.MustCompile(`[^\d]`)
	digitsOnly := re.ReplaceAllString(phone, "")

	if len(digitsOnly) != 11 {
		return 0, fmt.Errorf("invalid phone number, expected 11 digits, got %d", len(digitsOnly))
	}
	result := "8" + digitsOnly[1:]
	processedNumber, err := strconv.Atoi(result)
	if err != nil {
		return 0, err
	}
	return processedNumber, nil
}

// Set Order Reciever Contact Info. If corresponding slot is nil --> Fetch default Agent Contact Info by HREF from the body
func (u *ProcessedOrder) setOrderRecieverContactInfo(o unmarshaller.Order, r *moyskladapi.MoySkladProcessor) {
	if o.AttributesMap["Имя получателя"] != nil && o.AttributesMap["Телефон получателя"] != nil {
		u.RecieverName = o.AttributesMap["Имя получателя"].(string)
		recieverPhoneNumber := o.AttributesMap["Телефон получателя"].(string)
		tempNumber, err := processPhoneNumber(recieverPhoneNumber)
		if err != nil {
			panic(err)
		}
		u.RecieverPhoneNumber = tempNumber
	} else {
		reciever, err := unmarshaller.AgentUnmarshalling(r.FetchEntityByHREF(o.Agent.Meta.HREF))
		if err != nil {
			return
		}
		recieverPhoneNumber := reciever.Phone
		tempNumber, err := processPhoneNumber(recieverPhoneNumber)
		if err != nil {
			panic(err)
		}
		u.RecieverPhoneNumber = tempNumber
		if err != nil {
			panic(err)
		}
		u.RecieverName = reciever.Name
	}
}

// Get Shipment Address for the order, basically we can not check anything substantial here,
// as every order's address should be manually validated by an operator
func (u *ProcessedOrder) setOrderShipmentAddress(o unmarshaller.Order) {
	var b strings.Builder
	b.WriteString("Россия, ")
	b.WriteString(o.ShipmentAddress)
	u.ShipmentAddress = b.String()
}

// Set Order Delivery Interval, check for nil, as it's mandatory
func (u *ProcessedOrder) setOrderDeliveryInterval(o unmarshaller.Order) error {
	if o.AttributesMap["Интервал доставки"] != nil {
		tempInterval := o.AttributesMap["Интервал доставки"].(string)
		parts := strings.Split(tempInterval, "-")
		u.DeliveryIntervalFrom = parts[0]
		u.DeliveryIntervalUntil = parts[1]
	} else {
		return fmt.Errorf("%s - no interval", o.Name)
	}
	return nil
}

// Set Order Region, check for nil, in this case we can default it to "МСК", as 90% of deliveries are made in this region
func (u *ProcessedOrder) setOrderRegion(o unmarshaller.Order) {
	if o.AttributesMap["Регион доставки"] != nil {
		u.DeliveryRegion = o.AttributesMap["Регион доставки"].(string)
	} else {
		u.DeliveryRegion = "МСК"
	}
}

// Set Order Payment Method, no need to check for nil, as it is mandatory
func (u *ProcessedOrder) setOrderPaymentMethod(o unmarshaller.Order) {
	u.PaymentMethod = o.AttributesMap["Способ оплаты"].(string)
}

// Set Order Sum, will be converted to monetary amount during the xlsx file creation
func (u *ProcessedOrder) setOrderSum(o unmarshaller.Order) {
	u.Sum = float64(o.Sum / 100)
}

// Set Order Boxes Count, each can be nil, then default is set to 0. But not at the same time, as an order has to have at least 1 box
func (u *ProcessedOrder) setOrderBoxesCount(o unmarshaller.Order) {
	if o.AttributesMap["Кол-во коробок охл."] == nil && o.AttributesMap["Кол-во коробок зам."] == nil {
		return
	}
	if o.AttributesMap["Кол-во коробок охл."] != nil {
		tempChilledBoxesCount, err := strconv.Atoi(o.AttributesMap["Кол-во коробок охл."].(string))
		if err != nil {
			return
		}
		u.ChilledBoxes = uint8(tempChilledBoxesCount)
	} else {
		u.ChilledBoxes = 0
	}
	if o.AttributesMap["Кол-во коробок зам."] != nil {
		tempFrozenBoxesCount, err := strconv.Atoi(o.AttributesMap["Кол-во коробок зам."].(string))
		if err != nil {
			return
		}
		u.FrozenBoxes = uint8(tempFrozenBoxesCount)
	} else {
		u.FrozenBoxes = 0
	}
}

type Position struct {
	state       string
	quantity    float64
	weight      float64
	totalweight float64
}

func (p *Position) setState(pi *unmarshaller.ProductInfo) {
	if len(pi.Code) == 0 {
		return
	}
	runedID := []rune(pi.Code)
	switch runedID[0] {
	case '0':
		p.state = "frozen"
	case '1':
		p.state = "chilled"
	case '2':
		p.state = "any"
	}
}

func (p *Position) setWeight(pi *unmarshaller.ProductInfo) {
	p.weight = pi.Weight * 1000 //TODO move to constants
}

func (p *Position) setPositionTotalWeight() {
	p.totalweight = p.quantity * p.weight
}

func (u *ProcessedOrder) setPositions(o unmarshaller.Order, p *moyskladapi.MoySkladProcessor) error {
	chilled := make([]Position, 0)
	frozen := make([]Position, 0)
	services := make([]Position, 0)

	orderPositions, err := unmarshaller.PositionsUnmarshalling(p.FetchEntityByHREF(o.Positions.Meta.HREF))
	if err != nil {
		return err
	}

	for _, position := range orderPositions.Rows {
		var newPosition Position

		if position.Assortment.Meta.Type == "service" {
			newPosition.state = "any"
			newPosition.quantity = 1
			newPosition.weight = 0
			newPosition.totalweight = 0
			services = append(services, newPosition)
		} else {
			newPosition.quantity = position.Quantity

			positionSubInfo, err := unmarshaller.ProductInfoUnmarshalling(p.FetchEntityByHREF(position.Assortment.Meta.HREF))
			if err != nil {
				return err
			}
			newPosition.setState(positionSubInfo)
			newPosition.setWeight(positionSubInfo)
			newPosition.setPositionTotalWeight()
			switch newPosition.state {
			case "frozen":
				frozen = append(frozen, newPosition)
			case "chilled":
				chilled = append(chilled, newPosition)
			case "any":
				if chilled != nil {
					chilled = append(chilled, newPosition)
				} else {
					frozen = append(frozen, newPosition)
				}
			}
		}
	}
	u.frozenPositions = frozen
	u.chilledPositions = chilled
	u.servicePositions = services
	return nil
}

func (p *ProcessedOrder) setChilledTotalWeight() {
	var sumGramms float64
	for _, position := range p.chilledPositions {
		sumGramms += position.totalweight
	}
	roundedGramms := math.Ceil(sumGramms/500) * 500
	sumKG := roundedGramms / 1000
	if sumKG < 0.5 {
		p.ChilledWeight = 0.5
	} else {
		p.ChilledWeight = sumKG
	}
}

func (p *ProcessedOrder) setFrozenTotalWeight() {
	var sumGramms float64
	for _, position := range p.frozenPositions {
		sumGramms += position.totalweight
	}
	roundedGramms := math.Ceil(sumGramms/500) * 500
	sumKG := roundedGramms / 1000
	if sumKG < 0.5 {
		p.FrozenWeight = 0.5
	} else {
		p.FrozenWeight = sumKG
	}
}

func suitableForDelivery(o unmarshaller.Order) bool {
	slice := strings.Split(o.DeliveryPlannedMoment, " ")
	dayaftertomorrow := time.Now().AddDate(0, 0, 2)
	dayaftertomorrowstring := dayaftertomorrow.Format("2006-01-02")
	fmt.Println(dayaftertomorrowstring)
	fmt.Println(slice[0])
	if slice[0] == dayaftertomorrowstring && o.AttributesMap["Регион доставки"] == nil {
		fmt.Println("Filtered")
		fmt.Println(o.AttributesMap["Регион доставки"])
		return false
	}

	if slice[0] == dayaftertomorrowstring && o.AttributesMap["Регион доставки"] == "МСК" {
		fmt.Println("Filtered")
		fmt.Println(o.AttributesMap["Регион доставки"])
		return false
	}
	fmt.Println("Unfiltered")
	fmt.Println(o.AttributesMap["Регион доставки"])
	return true
}

func ProcessOrders(p *moyskladapi.MoySkladProcessor) (ordersCount int, orders *map[string]ProcessedOrder, err error) {
	storage := make(map[string]ProcessedOrder)
	unprocessedOrders, err := unmarshaller.BasicMoySkladResponseUnmarshalling(p.FetchDeliverableOrders())
	if err != nil {
		return 0, nil, err
	}

	refGoNumber := p.RefGoConfig.RGLatestOrder
	for _, order := range unprocessedOrders.Rows {

		if !suitableForDelivery(order) {
			continue
		}

		var newOrder ProcessedOrder
		order.UnmarshallOrderAttributes()
		newOrder.setOrderRefgoNumber(refGoNumber + 1)
		newOrder.setOrderHREF(order)
		newOrder.setOrderName(order)
		newOrder.setOrderDescription(order)
		newOrder.setOrderDeliveryPlannedDate(order)
		newOrder.setOrderShipmentAddress(order)
		newOrder.setOrderDeliveryInterval(order)
		newOrder.setOrderCounterpartyhref(order)
		newOrder.setOrderRecieverContactInfo(order, p)
		newOrder.setOrderRegion(order)
		newOrder.setOrderPaymentMethod(order)
		newOrder.setOrderSum(order)
		newOrder.setOrderBoxesCount(order)
		newOrder.setPositions(order, p)
		newOrder.setChilledTotalWeight()
		newOrder.setFrozenTotalWeight()
		storage[newOrder.HREF] = newOrder
		refGoNumber++
	}
	errs := config.ChangeRefGoLatest(refGoNumber)
	if errs != nil {
		fmt.Println(errs)
	}

	orders = &storage
	return ordersCount, orders, err
}
