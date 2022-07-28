package main

import (
	"errors"

	"github.com/go-redis/redis/v8"
)

func AddToStream(key string, v ...any) error {
	var _, err = client.XAdd(ctx, &redis.XAddArgs{
		Stream: key,
		MaxLen: 1000, // 流中允许存在的最大的消息数量
		ID:     "*",  // * 代表通过redis自己生成ID
		Values: v,
	}).Result()
	return err
}

func ReadFromStream(key string) ([]redis.XMessage, error) {
	streamMsgs, err := client.XRead(ctx, &redis.XReadArgs{
		Streams: []string{key, "0-0"},
		Count:   1,
		Block:   1000,
	}).Result()
	if err != nil {
		return nil, err
	}
	if len(streamMsgs) != 1 {
		return nil, errors.New("stream not found")
	}
	var streamMsg = streamMsgs[0]
	if streamMsg.Stream != key {
		return nil, errors.New("stream not match")
	}
	return streamMsg.Messages, nil
}

func RangeStream(key string) ([]redis.XMessage, error) {
	return client.XRange(ctx, key, "-", "+").Result()
}

func StreamAck(key, group, id string) (int64, error) {
	return client.XAck(ctx, key, group, id).Result()
}
