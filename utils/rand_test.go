package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCutString(t *testing.T) {

	src := []rune("抖音号：1954716971")

	pos := -1
	for i, r := range src {
		if r == '：' {
			pos = i
			break
		}
	}
	assert.NotEqual(t, -1, pos)
	douyinId := string(src[pos+1:])
	t.Log(douyinId)
}
