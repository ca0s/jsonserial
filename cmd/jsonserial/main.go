package main

import (
	hex "encoding/hex"
	"flag"
	"fmt"

	"jsonserial"
	"jsonserial/compression"
	"jsonserial/crypto"
	"jsonserial/encoding"
)

type ProbeData struct {
	SomeInt       int64  `json:"a"`
	SomeString    string `json:"b"`
	SomeUint      uint   `json:"c"`
	AnotherString string `json:"d"`
	SomeBool      bool   `json:"e"`
}

var key = "2077c3343664aa48fc33cebc4abf5ce0"

func main() {
	var (
		strCryptoMode  string
		strEncoding    string
		strCompression string

		decodeString string

		perfTest bool
	)

	flag.StringVar(&strCryptoMode, "crypto", "block", "Encryption mode, valid modes are: stream, block")
	flag.StringVar(&strEncoding, "encoding", "base32", "Final payload encoding. Valid modes are: base64, base32, hex")
	flag.StringVar(&strCompression, "compression", "flate", "Compression mode. Valid modes are: zlib, gzip, flate, lzw")
	flag.BoolVar(&perfTest, "perf-test", false, "Execute performance tests")
	flag.StringVar(&decodeString, "decode", "", "Decode given string")
	flag.Parse()

	if perfTest {
		testPerformance()
		return
	}

	key, err := hex.DecodeString(key)
	if err != nil {
		fmt.Printf("invalid key: %s\n", err)
		return
	}

	var crypt crypto.Crypto
	switch strCryptoMode {
	case "block":
		crypt = &crypto.BlockEncryptor{Key: key}
	case "stream":
		crypt = &crypto.StreamEncryptor{Key: key}
	default:
		fmt.Printf("invalid mode\n")
		return
	}

	var encoder encoding.Encoding
	switch strEncoding {
	case "base32":
		encoder = encoding.Base32Encoding
	case "base64":
		encoder = encoding.Base64Encoding
	case "hex":
		encoder = encoding.HexEncoding
	default:
		fmt.Printf("unknown encoding\n")
		return
	}

	var compress compression.Compression
	switch strCompression {
	case "zlib":
		compress = &compression.ZlibCompresser{}
	case "gzip":
		compress = &compression.GzipCompresser{}
	case "flate":
		compress = &compression.FlateCompresser{}
	case "lzw":
		compress = &compression.LzwCompresser{}
	default:
		fmt.Printf("unknown compression\n")
		return
	}

	jsonSerializer := jsonserial.NewJSONSerial().SetEncoder(
		encoder,
	).SetCompression(
		compress,
	).SetCrypto(
		crypt,
	)

	hostEncoder := encoding.HostnameEncoder{
		Zone: "zone.domain.tld.",
	}

	if decodeString == "" {
		data := ProbeData{
			SomeInt:       37337,
			SomeString:    "the brown fox",
			AnotherString: "ran fast i guess",
			SomeBool:      false,
			SomeUint:      123123,
		}

		enc, err := jsonSerializer.Dump(&data)
		if err != nil {
			fmt.Printf("error encrypting: %s\n", err)
			return
		}

		hname, err := hostEncoder.EncodeToHostname(enc)
		fmt.Printf("%s\nlen = %d\nerr=%s\n", hname, len(hname), err)

		decName := hostEncoder.DecodeHostname(hname)

		var dec ProbeData
		err = jsonSerializer.Load(decName, &dec)
		if err != nil {
			fmt.Printf("error decrypting: %s\n", err)
			return
		}

		fmt.Printf("decoded: %+v\n", dec)
	} else {
		decName := hostEncoder.DecodeHostname(decodeString)

		var dec ProbeData
		err = jsonSerializer.Load(decName, &dec)
		if err != nil {
			fmt.Printf("error decrypting: %s\n", err)
			return
		}

		fmt.Printf("decoded: %+v\n", dec)
	}
}
