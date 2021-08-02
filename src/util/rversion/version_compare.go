package rversion

import (
	"strconv"
	"strings"
)

type CompareResult int

const (
	Equal CompareResult = iota
	LessThan
	GreaterThan
)

func Compare(cmpVersion string, definedVersion string) CompareResult {
	if cmpVersion == definedVersion {
		return Equal
	}
	if cmpVersion == "" {
		return LessThan
	}
	if definedVersion == "" {
		return GreaterThan
	}
	args := strings.Split(cmpVersion, ".")
	lenCmp := len(args)
	defines := strings.Split(definedVersion, ".")
	lenDefined := len(defines)
	minLen := lenCmp
	if lenDefined < lenCmp {
		minLen = lenDefined
	}

	for i := 0; i < minLen; i++ {
		arg, err := strconv.ParseInt(args[i], 10, 64)
		if err != nil {
			return LessThan // 系统定义的version 一定是数字.数字的形式
		}
		defi, err := strconv.ParseInt(defines[i], 10, 64)
		if err != nil {
			return LessThan // 系统定义的version 一定是数字.数字的形式
		}
		if arg < defi {
			return LessThan
		} else if arg > defi {
			return GreaterThan
		}
	}
	if lenCmp < lenDefined {
		return LessThan
	} else if lenCmp == lenDefined {
		return Equal
	} else {
		return GreaterThan
	}
}
