package monitor

import (
	"fmt"

	appv1alpha1 "github.com/janekbaraniewski/kubeserial/pkg/apis/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateDaemonSet(cr *appv1alpha1.KubeSerial, monitorVersion string) *appsv1.DaemonSet {
	labels := map[string]string{
		"app": cr.Name + "-monitor",
	}
	return &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-monitor",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:      cr.Name + "-monitor",
					Namespace: cr.Namespace,
					Labels:    labels,
				},
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{
						{
							Name: "host-dev",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/dev",
								},
							},
						},
						{
							Name: "config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: cr.Name + "-monitor",
									},
									Items: []corev1.KeyToPath{
										{
											Key:  "98-devices.rules",
											Path: "98-devices.rules",
										},
									},
								},
							},
						},
					},
					Containers: []corev1.Container{
						{
							Name:            "udev-monitor",
							Image:           fmt.Sprintf("janekbaraniewski/kubeserial-device-monitor:%s", monitorVersion),
							ImagePullPolicy: "Always",
							SecurityContext: &corev1.SecurityContext{
								Privileged: &[]bool{true}[0],
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "host-dev",
									ReadOnly:  false,
									MountPath: "/dev",
								},
								{
									Name:      "config",
									ReadOnly:  false,
									MountPath: "/etc/udev/rules.d/98-devices.rules",
									SubPath:   "98-devices.rules",
								},
							},
						},
						{
							Name:            "device-monitor",
							Image:           fmt.Sprintf("janekbaraniewski/kubeserial-device-monitor:%s", monitorVersion),
							ImagePullPolicy: "Always",
							Command:         []string{"/bin/sh"},
							Args: []string{
								"-c",
								"./go/bin/device-monitor",
							},
							Env: []corev1.EnvVar{
								{
									Name:  "OPERATOR_NAMESPACE",
									Value: cr.Namespace,
								},
								{
									Name: "NODE_NAME",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "spec.nodeName",
										},
									},
								},
							},
							Ports: []corev1.ContainerPort{
								{
									Name:          "metrics",
									ContainerPort: 8080,
								},
							},
							SecurityContext: &corev1.SecurityContext{
								Privileged: &[]bool{true}[0],
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "host-dev",
									ReadOnly:  false,
									MountPath: "/dev",
								},
							},
						},
					},
					ServiceAccountName: "kubeserial",
				},
			},
		},
	}
}
