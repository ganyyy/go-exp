package db_test

import (
	"errors"
	"go-exp/mock/db"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestGetFromDB(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	m := db.NewMockDB(ctrl)
	// 如果key == tom, 会走到这个测试用例中
	m.EXPECT().Get(gomock.Eq("tom")).Return(100, errors.New("no exist"))
	// 如果key == hello, 会走到这个用例中
	m.EXPECT().Get(gomock.Eq("hello")).DoAndReturn(func(key string) (int, error) {
		t.Logf("get key:%v", key)
		if key == "hello" {
			return 100, nil
		} else {
			return 0, errors.New("cannot found key")
		}
	})
	// m.EXPECT().Get(gomock.Any()).Return(-1, errors.New("1234"))

	if v := db.GetFromDB(m, "tom"); v != -1 {
		t.Logf("expected -1, but got:%v", v)
		t.FailNow()
	}

	if v := db.GetFromDB(m, "hello"); v != 100 {
		t.Logf("expected 100, but got:%v", v)
		t.FailNow()
	}
}
