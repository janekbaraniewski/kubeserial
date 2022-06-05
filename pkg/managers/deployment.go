package managers

import (
	"path/filepath"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func (m *Manager) CreateDeployment(cr types.NamespacedName, deviceName string, includeCM bool) *appsv1.Deployment {
	labels := map[string]string{
		"app": m.GetName(cr.Name, deviceName),
	}
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.GetName(cr.Name, deviceName),
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:      m.GetName(cr.Name, deviceName),
					Namespace: cr.Namespace,
					Labels:    labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:    m.GetName(cr.Name, deviceName),
							Image:   m.Image,
							Command: []string{"/bin/sh"},
							Args: []string{
								"-c",
								"socat -d -d pty,raw,echo=0,b115200,link=/dev/device,perm=0660,group=tty tcp:" + strings.ToLower(deviceName+"-gateway") + ":3333 & " + m.RunCmnd, // TODO: make init container
							},
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									Protocol:      corev1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}
	if includeCM {
		deployment.Spec.Template.Spec.Volumes = append(deployment.Spec.Template.Spec.Volumes, corev1.Volume{
			Name: "config",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: m.GetName(cr.Name, deviceName),
					},
					Items: []corev1.KeyToPath{
						{
							Key:  filepath.Base(m.ConfigPath),
							Path: filepath.Base(m.ConfigPath),
						},
					},
				},
			},
		})

		container := deployment.Spec.Template.Spec.Containers[0]
		container.VolumeMounts = []corev1.VolumeMount{
			{
				Name:      "config",
				ReadOnly:  false,
				MountPath: m.ConfigPath,
				SubPath:   filepath.Base(m.ConfigPath),
			},
		}

		deployment.Spec.Template.Spec.Containers = []corev1.Container{container}
	}
	return deployment
}
