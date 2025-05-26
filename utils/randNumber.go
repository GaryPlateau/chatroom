package utils

import (
	"math/rand"
	"time"
)

/*
*
  - size 随机码的位数
  - kind 0    // 纯数字
    1    // 小写字母
    2    // 大写字母
    3    // 数字、大小写字母
*/
func CreateRandString(size int, kind int) string {
	rand.Seed(time.Now().UnixNano())
	ikind := kind
	kinds := [][]int{[]int{10, 48}, []int{26, 97}, []int{26, 65}}
	rebytes := make([]byte, size)
	isAll := kind > 2 || kind < 0
	for i := 0; i < size; i++ {
		if isAll {
			ikind = rand.Intn(3)
		}
		scope, base := kinds[ikind][0], kinds[ikind][1]
		rebytes[i] = uint8(base + rand.Intn(scope))
	}
	return string(rebytes)
}

func RandBGStyle() string {
	bgStyle := ""
	rand.Seed(time.Now().UnixNano())
	randNum := rand.Intn(5)
	switch randNum {
	case 0:
		bgStyle = "bg-primary"
	case 1:
		bgStyle = "bg-success"
	case 2:
		bgStyle = "bg-info"
	case 3:
		bgStyle = "bg-warning"
	case 4:
		bgStyle = "bg-danger"
	default:
		bgStyle = "bg-primary"
	}
	return bgStyle
}
