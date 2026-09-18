package main

import (
	"fmt"
	"os"

	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: verify-ipdb <ip2region.xdb>")
		os.Exit(2)
	}
	path := os.Args[1]
	info, err := os.Stat(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "stat xdb: %v\n", err)
		os.Exit(1)
	}
	if info.Size() < 1024 {
		fmt.Fprintf(os.Stderr, "xdb is too small: %d bytes\n", info.Size())
		os.Exit(1)
	}
	content, err := xdb.LoadContentFromFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load xdb: %v\n", err)
		os.Exit(1)
	}
	searcher, err := xdb.NewWithBuffer(content)
	if err != nil {
		fmt.Fprintf(os.Stderr, "create searcher: %v\n", err)
		os.Exit(1)
	}
	defer searcher.Close()
	for _, ip := range []uint32{0x08080808, 0x01010101, 0x7f000001} {
		if _, err := searcher.Search(ip); err != nil {
			fmt.Fprintf(os.Stderr, "query %08x: %v\n", ip, err)
			os.Exit(1)
		}
	}
	fmt.Printf("validated %s (%d bytes)\n", path, info.Size())
}
