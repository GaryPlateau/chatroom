package utils

import (
	"time"
)

// 时间戳转换成日期函数
func UnixToTime(timestamp int) string {
	t := time.UnixMilli(int64(timestamp))
	return t.Format("2006-01-02 15:04:05")
}

// 日期转换成时间戳
func DateToUnix(str string) int64 {
	template := "2006-01-02 15:04:05"
	t, err := time.ParseInLocation(template, str, time.Local)
	if err != nil {
		return 0
	}
	return t.UnixMilli()
}

// 获取当前时间戳
func GetUnix() int64 {
	return time.Now().UnixMilli()
}

// 获取当前日期
func GetDate() string {
	template := "2006-01-02 15:04:05"
	return time.Now().Format(template)
}

// 获取年月日
func GetDay() string {
	template := "20060102"
	return time.Now().Format(template)
}
