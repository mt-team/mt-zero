package response

import (
	"context"
	"runtime/debug"

	"github.com/tal-tech/go-zero/core/logx"
	"google.golang.org/grpc"
)

func RecoveryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			logx.WithContext(ctx).Errorf("Panic err: %v", string(debug.Stack()))
			err = ErrInternalFailed.grpcStatus
		}
	}()

	resp, err = handler(ctx, req)
	if biz, ok := err.(*BizResponse); ok {
		return resp, biz.grpcStatus
	}

	return resp, err
}
