package test

import (
	"context"
	"testing"

	"github.com/go-redis/redis/v8"
)

func TestRedisClient(t *testing.T) {
	var client = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})

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
			ret, _ := client.ZAdd(back, Key, &redis.Z{
				Score:  100,
				Member: "123",
			}, &redis.Z{
				Score:  200,
				Member: "123131",
			}).Result()

			t.Logf("set ret:%v", ret)
			logType(Key)
		}
	})

}
