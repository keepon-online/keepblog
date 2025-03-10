package cloudtag

import (
	"crypto/rand"
	"fmt"
	exmath "gitee.com/jieepre/go-site/pkg/math"
	"math/big"
	"strconv"
)

func CloudTags(minFontSize, maxFontSize float64, ratio int) string {
	size := minFontSize + ((maxFontSize - minFontSize) * float64(ratio))
	round := exmath.Round(size, 2)
	fontSize := strconv.FormatFloat(round, 'f', -1, 32)
	r, _ := rand.Int(rand.Reader, big.NewInt(200))
	g, _ := rand.Int(rand.Reader, big.NewInt(200))
	b, _ := rand.Int(rand.Reader, big.NewInt(200))
	rgb := fmt.Sprintf("rgb(%d,%d,%d)", r, g, b)
	style := fmt.Sprintf("font-size:%sem;color:%s", fontSize, rgb)
	return style
}
