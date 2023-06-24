package entity

const (
	ORDER_OPEN   = "OPEN"
	ORDER_BUY    = "BUY"
	ORDER_SELL   = "SELL"
	ORDER_CLOSED = "CLOSED"
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

func NewOrder(orderID string,
	investor *Investor,
	asset *Asset,
	shares int,
	price float64,
	ordertype string) *Order {
	return &Order{
		ID:            orderID,
		Investor:      investor,
		Asset:         asset,
		Shares:        shares,
		PendingShares: shares,
		Price:         price,
		OrderType:     ordertype,
		Status:        ORDER_OPEN,
		Transactions:  []*Transaction{},
	}
}
