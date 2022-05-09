package gateway

import (
	"strings"

	appv1alpha1 "github.com/janekbaraniewski/kubeserial/pkg/apis/kubeserial/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func CreateService(device *appv1alpha1.Device) *corev1.Service {
	labels := map[string]string{
		"app": device.Name + "-gateway",
	}
	name := strings.ToLower(device.Name + "-gateway")
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: device.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       "ser2net",
					Protocol:   corev1.ProtocolTCP,
					Port:       3333,
					TargetPort: intstr.FromInt(3333),
				},
			},
			Selector: map[string]string{
				"app": name,
			},
			Type: corev1.ServiceTypeClusterIP,
		},
	}

}
