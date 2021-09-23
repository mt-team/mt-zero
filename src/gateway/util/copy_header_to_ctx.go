package util

import (
	"context"
	"encoding/json"
	"net/http"

	"google.golang.org/grpc/metadata"
)

func CpoyHeaderToCtx(ctx context.Context, r *http.Request) context.Context {

	header := map[string]string{}
	for i, j := range r.Header {
		if firstIsX(i) {
			// 默认只允许一个value
			header[i] = j[0]
		}
	}

	headerStr, _ := json.Marshal(header)
	md := metadata.Pairs("x-header", string(headerStr))
	ctx = metadata.NewOutgoingContext(ctx, md)
	return ctx
}

func firstIsX(s string) bool {
	v := s[:1]
	if v == "X" || v == "x" {
		return true
	}

	return false
}
