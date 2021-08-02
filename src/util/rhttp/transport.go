package rhttp

import (
	"context"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/viki-org/dnscache"
)

var (
	resolver = dnscache.New(time.Minute * 3)
)

func Transport() *http.Transport {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
			Resolver: &net.Resolver{
				Dial: func(ctx context.Context, network string, addr string) (net.Conn, error) {
					separator := strings.LastIndex(addr, ":")
					ip, err := resolver.FetchOneString(addr[:separator])
					if err != nil {
						return nil, err
					}
					return net.Dial(network, ip+addr[separator:])
				},
			},
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
}
