package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func EncryptGCM(plaintext string, base64Key string) (string, error) {
	data := []byte(plaintext)

	key, err := base64.StdEncoding.DecodeString(base64Key)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonceSize := gcm.NonceSize()
	nonce := make([]byte, nonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ct := gcm.Seal(nil, nonce, data, nil)
	return base64.StdEncoding.EncodeToString(append(nonce, ct...)), nil
}

func DecryptGCM(ciphertextB64 string, base64Key string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(ciphertextB64)
	if err != nil {
		return ciphertextB64, nil
	}
	if len(data) < aes.BlockSize {
		return ciphertextB64, nil
	}
	key, err := base64.StdEncoding.DecodeString(base64Key)
	if err != nil {
		return ciphertextB64, nil
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return ciphertextB64, nil
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return ciphertextB64, nil
	}
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return ciphertextB64, nil
	}
	nonce, ct := data[:nonceSize], data[nonceSize:]
	ptBytes, err := gcm.Open(nil, nonce, ct, nil)
	if err != nil {
		return ciphertextB64, nil
	}
	return string(ptBytes), nil
}

func MaskIDCard(id string) string {
	runes := []rune(id)
	l := len(runes)
	if l <= 6 {
		return id
	}
	maskCount := l - 6
	return strings.Repeat("X", maskCount) + string(runes[maskCount:])
}
