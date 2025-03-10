package hash

import (
	"github.com/speps/go-hashids/v2"
)

func (h *hash) HashidsEncode(params []int) (string, error) {
	hd := hashids.NewData()
	hd.Salt = h.secret
	hd.MinLength = h.length
	hs, _ := hashids.NewWithData(hd)
	hashStr, _ := hs.Encode(params)
	return hashStr, nil
}

func (h *hash) HashidsDecode(hash string) ([]int, error) {
	hd := hashids.NewData()
	hd.Salt = h.secret
	hd.MinLength = h.length

	ids, err := hashids.NewWithData(hd)
	if err != nil {
		return nil, err
	}
	withError, err := ids.DecodeWithError(hash)
	if err != nil {
		return nil, err
	}
	return withError, nil
}
