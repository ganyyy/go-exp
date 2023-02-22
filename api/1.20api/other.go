package api

func Int64(v int64) *int64 {
	return &v
}

type SomeStruct struct {
	A, B, C, D, E, F, G int
}

func SomeStructPtr(data SomeStruct) *SomeStruct {
	return &data
}

//go:noinline
func monitorRPC(req *SomeStruct, rsp *int) {
	*rsp = req.A
}

//go:noinline
func CallRpc() {
	var req SomeStruct
	var rsp int
	monitorRPC(&req, &rsp)
	_ = rsp
}

func TestGetAddr() {
	var v int64
	var s SomeStruct
	pv := Int64(v)
	ps := SomeStructPtr(s)
	println(pv, ps)
}
