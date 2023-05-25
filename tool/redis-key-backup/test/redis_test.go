package test

import (
	"context"
	"redis-key-backup/config"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRedisClient(t *testing.T) {

	var client = config.GetClient()

	var back = context.Background()

	t.Run("key opt", func(t *testing.T) {
		{
			ret, err := client.Set(back, "Key", "Val", 0).Result()
			t.Logf("%+v, %v", ret, err)
		}

		{
			ret, err := client.Get(back, "Key").Result()
			t.Logf("%+v, %v", ret, err)
		}

		{
			ret, err := client.Del(back, "Key").Result()
			t.Logf("%+v, %v", ret, err)
		}
	})

	t.Run("key type", func(t *testing.T) {

		var logType = func(key string) {
			ret, _ := client.Type(back, key).Result()
			t.Logf("key %v type %v", key, ret)
		}

		{
			const Key = "key"
			_, _ = client.Set(back, "Key", "Val", 0).Result()
			logType(Key)
		}

		{
			const Key = "map"
			ret, _ := client.HSet(back, Key, "k1", "v1").Result()
			t.Logf("set ret:%v", ret)
			logType(Key)
		}

		{
			const Key = "rank"
			ret, _ := client.ZAdd(back, Key, redis.Z{
				Score:  100,
				Member: "123",
			}, redis.Z{
				Score:  200,
				Member: "123131",
			}).Result()

			t.Logf("set ret:%v", ret)
			logType(Key)
		}
	})

	t.Run("nil key", func(t *testing.T) {
		const NilKey = "nil_key"
		{
			var ret, err = client.Get(context.Background(), NilKey).Result()
			assert.Equal(t, err, redis.Nil)
			assert.Equal(t, ret, "")
			t.Logf("ret:%v, err:%v", ret, err)
		}
		{
			var ret, err = client.LRange(context.Background(), NilKey, 0, -1).Result()
			t.Logf("ret:%v, err:%v", ret, err)
		}
		{
			var ret, err = client.ZRangeWithScores(context.Background(), NilKey, 0, -1).Result()
			t.Logf("ret:%v, err:%v", ret, err)
		}
		{
			var ret, err = client.HGetAll(context.Background(), NilKey).Result()
			t.Logf("ret:%v, err:%v", ret, err)
		}
	})
}

func TestRedisFunction(t *testing.T) {
	var client = config.GetClient()

	const script = `#!lua name=mylib
local function hset(keys, args)
	return "Hello Redis 7.0"
end

redis.register_function('my_hset', hset)
	`

	msg, err := client.FunctionLoad(context.Background(), script).Result()
	if redis.HasErrorPrefix(err, "Library 'mylib' already exists") {
		err = nil
	}
	t.Logf("msg:%v, err:%v", msg, err)
	require.Nil(t, err)

	ret, err := client.FCall(context.Background(), "my_hset", []string{"0"}, "hello").Result()
	t.Logf("ret:%v, err:%v", ret, err)
	require.Nil(t, err)
}

func BenchmarkFuncCall(b *testing.B) {
	var client = config.GetClient()
	for i := 0; i < b.N; i++ {
		_, _ = client.FCall(context.Background(), "my_hset", []string{"0"}, "hello").Result()
	}
}
