package gateway

import (
	"strings"

	appv1alpha1 "github.com/janekbaraniewski/kubeserial/pkg/apis/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateDeployment(device *appv1alpha1.SerialDevice) *appsv1.Deployment {
	labels := map[string]string{
		"app": device.Name + "-gateway",
	}
	name := strings.ToLower(device.Name + "-gateway")
	return &appsv1.Deployment{ // TODO: add TCP probes
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: device.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: device.Namespace,
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
										Name: name,
									},
									Items: []corev1.KeyToPath{
										{
											Key:  "ser2net.conf",
											Path: "ser2net.conf",
										},
									},
								},
							},
						},
					},
					Containers: []corev1.Container{
						{
							Name:    "kubeserial-gateway",
							Image:   "janekbaraniewski/ser2net:latest",
							Command: []string{"/bin/sh"},
							Args: []string{
								"-c",
								"ser2net && sleep inf",
							},
							SecurityContext: &corev1.SecurityContext{
								Privileged: &[]bool{true}[0],
							},
							Ports: []corev1.ContainerPort{
								{
									Name:          "ser2net",
									Protocol:      corev1.ProtocolTCP,
									ContainerPort: 3333,
								},
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
									MountPath: "/etc/ser2net.conf",
									SubPath:   "ser2net.conf",
								},
							},
						},
					},
					NodeSelector: map[string]string{
						"kubernetes.io/hostname": device.Status.NodeName,
					},
				},
			},
		},
	}
}
