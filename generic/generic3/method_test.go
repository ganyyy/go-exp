package generic3

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type MyPanicTask struct {
	BaseSyncMethod[*MyModule, MyResp]
}

func (m *MyPanicTask) Do() {
	var rsp = m.GenResp()
	rsp.Age = 100
	panic("rand panic!")
}

func TestCommandSync(t *testing.T) {
	var req = MyReq{
		Name: "12313",
	}

	var m = &MyModule{
		name: "112233",
	}
	m.Init()
	m.Run()

	m.AddTask(&req)
	rsp, err := req.Resp()
	assert.Nil(t, err)

	assert.Equal(t, rsp.Name, "12313112233")

	var req2 MyPanicTask
	m.AddTask(&req2)
	rsp, err = req2.Resp()
	assert.NotNil(t, err)
	t.Logf("panic %+v", rsp)

	m.AddTask(&MyAsyncReq{Lala: "2131231"})
	time.Sleep(time.Second)
}

func TestParallelRun(t *testing.T) {

}
