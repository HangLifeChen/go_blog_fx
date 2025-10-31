package utils

import (
	"encoding/json"
	"math/rand"
	"time"
)

func MarshalToString(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

func UnmarshalFromString(s string, v interface{}) error {
	return json.Unmarshal([]byte(s), v)
}

// GenNumericCode 生成指定位数的纯数字验证码
// 参数：n 表示验证码位数
// 返回：长度为 n 的数字字符串（例如 "123456"）
//
// 特点：
// - 不会因为 n 太大而整型溢出（相比旧版本）
// - 每一位都是 0~9 的随机数字
// - 适合生成短信验证码、邮箱验证码、支付确认码等
func GenNumericCode(n int) string {
	if n <= 0 {
		// 如果 n 非法（例如 0 或负数），返回空字符串
		return ""
	}

	// 设置随机数种子（保证每次生成的结果都不同）
	rand.Seed(time.Now().UnixNano())

	// 创建一个字节切片用于存放生成的数字字符
	digits := make([]byte, n)

	// 循环生成 n 位随机数字
	for i := 0; i < n; i++ {
		// rand.Intn(10) 会返回 0~9 的随机整数
		digits[i] = byte('0' + rand.Intn(10))
	}

	// 将字节切片转换为字符串返回
	return string(digits)
}

// 不安全有栈溢出的风险
// generate a random numeric code with n digits
// func GenNumericCode(n int) string {
// 	if n <= 0 {
// 		return ""
// 	}
// 	rand.Seed(time.Now().UnixNano())

// 	// calculate the minimum and maximum values for the random number
// 	min := int32(1)
// 	for i := 1; i < n; i++ {
// 		min *= 10
// 	}
// 	max := min*10 - 1

// 	// generate a random number with n digits
// 	code := rand.Int31n(max-min+1) + min

// 	// format the random number as a string with n digits
// 	return fmt.Sprintf("%0*d", n, code)
// }
