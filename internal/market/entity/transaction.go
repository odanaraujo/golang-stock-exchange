package entity

import (
	"github.com/google/uuid"
	"time"
)

type Transaction struct {
	ID           string
	SellingOrder *Order
	BuyingOrder  *Order
	Shares       int
	Price        float64
	Total        float64
	DateTime     time.Time
}

func NewTransaction(sellingOrder *Order, buyingOrder *Order, shares int, price float64) *Transaction {
	total := float64(shares) * price
	return &Transaction{
		ID:           uuid.New().String(),
		SellingOrder: sellingOrder,
		BuyingOrder:  buyingOrder,
		Shares:       shares,
		Price:        price,
		Total:        total,
		DateTime:     time.Now(),
	}
}

func (t *Transaction) CalculateTotal(shares int, price float64) {
	t.Total = float64(shares) * price
}

func (t *Transaction) CloseBuyingOrder() {
	if t.BuyingOrder.PendingShares == 0 {
		t.BuyingOrder.Status = ORDER_CLOSED
	}
}

func (t *Transaction) CloseSellingOrder() {
	if t.SellingOrder.PendingShares == 0 {
		t.SellingOrder.Status = ORDER_CLOSED
	}
}

func (t *Transaction) AddBuyOrderPendingShares(shares int) {
	t.SellingOrder.PendingShares += shares
}

func (t *Transaction) AddSellOrderPendingShares(shares int) {
	t.SellingOrder.PendingShares += shares
}
