package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"unsafe"

	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"

	"redis-key-backup/config"
)

var (
	ErrNotFound = errors.New("cannot implement key operation")
)

type DumpOperation interface {
	Dump(client *redis.Client, key string) (string, error)
}

type RestoreOperation interface {
	Restore(client *redis.Client, key, val string) error
}

type KeyOperation interface {
	DumpOperation
	RestoreOperation
}

var m = make(map[string]KeyOperation)

func registerKeyOperation(keyType string, operation KeyOperation) {
	m[keyType] = operation
}

// 一些暴露的接口

func GetOperation(keyType string) (KeyOperation, error) {
	var opt, found = m[keyType]
	if !found {
		return nil, fmt.Errorf("keyType %s: %w", keyType, ErrNotFound)
	}
	return opt, nil
}

func StringToBytes(str string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{str, len(str)},
	))
}

func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func GetKeyType(key string) (string, error) {
	return config.GetClient().Type(context.Background(), key).Result()
}

func KeyExists(key string) (bool, error) {
	n, err := config.GetClient().Exists(context.Background(), key).Result()
	if checkRedisError(err) != nil {
		return false, err
	}
	return n == 1, nil
}

func ExportToFile(path, data string) error {
	return os.WriteFile(path, StringToBytes(data), 0666)
}

func ReadFromFile(path string) (string, error) {
	var data, err = os.ReadFile(path)
	return BytesToString(data), err
}

type SaveStruct struct {
	OldKey  string `json:"old,omitempty"`
	KeyType string `json:"kt,omitempty"`
	Val     string `json:"val,omitempty"`
}

func (s SaveStruct) String() string {
	var bs, _ = json.Marshal(s)
	return BytesToString(bs)
}

func (s *SaveStruct) FromVal(val string) error {
	if s == nil {
		return errors.New("invalid output object")
	}
	return json.Unmarshal(StringToBytes(val), s)
}

func (s *SaveStruct) DoDump(key string) (err error) {
	var keyType string
	var exists bool
	exists, err = KeyExists(key)
	if err != nil {
		return
	}
	if !exists {
		return redis.Nil
	}

	keyType, err = GetKeyType(key)
	if err != nil {
		return
	}
	s.OldKey = key
	s.KeyType = keyType
	operation, err := GetOperation(keyType)
	if err != nil {
		return
	}
	data, err := operation.Dump(config.GetClient(), key)
	if err != nil {
		return
	}
	s.Val = data
	return
}

func (s SaveStruct) DoRestore(key string) error {
	var exists, err = KeyExists(key)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("key %v exists", key)
	}

	operation, err := GetOperation(s.KeyType)
	if err != nil {
		return err
	}
	return operation.Restore(config.GetClient(), key, s.Val)
}

func checkRedisError(err error) error {
	// fix race
	return err
}

type zSetElement struct {
	Mem   string  `json:"m"`
	Score float64 `json:"s"`
}

type zSetElements []zSetElement

type listElements []string
