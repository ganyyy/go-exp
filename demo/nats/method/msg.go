package method

const (
	ReqTime = "ReqTime"
	ReqVal  = "ReqVal"

	PushTime  = "PushTime"
	PushValue = "PushValue"
)

type MethReqTime struct {
	Time int64
}

type MethodRspTime struct {
	Time string
}

type MethodReqVal struct {
	Old int64
}

type MethodRspVal struct {
	New int64
}

type PushTimeParam struct {
	Time int64
}

type PushValParam struct {
	Val int64
}
