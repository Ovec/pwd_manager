package password

import (
	"crypto/rand"
	"math/big"
)

func Generate(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-=_+"

	randomBytes := make([]byte, length)
	for i := range randomBytes {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		randomBytes[i] = charset[n.Int64()]
	}

	return string(randomBytes), nil
}
