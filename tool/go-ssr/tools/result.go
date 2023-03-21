package tools

type CodeType uint32

const (
	NormalCode CodeType = iota
	InfoCode
	WarnCode
	ErrorCode
)

type IResponse interface {
	Type() CodeType
	Msg() string
}

type responseWrap struct {
	codeType CodeType
	msg      string
}

// Msg implements IResponse
func (r *responseWrap) Msg() string {
	return r.msg
}

// Type implements IResponse
func (r *responseWrap) Type() CodeType {
	return r.codeType
}

var (
	_ IResponse = &responseWrap{}
)

func Response(msg string) IResponse {
	return &responseWrap{
		codeType: NormalCode,
		msg:      msg,
	}
}

func InfoResponse(msg string) IResponse {
	return &responseWrap{
		codeType: InfoCode,
		msg:      msg,
	}
}

func WarnResponse(msg string) IResponse {
	return &responseWrap{
		codeType: WarnCode,
		msg:      msg,
	}
}

func ErrorResponse(msg string) IResponse {
	return &responseWrap{
		codeType: ErrorCode,
		msg:      msg,
	}
}
