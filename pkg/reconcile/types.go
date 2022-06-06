package reconcile

import (
	"context"
	"fmt"
	"time"

	"github.com/engytita/engytita-operator/pkg/kubernetes/client"
	"github.com/go-logr/logr"
	_ "github.com/golang/mock/mockgen/model"
)

//go:generate mockgen -destination=reconcile_mocks.go -package=reconcile . Context,ContextProvider,Handler

// Pipeline for resource reconciliation
type Pipeline interface {
	// Process the pipeline
	// Returns true if processing should be repeated and optional error if occurred
	// important: even if error occurred it might not be needed to retry processing
	Process(i interface{}) (bool, time.Duration, error)
}

// Handler an individual stage in the pipeline
type Handler interface {
	Handle(i interface{}, ctx Context)
}

type HandlerFunc func(i interface{}, ctx Context)

func (f HandlerFunc) Handle(i interface{}, ctx Context) {
	f(i, ctx)
}

// Context of the pipeline, which is passed to each Handler
type Context interface {
	// Ctx the Pipeline's context.Context that should be passed to any functions requiring a context
	Ctx() context.Context

	// Client provides the client that should be used for all k8s resource CRUD operations
	Client() client.Client

	// Log the request logger associated with the resource
	Log() logr.Logger

	// Requeue indicates that the pipeline should stop once the current Handler has finished execution and
	// reconciliation should be requeued
	Requeue(reason error)

	// RequeueAfter indicates that the pipeline should stop once the current Handler has finished execution and
	// reconciliation should be requeued after delay time
	RequeueAfter(delay time.Duration, reason error)

	// Status the current status of a pipeline execution
	Status() FlowStatus

	// StopProcessing indicates that the pipeline should stop once the current Handler has finished execution
	StopProcessing(err error)
}

// ContextProvider returns a Context implementation for a given resource type
type ContextProvider interface {
	Get(i interface{}) (Context, error)
}

type ContextProviderFunc func(i interface{}) (Context, error)

func (f ContextProviderFunc) Get(i interface{}) (Context, error) {
	return f(i)
}

// FlowStatus Pipeline flow control
type FlowStatus struct {
	Retry bool
	Stop  bool
	Err   error
	Delay time.Duration
}

func (f *FlowStatus) String() string {
	return fmt.Sprintf("Requeue=%t, Stop=%t, Err=%s, Delay=%dms", f.Retry, f.Stop, f.Err.Error(), f.Delay.Milliseconds())
}

func (f *FlowStatus) Requeue(err error) {
	f.RequeueAfter(0, err)
}

func (f *FlowStatus) RequeueAfter(delay time.Duration, err error) {
	f.Retry = true
	f.Delay = delay
	f.StopProcessing(err)
}

func (f *FlowStatus) StopProcessing(err error) {
	f.Stop = true
	f.Err = err
}
