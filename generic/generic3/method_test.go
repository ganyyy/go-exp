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

	var m = &MyModule{}
	m.Init()
	m.Run()

	var err = m.AddTask(&req)
	assert.Nil(t, err)
	var rsp = req.Resp()

	assert.Equal(t, rsp.Name, req.rsp.Name)

	err = m.AddTask(&MyAsyncReq{Haha: "2131231"})
	assert.Nil(t, err)
	time.Sleep(time.Second)
}
