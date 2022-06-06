package pipeline

import "github.com/engytita/engytita-operator/pkg/reconcile"

type builder impl

func Builder() *builder {
	return &builder{}
}

func (b *builder) WithContextProvider(ctxProvider reconcile.ContextProvider) *builder {
	b.ctxProvider = ctxProvider
	return b
}

func (b *builder) WithHandlers(h ...reconcile.Handler) *builder {
	b.handlers = append(b.handlers, h...)
	return b
}

func (b *builder) Build() reconcile.Pipeline {
	return &impl{
		handlers:    b.handlers,
		ctxProvider: b.ctxProvider,
	}
}
