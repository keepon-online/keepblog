package config

import (
	"fmt"
	"testing"
)

func TestInitConfig(t *testing.T) {

	fmt.Println(Get().Http)
}
