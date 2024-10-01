package encoding

import (
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
)

type Encoding interface {
	DecodeString(s string) ([]byte, error)
	EncodeToString(src []byte) string
}

var Base32Encoding = base32.StdEncoding.WithPadding(base32.NoPadding)
var Base64Encoding = base64.RawURLEncoding

type hexEncoding struct{}

func (hexEncoding) DecodeString(s string) ([]byte, error) {
	return hex.DecodeString(s)
}

func (hexEncoding) EncodeToString(src []byte) string {
	return hex.EncodeToString(src)
}

var HexEncoding = hexEncoding{}
