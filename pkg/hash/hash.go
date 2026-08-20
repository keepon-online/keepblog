package hash

var _ Hash = (*hash)(nil)

// defaultSecret 是历史版本的内置 salt。作为默认值保留以兼容既有文章的
// hashids URL；可通过 config.yaml 的 hashids.salt 覆盖（见 hash.Configure）。
// 注意：更换 salt 会使已发布文章的 URL 全部失效。
const defaultSecret = "i1ydX9RtHyuJTrw7frcu"
const length = 12

// configuredSecret 运行期生效的 salt，由应用启动时通过 Configure 注入
var configuredSecret = defaultSecret

// Configure 覆盖默认 hashids salt，传入空串时保持默认值不变
func Configure(salt string) {
	if salt != "" {
		configuredSecret = salt
	}
}

type Hash interface {
	i()

	// HashidsEncode 加密
	HashidsEncode(params []int) (string, error)

	// HashidsDecode 解密
	HashidsDecode(hash string) ([]int, error)
}

type hash struct {
	secret string
	length int
}

func New() Hash {
	return NewWithOptions(configuredSecret, length)
}
func NewWithOptions(secret string, length int) Hash {
	return &hash{
		secret: secret,
		length: length,
	}
}

func (h *hash) i() {}
