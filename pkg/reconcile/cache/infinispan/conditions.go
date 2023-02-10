package infinispan

import (
	"fmt"
	"time"

	"github.com/gingersnap-project/operator/api/v1alpha1"
	binding "github.com/gingersnap-project/operator/pkg/apis/binding/v1beta1"
	"github.com/gingersnap-project/operator/pkg/reconcile/cache/context"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const conditionWait = time.Second * 2

func ConditionReady(c *v1alpha1.Cache, ctx *context.Context) {
	condition := c.Condition(v1alpha1.CacheConditionReady)

	update := func(status metav1.ConditionStatus, msg string) {
		condition.Status = status
		condition.Message = msg
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
		if condition.Type == binding.ServiceBindingConditionReady {
			if condition.Status == metav1.ConditionTrue {
				applicationBound = true
			} else {
				update(
					metav1.ConditionFalse,
					fmt.Sprintf("Cache ServiceBinding '%s' not Ready: '%s'", sbName, condition.Message),
				)
			}
			break
		}
	}
	if applicationBound {
		if c.Local() {
			ds := &appsv1.DaemonSet{}
			if err := ctx.Client().Load(c.Name, ds); client.IgnoreNotFound(err) != nil {
				ctx.Requeue(fmt.Errorf("unable to load DaemonSet for Available Condition check: %w", err))
				return
			} else if err != nil {
				notFound("DaemonSet", c.Name)
			} else {
				if ds.Status.NumberReady == ds.Status.DesiredNumberScheduled {
					update(
						metav1.ConditionTrue,
						"Expected number of DaemonSet pods are Ready",
					)
				} else {
					update(
						metav1.ConditionFalse,
						fmt.Sprintf("Required DaemonSet '%d' pods to be Ready, observed '%d'", ds.Status.DesiredNumberScheduled, ds.Status.NumberReady),
					)
				}
			}
		} else {
			deployment := &appsv1.Deployment{}
			if err := ctx.Client().Load(c.Name, deployment); client.IgnoreNotFound(err) != nil {
				ctx.Requeue(fmt.Errorf("unable to load Deployment for Available Condition check: %w", err))
				return
			} else if err != nil {
				notFound("Deployment", c.Name)
			} else {
				var deploymentAvailable bool
				for _, condition := range deployment.Status.Conditions {
					if condition.Type == appsv1.DeploymentAvailable {
						deploymentAvailable = condition.Status == corev1.ConditionTrue
						break
					}
				}

				if deploymentAvailable {
					update(
						metav1.ConditionTrue,
						"Expected number of Deployment pods are Ready",
					)
				} else {
					update(
						metav1.ConditionFalse,
						fmt.Sprintf("Required Deployment '%d' pods to be Ready, observed '%d'", deployment.Spec.Replicas, deployment.Status.ReadyReplicas),
					)
				}
			}
		}
	}

	c.SetCondition(condition)
	if err := ctx.Client().UpdateStatus(c); err != nil {
		ctx.Requeue(fmt.Errorf("unable to update Available condition: %w", err))
		return
	}

	if condition.Status == metav1.ConditionFalse {
		ctx.RequeueAfter(conditionWait, nil)
	}
}
