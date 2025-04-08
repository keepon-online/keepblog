package area

import (
	"fmt"
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"path"
	"runtime"
)

func Area(intIP uint32) string {
	// 1、从 dbPath 加载整个 xdb 到内存
	cBuff, err := xdb.LoadContentFromFile("./data/ip2region.xdb")
	if err != nil {
		fmt.Printf("failed: %s\n", err.Error())
		return ""
	}

	// 2、用全局的 cBuff 创建完全基于内存的查询对象。
	searcher, err := xdb.NewWithBuffer(cBuff)
	if err != nil {
		fmt.Printf("failed to create searcher with content: %s\n", err)
		return ""
	}

	// 备注：并发使用，用整个 xdb 缓存创建的 searcher 对象可以安全用于并发。
	search, err := searcher.Search(intIP)

	if err != nil {
		fmt.Printf("failed to search with content:%s\n", err)
		return ""
	}
	return search
}

// GetCurrAbPath 获取当前文件绝对路径
func GetCurrAbPath() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(filename)
	}
	return abPath
}
