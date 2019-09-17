package encrypt

import (
	"fmt"
	"testing"
)

func TestAesEncrypt(t *testing.T) {
	encode, _ := AesEncrypt("dsfads")
	fmt.Println(encode)
	decode, _ := AesDecrypt(encode)
	fmt.Println(decode)
}
