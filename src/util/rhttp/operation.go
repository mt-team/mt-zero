package rhttp

import (
	"context"
	"io"
	"net/http"

	"github.com/opentracing/opentracing-go"
)

func (c *Client) Get(ctx context.Context, url string) (resp *http.Response, err error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "http.Get")
	defer span.Finish()

	//start := time.Now()
	//defer Record(url, start)

	return c.Client.Get(url)
}

func (c *Client) Post(ctx context.Context, url, contentType string, body io.Reader) (resp *http.Response, err error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "http.Post")
	defer span.Finish()

	//start := time.Now()
	//defer Record(url, start)

	return c.Client.Post(url, contentType, body)
}

func (c *Client) Do(ctx context.Context, request *http.Request) (resp *http.Response, err error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "http.Do")
	defer span.Finish()

	//start := time.Now()
	//defer RecordRequest(request, start)

	return c.Client.Do(request)
}

func NewRequest(method, url string, body io.Reader) (request *http.Request, err error) {
	return http.NewRequest(method, url, body)
}
