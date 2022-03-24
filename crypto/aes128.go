package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"strings"
)

// errors
var (
	ErrAES128Length          = errors.New("plaintext length should be 16")
	ErrInternalDecryptLength = errors.New("decrypt data not formatted properly")
	ErrKeyLength             = errors.New("key should be at least 1 character")
)

type AES128 struct{}

// Encrypt the main encryption endpoint where the plaintext encrypted with the
// arbitrary length key and the constant initialization vector
func (a AES128) Encrypt(plaintext, key string, iv []byte) (string, error) {
	if err := a.validateEncrypt(plaintext); err != nil {
		return "", err
	}
	keyHex, err := Padding(key, 16)
	if err != nil {
		return "", err
	}
	return a.encrypt(plaintext, keyHex, iv)
}

// Decrypt the main decryption endpoint where the ciphertext decrypted
func (a AES128) Decrypt(ciphertext, key string, iv []byte) (string, error) {
	keyHex, err := Padding(key, 16)
	if err != nil {
		return "", err
	}
	return a.decrypt(ciphertext, keyHex, iv)
}

// HashSize ...
func (a AES128) HashSize() int {
	return 16
}

// the pure implementation of the AES-CBC-128 encryption
func (a AES128) encrypt(plaintext string, key string, iv []byte) (string, error) {
	keyBytes := []byte(key)
	plaintextBytes := []byte(plaintext)

	c, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	data := make([]byte, len(plaintextBytes))
	copy(data, plaintextBytes)

	cipher.NewCBCEncrypter(c, iv).CryptBlocks(data, data)
	return string(data), nil
}

// the pure implementation of the AES-CBC-128 decryption
func (a AES128) decrypt(ciphertext string, key string, iv []byte) (string, error) {
	keyBytes := []byte(key)
	ciphertextBytes := []byte(ciphertext)
	c, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	data := make([]byte, len(ciphertextBytes))
	copy(data, ciphertextBytes)

	if len(data) != 16 {
		return "", ErrInternalDecryptLength
	}

	cipher.NewCBCDecrypter(c, iv).CryptBlocks(data, data)
	return string(data), nil
}

// validate the parameters of the Encrypt, to make sure everything is right for a correct CBC encryption
func (a AES128) validateEncrypt(plaintext string) error {
	// Set of rules to ensure high entropy
	if len(plaintext) != 16 {
		return ErrAES128Length
	}

	return nil
}

func Padding(key string, length int) (string, error) {
	key = strings.ToLower(key)
	key = strings.ReplaceAll(key, "-", "")
	if len(key) > length {
		return key[0:length], nil
	}
	var keyPointer = 0
	var blockKey []byte
	if len(key) < 1 {
		return "", ErrKeyLength
	}
	for i := 0; i < length; i++ {
		blockKey = append(blockKey, key[keyPointer])
		keyPointer++
		if keyPointer > len(key)-1 {
			keyPointer = 0
		}
	}
	return string(blockKey), nil
}
