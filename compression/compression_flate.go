package compression

import (
	"bytes"
	"compress/flate"
	"io"
)

type FlateCompresser struct{}

func (fc FlateCompresser) Compress(data []byte) ([]byte, error) {
	var output []byte
	buf := bytes.NewBuffer(output)

	compresser, err := flate.NewWriter(buf, flate.BestCompression)
	if err != nil {
		return nil, err
	}

	if _, err := compresser.Write(data); err != nil {
		return nil, err
	}

	if err := compresser.Flush(); err != nil {
		return nil, err
	}

	if err := compresser.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (fc FlateCompresser) Decompress(data []byte) ([]byte, error) {
	buf := bytes.NewBuffer(data)

	decompresser := flate.NewReader(buf)

	output, err := io.ReadAll(decompresser)
	if err != nil {
		return nil, err
	}

	if err := decompresser.Close(); err != nil {
		return nil, err
	}

	return output, nil
}
