package cache

import (
	"fmt"
	"time"

	"github.com/gingersnap-project/operator/api/v1alpha1"
	binding "github.com/gingersnap-project/operator/pkg/apis/binding/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const conditionWait = time.Second * 2

func ConditionReady(c *v1alpha1.Cache, ctx *Context) {
	condition := &v1alpha1.CacheCondition{
		Type:   v1alpha1.CacheConditionReady,
		Status: metav1.ConditionFalse,
	}

	notFound := func(kind, name string) {
		condition.Status = metav1.ConditionFalse
		condition.Message = fmt.Sprintf("Cache %s '%s' %s", kind, name, metav1.StatusReasonNotFound)
	}

	sb := &binding.ServiceBinding{}
	sbName := c.CacheService().DataSourceServiceBinding()
	if err := ctx.Client().Load(sbName, sb); err != nil {
		if errors.IsNotFound(err) {
			notFound("ServiceBinding", sbName)
		} else {
			ctx.RequeueAfter(conditionWait, fmt.Errorf("unable to load ServiceBinding '%s': %w", sbName, err))
			return
		}
	}

	var applicationBound bool
	for _, condition := range sb.Status.Conditions {
		if condition.Type == "Ready" {
			if condition.Status == metav1.ConditionTrue {
				applicationBound = true
			} else {
				condition.Message = fmt.Sprintf("Cache ServiceBinding '%s' not Ready: '%s'", sbName, condition.Message)
			}
			break
		}
	}
	if applicationBound {
		if c.Local() {
			ds := &appsv1.DaemonSet{}
			if err := ctx.Client().Load(c.Name, ds); client.IgnoreNotFound(err) != nil {
				ctx.Requeue(fmt.Errorf("unable to load DaemonSet for %s Ready Condition check: %w", v1alpha1.KindCache, err))
				return
			} else if err != nil {
				notFound("DaemonSet", c.Name)
			} else {
				daemonSetStatus(ds, condition)
			}
		} else {
			deployment := &appsv1.Deployment{}
			if err := ctx.Client().Load(c.Name, deployment); client.IgnoreNotFound(err) != nil {
				ctx.Requeue(fmt.Errorf("unable to load Deployment for %s Ready Condition check: %w", v1alpha1.KindCache, err))
				return
			} else if err != nil {
				notFound("Deployment", c.Name)
			} else {
				deploymentStatus(deployment, condition)
			}
		}
	}

	if c.SetCondition(condition) {
		if err := ctx.Client().UpdateStatus(c); err != nil {
			ctx.Requeue(fmt.Errorf("unable to update Ready condition: %w", err))
			return
		}
	}

	if condition.Status == metav1.ConditionFalse {
		ctx.RequeueAfter(conditionWait, nil)
	}
}

func sbRootExists(container string, spec *corev1.PodSpec) bool {
	var env []corev1.EnvVar
	for _, c := range spec.Containers {
		if c.Name == container {
			env = c.Env
			break
		}
	}
	for _, e := range env {
		if e.Name == "SERVICE_BINDING_ROOT" {
			return true
		}
	}
	return false
}

func daemonSetStatus(ds *appsv1.DaemonSet, condition *v1alpha1.CacheCondition) {
	if !sbRootExists(cacheContainer, &ds.Spec.Template.Spec) {
		condition.Status = metav1.ConditionFalse
		condition.Message = "DaemonSet ServiceBinding not bound to Pod spec"
		return
	}

	// TODO make this a configurable percentage?
	if ds.Status.NumberAvailable == ds.Status.DesiredNumberScheduled {
		condition.Status = metav1.ConditionTrue
		condition.Message = "Daemon Pod(s) running on expected number of Nodes"
	} else {
		condition.Message = fmt.Sprintf("Required DaemonSet Pods to be Available on '%d' Nodes, observed '%d'", ds.Status.DesiredNumberScheduled, ds.Status.NumberAvailable)
	}
}

func deploymentStatus(deployment *appsv1.Deployment, condition *v1alpha1.CacheCondition) {
	if !sbRootExists(cacheContainer, &deployment.Spec.Template.Spec) {
		condition.Status = metav1.ConditionFalse
		condition.Message = "Deployment ServiceBinding not bound to Pod Spec"
		return
	}

	for _, c := range deployment.Status.Conditions {
		if c.Type == appsv1.DeploymentAvailable {
			if c.Status == corev1.ConditionTrue {
				condition.Status = metav1.ConditionTrue
				condition.Message = "Expected number of Deployment pods are Ready"
			} else {
				condition.Message = fmt.Sprintf("Required Deployment '%d' pods to be Ready, observed '%d'", *deployment.Spec.Replicas, deployment.Status.ReadyReplicas)
			}
			break
		}
	}
}
