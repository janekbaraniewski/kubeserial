package managers

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (m *Manager) CreateService(cr types.NamespacedName, deviceName string) *corev1.Service {
	labels := map[string]string{
		"app": m.GetName(cr.Name, deviceName),
	}
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.GetName(cr.Name, deviceName),
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       "http",
					Protocol:   corev1.ProtocolTCP,
					Port:       80,
					TargetPort: intstr.FromInt(80),
				},
			},
			Selector: map[string]string{
				"app": m.GetName(cr.Name, deviceName),
			},
			Type: corev1.ServiceTypeClusterIP,
		},
	}

}
