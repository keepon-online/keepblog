package hash

var _ Hash = (*hash)(nil)

const secret = "i1ydX9RtHyuJTrw7frcu"
const length = 12

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
	return NewWithOptions(secret, length)
}
func NewWithOptions(secret string, length int) Hash {
	return &hash{
		secret: secret,
		length: length,
	}
}

func (h *hash) i() {}
