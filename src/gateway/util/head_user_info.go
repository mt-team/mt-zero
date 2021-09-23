package util

import (
	"encoding/json"
	"net/http"
)

type userInfo struct {
	UserUuid string `json:"userUuid"`
}

func GetHeaderUserUuid(r *http.Request) string {
	userInfoStr := r.Header.Get("x-user-info")
	if userInfoStr == "" {
		return ""
	}

	user := &userInfo{}
	err := json.Unmarshal([]byte(userInfoStr), user)
	if err != nil {
		return ""
	}
	return user.UserUuid
}
