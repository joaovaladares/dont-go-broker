package entity

const (
	BuyOrder  = "BUY"
	SellOrder = "SELL"
)

type Order struct {
	ID            string
	Investor      *Investor
	Asset         *Asset
	Shares        int
	PendingShares int
	Price         float64
	OrderType     string
	Status        string
	Transactions  []*Transaction
}

func NewOrder(orderID string, investor *Investor, asset *Asset, shares int, price float64, orderType string) *Order {
	return &Order{
		ID:            orderID,
		Investor:      investor,
		Asset:         asset,
		Shares:        shares,
		PendingShares: shares,
		Price:         price,
		OrderType:     orderType,
		Status:        "OPEN",
		Transactions:  []*Transaction{},
	}
}

func getCounterOrderType(orderType string) string {
	if orderType == BuyOrder {
		return SellOrder
	}

	return BuyOrder
}

func isNotValidOrder(order *Order, orderQueue *OrderQueue) bool {
	if (order.OrderType == BuyOrder && orderQueue.Orders[0].Price > order.Price) ||
		(order.OrderType == SellOrder && orderQueue.Orders[0].Price < order.Price) {
		return true
	}

	return false
}
