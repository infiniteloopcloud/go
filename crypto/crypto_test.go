package crypto

import (
	"testing"
)

func TestArgon2id(t *testing.T) {
	var salt = RandomString(DefaultSaltLength)
	var password = "test_password"

	alg, err := Get(Argon2id)
	if err != nil {
		t.Error(err)
	}

	hash := alg.Hash(password, salt)
	if err := alg.Verify(password, salt, hash); err != nil {
		t.Error(err)
	}
}

func TestRandomString(t *testing.T) {
	t.Log(RandomString(10))
}

func BenchmarkArgon2id_Hash(b *testing.B) {
	var salt = RandomString(DefaultSaltLength)
	var password = "test_password"

	alg, err := Get(Argon2id)
	if err != nil {
		b.Error(err)
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		alg.Hash(password, salt)
	}
}

func BenchmarkRandStringBytes(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		RandomString(32)
	}
}
