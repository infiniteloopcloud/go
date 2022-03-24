package crypto

import (
	"encoding/base64"
	"encoding/json"
	"strings"
)

const (
	EncrypterVersion1 = "v1"
)

type Encrypter interface {
	Encrypt(data string) (string, error)
	Decrypt(data string) (string, error)
}

type result struct {
	Version string    `json:"version"`
	Length  uint32    `json:"length"`
	V1      *resultV1 `json:"v1,omitempty"`
}

type resultV1 struct {
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
	dataJSON, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}
	var r result
	if err := json.Unmarshal([]byte(dataJSON), &r); err != nil {
		return "", err
	}

	switch r.Version {
	case EncrypterVersion1:
		return e.decryptV1(r.V1, r.Length)
	}

	return "", nil
}

func (e encrypter) encryptV1(data string) (result, error) {
	var r = result{
		Version: EncrypterVersion1,
		Length:  uint32(len(data)),
		V1:      &resultV1{},
	}
	r.V1.chunks = e.splitToChunks(data, 16)
	r.Length = uint32(len(data))

	for _, chunk := range r.V1.chunks {
		if len(chunk) != 16 {
			chunk = e.appendTo(chunk, 16)
		}
		encryptedChunk, err := AES128{}.Encrypt(chunk, e.key, []byte(e.iv))
		if err != nil {
			return result{}, err
		}
		r.V1.EncryptedChunks = append(r.V1.EncryptedChunks, base64.StdEncoding.EncodeToString([]byte(encryptedChunk)))
	}

	return r, nil
}

func (e encrypter) decryptV1(data *resultV1, length uint32) (string, error) {
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
	return res[0:length], nil
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
