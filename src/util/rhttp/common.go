package rhttp

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/tal-tech/go-zero/core/logx"
)

func CommonFun(ctx context.Context, httpCli *Client, url string, data interface{}) ([]byte, error) {

	logx.WithContext(ctx)
	logx.Infof("url<%v> data<%+v>", url, data)
	reqData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	resp, err := httpCli.Post(ctx, url, "application/json; charset=utf-8", strings.NewReader(string(reqData)))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != StatusOK {
		return nil, fmt.Errorf("wrong reqeust")
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	logx.Infof("resp <%+v>", string(body))
	return body, nil
}
