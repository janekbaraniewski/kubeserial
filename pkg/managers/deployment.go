package managers

import (
	"path/filepath"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/janekbaraniewski/kubeserial/pkg/injector"
)

func (m *Manager) CreateDeployment(cr types.NamespacedName, device types.NamespacedName) *appsv1.Deployment {
	labels := map[string]string{
		"app": m.GetName(cr.Name, device.Name),
	}
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.GetName(cr.Name, device.Name),
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:      m.GetName(cr.Name, device.Name),
					Namespace: cr.Namespace,
					Labels:    labels,
				},
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{
						{
							Name: "config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: m.GetName(cr.Name, device.Name),
									},
									Items: []corev1.KeyToPath{
										{
											Key:  filepath.Base(m.ConfigPath),
											Path: filepath.Base(m.ConfigPath),
										},
									},
								},
							},
						},
					},
					Containers: []corev1.Container{
						{
							Name:    m.GetName(cr.Name, device.Name),
							Image:   m.Image,
							Command: []string{"/bin/sh"},
							Args: []string{
								"-c",
								m.RunCmnd,
							},
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									Protocol:      corev1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "config",
									ReadOnly:  false,
									MountPath: m.ConfigPath,
									SubPath:   filepath.Base(m.ConfigPath),
								},
							},
						},
					},
				},
			},
		},
	}
	injector.AddDeviceInjector(&deployment.Spec.Template.Spec, device.Name)
	return deployment
}
