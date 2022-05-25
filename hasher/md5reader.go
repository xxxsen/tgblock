package hasher

import (
	"crypto/md5"
	"encoding/hex"
	"hash"
	"io"
	"io/ioutil"
	"os"
)

type MD5Reader struct {
	r      io.Reader
	hasher hash.Hash
	sz     int
}

func NewMD5Reader(r io.Reader) *MD5Reader {
	hasher := md5.New()
	return &MD5Reader{r: r, hasher: hasher}
}

func (s *MD5Reader) Read(buf []byte) (int, error) {
	sz, err := s.r.Read(buf)
	if sz > 0 {
		s.sz += sz
		s.hasher.Write(buf[:sz])
	}
	return sz, err
}

func (s *MD5Reader) GetSum() string {
	return hex.EncodeToString(s.hasher.Sum(nil))
}

func (s *MD5Reader) GetSize() int {
	return s.sz
}

func CalcMD5(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()
	r := NewMD5Reader(f)
	if _, err := io.Copy(ioutil.Discard, r); err != nil {
		return "", err
	}
	return r.GetSum(), nil
}
