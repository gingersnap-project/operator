// Code generated by applyconfiguration-gen. DO NOT EDIT.

package applyconfigurations

import (
	v1beta1 "github.com/gingersnap-project/operator/pkg/apis/binding/v1beta1"
	v1alpha1 "github.com/gingersnap-project/operator/api/v1alpha1"
	gingersnapprojectv1alpha1 "github.com/gingersnap-project/operator/pkg/applyconfigurations/gingersnap-project/v1alpha1"
	monitoringv1 "github.com/gingersnap-project/operator/pkg/applyconfigurations/monitoring/v1"
	servicebindingv1beta1 "github.com/gingersnap-project/operator/pkg/applyconfigurations/servicebinding/v1beta1"
	v1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
)

// ForKind returns an apply configuration type for the given GroupVersionKind, or nil if no
// apply configuration type exists for the given GroupVersionKind.
func ForKind(kind schema.GroupVersionKind) interface{} {
	switch kind {
	// Group=gingersnap-project.io, Version=v1alpha1
	case v1alpha1.SchemeGroupVersion.WithKind("Cache"):
		return &gingersnapprojectv1alpha1.CacheApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("CacheDeploymentSpec"):
		return &gingersnapprojectv1alpha1.CacheDeploymentSpecApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("CacheSpec"):
		return &gingersnapprojectv1alpha1.CacheSpecApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("CacheStatus"):
		return &gingersnapprojectv1alpha1.CacheStatusApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("DataSourceSpec"):
		return &gingersnapprojectv1alpha1.DataSourceSpecApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("DBSyncerDeploymentSpec"):
		return &gingersnapprojectv1alpha1.DBSyncerDeploymentSpecApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("EagerCacheKey"):
		return &gingersnapprojectv1alpha1.EagerCacheKeyApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("EagerCacheRule"):
		return &gingersnapprojectv1alpha1.EagerCacheRuleApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("EagerCacheRuleSpec"):
		return &gingersnapprojectv1alpha1.EagerCacheRuleSpecApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("LazyCacheKey"):
		return &gingersnapprojectv1alpha1.LazyCacheKeyApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("LazyCacheRule"):
		return &gingersnapprojectv1alpha1.LazyCacheRuleApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("LazyCacheRuleSpec"):
		return &gingersnapprojectv1alpha1.LazyCacheRuleSpecApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("LocalObjectReference"):
		return &gingersnapprojectv1alpha1.LocalObjectReferenceApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("NamespacedObjectReference"):
		return &gingersnapprojectv1alpha1.NamespacedObjectReferenceApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("ResourceQuantity"):
		return &gingersnapprojectv1alpha1.ResourceQuantityApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("Resources"):
		return &gingersnapprojectv1alpha1.ResourcesApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("ServiceBinding"):
		return &gingersnapprojectv1alpha1.ServiceBindingApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("ServiceRef"):
		return &gingersnapprojectv1alpha1.ServiceRefApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("Value"):
		return &gingersnapprojectv1alpha1.ValueApplyConfiguration{}

		// Group=monitoring.coreos.com, Version=v1
	case v1.SchemeGroupVersion.WithKind("AlertingSpec"):
		return &monitoringv1.AlertingSpecApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("Alertmanager"):
		return &monitoringv1.AlertmanagerApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("AlertmanagerConfiguration"):
		return &monitoringv1.AlertmanagerConfigurationApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("AlertmanagerEndpoints"):
		return &monitoringv1.AlertmanagerEndpointsApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("AlertmanagerSpec"):
		return &monitoringv1.AlertmanagerSpecApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("AlertmanagerStatus"):
		return &monitoringv1.AlertmanagerStatusApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("APIServerConfig"):
		return &monitoringv1.APIServerConfigApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("ArbitraryFSAccessThroughSMsConfig"):
		return &monitoringv1.ArbitraryFSAccessThroughSMsConfigApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("AttachMetadata"):
		return &monitoringv1.AttachMetadataApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("Authorization"):
		return &monitoringv1.AuthorizationApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("BasicAuth"):
		return &monitoringv1.BasicAuthApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("CommonPrometheusFields"):
		return &monitoringv1.CommonPrometheusFieldsApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("EmbeddedObjectMetadata"):
		return &monitoringv1.EmbeddedObjectMetadataApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("EmbeddedPersistentVolumeClaim"):
		return &monitoringv1.EmbeddedPersistentVolumeClaimApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("Endpoint"):
		return &monitoringv1.EndpointApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("HostAlias"):
		return &monitoringv1.HostAliasApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("MetadataConfig"):
		return &monitoringv1.MetadataConfigApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("NamespaceSelector"):
		return &monitoringv1.NamespaceSelectorApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("OAuth2"):
		return &monitoringv1.OAuth2ApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("ObjectReference"):
		return &monitoringv1.ObjectReferenceApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("PodMetricsEndpoint"):
		return &monitoringv1.PodMetricsEndpointApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("PodMetricsEndpointTLSConfig"):
		return &monitoringv1.PodMetricsEndpointTLSConfigApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("PodMonitor"):
		return &monitoringv1.PodMonitorApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("PodMonitorSpec"):
		return &monitoringv1.PodMonitorSpecApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("Probe"):
		return &monitoringv1.ProbeApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("ProberSpec"):
		return &monitoringv1.ProberSpecApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("ProbeSpec"):
		return &monitoringv1.ProbeSpecApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("ProbeTargetIngress"):
		return &monitoringv1.ProbeTargetIngressApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("ProbeTargets"):
		return &monitoringv1.ProbeTargetsApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("ProbeTargetStaticConfig"):
		return &monitoringv1.ProbeTargetStaticConfigApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("ProbeTLSConfig"):
		return &monitoringv1.ProbeTLSConfigApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("Prometheus"):
		return &monitoringv1.PrometheusApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("PrometheusCondition"):
		return &monitoringv1.PrometheusConditionApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("PrometheusRule"):
		return &monitoringv1.PrometheusRuleApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("PrometheusRuleExcludeConfig"):
		return &monitoringv1.PrometheusRuleExcludeConfigApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("PrometheusRuleSpec"):
		return &monitoringv1.PrometheusRuleSpecApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("PrometheusSpec"):
		return &monitoringv1.PrometheusSpecApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("PrometheusStatus"):
		return &monitoringv1.PrometheusStatusApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("QuerySpec"):
		return &monitoringv1.QuerySpecApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("QueueConfig"):
		return &monitoringv1.QueueConfigApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("RelabelConfig"):
		return &monitoringv1.RelabelConfigApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("RemoteReadSpec"):
		return &monitoringv1.RemoteReadSpecApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("RemoteWriteSpec"):
		return &monitoringv1.RemoteWriteSpecApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("Rule"):
		return &monitoringv1.RuleApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("RuleGroup"):
		return &monitoringv1.RuleGroupApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("Rules"):
		return &monitoringv1.RulesApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("RulesAlert"):
		return &monitoringv1.RulesAlertApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("SafeAuthorization"):
		return &monitoringv1.SafeAuthorizationApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("SafeTLSConfig"):
		return &monitoringv1.SafeTLSConfigApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("SecretOrConfigMap"):
		return &monitoringv1.SecretOrConfigMapApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("ServiceMonitor"):
		return &monitoringv1.ServiceMonitorApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("ServiceMonitorSpec"):
		return &monitoringv1.ServiceMonitorSpecApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("ShardStatus"):
		return &monitoringv1.ShardStatusApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("Sigv4"):
		return &monitoringv1.Sigv4ApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("StorageSpec"):
		return &monitoringv1.StorageSpecApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("ThanosRuler"):
		return &monitoringv1.ThanosRulerApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("ThanosRulerSpec"):
		return &monitoringv1.ThanosRulerSpecApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("ThanosRulerStatus"):
		return &monitoringv1.ThanosRulerStatusApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("ThanosSpec"):
		return &monitoringv1.ThanosSpecApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("TLSConfig"):
		return &monitoringv1.TLSConfigApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("WebSpec"):
		return &monitoringv1.WebSpecApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("WebTLSConfig"):
		return &monitoringv1.WebTLSConfigApplyConfiguration{}

		// Group=servicebinding.io, Version=v1beta1
	case v1beta1.SchemeGroupVersion.WithKind("EnvMapping"):
		return &servicebindingv1beta1.EnvMappingApplyConfiguration{}
	case v1beta1.SchemeGroupVersion.WithKind("ServiceBinding"):
		return &servicebindingv1beta1.ServiceBindingApplyConfiguration{}
	case v1beta1.SchemeGroupVersion.WithKind("ServiceBindingSecretReference"):
		return &servicebindingv1beta1.ServiceBindingSecretReferenceApplyConfiguration{}
	case v1beta1.SchemeGroupVersion.WithKind("ServiceBindingServiceReference"):
		return &servicebindingv1beta1.ServiceBindingServiceReferenceApplyConfiguration{}
	case v1beta1.SchemeGroupVersion.WithKind("ServiceBindingSpec"):
		return &servicebindingv1beta1.ServiceBindingSpecApplyConfiguration{}
	case v1beta1.SchemeGroupVersion.WithKind("ServiceBindingStatus"):
		return &servicebindingv1beta1.ServiceBindingStatusApplyConfiguration{}
	case v1beta1.SchemeGroupVersion.WithKind("ServiceBindingWorkloadReference"):
		return &servicebindingv1beta1.ServiceBindingWorkloadReferenceApplyConfiguration{}

	}
	return nil
}
