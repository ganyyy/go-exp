package generic3

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCommandSync(t *testing.T) {
	var req = MyReq{
		Name: "12313",
	}

	var m = &MyModule{
		name: "112233",
	}
	m.Init()
	m.Run()

	var err = m.AddTask(&req)
	assert.Nil(t, err)
	rsp, err := req.Resp()
	assert.Nil(t, err)

	assert.Equal(t, rsp.Name, "12313112233")

	err = m.AddTask(&MyAsyncReq{Haha: "2131231"})
	assert.Nil(t, err)
	time.Sleep(time.Second)
}

func TestParallelRun(t *testing.T) {

}
