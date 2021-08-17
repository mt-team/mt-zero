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
		req, err := http.NewRequest(option.Method, option.Url, strings.NewReader(option.Data))
		if err != nil {
			logx.Errorf("Client.Call: http.NewRequest fail, err=%v", err)
			return nil, err
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

	if err != nil {
		return nil, err
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

func (c *Client) FetchWithJson(option *Option, typ reflect.Type) (interface{}, error) {
	resp, err := c.Call(option)
	if err != nil {
		logx.Errorf("Client.FetchWithJson: read fail, err=%v", err)
		return nil, err
	}

	var ret = reflect.New(typ)
	if err := json.Unmarshal(resp, ret.Interface()); err != nil {
		logx.Errorf("Client.FetchWithJson: json.Unmarshal fail, err=%v", err)
		return nil, err
	}

	return ret.Interface(), nil
}

func (c *Client) FetchWithString(option *Option) (string, error) {
	resp, err := c.Call(option)
	if err != nil {
		logx.Errorf("Client.FetchWithString: read fail, err=%v", err)
		return "", err
	}

	return string(resp), nil
}