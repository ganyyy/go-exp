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

	m.AddTask(&req)
	rsp, err := req.Resp()
	assert.Nil(t, err)

	assert.Equal(t, rsp.Name, "12313112233")

	m.AddTask(&MyAsyncReq{Lala: "2131231"})
	assert.Nil(t, err)
	time.Sleep(time.Second)
}

func TestParallelRun(t *testing.T) {

}
