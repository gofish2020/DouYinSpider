package utils

import (
	"math/rand"
	"time"
)

// 随机对象*rand
var r = rand.New(rand.NewSource(time.Now().UnixNano()))

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-=")

// 生成随机字符串
func RandString(n int) string {
	result := make([]rune, n)
	for i := 0; i < len(result); i++ {
		// 从letters中随机获取一个字符，保存到result中
		result[i] = rune(letters[r.Intn(len(letters))])
	}
	return string(result)
}

func GetDouYinId(data string) string {
	src := []rune(data)
	pos := -1
	for i, r := range src {
		if r == '：' {
			pos = i
			break
		}
	}

	if pos == -1 {
		return ""
	}
	douyinId := string(src[pos+1:])
	return douyinId
}
