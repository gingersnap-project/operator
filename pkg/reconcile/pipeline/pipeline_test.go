package pipeline_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/engytita/engytita-operator/pkg/reconcile"
	"github.com/engytita/engytita-operator/pkg/reconcile/pipeline"
	"github.com/go-logr/logr"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestBuilder(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Pipeline Suite")
}

var _ = Describe("Pipeline", func() {
	var (
		mockCtrl   *gomock.Controller
		ctx        *reconcile.MockContext
		defHandler = func() *reconcile.MockHandler {
			return reconcile.NewMockHandler(mockCtrl)
		}
		resource = struct{}{}
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		ctx = reconcile.NewMockContext(mockCtrl)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("should not retry if processing was successful", func() {
		h1 := defHandler()
		h1.EXPECT().Handle(resource, ctx)
		h2 := defHandler()
		h2.EXPECT().Handle(resource, ctx)
		builder := &pipeline.Builder{}
		p := builder.
			WithContextProvider(&ctxProvider{ctx: ctx}).
			WithHandlers(h1, h2).
			Build()

		ctx.EXPECT().Status().Return(reconcile.FlowStatus{}).Times(2)

		retry, delay, err := p.Process(resource)
		Expect(err).NotTo(HaveOccurred())
		Expect(delay).To(Equal(time.Duration(0)))
		Expect(retry).To(BeFalse())
	})

	It("should stop processing if retry requested and propagate error back to caller", func() {
		err := errors.New("foo")

		h1 := defHandler()
		h1.EXPECT().Handle(resource, ctx)
		h2 := func(resource interface{}, c reconcile.Context) {
			c.Requeue(err)
		}
		h3 := defHandler()
		builder := &pipeline.Builder{}
		p := builder.
			WithContextProvider(&ctxProvider{ctx: ctx}).
			WithHandlers(h1, reconcile.HandlerFunc(h2), h3).
			Build()

		ctx.EXPECT().Requeue(err)
		ctx.EXPECT().Status().Return(reconcile.FlowStatus{})
		ctx.EXPECT().Status().Return(reconcile.FlowStatus{Retry: true, Stop: true, Err: err})

		retry, delay, err := p.Process(resource)
		Expect(err).To(Equal(err))
		Expect(delay).To(Equal(time.Duration(0)))
		Expect(retry).To(BeTrue())
	})

	It("should retry processing if panic occurred in a handler", func() {
		err := fmt.Errorf("panic occurred: %v", "foo")

		h1 := defHandler()
		h1.EXPECT().Handle(resource, ctx)
		h2 := func(resource interface{}, c reconcile.Context) {
			panic("foo")
		}
		h3 := defHandler()
		builder := &pipeline.Builder{}
		p := builder.
			WithContextProvider(&ctxProvider{ctx: ctx}).
			WithHandlers(h1, reconcile.HandlerFunc(h2), h3).
			Build()

		ctx.EXPECT().Log().Return(logr.Discard())
		ctx.EXPECT().Requeue(err)
		ctx.EXPECT().Status().Return(reconcile.FlowStatus{})
		ctx.EXPECT().Status().Return(reconcile.FlowStatus{Retry: true, Stop: true, Err: err})

		retry, delay, err := p.Process(resource)
		Expect(err).To(Equal(err))
		Expect(delay).To(Equal(time.Duration(0)))
		Expect(retry).To(BeTrue())
	})

	It("should stop without retry and error and propagate that back to caller", func() {
		h1 := defHandler()
		h1.EXPECT().Handle(resource, ctx)
		h2 := func(resource interface{}, c reconcile.Context) {
			c.StopProcessing(nil)
		}
		h3 := defHandler()
		builder := &pipeline.Builder{}
		p := builder.
			WithContextProvider(&ctxProvider{ctx: ctx}).
			WithHandlers(h1, reconcile.HandlerFunc(h2), h3).
			Build()

		ctx.EXPECT().StopProcessing(nil)
		ctx.EXPECT().Status().Return(reconcile.FlowStatus{})
		ctx.EXPECT().Status().Return(reconcile.FlowStatus{Retry: false, Stop: true, Err: nil})

		retry, delay, err := p.Process(resource)
		Expect(err).NotTo(HaveOccurred())
		Expect(delay).To(Equal(time.Duration(0)))
		Expect(retry).To(BeFalse())
	})

	It("should retry processing if panic occurs when calling context provider", func() {
		h1 := defHandler()
		provider := reconcile.NewMockContextProvider(mockCtrl)
		provider.EXPECT().Get(resource).DoAndReturn(func(b interface{}) { panic("foo") })
		builder := &pipeline.Builder{}
		p := builder.
			WithContextProvider(provider).
			WithHandlers(h1).
			Build()

		retry, delay, err := p.Process(resource)
		Expect(err).To(Equal(fmt.Errorf("panic occurred: %v", "foo")))
		Expect(delay).To(Equal(time.Duration(0)))
		Expect(retry).To(BeTrue())
	})

	It("should retry processing after requested delay", func() {
		h1 := func(resource interface{}, c reconcile.Context) {
			c.RequeueAfter(time.Second, nil)
		}
		builder := &pipeline.Builder{}
		p := builder.
			WithContextProvider(&ctxProvider{ctx: ctx}).
			WithHandlers(reconcile.HandlerFunc(h1)).
			Build()

		ctx.EXPECT().RequeueAfter(time.Second, nil)
		ctx.EXPECT().Status().Return(reconcile.FlowStatus{Retry: true, Stop: true, Err: nil, Delay: time.Second})

		retry, delay, err := p.Process(resource)
		Expect(err).NotTo(HaveOccurred())
		Expect(delay).To(Equal(time.Second))
		Expect(retry).To(BeTrue())
	})
})

var _ reconcile.ContextProvider = &ctxProvider{}

type ctxProvider struct {
	ctx reconcile.Context
}

func (c *ctxProvider) Get(_ interface{}) (reconcile.Context, error) {
	return c.ctx, nil
}
