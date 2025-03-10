package pkg

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func HmacHash256(plain, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(plain))
	return hex.EncodeToString(h.Sum(nil))
}
