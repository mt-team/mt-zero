package rhttp

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type DemoObj struct {
	Name string `json:"name"`
	Sex  int    `json:"sex"`
}

var cli = NewHttpClient(10)
var addr = "127.0.0.1:20001"
var _ = mockHttp(addr)

func TestClient_FetchWithJson(t *testing.T) {
	a := assert.New(t)
	// 启动mock http服务

	ctx := context.Background()

	var demoObj DemoObj
	err := cli.FetchWithJson(ctx, &Option{
		Url:    fmt.Sprintf("http://%s/http/mock/json", addr),
		Method: "put",
	}, &demoObj)
	if err != nil {
		t.Errorf("FetchWithJson: fail, err=%v", err)
	}

	a.Equal("zhang san", demoObj.Name)
	a.Equal(1, demoObj.Sex)
}

func TestClient_FetchWithString(t *testing.T) {
	a := assert.New(t)

	ctx := context.Background()

	str, err := cli.FetchWithString(ctx, &Option{
		Url:    fmt.Sprintf("http://%s/http/mock/string", addr),
		Method: http.MethodGet,
	})
	if err != nil {
		t.Errorf("FetchWithJson: fail, err=%v", err)
	}

	a.Equal("Hello World", str)
}

func mockHttp(addr string) *http.Server {
	srv := &http.Server{Addr: addr}

	http.HandleFunc("/http/mock/json", func(w http.ResponseWriter, r *http.Request) {
		demoObj := DemoObj{
			Name: "zhang san",
			Sex:  1,
		}

		jsonB, _ := json.Marshal(demoObj)

		w.Write(jsonB)
	})

	http.HandleFunc("/http/mock/string", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("Httpserver: ListenAndServe() error: %s", err)
		}

	}()

	return srv
}
