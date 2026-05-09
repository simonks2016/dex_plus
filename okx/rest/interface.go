package rest

import "github.com/simonks2016/dex_plus/okx/Response"

type OKXRestAPI interface {
	GetInstruments(instType string, queryParams ...QueryParam) ([]Response.Instruments, error)
	GetPendingOrders(queryParams ...QueryParam) ([]Response.PendingOrder, error)
	GetPositions(...QueryParam) ([]Response.Position, error)
	GetBalance(...QueryParam) ([]Response.AccountBalance, error)
	GetOrderStatus(instId string, queryParams ...QueryParam) ([]Response.OrderStatus, error)

	PlaceOrder() error
	CancelOrder() error
}
