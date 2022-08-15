package hash

import (
	"crypto/sha1"
	"fmt"
	"math/rand"
	"time"
)

// SHA1Hasher uses SHA1 to hash passwords with provided salt
type SHA1Hasher struct {
	salt string
}

func NewSHA1Hasher() *SHA1Hasher {
	salt, err := salt()
	if err != nil {
		return nil
	}
	return &SHA1Hasher{salt: salt}
}

// Hash creates SHA1 hash of given password
func (h *SHA1Hasher) Hash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(h.salt)))
}

func salt() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	_, err := r.Read(b)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}