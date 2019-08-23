package scrypto

import (
	"bytes"
	"fmt"
)

type padding interface {
	Padding(cipherText []byte, blockSize int) []byte
	UnPadding(encrypt []byte) []byte
}

//PKCS5包装
type PKCS5Padding struct {
}

func (m *PKCS5Padding) Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

func (m *PKCS5Padding) UnPadding(encrypt []byte) []byte {
	padding := encrypt[len(encrypt)-1]
	return encrypt[:len(encrypt)-int(padding)]
}

func newPadding(mode string) (padding, error) {
	switch mode {
	case Padding_PKCS5:
		return &PKCS5Padding{}, nil
	}
	return nil, fmt.Errorf("not support padding mode")

}
