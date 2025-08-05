package tower

import "context"

type S[Req, Resp any] interface {
	Call(context.Context, Req) (Resp, error)
}

type L[Req, Resp any] interface {
	Layer(S[Req, Resp]) S[Req, Resp]
}

type Stack[Req, Resp any] struct {
	outer L[Req, Resp]
	inner L[Req, Resp]
}

func NewStack[Req, Resp any](outer L[Req, Resp], inner L[Req, Resp]) Stack[Req, Resp] {
	return Stack[Req, Resp]{outer: outer, inner: inner}
}

func (stack Stack[Req, Resp]) Layer(svc S[Req, Resp]) S[Req, Resp] {
	return stack.outer.Layer(stack.inner.Layer(svc))
}

type Identity[Req, Resp any] struct{}

func (receiver Identity[Req, Resp]) Layer(svc S[Req, Resp]) S[Req, Resp] {
	return svc
}

type ServiceBuilder[Req, Resp any] struct {
	layer L[Req, Resp]
}

func NewServiceBuilder[Req, Resp any]() ServiceBuilder[Req, Resp] {
	return ServiceBuilder[Req, Resp]{layer: Identity[Req, Resp]{}}
}

func (builder ServiceBuilder[Req, Resp]) WithLayer(layer L[Req, Resp]) ServiceBuilder[Req, Resp] {
	builder.layer = NewStack(layer, builder.layer)
	return builder
}

func (builder ServiceBuilder[Req, Resp]) Build(svc S[Req, Resp]) S[Req, Resp] {
	return builder.layer.Layer(svc)
}
