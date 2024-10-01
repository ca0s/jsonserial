package compression

import (
	"bytes"
	"compress/gzip"
	"io"
)

type GzipCompresser struct{}

func (GzipCompresser) Compress(data []byte) ([]byte, error) {
	var output []byte
	buf := bytes.NewBuffer(output)

	compresser, err := gzip.NewWriterLevel(buf, gzip.BestCompression)
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

func (GzipCompresser) Decompress(data []byte) ([]byte, error) {
	buf := bytes.NewBuffer(data)

	decompresser, err := gzip.NewReader(buf)
	if err != nil {
		return nil, err
	}

	output, err := io.ReadAll(decompresser)
	if err != nil {
		return nil, err
	}

	if err := decompresser.Close(); err != nil {
		return nil, err
	}

	return output, nil
}
