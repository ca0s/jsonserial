package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

type StreamEncryptor struct {
	Key []byte
}

func (je *StreamEncryptor) Encrypt(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(je.Key)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewOFB(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return ciphertext, nil
}

func (je *StreamEncryptor) Decrypt(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(je.Key)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < block.BlockSize() {
		return nil, ErrInsufficientData
	}

	iv := ciphertext[:block.BlockSize()]
	plaintext := make([]byte, len(ciphertext)-aes.BlockSize)

	stream := cipher.NewOFB(block, iv)
	stream.XORKeyStream(plaintext, ciphertext[aes.BlockSize:])

	return plaintext, nil
}
