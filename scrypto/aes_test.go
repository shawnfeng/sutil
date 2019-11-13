package scrypto

import (
	"encoding/base64"
	"encoding/hex"
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
	var tests = []struct {
		input string
	}{
		{""},
		{"86-12345678901"},
		{"abc"},
		{"121231223"},
	}
	for _, test := range tests {
		s, _ := AesEncryptString(key, iv, test.input)
		t.Log(s)
		s, _ = AesDecryptString(key, iv, s)
		t.Log(s)
		if s != test.input {
			t.Error(test.input)
		}
	}

}
func Test_AesECBEncryptString(t *testing.T) {
	keys, _ := base64.StdEncoding.DecodeString("aqPmuoOm1CliRUdn3TkTKQ==")
	var tests = []struct {
		input string
	}{
		{`{"pid":"0","id":"34","name":"测试部门"}`},
		//{"86-12345678901"},
		//{"abc"},
		//{"121231223"},
	}
	for _, test := range tests {
		sb, _ := ECBPKCS5PaddingAesEncrypt(keys, []byte(test.input))
		s := hex.EncodeToString(sb)
		t.Log(s)
		ss, _ := hex.DecodeString(s)
		sb, _ = ECBPKCS5PaddingAesDecrypt(keys, ss)
		s = string(sb)
		t.Log(s)
		if s != test.input {
			t.Error(test.input)
		}
	}

}
