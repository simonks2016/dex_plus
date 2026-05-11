package rest

import (
	"github.com/simonks2016/dex_plus/okx/param"
	"github.com/simonks2016/dex_plus/okx/response"
)

type OKXRestAPI interface {
	GetInstruments(instType string, queryParams ...QueryParam) ([]response.Instruments, error)
	GetPendingOrders(queryParams ...QueryParam) ([]response.PendingOrder, error)
	GetPositions(...QueryParam) ([]response.Position, error)
	GetBalance(...QueryParam) ([]response.AccountBalance, error)
	GetOrderStatus(instId string, queryParams ...QueryParam) ([]response.OrderStatus, error)

	PlaceOrder(...param.PlaceOrderParams) error
	CancelOrder(...param.CancelOrder) error
}
