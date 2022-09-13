package pipeline

import (
	"context"

	"github.com/engytita/engytita-operator/pkg/kubernetes/client"
	"github.com/engytita/engytita-operator/pkg/reconcile"
	"github.com/go-logr/logr"
)

var _ reconcile.Context = &ContextImpl{}

type ContextImpl struct {
	reconcile.FlowStatus
	ctx    context.Context
	client client.Client
	log    logr.Logger
}

func NewContext(ctx context.Context, log logr.Logger, client client.Client) reconcile.Context {
	return &ContextImpl{
		ctx:        ctx,
		FlowStatus: reconcile.FlowStatus{},
		client:     client,
		log:        log,
	}
}

func (i *ContextImpl) Ctx() context.Context {
	return i.ctx
}

func (i *ContextImpl) Client() client.Client {
	return i.client
}

func (i *ContextImpl) Status() reconcile.FlowStatus {
	return i.FlowStatus
}

func (i *ContextImpl) Log() logr.Logger {
	return i.log
}
