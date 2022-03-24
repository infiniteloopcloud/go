package crypto

import (
	"encoding/base64"
	"encoding/json"
	"strings"
)

type Encrypter interface {
	Encrypt(data string) (string, error)
	Decrypt(data string) (string, error)
}

func NewEncypter(key, iv string) Encrypter {
	return encrypter{
		key: key,
		iv:  iv,
	}
}

type result struct {
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
	r, err := e.encrypt(data)
	if err != nil {
		return "", err
	}
	return r.toString(), nil
}

func (e encrypter) Decrypt(data string) (string, error) {
	dataJSON, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}
	var r result
	if err := json.Unmarshal([]byte(dataJSON), &r); err != nil {
		return "", err
	}

	return e.decrypt(r)
}

func (e encrypter) encrypt(data string) (result, error) {
	var r = result{
		Length: uint32(len(data)),
	}
	r.chunks = e.splitToChunks(data, 16)
	r.Length = uint32(len(data))

	for _, chunk := range r.chunks {
		if len(chunk) != 16 {
			chunk = e.appendTo(chunk, 16)
		}
		encryptedChunk, err := AES128{}.Encrypt(chunk, e.key, []byte(e.iv))
		if err != nil {
			return result{}, err
		}
		r.EncryptedChunks = append(r.EncryptedChunks, base64.StdEncoding.EncodeToString([]byte(encryptedChunk)))
	}

	return r, nil
}

func (e encrypter) decrypt(data result) (string, error) {
	for _, chunk := range data.EncryptedChunks {
		decoded, err := base64.StdEncoding.DecodeString(chunk)
		if err != nil {
			return "", err
		}
		decrypted, err := AES128{}.Decrypt(string(decoded), e.key, []byte(e.iv))
		if err != nil {
			return "", err
		}
		data.chunks = append(data.chunks, decrypted)
	}
	res := strings.Join(data.chunks, "")
	return res[0:data.Length], nil
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

func (e encrypter) appendTo(chunk string, n int) string {
	var appendN = ""

	for i := 0; i < n-len(chunk); i++ {
		appendN += "x"
	}

	return chunk + appendN
}
