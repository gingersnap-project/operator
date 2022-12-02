package e2e

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var mysqlLabels = map[string]string{"app": "mysql"}

var MysqlConfigMap = &corev1.ConfigMap{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "mysql",
		Namespace: Namespace,
	},
	Data: map[string]string{"setup.sql": `
create schema debezium;
create table debezium.customer(id int not null, fullname varchar(255), email varchar(255), constraint primary key (id));
create table debezium.car_model(id int not null, model varchar(255), brand varchar(255), constraint primary key (id));
GRANT SELECT, RELOAD, SHOW DATABASES, REPLICATION SLAVE, REPLICATION CLIENT ON *.* TO 'gingersnap_user';
`},
}

var MysqlDeployment = &appsv1.Deployment{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "mysql",
		Namespace: Namespace,
		Labels:    mysqlLabels,
	},
	Spec: appsv1.DeploymentSpec{
		Selector: &metav1.LabelSelector{
			MatchLabels: mysqlLabels,
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: mysqlLabels,
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{{
					Name:  "mysql",
					Image: "mysql:latest",
					Args:  []string{"--default-authentication-plugin=mysql_native_password"},
					Env: []corev1.EnvVar{
						{Name: "MYSQL_ROOT_PASSWORD", Value: "root"},
						{Name: "MYSQL_USER", Value: "gingersnap_user"},
						{Name: "MYSQL_PASSWORD", Value: "password"},
					},
					Ports: []corev1.ContainerPort{{
						ContainerPort: 3306,
					}},
					VolumeMounts: []corev1.VolumeMount{{
						Name:      "init-db",
						MountPath: "/docker-entrypoint-initdb.d",
						ReadOnly:  true,
					}},
				}},
				Volumes: []corev1.Volume{{
					Name: "init-db",
					VolumeSource: corev1.VolumeSource{
						ConfigMap: &corev1.ConfigMapVolumeSource{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: MysqlConfigMap.Name,
							},
						},
					},
				}},
			},
		},
	},
}

var MysqlService = &corev1.Service{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "mysql",
		Namespace: Namespace,
	},
	Spec: corev1.ServiceSpec{
		Type: corev1.ServiceTypeClusterIP,
		Ports: []corev1.ServicePort{{
			Port:     3306,
			Protocol: corev1.ProtocolTCP,
		}},
		Selector: MysqlDeployment.Labels,
	},
}

var MysqlConnectionSecret = &corev1.Secret{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "db-credentials",
		Namespace: Namespace,
	},
	StringData: map[string]string{
		"host":     fmt.Sprintf("%s.%s.svc.cluster.local", MysqlService.Name, MysqlService.Namespace),
		"port":     "3306",
		"user":     "gingersnap_user",
		"password": "password",
	},
}
