package math

import (
	"math"
	"math/rand"

	"strconv"
	"strings"
	"time"
)

// RandFloat64 生成范围随机float64
func RandFloat64(min, max float64) float64 {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	minStr := strconv.FormatFloat(min, 'f', -1, 64)
	// 不包含小数点
	if strings.Index(minStr, ".") == -1 {
		return max
	}
	multipleNum := len(minStr) - (strings.Index(minStr, ".") + 1)
	multiple := math.Pow10(multipleNum)
	minMult := min * multiple
	maxMult := max * multiple
	randVal := RandInt64(int64(minMult), int64(maxMult))
	result := float64(randVal) / multiple
	return result
}

// RandInt64 随机整数
func RandInt64(min, max int64) int64 {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	rand.Seed(time.Now().UnixNano())
	return rand.Int63n(max-min+1) + min
}

// Round 四舍五入，ROUND_HALF_UP 模式实现
// 返回将 val 根据指定精度 precision（十进制小数点后数字的数目）进行四舍五入的结果。precision 也可以是负数或零。
func Round(val float64, precision int) float64 {
	if precision == 0 {
		return math.Round(val)
	}

	p := math.Pow10(precision)
	if precision < 0 {
		return math.Floor(val*p+0.5) * math.Pow10(-precision)
	}

	return math.Floor(val*p+0.5) / p
}
