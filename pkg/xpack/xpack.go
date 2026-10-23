//go:build !xpack

package xpack

func RemoveTamper(website string) {}

func LoadGpuInfo() []any {
	return nil
}
