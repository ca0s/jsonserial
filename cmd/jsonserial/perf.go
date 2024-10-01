package main

import (
	"encoding/hex"
	"fmt"
	"jsonserial"
	"time"

	"jsonserial/compression"
	"jsonserial/crypto"
	"jsonserial/encoding"
)

func testPerformance() {
	key, _ := hex.DecodeString(key)

	data := ProbeData{
		SomeInt:       37337,
		SomeString:    "the brown fox",
		AnotherString: "ran fast i guess",
		SomeBool:      false,
		SomeUint:      123123,
	}

	streamCrypto := &crypto.StreamEncryptor{Key: key}
	blockCrypto := &crypto.BlockEncryptor{Key: key}

	type cryptoTest struct {
		Name   string
		Crypto crypto.Crypto
	}

	type encodingTest struct {
		Name     string
		Encoding encoding.Encoding
	}

	type compressionTest struct {
		Name        string
		Compression compression.Compression
	}

	cryptoTests := []cryptoTest{
		{
			Name:   "stream",
			Crypto: streamCrypto,
		},
		{
			Name:   "block",
			Crypto: blockCrypto,
		},
	}

	encodingTests := []encodingTest{
		{
			Name:     "hex",
			Encoding: encoding.HexEncoding,
		},
		{
			Name:     "base32",
			Encoding: encoding.Base32Encoding,
		},
		{
			Name:     "base64",
			Encoding: encoding.Base64Encoding,
		},
	}

	compressionTests := []compressionTest{
		{
			Name:        "flate",
			Compression: compression.FlateCompresser{},
		},
		{
			Name:        "gzip",
			Compression: compression.GzipCompresser{},
		},
		{
			Name:        "lzw",
			Compression: compression.LzwCompresser{},
		},
		{
			Name:        "zlib",
			Compression: &compression.ZlibCompresser{},
		},
	}

	nTests := 10

	for _, et := range encodingTests {
		for _, ct := range compressionTests {
			for _, xt := range cryptoTests {
				var testEncAcc time.Duration
				var testDecAcc time.Duration
				var testSizeAcc int

				fmt.Printf("------\ntesting %s / %s / %s, ntests=%d\n", et.Name, ct.Name, xt.Name, nTests)

				for i := 0; i < nTests; i++ {
					jsonSerial := jsonserial.NewJSONSerial().SetEncoder(et.Encoding).SetCompression(ct.Compression).SetCrypto(xt.Crypto)

					encryptionStartedTs := time.Now()

					encData, err := jsonSerial.Dump(data)
					if err != nil {
						fmt.Printf("error encrypting: %s\n", err)
						return
					}

					encryptionFinishedTs := time.Now()

					//fmt.Printf("enc: %s\n", encData)

					var decData ProbeData

					decryptionStartedTs := time.Now()

					err = jsonSerial.Load(encData, &decData)
					if err != nil {
						fmt.Printf("error decrypting: %s\n", err)
					}

					decryptionFinishedTs := time.Now()

					testEncAcc += encryptionFinishedTs.Sub(encryptionStartedTs)
					testDecAcc += decryptionFinishedTs.Sub(decryptionStartedTs)
					testSizeAcc += len(encData)
				}

				//fmt.Printf("dec: %+v\n", decData)

				fmt.Printf("\n\tencryption took %s\n", testEncAcc/time.Duration(nTests))
				fmt.Printf("\tdecryption took %s\n", testDecAcc/time.Duration(nTests))
				fmt.Printf("\tencoded size: %d\n", testSizeAcc/nTests)
			}
		}
	}
}
