package jwt

import (
	"fmt"
	"testing"
)

func TestGenToken(t *testing.T) {
	token, _ := GenToken("zhaobin")
	fmt.Println(token)
	//解析token
	claims, _ := ParseToken(token)
	fmt.Println(claims)
}
