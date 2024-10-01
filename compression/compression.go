package compression

type Compression interface {
	Compress([]byte) ([]byte, error)
	Decompress([]byte) ([]byte, error)
}
