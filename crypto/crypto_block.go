package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

const aeadNonceSize = 12

type BlockEncryptor struct {
	Key []byte
}

func (je *BlockEncryptor) Encrypt(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(je.Key)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)

	return append(nonce, ciphertext...), nil
}

func (je *BlockEncryptor) Decrypt(enc []byte) ([]byte, error) {
	if len(enc) < aeadNonceSize {
		return nil, ErrInsufficientData
	}

	nonce := enc[:aeadNonceSize]
	ciphertext := enc[aeadNonceSize:]

	block, err := aes.NewCipher(je.Key)
	if err != nil {
		panic(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
