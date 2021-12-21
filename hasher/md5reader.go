package hasher

import (
	"crypto/md5"
	"encoding/hex"
	"hash"
	"io"
)

type MD5Reader struct {
	r      io.Reader
	hasher hash.Hash
}

func NewMD5Reader(r io.Reader) *MD5Reader {
	hasher := md5.New()
	return &MD5Reader{r: r, hasher: hasher}
}

func (s *MD5Reader) Read(buf []byte) (int, error) {
	sz, err := s.r.Read(buf)
	if sz > 0 {
		s.hasher.Write(buf[:sz])
	}
	return sz, err
}

func (s *MD5Reader) GetSum() string {
	return hex.EncodeToString(s.hasher.Sum(nil))
}
