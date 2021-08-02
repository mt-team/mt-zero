package rversion

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompareVersion(t *testing.T) {
	ast := assert.New(t)

	allTest := []struct {
		Cmp     string
		Defined string
		Result  CompareResult
	}{
		{"", "1.15", LessThan},
		{"1.15", "", GreaterThan},
		{"1.15", "1.15", Equal},
		{"1.15", "1.9", GreaterThan},
		{"1.1.9.1", "1.1.9", GreaterThan},
		{"1.1.9.1", "1.1.9.2", LessThan},
		{"1.1.9.1", "1.1.10", LessThan},
		{"1.1.9.1", "1.1.9.1.dev", LessThan},
		{"1.1.9.1.dev", "1.1.9.1", GreaterThan},
		{"1.1.9.1", "1.1.9.dev", LessThan},
		{"1.1.9.dev", "1.1.9.1", LessThan},
	}

	for _, tt := range allTest {
		ast.Equal(tt.Result, Compare(tt.Cmp, tt.Defined))
	}

}
