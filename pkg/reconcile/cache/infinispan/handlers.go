package infinispan

import (
	"fmt"

	"github.com/engytita/engytita-operator/api/v1alpha1"
	"github.com/engytita/engytita-operator/pkg/reconcile"
	apicorev1 "k8s.io/api/core/v1"
	appsv1 "k8s.io/client-go/applyconfigurations/apps/v1"
	corev1 "k8s.io/client-go/applyconfigurations/core/v1"
	metav1 "k8s.io/client-go/applyconfigurations/meta/v1"
)

const (
	containerName = "infinispan"
)

var (
	selectorLabels = map[string]string{
		"app": "infinispan",
	}
)

func Service(c *v1alpha1.Cache, ctx reconcile.Context) {
	service := corev1.
		Service(c.Name, c.Namespace).
		WithOwnerReferences(
			ctx.Client().OwnerReference(),
		).
		WithSpec(
			corev1.ServiceSpec().
				WithClusterIP(apicorev1.ClusterIPNone).
				WithType(apicorev1.ServiceTypeClusterIP).
				WithSelector(selectorLabels).
				WithPorts(
					corev1.ServicePort().WithName("infinispan").WithPort(11222),
				),
		)

	if err := ctx.Client().Apply(service); err != nil {
		ctx.Requeue(fmt.Errorf("unable to apply Infinispan Service: %w", err))
	}
}

func ConfigMap(c *v1alpha1.Cache, ctx reconcile.Context) {
	config := `
   <infinispan
     xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
     xsi:schemaLocation="urn:infinispan:config:13.0 https://infinispan.org/schemas/infinispan-config-13.0.xsd
                           urn:infinispan:server:13.0 https://infinispan.org/schemas/infinispan-server-13.0.xsd"
     xmlns="urn:infinispan:config:13.0"
     xmlns:server="urn:infinispan:server:13.0">

     <cache-container name="default" statistics="true">
         <local-cache name="airports" />
         <transport cluster="${infinispan.cluster.name:cluster}" stack="${infinispan.cluster.stack:tcp}" node-name="${infinispan.node.name:}"/>
         <security>
           <authorization/>
         </security>
     </cache-container>

     <server xmlns="urn:infinispan:server:13.0">
         <interfaces>
           <interface name="public">
               <inet-address value="${infinispan.bind.address:127.0.0.1}"/>
           </interface>
         </interfaces>

         <socket-bindings default-interface="public" port-offset="${infinispan.socket.binding.port-offset:0}">
           <socket-binding name="default" port="${infinispan.bind.port:11222}"/>
           <socket-binding name="memcached" port="11221"/>
         </socket-bindings>

         <security>
           <credential-stores>
               <credential-store name="credentials" path="credentials.pfx">
                 <clear-text-credential clear-text="secret"/>
               </credential-store>
           </credential-stores>
           <security-realms>
               <security-realm name="default">
                 <!-- Uncomment to enable TLS on the realm -->
                 <!-- server-identities>
                     <ssl>
                       <keystore path="application.keystore"
                                 password="password" alias="server"
                                 generate-self-signed-certificate-host="localhost"/>
                     </ssl>
                 </server-identities-->
                 <properties-realm groups-attribute="Roles">
                     <user-properties path="users.properties"/>
                     <group-properties path="groups.properties"/>
                 </properties-realm>
               </security-realm>
           </security-realms>
         </security>

         <endpoints socket-binding="default" security-realm="default" />
     </server>
   </infinispan>
`
	cm := corev1.
		ConfigMap(c.Name, c.Namespace).
		WithOwnerReferences(ctx.Client().OwnerReference()).
		WithData(
			map[string]string{
				"infinispan.xml": config,
			},
		)

	if err := ctx.Client().Apply(cm); err != nil {
		ctx.Requeue(fmt.Errorf("unable to apply Infinispan ConfigMap: %w", err))
	}
}

func DaemonSet(c *v1alpha1.Cache, ctx reconcile.Context) {
	ds := appsv1.
		DaemonSet(c.Name, c.Namespace).
		WithOwnerReferences(ctx.Client().OwnerReference()).
		WithSpec(appsv1.DaemonSetSpec().
			WithSelector(
				metav1.LabelSelector().WithMatchLabels(selectorLabels),
			).
			WithTemplate(corev1.PodTemplateSpec().
				WithName(containerName).
				WithLabels(selectorLabels).
				WithSpec(corev1.PodSpec().
					WithContainers(corev1.Container().
						WithName(containerName).
						WithImage("quay.io/infinispan/server:13.0").
						WithArgs("-c", "/config/infinispan.xml").
						WithPorts(
							corev1.ContainerPort().WithContainerPort(11222),
						).
						WithEnv(
							corev1.EnvVar().WithName("USER").WithValue("admin"),
							corev1.EnvVar().WithName("PASS").WithValue("password"),
						).
						WithVolumeMounts(
							corev1.VolumeMount().WithName("config").WithMountPath("/config").WithReadOnly(true),
						),
					).
					WithVolumes(corev1.Volume().
						WithName("config").
						WithConfigMap(
							corev1.ConfigMapVolumeSource().WithName(c.Name),
						),
					),
				),
			),
		)
	if err := ctx.Client().Apply(ds); err != nil {
		ctx.Requeue(fmt.Errorf("unable to apply Infinispan DaemonSet: %w", err))
	}
}
