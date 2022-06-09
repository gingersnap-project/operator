package pipeline

import "github.com/engytita/engytita-operator/pkg/reconcile"

type Builder impl

func (b *Builder) WithContextProvider(ctxProvider reconcile.ContextProvider) *Builder {
	b.ctxProvider = ctxProvider
	return b
}

func (b *Builder) WithHandlers(h ...reconcile.Handler) *Builder {
	b.handlers = append(b.handlers, h...)
	return b
}

func (b *Builder) Build() reconcile.Pipeline {
	return &impl{
		handlers:    b.handlers,
		ctxProvider: b.ctxProvider,
	}
}
