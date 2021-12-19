package processor

import (
	"crypto/md5"
	"encoding/hex"
	"hash"
	"io"
)

type ShaReader struct {
	r      io.Reader
	hasher hash.Hash
}

func NewShaReader(r io.Reader) *ShaReader {
	hasher := md5.New()
	return &ShaReader{r: r, hasher: hasher}
}

func (s *ShaReader) Read(buf []byte) (int, error) {
	sz, err := s.r.Read(buf)
	if sz > 0 {
		s.hasher.Write(buf[:sz])
	}
	return sz, err
}

func (s *ShaReader) GetSum() string {
	return hex.EncodeToString(s.hasher.Sum(nil))
}
