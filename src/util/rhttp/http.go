package rhttp

import (
	"net/http"
	"time"
)

var (
	StatusOK int
)

func init() {
	StatusOK = http.StatusOK
}

type Client struct {
	*http.Client
}

// 获取httpclient
func NewHttpClient(timeout int) *Client {
	client := &Client{
		&http.Client{
			Transport: Transport(),
			Timeout:   time.Duration(timeout) * time.Second,
		},
	}

	return client
}
