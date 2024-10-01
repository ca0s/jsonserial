package compression

import (
	"bytes"
	"compress/zlib"
	"io"
)

type ZlibCompresser struct{}

func (zc *ZlibCompresser) Compress(data []byte) ([]byte, error) {
	var output []byte
	buf := bytes.NewBuffer(output)

	compresser, err := zlib.NewWriterLevel(buf, zlib.BestCompression)
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

func (zc *ZlibCompresser) Decompress(data []byte) ([]byte, error) {
	buf := bytes.NewBuffer(data)

	decompresser, err := zlib.NewReader(buf)
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
