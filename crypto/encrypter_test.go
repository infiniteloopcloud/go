package crypto

import (
	"testing"
)

func TestEncrypter_Encrypt(t *testing.T) {
	dataToEncrypt := "testdatatestdatatestdatatestdatatestdatatestdata1234"
	e := encrypter{
		key: "bhXRirFB8IaQxxjm",
		iv:  "oZDRWCHryRtlVA1I",
	}
	resultData, err := e.Encrypt(dataToEncrypt)
	if err != nil {
		t.Fatal(err)
	}

	decrypted, err := e.Decrypt(resultData)
	if err != nil {
		t.Fatal(err)
	}
	if decrypted != dataToEncrypt {
		t.Errorf("Decrypted is not the original data, original: %s, decrypted: %s", dataToEncrypt, decrypted)
	}
}
