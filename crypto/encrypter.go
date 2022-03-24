package crypto

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

const (
	EncrypterVersion1 = "v1"
)

type Encrypter interface {
	Encrypt(data string) (string, error)
	Decrypt(data string) (string, error)
}

type result struct {
	Version         string   `json:"version"`
	Length          uint32   `json:"length"`
	EncryptedChunks []string `json:"encrypted_chunks"`
	chunks          []string
}

func (r result) toString() string {
	data, _ := json.Marshal(r)
	return base64.StdEncoding.EncodeToString(data)
}

type encrypter struct {
	key string
	iv  string
}

func (e encrypter) Encrypt(data string) (string, error) {
	r, err := e.encryptV1(data)
	if err != nil {
		return "", err
	}
	return r.toString(), nil
}

func (e encrypter) Decrypt(data string) (string, error) {
	return "", nil
}

func (e encrypter) encryptV1(data string) (result, error) {
	var r = result{}
	r.chunks = e.splitToChunks(data, 16)
	r.Length = uint32(len(data))
	fmt.Println(r.chunks)

	return r, nil
}

func (e encrypter) splitToChunks(s string, n int) []string {
	var result []string
	for i := 1; i < len(s); i++ {
		if i%n == 0 {
			result = append(result, s[:i])
			s = s[i:]
			i = 1
		}
	}
	result = append(result, s)
	return result
}
