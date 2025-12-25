package cloudtag

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"

	exmath "gitee.com/jieepre/go-site/pkg/math"
)

// 预定义的美观配色方案（HSL色相值）
// 使用鲜艳但不刺眼的颜色
var colorPalette = []string{
	"hsl(210, 70%, 50%)", // 蓝色
	"hsl(340, 65%, 55%)", // 粉红
	"hsl(160, 60%, 45%)", // 青绿
	"hsl(280, 55%, 55%)", // 紫色
	"hsl(25, 75%, 55%)",  // 橙色
	"hsl(180, 60%, 45%)", // 青色
	"hsl(45, 80%, 50%)",  // 金黄
	"hsl(200, 65%, 50%)", // 天蓝
	"hsl(320, 60%, 55%)", // 洋红
	"hsl(100, 55%, 45%)", // 草绿
	"hsl(0, 65%, 55%)",   // 红色
	"hsl(240, 55%, 55%)", // 靛蓝
}

// CloudTags 生成标签云样式
// minFontSize: 最小字体大小(em)
// maxFontSize: 最大字体大小(em)
// ratio: 权重比例(文章数量)
func CloudTags(minFontSize, maxFontSize float64, ratio int) string {
	// 计算字体大小，使用对数缩放更平滑
	normalizedRatio := float64(ratio)
	if normalizedRatio > 10 {
		normalizedRatio = 10
	}
	size := minFontSize + ((maxFontSize - minFontSize) * (normalizedRatio / 10))
	round := exmath.Round(size, 2)
	fontSize := strconv.FormatFloat(round, 'f', -1, 32)

	// 从预定义配色中随机选择
	colorIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(colorPalette))))
	color := colorPalette[colorIndex.Int64()]

	// 生成样式
	style := fmt.Sprintf("font-size:%sem;color:%s;transition:all .3s ease", fontSize, color)
	return style
}

// CloudTagsWithWeight 带权重计算的标签云样式
func CloudTagsWithWeight(weight, maxWeight int) string {
	if maxWeight <= 0 {
		maxWeight = 1
	}
	ratio := (weight * 10) / maxWeight
	return CloudTags(1.0, 1.8, ratio)
}
