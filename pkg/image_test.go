package pkg

import "testing"

func TestName(t *testing.T) {

	//total := Pixabay().Total
	//t.Log(total)
	image := GetPixabayImage()
	t.Logf("%v", image)

}
