package rest

import (
	"fmt"
	"time"

	"github.com/simonks2016/dex_plus/internal/httpClient"
	"github.com/simonks2016/dex_plus/okx/Response"
	"github.com/simonks2016/dex_plus/okx/internal"
)

type Client struct {
	client    *httpClient.Client
	auth      *internal.Auth
	BaseUrl   string
	isSandBox bool
}

// GetInstruments 获取交易产品基础信息
func (c *Client) GetInstruments(
	instType string,
	queryParams ...QueryParam,
) ([]Response.Instruments, error) {
	if c == nil || c.client == nil {
		return nil, fmt.Errorf("client is nil")
	}
	if instType == "" {
		return nil, fmt.Errorf("instType is required")
	}

	params := []QueryParam{
		WithInstType(instType),
	}
	params = append(params, queryParams...)
	// 生成请求路径
	path := buildPath("/api/v5/public/instruments", params...)
	// Get 请求
	return doGET[[]Response.Instruments](c.BaseUrl, path, c)
}

// GetPendingOrders 获取未成交订单
func (c *Client) GetPendingOrders(
	queryParams ...QueryParam,
) ([]Response.PendingOrder, error) {

	if c == nil || c.client == nil {
		return nil, fmt.Errorf("client is nil")
	}

	path := buildPath("/api/v5/trade/orders-pending", queryParams...)

	return doGET[[]Response.PendingOrder](c.BaseUrl, path, c)
}

// GetPositions 获取仓位信息
func (c *Client) GetPositions(queryParams ...QueryParam) ([]Response.Position, error) {
	// 生成请求路径
	path := buildPath("/api/v5/account/positions", queryParams...)
	// GET请求
	return doGET[[]Response.Position](c.BaseUrl, path, c)
}

// GetBalance 获取账户余额
func (c *Client) GetBalance(queryParams ...QueryParam) ([]Response.AccountBalance, error) {
	// 生成参数
	path := buildPath("/api/v5/account/balance", queryParams...)
	// GET请求
	return doGET[[]Response.AccountBalance](c.BaseUrl, path, c)
}

// GetOrderStatus 获取订单信息
func (c *Client) GetOrderStatus(
	instId string,
	queryParams ...QueryParam,
) ([]Response.OrderStatus, error) {

	if instId == "" {
		return nil, fmt.Errorf("instId is required")
	}
	params := []QueryParam{
		WithInstId(instId),
	}
	// 添加参数
	params = append(params, queryParams...)
	// 生成path
	path := buildPath("/api/v5/trade/order", params...)
	// GET请求
	return doGET[[]Response.OrderStatus](c.BaseUrl, path, c)
}

func (c *Client) PlaceOrder() error {
	//TODO implement me
	panic("implement me")
}

func (c *Client) CancelOrder() error {
	//TODO implement me
	panic("implement me")
}

func NewOKXRestClient(opts ...Option) OKXRestAPI {

	cli := &Client{
		client: httpClient.NewClient(httpClient.Config{
			WorkerSize: 10,
			QueueSize:  100,
			Timeout:    time.Second * time.Duration(30),
		}),
		auth:    nil,
		BaseUrl: "https://www.okx.com",
	}

	for _, opt := range opts {
		opt(cli)
	}
	// 启动client
	cli.client.Run()
	// 返回
	return cli
}

func (r *Client) Close() {
	//
	r.client.Close()
}
