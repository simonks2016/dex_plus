package httpClient

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Client struct {
	httpClient *http.Client

	queue chan Request

	workerSize int

	ctx    context.Context
	cancel context.CancelFunc

	wg sync.WaitGroup

	closeOnce sync.Once
	runOnce   sync.Once
}

type Config struct {
	WorkerSize int
	QueueSize  int

	Timeout  time.Duration
	ProxyURL string

	MaxIdleConn        int
	MaxIdleConnPerHost int
	IdleConnTimeout    time.Duration
}

func NewClient(cfg Config) *Client {
	if cfg.WorkerSize <= 0 {
		cfg.WorkerSize = 8
	}
	if cfg.QueueSize <= 0 {
		cfg.QueueSize = 1024
	}
	if cfg.Timeout <= 0 {
		cfg.Timeout = 5 * time.Second
	}
	if cfg.MaxIdleConn <= 0 {
		cfg.MaxIdleConn = 100
	}
	if cfg.MaxIdleConnPerHost <= 0 {
		cfg.MaxIdleConnPerHost = 100
	}
	if cfg.IdleConnTimeout <= 0 {
		cfg.IdleConnTimeout = 90 * time.Second
	}

	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,

		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,

		MaxIdleConns:        cfg.MaxIdleConn,
		MaxIdleConnsPerHost: cfg.MaxIdleConnPerHost,
		IdleConnTimeout:     cfg.IdleConnTimeout,
	}

	if cfg.ProxyURL != "" {
		proxyURL, err := url.Parse(cfg.ProxyURL)
		if err != nil {
			panic(err)
		}
		transport.Proxy = http.ProxyURL(proxyURL)
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &Client{
		httpClient: &http.Client{
			Timeout:   cfg.Timeout,
			Transport: transport,
		},
		queue:      make(chan Request, cfg.QueueSize),
		workerSize: cfg.WorkerSize,
		ctx:        ctx,
		cancel:     cancel,
	}
}

func (c *Client) Run() {
	c.runOnce.Do(func() {
		for i := 0; i < c.workerSize; i++ {
			// 开启worker
			c.wg.Go(func() {
				c.worker(i)
			})
		}
	})
}

func (c *Client) Close() {
	c.closeOnce.Do(func() {
		c.cancel()
		close(c.queue)
		// 等待
		c.wg.Wait()
		// 关闭http transport
		if transport, ok := c.httpClient.Transport.(*http.Transport); ok {
			transport.CloseIdleConnections()
		}
	})
}

func (c *Client) DoAsync(req Request) error {
	select {
	case <-c.ctx.Done():
		return c.ctx.Err()
	case c.queue <- req:
		return nil
	default:
		return fmt.Errorf("http client queue is full")
	}
}

func (c *Client) worker(workerID int) {
	for {
		select {
		case <-c.ctx.Done():
			return

		case req, ok := <-c.queue:
			if !ok {
				return
			}

			resp, err := c.doWithRetry(req)

			if req.Callback != nil {
				req.Callback(resp, err)
			}
		}
	}
}

func (c *Client) doWithRetry(req Request) (*Response, error) {
	if req.CreatedAt.IsZero() {
		req.CreatedAt = time.Now()
	}

	var lastResp *Response
	var lastErr error

	maxRetry := req.Retry
	if maxRetry < 0 {
		maxRetry = 0
	}

	for i := 0; i <= maxRetry; i++ {
		resp, err := c.do(req, i)
		if err == nil && resp != nil && resp.StatusCode < 500 {
			return resp, nil
		}

		lastResp = resp
		lastErr = err

		if i == maxRetry {
			break
		}

		interval := req.RetryInterval
		if interval <= 0 {
			interval = 200 * time.Millisecond
		}

		if req.EnableBackoff {
			interval = interval * time.Duration(1<<i)
		}

		select {
		case <-c.ctx.Done():
			return lastResp, c.ctx.Err()
		case <-time.After(interval):
		}
	}

	if lastResp != nil {
		lastResp.Err = lastErr
	}

	return lastResp, lastErr
}

func (c *Client) do(req Request, retryCount int) (*Response, error) {
	start := time.Now()

	ctx := req.Ctx
	if ctx == nil {
		ctx = c.ctx
	}

	timeout := req.Timeout
	if timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	requestURL, err := buildURL(req.URL, req.Query)
	if err != nil {
		return &Response{
			RequestID:  req.RequestId,
			Latency:    time.Since(start),
			RetryCount: retryCount,
			Err:        err,
		}, err
	}

	var body io.Reader
	if len(req.Body) > 0 {
		body = bytes.NewReader(req.Body)
	}

	httpReq, err := http.NewRequestWithContext(
		ctx,
		methodToString(req.Method),
		requestURL,
		body,
	)
	if err != nil {
		return &Response{
			RequestID:  req.RequestId,
			Latency:    time.Since(start),
			RetryCount: retryCount,
			Err:        err,
		}, err
	}

	for k, v := range req.Header {
		httpReq.Header.Set(k, v)
	}

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return &Response{
			RequestID:  req.RequestId,
			Latency:    time.Since(start),
			RetryCount: retryCount,
			Err:        err,
		}, err
	}
	defer httpResp.Body.Close()

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return &Response{
			RequestID:  req.RequestId,
			StatusCode: httpResp.StatusCode,
			Header:     httpResp.Header,
			Latency:    time.Since(start),
			RetryCount: retryCount,
			Err:        err,
		}, err
	}

	resp := &Response{
		RequestID:  req.RequestId,
		StatusCode: httpResp.StatusCode,
		Header:     httpResp.Header,
		Body:       respBody,
		Latency:    time.Since(start),
		RetryCount: retryCount,
	}

	return resp, nil
}

func buildURL(rawURL string, query map[string]string) (string, error) {
	if len(query) == 0 {
		return rawURL, nil
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	q := u.Query()
	for k, v := range query {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	return u.String(), nil
}

func methodToString(method Method) string {
	switch method {
	case GET:
		return http.MethodGet
	case POST:
		return http.MethodPost
	case PATCH:
		return http.MethodPatch
	case DELETE:
		return http.MethodDelete
	case PUT:
		return http.MethodPut
	default:
		return http.MethodGet
	}
}
