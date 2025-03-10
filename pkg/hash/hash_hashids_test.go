package hash

import "testing"

func TestHashidsEncode(t *testing.T) {
	str, _ := New().HashidsEncode([]int{1})
	t.Log(str)

	//GyV5pJqXvwAR
}

func TestHashidsDecode(t *testing.T) {
	ids, _ := New().HashidsDecode("3R74vLA8pABq")
	t.Log(ids[0])
}
