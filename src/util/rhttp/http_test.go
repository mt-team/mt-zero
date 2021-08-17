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

func TestClient_FetchWithJson(t *testing.T) {
	a := assert.New(t)
	// 启动mock http服务
	go mockHttp()

	obj, err := cli.FetchWithJson(&Option{
		Url:    "http://127.0.0.1:8080/http/mock/json",
		Method: "put",
	}, reflect.TypeOf(DemoObj{}))
	if err != nil {
		t.Errorf("FetchWithJson: fail, err=%v", err)
	}

	demoObj := obj.(*DemoObj)
	a.Equal("zhang san", demoObj.Name)
	a.Equal(1, demoObj.Sex)

}

func TestClient_FetchWithString(t *testing.T) {
	a := assert.New(t)
	// 启动mock http服务
	go mockHttp()

	str, err := cli.FetchWithString(&Option{
		Url:    "http://127.0.0.1:8080/http/mock/string",
		Method: http.MethodGet,
	})
	if err != nil {
		t.Errorf("FetchWithJson: fail, err=%v", err)
	}

	a.Equal("Hello World", str)
}

func mockHttp() {
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

	http.ListenAndServe("127.0.0.1:8080", nil)
}
