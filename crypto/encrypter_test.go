package crypto

import "testing"

func TestEncrypter_Encrypt(t *testing.T) {
	encrypter{
		key: "test",
		iv:  "test",
	}.Encrypt("testdatatestdatatestdatatestdatatestdatatestdata1234")
}
