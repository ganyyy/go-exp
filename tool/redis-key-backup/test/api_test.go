package test

import (
	"context"
	"encoding/json"
	"math/rand"
	"redis-key-backup/api"
	"redis-key-backup/config"
	"reflect"
	"strconv"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

const (
	Num = 10
)

func loopNum(f func(i int)) {
	for i := 0; i < Num; i++ {
		f(i)
	}
}

type testInterface interface {
	Prepare() error
	DoDump() (string, error)
	DoRestore(val string) (api.SaveStruct, error)
	After() error
	Equal(tt *testing.T, v1, v2 string) bool
}

type testBase struct {
	key  string
	rKey string
}

func (testBase) Prepare() error {
	return nil
}
func (t testBase) DoDump() (string, error) {
	var save api.SaveStruct
	err := save.DoDump(t.key)
	return save.String(), err
}
func (t testBase) DoRestore(val string) (ret api.SaveStruct, err error) {
	if err = json.Unmarshal(api.StringToBytes(val), &ret); err != nil {
		return
	}
	err = ret.DoRestore(t.rKey)
	if err != nil {
		return
	}
	err = ret.DoDump(t.rKey)
	if err != nil {
		return
	}
	return ret, nil
}
func (t testBase) After() error {
	_, err := config.GetClient().Del(context.Background(), t.key).Result()
	_, err = config.GetClient().Del(context.Background(), t.rKey).Result()
	return err
}
func (t testBase) Equal(tt *testing.T, v1, v2 string) bool {
	return assert.True(tt, v1 == v2)
}

type testString struct{ testBase }

func (t testString) Prepare() error {
	_, err := config.GetClient().Set(context.Background(), t.key, "rand value", 0).Result()
	return err
}

type testHmap struct{ testBase }

func (t testHmap) Prepare() error {
	const Num = 100

	var args = make([]interface{}, 0, Num*2)
	loopNum(func(i int) {
		args = append(args, "key"+strconv.Itoa(i), "val"+strconv.Itoa(rand.Int()%10))
	})
	_, err := config.GetClient().HMSet(context.Background(), t.key, args...).Result()
	return err
}
func (t testHmap) Equal(tt *testing.T, v1, v2 string) bool {
	var m1, m2 = make(map[string]string), make(map[string]string)
	err1 := json.Unmarshal(api.StringToBytes(v1), &m1)
	err2 := json.Unmarshal(api.StringToBytes(v2), &m2)
	assert.Nil(tt, err1)
	assert.Nil(tt, err2)
	return assert.Equal(tt, m1, m2)
}

type testZSet struct{ testBase }

func (t testZSet) Prepare() error {
	var args = make([]*redis.Z, 0, Num)
	loopNum(func(i int) {
		args = append(args, &redis.Z{
			Score:  float64(rand.Int() % 100),
			Member: "mem" + strconv.Itoa(i),
		})
	})
	_, err := config.GetClient().ZAdd(context.Background(), t.key, args...).Result()
	return err
}

type testList struct{ testBase }

func (t testList) Prepare() error {
	var args = make([]interface{}, 0, Num)
	loopNum(func(i int) {
		args = append(args, "list_ele"+strconv.Itoa(i))
	})
	_, err := config.GetClient().RPush(context.Background(), t.key, args...).Result()
	return err
}

func TestApiDump(t *testing.T) {
	t.Run("dump nil", func(t *testing.T) {
		const (
			NilKey = "nil key"
		)
		var saveStruct api.SaveStruct
		var err = saveStruct.DoDump(NilKey)
		assert.Equal(t, err, redis.Nil)
	})

	var testCases = []testInterface{
		testString{testBase: testBase{
			key:  "key1",
			rKey: "key222",
		}},
		testHmap{testBase: testBase{
			key:  "map1",
			rKey: "map22",
		}},
		testZSet{testBase: testBase{
			key:  "zset1",
			rKey: "zset22",
		}},
		testList{testBase: testBase{
			key:  "list1",
			rKey: "list22",
		}},
	}

	for _, tc := range testCases {
		t.Run(reflect.TypeOf(tc).Name(), func(t *testing.T) {
			err := tc.Prepare()
			assert.Nil(t, err)
			str, err := tc.DoDump()
			assert.Nil(t, err)
			var saveTemp api.SaveStruct
			err = json.Unmarshal(api.StringToBytes(str), &saveTemp)
			assert.Nil(t, err)
			restoreTmp, err := tc.DoRestore(str)
			assert.Nil(t, err)
			err = tc.After()
			assert.Nil(t, err)
			assert.True(t, tc.Equal(t, saveTemp.Val, restoreTmp.Val))
			t.Logf("save   :%+v", saveTemp)
			t.Logf("restore:%+v", restoreTmp)
		})
	}
}
