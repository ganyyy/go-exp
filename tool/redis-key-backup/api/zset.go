package api

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"

	"redis-key-backup/config"
)

func init() {
	registerKeyOperation(config.TypeZSet, zSetOperation{})
}

type zSetOperation struct{}

func (z zSetOperation) Dump(client *redis.Client, key string) (string, error) {
	var ret, err = client.ZRangeWithScores(context.Background(), key, 0, -1).Result()
	if checkRedisError(err) != nil {
		return "", err
	}
	var elemes = make([]zSetElement, 0, len(ret))
	for _, v := range ret {
		elemes = append(elemes, zSetElement{
			Mem:   v.Member.(string),
			Score: v.Score,
		})
	}
	bs, err := json.Marshal(elemes)
	return BytesToString(bs), err
}

func (z zSetOperation) Restore(client *redis.Client, key, val string) error {
	var elements zSetElements
	err := json.Unmarshal(StringToBytes(val), &elements)
	if err != nil {
		return err
	}
	var args = make([]*redis.Z, 0, len(elements))
	for _, ele := range elements {
		args = append(args, &redis.Z{
			Score:  ele.Score,
			Member: ele.Mem,
		})
	}

	n, err := client.ZAdd(context.Background(), key, args...).Result()
	if err != nil {
		return err
	}
	if int(n) != len(elements) {
		return errors.New(fmt.Sprintf("data elements: %v, zadd elements:%v", len(elements), n))
	}
	return nil
}
