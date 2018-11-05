package wrapper

import (
	"io"
	"fmt"
	"golang.org/x/crypto/pbkdf2"
	"crypto/sha256"
	"crypto/rand"
)

func Pbkdf2New(password []byte) ([]byte, error) {
	salt := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, fmt.Errorf("rand salt error:" + err.Error())
	}

	return append(salt, pbkdf2.Key(password, salt, 4096, 32, sha256.New)...), nil
}

func Pbkdf2Verify(password []byte, resultAndSalt []byte) bool {
	result := resultAndSalt[32:]
	salt := resultAndSalt[:32]
	newResult := pbkdf2.Key(password, salt, 4096, 32, sha256.New)
	if string(newResult) == string(result) {
		return true
	}
	return false
}
