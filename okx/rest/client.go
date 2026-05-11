package rest

import (
	"errors"
	"fmt"
	"time"

	"github.com/simonks2016/dex_plus/internal/httpClient"
	"github.com/simonks2016/dex_plus/okx/internal"
	"github.com/simonks2016/dex_plus/okx/param"
	"github.com/simonks2016/dex_plus/okx/response"
)

type Client struct {
	client    *httpClient.Client
	auth      *internal.Auth
	BaseUrl   string
	isSandBox bool
}

func (c *Client) PlaceOrder(params ...param.PlaceOrderParams) error {
	if len(params) == 0 {
		return errors.New("params is empty")
	}
	if len(params) > 20 {
		return errors.New("a maximum of 20 orders can be canceled at once")
	}
	// 生成path
	path := buildPath("/api/v5/trade/batch-orders")
	// POST请求
	if results, err := doPOST[[]response.ResultAsPlaceOrder](c.BaseUrl, path, c, params); err != nil {
		for _, r := range results {
			if r.SCode != "0" {
				return fmt.Errorf("%v, error_msg = %s", err, r.SMsg)
			}
		}
		return err
	} else {
		for _, r := range results {
			if r.SCode != "0" {
				return errors.New(r.SMsg)
			}
		}
		return nil
	}
}

func (c *Client) CancelOrder(params ...param.CancelOrder) error {
	if len(params) == 0 {
		return errors.New("params is empty")
	}
	if len(params) > 20 {
		return errors.New("a maximum of 20 orders can be canceled at once")
	}
	// 生成Path
	path := buildPath("/api/v5/trade/cancel-batch-orders")

	// 发送POST请求
	if results, err := doPOST[[]response.ResultAsCancelOrder](c.BaseUrl, path, c, params); err != nil {
		return err
	} else {
		for _, r := range results {
			if r.SCode != "0" {
				return errors.New(r.SMsg)
			}
		}
		return nil
	}
}

// GetInstruments 获取交易产品基础信息
func (c *Client) GetInstruments(
	instType string,
	queryParams ...QueryParam,
) ([]response.Instruments, error) {
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
	return doGET[[]response.Instruments](c.BaseUrl, path, c)
}

// GetPendingOrders 获取未成交订单
func (c *Client) GetPendingOrders(
	queryParams ...QueryParam,
) ([]response.PendingOrder, error) {

	if c == nil || c.client == nil {
		return nil, fmt.Errorf("client is nil")
	}

	path := buildPath("/api/v5/trade/orders-pending", queryParams...)

	return doGET[[]response.PendingOrder](c.BaseUrl, path, c)
}

// GetPositions 获取仓位信息
func (c *Client) GetPositions(queryParams ...QueryParam) ([]response.Position, error) {
	// 生成请求路径
	path := buildPath("/api/v5/account/positions", queryParams...)
	// GET请求
	return doGET[[]response.Position](c.BaseUrl, path, c)
}

// GetBalance 获取账户余额
func (c *Client) GetBalance(queryParams ...QueryParam) ([]response.AccountBalance, error) {
	// 生成参数
	path := buildPath("/api/v5/account/balance", queryParams...)
	// GET请求
	return doGET[[]response.AccountBalance](c.BaseUrl, path, c)
}

// GetOrderStatus 获取订单信息
func (c *Client) GetOrderStatus(
	instId string,
	queryParams ...QueryParam,
) ([]response.OrderStatus, error) {

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
	return doGET[[]response.OrderStatus](c.BaseUrl, path, c)
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

func (c *Client) Close() {
	//
	c.client.Close()
}
