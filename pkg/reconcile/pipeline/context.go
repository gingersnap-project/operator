package pipeline

import (
	"context"

	"github.com/gingersnap-project/operator/pkg/kubernetes/client"
	"github.com/gingersnap-project/operator/pkg/reconcile"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var _ reconcile.Context = &ContextImpl{}

type ContextImpl struct {
	reconcile.FlowStatus
	ctx            context.Context
	client         client.Client
	log            logr.Logger
	supportedTypes map[schema.GroupVersionKind]struct{}
}

func NewContext(ctx context.Context, log logr.Logger, supportedTypes map[schema.GroupVersionKind]struct{}, client client.Client) reconcile.Context {
	return &ContextImpl{
		ctx:            ctx,
		FlowStatus:     reconcile.FlowStatus{},
		client:         client,
		log:            log,
		supportedTypes: supportedTypes,
	}
}

func (i *ContextImpl) Ctx() context.Context {
	return i.ctx
}

func (i *ContextImpl) Client() client.Client {
	return i.client
}

func (i *ContextImpl) IsTypeSupported(gvk schema.GroupVersionKind) bool {
	_, ok := i.supportedTypes[gvk]
	return ok
}

func (i *ContextImpl) Status() reconcile.FlowStatus {
	return i.FlowStatus
}

func (i *ContextImpl) Log() logr.Logger {
	return i.log
}
