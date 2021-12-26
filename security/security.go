package security

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/xxxsen/tgblock/protos/gen/tgblock"
	"google.golang.org/protobuf/proto"
)

const (
	SigSecretId  = "Secret-Id"
	SigSecretTs  = "Secret-Ts"
	SigSecretSig = "Secret-Sig"
)

func CreateSig(secretid string, key string, expireTimestamp int64) (string, error) {
	sctx := &tgblock.SigContext{
		SecretId:  secretid,
		Timestamp: expireTimestamp,
	}
	data, err := proto.Marshal(sctx)
	if err != nil {
		return "", err
	}
	key = key + fmt.Sprintf("%d", expireTimestamp)
	data, err = EncryptByKey32(key, data)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(data), nil
}

func CheckSig(secretid string, key string, sig string, timestamp int64) (bool, error) {
	raw, err := base64.RawURLEncoding.DecodeString(sig)
	if err != nil {
		return false, err
	}
	key = key + fmt.Sprintf("%d", timestamp)
	data, err := DecryptByKey32(key, raw)
	if err != nil {
		return false, err
	}
	sctx := &tgblock.SigContext{}
	if err := proto.Unmarshal(data, sctx); err != nil {
		return false, err
	}
	if sctx.SecretId != secretid {
		return false, fmt.Errorf("token not match")
	}
	if sctx.Timestamp < time.Now().Unix() {
		return false, fmt.Errorf("token expire")
	}
	return true, nil
}
