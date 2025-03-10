package jwttoken

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {

	token := "1232"
	createToken, err := CreateToken(token)
	if err != nil {
		return
	}
	fmt.Println(createToken)

	parseToken, err := ParseToken(createToken)
	if err != nil {
		return
	}

	fmt.Println(parseToken)

}
