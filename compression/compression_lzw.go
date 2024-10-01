package compression

import (
	"bytes"
	"compress/lzw"
	"io"
)

type LzwCompresser struct{}

func (LzwCompresser) Compress(data []byte) ([]byte, error) {
	var output []byte
	buf := bytes.NewBuffer(output)

	compresser := lzw.NewWriter(buf, lzw.LSB, 8)

	if _, err := compresser.Write(data); err != nil {
		return nil, err
	}

	if err := compresser.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (LzwCompresser) Decompress(data []byte) ([]byte, error) {
	buf := bytes.NewBuffer(data)

	decompresser := lzw.NewReader(buf, lzw.LSB, 8)

	output, err := io.ReadAll(decompresser)
	if err != nil {
		return nil, err
	}

	if err := decompresser.Close(); err != nil {
		return nil, err
	}

	return output, nil
}
