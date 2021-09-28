package bizResponse

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BizResponse struct {
	Code       int         `json:"code"`
	Message    string      `json:"message"`
	Time       int64       `json:"time,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	ErrorIndex string      `json:"errorIndex,omitempty"`
	grpcStatus error
}

func (b *BizResponse) Error() string {
	return b.Message
}

func newBizResponse(code int, msg string) *BizResponse {
	resp := &BizResponse{
		Code:       code,
		Message:    msg,
		Time:       time.Now().Unix(),
		grpcStatus: status.Error(codes.Code(code), msg),
	}

	if code != 0 {
		resp.ErrorIndex = fmt.Sprintf("error_index_%d", code)
	}

	return resp
}

// WithMessage 用于从已定义的ErrResponse重载生成包含指定msg的BizResponse
func (b *BizResponse) WithMessage(msg string) *BizResponse {
	return newBizResponse(b.Code, msg)
}

// WithData 用于从已定义的ErrResponse重载生成包含指定返回的数据
func (b *BizResponse) WithData(data interface{}) *BizResponse {
	resp := newBizResponse(b.Code, b.Message)
	resp.Data = data

	return resp
}

func (b *BizResponse) EqualsTo(err error) bool {
	if err == nil {
		return false
	}
	c := FromError(err)
	if c == nil {
		return false
	}
	return (b == c) || (b.Code == c.Code && b.Message == c.Message)
}

func (b *BizResponse) generateKey() string {
	return fmt.Sprintf("%d%v", b.Code, b.Message)
}

func FromError(err error) (resp *BizResponse) {
	if err == nil {
		return nil
	}

	switch err {
	case context.Canceled:
		return ErrClientCancel
	case context.DeadlineExceeded:
		return ErrDeadlineExceed
	default:

	}

	s, ok := status.FromError(err)
	if !ok {
		if biz, ok := err.(*BizResponse); ok {
			return biz
		} else {
			return ErrInternalFailed.WithMessage(err.Error())
		}
	}

	if s.Code() < 100 { // grpc 通常都是100以内
		switch s.Code() {
		case codes.Canceled:
			return ErrClientCancel
		}
		return ErrInternalFailed.WithMessage(err.Error())
	}

	if br, ok := pool[fmt.Sprintf("%d%v", s.Code(), s.Message())]; ok {
		return br
	}

	return newBizResponse(int(s.Code()), s.Message())
}

// ********下面声明的变量需要在init()里放到pool里
var (
	Success = newBizResponse(0, "")

	// -1 未知错误
	ErrUnknown = newBizResponse(-1, "系统繁忙，此时请开发者稍候再试")

	// 20xxx 公共参数验证
	ErrInvalidArgs    = newBizResponse(20002, "参数错误，换个姿势再来一次吧")
	ErrInternalFailed = newBizResponse(20003, "服务器正在开小差呢，请稍后重试")
	ErrClientCancel   = newBizResponse(20004, "服务连接断开了, 请稍后重试")
	ErrDeadlineExceed = newBizResponse(20005, "服务器繁忙，请稍候重试")

	// 50xxx 音视频会议相关错误
	ErrTXMeetingInfo = newBizResponse(50001, "会议信息请求错误")

	// 90xxx app版本错误
	ErrAppCtl = newBizResponse(90001, "版本信息错误，请重试")
)

var (
	pool = make(map[string]*BizResponse)
)

func init() {
	pool[ErrUnknown.generateKey()] = ErrUnknown
	pool[ErrClientCancel.generateKey()] = ErrClientCancel
	pool[ErrInternalFailed.generateKey()] = ErrInternalFailed
	pool[ErrDeadlineExceed.generateKey()] = ErrDeadlineExceed
}
