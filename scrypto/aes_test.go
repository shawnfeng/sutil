package scrypto

import (
	"encoding/base64"
	"testing"
)

var key = []byte("1234567887654321")
var iv = []byte("abcdefghabcdefgh")

func Test_CBCPKCS5PaddingAesEncrypt(t *testing.T) {
	bytes, e := CBCPKCS5PaddingAesEncrypt(key, iv, []byte("86-15810007783"))
	s := base64.StdEncoding.EncodeToString(bytes)
	t.Log(s, e)
}
func Test_CBCPKCS5PaddingAesDecrypt(t *testing.T) {
	en := "DoGjETjeEyQLvrhqw1Kv2g=="
	bytes, err := base64.StdEncoding.DecodeString(en)
	t.Log(err)
	decrypt, e := CBCPKCS5PaddingAesDecrypt(key, iv, bytes)
	t.Log(string(decrypt), e)
}
func Test_AesEncryptString(t *testing.T) {
	src := "86-12345678901"
	s, e := AesEncryptString(key, iv, src)
	t.Log(s, e)
	s, e = AesDecryptString(key, iv, s)
	t.Log(s, e)

}
