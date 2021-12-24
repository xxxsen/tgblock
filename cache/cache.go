package cache

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"sync"

	"github.com/hnlq715/golang-lru/simplelru"
)

type Cache interface {
	Get(ctx context.Context, key string) (interface{}, error)
}

type SimpleCache struct {
	c    *Config
	mem  *simplelru.LRU
	file *simplelru.LRU
	lck  sync.RWMutex
}

func New(opts ...Option) (*SimpleCache, error) {
	c := &Config{}
	for _, opt := range opts {
		opt(c)
	}
	if c.caller == nil || c.codec == nil || len(c.dir) == 0 ||
		c.fileKeySize == 0 || c.memKeySize == 0 || c.valuetype == nil {
		return nil, fmt.Errorf("invalid params")
	}
	//程序重启后, 需要清理缓存目录
	dir := c.dir + "/simplecache"
	if err := os.RemoveAll(dir); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	c.dir = dir
	s := &SimpleCache{
		c: c,
	}
	s.mem, _ = simplelru.NewLRU(int(c.memKeySize), nil)
	s.file, _ = simplelru.NewLRU(int(c.fileKeySize), s.onFileClean)
	return s, nil
}

type OnCacheMissFunc func(ctx context.Context, key string) (interface{}, error)

func (s *SimpleCache) asHash(key string) string {
	v := md5.Sum([]byte(key))
	return hex.EncodeToString(v[:])
}

func (s *SimpleCache) hashToFilePath(name string) (string, string) {
	path := fmt.Sprintf("%s/%s/%s/", s.c.dir, name[:2], name[2:4])
	return path, name
}

func (s *SimpleCache) keyToFilePath(key string) (string, string) {
	name := s.asHash(key)
	return s.hashToFilePath(name)
}

func (s *SimpleCache) valueinst() interface{} {
	tp := reflect.TypeOf(s.c.valuetype)
	return reflect.New(tp).Interface()
}

func (s *SimpleCache) getFromFile(key string) (interface{}, bool) {
	path, name := s.keyToFilePath(key)
	filename := path + name
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("read cache but not found? file:%s", filename)
			return nil, false
		}
		log.Printf("read cache file fail, file:%s, err:%v", filename, err)
		return nil, false
	}
	inst := s.valueinst()
	if err := s.c.codec.Decode(raw, inst); err != nil {
		log.Printf("cache decode fail, file:%s, err:%v", filename, err)
		return nil, false
	}
	return inst, true
}

func (s *SimpleCache) get(key string) (interface{}, bool) {
	if v, ok := s.mem.Get(key); ok {
		return v, true
	}
	if _, ok := s.file.Get(s.asHash(key)); ok {
		return s.getFromFile(key)
	}
	return nil, false
}

func (s *SimpleCache) Get(ctx context.Context, key string) (interface{}, error) {
	s.lck.RLock()
	if v, ok := s.get(key); ok {
		s.lck.RUnlock()
		return v, nil
	}
	s.lck.RUnlock()
	//read from remote
	remoteValue, err := s.c.caller(ctx, key)
	if err != nil {
		return nil, err
	}
	//write to file
	hash, err := s.writeToFile(key, remoteValue)
	if err != nil {
		return nil, err
	}
	s.lck.Lock()
	defer s.lck.Unlock()
	s.mem.Add(key, remoteValue)
	s.file.Add(hash, true)
	return remoteValue, nil
}

func (s *SimpleCache) writeToFile(key string, v interface{}) (string, error) {
	raw, err := s.c.codec.Encode(v)
	if err != nil {
		return "", err
	}
	path, name := s.keyToFilePath(key)
	if err := os.MkdirAll(path, 0755); err != nil {
		return "", err
	}
	filename := path + name
	if err := ioutil.WriteFile(filename, raw, 0755); err != nil {
		return "", err
	}
	return name, nil
}

func (s *SimpleCache) onFileClean(key interface{}, value interface{}) {
	path, name := s.hashToFilePath(key.(string))
	filename := path + name
	if err := os.Remove(filename); err != nil {
		log.Printf("remove hash key fail, hash:%v, err:%v", key, err)
	}
}
