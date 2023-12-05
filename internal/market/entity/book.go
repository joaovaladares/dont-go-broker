package entity

import (
	"sync"
)

type Book struct {
	Order         []*Order
	Transactions  []*Transaction
	OrdersChanIn  chan *Order
	OrdersChanOut chan *Order
	Wg            *sync.WaitGroup
}

func NewBook(orderChanIn chan *Order, orderChanOut chan *Order, wg *sync.WaitGroup) *Book {
	return &Book{
		Order:         []*Order{},
		Transactions:  []*Transaction{},
		OrdersChanIn:  orderChanIn,
		OrdersChanOut: orderChanOut,
		Wg:            wg,
	}
}

func (b *Book) Trade() {
	orderQueues := map[string]*OrderQueue{
		BuyOrder:  NewOrderQueue(),
		SellOrder: NewOrderQueue(),
	}

	for _, queue := range orderQueues {
		initOrderQueue(queue)
	}

	processOrders := func(ownOrders *OrderQueue, counterOrders *OrderQueue, order *Order) {
		ownOrders.Push(order)

		if counterOrders.Len() == 0 || isNotValidOrder(order, counterOrders) {
			return
		}

		counterOrder := counterOrders.Pop().(*Order)
		if counterOrder.PendingShares <= 0 {
			return
		}

		var transaction *Transaction
		if order.OrderType == BuyOrder {
			transaction = NewTransaction(counterOrder, order, order.Shares, counterOrder.Price)
		}

		if order.OrderType == SellOrder {
			transaction = NewTransaction(order, counterOrder, order.Shares, counterOrder.Price)
		}

		b.AddTransaction(transaction, b.Wg)
		counterOrder.Transactions = append(counterOrder.Transactions, transaction)
		order.Transactions = append(order.Transactions, transaction)
		b.OrdersChanOut <- counterOrder
		b.OrdersChanOut <- order

		if counterOrder.PendingShares > 0 {
			counterOrders.Push(counterOrder)
		}
	}

	for order := range b.OrdersChanIn {
		if queue, ok := orderQueues[order.OrderType]; ok {
			processOrders(queue, orderQueues[getCounterOrderType(order.OrderType)], order)
		}
	}
}

func (b *Book) AddTransaction(transaction *Transaction, wg *sync.WaitGroup) {
	defer wg.Done()

	sellingShares := transaction.SellingOrder.PendingShares
	buyingShares := transaction.BuyingOrder.PendingShares

	minShares := min(sellingShares, buyingShares)

	transaction.SellingOrder.Investor.UpdateAssetPosition(transaction.SellingOrder.Asset.ID, -minShares)
	transaction.AddSellOrderPendingShares(-minShares)
	transaction.BuyingOrder.Investor.UpdateAssetPosition(transaction.BuyingOrder.Asset.ID, minShares)
	transaction.AddBuyOrderPendingShares(-minShares)

	transaction.CalculateTotal(transaction.Shares, transaction.BuyingOrder.Price)

	transaction.CloseBuyOrder()
	transaction.CloseSellOrder()

	b.Transactions = append(b.Transactions, transaction)
}
