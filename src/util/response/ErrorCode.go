package response

import (
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BizResponse struct {
	Code       int         `json:"errcode"`
	Message    string      `json:"errmsg"`
	Time       int64       `json:"time"`
	Data       interface{} `json:"data,omitempty"`
	grpcStatus error
}

func (b *BizResponse) Error() string {
	return b.Message
}

// 用于从已定义的ErrResponse重载生成包含指定msg的BizResponse
func (b *BizResponse) WithMessage(msg string) *BizResponse {
	return newBizResponse(b.Code, msg)
}

func newBizResponse(code int, msg string) *BizResponse {
	return &BizResponse{
		Code:       code,
		Message:    msg,
		Time:       time.Now().Unix(),
		grpcStatus: status.Error(codes.Code(code), msg),
	}
}

// ********下面声明的变量需要在init()里放到pool里
var (
	Success = newBizResponse(0, "")

	// -1 未知错误
	ErrUnknown = newBizResponse(-1, "系统繁忙，此时请开发者稍候再试")

	// 20xxx 公共参数验证
	ErrInvalidArgs = newBizResponse(20002, "参数错误，换个姿势再来一次吧")

	// 50xxx 音视频会议相关错误
	ErrTXMeetingInfo = newBizResponse(50001, "会议信息请求错误")

	// 90xxx app版本错误
	ErrAppCtl = newBizResponse(90001, "版本信息错误，请重试")
)
