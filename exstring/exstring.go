package exstring

import (
	"math/rand"
	"time"
	"errors"
)

func RandomStr(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}



/**
 *  截取字符串
 *  start 開始下标
 *  end   結束下标(若为负数，或比最大下标大，自动转换为最大下标)
 */
func SubStr(str string, start int, end int) (string,error) {
	if start < 0 {
		return "",errors.New("param error!")
	}

	rs := []rune(str)
	length := len(rs)

	if end < 0 || end > length-1 {
		end = length-1
	}

	if start > length || start > end {
		return "",errors.New("param error!")
	}

	return string(rs[start:end]),nil
}



/**
 * 字符串查找
 */
func InSlice(sStr []string,str string) bool {
	for _,value := range sStr {
		if value == str {
			return true
		}
	}
	return false
}