package monitor

import (
	"fmt"

	appv1alpha1 "github.com/janekbaraniewski/kubeserial/pkg/apis/kubeserial/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateConfigMap(cr *appv1alpha1.KubeSerial) *corev1.ConfigMap {
	rule := ""
	labels := map[string]string{
		"app": cr.Name + "-monitor",
	}

	for _, device := range cr.Spec.SerialDevices {
		rule += fmt.Sprintf(
			"SUBSYSTEM==\"tty\", ATTRS{idVendor}==\"%s\", ATTRS{idProduct}==\"%s\", SYMLINK+=\"%s\"\n",
			device.IdVendor,
			device.IdProduct,
			device.Name,
		)
	}

	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-monitor",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Data: map[string]string{
			"98-devices.rules": rule,
		},
	}
}
