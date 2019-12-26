package monitor

import (
	appv1alpha1 "github.com/janekbaraniewski/kubeserial/pkg/apis/app/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateConfigMap(cr *appv1alpha1.KubeSerial) *corev1.ConfigMap {
	rule := ""
	labels := map[string]string{
		"app": cr.Name + "-monitor",
	}

	for _, device := range cr.Spec.Devices {
		rule += "SUBSYSTEM==\"" + device.Subsystem + "\", ATTRS{idVendor}==\"" + device.IdVendor + "\", ATTRS{idProduct}==\"" + device.IdProduct + "\", SYMLINK+=\"" + device.Name + "\"\n"
	}

	return &corev1.ConfigMap {
		ObjectMeta:		metav1.ObjectMeta {
			Name: 		cr.Name + "-monitor",
			Namespace:	cr.Namespace,
			Labels:		labels,
		},
		Data:			map[string]string {
			"98-devices.rules": rule,
		},
	}
}
