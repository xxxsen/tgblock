package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

func makeKey(sz int, key string) []byte {
	buf := make([]byte, sz)
	copy(buf, []byte(key))
	return buf
}

func EncryptByKey32(skey string, in []byte) ([]byte, error) {
	if len(skey) == 0 {
		return in, nil
	}
	key := makeKey(32, skey)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("create cipher fail, err:%v", err)
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("create gcm failed")
	}
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("read nonce fail, err:%v", err)
	}
	ciphertext := aesGCM.Seal(nonce, nonce, in, nil)
	return ciphertext, nil
}

func DecryptByKey32(skey string, in []byte) ([]byte, error) {
	if len(skey) == 0 {
		return in, nil
	}
	key := makeKey(32, skey)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("create cipher fail, err:%v", err)
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("create gcm fail, err:%v", err)
	}
	nonceSize := aesGCM.NonceSize()
	if nonceSize > len(in) {
		return nil, fmt.Errorf("invalid input len, nonce size:%d, len(in):%d", nonceSize, len(in))
	}
	nonce, ciphertext := in[:nonceSize], in[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("decode fail, err:%v", err)
	}
	return plaintext, nil
}
