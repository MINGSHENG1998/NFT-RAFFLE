package helpers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
)

var (
	AesEncryptionHelper IAesEncrptionHelper = NewAesEncryptionHelper()

	secret_key string = DotEnvHelper.GetEnvVariable("AES_ENCRYTION_KEY")
)

type IAesEncrptionHelper interface {
	AesGCMEncrypt(text string) (string, error)
	AesGCMDecrypt(text string) (string, error)
}

type aesEncrytionHelperStruct struct{}

func NewAesEncryptionHelper() IAesEncrptionHelper {
	return &aesEncrytionHelperStruct{}
}

func (h *aesEncrytionHelperStruct) AesGCMEncrypt(text string) (string, error) {
	// When decoded the key should be 16 bytes (AES-128) or 32 (AES-256).
	block, err := aes.NewCipher([]byte(secret_key))
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	plainText := []byte(text)

	// generate nonce
	nonce := make([]byte, aesgcm.NonceSize())

	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	cipherText := aesgcm.Seal(nonce, nonce, plainText, nil)
	return fmt.Sprintf("%x", cipherText), nil
}

func (h *aesEncrytionHelperStruct) AesGCMDecrypt(text string) (string, error) {

	// When decoded the key should be 16 bytes (AES-128) or 32 (AES-256).
	block, err := aes.NewCipher([]byte(secret_key))
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	encryptedText, err := hex.DecodeString(text)
	if err != nil {
		return "", err
	}

	nonceSize := aesgcm.NonceSize()
	nonce := encryptedText[:nonceSize]
	cipherText := encryptedText[nonceSize:]

	plainText, err := aesgcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", err
	}

	// [:] operator allows you to create a slice from an array, optionally using start and end bounds
	// https://stackoverflow.com/questions/47722542/what-does-the-symbol-mean-in-go
	return string(plainText[:]), nil
}
