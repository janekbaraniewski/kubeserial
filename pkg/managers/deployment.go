package managers

import (
	"path/filepath"

	appv1alpha1 "github.com/janekbaraniewski/kubeserial/pkg/apis/app/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	v1beta2 "k8s.io/api/apps/v1beta2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (m *Manager)CreateDeployment(cr *appv1alpha1.KubeSerial, device *appv1alpha1.Device) *v1beta2.Deployment {
	labels := map[string]string{
		"app": m.GetName(cr.Name, device.Name),
	}
	return &v1beta2.Deployment {
		ObjectMeta:	metav1.ObjectMeta {
			Name: 	m.GetName(cr.Name, device.Name),
			Namespace:	cr.Namespace,
			Labels:		labels,
		},
		Spec:		v1beta2.DeploymentSpec {
			Selector:	&metav1.LabelSelector {
				MatchLabels:	labels,
			},
			Template:	corev1.PodTemplateSpec {
				ObjectMeta:	metav1.ObjectMeta{
					Name:		m.GetName(cr.Name, device.Name),
					Namespace:	cr.Namespace,
					Labels:		labels,
				},
				Spec: 		corev1.PodSpec{
					Volumes: 		[]corev1.Volume{
						{
							Name:			"config",
							VolumeSource:	corev1.VolumeSource{
								ConfigMap:		&corev1.ConfigMapVolumeSource {
									LocalObjectReference:	corev1.LocalObjectReference{
										Name: 		m.GetName(cr.Name, device.Name),
									},
									Items:					[]corev1.KeyToPath {
										{
											Key:	filepath.Base(m.ConfigPath),
											Path:	filepath.Base(m.ConfigPath),
										},
									},
								},
							},
						},
					},
					Containers: 	[]corev1.Container{
						{
							Name:				m.GetName(cr.Name, device.Name),
							Image:				m.Image,
							Command:			[]string {"/bin/sh"},
							Args:				[]string {
								"-c",
								"socat pty,wait-slave,link=/dev/device,perm=0660,group=tty tcp:" + cr.Name + "-" + device.Name + "-gateway:3333 & " + m.RunCmnd,  // TODO: make init container
							},
							Ports:				[]corev1.ContainerPort{
								{
									Name: 			"http",
									Protocol:		corev1.ProtocolTCP,
									ContainerPort:	80,
								},
							},
							VolumeMounts:		[]corev1.VolumeMount{
								{
									Name:		"config",
									ReadOnly:	false,
									MountPath:	m.ConfigPath,
									SubPath:	filepath.Base(m.ConfigPath),
								},
							},
						},
					},
				},
			},
		},
	}
}
