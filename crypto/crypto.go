package crypto

import "fmt"

const (
	BlockMode  = 0
	StreamMode = 1
)

type Crypto interface {
	Encrypt(plaintext []byte) ([]byte, error)
	Decrypt(ciphertext []byte) ([]byte, error)
}

var ErrInsufficientData = fmt.Errorf("insufficient data")
