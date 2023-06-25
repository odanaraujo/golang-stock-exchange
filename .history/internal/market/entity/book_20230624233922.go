package entity

import (
	"container/heap"
	"sync"
)

type Book struct {
	Order        []*Order
	Transaction  []*Transaction
	OrderChan    chan *Order
	OrdersChanOut chan *Order
	Wg           *sync.WaitGroup
}

func NewBook(orderChan chan *Order, orderChanOut chan *Order, wg *sync.WaitGroup) *Book {
	return &Book{
		Order:        []*Order{},
		Transaction:  []*Transaction{},
		OrderChan:    orderChan,
		OrdersChanOut: orderChanOut,
		Wg:           wg,
	}
}

func (b *Book) TradeBuyOrders() {

	buyOrders := make(map[string]*OrderQueue)

	for order := range b.OrderChan {
		asset := order.Asset.ID
		if buyOrders[asset] == nil {
			buyOrders[asset] = NewOrderQueue()
			heap.Init(buyOrders[asset])
		}
		buyOrders[asset].Push(order)
		if buyOrders[asset].Len() > 0 && buyOrders[asset].Orders[0].Price > order.Price {
			sellOrder := buyOrders[asset].Pop().(*Order)
			if sellOrder.PendingShares > 0 {
				b.bookTransaction(sellOrder, order)
				if sellOrder.PendingShares < 0 {
					buyOrders[asset].Push(sellOrder)
				}
			}
		}
	}
}

func (b *Book) TradeSellOrders() {

	sellOrders := make(map[string]*OrderQueue)

	for order := range b.OrderChan {
		asset := order.Asset.ID
		if sellOrders[asset] == nil {
			sellOrders[asset] = NewOrderQueue()
			heap.Init(sellOrders[asset])
		}
		sellOrders[asset].Push(order)
		if sellOrders[asset].Len() > 0 && sellOrders[asset].Orders[0].Price <= order.Price {
			buyOrder := sellOrders[asset].Pop().(*Order)
			if buyOrder.PendingShares > 0 {
				b.bookTransaction(buyOrder, order)
				if buyOrder.PendingShares > 0 {
					sellOrders[asset].Push(buyOrder)
				}
			}
		}
	}
}

func (b *Book) AddTransaction(transaction *Transaction, wg *sync.WaitGroup) {
	defer wg.Done()
	seelingShares := transaction.SellingOrder.PendingShares
	buingShares := transaction.BuyingOrder.PendingShares

	minShares := seelingShares

	if buingShares < minShares {
		minShares = buingShares
	}

	transaction.SellingOrder.Investor.UpdateAssertPosition(transaction.SellingOrder.Asset.ID, -minShares)
	transaction.AddSellOrderPendingShares(-minShares)
	transaction.BuyingOrder.Investor.UpdateAssertPosition(transaction.BuyingOrder.Asset.ID, +minShares)
	transaction.AddBuyOrderPendingShares(-minShares)

	transaction.CalculateTotal(transaction.BuyingOrder.Shares, transaction.BuyingOrder.Price)
	transaction.CloseBuyingOrder()
	transaction.CloseSellingOrder()

	b.Transaction = append(b.Transaction, transaction)
}

func (b *Book) bookTransaction(buyOrSellOrder *Order, order *Order) {
	transaction := NewTransaction(buyOrSellOrder, order, order.Shares, order.Price)
	b.AddTransaction(transaction, b.Wg)
	buyOrSellOrder.Transactions = append(buyOrSellOrder.Transactions, transaction)
	b.OrderChanOut <- buyOrSellOrder
	b.OrderChanOut <- order
}
