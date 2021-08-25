package rhttp

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/tal-tech/go-zero/core/logx"
)

var (
	StatusOK int
)

func init() {
	StatusOK = http.StatusOK
}

type (
	Client struct {
		*http.Client
	}

	Header struct {
		Key string
		Val string
	}

	Option struct {
		Url     string     `json:"url"`
		Method  string     `json:"method"`
		Type    string     `json:"type"`
		Data    string     `json:"data"`
		Headers []*Header  `json:"headers"`
		UrlVal  url.Values `json:"url_val"`
	}
)

// NewHttpClient 获取httpclient
func NewHttpClient(timeout int) *Client {
	client := &Client{
		&http.Client{
			Transport: Transport(),
			Timeout:   time.Duration(timeout) * time.Second,
		},
	}

	return client
}

func (c *Client) Call(option *Option) ([]byte, error) {
	var err error

	var resp *http.Response
	switch option.Method {
	case http.MethodPost:
		resp, err = c.Post(option.Url, option.Type, strings.NewReader(option.Data))
	case http.MethodGet:
		resp, err = c.Get(option.Url)
	case "FORM":
		resp, err = c.PostForm(option.Url, option.UrlVal)
	default:
		req, reqErr := http.NewRequest(option.Method, option.Url, strings.NewReader(option.Data))
		if reqErr != nil {
			logx.Errorf("Client.Call: http.NewRequest fail, err=%v", reqErr)
			return nil, reqErr
		}
		if len(option.Headers) > 0 {
			for _, header := range option.Headers {
				req.Header.Set(header.Key, header.Val)
			}
		}
		resp, err = c.Do(req)
	}

	if err != nil {
		logx.Errorf("Client.Call: json.Unmarshal fail, err=%v", err)
		return nil, err
	}

	if resp != nil {
		defer resp.Body.Close()
	}

	logx.Infof("Client.Call: http status, code=%d", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		logx.Errorf("Client.Call: http status not ok, code=%d", resp.StatusCode)
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logx.Errorf("Client.Call: read fail, err=%v", err)
		return nil, err
	}

	return body, nil
}

func (c *Client) FetchWithJson(ctx context.Context, option *Option, typ interface{}) error {
	resp, err := c.Call(ctx, option)
	if err != nil {
		logx.Errorf("Client.FetchWithJson: read fail, err=%v", err)
		return err
	}

	if err := json.Unmarshal(resp, typ); err != nil {
		logx.Errorf("Client.FetchWithJson: json.Unmarshal fail, err=%v", err)
		return err
	}

	return nil
}

func (c *Client) FetchWithString(ctx context.Context, option *Option) (string, error) {
	resp, err := c.Call(ctx, option)
	if err != nil {
		logx.Errorf("Client.FetchWithString: read fail, err=%v", err)
		return "", err
	}

	return string(resp), nil
}

func (c *Client) PostFetchWithJson(ctx context.Context, option *Option, typ interface{}) error {
	option.Method = "post"
	option.Headers = []Header{
		{
			Key: "Content-Type",
			Val: "application/json",
		},
	}
	return c.FetchWithJson(ctx, option, typ)
}

func (c *Client) GetFetchWithJson(ctx context.Context, option *Option, typ interface{}) error {
	option.Method = "get"
	option.Headers = []Header{
		{
			Key: "Content-Type",
			Val: "application/json",
		},
	}
	return c.FetchWithJson(ctx, option, typ)
}