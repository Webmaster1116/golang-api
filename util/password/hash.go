package password

import (
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

func Hash(password string) (hashHex string, err error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err == nil {
		hashHex = hex.EncodeToString(hash)
	}
	return hashHex, err
}

func Verify(password, hashHex string) bool {
	// get hash bytes
	hash, err := hex.DecodeString(hashHex)
	if err != nil {
		return false
	}
	// compare
	return bcrypt.CompareHashAndPassword(hash, []byte(password)) == nil
}
