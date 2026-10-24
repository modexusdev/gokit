package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
)

var (
	ErrInvalidKey = errors.New(
		"encryption key must be 32 bytes",
	)

	ErrInvalidCiphertext = errors.New(
		"invalid ciphertext",
	)
)

// Encrypt encrypts plaintext using AES-256-GCM.
//
// The key must be a base64 encoded 32-byte key.
// The returned ciphertext is base64 encoded and
// contains the random nonce followed by the
// encrypted data.
func Encrypt(
	plaintext string,
	key string,
) (string, error) {
	if plaintext == "" {
		return "", nil
	}

	keyBytes, err := decodeKey(key)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", fmt.Errorf(
			"create cipher: %w",
			err,
		)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf(
			"create gcm: %w",
			err,
		)
	}

	nonce := make(
		[]byte,
		gcm.NonceSize(),
	)

	if _, err := rand.Read(nonce); err != nil {
		return "", fmt.Errorf(
			"generate nonce: %w",
			err,
		)
	}

	ciphertext := gcm.Seal(
		nil,
		nonce,
		[]byte(plaintext),
		nil,
	)

	data := append(
		nonce,
		ciphertext...,
	)

	return base64.StdEncoding.EncodeToString(
		data,
	), nil
}

// Decrypt decrypts a value created by Encrypt.
//
// The key must be the same base64 encoded
// 32-byte key that was used for encryption.
func Decrypt(
	encrypted string,
	key string,
) (string, error) {
	if encrypted == "" {
		return "", nil
	}

	keyBytes, err := decodeKey(key)
	if err != nil {
		return "", err
	}

	data, err := base64.StdEncoding.DecodeString(
		encrypted,
	)
	if err != nil {
		return "", ErrInvalidCiphertext
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", fmt.Errorf(
			"create cipher: %w",
			err,
		)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf(
			"create gcm: %w",
			err,
		)
	}

	nonceSize := gcm.NonceSize()

	if len(data) < nonceSize {
		return "", ErrInvalidCiphertext
	}

	nonce := data[:nonceSize]
	ciphertext := data[nonceSize:]

	plaintext, err := gcm.Open(
		nil,
		nonce,
		ciphertext,
		nil,
	)
	if err != nil {
		return "", ErrInvalidCiphertext
	}

	return string(plaintext), nil
}

func decodeKey(
	key string,
) ([]byte, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(
		key,
	)
	if err != nil {
		return nil, ErrInvalidKey
	}

	if len(keyBytes) != 32 {
		return nil, ErrInvalidKey
	}

	return keyBytes, nil
}
