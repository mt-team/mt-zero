package rversion

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateVersion(t *testing.T) {
	ast := assert.New(t)

	testCase := []struct {
		Input  string
		Output int64
	}{
		{"2.7.2", 34013696},
		{"2.6.2", 33948160},
		{"3.7.2", 50790912},
	}
	for _, v := range testCase {
		ast.Equal(v.Output, CalculateRecommendVersion(v.Input))
	}
}
