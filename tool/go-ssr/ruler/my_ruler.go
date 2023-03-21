package ruler

import (
	"github.com/quasilyte/go-ruleguard/dsl"
)

func checkResponse(m dsl.Matcher) {
	m.Import("ganyyy.com/go-exp/tool/go-ssr/tools")

	m.Match("$_.InfoResponse($_)").
		Suggest("must use reture XXXResponse")
}

func checkPanicCatch(m dsl.Matcher) {
	m.Import("ganyyy.com/go-exp/tool/go-ssr/tools")
	m.Match("$p.PanicCatch()").
		Suggest("123")
}
