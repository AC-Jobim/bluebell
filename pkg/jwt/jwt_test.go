package jwt

import (
	"fmt"
	"testing"
)

func TestGenToken(t *testing.T) {
	token, _ := GenToken(33896775014158336, "zhaobin")
	fmt.Println(token)
	//解析token
	claims, _ := ParseToken(token)
	fmt.Println(claims)
}
