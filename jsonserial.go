package jsonserial

import (
	"encoding/json"

	"jsonserial/compression"
	"jsonserial/crypto"
	"jsonserial/encoding"
)

type JsonSerial struct {
	encoding    encoding.Encoding
	crypto      crypto.Crypto
	compression compression.Compression
}

func NewJSONSerial() *JsonSerial {
	return &JsonSerial{}
}

func (je *JsonSerial) Dump(data interface{}) (string, error) {
	var dump []byte
	var err error

	dump, err = json.Marshal(data)
	if err != nil {
		return "", err
	}

	if je.compression != nil {
		dump, err = je.compression.Compress(dump)
		if err != nil {
			return "", err
		}
	}

	if je.crypto != nil {
		dump, err = je.crypto.Encrypt(dump)
		if err != nil {
			return "", err
		}
	}

	if je.encoding != nil {
		encoded := je.encoding.EncodeToString(dump)
		return encoded, nil
	}

	return string(dump), nil
}

func (je *JsonSerial) Load(enc string, out interface{}) error {
	var decoded []byte
	var err error

	if je.encoding != nil {
		decoded, err = je.encoding.DecodeString(enc)
		if err != nil {
			return err
		}
	} else {
		decoded = []byte(enc)
	}

	if je.crypto != nil {
		decoded, err = je.crypto.Decrypt(decoded)
		if err != nil {
			return err
		}
	}

	if je.compression != nil {
		decoded, err = je.compression.Decompress(decoded)
		if err != nil {
			return err
		}
	}

	return json.Unmarshal(decoded, out)
}

func (je *JsonSerial) SetEncoder(e encoding.Encoding) *JsonSerial {
	je.encoding = e
	return je
}

func (je *JsonSerial) SetCrypto(c crypto.Crypto) *JsonSerial {
	je.crypto = c
	return je
}

func (je *JsonSerial) SetCompression(c compression.Compression) *JsonSerial {
	je.compression = c
	return je
}
