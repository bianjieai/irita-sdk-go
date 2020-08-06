package store

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

type AES struct{}

func (AES) Encrypt(text string, key string) (string, error) {
	plaintext := []byte(text)
	k := generateKey(key)

	block, err := aes.NewCipher(k)
	if err != nil {
		return "", err
	}

	cipherText := make([]byte, aes.BlockSize+len(plaintext))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plaintext)

	return base64.URLEncoding.EncodeToString(cipherText), nil
}

func (AES) Decrypt(cryptoText string, key string) (string, error) {
	cipherText, _ := base64.URLEncoding.DecodeString(cryptoText)
	k := generateKey(key)
	block, err := aes.NewCipher(k)
	if err != nil {
		return "", err
	}

	if len(cipherText) < aes.BlockSize {
		return "", err
	}
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)
	return fmt.Sprintf("%s", cipherText), nil
}

func generateKey(key string) (genKey []byte) {
	keyBz := []byte(key)
	genKey = make([]byte, 32)
	copy(genKey, keyBz)
	for i := 32; i < len(keyBz); {
		for j := 0; j < 32 && i < len(keyBz); j, i = j+1, i+1 {
			genKey[j] ^= keyBz[i]
		}
	}
	return genKey
}
