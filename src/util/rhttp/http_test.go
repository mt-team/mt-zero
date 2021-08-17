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

func TestClient_FetchWithJson(t *testing.T) {
	cli := NewHttpClient(10)
	a := assert.New(t)
	// 启动mock http服务
	go mockHttp()

	obj, err := cli.FetchWithJson(&Option{
		Url:    "http://127.0.0.1:8080/http/mock",
		Method: "put",
	}, reflect.TypeOf(DemoObj{}))
	if err != nil {
		t.Errorf("FetchWithJson: fail, err=%v", err)
	}

	demoObj := obj.(*DemoObj)
	a.Equal("zhang san", demoObj.Name)
	a.Equal(1, demoObj.Sex)

}

func mockHttp() {
	http.HandleFunc("/http/mock", func(w http.ResponseWriter, r *http.Request) {

		demoObj := DemoObj{
			Name: "zhang san",
			Sex:  1,
		}

		jsonB, _ := json.Marshal(demoObj)

		w.Write(jsonB)
	})

	http.ListenAndServe("127.0.0.1:8080", nil)
}
