package pipeline

import (
	"context"

	"github.com/engytita/engytita-operator/pkg/kubernetes/client"
	"github.com/engytita/engytita-operator/pkg/reconcile"
	"github.com/go-logr/logr"
)

var _ reconcile.Context = &contextImpl{}

type contextImpl struct {
	reconcile.FlowStatus
	ctx    context.Context
	client client.Client
	log    logr.Logger
}

func NewContext(ctx context.Context, log logr.Logger, client client.Client) reconcile.Context {
	return &contextImpl{
		ctx:        ctx,
		FlowStatus: reconcile.FlowStatus{},
		client:     client,
		log:        log,
	}
}

func (i *contextImpl) Ctx() context.Context {
	return i.ctx
}

func (i *contextImpl) Client() client.Client {
	return i.client
}

func (i *contextImpl) Status() reconcile.FlowStatus {
	return i.FlowStatus
}

func (i *contextImpl) Log() logr.Logger {
	return i.log
}
