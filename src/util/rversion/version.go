package rversion

import (
	"strconv"
	"strings"
)

func CalculateRecommendVersion(version string) int64 {
	clientVersion := int64(0)
	array := strings.Split(version, ".")
	if len(array) >= 3 {
		computVersionFlag := true
		versionL1, err := strconv.Atoi(array[0])
		if err != nil {
			computVersionFlag = false //错误不需要再计算
		}
		versionL2, err := strconv.Atoi(array[1])
		if err != nil {
			computVersionFlag = false
		}
		versionL3, err := strconv.Atoi(array[2])
		if err != nil {
			computVersionFlag = false
		}
		if computVersionFlag { //是否需要计算
			clientVersion = int64(versionL1<<24 | versionL2<<16 | versionL3<<8)
		}
	}
	return clientVersion
}
