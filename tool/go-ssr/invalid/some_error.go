package invalid

import t "ganyyy.com/go-exp/tool/go-ssr/tools"

func RunLogic() t.IResponse {

	t.WarnResponse("123")
	t.InfoResponse("123")
	t.ErrorResponse("456")

	return t.Response("123")

}

func NotCatch() {
	defer t.PanicCatch()
}
