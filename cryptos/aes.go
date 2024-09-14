package cryptos

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/json"
	"fmt"
	"log"
)

var iv = []byte("2622233964834367")

type aesCompatibleCryptoJS struct {
	key   []byte
	block cipher.Block
}

func NewAesCompatibleCryptoJS(key string) *aesCompatibleCryptoJS {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		log.Fatal(err)
	}

	return &aesCompatibleCryptoJS{
		key:   []byte(key),
		block: block,
	}
}

func (a *aesCompatibleCryptoJS) pad(data []byte) []byte {
	padding := a.block.BlockSize() - len(data)%a.block.BlockSize()
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

func (a *aesCompatibleCryptoJS) Encrypt(v interface{}) (string, error) {
	var b []byte
	switch t := v.(type) {
	case map[string]interface{}:
		out, err := json.Marshal(v)
		if err != nil {
			return "", err
		}
		b = out
	case string:
		b = []byte(t)
	case int, int8, int16, int32, int64, float32, float64:
		b = []byte(fmt.Sprintf("%v", v))
	default:
		return "", fmt.Errorf("unexpect type %T", t)
	}

	b = a.pad(b)

	blockMode := cipher.NewCBCEncrypter(a.block, iv)
	ciphertext := make([]byte, len(b))
	blockMode.CryptBlocks(ciphertext, b)
	return B64Encode(ciphertext), nil
}

func (a *aesCompatibleCryptoJS) unpad(b []byte) []byte {
	n := len(b)
	u := int(b[n-1])
	return b[:(n - u)]
}

func (a *aesCompatibleCryptoJS) Decrypt(enc string) (string, error) {
	b, err := B64Decode(enc)
	if err != nil {
		return "", err
	}

	blockMode := cipher.NewCBCDecrypter(a.block, iv)
	dec := make([]byte, len(b))
	blockMode.CryptBlocks(dec, b)
	return string(a.unpad(dec)), nil
}
