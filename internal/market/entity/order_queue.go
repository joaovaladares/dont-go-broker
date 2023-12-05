package entity

import "container/heap"

type OrderQueue struct {
	Orders []*Order
}

func (oq *OrderQueue) Less(i, j int) bool {
	return oq.Orders[i].Price < oq.Orders[j].Price
}

func (oq *OrderQueue) Swap(i, j int) {
	oq.Orders[i], oq.Orders[j] = oq.Orders[j], oq.Orders[i]
}

func (oq *OrderQueue) Len() int {
	return len(oq.Orders)
}

func (oq *OrderQueue) Push(x any) {
	oq.Orders = append(oq.Orders, x.(*Order))
}

func (oq *OrderQueue) Pop() any {
	oldOrders := oq.Orders
	numOldOrders := len(oldOrders)
	lastOrder := oldOrders[numOldOrders-1]
	oq.Orders = oldOrders[0 : numOldOrders-1]
	return lastOrder
}

func NewOrderQueue() *OrderQueue {
	return &OrderQueue{}
}

func initOrderQueue(orderQueue *OrderQueue) {
	heap.Init(orderQueue)
}
