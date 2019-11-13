package scrypto

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
)

const (
	BlockMode_CBC = "CBC"
	BlockMode_ECB = "ECB"
	Padding_PKCS5 = "PKCS5Padding"
)

type AesCryptor struct {
	key       []byte
	block     cipher.Block
	encrypter cipher.BlockMode
	decrypter cipher.BlockMode
	iv        []byte
	padding   padding
}

func NewAesCryptor(key, iv []byte, mode, pad string) (*AesCryptor, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	encrypter, err := getEncrypter(block, iv, mode)
	if err != nil {
		return nil, err
	}
	decrypter, err := getDecrypter(block, iv, mode)
	if err != nil {
		return nil, err
	}
	paddingMode, err := newPadding(pad)
	if err != nil {
		return nil, err
	}
	return &AesCryptor{
		key:       key,
		block:     block,
		encrypter: encrypter,
		decrypter: decrypter,
		iv:        iv,
		padding:   paddingMode,
	}, nil

}

//
func getEncrypter(block cipher.Block, iv []byte, mode string) (cipher.BlockMode, error) {
	switch mode {
	case BlockMode_CBC:
		return cipher.NewCBCEncrypter(block, iv), nil
	case BlockMode_ECB:
		return NewECBEncrypter(block), nil
	}
	return nil, fmt.Errorf("not support encrypter block mode: %s", mode)
}
func getDecrypter(block cipher.Block, iv []byte, mode string) (cipher.BlockMode, error) {
	switch mode {
	case BlockMode_CBC:
		return cipher.NewCBCDecrypter(block, iv), nil
	case BlockMode_ECB:
		return NewECBDecrypter(block), nil
	}
	return nil, fmt.Errorf("not support decrypter block mode: %s", mode)
}

//加密数据
func (m *AesCryptor) Encrypt(src []byte) ([]byte, error) {
	content := m.padding.Padding(src, m.block.BlockSize())
	encrypted := make([]byte, len(content))
	m.encrypter.CryptBlocks(encrypted, content)
	return encrypted, nil
}

//解密数据
func (m *AesCryptor) Decrypt(encryptedData []byte) ([]byte, error) {
	decrypted := make([]byte, len(encryptedData))
	m.decrypter.CryptBlocks(decrypted, encryptedData)
	return m.padding.UnPadding(decrypted), nil
}

//加密数据 CBC  PKCS5Padding
func CBCPKCS5PaddingAesEncrypt(key []byte, iv []byte, src []byte) ([]byte, error) {
	cryptor, err := NewAesCryptor(key, iv, BlockMode_CBC, Padding_PKCS5)
	if err != nil {
		return nil, err
	}
	return cryptor.Encrypt(src)
}

//解密数据  CBC  PKCS5Padding
func CBCPKCS5PaddingAesDecrypt(key []byte, iv []byte, encryptedData []byte) ([]byte, error) {
	cryptor, err := NewAesCryptor(key, iv, BlockMode_CBC, Padding_PKCS5)
	if err != nil {
		return nil, err
	}
	return cryptor.Decrypt(encryptedData)
}

//加密数据 ECB  PKCS5Padding
func ECBPKCS5PaddingAesEncrypt(key []byte, src []byte) ([]byte, error) {
	cryptor, err := NewAesCryptor(key, key, BlockMode_ECB, Padding_PKCS5)
	if err != nil {
		return nil, err
	}
	return cryptor.Encrypt(src)
}

//解密数据  ECB  PKCS5Padding
func ECBPKCS5PaddingAesDecrypt(key []byte, encryptedData []byte) ([]byte, error) {
	cryptor, err := NewAesCryptor(key, key, BlockMode_ECB, Padding_PKCS5)
	if err != nil {
		return nil, err
	}
	return cryptor.Decrypt(encryptedData)
}

//加密字符串数据
func AesEncryptString(key []byte, iv []byte, src string) (string, error) {
	bytes, err := CBCPKCS5PaddingAesEncrypt(key, iv, []byte(src))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
}

//解密字符串数据
func AesDecryptString(key []byte, iv []byte, encryptedData string) (string, error) {
	bytes, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", err
	}
	decrypt, err := CBCPKCS5PaddingAesDecrypt(key, iv, bytes)
	if err != nil {
		return "", err
	}
	return string(decrypt), nil
}
