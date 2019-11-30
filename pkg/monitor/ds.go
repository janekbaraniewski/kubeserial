package monitor

import (
	appv1alpha1 "github.com/janekbaraniewski/kubeserial/pkg/apis/app/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1beta2 "k8s.io/api/apps/v1beta2"
)

func CreateDaemonSet(cr *appv1alpha1.KubeSerial) *v1beta2.DaemonSet {
	labels := map[string]string{
		"app": cr.Name + "-monitor",
	}
	return &v1beta2.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:		cr.Name + "-monitor",
			Namespace: 	cr.Namespace,
			Labels:		labels,
		},
		Spec: v1beta2.DaemonSetSpec{
			Selector: &metav1.LabelSelector {
				MatchLabels:	labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta:	metav1.ObjectMeta{
					Name:		cr.Name + "-monitor",
					Namespace:	cr.Namespace,
					Labels:		labels,
				},
				Spec: 		corev1.PodSpec{
					Volumes: 	[]corev1.Volume{
						{
							Name: 			"host-dev",
							VolumeSource: 	corev1.VolumeSource{
								HostPath:		&corev1.HostPathVolumeSource{
									Path: 			"/dev",
								},
							},
						},
						{
							Name:			"config",
							VolumeSource:	corev1.VolumeSource{
								ConfigMap:		&corev1.ConfigMapVolumeSource {
									LocalObjectReference:	corev1.LocalObjectReference {
										Name: cr.Name + "-monitor",
									},
									Items:					[]corev1.KeyToPath {
										{
											Key:	"98-devices.rules",
											Path:	"98-devices.rules",
										},
									},
								},
							},
						},
					},
					Containers: []corev1.Container{
						{
							Name:				"udev-monitor",
							Image:				"janekbaraniewski/udev-monitor:latest",
							SecurityContext:	&corev1.SecurityContext{
								Privileged:	&[]bool{true}[0],
							},
							VolumeMounts:		[]corev1.VolumeMount{
								{
									Name:		"host-dev",
									ReadOnly: 	false,
									MountPath: 	"/dev",
								},
								{
									Name:		"config",
									ReadOnly:	false,
									MountPath:	"/etc/udev/rules.d/98-devices.rules",
									SubPath:	"98-devices.rules",
								},
							},
						},
						{
							Name:				"device-monitor",
							Image:				"janekbaraniewski/udev-monitor:latest",
							Command:			[]string {"/bin/sh"},
							Args:				[]string{
								"-c",
								"./go/bin/device-monitor",
							},
							Env:				[]corev1.EnvVar{
								{
									Name: 	"OPERATOR_NAMESPACE",
									Value:	cr.Namespace,
								},
								{
									Name: 		"NODE_NAME",
									ValueFrom:	&corev1.EnvVarSource{
										FieldRef:	&corev1.ObjectFieldSelector{
											FieldPath: "spec.nodeName",
										},
									},
								},
							},
							SecurityContext:	&corev1.SecurityContext{
								Privileged:	&[]bool{true}[0],
							},
							VolumeMounts:		[]corev1.VolumeMount{
								{
									Name:		"host-dev",
									ReadOnly: 	false,
									MountPath: 	"/dev",
								},
							},
						},
					},
					ServiceAccountName:			"kubeserial",  // TODO: add separate ServiceAccount for monitor
				},
			},
		},
	}
}
