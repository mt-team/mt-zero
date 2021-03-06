package bizResponse

import (
	"net/http"

	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/rest/httpx"
)

func Response(w http.ResponseWriter, resp interface{}) {
	switch err := resp.(type) {
	case *BizResponse:
		if err == nil {
			break
		}
		httpx.OkJson(w, err)
		logx.Errorf("%v", err)
		return
	case error:
		err = ErrUnknown
		httpx.OkJson(w, err)
		logx.Errorf("%v", err)
		return
	}

	// 成功的返回数据
	body := Success
	body.Data = resp
	httpx.OkJson(w, body)
}
